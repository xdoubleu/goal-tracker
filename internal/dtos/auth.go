package dtos

import "github.com/XDoubleU/essentia/pkg/validate"

type SignInDto struct {
	Email            string            `json:"email"`
	Password         string            `json:"password"`
	RememberMe       bool              `json:"rememberMe"`
	ValidationErrors map[string]string `json:"-"`
} //	@name	SignInDto

func (dto *SignInDto) Validate() *validate.Validator {
	v := validate.New()

	validate.Check(v, dto.Email, validate.IsNotEmpty, "email")
	validate.Check(v, dto.Password, validate.IsNotEmpty, "password")

	dto.ValidationErrors = v.Errors

	return v
}
