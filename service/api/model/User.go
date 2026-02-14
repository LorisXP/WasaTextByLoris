package model

import (
	"fmt"
	"strings"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
	"github.com/LorisXP/WasaTextByLoris/service/database/dml"
	"github.com/LorisXP/WasaTextByLoris/service/database/queries"
	"github.com/sirupsen/logrus"
)

func GetUser(userID int) (entity.User, error) {
	logrus.Debug("Entered in GetUser()")
	logrus.Infof("Getting user %d", userID)

	//Imposta i default
	var user entity.User
	var outErr error = fmt.Errorf("unable to get user %d", userID)

	user_found, err := queries.GetUserByID(userID)

	//Se non ci sono errori, prosegui
	if err == nil {

		user = user_found

		outErr = nil
		logrus.WithField("userID", userID).Info("User loaded")

	} else {

		outErr = fmt.Errorf("cannot load user %d: %w", userID, err)
		logrus.Error(outErr)
	}

	return user, outErr
}

// SetUserName imposta un nuovo nome utente
func SetUserName(u *entity.User, newUserName string) error {
	logrus.Debug("Entered in SetUserName()")
	var outErr error = fmt.Errorf("cannot update userName for user %d", u.UserID)

	//Rimuovi eventuali spazi
	newUserName = strings.TrimSpace(newUserName)

	//Aggiorna il nome
	err := dml.UpdateNameByUserID(u.UserID, newUserName)

	//Se non ci sono errori, prosegui
	if err == nil {

		u.Name = newUserName
		outErr = nil

		logrus.WithField("userID", u.UserID).Infof("userName updated with: %s", newUserName)

	} else {

		outErr = fmt.Errorf("cannot update userName for user %d: %w", u.UserID, err)
		logrus.Error(outErr)
	}

	return outErr
}

// SetPhoto aggiorna la foto dell'utente
func SetPhoto(u *entity.User, newPhoto string) error {
	logrus.Debug("Entered in SetPhoto()")
	var outErr error = fmt.Errorf("cannot update photo for user %d", u.UserID)

	//Effettua il caricamento della foto
	err := dml.UpdatePhotoByUserID(u.UserID, newPhoto)

	//Se non ci sono errori, prosegui
	if err == nil {

		u.Photo = newPhoto
		outErr = nil

		logrus.WithField("userID", u.UserID).Info("Photo updated")

	} else {

		outErr = fmt.Errorf("cannot update photo for user %d: %w", u.UserID, err)
		logrus.Error(outErr)
	}

	return outErr
}

// GetUsersByName cerca utenti per nome
func GetUsersByName(search string) ([]entity.User, error) {

	var users []entity.User
	var outErr error = fmt.Errorf("cannot search users by name '%s'", search)

	records, err := queries.GetUsersByName(search)

	//Se non ci sono errori, prosegui
	if err == nil {

		users = make([]entity.User, 0, len(records))

		for _, r := range records {
			users = append(users, entity.User{
				UserID: r.UserID,
				Name:   r.Name,
				Photo:  r.Photo,
			})
		}

		outErr = nil

		logrus.WithField("search", search).Info("Users retrieved")

	} else {

		outErr = fmt.Errorf("error searching users by name '%s': %w", search, err)
		logrus.WithField("search", search).Error(outErr)
	}

	return users, outErr
}
