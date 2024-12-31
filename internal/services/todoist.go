package services

import (
	"context"

	"goal-tracker/api/pkg/todoist"
)

type TodoistService struct {
	client    todoist.Client
	projectID string
}

func (service TodoistService) GetSections(
	ctx context.Context,
) ([]todoist.Section, map[string]string, error) {
	sections, err := service.client.GetAllSections(ctx, service.projectID)
	if err != nil {
		return nil, nil, err
	}

	sectionsIDNameMap := map[string]string{}
	for _, section := range sections {
		sectionsIDNameMap[section.ID] = section.Name
	}

	return sections, sectionsIDNameMap, nil
}

func (service TodoistService) GetTasks(ctx context.Context) ([]todoist.Task, error) {
	return service.client.GetActiveTasks(ctx, service.projectID)
}

func (service TodoistService) GetTaskByID(
	ctx context.Context,
	id string,
) (*todoist.Task, error) {
	return service.client.GetActiveTask(ctx, id)
}

func (service TodoistService) UpdateTask(ctx context.Context, id string, description string) error {
	_, err := service.client.UpdateTask(ctx, id, todoist.UpdateTaskDto{
		Description: &description,
	})

	return err
}
