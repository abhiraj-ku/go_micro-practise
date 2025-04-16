package main

import (
	"bytes"
	"html/template"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain     string
	Host       string
	Port       string
	Username   string
	Password   string
	Encryption string
	FromAddr   string
	FromName   string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]interface{}
}

func (m *Mail) SendMessageSMTP(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddr
	}
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]interface{}{
		"message": msg.Data,
	}

	msg.DataMap = data

	// send this data to template
	sendToTemplate, err := m.buildHTMLMessage(msg)
	if err != nil {
		return nil
	}

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	// set the impt configs for the smtp server
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpclient, err := server.Connect()
	if err != nil {
		return err
	}
	email := mail.NewMSG()
	email.SetFrom(msg.From)
	email.SetSubject(msg.Subject)
	email.AddTo(msg.To)
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, attachement := range msg.Attachments {
			email.AddAttachment(attachement)
		}
	}
	err = email.Send(smtpclient)
	if err != nil {
		return err
	}
	return nil
}

func getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSL
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}

// buildHTMLMessage
func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.plane.gohtml"

	t, err := template.New("email-plan").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}
	var tmpl bytes.Buffer
	if err := t.ExecuteTemplate(&tmpl, "body", msg.DataMap); err != nil {
		return "", nil
	}

	plainMessage := tmpl.String()
	return plainMessage, nil
}
