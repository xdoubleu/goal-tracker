package dtos

import (
	"time"

	"github.com/XDoubleU/essentia/pkg/validate"
)

type SubscribeMessageDto struct {
	Subject string `json:"subject"`
}

type StateMessageDto struct {
	LastRefresh  *time.Time `json:"lastRefresh"`
	IsRefreshing bool       `json:"isRefreshing"`
}

func (dto SubscribeMessageDto) Topic() string {
	return dto.Subject
}

func (dto SubscribeMessageDto) Validate() *validate.Validator {
	v := validate.New()

	return v
}
