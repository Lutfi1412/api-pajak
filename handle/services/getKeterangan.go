package services

import (
	"database/sql"
	"log"
	"strconv"
)

func GetKeterangan(db *sql.DB, plat string, tanggal string, jenis string) string {
	queryPemilik := `SELECT tenggat_thn, tenggat_bln FROM data_pemilik WHERE plat = $1 AND jenis = $2`
	var tenggatThn, tenggatBln int

	err := db.QueryRow(queryPemilik, plat, jenis).Scan(&tenggatThn, &tenggatBln)
	if err == sql.ErrNoRows {
		return "Plat nomor tidak terdaftar"
	} else if err != nil {
		log.Printf("Error querying data_pemilik: %v", err)
		return "Terjadi kesalahan saat validasi data"
	}

	// Konversi tanggal
	tahunStr := tanggal[:4] // Ambil yyyy
	tahun, _ := strconv.Atoi(tahunStr[2:])

	bulanStr := tanggal[5:7] // Ambil mm
	bulan, _ := strconv.Atoi(bulanStr)

	// Cek status pajak
	if tahun > tenggatThn {
		return "Belum bayar pajak"
	} else if tahun == tenggatThn && bulan > tenggatBln {
		return "Belum bayar pajak"
	} else {
		return "Sudah bayar pajak"
	}
}
