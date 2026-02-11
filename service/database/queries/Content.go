package queries

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
)

func GetContentByID(contentID int) (entity.Content, error) {

	err := fmt.Errorf("Unable to return content by id %d in the table Contents", contentID)

	var content entity.Content

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return content, err
}
