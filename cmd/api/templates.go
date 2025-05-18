package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/XDoubleU/essentia/pkg/context"
	"github.com/XDoubleU/essentia/pkg/parse"
	tpltools "github.com/XDoubleU/essentia/pkg/tpl"

	"goal-tracker/api/internal/constants"
	"goal-tracker/api/internal/models"
)

func (app *Application) templateRoutes(mux *http.ServeMux) {
	mux.Handle(
		"GET /images/",
		http.FileServerFS(app.images),
	)
	mux.HandleFunc(
		"GET /{$}",
		app.authTemplateAccess(app.rootHandler),
	)
	mux.HandleFunc(
		"GET /edit/{id}",
		app.authTemplateAccess(app.editHandler),
	)
	mux.HandleFunc(
		"GET /goals/{id}",
		app.authTemplateAccess(app.goalProgressHandler),
	)
}

func (app *Application) rootHandler(w http.ResponseWriter, r *http.Request) {
	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		panic(errors.New("not signed in"))
	}

	goals, err := app.services.Goals.GetAllGoalsGroupedByStateAndParentGoal(
		r.Context(),
		user.ID,
	)
	if err != nil {
		panic(err)
	}

	tpltools.RenderWithPanic(app.tpl, w, "root.html", goals)
}

type LinkTemplateData struct {
	Goal    models.Goal
	Sources []models.Source
	Tags    []string
}

func (app *Application) editHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := parse.URLParam[string](r, "id", nil)
	if err != nil {
		panic(err)
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		panic(errors.New("not signed in"))
	}

	goal, err := app.services.Goals.GetGoalByID(r.Context(), id, user.ID)
	if err != nil {
		panic(err)
	}

	tags, err := app.services.Goodreads.GetAllTags(r.Context(), user.ID)
	if err != nil {
		panic(err)
	}

	goalAndSources := LinkTemplateData{
		Goal:    *goal,
		Sources: models.Sources,
		Tags:    tags,
	}

	tpltools.RenderWithPanic(app.tpl, w, "edit.html", goalAndSources)
}

func (app *Application) goalProgressHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parse.URLParam[string](r, "id", nil)
	if err != nil {
		panic(err)
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		panic(errors.New("not signed in"))
	}

	goal, err := app.services.Goals.GetGoalByID(r.Context(), id, user.ID)
	if err != nil {
		panic(err)
	}

	viewType := models.Types[*goal.TypeID].ViewType
	switch viewType {
	case models.Graph:
		app.graphViewProgress(w, r, goal, user.ID)
	case models.List:
		app.listViewProgress(w, r, goal, user.ID)
	}
}

type GraphData struct {
	Goal                 models.Goal
	DateLabels           []string
	ProgressValues       []string
	TargetValues         []string
	CurrentProgressValue string
	CurrentTargetValue   string
	Details              []models.ListItem
}

func (app *Application) graphViewProgress(
	w http.ResponseWriter,
	r *http.Request,
	goal *models.Goal,
	userID string,
) {
	progressLabels, progressValues, err := app.services.Goals.GetProgressByTypeIDAndDates(
		r.Context(),
		*goal.TypeID,
		userID,
		goal.PeriodStart(),
		goal.PeriodEnd(),
	)
	if err != nil {
		panic(err)
	}

	details, err := app.services.Goals.GetListItemsByGoal(r.Context(), goal, userID)
	if err != nil {
		panic(err)
	}

	graphData := GraphData{
		Goal:           *goal,
		DateLabels:     progressLabels,
		ProgressValues: progressValues,
		Details:        details,
	}

	if len(progressValues) > 0 {
		startProgress, _ := strconv.ParseFloat(progressValues[0], 64)
		graphData.TargetValues = goal.AdaptiveTargetValues(int(startProgress))
		graphData.CurrentProgressValue = progressValues[len(progressValues)-1]
		graphData.CurrentTargetValue = graphData.TargetValues[len(progressValues)-1]
	}

	tpltools.RenderWithPanic(app.tpl, w, "graph.html", graphData)
}

type ListData struct {
	Goal      models.Goal
	ListItems []models.ListItem
}

func (app *Application) listViewProgress(
	w http.ResponseWriter,
	r *http.Request,
	goal *models.Goal,
	userID string,
) {
	listItems, err := app.services.Goals.GetListItemsByGoal(
		r.Context(),
		goal,
		userID,
	)
	if err != nil {
		panic(err)
	}

	listData := ListData{
		Goal:      *goal,
		ListItems: listItems,
	}

	tpltools.RenderWithPanic(app.tpl, w, "list.html", listData)
}
