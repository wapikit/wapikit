package conversation_service

import (
	"github.com/sarthakjdev/wapikit/api/services"
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
				{},
			},
		},
	}
}
