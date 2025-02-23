package handlers

import (
	"backend-pajak/database"
	"backend-pajak/services"
	"encoding/json"
	"net/http"
)

type DeleteRequest struct {
	ID []int `json:"id"`
}

func DeleteData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	db := database.InitializeDB()
	var req DeleteRequest

	// Decode request JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validasi apakah array kosong
	if len(req.ID) == 0 {
		http.Error(w, "No ID provided", http.StatusBadRequest)
		return
	}

	// Panggil service
	err := services.DeleteData(db, req.ID)
	if err != nil {
		http.Error(w, "Failed to delete data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Data deleted successfully"})
}
