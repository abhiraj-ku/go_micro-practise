package main

import "net/http"

type MailMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	var msgRequest MailMessage

	// read the json data from body
	err := app.readJSON(w, r, &msgRequest)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	message := MailMessage{
		From:    msgRequest.From,
		To:      msgRequest.To,
		Subject: msgRequest.Subject,
		Message: msgRequest.Message,
	}

	err = app.Mailer.SendMessageSMTP(message)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// send Response
	payload := jsonResponse{
		Error:   false,
		Message: "email sent to " + msgRequest.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}
