package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type RequestPayloadFormat struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "You've hit the Broker!",
	}

	_ = app.WriteJson(w, http.StatusOK, payload)

}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPayloadFormat

	// read the incoming data
	err := app.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)

	default:
		app.ErrorJSON(w, errors.New("unknown action"))

	}

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// Create some json data to send to auth microservice
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		app.ErrorJSON(w, errors.New("failed to auth payload"))
		return
	}

	// context to cancel the request
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// call the service
	request, err := http.NewRequestWithContext(ctx, "POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	// Use http connection pooling // this method inits the client pool when server startup
	// By this method , when ever this authenticate() func is called
	// we do not need new connections to handle the request
	response, err := app.httpClient.Do(request)

	if err != nil {
		slog.Error("failed to create request to auth service", "error", err)

		app.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// Log the responses
	slog.Info("auth service response received",
		slog.String("status", response.Status),
		slog.String("user", a.Email),
	)

	// get the desried status code
	if response.StatusCode == http.StatusUnauthorized {
		slog.Info("authentication failed", slog.String("user", a.Email))

		app.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	} else if response.StatusCode != http.StatusAccepted {
		slog.Error("unexpected status from auth service",
			slog.Int("status_code", response.StatusCode))
		app.ErrorJSON(w, errors.New("error calling the auth service"))
		return
	}

	// Read from the auth service response and send back to user via Broker
	var jsonFromAuth jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromAuth)
	if err != nil {
		slog.Error("failed to decode response from auth service", slog.Any("error", err))

		app.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	slog.Info("authentication successful", slog.String("user", a.Email))

	// send the response back to user in json

	var payload jsonResponse

	payload.Error = false
	payload.Message = "authenticated successfully"
	payload.Data = jsonFromAuth.Data

	app.WriteJson(w, http.StatusAccepted, payload)
}

// LogItem functions calls the logger microservice and logs the items
func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.ErrorJSON(w, err)
		return
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "logged"

	app.WriteJson(w, http.StatusAccepted, payload)

}
