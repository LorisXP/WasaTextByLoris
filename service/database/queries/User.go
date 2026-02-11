package queries

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
)

func GetUserByID(userID int) (entity.User, error) {

	err := fmt.Errorf("Unable to return user by id %d", userID)
	var u entity.User

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return u, err
}

func GetUserIDByName(username string) (int, error) {

	err := fmt.Errorf("Unable to return userID by name %s", username)
	var userID int

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return userID, err
}

func GetUsersIDByName(usernames []string) ([]int, error) {

	err := fmt.Errorf("Unable to return userID list by this names %v", usernames)
	var userID []int

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return userID, err
}

func GetUsersByName(search string) ([]entity.User, error) {

	err := fmt.Errorf("Unable to return users by name %s", search)
	var users []entity.User

	//query := "SELECT ..."

	//result = getDataByDatabase(query)

	return users, err
}
