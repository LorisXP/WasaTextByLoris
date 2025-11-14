package model

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/database/dml"
	"github.com/LorisXP/WasaTextByLoris/service/database/queries"
	"github.com/LorisXP/WasaTextByLoris/service/entity"
	"github.com/sirupsen/logrus"
)

// NewUser costruisce uno User dal database
func NewUser(userID int) (entity.User, error) {

	//Imposta i default
	var user entity.User
	var outErr error = fmt.Errorf("unable to create user %d", userID)

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

// SetPhoto aggiorna la foto dell'utente
func SetPhoto(u *entity.User, newPhoto string) error {

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
func GetUsersByName(search string) { //([]entity.User, error)

	// var users []User
	// var outErr error = fmt.Errorf("cannot search users by name '%s'", search)

	// records, err := queries.UsersByName(search)

	// //Se non ci sono errori, prosegui
	// if err == nil {

	// 	users = make([]User, 0, len(records))

	// 	for _, r := range records {
	// 		users = append(users, User{
	// 			userID: r.ID,
	// 			Name:   r.Name,
	// 			Photo:  r.Photo,
	// 		})
	// 	}

	// 	outErr = nil

	// 	logrus.WithField("search", search).Info("Users retrieved")

	// } else {

	// 	outErr = fmt.Errorf("error searching users by name '%s': %w", search, err)
	// 	logrus.WithField("search", search).Error(outErr)
	// }

	// return users, outErr
}
