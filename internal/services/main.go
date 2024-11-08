package services

import (
	"context"
	"log/slog"

	"goal-tracker/api/internal/config"
	"goal-tracker/api/internal/repositories"

	"github.com/supabase-community/gotrue-go"
)

type Services struct {
	Auth  AuthService
	Goals GoalService
}

func New(
	ctx context.Context,
	logger *slog.Logger,
	config config.Config,
	repositories repositories.Repositories,
	client gotrue.Client,
) Services {
	auth := AuthService{client: client}
	goals := GoalService{goals: repositories.Goals}

	return Services{
		Auth:  auth,
		Goals: goals,
	}
}
