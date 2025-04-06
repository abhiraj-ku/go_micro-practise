package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "You've hit the Broker!",
	}

	_ = app.WriteJson(w, http.StatusOK, payload)

}
