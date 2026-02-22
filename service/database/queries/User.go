package queries

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
)

// GetUserByID restituisce un utente dato il suo userID
func GetUserByID(userID int) (entity.User, error) {
	var u entity.User

	row := db.QueryRow("SELECT userID, name, photo FROM Users WHERE userID = ?", userID)

	var photo sql.NullString
	err := row.Scan(&u.UserID, &u.Name, &photo)
	if err != nil {
		return entity.User{}, fmt.Errorf("unable to return user by id %d: %w", userID, err)
	}

	if photo.Valid {
		u.Photo = photo.String
	}

	return u, nil
}

// GetOrCreateUserIDByName cerca un utente per nome. Se esiste, restituisce lo userID.
// Se non esiste, lo crea e restituisce il nuovo userID. newUser indica se l'utente è stato appena creato.
func GetOrCreateUserIDByName(username string) (int, bool, error) {
	var userID int
	var newUser bool

	// Cerca l'utente per nome
	row := db.QueryRow("SELECT userID FROM Users WHERE name = ?", username)
	err := row.Scan(&userID)

	if err == sql.ErrNoRows {
		// L'utente non esiste: crealo
		result, insertErr := db.Exec("INSERT INTO Users (name) VALUES (?)", username)
		if insertErr != nil {
			return 0, false, fmt.Errorf("unable to create user by name %s: %w", username, insertErr)
		}

		lastID, idErr := result.LastInsertId()
		if idErr != nil {
			return 0, false, fmt.Errorf("unable to get last insert id for user %s: %w", username, idErr)
		}

		userID = int(lastID)
		newUser = true
	} else if err != nil {
		return 0, false, fmt.Errorf("unable to return userID by name %s: %w", username, err)
	}

	return userID, newUser, nil
}

// GetUsersIDByName restituisce una lista di userID a partire da una lista di nomi utente
func GetUsersIDByName(usernames []string) ([]int, error) {
	if len(usernames) == 0 {
		return []int{}, nil
	}

	// Costruisci i placeholder (?, ?, ...)
	placeholders := make([]string, len(usernames))
	args := make([]interface{}, len(usernames))
	for i, name := range usernames {
		placeholders[i] = "?"
		args[i] = name
	}

	query := fmt.Sprintf("SELECT userID FROM Users WHERE name IN (%s)", strings.Join(placeholders, ", "))

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to return userID list by names %v: %w", usernames, err)
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error scanning userID: %w", err)
		}
		userIDs = append(userIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return userIDs, nil
}

// GetUsersByName cerca utenti il cui nome contiene la stringa di ricerca (LIKE)
func GetUsersByName(search string) ([]entity.User, error) {
	rows, err := db.Query("SELECT userID, name, photo FROM Users WHERE name LIKE ?", "%"+search+"%")
	if err != nil {
		return nil, fmt.Errorf("unable to return users by name %s: %w", search, err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		var photo sql.NullString
		if err := rows.Scan(&u.UserID, &u.Name, &photo); err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}
		if photo.Valid {
			u.Photo = photo.String
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users rows: %w", err)
	}

	return users, nil
}

// GetUserIDByName restituisce lo userID di un utente dato il suo nome esatto
func GetUserIDByName(name string) (int, error) {
	var userID int

	row := db.QueryRow("SELECT userID FROM Users WHERE name = ?", name)
	err := row.Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("unable to return userID by name %s: %w", name, err)
	}

	return userID, nil
}
