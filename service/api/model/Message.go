package model

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
	"github.com/LorisXP/WasaTextByLoris/service/database/dml"
	"github.com/LorisXP/WasaTextByLoris/service/database/queries"
	"github.com/sirupsen/logrus"
)

/*
Create a new message to send. Require:

- userID of user want to send the message

- conversationID of conversation where you want to send the message from

- between_users true if yes false if it is between users and group

- content of message

- type of message_type IN ('text', 'gif', 'photo')
*/
func CreateMessage(userID int, conversationID int, between_users bool, content string, content_type string) (entity.Message, error) {
	logrus.Debug("Entered in CreateMessage()")
	logrus.Infof("Creating a new message in conversationID %d, between users: %t", conversationID, between_users)

	// Imposta i default
	var outErr error
	var err error
	var message entity.Message

	// Crea la struct message da salvare nel DB
	message.ConversationID = conversationID
	message.Between_users = between_users
	message.Sender = userID
	message.Type = "standard"
	message.Status = "received"
	message.Content = content
	message.ContentType = content_type

	// Crea prima il contenuto tramite model.CreateContent
	contentEntity, err := CreateContent(content_type, content)
	if err != nil {
		outErr = fmt.Errorf("error during creating content for message: %w", err)
		logrus.Error(outErr)
		return message, outErr
	}
	logrus.Debug("Passed by CreateContent()")

	// Crea l'inserimento a DB in base a se il messaggio è tra utenti o utente->gruppo
	if between_users {
		logrus.Info("Creating a new message between two users")

		message.MessageID, message.Receiver, message.Sent_at, err = dml.CreateMessageBelongUsers(message, contentEntity.ContentID)
		logrus.Debug("Passed by CreateMessageBelongUsers()")
	} else {
		logrus.Info("Creating a new message between users and group")

		message.MessageID, message.Receiver, message.Sent_at, err = dml.CreateMessageBetweenUsersAndGroups(message, contentEntity.ContentID)
		logrus.Debug("Passed by CreateMessageBetweenUsersAndGroups()")
	}

	// Se non ci sono errori, prosegui
	if err == nil && message.MessageID != 0 {

		outErr = nil
		logrus.Info("Message created successfully")

	} else {
		outErr = fmt.Errorf("error during creating message: %w", err)
		logrus.Error(outErr)
	}

	return message, outErr
}

func GetMessageByID(messageID int, between_users bool) (entity.Message, error) {
	logrus.Debug("Entered in GetMessageByID() in package model/Message")

	// Imposta i default
	var outErr error
	var err error
	var message entity.Message

	// Ottieni le info a DB in base a se sono tra utenti o utente->gruppo
	if between_users {
		logrus.Infof("Getting message by id %d between_users", messageID)

		message, err = queries.GetMesByIDBelongUsers(messageID)
	} else {
		logrus.Infof("Getting message by id %d between users and groups", messageID)

		message, err = queries.GetMesByIDBelongGroups(messageID)
	}

	// Se non ci sono errori, prosegui
	if err == nil && message.MessageID != 0 {
		outErr = nil
		logrus.Info("Message obtained successfully")
	} else {
		outErr = fmt.Errorf("error during obtaining message by id %d: %w", messageID, err)
		logrus.Error(outErr)
	}

	return message, outErr
}

func ForwardMessage(userID int, conversationID int, between_users bool, sourceBetweenUsers bool, messageID int) (entity.Message, error) {
	logrus.Debug("Entered in ForwardMessage()")
	logrus.Infof("Forwarding messageID %d → conversationID %d (dest between_users: %t, src between_users: %t)", messageID, conversationID, between_users, sourceBetweenUsers)

	// Imposta i default
	var outErr error
	var err error
	var message entity.Message

	// Crea la struct message da salvare nel DB
	message.MessageID = messageID
	message.ConversationID = conversationID
	message.Between_users = between_users
	message.SourceBetweenUsers = sourceBetweenUsers
	message.Sender = userID
	message.Type = "forward"
	message.Status = "received"

	// Crea l'inserimento a DB in base a se il messagio è tra utenti o utente->gruppo
	if between_users {
		logrus.Info("Forwarding a new message belong two users")

		message.MessageID, message.Receiver, message.Sent_at, message.Content, message.ContentType, err = dml.ForwardMessageBelongUsers(message)
		logrus.Debug("Passed by ForwardMessageBelongUsers()")
	} else {
		logrus.Info("Forwarding a new message between users and group")

		message.MessageID, message.Receiver, message.Sent_at, message.Content, message.ContentType, err = dml.ForwardMessageBetweenUsersAndGroups(message)
		logrus.Debug("Passed by ForwardMessageBetweenUsersAndGroups()")
	}

	// Se non ci sono errori, prosegui
	if err == nil && message.MessageID != 0 {

		outErr = nil
		logrus.Info("Message forwarded successfully")

	} else {
		outErr = fmt.Errorf("error during forwarding message: %w", err)
		logrus.Error(outErr)
	}

	return message, outErr
}

func DeleteMessage(userID int, conversationID int, between_users bool, messageID int) error {
	logrus.Debug("Entered in DeleteMessage()")
	logrus.Warningf("userID %d wants delete messageID %d in conversationID %d between users %t", userID, messageID, conversationID, between_users)

	// Imposta i default
	var outErr error
	var err error

	// Cancella il messaggio sul DB nella tabella corretta in base all'utente
	if between_users {
		logrus.Info("Moving to trash a message belong two users")

		err = dml.DeleteMessageBelongUsers(messageID, userID, conversationID)
		logrus.Debug("Passed by DeleteMessageBelongUsers()")
	} else {
		logrus.Info("Moving to trash message between users and group")

		err = dml.DeleteMessageBetweenUsersAndGroups(messageID, userID, conversationID)
		logrus.Debug("Passed by DeleteMessageBetweenUsersAndGroups()")
	}

	// Se non ci sono errori, prosegui
	if err == nil {
		outErr = nil
		logrus.Infof("userID %d deleted messageID %d in conversationID %d between users %t successfully", userID, messageID, conversationID, between_users)

	} else {
		outErr = fmt.Errorf("error for userID %d during removing messageID %d in conversationID %d between users %t: %w", userID, messageID, conversationID, between_users, err)
		logrus.Error(outErr)
	}

	return outErr
}
