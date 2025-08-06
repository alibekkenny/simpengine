package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/alibekkenny/simpengine/cmd/config"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type application struct {
	errorLog  *log.Logger
	infoLog   *log.Logger
	accessLog *log.Logger
	config    *config.Config
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("ADDR")
	dsn := os.Getenv("DSN")
	jwt := os.Getenv("JWT_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	accessLog := log.New(os.Stdout, "ACCESS\t", log.Ldate|log.Ltime)

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed // before the main() function exits.
	defer db.Close()

	// Initialize a new instance of our application struct, containing the dependencies.
	config := &config.Config{
		JWTSecret: []byte(jwt),
		DB:        db,
	}

	app := &application{
		errorLog:  errorLog,
		infoLog:   infoLog,
		accessLog: accessLog,
		config:    config,
	}

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", addr)
	err = srv.ListenAndServe()

	errorLog.Fatal(err)
}

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
