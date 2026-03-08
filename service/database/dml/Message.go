package dml

import (
	"fmt"
	"strings"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
)

/*
CreateMessageBelongUsers inserisce un nuovo messaggio in una conversazione tra utenti.
Il contenuto deve essere già stato creato tramite CreateTextContent.
Passaggi:
1. Determina il destinatario dalla conversazione (l'altro utente)
2. Inserisce il messaggio nella tabella MessagesUser
3. Restituisce messageID, receiverID, sent_at
*/
func CreateMessageBelongUsers(message entity.Message, contentID int) (int, int, string, error) {
	// 1. Determina il destinatario dalla conversazione
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

	// 2. Inserisce il messaggio in MessagesUser
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

	// 3. Recupera il campo sent_at generato dal DEFAULT
	var sentAt string
	err = db.QueryRow(
		"SELECT sent_at FROM MessagesUser WHERE messageID = ?",
		int(messageID),
	).Scan(&sentAt)
	if err != nil {
		return 0, 0, "", fmt.Errorf("unable to get sent_at for messageID %d: %w", messageID, err)
	}

	// 4. Aggiorna lastMessageID della conversazione esistente
	_, err = db.Exec(
		"UPDATE ConversationsUser SET lastMessageID = ? WHERE conversationID = ?",
		int(messageID), message.ConversationID,
	)
	if err != nil {
		return 0, 0, "", fmt.Errorf("unable to update lastMessageID for conversation %d: %w", message.ConversationID, err)
	}

	return int(messageID), receiverID, sentAt, nil
}

