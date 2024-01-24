package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/anirudhsudhir/Bingo/internal/models"
)

type templateData struct {
	Snip        *models.Snip
	Snips       []*models.Snip
	CurrentYear int
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
