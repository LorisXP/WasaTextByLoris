package model

import (
	"encoding/base64"
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
	"github.com/LorisXP/WasaTextByLoris/service/database/dml"
	"github.com/LorisXP/WasaTextByLoris/service/database/queries"
	"github.com/sirupsen/logrus"
)

func CreateContent(content_type string, content string) (entity.Content, error) {
	logrus.Debug("Entered in CreateContent()")
	logrus.Info("Creating a new content")

	// Imposta i default
	var outErr error
	var err error
	var message_content entity.Content

	// Controlla il tipo
	if content_type == "text" || content_type == "gif" || content_type == "photo" {

		message_content.Type = content_type

		// Verifica che foto e gif siano in formato Base64
		if content_type == "photo" || content_type == "gif" {
			_, err := base64.StdEncoding.DecodeString(content)
			if err != nil {
				outErr = fmt.Errorf("error during creating content: content must be a valid Base64 string for type %s: %w", content_type, err)
				logrus.Error(outErr)
				return message_content, outErr
			}
		}

		// Inserisci il contenuto del messaggio
		message_content.Content = content

		// Crea l'inserimento a DB
		logrus.Info("Creating a new content in DB")

		message_content.ContentID, err = dml.CreateTextContent(message_content)
		logrus.Debug("Passed by CreateTextContent()")

		// Se non ci sono errori, prosegui
		if err == nil && message_content.ContentID != 0 {

			outErr = nil
			logrus.Info("Content inserted successfully")

		} else {
			outErr = fmt.Errorf("error during inserting content: %w", err)
			logrus.Error(outErr)
		}

	} else {
		outErr = fmt.Errorf("error during creating content: content type is not valid (%s). Expected 'text', 'gif' or 'photo'", content_type)
		logrus.Error(outErr)
	}

	return message_content, outErr
}

func GetContentByID(contentID int) (entity.Content, error) {
	logrus.Debug("Entered in GetContentByID() in package model/Message")

	// Imposta i default
	var outErr error
	var err error
	var content entity.Content

	// Ottieni le info a DB
	logrus.Infof("Getting content by id %d", contentID)
	content, err = queries.GetContentByID(contentID)

	// Se non ci sono errori, prosegui
	if err == nil && content.ContentID != 0 {
		outErr = nil
		logrus.Info("Content obtained successfully")
	} else {
		outErr = fmt.Errorf("error during obtaining content by id %d: %w", contentID, err)
		logrus.Error(outErr)
	}

	return content, outErr
}

func DeleteContentByID(contentID int) error {
	logrus.Debug("Entered in DeleteContentByID()")
	logrus.Warningf("contentID %d ready to delete", contentID)

	// Imposta i default
	var outErr error
	var err error

	// Cancella il contenuto dal DB
	logrus.Infof("Deleting contentID %d", contentID)

	err = dml.DeleteContentByID(contentID)
	logrus.Debug("Passed by DeleteContentByID()")

	// Se non ci sono errori, prosegui
	if err == nil {
		outErr = nil
		logrus.Infof("contentID %d removed from DB successfully", contentID)

	} else {
		outErr = fmt.Errorf("error during removing contentID %d: %w", contentID, err)
		logrus.Error(outErr)
	}

	return outErr
}
