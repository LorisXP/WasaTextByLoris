package queries

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/entity"
)

func GetByIDBelongUsers(conversationID int) (entity.Conversation, error) {

	err := fmt.Errorf("Unable to return conversation by id %d in the table ConversationUsers", conversationID)

	var conversation entity.Conversation

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return conversation, err
}

func GetByIDBelongGroups(conversationID int) (entity.Conversation, error) {

	err := fmt.Errorf("Unable to return conversation by id %d in the table ConversationGroups", conversationID)

	var conversation entity.Conversation

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return conversation, err
}
