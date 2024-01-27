package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/anirudhsudhir/Bingo/internal/models"
)

type templateData struct {
	Snip  *models.Snip
	Snips []*models.Snip
	Form  any
	Flash interface{}
}

func newTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{}

	files, err := filepath.Glob("./ui/html/*.html")
	if err != nil {
		return nil, err
	}

	humanDate := func(t time.Time) string {
		return t.Format("02 Jan 2006 at 15:04")
	}
	funcMap := template.FuncMap{
		"humanDate": humanDate,
	}

	for _, file := range files {
		name := filepath.Base(file)

		ts := template.New(name).Funcs(funcMap)
		ts, err := ts.ParseFiles("./ui/html/base.html", file)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}
		templateCache[name] = ts
	}
	return templateCache, nil
}

func (app *application) newTemplateData(w http.ResponseWriter, r *http.Request) (*templateData, error) {
	session, err := app.sessionStore.Get(r, "session")
	if err != nil {
		return nil, err
	}

	data := &templateData{}
	if flashes := session.Flashes(); len(flashes) > 0 {
		data.Flash = flashes[0]
	}
	session.Save(r, w)
	return data, nil
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
