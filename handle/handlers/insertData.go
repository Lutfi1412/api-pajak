package handlers

import (
	"backend-pajak/database"
	"backend-pajak/services"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
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

func InsertData(w http.ResponseWriter, r *http.Request) {
	db := database.InitializeDB()
	defer db.Close()

	var req DataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	err := services.InsertData(db, req.Nama, req.Kelamin, req.Alamat, req.Jenis, req.Plat, req.TenggatThn, req.TenggatBln)
	if err != nil {
		http.Error(w, "Failed to insert data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Sukses"})
}
