package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

// InitDB initializes and returns a connection to the database
func InitDB() (*gorm.DB, error) {
	host := os.Getenv("DATABASE_HOST")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	sslmode := "disable"
	port := os.Getenv("DATABASE_PORT")
	if port == "" {
		port = "5432" // default to 5432 if not specified
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", host, port, user, dbname, sslmode, password)

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
