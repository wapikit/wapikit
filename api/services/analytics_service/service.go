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
					Handler:                 handleGetAggregateDashboardStats,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/analytics/getConversationStats",
					Method:                  http.MethodPost,
					Handler:                 handleGetConversationStats,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/analytics/getMessageStats",
					Method:                  http.MethodPost,
					Handler:                 handleGetMessagingStats,
					IsAuthorizationRequired: true,
				},
			},
		},
	}
}

func handleGetConversationStats(context interfaces.CustomContext) error {
	return nil
}

func handleGetMessagingStats(context interfaces.CustomContext) error {
	return nil
}

func handleGetAggregateDashboardStats(context interfaces.CustomContext) error {
	return nil
}
