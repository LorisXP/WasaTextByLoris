package dml

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
)

/*
CreateMessageBelongUsers inserisce un nuovo messaggio in una conversazione tra utenti.
Il contenuto deve essere già stato creato tramite CreateTextContent.
Passi:
1. Determina il receiver dalla conversazione (l'altro utente)
2. Inserisce il messaggio nella tabella MessagesUser
3. Restituisce messageID, receiverID, sent_at
*/
func CreateMessageBelongUsers(message entity.Message, contentID int) (int, int, string, error) {
	// 1. Determina il receiver dalla conversazione
	var user1ID, user2ID int
	err := db.QueryRow(
		"SELECT user1ID, user2ID FROM ConversationsUser WHERE conversationID = ?",
		message.ConversationID,
	).Scan(&user1ID, &user2ID)
	if err != nil {
		return 0, 0, "", fmt.Errorf("unable to find conversation %d: %w", message.ConversationID, err)
	}

	receiverID := user1ID
	if message.Sender == user1ID {
		receiverID = user2ID
	}

	// 2. Inserisci il messaggio in MessagesUser
	msgResult, err := db.Exec(
		"INSERT INTO MessagesUser (conversationID, senderID, receiverID, status, type, contentID) VALUES (?, ?, ?, ?, ?, ?)",
		message.ConversationID, message.Sender, receiverID, "received", "standard", contentID,
	)
	if err != nil {
		return 0, 0, "", fmt.Errorf("unable to insert message belong users: %w", err)
	}

	messageID, err := msgResult.LastInsertId()
	if err != nil {
		return 0, 0, "", fmt.Errorf("unable to get messageID: %w", err)
	}

	// 3. Recupera il sent_at generato dal DEFAULT
	var sentAt string
	err = db.QueryRow(
		"SELECT sent_at FROM MessagesUser WHERE messageID = ?",
		int(messageID),
	).Scan(&sentAt)
	if err != nil {
		return 0, 0, "", fmt.Errorf("unable to get sent_at for messageID %d: %w", messageID, err)
	}

	// 4. Gestisci la tabella ConversationsUser:
	// Se la conversazione non esiste, la crea. Altrimenti aggiorna lastMessageID.
	var convCount int
	err = db.QueryRow(
		"SELECT COUNT(*) FROM ConversationsUser WHERE conversationID = ?",
		message.ConversationID,
	).Scan(&convCount)
	if err != nil {
		return 0, 0, "", fmt.Errorf("unable to check conversation existence %d: %w", message.ConversationID, err)
	}

	if convCount == 0 {
		// Recupera gli userID
		user1ID := message.Sender
		user2ID := receiverID
		// Crea la conversazione
		_, err := db.Exec(
			"INSERT INTO ConversationsUser (user1ID, user2ID, lastMessageID) VALUES (?, ?, ?)",
			user1ID, user2ID, int(messageID),
		)
		if err != nil {
			return 0, 0, "", fmt.Errorf("unable to create conversation for users %d-%d: %w", user1ID, user2ID, err)
		}
		// Aggiorna conversationID nel messaggio se serve (opzionale, dipende da logica)
		// message.ConversationID = int(res.LastInsertId())
	} else {
		// Aggiorna lastMessageID della conversazione esistente
		_, err = db.Exec(
			"UPDATE ConversationsUser SET lastMessageID = ? WHERE conversationID = ?",
			int(messageID), message.ConversationID,
		)
		if err != nil {
			return 0, 0, "", fmt.Errorf("unable to update lastMessageID for conversation %d: %w", message.ConversationID, err)
		}
	}

	return int(messageID), receiverID, sentAt, nil
}

