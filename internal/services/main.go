package services

import (
	"github.com/supabase-community/gotrue-go"

	"goal-tracker/api/internal/config"
	"goal-tracker/api/internal/repositories"
	"goal-tracker/api/pkg/todoist"
)

type Services struct {
	Auth    AuthService
	Goals   GoalService
	Todoist TodoistService
}

func New(
	config config.Config,
	repositories repositories.Repositories,
	supabaseClient gotrue.Client,
	todoistClient todoist.Client,
) Services {
	auth := AuthService{client: supabaseClient}
	todoist := TodoistService{client: todoistClient, projectID: config.TodoistProjectID}
	goals := GoalService{
		goals:    repositories.Goals,
		progress: repositories.Progress,
		todoist:  todoist,
	}

	return Services{
		Auth:    auth,
		Goals:   goals,
		Todoist: todoist,
	}
}
