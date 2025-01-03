package dtos

import (
	"github.com/XDoubleU/essentia/pkg/validate"
)

type LinkGoalDto struct {
	TypeID           int64             `json:"typeId"`
	TargetValue      *int64            `json:"targetValue"`
	Tag              *string           `json:"tag"`
	ValidationErrors map[string]string `json:"-"`
} //	@name	LinkGoalDto

func (dto *LinkGoalDto) Validate() *validate.Validator {
	v := validate.New()

	dto.ValidationErrors = v.Errors

	return v
}
