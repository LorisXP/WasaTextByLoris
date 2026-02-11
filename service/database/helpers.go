package database

import (
	"database/sql"
	"fmt"
)

/*
ExecuteQuery executes a SELECT query and returns multiple rows.
Use this for queries that return multiple results.

Example usage:

	rows, err := database.ExecuteQuery(db, "SELECT userID, name FROM Users WHERE name LIKE ?", "%john%")
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
func ExecuteQuery(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	return rows, nil
}

/*
ExecuteQueryRow executes a SELECT query and returns a single row.
Use this for queries that return exactly one result.

Example usage:

	var user entity.User
	err := database.ExecuteQueryRow(db, "SELECT userID, name, photo FROM Users WHERE userID = ?", userID).Scan(&user.UserID, &user.Name, &user.Photo)
	if err == sql.ErrNoRows {
		return entity.User{}, fmt.Errorf("user not found")
	}
	if err != nil {
		return entity.User{}, err
	}
*/
func ExecuteQueryRow(db *sql.DB, query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

/*
ExecuteInsert executes an INSERT statement and returns the last inserted ID.
Use this for inserting new records.

Example usage:

	lastID, err := database.ExecuteInsert(db, "INSERT INTO Users (name, photo) VALUES (?, ?)", username, photo)
	if err != nil {
		return 0, fmt.Errorf("error inserting user: %w", err)
	}
	return int(lastID), nil
*/
func ExecuteInsert(db *sql.DB, query string, args ...interface{}) (int64, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("error executing insert: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %w", err)
	}

	return lastID, nil
}

/*
ExecuteUpdate executes an UPDATE statement and returns the number of affected rows.
Use this for updating existing records.

Example usage:

	rowsAffected, err := database.ExecuteUpdate(db, "UPDATE Users SET photo = ? WHERE userID = ?", photo, userID)
	if err != nil {
		return fmt.Errorf("error updating user photo: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
*/
func ExecuteUpdate(db *sql.DB, query string, args ...interface{}) (int64, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("error executing update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error getting rows affected: %w", err)
	}

	return rowsAffected, nil
}

/*
ExecuteDelete executes a DELETE statement and returns the number of affected rows.
Use this for deleting records.

Example usage:

	rowsAffected, err := database.ExecuteDelete(db, "DELETE FROM Users WHERE userID = ?", userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
*/
func ExecuteDelete(db *sql.DB, query string, args ...interface{}) (int64, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("error executing delete: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error getting rows affected: %w", err)
	}

	return rowsAffected, nil
}

/*
GetDB returns the underlying *sql.DB connection from AppDatabase.
Use this to access the raw DB connection when calling helper functions.

Example usage:

	db := database.GetDB(appDB)
	rows, err := database.ExecuteQuery(db, "SELECT * FROM Users")
*/
func GetDB(appDB AppDatabase) *sql.DB {
	// Type assertion to get the underlying implementation
	if impl, ok := appDB.(*appdbimpl); ok {
		return impl.c
	}
	return nil
}
