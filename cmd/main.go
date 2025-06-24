package main

import (
	"boiler-platecode/src/apis"
	"boiler-platecode/src/apis/auth"
	"boiler-platecode/src/apis/user"
	"boiler-platecode/src/core/config"
	"boiler-platecode/src/core/database"
	"boiler-platecode/src/common/validator"
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
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	// Get port
	port := config.AppConfig.PORT
	addr := ":" + port

	// Initialize DB
	database.InitDB()
	db := database.GetDB()
	validator.RegisterCustomValidations()

	// Create services and controllers
	userController := user.InitUserController(db)
	authController := auth.InitAuthController(db)
	apiController := apis.NewApiController(userController, authController)
	// Setup Gin router
	r := gin.Default()

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Add routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, world!"})
	})

	// Register all routes
	apiController.RegisterRoutes(r)
	// Create HTTP server
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Run server in goroutine
	go func() {
		log.Printf("Starting server on http://localhost%s (mode: %s)...", addr, mode)
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
