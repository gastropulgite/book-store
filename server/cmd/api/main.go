package main

import (
	"book-api/internal/driver"
	"fmt"
	"log"
	"net/http"
	"os"
)

type config struct {
	port int
}

type application struct {
	config config
	infoLog *log.Logger
	errorLog *log.Logger
	db *driver.DB
}

type Payload struct {
	Okay  bool `json:"okay"`
	Message string `json:"message"`
}

func main() {
	var config config
	config.port = 8081

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate | log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile)

	dsn := os.Getenv("DSN")

	db, err := driver.ConnectPostgres(dsn)

	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	defer db.SQL.Close()
	
	app := &application{
		config: config,
		infoLog: infoLog,
		errorLog: errorLog,
		db: db,
	}

	err = app.serve()

	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) serve() error{
	app.infoLog.Println("API's listening on port", app.config.port)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}
	return srv.ListenAndServe()
}