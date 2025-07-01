package main

import (
	"boiler-platecode/src/apis"
	"boiler-platecode/src/common/lib/logger"
	"boiler-platecode/src/common/validator"
	"boiler-platecode/src/core/config"
	"boiler-platecode/src/core/database"
	"boiler-platecode/src/core/redis"

	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	module = "main"
)

func main() {
	// Load env variables
	config.LoadEnv()

	// Set Gin mode
	mode := config.AppConfig.GinMode
	if mode == "" {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	// Get port
	port := config.AppConfig.PORT
	addr := ":" + port

	// Initialize DB and Validator
	database.InitDB()
	db := database.GetDB()

	validator.RegisterCustomValidations()

	// Setup Gin
	r := gin.Default()
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		logger.Error(module, "SetTrustedProxies", err)
		os.Exit(1)
	}

	// Initialize Redis connection
	redis.Init()
	redisService := redis.GetRedisService()

	// Initialize API controller
	apiController := apis.InitApiController(db, &redisService)

	// Register versioned API routes
	apiController.RegisterRoutes(r)

	// Root Health Route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the API",
			"version": "v1",
		})
	})

	// Create and run HTTP server
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		logger.Info(module, "ServerStart", "Server is running on http://localhost"+addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(module, "ListenAndServe", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(module, "Shutdown", "Received shutdown signal. Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(module, "Shutdown", errors.New("Server forced to shutdown: "+err.Error()))
		os.Exit(1)
	}

	// Close Redis connection
	if err := redis.Close(); err != nil {
		logger.Error(module, "redisClose", errors.New("Redis shutdown  failed: "+err.Error()))
	} else {
		logger.Info(module, "Shutdown", "Redis connection closed.")
	}

	// Close Database connection
	if err := database.CloseDB(); err != nil {
		logger.Error(module, "databaseClose", errors.New("database shutdown  failed: "+err.Error()))
	} else {
		logger.Info(module, "Shutdown", "Database connection closed.")
	}

	logger.Info(module, "Shutdown", "Server stopped gracefully.")
}
