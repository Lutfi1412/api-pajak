package handlers

import (
	"backend-pajak/services"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// DetectHandler menangani request POST untuk deteksi
func DetectPlate(w http.ResponseWriter, r *http.Request) {
	log.Println("📩 Menerima request deteksi...")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No image file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		http.Error(w, "Invalid file format.", http.StatusBadRequest)
		return
	}

	tempFile := fmt.Sprintf("/tmp/temp_%d%s", time.Now().UnixNano(), ext)
	out, err := os.Create(tempFile)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	defer os.Remove(tempFile)

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to write image file", http.StatusInternalServerError)
		return
	}

	detections, err := services.DetectPlates(tempFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Detection error: %v", err), http.StatusInternalServerError)
		return
	}

	services.SaveToDatabase(detections)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   detections,
	})
}
