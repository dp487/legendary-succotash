package database

import (
	"fmt"
	"log"
	"os"

	"github.com/dp487/legendary-succotash/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database holds the database connection
type Database struct {
	DB *gorm.DB
}

// Connect creates a new database connection and returns a Database instance
func Connect() (*Database, error) {
	dsn := getDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("connection to database failed: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.UserSessions{}); err != nil {
		return nil, fmt.Errorf("AutoMigrate failed: %w", err)
	}

	return &Database{DB: db}, nil
}

// getDSN constructs the Data Source Name for connecting to the database
func getDSN() string {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Failed to load ENV")
	}
	POSTGRES_HOST := os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT := os.Getenv("POSTGRES_PORT")
	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DBNAME := os.Getenv("POSTGRES_DBNAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DBNAME)
}
