package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
)

const (
	webPort = "80"
)

type Config struct {
	Mailer Mail
	logger *slog.Logger
}

func main() {
	app := Config{
		Mailer: createMail(),
	}

	slog.Info("starting mail service on port", fmt.Sprintf("mail-server"), webPort)

	server := http.Server{
		Addr:    ":" + webPort,
		Handler: app.Routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("mailer service failed to start", err)
	}

}

func createMail() Mail {
	port := os.Getenv("MAILER_PORT")
	m := Mail{
		Domain:     os.Getenv("MAIL_DOMAIN"),
		Host:       os.Getenv("MAIL_HOST"),
		Port:       port,
		Username:   os.Getenv("MAIL_USERNAME"),
		Password:   os.Getenv("MAIL_PASSWORD"),
		Encryption: os.Getenv("MAIL_ENCRYPTION"),
		From:       os.Getenv("MAIL_FROM"),
		FromName:   os.Getenv("FROM_ADDRESS"),
	}
	return m
}
