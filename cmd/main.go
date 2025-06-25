package main

import (
	"boiler-platecode/src/apis"
	"boiler-platecode/src/common/validator"
	"boiler-platecode/src/core/config"
	"boiler-platecode/src/core/database"

	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
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
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}
	
	// Initialize API controller
	apiController := apis.InitApiController(db)

	// Register versioned API routes
	apiController.RegisterRoutes(r)

	// Root Health Route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "Welcome to the API",
			"version":     "v1",
		})
	})


	// Create and run HTTP server
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("[START] Server is running on http://localhost%s ", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Received shutdown signal. Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully.")
}
