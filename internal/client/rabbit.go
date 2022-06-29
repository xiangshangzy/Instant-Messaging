package client

import (
	"fmt"
	"github.com/streadway/amqp"
)

type Mail struct {
	To      string
	Tittle  string
	Content string
	MsgId   string
}

func mq() {
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672")
	if err != nil {

	}
	channel, err := conn.Channel()
	delivery, err := channel.Consume(
		"queueName", // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	for msg := range delivery {
		fmt.Printf("d.Body: %v\n", string(msg.Body))
		msg.Ack(false)
	}

}

func rec() {
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672")
	if err != nil {
	}
	channel, err := conn.Channel()
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("body"),
	}
	channel.Publish("email.exchange", "email.key", false, false, msg)
}
