package dto

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserLoginDTO struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
}

// do we combine validate + bind in one?
func (u *UserLoginDTO) Bind(r *http.Request) error {

	return nil
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
