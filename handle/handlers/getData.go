package handlers

import (
	"backend-pajak/database"
	"backend-pajak/services"
	"encoding/json"
	"net/http"
	"time"
)

func GetData(w http.ResponseWriter, r *http.Request) {
	db := database.InitializeDB()
	defer db.Close()

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	tanggal := now.Format("2006-01-02")

	dataList, err := services.GetData(db, tanggal)
	if err != nil {
		http.Error(w, "Failed to get data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": dataList})
}
