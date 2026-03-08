package app

import (
	"fmt"
	"os"

	"github.com/dp487/legendary-succotash/database"
	"github.com/dp487/legendary-succotash/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

var SecretKey string // Declare the secret key variable

func SetupAndRunApp() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("loading environment variables failed: %w", err)
	}

	// Load and check for the JWT secret key
	SecretKey = os.Getenv("TOKEN_SECRET")
	if SecretKey == "" {
		return fmt.Errorf("environment variable TOKEN_SECRET is required")
	}

	// Validate required environment variables
	requiredVars := []string{"APP_HOST", "APP_PORT", "POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DBNAME"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return fmt.Errorf("environment variable %s is required", v)
		}
	}

	// Check Connection to the database
	db, err := database.Connect() // Get the DB connection
	if err != nil {
		return err
	}

	// Create app
	app := fiber.New()

	app.Use(cors.New())

	// Attach middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))

	// Setup routes
	router.SetupRoutes(app, db) // Pass the db instance

	// Get the port and start
	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")

	return app.Listen(fmt.Sprintf("%s:%s", host, port))
}
