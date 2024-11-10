package services

import (
	"goal-tracker/api/pkg/todoist"
)

type TodoistService struct {
	client todoist.Client
}

type SectionNameTasksPair struct {
	SectionName string
	Tasks       []todoist.Task
}

func (service TodoistService) GetTasksFromProjectGroupedBySection(projectID string) ([]SectionNameTasksPair, error) {
	sections, err := service.client.GetAllSections(projectID)
	if err != nil {
		return nil, err
	}

	tasks, err := service.client.GetActiveTasks(projectID)
	if err != nil {
		return nil, err
	}

	tasksMap := map[string][]todoist.Task{}
	for _, t := range *tasks {
		tasksMap[t.SectionId] = append(tasksMap[t.SectionId], t)
	}

	result := []SectionNameTasksPair{}
	for _, s := range *sections {
		pair := SectionNameTasksPair{
			SectionName: s.Name,
			Tasks:       tasksMap[s.Id],
		}
		result = append(result, pair)
	}

	return result, nil
}
