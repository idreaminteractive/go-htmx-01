package dto

type CreateNoteDTO struct {
	Content  string `form:"content" validate:"required"`
	IsPublic string `form:"is_public"`
}

type UpdateNoteDTO struct {
	CreateNoteDTO
}
