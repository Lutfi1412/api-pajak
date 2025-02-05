package handlers

import (
	"backend-pajak/database"
	"backend-pajak/services"
	"log"
	"net/http"
	"time"

	. "github.com/tbxark/g4vercel"
)

func GetData(c *Context) {
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

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	tanggal := now.Format("2006-01-02")

	// Panggil service untuk mengambil data dengan pengurutan
	dataList, err := services.GetData(db, tanggal)
	if err != nil {
		log.Printf("Error getting data: %v", err)
		c.JSON(http.StatusInternalServerError, H{
			"error": "Failed to get data",
		})
		return
	}

	c.JSON(200, H{
		"data": dataList,
	})
}
