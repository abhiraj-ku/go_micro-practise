package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EnquiryMailPayloadUsingSendgrid struct {
	To               string    `json:"to"`
	ToName           string    `json:"to_name"`
	Subject          string    `json:"subject"`
	PropertyName     string    `json:"name"`
	PropertyLocation string    `json:"location"`
	Timestamp        time.Time `json:"timestamp"`
}
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

func (app *Config) SendMailViaSendGrid(w http.ResponseWriter, r *http.Request) {
	var payload EnquiryMailPayloadUsingSendgrid
	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = SendEmailWithSendGrid(payload.To, payload.ToName, payload.Subject, payload.PropertyName, payload.PropertyLocation)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	log.Println("sent email using sendgrid:[DEBUG:SendMailViaSendGrid]")

	_payload := jsonResponse{
		Error:   false,
		Message: "sent email using sendgrid to " + payload.To,
	}

	app.writeJSON(w, http.StatusAccepted, _payload)
}

func SendEmailWithSendGrid(toEmail, toName, mailSubject, propertyName, propertyLocation string) error {
	from := mail.NewEmail("Enquiry Manager", "anand.japan896@icloud.com")
	subject := mailSubject
	to := mail.NewEmail(toName, toEmail)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := fmt.Sprintf(`
		<html>
		  <body style="font-family: Arial, sans-serif; padding: 20px;">
			<p>Thank you for your enquiry for the property <strong>%s</strong> at <strong>%s</strong>.</p>
			<div style="margin-top: 20px;">
			  <a href="http://localhost/make-reservation" style="display: inline-block; padding: 10px 20px; margin-right: 10px; background-color: black; color: white; text-decoration: none; border-radius: 5px;">I want to Visit property</a>
			  <a href="http://localhost/request-callback" style="display: inline-block; padding: 10px 20px; background-color: black; color: white; text-decoration: none; border-radius: 5px;">I want a call again</a>
			</div>
		  </body>
		</html>
	`, propertyName, propertyLocation)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API"))
	response, err := client.Send(message)
	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Println("Email sent successfully. Status Code:", response.StatusCode)
	log.Println("Response Body:", response.Body)
	log.Println("Response Headers:", response.Headers)
	return nil
}