/*
CreateMessageBetweenUsersAndGroups inserisce un nuovo messaggio in una conversazione utente-gruppo.
Il contenuto deve essere già stato creato tramite CreateTextContent.
Passi:
1. Determina il groupID (receiver) dalla conversazione
2. Inserisce il messaggio nella tabella MessagesGroup
3. Restituisce messageID, groupID (receiver), sent_at
*/
func CreateMessageBetweenUsersAndGroups(message entity.Message, contentID int) (int, int, string, error) {
	// 1. Determina il groupID dalla conversazione
	var groupID int
	err := db.QueryRow(
		"SELECT groupID FROM ConversationsGroup WHERE conversationID = ?",
		message.ConversationID,
	).Scan(&groupID)
	if err != nil {
		return 0, 0, "", fmt.Errorf("unable to find group conversation %d: %w", message.ConversationID, err)
	}

	// 2. Inserisci il messaggio in MessagesGroup
	msgResult, err := db.Exec(
		"INSERT INTO MessagesGroup (conversationID, senderID, status, type, contentID) VALUES (?, ?, ?, ?, ?)",
		message.ConversationID, message.Sender, "received", "standard", contentID,
	)
	if err != nil {
		return 0, 0, "", fmt.Errorf("unable to insert message between user and group: %w", err)
	}

	messageID, err := msgResult.LastInsertId()
	if err != nil {
		return 0, 0, "", fmt.Errorf("unable to get messageID: %w", err)
	}

	// 3. Recupera il sent_at generato dal DEFAULT
	var sentAt string
	err = db.QueryRow(
		"SELECT sent_at FROM MessagesGroup WHERE messageID = ?",
		int(messageID),
	).Scan(&sentAt)
	if err != nil {
		return 0, 0, "", fmt.Errorf("unable to get sent_at for messageID %d: %w", messageID, err)
	}

	// 4. Gestisci la tabella ConversationsGroup:
	// Se la conversazione non esiste, la crea. Altrimenti aggiorna lastMessageID.
	var convCount int
	err = db.QueryRow(
		"SELECT COUNT(*) FROM ConversationsGroup WHERE conversationID = ?",
		message.ConversationID,
	).Scan(&convCount)
	if err != nil {
		return 0, 0, "", fmt.Errorf("unable to check group conversation existence %d: %w", message.ConversationID, err)
	}

	if convCount == 0 {
		// Crea la conversazione
		_, err := db.Exec(
			"INSERT INTO ConversationsGroup (userID, groupID, lastMessageID) VALUES (?, ?, ?)",
			message.Sender, groupID, int(messageID),
		)
		if err != nil {
			return 0, 0, "", fmt.Errorf("unable to create group conversation for user %d and group %d: %w", message.Sender, groupID, err)
		}
		// Aggiorna conversationID nel messaggio se serve (opzionale)
		// message.ConversationID = int(res.LastInsertId())
	} else {
		// Aggiorna lastMessageID della conversazione esistente
		_, err = db.Exec(
			"UPDATE ConversationsGroup SET lastMessageID = ? WHERE conversationID = ?",
			int(messageID), message.ConversationID,
		)
		if err != nil {
			return 0, 0, "", fmt.Errorf("unable to update lastMessageID for group conversation %d: %w", message.ConversationID, err)
		}
	}

	return int(messageID), groupID, sentAt, nil
}

/*
DeleteMessageBelongUsers cancella un messaggio tra utenti.
Verifica che l'utente sia il sender del messaggio.
Cancella anche il contenuto associato tramite DeleteContentByID.
*/
func DeleteMessageBelongUsers(messageID int, userID int, conversationID int) error {
	// Recupera il contentID associato al messaggio prima di cancellarlo
	var contentID int
	err := db.QueryRow(
		"SELECT contentID FROM MessagesUser WHERE messageID = ? AND senderID = ? AND conversationID = ?",
		messageID, userID, conversationID,
	).Scan(&contentID)
	if err != nil {
		return fmt.Errorf("unable to find messageID %d for user %d in conversation %d: %w", messageID, userID, conversationID, err)
	}

	// Cancella il messaggio
	result, err := db.Exec(
		"DELETE FROM MessagesUser WHERE messageID = ? AND senderID = ? AND conversationID = ?",
		messageID, userID, conversationID,
	)
	if err != nil {
		return fmt.Errorf("unable to delete messageID %d: %w", messageID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for messageID %d: %w", messageID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("messageID %d not found or user %d is not the sender", messageID, userID)
	}

	// Cancella il contenuto associato tramite DeleteContentByID
	err = DeleteContentByID(contentID)
	if err != nil {
		return fmt.Errorf("unable to delete content for messageID %d: %w", messageID, err)
	}

	return nil
}

/*
DeleteMessageBetweenUsersAndGroups cancella un messaggio in una conversazione utente-gruppo.
Verifica che l'utente sia il sender del messaggio.
Cancella anche il contenuto associato tramite DeleteContentByID.
*/
func DeleteMessageBetweenUsersAndGroups(messageID int, userID int, conversationID int) error {
	// Recupera il contentID associato al messaggio prima di cancellarlo
	var contentID int
	err := db.QueryRow(
		"SELECT contentID FROM MessagesGroup WHERE messageID = ? AND senderID = ? AND conversationID = ?",
		messageID, userID, conversationID,
	).Scan(&contentID)
	if err != nil {
		return fmt.Errorf("unable to find messageID %d for user %d in group conversation %d: %w", messageID, userID, conversationID, err)
	}

	// Cancella il messaggio
	result, err := db.Exec(
		"DELETE FROM MessagesGroup WHERE messageID = ? AND senderID = ? AND conversationID = ?",
		messageID, userID, conversationID,
	)
	if err != nil {
		return fmt.Errorf("unable to delete messageID %d: %w", messageID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for messageID %d: %w", messageID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("messageID %d not found or user %d is not the sender", messageID, userID)
	}

	// Cancella il contenuto associato tramite DeleteContentByID
	err = DeleteContentByID(contentID)
	if err != nil {
		return fmt.Errorf("unable to delete content for messageID %d: %w", messageID, err)
	}

	return nil
}

/*
ForwardMessageBelongUsers inoltra un messaggio in una conversazione tra utenti.
Passi:
1. Recupera il contenuto del messaggio originale (message.MessageID)
2. Copia il contenuto tramite CreateTextContent
3. Determina il receiver dalla conversazione di destinazione
4. Inserisce il nuovo messaggio con type='forward'
5. Restituisce messageID, receiver, sent_at, content, contentType
*/
func ForwardMessageBelongUsers(message entity.Message) (int, int, string, string, string, error) {
	// 1. Recupera il contenuto del messaggio originale
	var content, contentType string

	// Cerca prima in MessagesUser, poi in MessagesGroup
	err := db.QueryRow(
		"SELECT c.content, c.type FROM MessagesUser mu JOIN Contents c ON mu.contentID = c.contentID WHERE mu.messageID = ?",
		message.MessageID,
	).Scan(&content, &contentType)
	if err != nil {
		// Prova in MessagesGroup
		err = db.QueryRow(
			"SELECT c.content, c.type FROM MessagesGroup mg JOIN Contents c ON mg.contentID = c.contentID WHERE mg.messageID = ?",
			message.MessageID,
		).Scan(&content, &contentType)
		if err != nil {
			return 0, 0, "", "", "", fmt.Errorf("unable to find original message %d to forward: %w", message.MessageID, err)
		}
	}

	// 2. Copia il contenuto tramite CreateTextContent
	newContent := entity.Content{Type: contentType, Content: content}
	newContentID, err := CreateTextContent(newContent)
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("unable to copy content for forwarded message: %w", err)
	}

	// 3. Determina il receiver dalla conversazione di destinazione
	var user1ID, user2ID int
	err = db.QueryRow(
		"SELECT user1ID, user2ID FROM ConversationsUser WHERE conversationID = ?",
		message.ConversationID,
	).Scan(&user1ID, &user2ID)
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("unable to find conversation %d: %w", message.ConversationID, err)
	}

	receiverID := user1ID
	if message.Sender == user1ID {
		receiverID = user2ID
	}

	// 4. Inserisci il messaggio con type='forward'
	msgResult, err := db.Exec(
		"INSERT INTO MessagesUser (conversationID, senderID, receiverID, status, type, contentID) VALUES (?, ?, ?, ?, ?, ?)",
		message.ConversationID, message.Sender, receiverID, "received", "forward", newContentID,
	)
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("unable to insert forwarded message: %w", err)
	}

	newMessageID, err := msgResult.LastInsertId()
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("unable to get new messageID: %w", err)
	}

	// 5. Recupera il sent_at
	var sentAt string
	err = db.QueryRow(
		"SELECT sent_at FROM MessagesUser WHERE messageID = ?",
		int(newMessageID),
	).Scan(&sentAt)
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("unable to get sent_at for messageID %d: %w", newMessageID, err)
	}

	return int(newMessageID), receiverID, sentAt, content, contentType, nil
}

