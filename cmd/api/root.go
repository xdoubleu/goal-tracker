package main

import (
	"errors"
	"goal-tracker/api/internal/constants"
	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/tplhelper"
	"net/http"

	"github.com/XDoubleU/essentia/pkg/context"
)

func (app *Application) templateRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /",
		app.authTemplateAccess(app.rootHandler),
	)
	mux.HandleFunc(
		"GET /favicon.ico",
		app.serveFaviconHandler,
	)
}

func (app *Application) rootHandler(w http.ResponseWriter, r *http.Request) {
	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		panic(errors.New("not signed in"))
	}

	goals, err := app.services.Goals.GetAllGroupedByState(r.Context(), user.ID)
	if err != nil {
		panic(err)
	}

	tplhelper.RenderWithPanic(app.tpl, w, "root.html", goals)
}

func (app *Application) serveFaviconHandler(w http.ResponseWriter, r *http.Request) {
	//TODO
}
