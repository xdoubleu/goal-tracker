package main

import (
	"goal-tracker/api/internal/tplhelper"
	"net/http"
)

func (app *Application) rootRoute(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /",
		app.authTemplateAccess(app.rootHandler),
	)
}

func (app *Application) rootHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: if user has no todoist coupled, ask them to sign in
	//TODO: if user has no steam coupled ask them to sign in
	//TODO: if user has no goodreads coupled, ask them for their goodreads UR

	tplhelper.RenderWithPanic(app.tpl, w, "root.html", nil)
}
