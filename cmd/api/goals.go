package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/XDoubleU/essentia/pkg/context"
	"github.com/XDoubleU/essentia/pkg/parse"

	"goal-tracker/api/internal/constants"
	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/temptools"
)

func (app *Application) goalsRoutes(prefix string, mux *http.ServeMux) {
	mux.HandleFunc(
		fmt.Sprintf("POST %s/goals/{id}/link", prefix),
		app.authAccess(app.linkGoalHandler),
	)
	mux.HandleFunc(
		fmt.Sprintf("GET %s/goals/{id}/unlink", prefix),
		app.authAccess(app.unlinkGoalHandler),
	)
}

func (app *Application) linkGoalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parse.URLParam[string](r, "id", nil)
	if err != nil {
		panic(err)
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		panic(errors.New("not signed in"))
	}

	var linkGoalDto dtos.LinkGoalDto

	err = temptools.ReadForm(r, &linkGoalDto)
	if err != nil {
		temptools.RedirectWithError(w, r, fmt.Sprintf("/link/%s", id), err)
		return
	}

	err = app.services.Goals.LinkGoal(r.Context(), id, user.ID, &linkGoalDto)
	if err != nil {
		temptools.RedirectWithError(w, r, fmt.Sprintf("/link/%s", id), err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/goals/%s", id), http.StatusSeeOther)
}

func (app *Application) unlinkGoalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parse.URLParam[string](r, "id", nil)
	if err != nil {
		panic(err)
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		panic(errors.New("not signed in"))
	}

	err = app.services.Goals.UnlinkGoal(r.Context(), id, user.ID)
	if err != nil {
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
