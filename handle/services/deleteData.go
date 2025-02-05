package services

import (
	"database/sql"
	"fmt"
)

func DeleteData(db *sql.DB, ids []int) error {
	// Query DELETE tetap sama
	placeholder := ""
	for i := range ids {
		if i > 0 {
			placeholder += ", "
		}
		placeholder += fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(`
		DELETE FROM data_pemilik
		WHERE id IN (%s)
	`, placeholder)

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err := db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
