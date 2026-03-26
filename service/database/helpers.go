package database

import (
	"database/sql"
	"fmt"
)

/*
ExecuteQuery executes a SELECT query and returns multiple rows.
Use this for queries that return multiple results.

Example usage:

	rows, err := appDB.ExecuteQuery("SELECT userID, name FROM Users WHERE name LIKE ?", "%john%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		err := rows.Scan(&u.UserID, &u.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
*/
func (db *appdbimpl) ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	// Imposta i default
	var outErr error
	var rows *sql.Rows

	// Esegui la query
	rows, err := db.c.Query(query, args...)

	// Se non ci sono errori, prosegui
	if err == nil {
	} else {
		outErr = fmt.Errorf("error executing query: %w", err)
	}

	return rows, outErr
}

/*
ExecuteQueryRow executes a SELECT query and returns a single row.
Use this for queries that return exactly one result.

Example usage:

	var user entity.User
	err := appDB.ExecuteQueryRow("SELECT userID, name, photo FROM Users WHERE userID = ?", userID).Scan(&user.UserID, &user.Name, &user.Photo)
	if err == sql.ErrNoRows {
		return entity.User{}, fmt.Errorf("user not found")
	}
	if err != nil {
		return entity.User{}, err
	}
*/
func (db *appdbimpl) ExecuteQueryRow(query string, args ...interface{}) *sql.Row {
	return db.c.QueryRow(query, args...)
}

/*
ExecuteInsert executes an INSERT statement and returns the last inserted ID.
Use this for inserting new records.

Example usage:

	lastID, err := appDB.ExecuteInsert("INSERT INTO Users (name, photo) VALUES (?, ?)", username, photo)
	if err != nil {
		return 0, fmt.Errorf("error inserting user: %w", err)
	}
	return int(lastID), nil
*/
func (db *appdbimpl) ExecuteInsert(query string, args ...interface{}) (int64, error) {
	// Imposta i default
	var outErr error
	var lastID int64

	// Esegui l'insert
	result, err := db.c.Exec(query, args...)

	// Se non ci sono errori, prosegui
	if err == nil {
		lastID, err = result.LastInsertId()
		if err == nil {
		} else {
			outErr = fmt.Errorf("error getting last insert ID: %w", err)
		}
	} else {
		outErr = fmt.Errorf("error executing insert: %w", err)
	}

	return lastID, outErr
}

/*
ExecuteUpdate executes an UPDATE statement and returns the number of affected rows.
Use this for updating existing records.

Example usage:

	rowsAffected, err := appDB.ExecuteUpdate("UPDATE Users SET photo = ? WHERE userID = ?", photo, userID)
	if err != nil {
		return fmt.Errorf("error updating user photo: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
*/
func (db *appdbimpl) ExecuteUpdate(query string, args ...interface{}) (int64, error) {
	// Imposta i default
	var outErr error
	var rowsAffected int64

	// Esegui l'update
	result, err := db.c.Exec(query, args...)

	// Se non ci sono errori, prosegui
	if err == nil {
		rowsAffected, err = result.RowsAffected()
		if err == nil {
		} else {
			outErr = fmt.Errorf("error getting rows affected: %w", err)
		}
	} else {
		outErr = fmt.Errorf("error executing update: %w", err)
	}

	return rowsAffected, outErr
}

/*
ExecuteDelete executes a DELETE statement and returns the number of affected rows.
Use this for deleting records.

Example usage:

	rowsAffected, err := appDB.ExecuteDelete("DELETE FROM Users WHERE userID = ?", userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
*/
func (db *appdbimpl) ExecuteDelete(query string, args ...interface{}) (int64, error) {
	// Imposta i default
	var outErr error
	var rowsAffected int64

	// Esegui il delete
	result, err := db.c.Exec(query, args...)

	// Se non ci sono errori, prosegui
	if err == nil {
		rowsAffected, err = result.RowsAffected()
		if err == nil {
		} else {
			outErr = fmt.Errorf("error getting rows affected: %w", err)
		}
	} else {
		outErr = fmt.Errorf("error executing delete: %w", err)
	}

	return rowsAffected, outErr
}
