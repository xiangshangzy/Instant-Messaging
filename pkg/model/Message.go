package model

type Message struct {
	Receiver string `json:"receiver"`
	Sender   string `json:"sender"`
	Body     string `json:"body"`
	Type     string `json:"type"`
	Id       string `json:"id"`
}
