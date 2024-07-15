package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/darkphotonKN/virtual-credit-card-server/internal/driver"
	"github.com/darkphotonKN/virtual-credit-card-server/internal/models"
	"github.com/joho/godotenv"
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
	DB       models.DBModel
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
	// database information
	flag.StringVar(&cfg.db.dsn, "dsn", "root:123456@tcp(localhost:3307)/virtual_terminal_db?parseTime=true&tls=false", "DSN")

	flag.Parse()

	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	fmt.Println("stripe key from env:", cfg.stripe.key)
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")
	fmt.Println("stripe secret from env:", cfg.stripe.secret)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	version := "1.0"

	// Connecting to DB with gorm
	db, err := driver.OpenDB(cfg.db.dsn)

	if err != nil {
		log.Fatal("DB could not be connected to.")
	}

	fmt.Println("DB connected.")

	// auto migration for tables
	err = db.AutoMigrate(&models.Product{})

	if err != nil {
		log.Fatalf("Could not initialize DB table products.")
	}

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
		DB:       models.DBModel{DB: db},
	}

	err = app.serve()

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Product{})
}
