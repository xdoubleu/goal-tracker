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
) ([]todoist.Section, error) {
	sections, err := service.client.GetAllSections(ctx, service.projectID)
	if err != nil {
		return nil, err
	}

	return sections, nil
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

func (service TodoistService) UpdateTask(
	ctx context.Context,
	id string,
	description string,
) error {
	//nolint:godox //I know
	//TODO: other fields

	//nolint:exhaustruct //other fields are skipped for now
	_, err := service.client.UpdateTask(ctx, id, todoist.UpdateTaskDto{
		Description: &description,
	})

	return err
}
