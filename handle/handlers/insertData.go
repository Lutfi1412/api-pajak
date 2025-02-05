package handlers

import (
	"backend-pajak/database"
	"backend-pajak/services"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	. "github.com/tbxark/g4vercel"
)

type DataRequest struct {
	Nama       string `json:"nama"`
	Kelamin    string `json:"kelamin"`
	Alamat     string `json:"alamat"`
	Jenis      string `json:"jenis"`
	Plat       string `json:"plat"`
	TenggatThn int    `json:"tenggat_thn"`
	TenggatBln int    `json:"tenggat_bln"`
}

func InsertData(c *Context) {
	db := database.InitializeDB()
	defer db.Close()

	var req DataRequest

	// Membaca dan mendekodekan JSON dari request body
	if err := json.NewDecoder(c.Req.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		c.JSON(http.StatusBadRequest, H{
			"error": "Invalid JSON format",
		})
		return
	}

	// Memasukkan data ke dalam database
	err := services.InsertData(db, req.Nama, req.Kelamin, req.Alamat, req.Jenis, req.Plat, req.TenggatThn, req.TenggatBln)
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		c.JSON(http.StatusInternalServerError, H{
			"error": "Failed to insert data",
		})
		return
	}

	// Response sukses
	response := map[string]string{"status": "Sukses"}
	c.JSON(http.StatusOK, H{
		"data": response,
	})
}
