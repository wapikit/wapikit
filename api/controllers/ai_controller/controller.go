package analytics_service

import (
	"net/http"

	controller "github.com/wapikit/wapikit/api/controllers"
	"github.com/wapikit/wapikit/internal/api_types"
	"github.com/wapikit/wapikit/internal/interfaces"
)

type IntegrationController struct {
	controller.BaseController `json:"-,inline"`
}

func NewIntegrationService() *IntegrationController {
	return &IntegrationController{
		BaseController: controller.BaseController{
			Name:        "AI Service",
			RestApiPath: "/api/ai",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/ai/chats",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetChats),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
				{
					Path:                    "/api/ai/chat",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(handleReplyToChat),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
			},
		},
	}
}

func handleGetChats(context interfaces.ContextWithSession) error {
	return nil

}

func handleReplyToChat(context interfaces.ContextWithSession) error {
	return nil
}
