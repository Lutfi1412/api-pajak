package services

import (
	"database/sql"
	"log"
)

func UpdateData(db *sql.DB, id int, tenggatThn int, tenggatBln int) error {
	// Query untuk update kolom jenis
	query := `
		UPDATE data_pemilik 
		SET tenggat_thn = $1, tenggat_bln = $2
		WHERE id = $3
	`

	// Eksekusi query
	_, err := db.Exec(query, tenggatThn, tenggatBln, id)
	if err != nil {
		log.Printf("Error executing update query: %v", err)
		return err
	}

	log.Printf("Jenis updated successfully for ID: %d", id)
	return nil
}
