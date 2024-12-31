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
		goalsMap[goal.Goal.State] = append(goalsMap[goal.Goal.State], goal)
	}

	//nolint:godox //I know
	// TODO deal with this in a better way
	// important that the order is static (maybe configureable)
	// important that the actual string can be dynamic
	// important that we don't do unnecessary API calls to obtain this
	states := []string{
		"In Progress",
		"Planned",
		"Backlog",
	}

	result := []StateGoalsPair{}
	for _, state := range states {
		pair := StateGoalsPair{
			State: state,
			Goals: goalsMap[state],
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
	_, sectionsIDNameMap, err := service.todoist.GetSections(ctx)
	if err != nil {
		return err
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

		_, err = service.goals.Update(ctx, sectionsIDNameMap, goal, task)
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
			sectionsIDNameMap[task.SectionID],
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

func (service GoalService) FetchProgress(
	ctx context.Context,
	typeID int64,
) ([]string, []int64, error) {
	progresses, err := service.progress.Fetch(ctx, typeID)
	if err != nil {
		return nil, nil, err
	}

	progressLabels := []string{}
	progressValues := []int64{}

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
	progressValues []int64,
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
