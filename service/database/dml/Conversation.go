package dml

import (
	"fmt"
)

// CreateConversationBetweenUsers crea una nuova conversazione tra due utenti
// e restituisce il conversationID generato
func CreateConversationBetweenUsers(sender int, receiver int) (int, error) {
	result, err := db.Exec(
		"INSERT INTO ConversationsUser (user1ID, user2ID) VALUES (?, ?)",
		sender, receiver,
	)
	if err != nil {
		return 0, fmt.Errorf("unable to insert conversation between users %d and %d: %w", sender, receiver, err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("unable to get last insert id for conversation: %w", err)
	}

	return int(lastID), nil
}

// CreateConversationBetweenGroups crea una nuova conversazione tra un utente e un gruppo
// e restituisce il conversationID generato
func CreateConversationBetweenGroups(userID int, groupID int) (int, error) {
	result, err := db.Exec(
		"INSERT INTO ConversationsGroup (userID, groupID) VALUES (?, ?)",
		userID, groupID,
	)
	if err != nil {
		return 0, fmt.Errorf("unable to insert conversation between user %d and group %d: %w", userID, groupID, err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("unable to get last insert id for group conversation: %w", err)
	}

	return int(lastID), nil
}

// UpdateLastMessageIDforUser aggiorna il lastMessageID di una conversazione tra utenti
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

// UpdateLastMessageIDforGroup aggiorna il lastMessageID di una conversazione utente-gruppo
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
