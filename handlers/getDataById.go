package handlers

import (
	"backend-pajak/database"
	"backend-pajak/services"
	"encoding/json"
	"net/http"
	"time"
)

type DataById struct {
	ID int `json:"id"`
}

func GetDataById(w http.ResponseWriter, r *http.Request) {
	db := database.InitializeDB()
	defer db.Close()

	var req DataById
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	tanggal := now.Format("2006-01-02")

	dataList, err := services.GetDataById(db, req.ID, tanggal)
	if err != nil {
		http.Error(w, "Failed to get data", http.StatusInternalServerError)
		return
	}

	if dataList == nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "data": dataList})
}
