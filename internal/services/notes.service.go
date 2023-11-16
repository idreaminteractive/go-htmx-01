package services

import (
	"context"
	"encoding/json"
	"log"

	"main/internal/db"
	"main/internal/views/dto"

	"github.com/sirupsen/logrus"
)

type NotesService struct {
	queries *db.Queries
	sl      *ServiceLocator
}

func InitNotesService(sl *ServiceLocator, q *db.Queries) *NotesService {
	return &NotesService{sl: sl, queries: q}
}

type AggregateNoteItem struct {
	NoteId  int    `json:"noteId"`
	Content string `json:"content"`
}

type AggregateUserNotes struct {
	ID    int64
	Email string
	Notes []AggregateNoteItem
}

// how we can parse the data and get it into a proper struct using our json agg functions
func (ns *NotesService) GetAggregatedNotes() (*[]AggregateUserNotes, error) {
	// testing....
	ctx := context.Background()
	notesAgg, err := ns.queries.GetUserNoteAggregate(ctx)
	logrus.WithFields(logrus.Fields{"notesAgg": notesAgg}).Info("Notes")
	if err != nil {
		return nil, Errorf(EINTERNAL, "Error when running query %v", err)
	}
	// ok - create a slice of these
	items := []AggregateUserNotes{}
	for _, item := range notesAgg {
		aggregatedNotes := []AggregateNoteItem{}

		if err := json.Unmarshal([]byte(item.Notes.(string)), &aggregatedNotes); err != nil {
			log.Fatal(err)
		}
		items = append(items, AggregateUserNotes{
			ID: item.ID,

			Email: item.Email,
			Notes: aggregatedNotes,
		})
		logrus.WithField("wow", items).Info("pls")
	}
	return &items, nil
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

func (ns *NotesService) GetNoteById(userId, noteId int) (*db.Note, error) {
	ctx := context.Background()
	note, err := ns.queries.GetNoteById(ctx, int64(noteId))
	if err != nil {
		return nil, Errorf(EINTERNAL, "Error when running query %v", err)
	}

	return &note, nil
}

func (ns *NotesService) UpdateNote(userId, noteId int, dto *dto.UpdateNoteDTO) error {

	ctx := context.Background()
	err := ns.queries.UpdateNote(ctx, db.UpdateNoteParams{ID: int64(noteId), Content: dto.Content, IsPublic: dto.IsPublic == "on"})
	if err != nil {
		return Errorf(EINTERNAL, "Error when updating note %v", err)
	}

	return nil
}

func (ns *NotesService) GetPublicNotes() (*[]db.GetPublicNotesRow, error) {
	ctx := context.Background()
	notes, err := ns.queries.GetPublicNotes(ctx)
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

func (ns *NotesService) DeleteNote(userId int, noteId int) error {
	// check and make sure the user owns the note

	ctx := context.Background()
	err := ns.queries.DeleteNote(ctx, int64(noteId))
	if err != nil {
		return Errorf(EINTERNAL, "Error when running query %v", err)
	}
	return nil
}
