package handlers

import (
	"backend-pajak/database"
	"backend-pajak/services"
	"log"
	"net/http"

	. "github.com/tbxark/g4vercel"
)

func GetLaporan(c *Context) {
	// Inisialisasi database
	db := database.InitializeDB()
	if db == nil {
		log.Println("Failed to initialize database")
		c.JSON(http.StatusInternalServerError, H{
			"error": "Database connection failed",
		})
		return
	}
	defer db.Close()

	dataList, err := services.GetLaporan(db)
	if err != nil {
		log.Printf("Error getting data: %v", err)
		c.JSON(http.StatusInternalServerError, H{
			"error": "Failed to get data",
		})
		return
	}

	// Respons dengan data
	c.JSON(200, H{
		"data": dataList,
	})
}
