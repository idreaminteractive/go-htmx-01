package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateNoteDTO struct {
	Content  string `form:"content" validate:"required,len=5" errmsg:"Enter valid data for your content"`
	IsPublic string `form:"is_public"`
}

// create zod like rules?
func (a CreateNoteDTO) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Content, validation.Required, validation.Length(1, 0)),
		validation.Field(&a.IsPublic, validation.Required),
	)
}

type UpdateNoteDTO struct {
	CreateNoteDTO
}
