package main 

import (
	"os"
	"log"
	"flag"
	"net/http"
	"database/sql"
	"html/template"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

const version = "1.0.0"

type Application struct {

	config struct {
		port string      // to be read from the command line flag 
		env string       // to be read from the command line flag 
		api string       // to be read from the command line flag 

		stripe struct {
			secret string  // to be read from the environment variable 
			key string     // to be read from the environment variable
		}

	} 

	infoLog   *log.Logger
	errorLog  *log.Logger 
	templateCache map[string]*template.Template
	version   string 

	DB *sql.DB
	Router *mux.Router
	Channel *amqp.Channel

}

func (app *Application) Initialize(STRIPE_KEY string, STRIPE_SECRET string, DB_NAME string) {
	cfg := app.config
	// Defining command line flags to (used as --port <value>)
	flag.StringVar(&cfg.port, "port", "8000", "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production}")
	flag.StringVar(&cfg.api, "api", "http://localhost:4000", "URL to api")
	// Parse parses the command-line flags from os.Args[1:]
	flag.Parse()

	cfg.stripe.key     = STRIPE_KEY
	cfg.stripe.secret  = STRIPE_SECRET

	iLog  := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	eLog  := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	tc    := make(map[string]*template.Template)

	
	app.config         = cfg 
	app.infoLog        = iLog
	app.errorLog       = eLog
	app.Router         = mux.NewRouter()
	app.templateCache  = tc
	app.version        = version
	
	/* RABBIT MQ SETUP */
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	app.logForError("Error while connecting to RabbitMQ", err)
	
	ch, err := conn.Channel()
	app.logForError("Failed to open a channel", err)
	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	app.logForError("Failed to declare an exchange", err)

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	app.logForError("Failed to declare queue", err)

	err = ch.QueueBind(
		q.Name, // queue name
		"key123",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	app.logForError("Failed to bind to queue", err)

	app.Channel = ch
	
	app.InitializeRoutes()
}

func (app *Application) Run() {
	app.infoLog.Printf("Starting in %v mode on port %v", app.config.env , app.config.port)
	err := http.ListenAndServe(":" + (app.config.port), app.Router)
	if err != nil {
		app.errorLog.Fatalf("Error while listening on port %v \n %v", app.config.port, err)
	}
}

func (a * Application) logForError(message string, err error) {
	if err != nil {
		a.errorLog.Printf("%s: %s", message, err)
	}
}