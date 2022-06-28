package requests

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type GettingStock struct {
	Id        string `form:"id" json:"id"`
	ClassCode string `form:"classCode" json:"classCode"`
}

func (data GettingStock) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Id, validation.Required),
		validation.Field(&data.ClassCode, validation.Required, validation.In("TQBR", "SPBMX")),
	)
}
