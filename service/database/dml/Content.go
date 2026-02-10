package dml

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/entity"
)

func CreateTextContent(message_content entity.Content) (int, error) {
	//Imposta i default
	outErr := fmt.Errorf("Unable to insert the new content in the database")
	contentID := 0

	//dml := "INSERT ..."

	//err = InsertDataToDatabase(query)

	return contentID, outErr
}

func DeleteContentByID(contentID int) (error) {
	//Imposta i default
	outErr := fmt.Errorf("Unable to remove the contentID %d from database", contentID)

	//dml := "DELETE ..."
	//err = DeleteDataToDatabase(query)

	return outErr
}