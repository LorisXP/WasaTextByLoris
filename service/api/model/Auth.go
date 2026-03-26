package model

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/database/queries"
	"github.com/sirupsen/logrus"
)

// Se utente esite -> restituisce userID e photo
// Altrimenti -> nuovo utente, nuovo userID
func AuthUser(username string) (int, string, bool, error) {
	logrus.Debug("Entered in AuthUser()")
	logrus.Infof("The user %s is authenticating now...", username)

	// Imposta i default
	var newUser bool = false
	var photo string
	var outErr error

	// Fai una query di ricerca
	userID, isNew, err := queries.GetOrCreateUserIDByName(username)
	logrus.Debug("Passed by GetOrCreateUserIDByName()")

	// Se non ci sono errori, prosegui
	if err == nil && userID != 0 {
		newUser = isNew

		if newUser {
			logrus.Infof("%s is a new user", username)
		} else {
			logrus.Infof("%s authenticated successfully", username)
			// Carica la foto profilo dell'utente esistente
			u, errUser := queries.GetUserByID(userID)
			if errUser == nil {
				photo = u.Photo
			}
		}

	} else {
		outErr = fmt.Errorf("error during authenticating userName %s: %w", username, err)
		logrus.Error(outErr)
	}

	// Restituisci
	return userID, photo, newUser, outErr
}
