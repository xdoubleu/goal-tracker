package services

import (
	"github.com/supabase-community/gotrue-go"

	"goal-tracker/api/internal/config"
	"goal-tracker/api/internal/repositories"
	"goal-tracker/api/internal/temptools"
	"goal-tracker/api/pkg/steam"
	"goal-tracker/api/pkg/todoist"
)

type Services struct {
	Auth      AuthService
	Goals     GoalService
	Todoist   TodoistService
	Steam     SteamService
	WebSocket *WebSocketService
}

func New(
	config config.Config,
	jobQueue *temptools.JobQueue,
	repositories repositories.Repositories,
	supabaseClient gotrue.Client,
	todoistClient todoist.Client,
	steamClient steam.Client,
) Services {
	auth := AuthService{client: supabaseClient}
	todoist := TodoistService{client: todoistClient, projectID: config.TodoistProjectID}
	steam := SteamService{client: steamClient, userID: config.SteamUserID}
	goals := GoalService{
		webURL:   config.WebURL,
		goals:    repositories.Goals,
		progress: repositories.Progress,
		todoist:  todoist,
	}

	return Services{
		Auth:      auth,
		Goals:     goals,
		Todoist:   todoist,
		Steam:     steam,
		WebSocket: NewWebSocketService([]string{config.WebURL}, jobQueue),
	}
}
