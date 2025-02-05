package handlers

import (
	"backend-pajak/database"
	"backend-pajak/services"
	"encoding/json"
	"log"
	"net/http"

	. "github.com/tbxark/g4vercel"
)

type UpdateJenisRequest struct {
	ID         int `json:"id"` // ID yang ingin diubah
	TenggatThn int `json:"tenggat_thn"`
	TenggatBln int `json:"tenggat_bln"` // Nilai baru untuk kolom jenis
}

func UpdateData(c *Context) {
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

	var req UpdateJenisRequest
	if err := json.NewDecoder(c.Req.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		c.JSON(http.StatusBadRequest, H{
			"error": "Invalid JSON format",
		})
		return
	}

	// Panggil service untuk update data
	err := services.UpdateData(db, req.ID, req.TenggatThn, req.TenggatBln)
	if err != nil {
		log.Printf("Error updating jenis: %v", err)
		c.JSON(http.StatusInternalServerError, H{
			"error": "Failed to update jenis",
		})
		return
	}

	c.JSON(http.StatusOK, H{
		"message": "Jenis updated successfully",
	})
}
