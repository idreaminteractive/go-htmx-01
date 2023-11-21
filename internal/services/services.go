package services

type ServiceLocator struct {
	// each service will be here, but it will also
	// be able to reference each other.

	SessionService        ISessionService
	AuthenticationService *AuthenticationService

	// add our chat service
}
