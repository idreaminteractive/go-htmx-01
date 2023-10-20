package services

import (
	"context"
	"main/internal/db"
	"main/internal/views/dto"

	"github.com/sirupsen/logrus"
)

type NotesService struct {
	queries *db.Queries
}

func InitNotesService(q *db.Queries) *NotesService {
	return &NotesService{queries: q}
}

func (ns *NotesService) GetNotesForUserId(userId int) (*[]db.ListNotesForUserRow, error) {
	ctx := context.Background()
	notes, err := ns.queries.ListNotesForUser(ctx, int64(userId))
	logrus.WithFields(logrus.Fields{"userId": userId, "notes": notes}).Info("Notes")
	if err != nil {
		return nil, Errorf(EINTERNAL, "Error when running query %v", err)
	}
	return &notes, nil
}

func (ns *NotesService) CreateNewNote(userId int, dto *dto.CreateNoteDTO) (*db.Note, error) {

	ctx := context.Background()
	note, err := ns.queries.CreateNote(ctx, db.CreateNoteParams{Content: dto.Content, IsPublic: dto.IsPublic == "on", UserID: int64(userId)})
	logrus.WithFields(logrus.Fields{"userId": userId, "notes": note}).Info("Created note")
	if err != nil {
		return nil, Errorf(EINTERNAL, "Error when running query %v", err)
	}
	return &note, nil
}
