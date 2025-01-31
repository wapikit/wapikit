package event_controller

import (
	"encoding/json"
	"net/http"

	// . "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	controller "github.com/wapikit/wapikit/api/controllers"
	"github.com/wapikit/wapikit/interfaces"
)

type EventController struct {
	controller.BaseController `json:"-,inline"`
}

func NewEventController() *EventController {
	return &EventController{
		BaseController: controller.BaseController{
			Name:        "Event Controller",
			RestApiPath: "/api/events",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/events",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleEventsSubscription),
					IsAuthorizationRequired: true,
				},
			},
		},
	}
}

func handleEventsSubscription(context interfaces.ContextWithSession) error {
	logger := context.App.Logger
	eventService := context.App.EventService

	userUuid, err := uuid.Parse(context.Session.User.UniqueId)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Invalid user id found")
	}

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Invalid organization id found")
	}

	logger.Info("User %v subscribed to events", userUuid, nil)
	logger.Info("Organization %v subscribed to events", orgUuid, nil)

	eventChannel := eventService.HandleApiServerEvents(context.Request().Context())

	if err != nil {
		logger.Error("Error inserting ai message: %v", err.Error(), nil)
		return context.JSON(http.StatusInternalServerError, "Error inserting ai message to database")
	}

	context.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)
	context.Response().WriteHeader(http.StatusOK)
	enc := json.NewEncoder(context.Response())

	for response := range eventChannel {
		if err := enc.Encode(response); err != nil {
			return err
		}
		context.Response().Flush()
	}

	return nil

}
