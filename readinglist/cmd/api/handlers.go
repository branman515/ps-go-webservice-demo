package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"readinglist.github.io/internal/data"
)

// Return a health check in a json format via manual creation of the json message
func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	//Ensure this is a get method
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	data := map[string]string{ //set message
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	js, err := json.Marshal(data) //encode into json

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	js = append(js, '\n') //newline

	w.Header().Set("Content-Type", "application.json") //set the format type as json

	w.Write(js) //write the json response
}

func (app *application) getCreateBooksHandler(w http.ResponseWriter, r *http.Request) {
	//Ensure this is a get method
	if r.Method == http.MethodGet {
		books := []data.Book{
			{
				ID:        1,
				CreatedAt: time.Now(),
				Title:     "The Darkening of Tristram",
				Published: 1998,
				Pages:     300,
				Genres:    []string{"Fiction", "Thriller"},
				Rating:    4.5,
				Version:   1,
			},
			{
				ID:        2,
				CreatedAt: time.Now(),
				Title:     "The Legecy of Deckard Cain",
				Published: 2007,
				Pages:     432,
				Genres:    []string{"Fiction", "Adventure"},
				Rating:    4.9,
				Version:   1,
			},
			{
				ID:        3,          // system generated
				CreatedAt: time.Now(), // system generated
				Title:     "The Black Soulstone",
				Version:   1, // system generated
			},
		}

		if err := app.writeJSON(w, http.StatusOK, envelope{"books": books}); err != nil { //envelop the json response with books:[]
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		var input struct {
			Title     string   `json:"title"`
			Published int      `json:"published"`
			Pages     int      `json:"pages"`
			Genres    []string `json:"genres"`
			Rating    float32  `json:"rating"`
		}

		err := app.readJSON(w, r, &input)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "%v\n", input)

	}

}

func (app *application) getUpdateDeleteBooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getBook(w, r)
	case http.MethodPut:
		app.updateBook(w, r)
	case http.MethodDelete:
		app.deleteBook(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(id, 10, 64) //parse the book id (base 10, 64bit size)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	book := data.Book{ //define a hard coded entry for now
		ID:        idInt,
		CreatedAt: time.Now(),
		Title:     "Echoes in the Darkness",
		Published: 2019,
		Pages:     300,
		Genres:    []string{"Fiction", "Thriller"},
		Rating:    4.5,
		Version:   1,
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"book": book}); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) updateBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	//marshal: Encode an object into a byte slice
	var input struct { //Defines how to unmarshal (decode into an object)
		Title     *string  `json:"title"`
		Published *int     `json:"published"`
		Pages     *int     `json:"pages"`
		Genres    []string `json:"genres"`
		Rating    *float32 `json:"rating"`
	}

	book := data.Book{ //create a mock object
		ID:        idInt,
		CreatedAt: time.Now(),
		Title:     "Echoes in the Darkness",
		Published: 2020,
		Pages:     325,
		Genres:    []string{"Fiction", "Thriller"},
		Rating:    5.5,
		Version:   1,
	}

	err = app.readJSON(w, r, &input) //read in the body to parse
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//checks if any new information is in the decoded request object
	if input.Title != nil {
		book.Title = *input.Title
	}

	if input.Published != nil {
		book.Published = *input.Published
	}

	if input.Pages != nil {
		book.Pages = *input.Pages
	}

	if len(input.Genres) > 0 {
		book.Genres = input.Genres
	}

	if input.Rating != nil {
		book.Rating = *input.Rating
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"book": book}); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Delete a specific book with ID: %d", idInt)
}