/*
ForwardMessageBetweenUsersAndGroups inoltra un messaggio in una conversazione utente-gruppo.
Stessa logica di ForwardMessageBelongUsers ma per MessagesGroup.
*/
func ForwardMessageBetweenUsersAndGroups(message entity.Message) (int, int, string, string, string, error) {
	// 1. Recupera il contenuto del messaggio originale
	var content, contentType string

	// Cerca prima in MessagesUser, poi in MessagesGroup
	err := db.QueryRow(
		"SELECT c.content, c.type FROM MessagesUser mu JOIN Contents c ON mu.contentID = c.contentID WHERE mu.messageID = ?",
		message.MessageID,
	).Scan(&content, &contentType)
	if err != nil {
		err = db.QueryRow(
			"SELECT c.content, c.type FROM MessagesGroup mg JOIN Contents c ON mg.contentID = c.contentID WHERE mg.messageID = ?",
			message.MessageID,
		).Scan(&content, &contentType)
		if err != nil {
			return 0, 0, "", "", "", fmt.Errorf("unable to find original message %d to forward: %w", message.MessageID, err)
		}
	}

	// 2. Copia il contenuto tramite CreateTextContent
	newContent := entity.Content{Type: contentType, Content: content}
	newContentID, err := CreateTextContent(newContent)
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("unable to copy content for forwarded group message: %w", err)
	}

	// 3. Determina il groupID dalla conversazione
	var groupID int
	err = db.QueryRow(
		"SELECT groupID FROM ConversationsGroup WHERE conversationID = ?",
		message.ConversationID,
	).Scan(&groupID)
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("unable to find group conversation %d: %w", message.ConversationID, err)
	}

	// 4. Inserisci il messaggio con type='forward'
	msgResult, err := db.Exec(
		"INSERT INTO MessagesGroup (conversationID, senderID, status, type, contentID) VALUES (?, ?, ?, ?, ?)",
		message.ConversationID, message.Sender, "received", "forward", newContentID,
	)
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("unable to insert forwarded group message: %w", err)
	}

	newMessageID, err := msgResult.LastInsertId()
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("unable to get new messageID: %w", err)
	}

	// 5. Recupera il sent_at
	var sentAt string
	err = db.QueryRow(
		"SELECT sent_at FROM MessagesGroup WHERE messageID = ?",
		int(newMessageID),
	).Scan(&sentAt)
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("unable to get sent_at for messageID %d: %w", newMessageID, err)
	}

	return int(newMessageID), groupID, sentAt, content, contentType, nil
}

