package queries

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
)

// GetContentByID restituisce un contenuto dato il suo contentID
func GetContentByID(contentID int) (entity.Content, error) {
	var content entity.Content

	row := db.QueryRow("SELECT contentID, type, content FROM Contents WHERE contentID = ?", contentID)
	err := row.Scan(&content.ContentID, &content.Type, &content.Content)
	if err != nil {
		return entity.Content{}, fmt.Errorf("unable to return content by id %d: %w", contentID, err)
	}

	return content, nil
}
