package queries

import (
	"database/sql"
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
)

// GetConvByIDBelongUsers restituisce una conversazione tra utenti dato il conversationID
func GetConvByIDBelongUsers(conversationID int) (entity.Conversation, error) {
	var conv entity.Conversation
	var lastMessageID sql.NullInt64

	row := db.QueryRow(
		"SELECT conversationID, user1ID, user2ID, lastMessageID FROM ConversationsUser WHERE conversationID = ?",
		conversationID,
	)
	err := row.Scan(&conv.ConversationID, &conv.Sender, &conv.Receiver, &lastMessageID)
	if err != nil {
		return entity.Conversation{}, fmt.Errorf("unable to return conversation by id %d in ConversationsUser: %w", conversationID, err)
	}

	conv.Between_users = true
	if lastMessageID.Valid {
		conv.LastMessageID = int(lastMessageID.Int64)
	}

	return conv, nil
}

// GetConvByIDBelongGroups restituisce una conversazione utente-gruppo dato il conversationID
func GetConvByIDBelongGroups(conversationID int) (entity.Conversation, error) {
	var conv entity.Conversation
	var lastMessageID sql.NullInt64

	row := db.QueryRow(
		"SELECT conversationID, userID, groupID, lastMessageID FROM ConversationsGroup WHERE conversationID = ?",
		conversationID,
	)
	err := row.Scan(&conv.ConversationID, &conv.Sender, &conv.Receiver, &lastMessageID)
	if err != nil {
		return entity.Conversation{}, fmt.Errorf("unable to return conversation by id %d in ConversationsGroup: %w", conversationID, err)
	}

	conv.Between_users = false
	if lastMessageID.Valid {
		conv.LastMessageID = int(lastMessageID.Int64)
	}

	return conv, nil
}

