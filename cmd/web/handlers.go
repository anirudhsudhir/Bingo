package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/anirudhsudhir/Bingo/internal/models"
	"github.com/go-chi/chi/v5"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	rows, err := app.snipModel.GetLatestSnips()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := newTemplateData()
	data.Snips = rows
	app.renderTemplate(w, "home.html", data)
}

func (app *application) viewSnip(w http.ResponseWriter, r *http.Request) {
	idr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idr)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snip, err := app.snipModel.ReadSnip(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		} else {
			app.serverError(w, err)
			return
		}
	}

	data := newTemplateData()
	data.Snip = snip
	app.renderTemplate(w, "view.html", data)
}

func (app *application) createSnip(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Will display html form soon"))
}

func (app *application) createSnipPost(w http.ResponseWriter, r *http.Request) {
	title := "Test snip 1"
	content := "content of test snip"
	expires := 7

	id, err := app.snipModel.InsertSnip(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snip/view/%d", id), http.StatusSeeOther)
}
