package dtos

import (
	"goal-tracker/api/internal/models"
	"time"

	"github.com/XDoubleU/essentia/pkg/validate"
)

type CreateGoalDto struct {
	Name             string            `json:"name"`
	Description      *string           `json:"description"`
	Date             *time.Time        `json:"date"`
	Value            *int64            `json:"value"`
	SourceID         *int64            `json:"sourceId"`
	TypeID           *int64            `json:"typeId"`
	Score            int64             `json:"score"`
	StateID          int64             `json:"stateId"`
	ValidationErrors map[string]string `json:"-"`
} //	@name	CreateGoalDto

type UpdateGoalDto struct {
	Name             *string           `json:"name"`
	Description      *string           `json:"description"`
	Date             *time.Time        `json:"date"`
	Value            *int64            `json:"value"`
	SourceID         *int64            `json:"sourceId"`
	TypeID           *int64            `json:"typeId"`
	Score            *int64            `json:"score"`
	StateID          *int64            `json:"stateId"`
	ValidationErrors map[string]string `json:"-"`
} //	@name	UpdateGoalDto

func (dto *CreateGoalDto) Validate() *validate.Validator {
	v := validate.New()

	validate.Check(v, dto.Name, validate.IsNotEmpty, "name")
	validate.Check(v, dto.Score, validate.IsGreaterThanFunc(int64(0)), "score")
	validate.Check(v, dto.StateID, checkOptions, "stateId")

	dto.ValidationErrors = v.Errors

	return v
}

func (dto *UpdateGoalDto) Validate() *validate.Validator {
	v := validate.New()

	if dto.Name != nil {
		validate.Check(v, *dto.Name, validate.IsNotEmpty, "name")
	}

	if dto.Score != nil {
		validate.Check(v, *dto.Score, validate.IsGreaterThanFunc(int64(0)), "score")
	}

	if dto.StateID != nil {
		validate.Check(v, *dto.StateID, checkOptions, "state")
	}

	dto.ValidationErrors = v.Errors

	return v
}

func checkOptions(value int64) (bool, string) {
	switch value {
	case models.BacklogState.ID, models.PlannedState.ID, models.InProgressState.ID, models.DoneState.ID:
		return true, ""
	default:
		return false, "must be an existing state"
	}
}
