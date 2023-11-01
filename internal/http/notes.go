package http

import (
	"fmt"
	"main/internal/views"
	"main/internal/views/dto"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) registerNoteRoutes(group *echo.Group) {
	group.POST("", s.handleCreateNote)
	group.DELETE("/:id", s.handleDeleteNote)
	group.GET("/:id/edit", s.handleGetEditForm)
	group.PUT("/:id/edit", s.handlePutEditForm)

}

func (s *Server) handlePutEditForm(c echo.Context) error {
	// return the note rendered in place!
	sp, err := s.sessionService.ReadSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Could not read session")
	}

	var notePayload dto.UpdateNoteDTO

	if err := c.Bind(&notePayload); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(notePayload); err != nil {
		logrus.WithField("e", err).Error("Error on validate")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}

	noteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	// prob don't need to pass in ref hjere?
	err = s.notesService.UpdateNote(sp.UserId, noteId, &notePayload)
	// note, we could target the individual note to show
	if err != nil {
		logrus.Error("Error in updating note")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

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

func (s *Server) handleGetEditForm(c echo.Context) error {
	// get our note by id
	sp, err := s.sessionService.ReadSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Could not read session")

	}
	noteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	userNote, err := s.notesService.GetNoteById(sp.UserId, noteId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not fetch note edit form")
	}
	component := views.EditNoteForm(userNote)
	// no need to redirect here
	renderComponent(component, c)

	return nil
}

func (s *Server) handleCreateNote(c echo.Context) error {
	sp, err := s.sessionService.ReadSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Could not read session")

	}

	var notePayload dto.CreateNoteDTO

	if err := c.Bind(&notePayload); err != nil {
		logrus.WithField("e", err).Error("Error on bind")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// csrf_value := getCSRFValueFromContext(c)
	if err := c.Validate(notePayload); err != nil {
		logrus.WithField("e", err).Error("Error on validate")

		/// do tiT?
		errs := err.(validator.ValidationErrors)
		rsf := reflect.TypeOf(notePayload)

		fmt.Printf("Reflect: %v", rsf)
		for _, e := range errs {
			snp := e.StructNamespace()
			errmgs := errorTagFunc[dto.CreateNoteDTO](notePayload, snp, e.Field(), e.ActualTag())
			fmt.Printf("err %v\n", errmgs)
			// can translate each error one at a time.
			fmt.Printf("Namespace %v\n", e.Namespace())
			fmt.Printf("Field %v\n", e.Field())
			fmt.Printf("StructNamespace %v\n", e.StructNamespace())
			fmt.Printf("StructField %v\n", e.StructField())
			fmt.Printf("Tag %v\n", e.Tag())
			fmt.Printf("ActualTag %v\n", e.ActualTag())
			fmt.Printf("Type %v\n", e.Type())
			fmt.Printf("Type %v\n", e.Type())
			fmt.Printf("Value %v\n", e.Value())
			fmt.Printf("Param %v\n", e.Param())
			fmt.Println()
		}
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
