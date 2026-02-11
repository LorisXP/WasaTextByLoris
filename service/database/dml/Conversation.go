package dml

import (
	"fmt"
)

func CreateConversationBetweenUsers(sender int, receiver int) (int, error) {
	//Imposta i default
	outErr := fmt.Errorf("Unable to insert the new conversation between users in the database")
	conversationID := 0

	//dml := "INSERT ..."

	//err = InsertDataToDatabase(query)

	return conversationID, outErr
}

func CreateConversationBetweenGroups(userID int, groupID int) (int, error) {
	//Imposta i default
	outErr := fmt.Errorf("Unable to insert the new conversation between a user and group in the database")
	conversationID := 0

	//dml := "INSERT ..."

	//err = InsertDataToDatabase(query)

	return conversationID, outErr
}

func UpdateLastMessageIDforUser(conversationID int, messageID int) error {
	//Imposta i default
	outErr := fmt.Errorf("Unable to update lastMessageID %d value in the table ConversationsUser by conversationID %d", messageID, conversationID)

	//dml := "UPDATE ..."

	//err = UpdateDataToDatabase(query)

	return outErr
}

func UpdateLastMessageIDforGroup(conversationID int, messageID int) error {
	//Imposta i default
	outErr := fmt.Errorf("Unable to update lastMessageID %d value in the table ConversationsGroup by conversationID %d", messageID, conversationID)

	//dml := "UPDATE ..."

	//err = UpdateDataToDatabase(query)

	return outErr
}
