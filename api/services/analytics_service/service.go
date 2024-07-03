package analytics_service

import (
	"net/http"

	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

type AnalyticsService struct {
	services.BaseService `json:"-,inline"`
}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{
		BaseService: services.BaseService{
			Name:        "Analytics Service",
			RestApiPath: "/api/analytics",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/analytics/getAggregateDashboardStats",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(handleGetAggregateDashboardStats),
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/analytics/getConversationStats",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(handleGetConversationStats),
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/analytics/getMessageStats",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(handleGetMessagingStats),
					IsAuthorizationRequired: true,
				},
			},
		},
	}
}

func handleGetConversationStats(context interfaces.ContextWithSession) error {
	return nil
}

func handleGetMessagingStats(context interfaces.ContextWithSession) error {
	return nil
}

func handleGetAggregateDashboardStats(context interfaces.ContextWithSession) error {
	return nil
}
