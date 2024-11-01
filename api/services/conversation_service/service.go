package conversation_service

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/internal/api_types"
	"github.com/sarthakjdev/wapikit/internal/core/api_server_events"
	"github.com/sarthakjdev/wapikit/internal/core/utils"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

type ConversationService struct {
	services.BaseService `json:"-,inline"`
}

func NewConversationService() *ConversationService {
	return &ConversationService{
		BaseService: services.BaseService{
			Name:        "Conversation Service",
			RestApiPath: "/api/conversation",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/conversation",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetConversations),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    100,
							WindowTimeInMs: time.Hour.Milliseconds(),
						},
					},
				},
				{
					Path:                    "/api/conversation/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetConversationById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    100,
							WindowTimeInMs: time.Hour.Milliseconds(),
						},
					},
				},
				{
					Path:                    "/api/conversation/:id",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(handleUpdateConversationById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    100,
							WindowTimeInMs: time.Hour.Milliseconds(),
						},
					},
				},
				{
					Path:                    "/api/conversation/:id",
					Method:                  http.MethodDelete,
					Handler:                 interfaces.HandlerWithSession(handleDeleteConversationById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    100,
							WindowTimeInMs: time.Hour.Milliseconds(),
						},
					},
				},
				{
					Path:                    "/api/conversation/:id/assign",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleAssignConversation),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    100,
							WindowTimeInMs: time.Hour.Milliseconds(),
						},
					},
				},
				{
					Path:                    "/api/conversation/:id/unassign",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleUnassignConversation),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    100,
							WindowTimeInMs: time.Hour.Milliseconds(),
						},
					},
				},
				{
					Path:                    "/api/conversation/:id/messages",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetConversationMessages),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    100,
							WindowTimeInMs: time.Hour.Milliseconds(),
						},
					},
				},
			},
		},
	}
}

func handleGetConversations(context interfaces.ContextWithSession) error {
	queryParams := new(api_types.GetConversationsParams)

	if err := utils.BindQueryParams(context, &queryParams); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func handleGetConversationById(context interfaces.ContextWithSession) error {
	conversationId := context.Param("id")
	if conversationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "conversation id is required")
	}
	conversationUuid, err := uuid.Parse(conversationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid conversation id")
	}

	context.App.Logger.Info("conversation id: %v", conversationUuid)

	return nil
}

func handleUpdateConversationById(context interfaces.ContextWithSession) error {
	conversationId := context.Param("id")
	if conversationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "conversation id is required")
	}
	conversationUuid, err := uuid.Parse(conversationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid conversation id")
	}

	context.App.Logger.Info("conversation id: %v", conversationUuid)

	return nil
}

func handleDeleteConversationById(context interfaces.ContextWithSession) error {
	conversationId := context.Param("id")
	if conversationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "conversation id is required")
	}
	conversationUuid, err := uuid.Parse(conversationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid conversation id")
	}

	context.App.Logger.Info("conversation id: %v", conversationUuid)

	return nil
}

func handleGetConversationMessages(context interfaces.ContextWithSession) error {
	conversationId := context.Param("id")
	if conversationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "conversation id is required")
	}
	conversationUuid, err := uuid.Parse(conversationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid conversation id")
	}

	context.App.Logger.Info("conversation id: %v", conversationUuid)

	return nil
}

func handleAssignConversation(context interfaces.ContextWithSession) error {

	redis := context.App.Redis

	event := api_server_events.BaseApiServerEvent{
		EventType: api_server_events.ApiServerChatAssignmentEvent,
	}

	redis.PublishMessageToRedisChannel(context.App.Constants.RedisEventChannelName, string(event.ToJson()))

	return nil
}

func handleUnassignConversation(context interfaces.ContextWithSession) error {

	redis := context.App.Redis

	event := api_server_events.BaseApiServerEvent{
		EventType: api_server_events.ApiServerChatUnAssignmentEvent,
	}

	redis.PublishMessageToRedisChannel(context.App.Constants.RedisEventChannelName, string(event.ToJson()))

	return nil
}
