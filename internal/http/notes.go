package http

import (
	"main/internal/views"
	"main/internal/views/dto"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) registerNoteRoutes(group *echo.Group) {
	group.POST("", s.handleCreateNote)
	group.DELETE("/:id", s.handleDeleteNote)

}

func (s *Server) handleCreateNote(c echo.Context) error {
	sp, err := s.sessionService.ReadSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Could not read session")

	}

	var notePayload dto.CreateNoteDTO

	if err := c.Bind(&notePayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// csrf_value := getCSRFValueFromContext(c)
	if err := c.Validate(notePayload); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}

	logrus.WithField("Note:", notePayload).Info("Crearting....")

	// prob don't need to pass in ref hjere?
	_, err = s.notesService.CreateNewNote(sp.UserId, &notePayload)
	if err != nil {
		logrus.Error("Error in creating note")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// return our notes template w/ htmx ONLY...
	userNotes, err := s.notesService.GetNotesForUserId(sp.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not fetch notes for user")
	}
	component := views.NotesListing(userNotes)
	// no need to redirect here
	c.Response().Header().Set("HX-Push-Url", "/dashboard")
	renderComponent(component, c)

	return nil

}

func (s *Server) handleDeleteNote(c echo.Context) error {
	sp, err := s.sessionService.ReadSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Could not read session")

	}

	noteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error("Oof")
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	logrus.WithField("c", noteId).Info("hey")
	// todo-  delete the note + return the notes list remaining
	err = s.notesService.DeleteNote(sp.UserId, noteId)
	if err != nil {
		logrus.WithField("noteId", noteId).Error("Could note delete")
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// return our notes template w/ htmx ONLY...
	userNotes, err := s.notesService.GetNotesForUserId(sp.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not fetch notes for user")
	}
	component := views.NotesListing(userNotes)
	c.Response().Header().Set("HX-Push-Url", "/dashboard")
	renderComponent(component, c)

	return nil

}
