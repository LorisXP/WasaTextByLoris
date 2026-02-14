package model

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/database/queries"
	"github.com/sirupsen/logrus"
)

//Se utente esite -> restituisce userID
//Altrimenti -> nuovo utente, nuovo userID
func AuthUser(username string) (int, error) {
	logrus.Debug("Entered in AuthUser()")
	logrus.Infof("The user %s is authenticating now...", username)

	//Imposta i default
	outErr := fmt.Errorf("Unable to return userID for username %s", username)

	//Fai una query di ricerca
	userID, newUser, err := queries.GetOrCreateUserIDByName(username)
	logrus.Debug("Passed by GetUserIDByName()")

	//Se non ci sono errori, prosegui
	if err == nil && userID != 0 {
		outErr = nil

		if newUser {
			logrus.Info("%s is a new user", username)
		} else{
			logrus.Info("%s authenticated succesfully")
		}

	} else {
		outErr = fmt.Errorf("error during authenticating userName %s: %w", username, err)
		logrus.Error(outErr)
	}

	//Resituisci
	return userID, outErr
}
