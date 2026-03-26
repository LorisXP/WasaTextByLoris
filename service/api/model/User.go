package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
	"github.com/LorisXP/WasaTextByLoris/service/database/dml"
	"github.com/LorisXP/WasaTextByLoris/service/database/queries"
	"github.com/sirupsen/logrus"
)

// ErrUserNameAlreadyInUse viene restituito quando il nome utente richiesto è già presente nel sistema
var ErrUserNameAlreadyInUse = errors.New("userName already in use")

// ErrUserNotFound viene restituito quando la ricerca non produce risultati
var ErrUserNotFound = errors.New("no users found")

func GetUser(userID int) (entity.User, error) {
	logrus.Debug("Entered in GetUser()")
	logrus.Infof("Getting user %d", userID)

	// Imposta i default
	var user entity.User
	var outErr error

	user_found, err := queries.GetUserByID(userID)

	// Se non ci sono errori, prosegui
	if err == nil {

		user = user_found

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
	var outErr error

	// Rimuovi eventuali spazi
	newUserName = strings.TrimSpace(newUserName)

	// Verifica se il nome è già in uso da un altro utente
	existingID, errCheck := queries.GetUserIDByName(newUserName)

	if errCheck == nil && existingID != u.UserID {
		// Il nome è già usato da un altro utente
		outErr = fmt.Errorf("%w: '%s'", ErrUserNameAlreadyInUse, newUserName)
		logrus.WithField("userID", u.UserID).Warnf("userName '%s' is already in use by userID %d", newUserName, existingID)

	} else {

		// Aggiorna il nome
		err := dml.UpdateNameByUserID(u.UserID, newUserName)

		// Se non ci sono errori, prosegui
		if err == nil {

			u.Name = newUserName

			logrus.WithField("userID", u.UserID).Infof("userName updated with: %s", newUserName)

		} else {

			outErr = fmt.Errorf("cannot update userName for user %d: %w", u.UserID, err)
			logrus.Error(outErr)
		}
	}

	return outErr
}

// SetPhoto aggiorna la foto dell'utente
func SetPhoto(u *entity.User, newPhoto string) error {
	logrus.Debug("Entered in SetPhoto()")
	var outErr error

	// Effettua il caricamento della foto
	err := dml.UpdatePhotoByUserID(u.UserID, newPhoto)

	// Se non ci sono errori, prosegui
	if err == nil {

		u.Photo = newPhoto

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
	var outErr error

	records, err := queries.GetUsersByName(search)

	// Se non ci sono errori, prosegui
	if err == nil {

		users = make([]entity.User, 0, len(records))

		for _, r := range records {
			users = append(users, entity.User{
				UserID: r.UserID,
				Name:   r.Name,
				Photo:  r.Photo,
			})
		}

		if len(users) == 0 {
			outErr = fmt.Errorf("%w matching '%s'", ErrUserNotFound, search)
			logrus.WithField("search", search).Warn("No users found")
		} else {
			logrus.WithField("search", search).Infof("%d users retrieved", len(users))
		}

	} else {

		outErr = fmt.Errorf("error searching users by name '%s': %w", search, err)
		logrus.WithField("search", search).Error(outErr)
	}

	return users, outErr
}

func GetUserIdByName(userName string) (int, error) {
	logrus.Debug("Entered in GetUserIdByName()")
	var outErr error

	// Rimuovi eventuali spazi
	userName = strings.TrimSpace(userName)

	// Ottieni lo userID
	userID, errUserID := queries.GetUserIDByName(userName)

	// Se la ricerca è stata eseguita
	if errUserID == nil {

		logrus.Infof("Obtained userID by userName %s", userName)

	} else {
		outErr = fmt.Errorf("error during retrieving userID for userName %s: %w", userName, errUserID)
		logrus.Error(outErr)
	}

	return userID, outErr
}
