package model

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/database/dml"
	"github.com/LorisXP/WasaTextByLoris/service/database/queries"
	"github.com/LorisXP/WasaTextByLoris/service/entity"
	"github.com/sirupsen/logrus"
)

/*
Create a new conversation between user (true) or group (false) by userID of sender and userID/groupID of receiver.
*/
func Create(between_users bool, sender int, receiver int) (entity.Conversation, error) {
	logrus.Debug("Entered in Create() in package model/Conversation")
	logrus.Infof("Creating a new conversation")

	//Imposta i default
	outErr := fmt.Errorf("Unable to create a new conversation")
	var err error = nil
	var conversation entity.Conversation

	//Crea l'inserimento a DB in base a se sono tra utenti o utente->gruppo
	if between_users {
		logrus.Infof("Creating a new conversation between two users")
		logrus.Debugf("They are '%d' and '%d'", sender, receiver)

		conversation.ConversationID, err = dml.CreateConversationBetweenUsers(sender, receiver)
		logrus.Debug("Passed by CreateConversationBetweenUsers()")
	} else {
		logrus.Infof("Creating a new conversation between user and groupID %d", receiver)
		logrus.Debugf("They are '%d' and '%d'", sender, receiver)

		conversation.ConversationID, err = dml.CreateConversationBetweenGroups(sender, receiver)
		logrus.Debug("Passed by CreateConversationBetweenGroups()")
	}

	//Se non ci sono errori, prosegui
	if err == nil && conversation.ConversationID != 0 {

		//Crea la struct conversation
		conversation.Between_users = between_users
		conversation.Sender = sender
		conversation.Receiver = receiver
		conversation.LastMessageID = 0

		outErr = nil

		logrus.Info("Conversation created succesfully")

	} else {
		outErr = fmt.Errorf("error during creating conversation: %w", err)
		logrus.Error(outErr)
	}

	return conversation, outErr
}

/*
Update lastMessageID of a conversation
*/
func SetLastMessageID(conversation *entity.Conversation, messageID int) error {
	logrus.Debug("Entered in SetLastMessageID()")
	logrus.Infof("Updating last message id (%d) in conversationID %d", messageID, conversation.ConversationID)

	//Imposta i default
	outErr := fmt.Errorf("Unable to update last message id (%d) in conversationID %d", messageID, conversation.ConversationID)
	err := fmt.Errorf("Unable to update last message id (%d) in conversationID %d", messageID, conversation.ConversationID)

	//Crea l'inserimento a DB in base a se sono tra utenti o utente->gruppo
	if conversation.Between_users {
		logrus.Infof("Updating lastMessageID %d in a conversationID %d between two users", messageID, conversation.ConversationID)

		err = dml.UpdateLastMessageIDforUser(conversation.ConversationID, messageID)
		logrus.Debug("Passed by UpdateLastMessageIDforUser()")
	} else {
		logrus.Infof("Updating lastMessageID %d in a conversationID %d between user and group", messageID, conversation.ConversationID)

		err = dml.UpdateLastMessageIDforGroup(conversation.ConversationID, messageID)
		logrus.Debug("Passed by UpdateLastMessageIDforGroup()")
	}

	//Se non ci sono errori, prosegui
	if err == nil {
		conversation.LastMessageID = messageID
		outErr = nil
		logrus.Info("last message ID updated succesfully")
	} else {
		outErr = fmt.Errorf("error during updating last message ID (%d): %w", messageID, err)
		logrus.Error(outErr)
	}

	return outErr
}

/*
Return a conversation struct in WasaText by his conversationID
*/
func GetByID(conversationID int, between_users bool) (entity.Conversation, error) {
	logrus.Debug("Entered in GetByID() in package model/Conversation")

	//Imposta i default
	outErr := fmt.Errorf("Unable to retrieve conversation by id %d ", conversationID)
	err := fmt.Errorf("Unable to retrieve conversation by id %d ", conversationID)
	var conversation entity.Conversation

	//Ottieni le info a DB in base a se sono tra utenti o utente->gruppo
	if between_users {
		logrus.Infof("Getting conversation by id %d between_users", conversationID)

		conversation, err = queries.GetByIDBelongUsers(conversationID)
	} else {
		logrus.Infof("Getting conversation by id %d between users and groups", conversationID)

		conversation, err = queries.GetByIDBelongGroups(conversationID)
	}

	//Se non ci sono errori, prosegui
	if err == nil && conversation.ConversationID != 0 {
		outErr = nil
		logrus.Info("conversation obtained succesfully")
	} else {
		outErr = fmt.Errorf("error during obtaining conversation by id %d: %w", conversationID, err)
		logrus.Error(outErr)
	}

	return conversation, outErr
}

/*
Return a list of messages of specific conversation (belonging users or group)
*/
func GetListMessages(conversation *entity.Conversation) {
}

/*
Restituisce le conversazioni di un utente tra altri utenti e gruppi, attraverso il suo userID
*/
func GetByUser(conversation *entity.Conversation) {

}
