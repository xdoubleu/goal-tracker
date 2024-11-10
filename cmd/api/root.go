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
	tasks, err := app.services.Todoist.GetTasksFromProjectGroupedBySection(app.config.TodoistProjectID)
	if err != nil {
		panic(err)
	}

	tplhelper.RenderWithPanic(app.tpl, w, "root.html", tasks)
}
