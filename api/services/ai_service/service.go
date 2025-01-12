package auth_service

import (
	"net/http"

	"github.com/wapikit/wapikit/api/services"
	"github.com/wapikit/wapikit/internal/interfaces"

	. "github.com/go-jet/jet/v2/postgres"
)

type AiService struct {
	services.BaseService `json:"-,inline"`
}

func NewAiService() *AiService {
	return &AiService{
		BaseService: services.BaseService{
			Name:        "AI Service",
			RestApiPath: "/api/ai",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/ai/completion",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(chatCompletion),
					IsAuthorizationRequired: false,
				},
				{
					Path:                    "/api/ai/conversations",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getAiConversations),
					IsAuthorizationRequired: false,
				},
			},
		},
	}
}

func chatCompletion(context interfaces.ContextWithSession) error {
	return nil
}

func getAiConversations(context interfaces.ContextWithSession) error {
	return nil

}
