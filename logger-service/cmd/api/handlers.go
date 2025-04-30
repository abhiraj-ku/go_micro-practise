package main

import (
	"net/http"

	"github.com/abhiraj-ku/go_micro-practise/data"
)

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var reqPay LogPayload
	_ = app.readJSON(w, r, &reqPay)

	event := data.LogEntry{
		Name: reqPay.Name,
		Data: reqPay.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}
	app.writeJSON(w, http.StatusAccepted, resp)

}
