package main

import (
	"auth_service/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const webPort = "9090"

type Config struct {
	m data.Models
}

var db *sql.DB
var noOfConn int32

func main() {

	// slog.Info("Auth service main func", fmt.Sprintf("Auth"))
	// Connect to db
	conn := connectDB()
	if conn != nil {
		log.Panic("Can't connect to db")
	}
	models := data.NewModels(db)
	// config:
	cfg := Config{
		m: models,
	}

	// models := data.NewModels(db)

	// Server setup

	// define the server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: cfg.Routes(),
	}

	// start the server
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// Db connection

func openConn(uri string) (*sql.DB, error) {
	db, err := sql.Open("pgx", uri)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectDB() *sql.DB {
	dsn := os.Getenv("DB_URI")

	for {
		connection, err := openConn(dsn)
		if err != nil {
			fmt.Println("DB not ready yet..")
			noOfConn++
		} else {
			log.Println("DB connected...")
			return connection
		}
		if noOfConn > 10 {
			fmt.Println(err)
			return nil
		}

		log.Println("Back of for 2 sec")
		time.Sleep(2 * time.Second)
		continue
	}
}
