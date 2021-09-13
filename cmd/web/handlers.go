package main

import (
	"net/http"
	"github.com/streadway/amqp"
)

func (app *Application) RootHandler(rw http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Root handler invoked")
}

func (app *Application) PublishEvent(rw http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Publishing an event")
	app.Channel.Publish(
		"logs",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte("{ data : Test Data }"),
		},
	)
	app.infoLog.Println("Event successfully published on the channel")
}