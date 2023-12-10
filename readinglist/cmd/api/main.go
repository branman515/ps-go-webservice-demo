package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"readinglist.github.io/internal/data"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	dsn  string
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	//pass in config values
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|stage|prod)")
	flag.StringVar(&cfg.dsn, "db-dsn", os.Getenv("READINGLIST_DB_DSN"), "PostgreSQL DSN")
	flag.Parse()

	//define the logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//define database
	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("database connection pool established")

	//define an app object to store information for each handler
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	//set the listening port/endpoint
	addr := fmt.Sprintf(":%d", cfg.port)

	//create our own mux to prevent modification of global handler
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.route(), //Mux handler
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	//start up the web service using the Server
	logger.Printf("Starting %s server on %s", cfg.env, addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)

}
