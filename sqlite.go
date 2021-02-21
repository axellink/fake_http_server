package main

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

/*
Data processing error structure
Type is an integer following this mapping :
 - 1   : Cannot create database file
 - 2   : Cannot open database file
 - 3   : Cannot create Tables
 - 4   : Cannot insert new session
 - 5   : Cannot insert new request
 - 6   : Cannot purge old sessions
 - 100 : Session already exist in DB
 - 101 : Session does not exist in DB
*/
type DataError struct {
	Type    int
	Message string
}

/*
Let's meet Error interface contract
*/
func (err *DataError) Error() string {
	return err.Message
}

/*
Wrapper type for SQL database
Just because I prefer to do db.myfunction
TODO : add Lock to avoid "databse is locked" error while multiple access is going on (we will use a goroutine for the purge sooooo
*/
type DataBase struct {
	DB *sql.DB
}

/*
Open the sqlite database.
If file does not exist yet, creates it
If it does exist, will purge out of date data
*/
func OpenDB(filename string) (*DataBase, error) {
	// create file if it does not exist
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return nil, &DataError{1, "Cannot create database file : " + err.Error()}
		}
		file.Close()
	}
	sqldb, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, &DataError{2, "Cannot open database file : " + err.Error()}
	}
	db := &DataBase{DB: sqldb}
	err = db.createTables()
	return db, err
}

/*
Create database structure on file
Used by OpenDB function
*/
func (db *DataBase) createTables() error {
	SessionTable := `CREATE TABLE IF NOT EXISTS session (
			"id" TEXT NOT NULL UNIQUE,
			"creation_date" INTEGER NOT NULL
		);`
	statement, err := db.DB.Prepare(SessionTable)
	if err != nil {
		return &DataError{3, "Cannot create Session table : " + err.Error()}
	}
	_, err = statement.Exec()
	if err != nil {
		return &DataError{3, "Cannot create Session table : " + err.Error()}
	}

	RequestTable := `CREATE TABLE IF NOT EXISTS request (
			"id_session" TEXT NOT NULL,
			"data" TEXT NOT NULL
		);`
	statement, err = db.DB.Prepare(RequestTable)
	if err != nil {
		return &DataError{3, "Cannot create Request table : " + err.Error()}
	}
	statement.Exec()
	if err != nil {
		return &DataError{3, "Cannot create Session table : " + err.Error()}
	}

	return nil
}

/*
Insert a new session into database
If session exists already, will return a DataError Type 100
*/
func (db *DataBase) InsertSession(id string) error {
	// Check if id exists already
	rows, err := db.DB.Query("SELECT id FROM session WHERE id=?", id)
	if err != nil {
		rows.Close()
		return &DataError{4, "Cannot select from session to check if id exists : " + err.Error()}
	}

	if rows.Next() {
		rows.Close()
		return &DataError{100, "Session id " + id + " already exists"}
	}

	rows.Close()

	// actual insertion
	now := time.Now().Unix()

	statement, err := db.DB.Prepare("INSERT INTO session VALUES (?,?)")
	if err != nil {
		return &DataError{4, "Cannot insert new session : " + err.Error()}
	}

	_, err = statement.Exec(id, now)
	if err != nil {
		return &DataError{4, "Cannot insert new session : " + err.Error()}
	}

	return nil
}

/*
Insert a new request into database
If session does not exist, will return a DataError Type 101
*/
func (db *DataBase) InsertRequest(idSession string, request string) error {
	// Check if id exists
	rows, err := db.DB.Query("SELECT id FROM session WHERE id=?", idSession)
	if err != nil {
		rows.Close()
		return &DataError{5, "Cannot select from session to check if id exists : " + err.Error()}
	}

	if !rows.Next() {
		rows.Close()
		return &DataError{101, "Session id " + idSession + " does not exist"}
	}

	rows.Close()

	// So id exists, let's insert the request
	statement, err := db.DB.Prepare("INSERT INTO request VALUES (?,?)")
	if err != nil {
		return &DataError{5, "Cannot create request : " + err.Error()}
	}

	_, err = statement.Exec(idSession, request)
	if err != nil {
		return &DataError{5, "Cannot create request : " + err.Error()}
	}

	return nil
}

/*
Checks for too old sessions, deletes them if found with corresponding Requests
*/
func (db *DataBase) PurgeSessions(maxDuration int64) error {
	limitDate := time.Now().Unix() - maxDuration
	rows, err := db.DB.Query("SELECT id FROM session WHERE creation_date<?", limitDate)
	if err != nil {
		rows.Close()
		return &DataError{6, "Cannot purge old sessions : " + err.Error()}
	}

	var ids []string
	for rows.Next() {
		var idSession string
		rows.Scan(&idSession)
		ids = append(ids, idSession)
	}
	rows.Close()

	for _, idSession := range ids {
		statement, err := db.DB.Prepare("DELETE FROM request WHERE id_session = ?")
		if err != nil {
			return &DataError{6, "Cannot purge old sessions : " + err.Error()}
		}
		_, err = statement.Exec(idSession)
		if err != nil {
			return &DataError{6, "Cannot purge old sessions : " + err.Error()}
		}

		statement, err = db.DB.Prepare("DELETE FROM session WHERE id= ? ")
		if err != nil {
			return &DataError{6, "Cannot purge old sessions : " + err.Error()}
		}
		_, err = statement.Exec(idSession)
		if err != nil {
			return &DataError{6, "Cannot purge old sessions : " + err.Error()}
		}
	}
	return nil
}
