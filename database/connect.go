package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitializeDB() *sql.DB {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables.")
	}

	// Ambil variabel environment
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=require",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Buka koneksi database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	// Cek apakah koneksi berhasil
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
	} else {
		log.Println("Successfully connected to the database!")
	}

	return db
}
