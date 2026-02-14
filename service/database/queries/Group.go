package queries

import (
	"database/sql"
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
)

// GetGroupByID restituisce un gruppo dato il suo groupID
func GetGroupByID(groupID int) (entity.Group, error) {
	var group entity.Group

	row := db.QueryRow("SELECT groupID, name, photo, adminID FROM Groups WHERE groupID = ?", groupID)

	var photo sql.NullString
	err := row.Scan(&group.GroupID, &group.Name, &photo, &group.AdminID)
	if err != nil {
		return entity.Group{}, fmt.Errorf("unable to return group by id %d: %w", groupID, err)
	}

	if photo.Valid {
		group.Photo = photo.String
	}

	return group, nil
}
