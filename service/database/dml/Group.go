package dml

import (
	"fmt"
)

// CreateGroup inserisce un nuovo gruppo e restituisce il groupID generato.
// Inserisce inoltre l'admin come primo membro e registra l'evento "entered".
func CreateGroup(name string, photo string, admin int) (int, error) {
	// Inserisci il gruppo (photo = nil se vuota, per rispettare il CHECK constraint)
	var photoParam interface{}
	if photo != "" {
		photoParam = photo
	}
	result, err := db.Exec("INSERT INTO Groups (name, photo, adminID) VALUES (?, ?, ?)", name, photoParam, admin)
	if err != nil {
		return 0, fmt.Errorf("unable to insert new group %s: %w", name, err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("unable to get last insert id for group %s: %w", name, err)
	}

	groupID := int(lastID)

	// Aggiunge l'admin come primo membro del gruppo
	_, err = db.Exec("INSERT INTO Members (groupID, userID) VALUES (?, ?)", groupID, admin)
	if err != nil {
		return 0, fmt.Errorf("unable to add admin %d as member of group %d: %w", admin, groupID, err)
	}

	// Registra l'evento di ingresso dell'admin (timestamp sub-secondo per evitare
	// collisioni sulla PK (groupID, dt_event) se più utenti entrano nello stesso secondo)
	_, err = db.Exec("INSERT INTO Events (groupID, type, user, dt_event) VALUES (?, 'entered', ?, strftime('%Y-%m-%d %H:%M:%f', 'now'))", groupID, admin)
	if err != nil {
		return 0, fmt.Errorf("unable to register entered event for admin %d in group %d: %w", admin, groupID, err)
	}

	return groupID, nil
}

// LeaveGroup rimuove un utente dal gruppo e registra l'evento "leave".
func LeaveGroup(userID int, groupID int) error {
	// Rimuovi l'utente dalla tabella Members
	result, err := db.Exec("DELETE FROM Members WHERE groupID = ? AND userID = ?", groupID, userID)
	if err != nil {
		return fmt.Errorf("unable to remove member row for groupID %d and userID %d: %w", groupID, userID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for leaving group %d: %w", groupID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user %d is not a member of group %d", userID, groupID)
	}

	// Cancella la conversazione (e i messaggi via CASCADE)
	_, err = db.Exec("DELETE FROM ConversationsGroup WHERE userID = ? AND groupID = ?", userID, groupID)
	if err != nil {
		return fmt.Errorf("unable to delete conversation for userID %d in group %d: %w", userID, groupID, err)
	}

	// Registra l'evento di uscita (timestamp sub-secondo)
	_, err = db.Exec("INSERT INTO Events (groupID, type, user, dt_event) VALUES (?, 'leave', ?, strftime('%Y-%m-%d %H:%M:%f', 'now'))", groupID, userID)
	if err != nil {
		return fmt.Errorf("unable to register leave event for userID %d in group %d: %w", userID, groupID, err)
	}

	return nil
}

// AddToGroup aggiunge una lista di utenti al gruppo e registra gli eventi "entered".
func AddToGroup(adminID int, userID []int, groupID int) error {
	for _, uid := range userID {
		// Inserisci il membro
		_, err := db.Exec("INSERT INTO Members (groupID, userID) VALUES (?, ?)", groupID, uid)
		if err != nil {
			return fmt.Errorf("unable to add userID %d to group %d: %w", uid, groupID, err)
		}

		// Registra l'evento di ingresso (timestamp sub-secondo per evitare
		// collisioni sulla PK (groupID, dt_event) se più utenti entrano nello stesso secondo)
		_, err = db.Exec("INSERT INTO Events (groupID, type, user, dt_event) VALUES (?, 'entered', ?, strftime('%Y-%m-%d %H:%M:%f', 'now'))", groupID, uid)
		if err != nil {
			return fmt.Errorf("unable to register entered event for userID %d in group %d: %w", uid, groupID, err)
		}
	}

	return nil
}

// UpdateNameGroup aggiorna il nome del gruppo, verificando che l'utente sia l'admin.
func UpdateNameGroup(groupID int, adminID int, new_name string) error {
	result, err := db.Exec("UPDATE Groups SET name = ? WHERE groupID = ? AND adminID = ?", new_name, groupID, adminID)
	if err != nil {
		return fmt.Errorf("unable to update the name of groupID %d: %w", groupID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for group %d: %w", groupID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group %d not found or user %d is not the admin", groupID, adminID)
	}

	return nil
}

// UpdatePhotoGroup aggiorna la foto del gruppo, verificando che l'utente sia l'admin.
func UpdatePhotoGroup(groupID int, adminID int, new_photo string) error {
	result, err := db.Exec("UPDATE Groups SET photo = ? WHERE groupID = ? AND adminID = ?", new_photo, groupID, adminID)
	if err != nil {
		return fmt.Errorf("unable to update the photo of groupID %d: %w", groupID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for group %d: %w", groupID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group %d not found or user %d is not the admin", groupID, adminID)
	}

	return nil
}

// DeleteGroup elimina il gruppo. Grazie al CASCADE vengono rimossi anche Members, Events, ecc.
// Verifica che l'utente sia l'admin del gruppo.
func DeleteGroup(groupID int, userID int) error {
	result, err := db.Exec("DELETE FROM Groups WHERE groupID = ? AND adminID = ?", groupID, userID)
	if err != nil {
		return fmt.Errorf("unable to remove group %d from database: %w", groupID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for group %d: %w", groupID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group %d not found or user %d is not the admin", groupID, userID)
	}

	return nil
}

// KickFromGroup rimuove un utente dal gruppo e registra l'evento "kick".
func KickFromGroup(groupID int, adminID int, userID int) error {
	// Rimuovi l'utente dalla tabella Members
	result, err := db.Exec("DELETE FROM Members WHERE groupID = ? AND userID = ?", groupID, userID)
	if err != nil {
		return fmt.Errorf("unable to kick userID %d from group %d: %w", userID, groupID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for kick in group %d: %w", groupID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user %d is not a member of group %d", userID, groupID)
	}

	// Cancella la conversazione (e i messaggi via CASCADE)
	_, err = db.Exec("DELETE FROM ConversationsGroup WHERE userID = ? AND groupID = ?", userID, groupID)
	if err != nil {
		return fmt.Errorf("unable to delete conversation for userID %d in group %d: %w", userID, groupID, err)
	}

	// Registra l'evento di kick (timestamp sub-secondo)
	_, err = db.Exec("INSERT INTO Events (groupID, type, user, dt_event) VALUES (?, 'kick', ?, strftime('%Y-%m-%d %H:%M:%f', 'now'))", groupID, userID)
	if err != nil {
		return fmt.Errorf("unable to register kick event for userID %d in group %d: %w", userID, groupID, err)
	}

	return nil
}
