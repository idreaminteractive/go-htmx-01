package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserLoginDTO struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
}

func (a UserLoginDTO) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email,
			validation.Required.Error("Email is required"),
			is.Email.Error("Must be a valid email"),
		),
		validation.Field(&a.Password, validation.Required.Error("Password is required")),
	)
}

type UserLoginFormErrors struct {
	Message string
}
