package queries

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/entity"
)

func GetGroupByID(groupID int) (entity.Group, error) {

	err := fmt.Errorf("Unable to return group by id %d", groupID)

	var group entity.Group

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return group, err
}
