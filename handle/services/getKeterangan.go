package services

import (
	"database/sql"
	"log"
	"strconv"
)

func GetKeterangan(db *sql.DB, plat string, tanggal string, jenis string) (string, error) {
	queryPemilik := `SELECT tenggat_thn, tenggat_bln FROM data_pemilik WHERE plat = $1 AND jenis = $2`
	var tenggatThn, tenggatBln int

	err := db.QueryRow(queryPemilik, plat, jenis).Scan(&tenggatThn, &tenggatBln)
	if err == sql.ErrNoRows {
		return "Plat nomor tidak terdaftar", nil // Jangan return error supaya tidak gagal saat insert
	} else if err != nil {
		log.Printf("Error querying data_pemilik: %v", err)
		return "Terjadi kesalahan saat validasi data", err
	}

	// Konversi tanggal
	tahunStr := tanggal[:4] // Ambil yyyy
	tahun, _ := strconv.Atoi(tahunStr[2:])

	bulanStr := tanggal[5:7] // Ambil mm
	bulan, _ := strconv.Atoi(bulanStr)

	// Cek status pajak
	if tahun > tenggatThn {
		return "Belum bayar pajak", nil
	} else if tahun == tenggatThn && bulan > tenggatBln {
		return "Belum bayar pajak", nil
	} else {
		return "Sudah bayar pajak", nil
	}
}
