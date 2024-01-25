package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/anirudhsudhir/Bingo/internal/validators"
)

type FormData struct {
	Title   string
	Content string
	Expires int
	validators.Validator
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLogger.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	app.notFound(w)
}

func (app *application) renderTemplate(w http.ResponseWriter, page string, status int, templateData *templateData) {
	ts, found := app.templateCache[page]
	if !found {
		err := fmt.Errorf("no template present for %s page", page)
		app.serverError(w, err)
		return
	}

	buf := &bytes.Buffer{}
	w.WriteHeader(status)
	err := ts.ExecuteTemplate(buf, "base", templateData)
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

func newTemplateData() *templateData {
	return &templateData{}
}
