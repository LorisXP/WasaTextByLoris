package dml

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
)

// CreateTextContent inserisce un nuovo contenuto nella tabella Contents
// e restituisce il contentID generato
func CreateTextContent(message_content entity.Content) (int, error) {
	result, err := db.Exec(
		"INSERT INTO Contents (type, content) VALUES (?, ?)",
		message_content.Type, message_content.Content,
	)
	if err != nil {
		return 0, fmt.Errorf("unable to insert content in the database: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("unable to get contentID: %w", err)
	}

	return int(lastID), nil
}

// DeleteContentByID cancella un contenuto dalla tabella Contents dato il contentID
func DeleteContentByID(contentID int) error {
	result, err := db.Exec("DELETE FROM Contents WHERE contentID = ?", contentID)
	if err != nil {
		return fmt.Errorf("unable to remove contentID %d from database: %w", contentID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for contentID %d: %w", contentID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("contentID %d not found", contentID)
	}

	return nil
}