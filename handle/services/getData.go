package services

import (
	"database/sql"
	"log"
	"strconv"
)

type DataPemilik struct {
	ID         int    `json:"id"`
	Jenis      string `json:"jenis"`
	TenggatThn int    `json:"tenggat_thn"`
	TenggatBln int    `json:"tenggat_bln"`
	Keterangan string `json:"keterangan"`
}

// GetData mengambil data dari database dengan pengurutan dinamis
func GetData(db *sql.DB, tanggal string) ([]DataPemilik, error) {
	// Default query untuk pengurutan berdasarkan ID
	queryPemilik := `SELECT id, jenis, tenggat_thn, tenggat_bln FROM data_pemilik`
	rows, err := db.Query(queryPemilik)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	// List untuk menyimpan data
	var dataList []DataPemilik

	tahunStr := tanggal[:4] // Ambil yyyy
	tahun, _ := strconv.Atoi(tahunStr[2:])

	bulanStr := tanggal[5:7] // Ambil mm
	bulan, _ := strconv.Atoi(bulanStr)

	for rows.Next() {
		var data DataPemilik
		err := rows.Scan(&data.ID, &data.Jenis, &data.TenggatThn, &data.TenggatBln)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		// Logika untuk menentukan keterangan
		if data.TenggatThn > tahun || (data.TenggatThn == tahun && data.TenggatBln >= bulan) {
			data.Keterangan = "sudah pajak"
		} else {
			data.Keterangan = "belum pajak"
		}

		// Tambahkan data ke list
		dataList = append(dataList, data)
	}

	log.Printf("Total rows fetched: %d", len(dataList)) // Log jumlah data
	return dataList, nil
}
