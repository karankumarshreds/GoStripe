package main 

import (
	"os"
	"log"
	"flag"
	"net/http"
	"database/sql"
	"html/template"
	"github.com/gorilla/mux"
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

	app.InitializeRoutes()
}

func (app *Application) Run() {
	err := http.ListenAndServe(":" + (app.config.port), app.Router)
	if err != nil {
		app.errorLog.Fatalf("Error while listening on port %v \n %v", app.config.port, err)
	}
}