/*
GetConvBelongUsers restituisce la lista delle conversazioni tra utenti per un dato userID.
Ogni elemento contiene: name, lastMessage (messageID, messageType, preview, timestamp), type, photo, conversationID.
Join tra ConversationsUser, MessagesUser, Contents e Users per ottenere il nome e la foto dell'altro utente
e l'anteprima dell'ultimo messaggio.
*/
func GetConvBelongUsers(userID int) ([]map[string]interface{}, error) {
	query := `
		SELECT
			cu.conversationID,
			CASE WHEN cu.user1ID = ? THEN u2.name ELSE u1.name END AS otherName,
			CASE WHEN cu.user1ID = ? THEN u2.photo ELSE u1.photo END AS otherPhoto,
			mu.messageID,
			c.type AS messageType,
			c.content AS preview,
			mu.sent_at AS timestamp
		FROM ConversationsUser cu
		JOIN Users u1 ON cu.user1ID = u1.userID
		JOIN Users u2 ON cu.user2ID = u2.userID
		LEFT JOIN MessagesUser mu ON cu.lastMessageID = mu.messageID
		LEFT JOIN Contents c ON mu.contentID = c.contentID
		WHERE cu.user1ID = ? OR cu.user2ID = ?
		ORDER BY mu.sent_at DESC
	`

	rows, err := db.Query(query, userID, userID, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("unable to return conversations belong users by userID %d: %w", userID, err)
	}
	defer rows.Close()

	conversations := []map[string]interface{}{}
	for rows.Next() {
		var convID int
		var name string
		var photo sql.NullString
		var msgID sql.NullInt64
		var msgType sql.NullString
		var preview sql.NullString
		var timestamp sql.NullString

		if err := rows.Scan(&convID, &name, &photo, &msgID, &msgType, &preview, &timestamp); err != nil {
			return nil, fmt.Errorf("error scanning conversation row: %w", err)
		}

		conv := map[string]interface{}{
			"conversationID": convID,
			"name":           name,
			"type":           "user",
		}

		if photo.Valid {
			conv["photo"] = photo.String
		} else {
			conv["photo"] = ""
		}

		lastMessage := map[string]interface{}{}
		if msgID.Valid {
			lastMessage["messageID"] = int(msgID.Int64)
			lastMessage["messageType"] = msgType.String
			lastMessage["preview"] = preview.String
			lastMessage["timestamp"] = timestamp.String
		}
		conv["lastMessage"] = lastMessage

		conversations = append(conversations, conv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating conversation rows: %w", err)
	}

	return conversations, nil
}

/*
GetConvBetweenUsersAndGroups restituisce la lista delle conversazioni utente-gruppo per un dato userID.
Ogni elemento contiene: name (del gruppo), lastMessage, type, photo, conversationID.
Join tra ConversationsGroup, Groups, MessagesGroup e Contents.
*/
func GetConvBetweenUsersAndGroups(userID int) ([]map[string]interface{}, error) {
	query := `
		SELECT
			cg.conversationID,
			g.name AS groupName,
			g.photo AS groupPhoto,
			mg.messageID,
			c.type AS messageType,
			c.content AS preview,
			mg.sent_at AS timestamp
		FROM ConversationsGroup cg
		JOIN Groups g ON cg.groupID = g.groupID
		LEFT JOIN MessagesGroup mg ON cg.lastMessageID = mg.messageID
		LEFT JOIN Contents c ON mg.contentID = c.contentID
		WHERE cg.userID = ?
		ORDER BY mg.sent_at DESC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("unable to return conversations between users and groups by userID %d: %w", userID, err)
	}
	defer rows.Close()

	conversations := []map[string]interface{}{}
	for rows.Next() {
		var convID int
		var groupName string
		var groupPhoto sql.NullString
		var msgID sql.NullInt64
		var msgType sql.NullString
		var preview sql.NullString
		var timestamp sql.NullString

		if err := rows.Scan(&convID, &groupName, &groupPhoto, &msgID, &msgType, &preview, &timestamp); err != nil {
			return nil, fmt.Errorf("error scanning group conversation row: %w", err)
		}

		conv := map[string]interface{}{
			"conversationID": convID,
			"name":           groupName,
			"type":           "group",
		}

		if groupPhoto.Valid {
			conv["photo"] = groupPhoto.String
		} else {
			conv["photo"] = ""
		}

		lastMessage := map[string]interface{}{}
		if msgID.Valid {
			lastMessage["messageID"] = int(msgID.Int64)
			lastMessage["messageType"] = msgType.String
			lastMessage["preview"] = preview.String
			lastMessage["timestamp"] = timestamp.String
		}
		conv["lastMessage"] = lastMessage

		conversations = append(conversations, conv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating group conversation rows: %w", err)
	}

	return conversations, nil
}

/*
GetMessagesUserList restituisce la lista dei messaggi di una conversazione tra utenti.
Per ogni messaggio restituisce: content, content_type, sent_at, senderID, userNameSenderID,
status, type, messageID. I commenti vengono aggregati separatamente.
Join tra MessagesUser, Contents e Users (per il nome del sender).
Per ogni messaggio, recupera anche i commenti associati dalla tabella Comments.
*/
func GetMessagesUserList(conversationID int) ([]map[string]interface{}, error) {
	// Prima recupera i messaggi
	query := `
		SELECT
			mu.messageID,
			c.content,
			c.type AS content_type,
			mu.sent_at,
			mu.senderID,
			u.name AS userNameSenderID,
			mu.status,
			mu.type
		FROM MessagesUser mu
		JOIN Contents c ON mu.contentID = c.contentID
		JOIN Users u ON mu.senderID = u.userID
		WHERE mu.conversationID = ?
		ORDER BY mu.sent_at ASC
	`

	rows, err := db.Query(query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("unable to return message list belong users by conversationID %d: %w", conversationID, err)
	}
	defer rows.Close()

	var messages []map[string]interface{}
	var messageIDs []int

	for rows.Next() {
		var msgID, senderID int
		var content, contentType, sentAt, userName, status, msgType string

		if err := rows.Scan(&msgID, &content, &contentType, &sentAt, &senderID, &userName, &status, &msgType); err != nil {
			return nil, fmt.Errorf("error scanning message row: %w", err)
		}

		msg := map[string]interface{}{
			"messageID":        msgID,
			"content":          content,
			"content_type":     contentType,
			"sent_at":          sentAt,
			"senderID":         senderID,
			"userNameSenderID": userName,
			"status":           status,
			"type":             msgType,
			"comment":          nil,
			"replyToMessageID": nil,
		}

		messages = append(messages, msg)
		messageIDs = append(messageIDs, msgID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating message rows: %w", err)
	}

	// Ora recupera i commenti per ogni messaggio
	for i, msgID := range messageIDs {
		comments, err := getCommentsForUserMessage(msgID)
		if err != nil {
			return nil, fmt.Errorf("error getting comments for messageID %d: %w", msgID, err)
		}
		if len(comments) > 0 {
			messages[i]["comment"] = comments
		}
	}

	return messages, nil
}

/*
GetMessagesGroupList restituisce la lista dei messaggi di una conversazione in un gruppo.
Stessa struttura di GetMessagesUserList ma usa MessagesGroup.
*/
func GetMessagesGroupList(conversationID int) ([]map[string]interface{}, error) {
	// Recupera i messaggi
	query := `
		SELECT
			mg.messageID,
			c.content,
			c.type AS content_type,
			mg.sent_at,
			mg.senderID,
			u.name AS userNameSenderID,
			mg.status,
			mg.type
		FROM MessagesGroup mg
		JOIN Contents c ON mg.contentID = c.contentID
		JOIN Users u ON mg.senderID = u.userID
		WHERE mg.conversationID = ?
		ORDER BY mg.sent_at ASC
	`

	rows, err := db.Query(query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("unable to return message list between users and groups by conversationID %d: %w", conversationID, err)
	}
	defer rows.Close()

	var messages []map[string]interface{}
	var messageIDs []int

	for rows.Next() {
		var msgID, senderID int
		var content, contentType, sentAt, userName, status, msgType string

		if err := rows.Scan(&msgID, &content, &contentType, &sentAt, &senderID, &userName, &status, &msgType); err != nil {
			return nil, fmt.Errorf("error scanning group message row: %w", err)
		}

		msg := map[string]interface{}{
			"messageID":        msgID,
			"content":          content,
			"content_type":     contentType,
			"sent_at":          sentAt,
			"senderID":         senderID,
			"userNameSenderID": userName,
			"status":           status,
			"type":             msgType,
			"comment":          nil,
			"replyToMessageID": nil,
		}

		messages = append(messages, msg)
		messageIDs = append(messageIDs, msgID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating group message rows: %w", err)
	}

	// Recupera i commenti per ogni messaggio di gruppo
	for i, msgID := range messageIDs {
		comments, err := getCommentsForGroupMessage(msgID)
		if err != nil {
			return nil, fmt.Errorf("error getting comments for group messageID %d: %w", msgID, err)
		}
		if len(comments) > 0 {
			messages[i]["comment"] = comments
		}
	}

	return messages, nil
}

// getCommentsForUserMessage recupera i commenti associati a un messaggio tra utenti
func getCommentsForUserMessage(messageID int) ([]interface{}, error) {
	query := `
		SELECT
			co.commentID,
			co.reaction,
			mu.senderID,
			u.name AS userNameSenderID
		FROM Comments co
		JOIN MessagesUser mu ON co.messageUserID = mu.messageID
		JOIN Users u ON mu.senderID = u.userID
		WHERE co.messageUserID = ?
	`

	rows, err := db.Query(query, messageID)
	if err != nil {
		return nil, fmt.Errorf("error querying comments for user messageID %d: %w", messageID, err)
	}
	defer rows.Close()

	var comments []interface{}
	for rows.Next() {
		var commentID, senderID int
		var reaction, userName string

		if err := rows.Scan(&commentID, &reaction, &senderID, &userName); err != nil {
			return nil, fmt.Errorf("error scanning comment row: %w", err)
		}

		comment := map[string]interface{}{
			"commentID":       commentID,
			"content":         reaction,
			"senderID":        senderID,
			"userNameSenderID": userName,
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

// getCommentsForGroupMessage recupera i commenti associati a un messaggio di gruppo
func getCommentsForGroupMessage(messageID int) ([]interface{}, error) {
	query := `
		SELECT
			co.commentID,
			co.reaction,
			mg.senderID,
			u.name AS userNameSenderID
		FROM Comments co
		JOIN MessagesGroup mg ON co.messageGroupID = mg.messageID
		JOIN Users u ON mg.senderID = u.userID
		WHERE co.messageGroupID = ?
	`

	rows, err := db.Query(query, messageID)
	if err != nil {
		return nil, fmt.Errorf("error querying comments for group messageID %d: %w", messageID, err)
	}
	defer rows.Close()

	var comments []interface{}
	for rows.Next() {
		var commentID, senderID int
		var reaction, userName string

		if err := rows.Scan(&commentID, &reaction, &senderID, &userName); err != nil {
			return nil, fmt.Errorf("error scanning group comment row: %w", err)
		}

		comment := map[string]interface{}{
			"commentID":       commentID,
			"content":         reaction,
			"senderID":        senderID,
			"userNameSenderID": userName,
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

