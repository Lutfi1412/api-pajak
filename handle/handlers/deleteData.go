package handlers

import (
	"backend-pajak/database"
	"backend-pajak/services"
	"encoding/json"
	"net/http"

	. "github.com/tbxark/g4vercel"
)

type DeleteRequest struct {
	ID []int `json:"id"` // Masih array, tapi properti JSON bernama "id"
}

func DeleteData(c *Context) {

	db := database.InitializeDB()
	var req DeleteRequest
	// Decode request JSON
	if err := json.NewDecoder(c.Req.Body).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, H{
			"error": "Invalid JSON format",
		})
		return
	}

	// Validasi apakah array kosong
	if len(req.ID) == 0 {
		c.JSON(http.StatusBadRequest, H{
			"error": "No ID provided",
		})
		return
	}

	// Panggil service
	err := services.DeleteData(db, req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, H{
			"error": "Failed to delete data",
		})
		return
	}

	c.JSON(http.StatusOK, H{
		"message": "Data deleted successfully",
	})
}
