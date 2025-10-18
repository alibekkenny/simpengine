package main

// @title           SimpEngine API
// @version         1.0
// @description     API for managing romantic events in SimpEngine project
// @host
// @BasePath        /
// @schemes https http

// @securityDefinitions.apikey BearerAuth
// @in cookie
// @name jwt

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/alibekkenny/simpengine/cmd/config"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type application struct {
	errorLog  *log.Logger
	infoLog   *log.Logger
	accessLog *log.Logger
	config    *config.Config
}

func main() {
	addr := os.Getenv("ADDR")
	dsn := os.Getenv("DSN")
	jwt := os.Getenv("JWT_SECRET")
	minioUrl := os.Getenv("MINIO_ENDPOINT")
	minioBucketName := os.Getenv("MINIO_BUCKET")
	minioAccesKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	frontEndHost := os.Getenv("FRONTEND_HOST")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	accessLog := log.New(os.Stdout, "ACCESS\t", log.Ldate|log.Ltime)

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed // before the main() function exits.
	defer db.Close()

	minioClient, err := minio.New(minioUrl, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccesKey, minioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new instance of our application struct, containing the dependencies.
	config := &config.Config{
		JWTSecret:       []byte(jwt),
		DB:              db,
		MinioClient:     minioClient,
		MinioBucketName: minioBucketName,
		FrontEndHost:    frontEndHost,
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
