package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

// Config
type config struct {
	port int
	env  string
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

// Serve
func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Println((fmt.Sprintf("Starting Backend API Server in %s mode on port %d", app.config.env, app.config.port)))

	return srv.ListenAndServe()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var cfg config

	// flags for command line arguments
	// i.e. go run main.go -port=8000 -env=production -api=http://example.com/api
	// if not provided defaults will be chosen from the ones provided here
	flag.IntVar(&cfg.port, "port", 6060, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production|maintenance}")

	flag.Parse()

	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	fmt.Println("stripe key from env:", cfg.stripe.key)
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")
	fmt.Println("stripe secret from env:", cfg.stripe.secret)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	version := "1.0"

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
	}

	err = app.serve()

	if err != nil {
		log.Fatal(err)
	}
}
