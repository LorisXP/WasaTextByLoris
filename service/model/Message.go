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

func GetMessage(){
	
}

/* !!!! DA FINIRE ANCORA !!!! */
func ForwardMessage(userID int, conversationID int, between_users bool, messageID int) (entity.Message, error){
	logrus.Debug("Entered in ForwardMessage()")
	logrus.Infof("Forwarding a new message in conversationID %d, between users: %t", conversationID, between_users)

	//Imposta i default
	outErr, err := fmt.Errorf("Unable to forward a new message in conversationID %d, between users: %t", conversationID, between_users)
	var message entity.Message

	//Crea la struct message da salvare nel DB
	message.ConversationID = conversationID
	message.Between_users = between_users
	message.Sender = userID
	message.Status = "received"

	//Crea l'inserimento a DB in base a se il messagio è tra utenti o utente->gruppo
	if between_users {
		logrus.Infof("Creating a new message belong two users")

		/*Il DB deve restituire:
		1. Content
		2. Type

		*/
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

func DeleteMessage(userID int, conversationID int, between_users bool, messageID int) (error) {
	logrus.Debug("Entered in DeleteMessage()")
	logrus.Warningf("userID %d wants delete messageID %d in conversationID %d between users %t", userID, messageID, conversationID, between_users)

	//Imposta i default
	outErr, err := fmt.Errorf("Unable for userID %d delete messageID %d in conversationID %d between users %t", userID, messageID, conversationID, between_users)
	logrus.Debug("Passed by DeleteGroup()")

	//Cancella il messagio sul DB nella tabella corretta in base all'utente
	if between_users {
		logrus.Infof("Moving to trash a message belong two users")

		err = dml.DeleteMessageBelongUsers(message)
		logrus.Debug("Passed by DeleteMessageBelongUsers()")
	} else {
		logrus.Infof("Moving to trash message between users and group")

		err = dml.DeleteMessageBetweenUsersAndGroups(message)
		logrus.Debug("Passed by DeleteMessageBetweenUsersAndGroups()")
	}

	//Se non ci sono errori, prosegui
	if err == nil {
		outErr = nil
		logrus.Infof("userID %d SUCCESFULLY delete messageID %d in conversationID %d between users %t", userID, messageID, conversationID, between_users)

	} else {
		outErr = fmt.Errorf("error for userID %d during removing messageID %d in conversationID %d between users %t: %w", userID, messageID, conversationID, between_users,err)
		logrus.Error(outErr)
	}

	return outErr
}