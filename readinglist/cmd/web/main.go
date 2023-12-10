package main

import (
	"flag"
	"log"
	"net/http"

	"readinglist.github.io/internal/models"
)

type application struct {
	readinglist *models.ReadingListModel
}

func main() {
	addr := flag.String("addr", ":80", "HTTP network address")
	endpoint := flag.String("endpoint", "http://localhost:4000/v1/books", "Endpoint for the readlingList web service")

	app := &application{
		readinglist: &models.ReadingListModel{Endpoint: *endpoint},
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	log.Printf("Starting the server on %s", *addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
