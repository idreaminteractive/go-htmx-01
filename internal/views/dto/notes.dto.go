package dto

type CreateNoteDTO struct {
	Content  string `form:"content" validate:"required,len=5"`
	IsPublic string `form:"is_public"`
}

type UpdateNoteDTO struct {
	CreateNoteDTO
}
