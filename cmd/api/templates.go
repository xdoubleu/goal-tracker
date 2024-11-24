package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/XDoubleU/essentia/pkg/context"
	"github.com/XDoubleU/essentia/pkg/parse"

	"goal-tracker/api/internal/constants"
	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/tplhelper"
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
	mux.HandleFunc(
		"GET /{id}",
		app.authTemplateAccess(app.graphHandler),
	)
}

func (app *Application) rootHandler(w http.ResponseWriter, r *http.Request) {
	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		panic(errors.New("not signed in"))
	}

	// note: todoist will only use 4 indent levels
	// (0: parent, 1: sub, 2: 2*sub, 3: 3*sub, 4: 4*sub)
	goals, err := app.services.Goals.GetAllGroupedByStateAndParentGoal(
		r.Context(),
		user.ID,
	)
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

type GoalAndProgress struct {
	Goal           models.Goal
	ProgressLabels []string
	ProgressValues []int64
}

func (app *Application) graphHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parse.URLParam[string](r, "id", nil)
	if err != nil {
		panic(err)
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		panic(errors.New("not signed in"))
	}

	goal, err := app.services.Goals.GetByID(r.Context(), id, *user)
	if err != nil {
		panic(err)
	}

	//nolint:godox //i'm aware
	//TODO: fetch progress
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	yesterdaySquared := yesterday.Add(-24 * time.Hour)
	format := "2006-01-02"

	progressLabels := []string{
		yesterdaySquared.Format(format),
		yesterday.Format(format),
		now.Format(format),
	}

	progressValues := []int64{
		10,
		5,
		10,
	}

	goalAndProgress := GoalAndProgress{
		Goal:           *goal,
		ProgressLabels: progressLabels,
		ProgressValues: progressValues,
	}

	tplhelper.RenderWithPanic(app.tpl, w, "graph.html", goalAndProgress)
}
