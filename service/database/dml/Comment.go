package dml

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/entity"
)

func AddCommentBelongUsers(comment entity.Comment) (int, error) {
	//Imposta i default
	outErr := fmt.Errorf("Unable to insert reaction to message belong users in the database")
	commentID := 0

	//dml := "INSERT ..."

	//err = InsertDataToDatabase(query)

	return commentID, outErr
}

func AddCommentBetweenUsersAndGroups(message entity.Message) (int, error) {
	//Imposta i default
	outErr := fmt.Errorf("Unable to insert reaction to message between a user and group in the database")
	commentID := 0

	//dml := "INSERT ..."

	//err = InsertDataToDatabase(query)

	return commentID, outErr
}

func DeleteCommentBelongUsers(userID int, commentID int, messageID int) (error) {
	//Imposta i default
	outErr := fmt.Errorf("Unable to remove the commentID %d of message belong users in the database", commentID)

	//Partendo dal commentID, bisogna fare una JOIN tra message e user per vedere se l'utente è abilitato alla cancellazione

	//dml := "DELETE ..."
	//err = DeleteDataToDatabase(query)

	return outErr
}

func DeleteCommentBetweenUsersAndGroups(userID int, commentID int, messageID int) (error) {
	//Imposta i default
	outErr := fmt.Errorf("Unable to remove the commentID %d of message between users and groups in the database", commentID)

	//Partendo dal commentID, bisogna fare una JOIN tra message, group e user per vedere se l'utente è abilitato alla cancellazione

	//dml := "DELETE ..."
	//err = DeleteDataToDatabase(query)

	return outErr
}