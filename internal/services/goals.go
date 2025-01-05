package services

import (
	"context"
	"fmt"
	"time"

	"github.com/XDoubleU/essentia/pkg/errors"

	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/helper"
	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/repositories"
	"goal-tracker/api/pkg/todoist"
)

type GoalService struct {
	webURL    string
	goals     *repositories.GoalRepository
	states    *repositories.StateRepository
	progress  *repositories.ProgressRepository
	listItems *repositories.ListItemRepository
	todoist   *TodoistService
}

type StateGoalsPair struct {
	State string
	Goals []helper.GoalWithSubGoals
}

func (service *GoalService) GetAllGoalsGroupedByStateAndParentGoal(
	ctx context.Context,
	userID string,
) ([]StateGoalsPair, error) {
	goals, err := service.goals.GetAll(ctx, userID)
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

	states, err := service.states.GetAll(ctx, userID)
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

func (service *GoalService) GetGoalByID(
	ctx context.Context,
	id string,
	userID string,
) (*models.Goal, error) {
	return service.goals.GetByID(ctx, id, userID)
}

func (service *GoalService) GetGoalsByTypeID(
	ctx context.Context,
	id int64,
	userID string,
) ([]models.Goal, error) {
	return service.goals.GetByTypeID(ctx, id, userID)
}

func (service *GoalService) ImportStatesFromTodoist(
	ctx context.Context,
	userID string,
) error {
	sections, err := service.todoist.GetSections(ctx)
	if err != nil {
		return err
	}

	sectionsMap := map[string]todoist.Section{}
	for _, section := range sections {
		sectionsMap[section.ID] = section
	}

	existingStates, err := service.states.GetAll(ctx, userID)
	if err != nil {
		return err
	}

	for _, state := range existingStates {
		_, ok := sectionsMap[state.ID]

		if ok {
			continue
		}

		err = service.states.Delete(ctx, &state, userID)
		if err != nil {
			return err
		}
	}

	for _, section := range sections {
		_, err = service.states.Upsert(
			ctx,
			section.ID,
			userID,
			section.Name,
			section.Order,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *GoalService) ImportGoalsFromTodoist(
	ctx context.Context,
	userID string,
) error {
	tasks, err := service.todoist.GetTasks(ctx)
	if err != nil {
		return err
	}

	tasksMap := map[string]todoist.Task{}
	for _, task := range tasks {
		tasksMap[task.ID] = task
	}

	existingGoals, err := service.goals.GetAll(ctx, userID)
	if err != nil {
		return err
	}

	for _, goal := range existingGoals {
		_, ok := tasksMap[goal.ID]

		if ok {
			continue
		}

		err = service.goals.Delete(ctx, &goal, userID)
		if err != nil {
			return err
		}
	}

	for _, task := range tasksMap {
		_, err = service.goals.Upsert(
			ctx,
			task.ID,
			userID,
			task.ParentID,
			task.Content,
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

func (service *GoalService) LinkGoal(
	ctx context.Context,
	id string,
	userID string,
	linkGoalDto *dtos.LinkGoalDto,
) error {
	if v := linkGoalDto.Validate(); !v.Valid() {
		return errors.ErrFailedValidation
	}

	goal, err := service.goals.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	err = service.goals.Link(
		ctx,
		goal,
		userID,
		*linkGoalDto,
	)
	if err != nil {
		return err
	}

	return service.todoist.UpdateTask(
		ctx,
		goal.ID,
		fmt.Sprintf("%s/goals/%s", service.webURL, goal.ID),
	)
}

func (service *GoalService) UnlinkGoal(
	ctx context.Context,
	id string,
	userID string,
) error {
	goal, err := service.goals.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	err = service.goals.Unlink(
		ctx,
		*goal,
		userID,
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

func (service *GoalService) GetProgressByTypeIDAndDates(
	ctx context.Context,
	typeID int64,
	userID string,
	dateStart time.Time,
	dateEnd time.Time,
) ([]string, []string, error) {
	progresses, err := service.progress.GetByTypeIDAndDates(
		ctx,
		typeID,
		userID,
		dateStart,
		dateEnd,
	)
	if err != nil {
		return nil, nil, err
	}

	progressLabels := []string{}
	progressValues := []string{}

	for _, progress := range progresses {
		progressLabels = append(
			progressLabels,
			progress.Date.Format(models.ProgressDateFormat),
		)
		progressValues = append(progressValues, progress.Value)
	}

	return progressLabels, progressValues, nil
}

func (service *GoalService) SaveProgress(
	ctx context.Context,
	typeID int64,
	userID string,
	progressLabels []string,
	progressValues []string,
) error {
	for i := 0; i < len(progressLabels); i++ {
		_, err := service.progress.Upsert(
			ctx,
			typeID,
			userID,
			progressLabels[i],
			progressValues[i],
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *GoalService) GetListItemsByGoalID(
	ctx context.Context,
	goalID string,
	userID string,
) ([]models.ListItem, error) {
	return service.listItems.GetByGoalID(ctx, goalID, userID)
}

func (service *GoalService) SaveListItem(
	ctx context.Context,
	id int64,
	userID string,
	goalID string,
	value string,
	completed bool,
) (*models.ListItem, error) {
	return service.listItems.Upsert(ctx, id, userID, goalID, value, completed)
}
