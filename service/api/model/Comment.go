package model

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
	"github.com/LorisXP/WasaTextByLoris/service/database/dml"
	"github.com/LorisXP/WasaTextByLoris/service/database/queries"
	"github.com/sirupsen/logrus"
)

/*
Create a new comment to message. Require:

- userID of user want to add a reaction

- messageID of message where to add the reaction

- between_users true if yes false if it is between users and group

- reaction emoji string
*/
func CreateComment(userID int, messageID int, between_users bool, reaction string) (entity.Comment, error) {
	logrus.Debug("Entered in CreateComment()")
	logrus.Infof("Adding a reaction to messageID %d, between users: %t", messageID, between_users)

	//Imposta i default
	outErr := fmt.Errorf("Unable to add a new reaction to messageID %d, between users: %t", messageID, between_users)
	err := fmt.Errorf("Unable to add a new reaction to messageID %d, between users: %t", messageID, between_users)
	var comment entity.Comment

	//Crea la struct comment da salvare nel DB
	comment.Reaction = reaction

	//Crea l'inserimento a DB in base a se il commento del messagio è tra utenti o utente->gruppo
	if between_users {
		comment.MessageUserID = messageID
		logrus.Infof("Adding a new reaction to message belong two users")
		
		comment.CommentID, err = dml.AddCommentBelongUsers(comment)
		logrus.Debug("Passed by AddCommentBelongUsers()")
	} else {
		comment.MessageGroupID = messageID
		logrus.Infof("Adding a new reaction to message between users and group")

		comment.CommentID, err = dml.AddCommentBetweenUsersAndGroups(comment)
		logrus.Debug("Passed by AddCommentBetweenUsersAndGroups()")
	}

	//Se non ci sono errori, prosegui
	if err == nil && comment.CommentID != 0 {
		outErr = nil
		logrus.Info("Added reaction to message succesfully")

	} else {
		outErr = fmt.Errorf("error during adding reaction to message: %w", err)
		logrus.Error(outErr)
	}

	return comment, outErr
}

func GetCommentByID(commentID int) (entity.Comment, error) {
	logrus.Infof("Getting a reaction with commentID %d", commentID)

	//Imposta i default
	outErr := fmt.Errorf("Unable to retrieve comment by id %d ", commentID)
	err := fmt.Errorf("Unable to retrieve comment by id %d ", commentID)
	var comment entity.Comment

	//Ottieni le info a DB in base a se sono tra utenti o utente->gruppo
	logrus.Infof("Getting comment by id %d", commentID)
	comment, err = queries.GetCommentByID(commentID)

	//Se non ci sono errori, prosegui
	if err == nil && comment.CommentID != 0 {
		outErr = nil
		logrus.Info("comment obtained succesfully")
	} else {
		outErr = fmt.Errorf("error during obtaining comment by id %d: %w", commentID, err)
		logrus.Error(outErr)
	}

	return comment, outErr
}

func DeleteComment(userID int, commentID int, messageID int, between_users bool, ) (error) {
	logrus.Debug("Entered in DeleteComment()")
	logrus.Warningf("userID %d wants delete commentID %d of messageID %d between users %t", userID, commentID, messageID, between_users)

	//Imposta i default
	outErr := fmt.Errorf("Unable for userID %d delete commentID %d of messageID %d between users %t", userID, commentID, messageID, between_users)
	err := fmt.Errorf("Unable for userID %d delete commentID %d of messageID %d between users %t", userID, commentID, messageID, between_users)

	//Cancella il messagio sul DB nella tabella corretta in base all'utente
	if between_users {
		logrus.Infof("Moving to trash a comment belong two users")

		err = dml.DeleteCommentBelongUsers(userID, commentID, messageID)
		logrus.Debug("Passed by DeleteCommentBelongUsers()")
	} else {
		logrus.Infof("Moving to trash comment between users and group")

		err = dml.DeleteCommentBetweenUsersAndGroups(userID, commentID, messageID)
		logrus.Debug("Passed by DeleteCommentBetweenUsersAndGroups()")
	}

	//Se non ci sono errori, prosegui
	if err == nil {
		outErr = nil
		logrus.Infof("userID %d SUCCESFULLY delete commentID %d of messageID %d between users %t", userID, commentID, messageID, between_users)

	} else {
		outErr = fmt.Errorf("error for userID %d during removing commentID %d of messageID %d between users %t: %w", userID, commentID, messageID, between_users,err)
		logrus.Error(outErr)
	}

	return outErr
}