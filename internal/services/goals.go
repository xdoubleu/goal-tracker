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
}

func (service GoalService) GetPage(ctx context.Context, user models.User, pageSize int, getAfterID *string) ([]*models.Goal, error) {
	return service.goals.GetPage(ctx, user.ID, pageSize, getAfterID)
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

	return service.goals.Create(ctx, user.ID, createGoalDto.Name, createGoalDto.Description, createGoalDto.Date, createGoalDto.Value, createGoalDto.SourceID, createGoalDto.TypeID, createGoalDto.Score, createGoalDto.StateID)
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

	goal, err := service.GetByID(ctx, id, user)
	if err != nil {
		return nil, err
	}

	_, err = service.progress.Create(ctx, id, *updateGoalDto.Value)
	if err != nil {
		return nil, err
	}

	return service.goals.Update(ctx, *goal, updateGoalDto)
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