/*
CreateMessageBetweenUsersAndGroups inserisce un nuovo messaggio in una conversazione utente-gruppo.
Il contenuto deve essere già stato creato tramite CreateTextContent.
Passaggi:
1. Determina il groupID (destinatario) dalla conversazione
2. Inserisce il messaggio nella tabella MessagesGroup
3. Restituisce messageID, groupID (destinatario), sent_at
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

	// 2. Inserisce il messaggio in MessagesGroup
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

	// 3. Recupera il campo sent_at generato dal DEFAULT
	var sentAt string
	err = db.QueryRow(
		"SELECT sent_at FROM MessagesGroup WHERE messageID = ?",
		int(messageID),
	).Scan(&sentAt)
	if err != nil {
		return 0, 0, "", fmt.Errorf("unable to get sent_at for messageID %d: %w", messageID, err)
	}

	// 4. Aggiorna lastMessageID per TUTTE le conversazioni del gruppo (tutti i membri)
	_, err = db.Exec(
		"UPDATE ConversationsGroup SET lastMessageID = ? WHERE groupID = ?",
		int(messageID), groupID,
	)
	if err != nil {
		return 0, 0, "", fmt.Errorf("unable to update lastMessageID for group %d: %w", groupID, err)
	}

	return int(messageID), groupID, sentAt, nil
}

/*
DeleteMessageBelongUsers cancella un messaggio tra utenti.
Verifica che l'utente sia il mittente del messaggio.
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

	// Aggiorna lastMessageID se il messaggio cancellato era l'ultimo della conversazione
	var currentLastMessageID *int
	err = db.QueryRow(
		"SELECT lastMessageID FROM ConversationsUser WHERE conversationID = ?",
		conversationID,
	).Scan(&currentLastMessageID)
	if err == nil && currentLastMessageID != nil && *currentLastMessageID == messageID {
		// Trova il nuovo ultimo messaggio (il più recente rimasto)
		var newLastMessageID *int
		err = db.QueryRow(
			"SELECT MAX(messageID) FROM MessagesUser WHERE conversationID = ?",
			conversationID,
		).Scan(&newLastMessageID)
		if err != nil {
			return fmt.Errorf("unable to find new lastMessageID for conversation %d: %w", conversationID, err)
		}
		// newLastMessageID sarà NULL se non ci sono più messaggi
		_, err = db.Exec(
			"UPDATE ConversationsUser SET lastMessageID = ? WHERE conversationID = ?",
			newLastMessageID, conversationID,
		)
		if err != nil {
			return fmt.Errorf("unable to update lastMessageID for conversation %d: %w", conversationID, err)
		}
	}

	return nil
}

/*
DeleteMessageBetweenUsersAndGroups cancella un messaggio in una conversazione utente-gruppo.
Verifica che l'utente sia il mittente del messaggio.
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

	// Aggiorna lastMessageID se il messaggio cancellato era l'ultimo della conversazione
	var currentLastMessageID *int
	err = db.QueryRow(
		"SELECT lastMessageID FROM ConversationsGroup WHERE conversationID = ?",
		conversationID,
	).Scan(&currentLastMessageID)
	if err == nil && currentLastMessageID != nil && *currentLastMessageID == messageID {
		// Trova il nuovo ultimo messaggio (il più recente rimasto)
		var newLastMessageID *int
		err = db.QueryRow(
			"SELECT MAX(messageID) FROM MessagesGroup WHERE conversationID = ?",
			conversationID,
		).Scan(&newLastMessageID)
		if err != nil {
			return fmt.Errorf("unable to find new lastMessageID for group conversation %d: %w", conversationID, err)
		}
		// newLastMessageID sarà NULL se non ci sono più messaggi
		_, err = db.Exec(
			"UPDATE ConversationsGroup SET lastMessageID = ? WHERE conversationID = ?",
			newLastMessageID, conversationID,
		)
		if err != nil {
			return fmt.Errorf("unable to update lastMessageID for group conversation %d: %w", conversationID, err)
		}
	}

	return nil
}

/*
ForwardMessageBelongUsers inoltra un messaggio in una conversazione tra utenti.
Passaggi:
1. Recupera il contenuto del messaggio originale (message.MessageID)
2. Copia il contenuto tramite CreateTextContent
3. Determina il destinatario dalla conversazione di destinazione
4. Inserisce il nuovo messaggio con type='forward'
5. Restituisce messageID, receiver, sent_at, content, contentType
*/
func ForwardMessageBelongUsers(message entity.Message) (int, int, string, string, string, error) {
	// 1. Recupera il contenuto del messaggio originale verificando che il sender abbia accesso
	var content, contentType string
	var err error

	if message.SourceBetweenUsers {
		// Il sender deve essere uno dei due partecipanti della conversazione originale
		err = db.QueryRow(
			`SELECT c.content, c.type
			FROM MessagesUser mu
			JOIN Contents c ON mu.contentID = c.contentID
			JOIN ConversationsUser cu ON mu.conversationID = cu.conversationID
			WHERE mu.messageID = ? AND (cu.user1ID = ? OR cu.user2ID = ?)`,
			message.MessageID, message.Sender, message.Sender,
		).Scan(&content, &contentType)
	} else {
		// Il sender deve essere membro del gruppo a cui appartiene la conversazione originale
		err = db.QueryRow(
			`SELECT c.content, c.type
			FROM MessagesGroup mg
			JOIN Contents c ON mg.contentID = c.contentID
			JOIN ConversationsGroup cg ON mg.conversationID = cg.conversationID
			JOIN Members m ON cg.groupID = m.groupID AND m.userID = ?
			WHERE mg.messageID = ?`,
			message.Sender, message.MessageID,
		).Scan(&content, &contentType)
	}
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("message %d not found or user %d has no access: %w", message.MessageID, message.Sender, err)
	}

	// 2. Copia il contenuto tramite CreateTextContent
	newContent := entity.Content{Type: contentType, Content: content}
	newContentID, err := CreateTextContent(newContent)
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("unable to copy content for forwarded message: %w", err)
	}

	// 3. Determina il destinatario dalla conversazione di destinazione
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

	// 4. Inserisce il messaggio con type='forward'
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

	// 5. Recupera il campo sent_at
	var sentAt string
	err = db.QueryRow(
		"SELECT sent_at FROM MessagesUser WHERE messageID = ?",
		int(newMessageID),
	).Scan(&sentAt)
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("unable to get sent_at for messageID %d: %w", newMessageID, err)
	}

	// 6. Aggiorna lastMessageID della conversazione
	_, err = db.Exec(
		"UPDATE ConversationsUser SET lastMessageID = ? WHERE conversationID = ?",
		int(newMessageID), message.ConversationID,
	)
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("unable to update lastMessageID for conversation %d: %w", message.ConversationID, err)
	}

	return int(newMessageID), receiverID, sentAt, content, contentType, nil
}

