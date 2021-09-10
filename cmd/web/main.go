package main 

import (
	"log"
	"flag"
	"html/template"
)

const version     = "1.0.0"
const cssVersion  = "1"

type config struct {
	port int         // to be read from the command line flag 
	env  string      // to be read from the command line flag 
	api  string      // to be read from the command line flag 
	db   struct {
		dbName string 
	}
	stripe struct {
		secret  string 
		key     string 
	}
}

type application struct {
	config    config 
	infoLog   *log.Logger
	errorLog  *log.Logger 
	templateCache map[string]*template.Template
	version   string 
}

func main() {
	var cfg config 
	// Defining command line flags to (used as --port <value>)
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production}")
	flag.StringVar(&cfg.api, "api", "http://localhost:4000", "URL to api")

	// Parse parses the command-line flags from os.Args[1:]
	flag.Parse()

}