package main

import (
	"errors"
	"goal-tracker/api/internal/constants"
	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/tplhelper"
	"net/http"

	"github.com/XDoubleU/essentia/pkg/context"
	"github.com/XDoubleU/essentia/pkg/parse"
)

func (app *Application) templateRoutes(mux *http.ServeMux) {
	mux.Handle(
		"GET /images/",
		http.FileServerFS(app.images),
	)
	mux.HandleFunc(
		"GET /",
		app.authTemplateAccess(app.rootHandler),
	)
	mux.HandleFunc(
		"GET /link/{id}",
		app.authTemplateAccess(app.linkHandler),
	)
}

func (app *Application) rootHandler(w http.ResponseWriter, r *http.Request) {
	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		panic(errors.New("not signed in"))
	}

	// note: todoist will only use 4 indent levels (0: parent, 1: sub, 2: 2*sub, 3: 3*sub, 4: 4*sub)
	goals, err := app.services.Goals.GetAllGroupedByStateAndParentGoal(r.Context(), user.ID)
	if err != nil {
		panic(err)
	}

	tplhelper.RenderWithPanic(app.tpl, w, "root.html", goals)
}

type GoalAndSources struct {
	Goal    models.Goal
	Sources []models.Source
}

func (app *Application) linkHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parse.URLParam[string](r, "id", nil)
	if err != nil {
		panic(err)
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		panic(errors.New("not signed in"))
	}

	task, err := app.services.Todoist.GetTaskByID(r.Context(), id)
	if err != nil {
		panic(err)
	}

	goalAndSources := GoalAndSources{
		Goal:    models.NewGoalFromTask(*task, user.ID, ""),
		Sources: models.Sources,
	}
	tplhelper.RenderWithPanic(app.tpl, w, "link.html", goalAndSources)
}
