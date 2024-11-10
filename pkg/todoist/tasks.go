package todoist

import (
	"fmt"
	"net/http"
)

var TASKS_ENDPOINT = "tasks"

type Task struct {
	Id           string   `json:"id"`
	ProjectId    string   `json:"project_id"`
	SectionId    string   `json:"section_id"`
	Content      string   `json:"content"`
	Description  string   `json:"description"`
	IsCompleted  bool     `json:"is_completed"`
	Labels       []string `json:"labels"`
	ParentId     string   `json:"parent_id"`
	Order        int      `json:"order"`
	Priority     int      `json:"priority"`
	Due          Due      `json:"due"`
	Url          string   `json:"url"`
	CommentCount int      `json:"comment_count"`
	CreatedAt    string   `json:"created_at"`
	CreatorId    string   `json:"creator_id"`
	AssigneeId   string   `json:"assignee_id"`
	AssignerId   string   `json:"assigner_id"`
	Duration     Duration `json:"duration"`
}

type Due struct {
	String      string `json:"string"`
	Date        string `json:"date"`
	IsRecurring bool   `json:"is_recurring"`
	Datetime    string `json:"datetime"`
	Timezone    string `json:"timezone"`
}

type Duration struct {
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
}

func (client Client) GetActiveTasks(projectId string) (*[]Task, error) {
	query := fmt.Sprintf("project_id=%s", projectId)

	var tasks *[]Task
	err := client.sendRequest(http.MethodGet, TASKS_ENDPOINT, query, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
