package main

import (
	"net/http"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
	"github.com/XDoubleU/essentia/pkg/context"
	"github.com/XDoubleU/essentia/pkg/parse"

	"goal-tracker/api/internal/constants"
	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/models"
)

func (app *Application) goalsRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /goals/states",
		app.authAccess(app.getStatesHandler),
	)
	mux.HandleFunc(
		"GET /goals/sources",
		app.authAccess(app.getSourcesHandler),
	)
	mux.HandleFunc(
		"GET /goals/types",
		app.authAccess(app.getTypesHandler),
	)
	mux.HandleFunc(
		"GET /goals",
		app.authAccess(app.getPagedGoalsHandler),
	)
	mux.HandleFunc(
		"POST /goals",
		app.authAccess(app.createGoalHandler),
	)
	mux.HandleFunc(
		"PATCH /goals/{id}",
		app.authAccess(app.updateGoalHandler),
	)
	mux.HandleFunc(
		"DELETE /goals/{id}",
		app.authAccess(app.deleteGoalHandler),
	)
}

// @Summary	Get all states
// @Tags		goals
// @Success	200		{object}	[]State
// @Failure	400		{object}	ErrorDto
// @Failure	401		{object}	ErrorDto
// @Failure	500		{object}	ErrorDto
// @Router		/goals/states [get].
func (app *Application) getStatesHandler(w http.ResponseWriter,
	r *http.Request) {
	err := httptools.WriteJSON(w, http.StatusOK, models.States, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Get all sources
// @Tags		goals
// @Success	200		{object}	[]Source
// @Failure	400		{object}	ErrorDto
// @Failure	401		{object}	ErrorDto
// @Failure	500		{object}	ErrorDto
// @Router		/goals/sources [get].
func (app *Application) getSourcesHandler(w http.ResponseWriter,
	r *http.Request) {

	err := httptools.WriteJSON(w, http.StatusOK, models.Sources, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Get all types
// @Tags		goals
// @Success	200		{object}	[]Type
// @Failure	400		{object}	ErrorDto
// @Failure	401		{object}	ErrorDto
// @Failure	500		{object}	ErrorDto
// @Router		/goals/types [get].
func (app *Application) getTypesHandler(w http.ResponseWriter,
	r *http.Request) {

	err := httptools.WriteJSON(w, http.StatusOK, models.Types, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Get all goals paged
// @Tags		goals
// @Param		page	query		int	false	"Page to fetch"
// @Success	200		{object}	[]Goal
// @Failure	400		{object}	ErrorDto
// @Failure	401		{object}	ErrorDto
// @Failure	500		{object}	ErrorDto
// @Router		/goals [get].
func (app *Application) getPagedGoalsHandler(w http.ResponseWriter,
	r *http.Request) {
	var pageSize int = 4

	last, err := parse.QueryParam[*string](r, "last", nil, nil)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	page, err := app.services.Goals.GetPage(r.Context(), *user, pageSize, last)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, page, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Create goal
// @Tags		goals
// @Param		createGoalDto	body		CreateGoalDto	true	"CreateGoalDto"
// @Success	201			{object}	Goal
// @Failure	400			{object}	ErrorDto
// @Failure	401			{object}	ErrorDto
// @Failure	409			{object}	ErrorDto
// @Failure	500			{object}	ErrorDto
// @Router		/goals [post].
func (app *Application) createGoalHandler(w http.ResponseWriter, r *http.Request) {
	var createGoalDto *dtos.CreateGoalDto

	err := httptools.ReadJSON(r.Body, &createGoalDto)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	school, err := app.services.Goals.Create(r.Context(), *user, createGoalDto)
	if err != nil {
		httptools.HandleError(w, r, err, createGoalDto.ValidationErrors)
		return
	}

	err = httptools.WriteJSON(w, http.StatusCreated, school, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Update goal
// @Tags		goals
// @Param		id			path		string			true	"Goal ID"
// @Param		updateGoalDto	body		UpdateGoalDto	true	"UpdateGoalDto"
// @Success	200			{object}	Goal
// @Failure	400			{object}	ErrorDto
// @Failure	401			{object}	ErrorDto
// @Failure	409			{object}	ErrorDto
// @Failure	500			{object}	ErrorDto
// @Router		/goals/{id} [patch].
func (app *Application) updateGoalHandler(w http.ResponseWriter, r *http.Request) {
	var updateGoalDto *dtos.UpdateGoalDto

	id, err := parse.URLParam(r, "id", parse.UUID)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	err = httptools.ReadJSON(r.Body, &updateGoalDto)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	school, err := app.services.Goals.Update(r.Context(), *user, id, updateGoalDto)
	if err != nil {
		httptools.HandleError(w, r, err, updateGoalDto.ValidationErrors)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, school, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Delete goal
// @Tags		goals
// @Param		id	path		string	true	"Goal ID"
// @Success	200	{object}	Goal
// @Failure	400	{object}	ErrorDto
// @Failure	401	{object}	ErrorDto
// @Failure	404	{object}	ErrorDto
// @Failure	500	{object}	ErrorDto
// @Router		/goals/{id} [delete].
func (app *Application) deleteGoalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parse.URLParam(r, "id", parse.UUID)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	goal, err := app.services.Goals.Delete(r.Context(), *user, id)
	if err != nil {
		httptools.HandleError(w, r, err, nil)
	}

	err = httptools.WriteJSON(w, http.StatusOK, goal, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}
