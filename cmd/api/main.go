package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chrisganov/go-game-api/internal/data"
	_ "github.com/lib/pq"
)

type application struct {
	models data.Models
	logger *log.Logger
}

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	port := os.Getenv("PORT")
	databaseConnectionString := os.Getenv("GAME_DB_CONNECTION")

  if port == "" {
    logger.Fatal("Port is missing")
  }

	if databaseConnectionString == "" {
		logger.Fatal("DB Connection string missing")
	}

	db, err := openDB(databaseConnectionString)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	app := &application{
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", port),
		Handler: app.routes(),
	}

	logger.Printf("Starting server on %s", port)

	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// TODO:
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
