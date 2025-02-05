package handlers

import (
	"backend-pajak/database"
	"backend-pajak/services"
	"encoding/json"
	"log"
	"net/http"
	"time"

	. "github.com/tbxark/g4vercel"
)

type DataById struct {
	ID int `json:"id"`
}

func GetDataById(c *Context) {
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

	var req DataById

	if err := json.NewDecoder(c.Req.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		c.JSON(http.StatusBadRequest, H{
			"error": "Invalid JSON format",
		})
		return
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	tanggal := now.Format("2006-01-02")

	dataList, err := services.GetDataById(db, req.ID, tanggal)
	if err != nil {
		log.Printf("Error getting data: %v", err)
		c.JSON(http.StatusInternalServerError, H{
			"error": "Failed to get data",
		})
		return
	}

	if dataList == nil {
		c.JSON(http.StatusNotFound, H{
			"error": "Data not found",
		})
		return
	}

	c.JSON(http.StatusOK, H{
		"status": "success",
		"data":   dataList,
	})
}
