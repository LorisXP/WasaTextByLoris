package queries

import (
	"database/sql"
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
)

// GetCommentByID restituisce un commento (reaction) dato il suo commentID
func GetCommentByID(commentID int) (entity.Comment, error) {
	var comment entity.Comment
	var msgUserID sql.NullInt64
	var msgGroupID sql.NullInt64

	row := db.QueryRow(
		"SELECT commentID, messageUserID, messageGroupID, reaction FROM Comments WHERE commentID = ?",
		commentID,
	)
	err := row.Scan(&comment.CommentID, &msgUserID, &msgGroupID, &comment.Reaction)
	if err != nil {
		return entity.Comment{}, fmt.Errorf("unable to return comment by id %d: %w", commentID, err)
	}

	if msgUserID.Valid {
		comment.MessageUserID = int(msgUserID.Int64)
	}
	if msgGroupID.Valid {
		comment.MessageGroupID = int(msgGroupID.Int64)
	}

	return comment, nil
}
