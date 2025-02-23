package handlers

import (
	"backend-pajak/database"
	"backend-pajak/services"
	"encoding/json"
	"net/http"
)

func DeleteLaporan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	db := database.InitializeDB()
	defer db.Close()

	var req DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if len(req.ID) == 0 {
		http.Error(w, "No ID provided", http.StatusBadRequest)
		return
	}

	err := services.DeleteLaporan(db, req.ID)
	if err != nil {
		http.Error(w, "Failed to delete data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Data deleted successfully"})
}
