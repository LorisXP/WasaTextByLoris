package dml

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
)

// AddCommentBelongUsers inserisce una reaction (commento) a un messaggio tra utenti.
// Il campo messageUserID dell'entity deve essere valorizzato.
func AddCommentBelongUsers(comment entity.Comment) (int, error) {
	result, err := db.Exec(
		"INSERT INTO Comments (messageUserID, reaction) VALUES (?, ?)",
		comment.MessageUserID, comment.Reaction,
	)
	if err != nil {
		return 0, fmt.Errorf("unable to insert reaction to message %d belong users: %w", comment.MessageUserID, err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("unable to get commentID: %w", err)
	}

	return int(lastID), nil
}

// AddCommentBetweenUsersAndGroups inserisce una reaction (commento) a un messaggio di gruppo.
// Il campo messageGroupID dell'entity deve essere valorizzato.
func AddCommentBetweenUsersAndGroups(comment entity.Comment) (int, error) {
	result, err := db.Exec(
		"INSERT INTO Comments (messageGroupID, reaction) VALUES (?, ?)",
		comment.MessageGroupID, comment.Reaction,
	)
	if err != nil {
		return 0, fmt.Errorf("unable to insert reaction to message %d between user and group: %w", comment.MessageGroupID, err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("unable to get commentID: %w", err)
	}

	return int(lastID), nil
}

/*
DeleteCommentBelongUsers cancella un commento su un messaggio tra utenti.
Verifica che l'utente sia il sender del messaggio a cui appartiene il commento
tramite una JOIN tra Comments e MessagesUser.
*/
func DeleteCommentBelongUsers(userID int, commentID int, messageID int) error {
	// Verifica che il commento appartenga al messaggio e che l'utente sia il sender del messaggio
	var exists int
	err := db.QueryRow(`
		SELECT 1 FROM Comments co
		JOIN MessagesUser mu ON co.messageUserID = mu.messageID
		WHERE co.commentID = ? AND co.messageUserID = ? AND mu.senderID = ?
	`, commentID, messageID, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("commentID %d not found on messageID %d for user %d, or user not authorized: %w", commentID, messageID, userID, err)
	}

	// Cancella il commento
	result, err := db.Exec("DELETE FROM Comments WHERE commentID = ?", commentID)
	if err != nil {
		return fmt.Errorf("unable to delete commentID %d: %w", commentID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for commentID %d: %w", commentID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("commentID %d not found", commentID)
	}

	return nil
}

/*
DeleteCommentBetweenUsersAndGroups cancella un commento su un messaggio di gruppo.
Verifica che l'utente sia il sender del messaggio a cui appartiene il commento
tramite una JOIN tra Comments, MessagesGroup e ConversationsGroup.
*/
func DeleteCommentBetweenUsersAndGroups(userID int, commentID int, messageID int) error {
	// Verifica che il commento appartenga al messaggio e che l'utente sia il sender del messaggio
	var exists int
	err := db.QueryRow(`
		SELECT 1 FROM Comments co
		JOIN MessagesGroup mg ON co.messageGroupID = mg.messageID
		WHERE co.commentID = ? AND co.messageGroupID = ? AND mg.senderID = ?
	`, commentID, messageID, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("commentID %d not found on messageID %d for user %d, or user not authorized: %w", commentID, messageID, userID, err)
	}

	// Cancella il commento
	result, err := db.Exec("DELETE FROM Comments WHERE commentID = ?", commentID)
	if err != nil {
		return fmt.Errorf("unable to delete commentID %d: %w", commentID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected for commentID %d: %w", commentID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("commentID %d not found", commentID)
	}

	return nil
}