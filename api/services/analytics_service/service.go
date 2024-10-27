package analytics_service

import (
	"net/http"
	"time"

	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/internal/api_types"
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
					Path:                    "/api/analytics/primary",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(handlePrimaryAnalyticsDashboardData),
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/analytics/secondary",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(handleSecondaryAnalyticsDashboardData),
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/analytics/campaigns/:campaignId",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetCampaignAnalyticsById),
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/analytics/campaigns",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(handleGetCampaignAnalytics),
					IsAuthorizationRequired: true,
				},
			},
		},
	}
}

func handlePrimaryAnalyticsDashboardData(context interfaces.ContextWithSession) error {

	responseToReturn := api_types.PrimaryAnalyticsResponseSchema{
		AggregateAnalytics: api_types.AggregateAnalyticsSchema{
			CampaignStats: api_types.AggregateCampaignStatsDataPointsSchema{
				Cancelled: 0,
				Running:   0,
				Draft:     0,
				Finished:  0,
				Paused:    0,
				Scheduled: 0,
			},
			ContactStats: api_types.AggregateContactStatsDataPointsSchema{
				Active:   0,
				Blocked:  0,
				Inactive: 0,
				Total:    0,
			},
			ConversationStats: api_types.AggregateConversationStatsDataPointsSchema{
				Active:  0,
				Closed:  0,
				Pending: 0,
				Total:   0,
			},
			MessageStats: api_types.AggregateMessageStatsDataPointsSchema{
				Delivered:   0,
				Failed:      0,
				Read:        0,
				Sent:        0,
				Total:       0,
				Undelivered: 0,
				Unread:      0,
			},
		},
		LinkClickAnalytics: []api_types.LinkClicksGraphDataPointSchema{},
		MessageAnalytics:   []api_types.MessageAnalyticGraphDataPointSchema{},
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func handleSecondaryAnalyticsDashboardData(context interfaces.ContextWithSession) error {

	responseToReturn := api_types.SecondaryAnalyticsDashboardResponseSchema{
		ConversationsAnalytics:                  []api_types.ConversationAnalyticsDataPointSchema{},
		MessageTypeTrafficDistributionAnalytics: []api_types.MessageTypeDistributionGraphDataPointSchema{},
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func handleGetCampaignAnalyticsById(context interfaces.ContextWithSession) error {
	return nil
}

func handleGetCampaignAnalytics(context interfaces.ContextWithSession) error {
	return nil
}

func _getMessagingStats(organizationId string, from, to time.Time) {

}

func _getConversationAnalytics(organizationId string, from, to time.Time) {

}

func _getAggregateDashboardStats(organizationId string, from, to time.Time) {

}

func _getLinkClicks(organizationId string, from, to time.Time) {

}
