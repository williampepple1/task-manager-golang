package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// InitDB initializes and returns a connection to the database
func InitDB() (*gorm.DB, error) {
	// The connection string should be retrieved from environment variables or a secure config for production apps!
	const connStr = "host=localhost user=postgres dbname=taskmanager sslmode=disable password=postgres"

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
