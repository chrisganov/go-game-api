package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"go_game_api.com/internal/data"
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "admin"
	dbPassword = "12345"
	dbName     = "game-api"
)

type application struct {
	models data.Models
	logger *log.Logger
}

func main() {
	addr := flag.String("addr", ":8080", "Port of the application")

	flag.Parse()

	// var databasePass string
	// databasePass = os.Getenv("GAME_DB_CONNECTION")
	// fmt.Printf("Database: %s\n", databasePass)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(psqlInfo)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	app := &application{
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost%s", *addr),
		Handler: app.routes(),
	}

	logger.Printf("Starting server on %s", *addr)

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
