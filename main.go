package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

type Page struct {
	id    int
	Title string
	Body  string
}

var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := getPageByTitle(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", &p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := getPageByTitle(title)
	if err != nil {
		p = Page{Title: title}
	}
	renderTemplate(w, "edit", &p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	_, err := saveBodyForTitle(title, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getPageByTitle(title string) (Page, error) {
	db, err := sql.Open("sqlite3", "data/main.db")
	if err != nil {
		log.Panic(err)
	}

	var body Page
	row := db.QueryRow("SELECT body, title FROM entries WHERE title = ?", title)
	if err := row.Scan(&body.Body, &body.Title); err != nil {
		return body, err
	}

	return body, nil
}

func saveBodyForTitle(title string, body string) (int64, error) {
	db, err := sql.Open("sqlite3", "data/main.db")
	if err != nil {
		log.Panic(err)
	}

	_, err = getPageByTitle(title)

	if err == sql.ErrNoRows {
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
