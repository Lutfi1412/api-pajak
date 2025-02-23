package services

import (
	"backend-pajak/database"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Struct hasil deteksi
type Detection struct {
	Jenis string `json:"jenis"`
	Plat  string `json:"plat"`
}

// DetectPlates mengirim gambar ke API dan mengembalikan hasil deteksi
func DetectPlates(imagePath string) ([]Detection, error) {
	apiURL := os.Getenv("API_URL")
	apiKey := os.Getenv("API_KEY")

	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("❌ Gagal membuka file: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("upload", filepath.Base(imagePath))
	if err != nil {
		return nil, fmt.Errorf("❌ Gagal membuat form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("❌ Gagal menyalin data file: %v", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return nil, fmt.Errorf("❌ Gagal membuat request: %v", err)
	}
	req.Header.Set("Authorization", "Token "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("❌ API request gagal: %v", err)
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	log.Printf("📄 API response body: %s", responseBody)

	var apiResponse struct {
		Results []struct {
			Vehicle struct {
				Type string `json:"type"`
			} `json:"vehicle"`
			Plate string `json:"plate"`
		} `json:"results"`
	}

	err = json.Unmarshal(responseBody, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("❌ Gagal decode response API: %v", err)
	}

	log.Printf("📊 Parsed API Response: %+v", apiResponse)

	var detections []Detection
	for _, result := range apiResponse.Results {
		jenis := mapVehicleType(result.Vehicle.Type)
		plat := result.Plate

		if jenis != "Tidak terdeteksi" && plat != "" {
			detections = append(detections, Detection{Jenis: jenis, Plat: plat})
		}
	}

	if len(detections) == 0 {
		log.Println("⚠️ Tidak ada kendaraan atau plat nomor yang terdeteksi.")
	}

	return detections, nil
}

// mapVehicleType menyederhanakan jenis kendaraan
func mapVehicleType(apiType string) string {
	apiType = strings.ToLower(apiType)
	switch apiType {
	case "sedan", "suv", "van", "hatchback", "pickup", "truck", "jeep", "wagon":
		return "mobil"
	case "motorcycle", "bike", "scooter":
		return "motor"
	default:
		return "Tidak terdeteksi"
	}
}

// SaveToDatabase menyimpan hasil deteksi ke database
func SaveToDatabase(detections []Detection) error {
	db := database.InitializeDB()
	tanggal := time.Now().Format("2006-01-02")
	jam := time.Now().Format("15:04:05")

	for _, deteksi := range detections {
		_, err := db.Exec(
			"INSERT INTO laporan (tanggal, jam, jenis, plat) VALUES ($1, $2, $3, $4)",
			tanggal, jam, deteksi.Jenis, deteksi.Plat,
		)
		if err != nil {
			return fmt.Errorf("❌ Gagal memasukkan data ke database: %v", err)
		}
	}
	return nil
}
