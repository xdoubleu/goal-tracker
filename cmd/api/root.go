package main

import "net/http"

func (app *Application) rootRoute(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /",
		app.rootHandler,
	)
}

func (app *Application) rootHandler(w http.ResponseWriter, r *http.Request) {
	// use golang html template
}
