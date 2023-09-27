package main

import (
	"database/sql"
	"fmt"
	"log"
)

var db = initDataBaseAccess()

func initDataBaseAccess() *sql.DB {
	db, err := sql.Open("sqlite3", "data/main.db")
	if err != nil {
		log.Panic(err)
	}
	return db
}

func getPageByTitle(title string) (Page, error) {
	var body Page
	row := db.QueryRow("SELECT body, title FROM entries WHERE title = ?", title)
	if err := row.Scan(&body.Body, &body.Title); err != nil {
		return body, err
	}

	return body, nil
}

func saveBodyForTitle(title string, body string) (int64, error) {
	if _, err := getPageByTitle(title); err == sql.ErrNoRows {
		return insertBodyForTitle(title, body)
	}
	return updateBodyForTitle(title, body)
}

func insertBodyForTitle(title string, body string) (int64, error) {
	result, err := db.Exec("INSERT INTO entries (title, body) VALUES (?,?)", title, body)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

func updateBodyForTitle(title string, body string) (int64, error) {
	result, err := db.Exec("UPDATE entries SET body = ? WHERE title = ?", body, title)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}
