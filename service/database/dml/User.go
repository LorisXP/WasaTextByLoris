package dml

import "fmt"

// UpdatePhotoByUserID aggiorna la foto di un utente dato il suo userID
func UpdatePhotoByUserID(userID int, photo string) error {
	result, err := db.Exec("UPDATE Users SET photo = ? WHERE userID = ?", photo, userID)
	if err != nil {
		return fmt.Errorf("unable to update photo for user %d: %w", userID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for user %d: %w", userID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user %d not found", userID)
	}

	return nil
}

// UpdateNameByUserID aggiorna il nome di un utente dato il suo userID
func UpdateNameByUserID(userID int, userName string) error {
	result, err := db.Exec("UPDATE Users SET name = ? WHERE userID = ?", userName, userID)
	if err != nil {
		return fmt.Errorf("unable to update user name for user %d: %w", userID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for user %d: %w", userID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user %d not found", userID)
	}

	return nil
}
