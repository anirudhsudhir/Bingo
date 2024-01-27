package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/anirudhsudhir/Bingo/internal/models"
	"github.com/anirudhsudhir/Bingo/internal/validators"
	"github.com/go-chi/chi/v5"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	rows, err := app.snipModel.GetLatestSnips()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data, err := app.newTemplateData(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data.Snips = rows
	app.renderTemplate(w, "home.html", http.StatusOK, data)
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

	data, err := app.newTemplateData(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data.Snip = snip
	app.renderTemplate(w, "view.html", http.StatusOK, data)
}

func (app *application) createSnip(w http.ResponseWriter, r *http.Request) {
	data := &templateData{Form: FormData{Expires: 365}}
	app.renderTemplate(w, "create.html", http.StatusOK, data)
}

func (app *application) createSnipPost(w http.ResponseWriter, r *http.Request) {
	formData := &FormData{}
	err := app.parseForm(formData, r)
	if err != nil {
		app.serverError(w, err)
		return
	}

	formData.ValidateElement(validators.NoContent(formData.Title), "title", "The title field cannot be empty")
	formData.ValidateElement(validators.NoContent(formData.Content), "content", "The content field cannot be empty")
	formData.ValidateElement(validators.MaxLen(formData.Title, 100), "title", "The title field cannot contain more than 100 characters")
	formData.ValidateElement(validators.AllowedValues(formData.Expires, 1, 7, 365), "expires", "Expiry duration must be 1, 7 or 365 days")

	if !formData.ValidForm() {
		data := &templateData{Form: formData}
		app.renderTemplate(w, "create.html", http.StatusUnprocessableEntity, data)
		return
	}

	id, err := app.snipModel.InsertSnip(formData.Title, formData.Content, formData.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	session, err := app.sessionStore.Get(r, "session")
	if err != nil {
		app.serverError(w, err)
		return
	}
	session.AddFlash("The snip has been added successfully!")
	err = session.Save(r, w)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snip/view/%d", id), http.StatusSeeOther)
}
