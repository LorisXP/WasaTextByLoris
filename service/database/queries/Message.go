package queries

import (
	"database/sql"
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
)

// GetMesByIDBelongUsers restituisce un messaggio tra utenti dato il messageID.
// JOIN con Contents per ottenere il contenuto e con Users per il sender.
func GetMesByIDBelongUsers(messageID int) (entity.Message, error) {
	var msg entity.Message

	row := db.QueryRow(`
		SELECT
			mu.messageID,
			mu.conversationID,
			mu.senderID,
			mu.receiverID,
			mu.sent_at,
			mu.status,
			mu.type,
			c.type AS contentType,
			c.content
		FROM MessagesUser mu
		JOIN Contents c ON mu.contentID = c.contentID
		WHERE mu.messageID = ?
	`, messageID)

	err := row.Scan(
		&msg.MessageID,
		&msg.ConversationID,
		&msg.Sender,
		&msg.Receiver,
		&msg.Sent_at,
		&msg.Status,
		&msg.Type,
		&msg.ContentType,
		&msg.Content,
	)
	if err != nil {
		return entity.Message{}, fmt.Errorf("unable to return message by id %d in MessagesUser: %w", messageID, err)
	}

	msg.Between_users = true

	return msg, nil
}

// GetMesByIDBelongGroups restituisce un messaggio di gruppo dato il messageID.
// JOIN con Contents per ottenere il contenuto.
// Il Receiver sarà il groupID dalla conversazione.
func GetMesByIDBelongGroups(messageID int) (entity.Message, error) {
	var msg entity.Message

	row := db.QueryRow(`
		SELECT
			mg.messageID,
			mg.conversationID,
			mg.senderID,
			cg.groupID,
			mg.sent_at,
			mg.status,
			mg.type,
			c.type AS contentType,
			c.content
		FROM MessagesGroup mg
		JOIN Contents c ON mg.contentID = c.contentID
		JOIN ConversationsGroup cg ON mg.conversationID = cg.conversationID
		WHERE mg.messageID = ?
	`, messageID)

	var groupID sql.NullInt64
	err := row.Scan(
		&msg.MessageID,
		&msg.ConversationID,
		&msg.Sender,
		&groupID,
		&msg.Sent_at,
		&msg.Status,
		&msg.Type,
		&msg.ContentType,
		&msg.Content,
	)
	if err != nil {
		return entity.Message{}, fmt.Errorf("unable to return message by id %d in MessagesGroup: %w", messageID, err)
	}

	msg.Between_users = false
	if groupID.Valid {
		msg.Receiver = int(groupID.Int64)
	}

	return msg, nil
}
