package services

import (
	"context"
	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/repositories"

	"github.com/XDoubleU/essentia/pkg/errors"
)

type GoalService struct {
	goals    repositories.GoalRepository
	progress repositories.ProgressRepository
	todoist  TodoistService
}

type StateGoalsPair struct {
	State string
	Goals []models.Goal
}

func (service GoalService) GetAllGroupedByState(ctx context.Context, userID string) ([]StateGoalsPair, error) {
	sections, sectionsIdNameMap, err := service.todoist.GetSections(ctx)
	if err != nil {
		return nil, err
	}

	tasks, err := service.todoist.GetTasks(ctx)
	if err != nil {
		return nil, err
	}

	configuredGoals, err := service.goals.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	goalsMap := map[string][]models.Goal{}
	for _, configuredGoal := range configuredGoals {
		goalsMap[configuredGoal.State] = append(goalsMap[configuredGoal.State], *configuredGoal)
	}

	for _, task := range *tasks {
		goal := models.NewGoalFromTask(task, userID, sectionsIdNameMap[task.SectionId])
		goalsMap[goal.State] = append(goalsMap[goal.State], goal)
	}

	result := []StateGoalsPair{}
	for _, section := range *sections {
		pair := StateGoalsPair{
			State: section.Name,
			Goals: goalsMap[section.Name],
		}
		result = append(result, pair)
	}

	return result, nil
}

func (service GoalService) GetByID(ctx context.Context, id string, user models.User) (*models.Goal, error) {
	return service.goals.GetByID(ctx, id, user.ID)
}

func (service GoalService) Link(
	ctx context.Context,
	id string,
	user models.User,
	linkGoalDto *dtos.LinkGoalDto,
) error {
	if v := linkGoalDto.Validate(); !v.Valid() {
		return errors.ErrFailedValidation
	}

	_, sectionsMap, err := service.todoist.GetSections(ctx)
	if err != nil {
		return err
	}

	task, err := service.todoist.GetTaskByID(ctx, id)
	if err != nil {
		return err
	}

	_, err = service.goals.Create(ctx, id, user.ID, task.Content, true, linkGoalDto.TargetValue, linkGoalDto.TypeID, sectionsMap[task.SectionId])
	return err
}
