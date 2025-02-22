package services

import (
	"log/slog"

	"github.com/XDoubleU/essentia/pkg/threading"
	"github.com/supabase-community/gotrue-go"

	"goal-tracker/api/internal/config"
	"goal-tracker/api/internal/repositories"
	"goal-tracker/api/pkg/goodreads"
	"goal-tracker/api/pkg/steam"
	"goal-tracker/api/pkg/todoist"
)

type Services struct {
	Auth      *AuthService
	Goals     *GoalService
	Todoist   *TodoistService
	Steam     *SteamService
	Goodreads *GoodreadsService
	WebSocket *WebSocketService
}

func New(
	logger *slog.Logger,
	config config.Config,
	jobQueue *threading.JobQueue,
	repositories *repositories.Repositories,
	supabaseClient gotrue.Client,
	todoistClient todoist.Client,
	steamClient steam.Client,
	goodreadsClient goodreads.Client,
) *Services {
	auth := &AuthService{supabaseUserID: config.SupabaseUserID, client: supabaseClient}
	goodreads := &GoodreadsService{
		logger:     logger,
		profileURL: config.GoodreadsURL,
		goodreads:  repositories.Goodreads,
		client:     goodreadsClient,
	}
	todoist := &TodoistService{
		client:    todoistClient,
		projectID: config.TodoistProjectID,
	}
	steam := &SteamService{
		logger: logger,
		client: steamClient,
		userID: config.SteamUserID,
		steam:  repositories.Steam,
	}
	goals := &GoalService{
		webURL:    config.WebURL,
		states:    repositories.States,
		goals:     repositories.Goals,
		progress:  repositories.Progress,
		todoist:   todoist,
		goodreads: goodreads,
	}

	return &Services{
		Auth:      auth,
		Goals:     goals,
		Todoist:   todoist,
		Steam:     steam,
		Goodreads: goodreads,
		WebSocket: NewWebSocketService(logger, []string{config.WebURL}, jobQueue),
	}
}
