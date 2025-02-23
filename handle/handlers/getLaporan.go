package handlers

import (
	"backend-pajak/database"
	"backend-pajak/services"
	"encoding/json"
	"net/http"
)

func GetLaporan(w http.ResponseWriter, r *http.Request) {
	db := database.InitializeDB()
	defer db.Close()

	dataList, err := services.GetLaporan(db)
	if err != nil {
		http.Error(w, "Failed to get data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": dataList})
}
