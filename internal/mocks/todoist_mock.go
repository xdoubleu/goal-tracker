package mocks

import (
	"context"
	"goal-tracker/api/pkg/todoist"
)

type MockTodoistClient struct {
}

func NewMockTodoistClient() todoist.Client {
	return MockTodoistClient{}
}

func (client MockTodoistClient) GetActiveTask(ctx context.Context, id string) (*todoist.Task, error) {
	return &todoist.Task{}, nil
}

func (client MockTodoistClient) GetActiveTasks(ctx context.Context, projectId string) ([]todoist.Task, error) {
	return []todoist.Task{}, nil
}

func (client MockTodoistClient) GetAllProjects(ctx context.Context) ([]todoist.Project, error) {
	return []todoist.Project{}, nil
}

func (client MockTodoistClient) GetAllSections(ctx context.Context, projectId string) ([]todoist.Section, error) {
	return []todoist.Section{}, nil
}
