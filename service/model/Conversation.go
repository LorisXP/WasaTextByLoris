package model

import (
	"fmt"
	"sort"
	"time"

	"github.com/LorisXP/WasaTextByLoris/service/database/dml"
	"github.com/LorisXP/WasaTextByLoris/service/database/queries"
	"github.com/LorisXP/WasaTextByLoris/service/entity"
	"github.com/sirupsen/logrus"
)

/*
Create a new conversation between user (true) or group (false) by userID of sender and userID/groupID of receiver.
*/
func CreateConversation(between_users bool, sender int, receiver int) (entity.Conversation, error) {
	logrus.Debug("Entered in CreateConversation()")
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
func GetConversationByID(conversationID int, between_users bool) (entity.Conversation, error) {
	logrus.Debug("Entered in GetConversationByID() in package model/Conversation")

	//Imposta i default
	outErr := fmt.Errorf("Unable to retrieve conversation by id %d ", conversationID)
	err := fmt.Errorf("Unable to retrieve conversation by id %d ", conversationID)
	var conversation entity.Conversation

	//Ottieni le info a DB in base a se sono tra utenti o utente->gruppo
	if between_users {
		logrus.Infof("Getting conversation by id %d between_users", conversationID)

		conversation, err = queries.GetConvByIDBelongUsers(conversationID)
	} else {
		logrus.Infof("Getting conversation by id %d between users and groups", conversationID)

		conversation, err = queries.GetConvByIDBelongGroups(conversationID)
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
{
  "messages": [
    {
      "content": "aGVsbG8=",
      "timestamp": "2025-10-25T17:40:11Z",
      "sender": {
        "userName": "loris2155519",
        "userID": 1
      },
      "status": "received",
      "type": "standard",
      "comments": [
        {
          "content": "👍",
          "sender": {
            "userName": "loris2155519",
            "userID": 1
          },
          "commentID": 1
        }
      ],
      "messageID": 1,
      "replyToMessageID": 12
    }
  ]
}
*/
func GetListMessages(conversation *entity.Conversation) (map[string]interface{}, error) {
	logrus.Debug("Entered in GetListMessages()")
	logrus.Infof("Getting list of messages")

	//Imposta i default
	outErr := fmt.Errorf("Unable to get conversation message list")
	err := fmt.Errorf("Unable to get conversation message list")
	var result map[string]interface{}
	var messagesList []map[string]interface{}

	if conversation.Between_users {
		//Recupera tutti i messaggi con conversationID = conversation.ConversationID
		logrus.Infof("Retrieving list of messages of conversationID %d between two users", conversation.ConversationID)
		
		// GetMessagesUserList Ritornerà una lista di dict con:
		// content, content_type, sent_at, senderID, userNameSenderID, 
		// status, type, comment (content, senderID, userNameSenderID, commentID),
		// messageID, replyToMessageID
		messagesList, err = queries.GetMessagesUserList(conversation.ConversationID)
		logrus.Debug("Passed by GetMessagesUserList()")

	} else {
		//Recupera tutti i messaggi con conversationID = conversation.ConversationID
		logrus.Infof("Retrieving list of messages of conversationID %d between users and groups", conversation.ConversationID)
		
		// GetMessagesGroupList Ritornerà una lista di dict con:
		// content, content_type, sent_at, senderID, userNameSenderID, 
		// status, type, comment (content, senderID, userNameSenderID, commentID),
		// messageID, replyToMessageID
		messagesList, err = queries.GetMessagesGroupList(conversation.ConversationID)
		logrus.Debug("Passed by GetMessagesGroupList()")
	}

	//Se non ci sono errori, prosegui
	if err == nil && len(messagesList) > 0 {
		
		//Crea l'oggetto in output
		formattedMessages := []map[string]interface{}{}
		
		for _, msg := range messagesList {
			//Crea la struttura del sender
			sender := map[string]interface{}{
				"userName": msg["userNameSenderID"],
				"userID":   msg["senderID"],
			}
			
			//Prepara la lista dei commenti
			comments := []map[string]interface{}{}
			if msg["comment"] != nil {
				//Se ci sono commenti, processali (assumendo che sia una lista)
				if commentList, ok := msg["comment"].([]interface{}); ok {
					for _, c := range commentList {
						if comment, ok := c.(map[string]interface{}); ok {
							commentSender := map[string]interface{}{
								"userName": comment["userNameSenderID"],
								"userID":   comment["senderID"],
							}
							formattedComment := map[string]interface{}{
								"content":   comment["content"],
								"sender":    commentSender,
								"commentID": comment["commentID"],
							}
							comments = append(comments, formattedComment)
						}
					}
				}
			}
			
			//Crea il messaggio formattato
			formattedMsg := map[string]interface{}{
				"content":   msg["content"],
				"timestamp": msg["sent_at"],
				"sender":    sender,
				"status":    msg["status"],
				"type":      msg["type"],
				"comments":  comments,
				"messageID": msg["messageID"],
			}
			
			//Aggiungi replyToMessageID solo se presente
			if msg["replyToMessageID"] != nil {
				formattedMsg["replyToMessageID"] = msg["replyToMessageID"]
			}
			
			formattedMessages = append(formattedMessages, formattedMsg)
		}
		
		//Crea il risultato finale
		result = map[string]interface{}{
			"messages": formattedMessages,
		}
		
		outErr = nil
		logrus.Info("List of messages obtained succesfully")

	} else if err == nil && len(messagesList) == 0 {

		// messagesList è vuota
		outErr = nil
		logrus.Infof("No messages found for this conversationID %d", conversation.ConversationID)
		result = map[string]interface{}{
			"messages": []map[string]interface{}{},
		}

	} else {
		outErr = fmt.Errorf("error during obtaining messages by conversationID %d: %w", conversation.ConversationID, err)
		logrus.Error(outErr)
	}

	return result, outErr
}

/*
Returns a user's conversations (list) with other users and groups, using their user ID.
The list will contain an object like this:

	{
		"name": "loris2155519",
		"lastMessage":
		{
			"messageID": 1,
			"messageType": "text",
			"preview": "Prova messaggio",
			"timestamp": "2025-11-05T14:35:10Z"
		},
		"type": "group",
		"photo": "aGVsbG8=",
		"conversationID": 1
	}
*/
func GetConversationByUserID(userID int) ([]map[string]interface{}, error) {
	logrus.Debug("Entered in GetConversationByUserID() in package model/Conversation")
	logrus.Info("Returning conversation of user")
	logrus.Debugf("Returning conversation of userID %d", userID)

	//Imposta i default
	outErr := fmt.Errorf("Unable to retrieve conversations belong users of userID %d", userID)
	//Crea una lista di dict
	conversations := []map[string]interface{}{}

	//Ottieni le conversazioni prima tra utenti e poi nei gruppi, con ordine cronologico inverso
	conversations_users, err := queries.GetConvBelongUsers(userID)

	//Se non ci sono errori, prosegui
	if err == nil {
		logrus.Info("conversation belong users obtained succesfully")

		//Ora ottieni le conversazioni tra utente e gruppi
		conversations_groups, err := queries.GetConvBetweenUsersAndGroups(userID)

		//Se non ci sono errori, prosegui
		if err == nil {
			logrus.Info("conversation between users and groups obtained succesfully")

			//Ora uniscile
			conversations = append(conversations, conversations_users...)
			conversations = append(conversations, conversations_groups...)

			//Effetua l'ordinamento cronologico inverso per timestamp del messaggio
			sort.Slice(conversations, func(i, j int) bool {

				// Estrai timestamp i-esimo
				ts1 := conversations[i]["lastMessage"].(map[string]interface{})["timestamp"].(string)
				ts2 := conversations[j]["lastMessage"].(map[string]interface{})["timestamp"].(string)

				// Converti in time.Time
				t1, _ := time.Parse(time.RFC3339, ts1)
				t2, _ := time.Parse(time.RFC3339, ts2)

				// Ordine cronologico inverso → più recente prima
				return t1.After(t2)
			})

			logrus.Infof("conversation of userID %d, sorted succesfully", userID)

		} else {
			outErr = fmt.Errorf("error during obtaining conversation between users and groups by userID %d: %w", userID, err)
			logrus.Error(outErr)
		}

	} else {
		outErr = fmt.Errorf("error during obtaining conversation belong users by userID %d: %w", userID, err)
		logrus.Error(outErr)
	}

	return conversations, outErr
}
