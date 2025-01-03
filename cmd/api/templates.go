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
	)
	if err != nil {
		panic(err)
	}

	tplhelper.RenderWithPanic(app.tpl, w, "root.html", goals)
}

type LinkTemplateData struct {
	Goal    models.Goal
	Sources []models.Source
	Tags    []string
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

	goal, err := app.services.Goals.GetGoalByID(r.Context(), id)
	if err != nil {
		panic(err)
	}

	tags, err := app.services.Goodreads.GetAllTags(r.Context())
	if err != nil {
		panic(err)
	}

	goalAndSources := LinkTemplateData{
		Goal:    *goal,
		Sources: models.Sources,
		Tags:    tags,
	}
	tplhelper.RenderWithPanic(app.tpl, w, "link.html", goalAndSources)
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

	goal, err := app.services.Goals.GetGoalByID(r.Context(), id)
	if err != nil {
		panic(err)
	}

	viewType := models.Types[*goal.TypeID].ViewType
	switch viewType {
	case models.Graph:
		app.graphViewProgress(w, r, goal)
	case models.List:
		app.listViewProgress(w, r, goal)
	}
}

type GraphData struct {
	Goal           models.Goal
	ProgressLabels []string
	ProgressValues []string
}

func (app *Application) graphViewProgress(
	w http.ResponseWriter,
	r *http.Request,
	goal *models.Goal,
) {
	//nolint:godox //I know
	// TODO make this based on duration of goal

	var dateStart time.Time
	var dateEnd time.Time

	switch *goal.TypeID {
	case models.SteamCompletionRate.ID:
		// only get the past year
		dateStart = time.Date(
			time.Now().Year()-1,
			time.Now().Month(),
			time.Now().Day(),
			0,
			0,
			0,
			0,
			time.UTC,
		)
		dateEnd = time.Date(
			time.Now().Year(),
			time.Now().Month(),
			time.Now().Day(),
			0,
			0,
			0,
			0,
			time.UTC,
		)
	case models.FinishedBooksThisYear.ID:
		// only get this year
		dateStart = time.Date(
			time.Now().Year(),
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		)
		dateEnd = time.Date(
			time.Now().Year(),
			12,
			31,
			0,
			0,
			0,
			0,
			time.UTC,
		)
	default:
		panic("not implemented")
	}

	progressLabels, progressValues, err := app.services.Goals.GetProgressByTypeIDAndDates(
		r.Context(),
		*goal.TypeID,
		dateStart,
		dateEnd,
	)
	if err != nil {
		panic(err)
	}

	graphData := GraphData{
		Goal:           *goal,
		ProgressLabels: progressLabels,
		ProgressValues: progressValues,
	}

	tplhelper.RenderWithPanic(app.tpl, w, "graph.html", graphData)
}

type ListData struct {
	Goal      models.Goal
	ListItems []models.ListItem
}

func (app *Application) listViewProgress(
	w http.ResponseWriter,
	r *http.Request,
	goal *models.Goal,
) {
	listItems, err := app.services.Goals.GetListItemsByGoalID(r.Context(), goal.ID)
	if err != nil {
		panic(err)
	}

	listData := ListData{
		Goal:      *goal,
		ListItems: listItems,
	}
	tplhelper.RenderWithPanic(app.tpl, w, "list.html", listData)
}
