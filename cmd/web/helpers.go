package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

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

func (app *application) renderTemplate(w http.ResponseWriter, page string, templateData *templateData) {
	ts, found := app.templateCache[page]
	if !found {
		err := fmt.Errorf("no template present for %s page", page)
		app.serverError(w, err)
		return
	}

	buf := &bytes.Buffer{}

	err := ts.ExecuteTemplate(buf, "base", templateData)
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

func newTemplateData() *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}
