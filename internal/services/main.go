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
	goals := GoalService{goals: repositories.Goals}
	todoist := TodoistService{client: todoistClient}

	return Services{
		Auth:    auth,
		Goals:   goals,
		Todoist: todoist,
	}
}
