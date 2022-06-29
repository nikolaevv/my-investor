package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type BuyingShare struct {
	Id        string `form:"id" json:"id"`
	ClassCode string `form:"classCode" json:"classCode"`
	Quantity  int    `form:"quantity" json:"quantity"`
}

func (data BuyingShare) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Id, validation.Required),
		validation.Field(&data.ClassCode, validation.Required, validation.In("TQBR", "SPBMX")),
		validation.Field(&data.Quantity, validation.Required, validation.Min(1)),
	)
}
