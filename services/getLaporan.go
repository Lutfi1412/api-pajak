package services

import (
	"database/sql"
	"log"
	"time"
)

type DataLaporan struct {
	ID         int    `json:"id"`
	Tanggal    string `json:"tanggal"`
	Jam        string `json:"jam"`
	Jenis      string `json:"jenis"`
	Plat       string `json:"plat"`
	Keterangan string `json:"keterangan"`
}

// GetLaporan mengambil data dari database dengan urutan yang ditentukan berdasarkan kolom dan arah
func GetLaporan(db *sql.DB) ([]DataLaporan, error) {

	queryLaporan := `SELECT * FROM laporan`
	rows, err := db.Query(queryLaporan)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	// List untuk menyimpan data
	var dataList []DataLaporan

	for rows.Next() {
		var data DataLaporan
		var tanggal time.Time
		var jam time.Time

		// Ambil data dari database
		err := rows.Scan(&data.ID, &tanggal, &jam, &data.Jenis, &data.Plat, &data.Keterangan)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		// Format tanggal dan jam ke format yang diinginkan
		data.Tanggal = tanggal.Format("2006-01-02") // Format yyyy-mm-dd
		data.Jam = jam.Format("15:04:05")           // Format hh:mm:ss

		// Tambahkan data ke list
		dataList = append(dataList, data)
	}

	log.Printf("Total rows fetched: %d", len(dataList)) // Log jumlah data
	return dataList, nil
}