/*
ForwardMessageBetweenUsersAndGroups inoltra un messaggio in una conversazione utente-gruppo.
Stessa logica di ForwardMessageBelongUsers ma per MessagesGroup.
*/
func ForwardMessageBetweenUsersAndGroups(message entity.Message) (int, int, string, string, string, error) {
	// 1. Recupera il contenuto del messaggio originale verificando che il sender abbia accesso
	var content, contentType string
	var err error

	if message.SourceBetweenUsers {
		// Il sender deve essere uno dei due partecipanti della conversazione originale
		err = db.QueryRow(
			`SELECT c.content, c.type
			FROM MessagesUser mu
			JOIN Contents c ON mu.contentID = c.contentID
			JOIN ConversationsUser cu ON mu.conversationID = cu.conversationID
			WHERE mu.messageID = ? AND (cu.user1ID = ? OR cu.user2ID = ?)`,
			message.MessageID, message.Sender, message.Sender,
		).Scan(&content, &contentType)
	} else {
		// Il sender deve essere membro del gruppo a cui appartiene la conversazione originale
		err = db.QueryRow(
			`SELECT c.content, c.type
			FROM MessagesGroup mg
			JOIN Contents c ON mg.contentID = c.contentID
			JOIN ConversationsGroup cg ON mg.conversationID = cg.conversationID
			JOIN Members m ON cg.groupID = m.groupID AND m.userID = ?
			WHERE mg.messageID = ?`,
			message.Sender, message.MessageID,
		).Scan(&content, &contentType)
	}
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("message %d not found or user %d has no access: %w", message.MessageID, message.Sender, err)
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

	// 4. Inserisce il messaggio con type='forward'
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

	// 5. Recupera il campo sent_at
	var sentAt string
	err = db.QueryRow(
		"SELECT sent_at FROM MessagesGroup WHERE messageID = ?",
		int(newMessageID),
	).Scan(&sentAt)
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("unable to get sent_at for messageID %d: %w", newMessageID, err)
	}

	// 6. Aggiorna lastMessageID per TUTTE le conversazioni del gruppo (tutti i membri)
	_, err = db.Exec(
		"UPDATE ConversationsGroup SET lastMessageID = ? WHERE groupID = ?",
		int(newMessageID), groupID,
	)
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("unable to update lastMessageID for group %d: %w", groupID, err)
	}

	return int(newMessageID), groupID, sentAt, content, contentType, nil
}

/*
MarkGroupMessagesAsReadByUser registra la lettura dei messaggi di gruppo da parte di un utente.
Per ogni messageID inserisce una riga in GroupMessageReads (INSERT OR IGNORE per idempotenza).
Poi aggiorna status='read' in MessagesGroup per i messaggi in cui tutti i membri
(escluso il sender) hanno letto.
*/
func MarkGroupMessagesAsReadByUser(messageIDs []int, userID int, groupID int) error {
	if len(messageIDs) == 0 {
		return nil
	}

	// 1. Bulk INSERT OR IGNORE delle read receipts
	placeholders := strings.Repeat("(?,?),", len(messageIDs))
	placeholders = placeholders[:len(placeholders)-1]

	args := make([]interface{}, 0, len(messageIDs)*2)
	for _, id := range messageIDs {
		args = append(args, id, userID)
	}

	_, err := db.Exec(
		"INSERT OR IGNORE INTO GroupMessageReads (messageID, userID) VALUES "+placeholders,
		args...,
	)
	if err != nil {
		return fmt.Errorf("unable to insert group message reads: %w", err)
	}

	// 2. Aggiorna status='read' per i messaggi in cui tutti i membri (escluso il sender) hanno letto
	ph := strings.Repeat("?,", len(messageIDs))
	ph = ph[:len(ph)-1]

	updateArgs := make([]interface{}, 0, len(messageIDs)+1)
	for _, id := range messageIDs {
		updateArgs = append(updateArgs, id)
	}
	updateArgs = append(updateArgs, groupID)

	_, err = db.Exec(
		`UPDATE MessagesGroup SET status = 'read'
		 WHERE messageID IN (`+ph+`)
		   AND status = 'received'
		   AND (
		       SELECT COUNT(*) FROM GroupMessageReads gmr WHERE gmr.messageID = MessagesGroup.messageID
		   ) >= (
		       SELECT COUNT(*) FROM Members m WHERE m.groupID = ? AND m.userID != MessagesGroup.senderID
		   )`,
		updateArgs...,
	)
	if err != nil {
		return fmt.Errorf("unable to update group messages status to read: %w", err)
	}

	return nil
}

/*
MarkMessagesAsRead aggiorna lo status da "received" a "read" per una lista di messageID.
Viene eseguita una singola query UPDATE con una clausola IN costruita dinamicamente.
*/
func MarkMessagesAsRead(messageIDs []int) error {
	if len(messageIDs) == 0 {
		return nil
	}

	placeholders := strings.Repeat("?,", len(messageIDs))
	placeholders = placeholders[:len(placeholders)-1] // rimuove la virgola finale

	args := make([]interface{}, len(messageIDs))
	for i, id := range messageIDs {
		args[i] = id
	}

	_, err := db.Exec(
		"UPDATE MessagesUser SET status = 'read' WHERE messageID IN ("+placeholders+")",
		args...,
	)
	if err != nil {
		return fmt.Errorf("unable to mark messages as read: %w", err)
	}
	return nil
}
