package services

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
)

type DataById struct {
	ID         int    `json:"id"`
	Nama       string `json:"nama"`
	Kelamin    string `json:"kelamin"`
	Plat       string `json:"plat"`
	Alamat     string `json:"alamat"`
	Jenis      string `json:"jenis"`
	TenggatThn int    `json:"tenggat_thn"`
	TenggatBln int    `json:"tenggat_bln"`
	Keterangan string `json:"keterangan"`
}

// GetDataById mengambil data berdasarkan ID
func GetDataById(db *sql.DB, id int, tanggal string) (*DataById, error) {
	// Query untuk mengambil data berdasarkan ID
	query := `SELECT id, nama, kelamin, plat, alamat, jenis, tenggat_thn, tenggat_bln FROM data_pemilik WHERE id = $1`

	// Menggunakan QueryRow karena kita hanya mengambil satu baris data berdasarkan ID
	var data DataById

	// Eksekusi query untuk satu data berdasarkan ID
	err := db.QueryRow(query, id).Scan(&data.ID, &data.Nama, &data.Kelamin, &data.Plat, &data.Alamat, &data.Jenis, &data.TenggatThn, &data.TenggatBln)
	if err != nil {
		if err == sql.ErrNoRows {
			// Jika data tidak ditemukan
			log.Printf("No data found for ID: %d", id)
			return nil, nil // Tidak ditemukan
		}
		log.Printf("Error executing query: %v", err)
		return nil, err
	}

	log.Printf("tanggal: %v", tanggal)

	// Ambil tahun dan bulan dari tanggal yang diterima (seharusnya dalam format yyyy-mm-dd)
	tahunStr := tanggal[2:4] // Ambil tahun
	tahun, err := strconv.Atoi(tahunStr)
	if err != nil {
		log.Printf("Error parsing year: %v", err)
		return nil, err
	}

	// Ambil bulan dan pastikan formatnya benar, jika bulan memiliki angka di depan (misalnya "01", "02"), hilangkan angka 0 di depan
	bulanStr := tanggal[5:7] // Ambil bulan
	log.Printf("bulanStr: %v", bulanStr)
	bulan, err := strconv.Atoi(strings.TrimLeft(bulanStr, "0")) // Menghapus 0 di depan bulan
	if err != nil {
		log.Printf("Error parsing month: %v", err)
		return nil, err
	}

	// Perbandingan tahun dan bulan
	if data.TenggatThn > tahun || (data.TenggatThn == tahun && data.TenggatBln >= bulan) {
		data.Keterangan = "sudah pajak"
	} else {
		data.Keterangan = "belum pajak"
	}

	// Mengembalikan data yang telah terisi
	return &data, nil
}
