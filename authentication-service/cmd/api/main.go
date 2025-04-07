package main

import (
	"auth_service/data"
	"database/sql"
	"fmt"
)

const webPort = "9090"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	fmt.Println("Hello, world!")
}
