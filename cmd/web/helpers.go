package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/anirudhsudhir/Bingo/internal/validators"
)

type FormData struct {
	Title                string `schema:"title"`
	Content              string `schema:"content"`
	Expires              int    `schema:"expires"`
	validators.Validator `schema:"-"`
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

func (app *application) parseForm(dst any, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		return err
	}
	return nil
}

func generateSessionKey(length int) (key string, err error) {
	randomBytes := make([]byte, length)
	_, err = rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	key = base64.StdEncoding.EncodeToString(randomBytes)
	return key, nil
}
