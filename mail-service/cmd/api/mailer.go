package main

type Mail struct {
	Domain     string
	Host       string
	Port       string
	Username   string
	Password   string
	Encryption string
	From       string
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

func (m *Mail) SendMessageSMTP(msg Message) {

}
