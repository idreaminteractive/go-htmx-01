package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateNoteDTO struct {
	Content  string `form:"content" validate:"required,len=5" errmsg:"Enter valid data for your content"`
	IsPublic string `form:"is_public"`
}

// Similar to zod stuff
func (a CreateNoteDTO) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Content,
			validation.Required.Error("Content is required"),
			validation.Length(2, 0).Error("Content must be at least 2 characters")),
		validation.Field(&a.IsPublic, validation.Required),
	)
}

type UpdateNoteDTO struct {
	CreateNoteDTO
}
