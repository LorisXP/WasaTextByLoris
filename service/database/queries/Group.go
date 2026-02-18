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

// IsUserInGroup verifica se un utente è membro o admin di un gruppo
func IsUserInGroup(userID int, groupID int) (bool, error) {
	var count int

	// Controlla sia nella tabella Members che se è admin del gruppo
	err := db.QueryRow(
		`SELECT COUNT(*) FROM (
			SELECT userID FROM Members WHERE groupID = ? AND userID = ?
			UNION
			SELECT adminID FROM Groups WHERE groupID = ? AND adminID = ?
		)`,
		groupID, userID, groupID, userID,
	).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("unable to check membership for userID %d in groupID %d: %w", userID, groupID, err)
	}

	return count > 0, nil
}

// GetGroupMembers restituisce la lista dei membri di un gruppo (userName + photo)
func GetGroupMembers(groupID int) ([]entity.User, error) {
	rows, err := db.Query(
		`SELECT u.userID, u.name, u.photo
		 FROM Members m
		 JOIN Users u ON m.userID = u.userID
		 WHERE m.groupID = ?`, groupID)
	if err != nil {
		return nil, fmt.Errorf("unable to get members for groupID %d: %w", groupID, err)
	}
	defer rows.Close()

	var members []entity.User
	for rows.Next() {
		var u entity.User
		var photo sql.NullString
		err := rows.Scan(&u.UserID, &u.Name, &photo)
		if err != nil {
			return nil, fmt.Errorf("error scanning member for groupID %d: %w", groupID, err)
		}
		if photo.Valid {
			u.Photo = photo.String
		}
		members = append(members, u)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating members for groupID %d: %w", groupID, err)
	}

	return members, nil
}
