package services

import "main/internal/sse"

type ServiceLocator struct {
	// each service will be here, but it will also
	// be able to reference each other.

	SSEEventBus *sse.Handler

	SessionService        ISessionService
	AuthenticationService *AuthenticationService
	ChatService           *ChatService

	// add our chat service
}
