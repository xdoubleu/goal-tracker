package services

import (
	"context"

	"github.com/XDoubleU/essentia/pkg/errors"

	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/helper"
	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/repositories"
)

type GoalService struct {
	goals    repositories.GoalRepository
	progress repositories.ProgressRepository
	todoist  TodoistService
}

type StateGoalsPair struct {
	State string
	Goals []helper.GoalWithSubGoals
}

func (service GoalService) GetAllGroupedByStateAndParentGoal(
	ctx context.Context,
	userID string,
) ([]StateGoalsPair, error) {
	sections, sectionsIDNameMap, err := service.todoist.GetSections(ctx)
	if err != nil {
		return nil, err
	}

	tasks, err := service.todoist.GetTasks(ctx)
	if err != nil {
		return nil, err
	}

	goals, err := service.goals.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		goal := models.NewGoalFromTask(task, userID, sectionsIDNameMap[task.SectionID])
		goals = append(goals, goal)
	}

	goalTree := helper.NewGoalTree()
	for _, goal := range goals {
		if !goalTree.TryAdd(goal) {
			return nil, err
		}
	}

	goalsMap := map[string][]helper.GoalWithSubGoals{}
	for _, goal := range goalTree.ToSlice() {
		goalsMap[goal.Goal.State] = append(goalsMap[goal.Goal.State], goal)
	}

	result := []StateGoalsPair{}
	for _, section := range sections {
		pair := StateGoalsPair{
			State: section.Name,
			Goals: goalsMap[section.Name],
		}
		result = append(result, pair)
	}

	return result, nil
}

func (service GoalService) GetByID(
	ctx context.Context,
	id string,
	user models.User,
) (*models.Goal, error) {
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

	_, err = service.goals.Create(
		ctx,
		id,
		task.ParentID,
		user.ID,
		task.Content,
		true,
		linkGoalDto.TargetValue,
		linkGoalDto.TypeID,
		sectionsMap[task.SectionID],
	)
	return err
}
