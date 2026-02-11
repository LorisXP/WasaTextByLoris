package dml

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
)

func CreateMessageBelongUsers(message entity.Message) (int, int, string, error) {
	//Imposta i default
	outErr := fmt.Errorf("Unable to insert the new message belong users in the database")
	messageID := 0
	receiver := 0
	sent_at := ""

	//dml := "INSERT ..."

	//err = InsertDataToDatabase(query)

	return messageID, receiver, sent_at, outErr
}

func CreateMessageBetweenUsersAndGroups(message entity.Message) (int, int, string, error) {
	//Imposta i default
	outErr := fmt.Errorf("Unable to insert the the new message between a user and group in the database")
	messageID := 0
	receiver := 0
	sent_at := ""

	//dml := "INSERT ..."

	//err = InsertDataToDatabase(query)

	return messageID, receiver, sent_at, outErr
}

func DeleteMessageBelongUsers(messageID int, userID int, conversationID int) (error) {
	//Imposta i default
	outErr := fmt.Errorf("Unable to remove the messageID %d belong users in the database", messageID)

	//dml := "DELETE ..."
	//err = DeleteDataToDatabase(query)

	return outErr
}

func DeleteMessageBetweenUsersAndGroups(messageID int, userID int, conversationID int) (error) {
	//Imposta i default
	outErr := fmt.Errorf("Unable to remove the messageID %d between users and groups in the database", messageID)

	//dml := "DELETE ..."
	//err = DeleteDataToDatabase(query)

	return outErr
}

func ForwardMessageBelongUsers(message entity.Message) (int, int, string, string, string, error) {
	//Imposta i default
	outErr := fmt.Errorf("Unable to insert the forwarded message belong users in the database")
	messageID := 0
	receiver := 0
	sent_at := ""
	content := ""
	contentType := ""

	//dml := "INSERT ..."

	//err = InsertDataToDatabase(query)

	return messageID, receiver, sent_at, content, contentType, outErr
}

func ForwardMessageBetweenUsersAndGroups(message entity.Message) (int, int, string, string, string, error) {
	//Imposta i default
	outErr := fmt.Errorf("Unable to insert the the forwarded message between a user and group in the database")
	messageID := 0
	receiver := 0
	sent_at := ""
	content := ""
	contentType := ""

	//dml := "INSERT ..."

	//err = InsertDataToDatabase(query)

	return  messageID, receiver, sent_at, content, contentType, outErr
}

