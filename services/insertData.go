package services

import (
	"database/sql"
	"log"
)

func InsertData(db *sql.DB, nama string, kelamin string, alamat string, jenis string, plat string, tenggatThn int, tenggatBln int) error {
	log.Printf("Inserting data: %v, %v, %v, %v, %v, %v, %v", nama, kelamin, alamat, jenis, plat, tenggatThn, tenggatBln)
	queryData := `INSERT INTO data_pemilik (nama, kelamin, alamat, jenis, plat, tenggat_thn, tenggat_bln) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(queryData, nama, kelamin, alamat, jenis, plat, tenggatThn, tenggatBln)
	if err != nil {
		log.Printf("Error inserting data into laporan: %v", err)
		return err
	}
	log.Println("Data pemilik berhasil dimasukkan!")
	return nil
}
