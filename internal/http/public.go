package http

import (
	"fmt"
	"main/internal/services"
	"main/internal/views/base"
	"main/internal/views/dto"
	"main/internal/views/loggedin"
	"main/internal/views/login"
	"main/internal/views/register"
	"main/internal/views/root"
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (s *Server) handleLogoutGet(w http.ResponseWriter, r *http.Request) {
	// kill session + redirect (i should not need to post anywhere)
	// write a blank session
	s.services.SessionService.WriteSession(w, r, services.SessionPayload{})
	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

// will be the main page of the system
// let's mirror our current live version that pulls in the stuff
func (s *Server) handleRootGet(w http.ResponseWriter, r *http.Request) {

	// get our message count

	csrf_value := csrfFromRequest(r)
	count, err := s.services.ChatService.GetTotalMessagCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body := root.Root(count)

	// renderComponent(
	base := base.Base(
		base.BaseData{
			Body:  body,
			CSRF:  csrf_value,
			Title: "Chat App",
		},
	)
	base.Render(r.Context(), w)

}

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	htmx.NewResponse().RenderTempl(r.Context(), w, base.NotFoundComponent())
}

func (s *Server) handleLoginPost(w http.ResponseWriter, r *http.Request) {

	var user dto.UserLoginDTO
	if err := render.Bind(r, &user); err != nil {
		base.InternalServerError(err).Render(r.Context(), w)
		return

	}
	var formErrors validation.Errors
	if err := user.Validate(); err != nil {
		formErrors = err.(validation.Errors)
		component := login.LoginScreen(login.LoginScreenProps{
			LastSubmission: user,
			Errors:         formErrors,
		})

		htmx.NewResponse().
			Retarget("#loginScreen").
			Reswap(htmx.SwapOuterHTML).
			RenderTempl(r.Context(), w, component)

		return
	}

	// create our user + id
	results, err := s.services.AuthenticationService.Authenticate(user)
	// auth fails
	if err != nil {
		component := login.LoginScreen(login.LoginScreenProps{
			LastSubmission: user,
			Errors: map[string]error{
				"email": validation.NewError("", "Invalid email or password"),
			},
		})
		htmx.NewResponse().
			Retarget("#loginScreen").
			Reswap(htmx.SwapOuterHTML).
			RenderTempl(r.Context(), w, component)
		return
	}

	// successful login, lesgo.

	// create our session + stuff
	fmt.Println("Writing session")
	s.services.SessionService.WriteSession(w, r, services.SessionPayload{UserId: int(results.ID), Email: user.Email})
	fmt.Printf("Done writing session, redirecterroo")

	htmx.NewResponse().
		Redirect("/loggedin").Write(w)

}

func (s *Server) handleLoginGet(w http.ResponseWriter, r *http.Request) {
	// no errors or anything on initial bits.
	csrf_value := csrfFromRequest(r)

	component := login.LoginScreen(login.LoginScreenProps{})
	base := base.Base(base.BaseData{Body: component, CSRF: csrf_value, Title: "Login"})
	base.Render(r.Context(), w)

}

// func (s *Server) handleMessageCountGet(c echo.Context) error {

// 	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
// 	c.Response().Header().Set("Cache-Control", "no-cache")
// 	c.Response().Header().Set("Connection", "keep-alive")
// 	c.Response().Header().Set("Content-Type", "text/event-stream")
// 	component := views.MessageCount(13)
// 	c.Response().Writer.WriteHeader(200)
// 	templ.Handler(component).ServeHTTP(c.Response().Writer, c.Request())
// 	// renderComponent(component, c)
// 	return nil
// }

func (s *Server) handleRegisterGet(w http.ResponseWriter, r *http.Request) {
	// no errors or anything on initial bits.
	component := register.RegisterForm(register.RegisterFormData{})
	base := base.Base(base.BaseData{Body: component, CSRF: csrfFromRequest(r), Title: "Register"})
	htmx.NewResponse().RenderTempl(r.Context(), w, base)

}

func (s *Server) handleRegisterPost(w http.ResponseWriter, r *http.Request) {
	var reg dto.RegisterDTO

	var formErrors validation.Errors

	if err := render.Bind(r, &reg); err != nil {
		base.InternalServerError(err).Render(r.Context(), w)
		return

	}

	if err := reg.Validate(); err != nil {
		formErrors = err.(validation.Errors)
	} else {
		// also check that password  === confirm password
		fmt.Printf("%q - %q\n", reg.Password, reg.ConfirmPassword)
		if reg.Password != reg.ConfirmPassword {
			// err.
			formErrors = map[string]error{
				"confirm_password": validation.NewError("", "The confirmed password must match your password"),
			}
		}
	}

	if formErrors.Filter() != nil {
		component := register.RegisterForm(register.RegisterFormData{Previous: reg, Errors: formErrors})
		htmx.NewResponse().
			Retarget("#registerForm").
			Reswap(htmx.SwapOuterHTML).
			RenderTempl(r.Context(), w, component)
		return

	}
	// if there are form errors, don;t do the following!
	// let's hash our password + then check and see if the user already exists (we hash first to prevent timing attacks)
	user, err := s.services.AuthenticationService.Register(reg)
	if err != nil {
		s.logger.Error("Already exists", err)
		// actually check if it exists or not
		formErrors = map[string]error{
			"email": validation.NewError("", "That user already exists "),
		}
		component := register.RegisterForm(register.RegisterFormData{Previous: reg, Errors: formErrors})
		htmx.NewResponse().
			Retarget("#registerForm").
			Reswap(htmx.SwapOuterHTML).
			RenderTempl(r.Context(), w, component)
		return
	}
	// ok - return success!
	s.logger.Info("Successful registration!")
	s.services.SessionService.WriteSession(w, r, services.SessionPayload{UserId: int(user.ID), Email: user.Email})

	htmx.NewResponse().Redirect("/loggedin").Write(w)

}

func (s *Server) handleLoggedInGet(w http.ResponseWriter, r *http.Request) {
	// kill session + redirect (i should not need to post anywhere)
	// write a blank session
	userId := s.getUserIdFromCTX(r)

	component := loggedin.LoggedInView(userId)
	base := base.Base(base.BaseData{Body: component, CSRF: csrfFromRequest(r), Title: "Register"})
	htmx.NewResponse().
		RenderTempl(r.Context(), w, base)

}
