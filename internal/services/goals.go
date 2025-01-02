package services

import (
	"context"
	"fmt"

	"github.com/XDoubleU/essentia/pkg/errors"

	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/helper"
	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/repositories"
	"goal-tracker/api/pkg/todoist"
)

type GoalService struct {
	webURL   string
	goals    repositories.GoalRepository
	states   repositories.StateRepository
	progress repositories.ProgressRepository
	todoist  TodoistService
}

type StateGoalsPair struct {
	State string
	Goals []helper.GoalWithSubGoals
}

func (service GoalService) GetAllGroupedByStateAndParentGoal(
	ctx context.Context,
) ([]StateGoalsPair, error) {
	goals, err := service.goals.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	goalTree := helper.NewGoalTree()
	for _, goal := range goals {
		if !goalTree.TryAdd(goal) {
			return nil, fmt.Errorf("failed to add goal %s to tree", goal.ID)
		}
	}

	goalsMap := map[string][]helper.GoalWithSubGoals{}
	for _, goal := range goalTree.ToSlice() {
		goalsMap[goal.Goal.StateID] = append(goalsMap[goal.Goal.StateID], goal)
	}

	states, err := service.states.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := []StateGoalsPair{}
	for _, state := range states {
		pair := StateGoalsPair{
			State: state.Name,
			Goals: goalsMap[state.ID],
		}
		result = append(result, pair)
	}

	return result, nil
}

func (service GoalService) GetByID(
	ctx context.Context,
	id string,
) (*models.Goal, error) {
	return service.goals.GetByID(ctx, id)
}

func (service GoalService) GetByTypeID(
	ctx context.Context,
	id int64,
) ([]models.Goal, error) {
	return service.goals.GetByTypeID(ctx, id)
}

func (service GoalService) ImportFromTodoist(ctx context.Context) error {
	states, err := service.states.GetAll(ctx)
	if err != nil {
		return err
	}

	if len(states) == 0 {
		sections, err := service.todoist.GetSections(ctx)
		if err != nil {
			return err
		}

		for _, section := range sections {
			state, err := service.states.Create(ctx, section.ID, section.Name, section.Order)
			if err != nil {
				return err
			}

			states = append(states, *state)
		}
	}

	tasks, err := service.todoist.GetTasks(ctx)
	if err != nil {
		return err
	}

	tasksMap := map[string]todoist.Task{}
	for _, task := range tasks {
		tasksMap[task.ID] = task
	}

	existingGoals, err := service.goals.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, goal := range existingGoals {
		task, ok := tasksMap[goal.ID]

		if !ok {
			err = service.goals.Delete(ctx, &goal)
			if err != nil {
				return err
			}
			continue
		}

		_, err = service.goals.Update(ctx, goal, task)
		if err != nil {
			return err
		}

		delete(tasksMap, goal.ID)
	}

	// only new tasks remain
	for _, task := range tasksMap {
		_, err = service.goals.Create(
			ctx,
			task.ID,
			task.ParentID,
			task.Content,
			false,
			nil,
			nil,
			task.SectionID,
			task.Due,
			task.Order,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service GoalService) Link(
	ctx context.Context,
	id string,
	linkGoalDto *dtos.LinkGoalDto,
) error {
	if v := linkGoalDto.Validate(); !v.Valid() {
		return errors.ErrFailedValidation
	}

	goal, err := service.goals.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = service.goals.Link(
		ctx,
		*goal,
		*linkGoalDto,
	)
	if err != nil {
		return err
	}

	return service.todoist.UpdateTask(
		ctx,
		goal.ID,
		fmt.Sprintf("%s/%s", service.webURL, goal.ID),
	)
}

func (service GoalService) Unlink(
	ctx context.Context,
	id string,
) error {
	goal, err := service.goals.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = service.goals.Unlink(
		ctx,
		*goal,
	)
	if err != nil {
		return err
	}

	return service.todoist.UpdateTask(
		ctx,
		goal.ID,
		"",
	)
}

func (service GoalService) FetchProgress(
	ctx context.Context,
	typeID int64,
) ([]string, []string, error) {
	progresses, err := service.progress.Fetch(ctx, typeID)
	if err != nil {
		return nil, nil, err
	}

	progressLabels := []string{}
	progressValues := []string{}

	for _, progress := range progresses {
		progressLabels = append(
			progressLabels,
			progress.Date.Time.Format(models.ProgressDateFormat),
		)
		progressValues = append(progressValues, progress.Value)
	}

	return progressLabels, progressValues, nil
}

func (service GoalService) SaveProgress(
	ctx context.Context,
	typeID int64,
	progressLabels []string,
	progressValues []string,
) error {
	for i := 0; i < len(progressLabels); i++ {
		_, err := service.progress.Save(
			ctx,
			typeID,
			progressLabels[i],
			progressValues[i],
		)
		if err != nil {
			return err
		}
	}
	return nil
}
