package queries

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
)

func GetCommentByID(commentID int) (entity.Comment, error) {

	err := fmt.Errorf("Unable to return comment by id %d in the table Comments", commentID)

	var comment entity.Comment

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return comment, err
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
