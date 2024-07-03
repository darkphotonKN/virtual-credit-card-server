package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Config
type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

// Application
type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
}

func main() {
	var cfg config

	// flags for command line arguments
	// i.e. go run main.go -port=8000 -env=production -api=http://example.com/api
	// if not provided defaults will be chosen from the ones provided here
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production|maintenance}")

	flag.Parse()

	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	fmt.Println("stripe key from env:", cfg.stripe.key)
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application
}
