package dtos

import (
	"github.com/XDoubleU/essentia/pkg/validate"
)

type LinkGoalDto struct {
	ID               string            `json:"id"`
	TargetValue      int64             `json:"targetValue"`
	TypeID           int64             `json:"typeId"`
	ValidationErrors map[string]string `json:"-"`
} //	@name	CreateGoalDto

type UpdateGoalDto struct {
	TargetValue      *int64            `json:"targetValue"`
	SourceID         *int64            `json:"sourceId"`
	TypeID           *int64            `json:"typeId"`
	ValidationErrors map[string]string `json:"-"`
} //	@name	UpdateGoalDto

func (dto *LinkGoalDto) Validate() *validate.Validator {
	v := validate.New()

	dto.ValidationErrors = v.Errors

	return v
}

func (dto *UpdateGoalDto) Validate() *validate.Validator {
	v := validate.New()

	dto.ValidationErrors = v.Errors

	return v
}
