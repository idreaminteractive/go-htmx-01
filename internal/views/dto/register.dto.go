package dto

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RegisterDTO struct {
	Handle          string `form:"handle"`
	Email           string `form:"email"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirm_password"`
}

func (reg *RegisterDTO) Bind(r *http.Request) error {

	return nil
}

func (a RegisterDTO) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Handle,
			validation.Required.Error("Chat handle is required"),
			validation.Length(2, 0).Error("Chat handle must be at least 2 characters")),
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, validation.Required.Error("Password is required")),
		validation.Field(&a.ConfirmPassword, validation.Required.Error("Confirmed password is required")),
	)
}
