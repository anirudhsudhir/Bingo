package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", app.home)
	mux.Route("/snip", func(r chi.Router) {
		mux.Get("/snip/create", app.createSnip)
		mux.Post("/snip/create", app.createSnipPost)
		mux.Get("/snip/view/{id}", app.viewSnip)
	})
	mux.HandleFunc("/*", app.notFoundHandler)

	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Method(http.MethodGet, "/static/*", http.StripPrefix("/static", fileserver))

	return app.recoverPanic(app.logRequests(secureHeaders(mux)))
}
