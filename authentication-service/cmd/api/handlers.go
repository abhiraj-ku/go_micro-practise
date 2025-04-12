package main

import (
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

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}

	app.WriteJson(w, http.StatusAccepted, payload)
}
