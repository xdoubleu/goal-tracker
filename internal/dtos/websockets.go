package dtos

import (
	"time"

	"github.com/XDoubleU/essentia/pkg/validate"
)

type SubscribeMessageDto struct {
	Subject string `json:"subject"`
} //	@name	SubscribeMessageDto

type StateMessageDto struct {
	LastRefresh  *time.Time `json:"lastRefresh"`
	IsRefreshing bool       `json:"isRefreshing"`
} // @name StateMessageDto

func (dto SubscribeMessageDto) Topic() string {
	return dto.Subject
}

func (dto SubscribeMessageDto) Validate() *validate.Validator {
	v := validate.New()

	return v
}
