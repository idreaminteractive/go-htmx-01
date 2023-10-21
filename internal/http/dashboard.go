package http

import (
	"main/internal/views"
	"main/internal/views/dto"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) requireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.Info("Check")
		sess, err := s.sessionService.ReadSession(c)
		if err != nil {
			logrus.WithField("err", err).Error("Error in getting session")
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if sess.UserId == 0 {
			logrus.Error("Not logged in")
			return c.Redirect(http.StatusFound, "/login")
		}
		return next(c)
	}
}

func (s *Server) registerLoggedInRoutes(group *echo.Group) {
	group.GET("/", s.handleDashboard)
	group.POST("/create-note/", s.handleCreateNote)

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
	c.Response().Header().Set("HX-Push-Url", "/dashboard")
	renderComponent(component, c)

	return nil

	// full redirect back to home
	// return c.Redirect(http.StatusFound, "/dashboard")

}

// will be the main page of the system
func (s *Server) handleDashboard(c echo.Context) error {
	// find our logged in user to get their personal notes
	sp, err := s.sessionService.ReadSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read session")

	}

	userNotes, err := s.notesService.GetNotesForUserId(sp.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not fetch notes for user")
	}
	csrf_value := getCSRFValueFromContext(c)
	component := views.Dashboard(csrf_value, userNotes)
	base := views.Base(component)
	renderComponent(base, c)
	return nil
}
