package handlers

import (
	"backend-pajak/database"
	"backend-pajak/services"
	"encoding/json"
	"net/http"
)

type UpdateJenisRequest struct {
	ID         int `json:"id"` // ID yang ingin diubah
	TenggatThn int `json:"tenggat_thn"`
	TenggatBln int `json:"tenggat_bln"` // Nilai baru untuk kolom jenis
}

func UpdateData(w http.ResponseWriter, r *http.Request) {
	db := database.InitializeDB()
	defer db.Close()

	var req UpdateJenisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	err := services.UpdateData(db, req.ID, req.TenggatThn, req.TenggatBln)
	if err != nil {
		http.Error(w, "Failed to update data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Data updated successfully"})
}
