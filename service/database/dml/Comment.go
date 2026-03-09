package dml

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
)

// AddCommentBelongUsers inserisce una reazione (commento) a un messaggio tra utenti.
// Il campo MessageUserID dell'entity deve essere valorizzato.
func AddCommentBelongUsers(comment entity.Comment) (int, error) {
	result, err := db.Exec(
		"INSERT INTO Comments (senderID, messageUserID, reaction) VALUES (?, ?, ?)",
		comment.SenderID, comment.MessageUserID, comment.Reaction,
	)
	if err != nil {
		return 0, fmt.Errorf("unable to insert reaction to user message %d: %w", comment.MessageUserID, err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("unable to get commentID: %w", err)
	}

	return int(lastID), nil
}

// AddCommentBetweenUsersAndGroups inserisce una reazione (commento) a un messaggio di gruppo.
// Il campo MessageGroupID dell'entity deve essere valorizzato.
func AddCommentBetweenUsersAndGroups(comment entity.Comment) (int, error) {
	result, err := db.Exec(
		"INSERT INTO Comments (senderID, messageGroupID, reaction) VALUES (?, ?, ?)",
		comment.SenderID, comment.MessageGroupID, comment.Reaction,
	)
	if err != nil {
		return 0, fmt.Errorf("unable to insert reaction to group message %d: %w", comment.MessageGroupID, err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("unable to get commentID: %w", err)
	}

	return int(lastID), nil
}

/*
DeleteCommentBelongUsers cancella un commento su un messaggio tra utenti.
Verifica che l'utente sia l'autore della reazione (co.senderID).
*/
func DeleteCommentBelongUsers(userID int, commentID int, messageID int) error {
	// Verifica che il commento appartenga al messaggio e che l'utente sia l'autore della reazione
	var exists int
	err := db.QueryRow(`
		SELECT 1 FROM Comments co
		WHERE co.commentID = ? AND co.messageUserID = ? AND co.senderID = ?
	`, commentID, messageID, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("commentID %d not found for messageID %d or user %d not authorized: %w", commentID, messageID, userID, err)
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
Verifica che l'utente sia l'autore della reazione (co.senderID).
*/
func DeleteCommentBetweenUsersAndGroups(userID int, commentID int, messageID int) error {
	// Verifica che il commento appartenga al messaggio e che l'utente sia l'autore della reazione
	var exists int
	err := db.QueryRow(`
		SELECT 1 FROM Comments co
		WHERE co.commentID = ? AND co.messageGroupID = ? AND co.senderID = ?
	`, commentID, messageID, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("commentID %d not found for messageID %d or user %d not authorized: %w", commentID, messageID, userID, err)
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
