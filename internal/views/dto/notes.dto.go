package dto

type CreateNoteDTO struct {
	Content  string `form:"content" validate:"required,len=5" errmsg:"Enter valid data for your content"`
	IsPublic string `form:"is_public"`
}

type UpdateNoteDTO struct {
	CreateNoteDTO
}
