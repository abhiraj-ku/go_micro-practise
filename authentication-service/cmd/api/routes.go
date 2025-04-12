package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (m *Config) Routes() http.Handler {
	mux := chi.NewRouter()

	// User cors to allow * everyone as of Now
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://* "},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("ping"))

	return mux
}
