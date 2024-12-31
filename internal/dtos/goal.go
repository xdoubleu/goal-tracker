package dtos

import (
	"github.com/XDoubleU/essentia/pkg/validate"
)

type LinkGoalDto struct {
	TargetValue      int64             `json:"targetValue"`
	TypeID           int64             `json:"typeId"`
	ValidationErrors map[string]string `json:"-"`
} //	@name	LinkGoalDto

type RelinkGoalDto struct {
	TargetValue      *int64            `json:"targetValue"`
	TypeID           *int64            `json:"typeId"`
	ValidationErrors map[string]string `json:"-"`
} //	@name	RelinkGoalDto

func (dto *LinkGoalDto) Validate() *validate.Validator {
	v := validate.New()

	dto.ValidationErrors = v.Errors

	return v
}

func (dto *RelinkGoalDto) Validate() *validate.Validator {
	v := validate.New()

	dto.ValidationErrors = v.Errors

	return v
}
