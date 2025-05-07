package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type UserInputs struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var uInput UserInputs

	// read the json data from body
	err := app.ReadJSON(w, r, &uInput)
	if err != nil {
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	user, err := app.m.Users.GetByEmail(uInput.Email)
	if err != nil {
		app.ErrorJSON(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}

	// match the password
	isPassValid, err := user.PasswordMatches(uInput.Password)

	if err != nil || !isPassValid {
		app.ErrorJSON(w, errors.New("invalid password"), http.StatusBadRequest)
		return
	}

	// log authentication to check when someone log's in
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}

	app.WriteJson(w, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil

}
