package queries

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/entity"
)

func GetMesByIDBelongUsers(messageID int) (entity.Message, error) {

	err := fmt.Errorf("Unable to return message by id %d in the table MessagesUser", messageID)

	var message entity.Message

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return message, err
}

func GetMesByIDBelongGroups(messageID int) (entity.Message, error) {

	err := fmt.Errorf("Unable to return message by id %d in the table MessagesGroup", messageID)

	var message entity.Message

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return message, err
}

/*
func GetBelongUsers(userID int) ([]map[string]interface{}, error) {

	err := fmt.Errorf("Unable to return conversation belong users by userID %d from database", userID)

	conversations_users := []map[string]interface{}{}

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return conversations_users, err
}

func GetBetweenUsersAndGroups(userID int) ([]map[string]interface{}, error) {

	err := fmt.Errorf("Unable to return conversation between users and groups by userID %d from database", userID)

	conversations_groups := []map[string]interface{}{}

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return conversations_groups, err
}

*/
