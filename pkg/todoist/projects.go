package todoist

import (
	"context"
	"net/http"
)

var PROJECTS_ENDPOINT = "projects"

type Project struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Color          string `json:"color"`
	ParentId       string `json:"parent_id"`
	Order          int    `json:"order"`
	CommentCount   int    `json:"comment_count"`
	IsShared       bool   `json:"is_shared"`
	IsFavorite     bool   `json:"is_favorite"`
	IsInboxProject bool   `json:"is_inbox_project"`
	IsTeamInbox    bool   `json:"is_team_inbox"`
	ViewStyle      string `json:"view_style"`
	Url            string `json:"url"`
}

func (client Client) GetAllProjects(ctx context.Context) (*[]Project, error) {
	var projects *[]Project
	err := client.sendRequest(ctx, http.MethodGet, PROJECTS_ENDPOINT, "", &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}
