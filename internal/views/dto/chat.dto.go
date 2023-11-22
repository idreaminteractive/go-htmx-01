package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ChatMessageDTO struct {
	Message string `form:"messageInput"`
}

func (a ChatMessageDTO) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Message,
			validation.Required.Error("Message is required"),
			validation.Length(2, 0).Error("Message  must be at least 2 characters")),
	)
}
