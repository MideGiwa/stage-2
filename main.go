package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"stage-2/config"
	"stage-2/controllers"
	"stage-2/models"
	"stage-2/services"
	"stage-2/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	log.Println("Attempting to load environment variables...")
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables from host environment.")
	} else {
		log.Println("Environment variables loaded from .env file.")
	}

	// Initialize database
	log.Println("Attempting to connect to the database...")
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection successful.")

	// Auto-migrate models
	log.Println("Attempting to auto-migrate database models...")
	err = db.AutoMigrate(&models.Country{}, &models.Status{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	log.Println("Database auto-migration successful.")

	// Initialize services
	countryService := services.NewCountryService(db)
	statusService := services.NewStatusService(db)

	// Initialize controllers
	countryController := controllers.NewCountryController(countryService, statusService)
	statusController := controllers.NewStatusController(statusService)

	// Set up Gin router
	router := gin.Default()

	// Middleware to handle external API errors
	router.Use(utils.ExternalAPIErrorMiddleware())

	// Routes
	router.POST("/countries/refresh", countryController.RefreshCountries)
	router.GET("/countries", countryController.GetCountries)
	router.GET("/countries/:name", countryController.GetCountryByName)
	router.DELETE("/countries/:name", countryController.DeleteCountry)
	router.GET("/status", statusController.GetStatus)
	router.GET("/countries/image", countryController.ServeSummaryImage)

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Country Data API!"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Starting server on port %s...", port)
	server := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
	log.Println("Server stopped.")
}
