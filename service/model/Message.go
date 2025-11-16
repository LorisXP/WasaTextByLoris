package model

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/database/dml"
	//"github.com/LorisXP/WasaTextByLoris/service/database/queries"
	"github.com/LorisXP/WasaTextByLoris/service/entity"
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
func CreateMessage(userID int, conversationID int, between_users bool, content string, message_type string) (entity.Message, error) {
	logrus.Debug("Entered in CreateMessage()")
	logrus.Infof("Creating a new message in conversationID %d, between users: %t", conversationID, between_users)

	//Imposta i default
	outErr := fmt.Errorf("Unable to create a new message in conversationID %d, between users: %t", conversationID, between_users)
	err := fmt.Errorf("Unable to create a new message in conversationID %d, between users: %t", conversationID, between_users)
	var message entity.Message

	//Crea la struct message da salvare nel DB
	message.ConversationID = conversationID
	message.Between_users = between_users
	message.Sender = userID
	message.Status = "received"
	message.Content = content
	message.Type = message_type

	//Crea l'inserimento a DB in base a se il messagio è tra utenti o utente->gruppo
	if between_users {
		logrus.Infof("Creating a new message belong two users")

		message.MessageID, message.Receiver, message.Sent_at, err = dml.CreateMessageBelongUsers(message)
		logrus.Debug("Passed by CreateMessageBelongUsers()")
	} else {
		logrus.Infof("Creating a new message between users and group")

		message.MessageID, message.Receiver, message.Sent_at, err = dml.CreateMessageBetweenUsersAndGroups(message)
		logrus.Debug("Passed by CreateMessageBetweenUsersAndGroups()")
	}

	//Se non ci sono errori, prosegui
	if err == nil && message.MessageID != 0 {

		outErr = nil
		logrus.Info("Message created succesfully")

	} else {
		outErr = fmt.Errorf("error during creating message: %w", err)
		logrus.Error(outErr)
	}

	return message, outErr
}
