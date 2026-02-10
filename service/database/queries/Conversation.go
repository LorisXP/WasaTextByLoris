package queries

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/entity"
)

func GetConvByIDBelongUsers(conversationID int) (entity.Conversation, error) {

	err := fmt.Errorf("Unable to return conversation by id %d in the table ConversationUsers", conversationID)

	var conversation entity.Conversation

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return conversation, err
}

func GetConvByIDBelongGroups(conversationID int) (entity.Conversation, error) {

	err := fmt.Errorf("Unable to return conversation by id %d in the table ConversationGroups", conversationID)

	var conversation entity.Conversation

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return conversation, err
}

func GetConvBelongUsers(userID int) ([]map[string]interface{}, error) {

	err := fmt.Errorf("Unable to return conversation belong users by userID %d from database", userID)

	conversations_users := []map[string]interface{}{}

	/*
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

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return conversations_users, err
}

func GetConvBetweenUsersAndGroups(userID int) ([]map[string]interface{}, error) {

	err := fmt.Errorf("Unable to return conversation between users and groups by userID %d from database", userID)

	conversations_groups := []map[string]interface{}{}

	/*
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

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return conversations_groups, err
}
