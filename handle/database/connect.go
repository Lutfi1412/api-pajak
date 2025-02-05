package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitializeDB() *sql.DB {
	dbUser := os.Getenv("DB_USER")         // Your database user
	dbPassword := os.Getenv("DB_PASSWORD") // Your database password
	dbHost := os.Getenv("DB_HOST")         // Your database host
	dbPort := os.Getenv("DB_PORT")         // Your database port
	dbName := os.Getenv("DB_NAME")         // Your database name

	// If any of these environment variables are missing, log an error
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("One or more environment variables are missing. Please check your environment setup.")
	}

	// Create connection string using the environment variables
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=require", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open connection to PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
	} else {
		log.Println("Successfully connected to the database!")
	}

	// Return the database connection object
	return db
}
