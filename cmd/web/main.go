package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snip/view", app.viewSnip)
	mux.HandleFunc("/snip/create", app.createSnip)

	server := &http.Server{
		Addr:     *addr,
		Handler:  mux,
		ErrorLog: errorLogger,
	}

	infoLogger.Printf("Listening on port %s", *addr)
	err := server.ListenAndServe()
	if err != nil {
		errorLogger.Fatal(err)
	}
}
