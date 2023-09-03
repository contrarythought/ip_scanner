package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"ip_scanner/app"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Creds struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

func getCreds() (*Creds, error) {
	var creds Creds
	file, err := os.Open("creds.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(jsonData, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}

func main() {
	creds, err := getCreds()
	if err != nil {
		log.Fatal(err)
	}

	connStr := "postgres://" + creds.User + ":" + creds.Pass + "@localhost/ip_scanner?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	errCh := make(chan error)
	logfile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 644)
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()
	logger := log.New(logfile, "log", log.LstdFlags|log.Lshortfile)
	app.Scan(db, errCh, logger)
}
