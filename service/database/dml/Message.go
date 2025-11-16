package dml

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/entity"
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
