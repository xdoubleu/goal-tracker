package services

import (
	"context"
	"log/slog"

	"goal-tracker/api/internal/config"
	"goal-tracker/api/internal/repositories"
	"goal-tracker/api/pkg/todoist"

	"github.com/supabase-community/gotrue-go"
)

type Services struct {
	Auth    AuthService
	Goals   GoalService
	Todoist TodoistService
}

func New(
	ctx context.Context,
	logger *slog.Logger,
	config config.Config,
	repositories repositories.Repositories,
	supabaseClient gotrue.Client,
	todoistClient todoist.Client,
) Services {
	auth := AuthService{client: supabaseClient}
	todoist := TodoistService{client: todoistClient, projectID: config.TodoistProjectID}
	goals := GoalService{goals: repositories.Goals, todoist: todoist}

	return Services{
		Auth:    auth,
		Goals:   goals,
		Todoist: todoist,
	}
}
