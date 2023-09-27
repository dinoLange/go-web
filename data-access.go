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

func (p *Page) save() (int64, error) {
	if _, err := getPageByTitle(p.Title); err == sql.ErrNoRows {
		return p.insert()
	}
	return p.update()
}

func (p *Page) insert() (int64, error) {
	result, err := db.Exec("INSERT INTO entries (title, body) VALUES (?,?)", p.Title, p.Body)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

func (p *Page) update() (int64, error) {
	result, err := db.Exec("UPDATE entries SET body = ? WHERE title = ?", p.Body, p.Title)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}
