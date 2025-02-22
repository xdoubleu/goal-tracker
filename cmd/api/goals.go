package main

import (
	"errors"
	"fmt"
	"net/http"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
	"github.com/XDoubleU/essentia/pkg/context"
	"github.com/XDoubleU/essentia/pkg/parse"

	"goal-tracker/api/internal/constants"
	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/models"
)

func (app *Application) goalsRoutes(prefix string, mux *http.ServeMux) {
	mux.HandleFunc(
		fmt.Sprintf("POST %s/goals/{id}/edit", prefix),
		app.authAccess(app.editGoalHandler),
	)
	mux.HandleFunc(
		fmt.Sprintf("GET %s/goals/{id}/unlink", prefix),
		app.authAccess(app.unlinkGoalHandler),
	)
	mux.HandleFunc(
		fmt.Sprintf("GET %s/goals/{id}/complete", prefix),
		app.authAccess(app.completeGoalHandler),
	)
}

func (app *Application) editGoalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parse.URLParam[string](r, "id", nil)
	if err != nil {
		panic(err)
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		panic(errors.New("not signed in"))
	}

	var linkGoalDto dtos.LinkGoalDto

	err = httptools.ReadForm(r, &linkGoalDto)
	if err != nil {
		httptools.RedirectWithError(w, r, fmt.Sprintf("/edit/%s", id), err)
		return
	}

	if ok, errs := linkGoalDto.Validate(); !ok {
		httptools.FailedValidationResponse(w, r, errs)
		return
	}

	err = app.services.Goals.LinkGoal(r.Context(), id, user.ID, &linkGoalDto)
	if err != nil {
		httptools.RedirectWithError(w, r, fmt.Sprintf("/edit/%s", id), err)
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
		panic(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) completeGoalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parse.URLParam[string](r, "id", nil)
	if err != nil {
		panic(err)
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		panic(errors.New("not signed in"))
	}

	err = app.services.Goals.CompleteGoal(r.Context(), id, user.ID)
	if err != nil {
		panic(err)
	}

	err = app.services.Goals.ImportGoalsFromTodoist(r.Context(), user.ID)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
