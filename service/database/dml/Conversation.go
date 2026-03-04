package dml

import (
	"fmt"
)

// CreateConversationBetweenUsers crea una nuova conversazione tra due utenti (se non esiste già)
// e restituisce il conversationID. I due ID vengono normalizzati (min, max) per garantire unicità
// indipendentemente dall'ordine in cui sender e receiver vengono passati.
func CreateConversationBetweenUsers(sender int, receiver int) (int, error) {
	// Normalizza: user1ID è sempre il minore, così (A,B) e (B,A) puntano alla stessa riga
	user1, user2 := sender, receiver
	if user1 > user2 {
		user1, user2 = user2, user1
	}

	result, err := db.Exec(
		"INSERT OR IGNORE INTO ConversationsUser (user1ID, user2ID) VALUES (?, ?)",
		user1, user2,
	)
	if err != nil {
		return 0, fmt.Errorf("unable to insert conversation between users %d and %d: %w", sender, receiver, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("unable to get rows affected for conversation between users %d and %d: %w", sender, receiver, err)
	}

	if rowsAffected > 0 {
		lastID, err := result.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("unable to get last insert id for conversation: %w", err)
		}
		return int(lastID), nil
	}

	// La conversazione esiste già: la recupera
	var conversationID int
	err = db.QueryRow(
		"SELECT conversationID FROM ConversationsUser WHERE user1ID = ? AND user2ID = ?",
		user1, user2,
	).Scan(&conversationID)
	if err != nil {
		return 0, fmt.Errorf("unable to retrieve existing conversation between users %d and %d: %w", sender, receiver, err)
	}

	return conversationID, nil
}

// CreateConversationBetweenGroups crea una nuova conversazione tra un utente e un gruppo (se non esiste già)
// e restituisce il conversationID.
func CreateConversationBetweenGroups(userID int, groupID int) (int, error) {
	result, err := db.Exec(
		"INSERT OR IGNORE INTO ConversationsGroup (userID, groupID) VALUES (?, ?)",
		userID, groupID,
	)
	if err != nil {
		return 0, fmt.Errorf("unable to insert conversation between user %d and group %d: %w", userID, groupID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("unable to get rows affected for conversation between user %d and group %d: %w", userID, groupID, err)
	}

	if rowsAffected > 0 {
		lastID, err := result.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("unable to get last insert id for group conversation: %w", err)
		}
		return int(lastID), nil
	}

	// La conversazione esiste già: la recupera
	var conversationID int
	err = db.QueryRow(
		"SELECT conversationID FROM ConversationsGroup WHERE userID = ? AND groupID = ?",
		userID, groupID,
	).Scan(&conversationID)
	if err != nil {
		return 0, fmt.Errorf("unable to retrieve existing conversation between user %d and group %d: %w", userID, groupID, err)
	}

	return conversationID, nil
}

// UpdateLastMessageIDforUser aggiorna il campo lastMessageID di una conversazione tra utenti
func UpdateLastMessageIDforUser(conversationID int, messageID int) error {
	result, err := db.Exec(
		"UPDATE ConversationsUser SET lastMessageID = ? WHERE conversationID = ?",
		messageID, conversationID,
	)
	if err != nil {
		return fmt.Errorf("unable to update lastMessageID %d for conversationID %d in ConversationsUser: %w", messageID, conversationID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for conversationID %d: %w", conversationID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("conversation %d not found in ConversationsUser", conversationID)
	}

	return nil
}

// UpdateLastMessageIDforGroup aggiorna il campo lastMessageID di una conversazione utente-gruppo
func UpdateLastMessageIDforGroup(conversationID int, messageID int) error {
	result, err := db.Exec(
		"UPDATE ConversationsGroup SET lastMessageID = ? WHERE conversationID = ?",
		messageID, conversationID,
	)
	if err != nil {
		return fmt.Errorf("unable to update lastMessageID %d for conversationID %d in ConversationsGroup: %w", messageID, conversationID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for conversationID %d: %w", conversationID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("conversation %d not found in ConversationsGroup", conversationID)
	}

	return nil
}
