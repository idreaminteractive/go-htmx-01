package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateNoteDTO struct {
	Content  string `form:"content" validate:"required,len=5" errmsg:"Enter valid data for your content"`
	IsPublic string `form:"is_public"`
}

// create zod like rules?
// i want this to be able to be parsed + built into something that
// we can pass into a form var to parse out properly + display as err messages
func (a CreateNoteDTO) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Content, validation.Required.Error("WWWWW"), validation.Length(4, 0).Error("ffsdf")),
		validation.Field(&a.IsPublic, validation.Required),
	)
}

type UpdateNoteDTO struct {
	CreateNoteDTO
}
