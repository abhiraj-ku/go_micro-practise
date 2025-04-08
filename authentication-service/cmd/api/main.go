package main

import (
	"auth_service/data"
	"database/sql"
	"log"
)

const webPort = "9090"

type Config struct {
	DB     *sql.DB
	Models data.UserModel
}

func main() {
	db, err := sql.Open("postgres", "")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// models := data.NewModels(db)
}
