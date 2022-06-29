package request

import validation "github.com/go-ozzo/ozzo-validation"

type UserAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (data UserAuth) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Login, validation.Required, validation.Length(4, 20)),
		validation.Field(&data.Password, validation.Required, validation.Length(6, 100)),
	)
}
