package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (cfg *Config) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(data); err != nil {
		return err
	}

	// Also check if there is no extra data
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must contain only one json data")
	}

	return nil

}

func (cfg *Config) WriteJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

// error Json

func (cfg *Config) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return cfg.WriteJson(w, statusCode, payload)
}
