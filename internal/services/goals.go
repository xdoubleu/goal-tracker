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

func (service GoalService) Create(
	ctx context.Context,
	user models.User,
	createGoalDto *dtos.CreateGoalDto,
) (*models.Goal, error) {
	if v := createGoalDto.Validate(); !v.Valid() {
		return nil, errors.ErrFailedValidation
	}

	return nil, nil
	//return service.goals.Create(ctx, user.ID, createGoalDto.Name, createGoalDto.Description, createGoalDto.Date, createGoalDto.Value, createGoalDto.SourceID, createGoalDto.TypeID, createGoalDto.Score, createGoalDto.StateID)
}

func (service GoalService) Update(
	ctx context.Context,
	user models.User,
	id string,
	updateGoalDto *dtos.UpdateGoalDto,
) (*models.Goal, error) {
	if v := updateGoalDto.Validate(); !v.Valid() {
		return nil, errors.ErrFailedValidation
	}

	return nil, nil

	/*
		goal, err := service.GetByID(ctx, id, user)
		if err != nil {
			return nil, err
		}

		_, err = service.progress.Create(ctx, id, *updateGoalDto.Value)
		if err != nil {
			return nil, err
		}

		return service.goals.Update(ctx, *goal, updateGoalDto)*/
}

func (service GoalService) Delete(
	ctx context.Context,
	user models.User,
	id string,
) (*models.Goal, error) {
	goal, err := service.GetByID(ctx, id, user)
	if err != nil {
		return nil, err
	}

	err = service.goals.Delete(ctx, goal)
	if err != nil {
		return nil, err
	}

	return goal, nil
}
