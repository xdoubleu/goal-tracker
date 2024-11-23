package todoist

import (
	"context"
	"fmt"
	"net/http"
)

const TasksEndpoint = "tasks"

type Task struct {
	ID           string   `json:"id"`
	ProjectID    string   `json:"project_id"`
	SectionID    string   `json:"section_id"`
	Content      string   `json:"content"`
	Description  string   `json:"description"`
	IsCompleted  bool     `json:"is_completed"`
	Labels       []string `json:"labels"`
	ParentID     *string  `json:"parent_id"`
	Order        int      `json:"order"`
	Priority     int      `json:"priority"`
	Due          Due      `json:"due"`
	URL          string   `json:"url"`
	CommentCount int      `json:"comment_count"`
	CreatedAt    string   `json:"created_at"`
	CreatorID    string   `json:"creator_id"`
	AssigneeID   string   `json:"assignee_id"`
	AssignerID   string   `json:"assigner_id"`
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

func (client client) GetActiveTasks(
	ctx context.Context,
	projectID string,
) ([]Task, error) {
	query := fmt.Sprintf("project_id=%s", projectID)

	var tasks []Task
	err := client.sendRequest(ctx, http.MethodGet, TasksEndpoint, query, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (client client) GetActiveTask(ctx context.Context, taskID string) (*Task, error) {
	endpoint := fmt.Sprintf("%s/%s", TasksEndpoint, taskID)

	var task *Task
	err := client.sendRequest(ctx, http.MethodGet, endpoint, "", &task)
	if err != nil {
		return nil, err
	}

	return task, nil
}
