package model

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/database/queries"
	"github.com/sirupsen/logrus"
)

func GetUserId(username string) (int, error) {
	logrus.Debug("Entered in GetUserId()")

	//Imposta i default
	outErr := fmt.Errorf("Unable to get userID by name of the user %s", username)

	//Opera
	userID, err := queries.GetUserIDByName(username)
	logrus.Debug("Passed by GetUserIDByName()")

	//Se non ci sono errori, prosegui
	if err == nil && userID != 0 {
		outErr = nil
		logrus.Info("Completed search of userID by username")
	} else {
		outErr = fmt.Errorf("userName %s not found: %w", username, err)
		logrus.Error(outErr)
	}

	//Resituisci
	return userID, outErr
}
