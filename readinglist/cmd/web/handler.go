package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { //ensures users can't get her without right path
		http.NotFound(w, r)
		return
	}

	books, err := app.readinglist.GetAll() //queries web service
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "<html><head><title>Reading List</title></head><body><h1>Reading List</h1><ul>") //web page will render html
	for _, book := range *books {
		fmt.Fprintf(w, "<li>%s (%d)</li>", book.Title, book.Pages) //prints out each book
	}
	fmt.Fprintf(w, "</ul></body></html>")
}

func (app *application) bookView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id")) //gets the ID
	if err != nil || id < 1 {                        //check if ID is valid
		http.NotFound(w, r)
		return
	}

	book, err := app.readinglist.Get(int64(id)) //get the book
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s (%d)\n", book.Title, book.Pages)
}

func (app *application) bookCreate(w http.ResponseWriter, r *http.Request) { //add a new book to the entry
	switch r.Method { //creation form vs API call to create entry
	case http.MethodGet:
		app.bookCreateForm(w, r)
	case http.MethodPost:
		app.bookCreateProcess(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) bookCreateForm(w http.ResponseWriter, r *http.Request) { //create the form
	fmt.Fprintf(w, "<html><head><title>Create Book</title></head>"+
		"<body><h1>Create Book</h1><form action=\"/book/create\" method=\"post\">"+
		"<label for=\"title\">Title</label><input type=\"text\" name=\"title\" id=\"title\">"+
		"<label for=\"pages\">Pages</label><input type=\"number\" name=\"pages\" id=\"pages\">"+
		"<label for=\"published\">Published</label><input type=\"number\" name=\"published\" id=\"published\">"+
		"<label for=\"genres\">Genres</label><input type=\"text\" name=\"genres\" id=\"genres\">"+
		"<label for=\"rating\">Rating</label><input type=\"number\" step=\"0.1\" name=\"rating\" id=\"rating\">"+
		"<button type=\"submit\">Create</button></form></body></html>")
}

func (app *application) bookCreateProcess(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title") //verify for each required form entry (could you pair with better input validation like in .net?)
	if title == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	pages, err := strconv.Atoi(r.PostFormValue("pages")) //atoI means to convert to int
	if err != nil || pages < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	published, err := strconv.Atoi(r.PostFormValue("published"))
	if err != nil || published < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	genres := strings.Split(r.PostFormValue("genres"), " ")

	ratingFloat, err := strconv.ParseFloat(r.PostFormValue("rating"), 32)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	rating := float32(ratingFloat) //convert float64 to float32

	book := struct { //create struct to marshal
		Title     string   `json:"title"`
		Pages     int      `json:"pages"`
		Published int      `json:"published"`
		Genres    []string `json:"genres"`
		Rating    float32  `json:"rating"`
	}{
		Title:     title,
		Pages:     pages,
		Published: published,
		Genres:    genres,
		Rating:    rating,
	}

	data, err := json.Marshal(book) //encode into json
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	req, _ := http.NewRequest("POST", app.readinglist.Endpoint, bytes.NewBuffer(data)) //create post request to api
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req) //calls API
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close() //defer the closing of connection

	if resp.StatusCode != http.StatusCreated { //display if bad response
		log.Printf("unexpected status: %s", resp.Status)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther) //redirect to home
}
