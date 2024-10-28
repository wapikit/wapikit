package analytics_service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/internal/api_types"
	"github.com/sarthakjdev/wapikit/internal/core/utils"
	"github.com/sarthakjdev/wapikit/internal/interfaces"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/sarthakjdev/wapikit/.db-generated/model"
	table "github.com/sarthakjdev/wapikit/.db-generated/table"
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
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handlePrimaryAnalyticsDashboardData),
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/analytics/secondary",
					Method:                  http.MethodGet,
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
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetCampaignAnalytics),
					IsAuthorizationRequired: true,
				},
			},
		},
	}
}

func handlePrimaryAnalyticsDashboardData(context interfaces.ContextWithSession) error {

	fmt.Println("Primary analytics dashboard data")

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Invalid organization id")
	}

	var aggregateAnalyticsData api_types.AggregateAnalyticsSchema

	CampaignStatsCte := CTE("CampaignStats")
	ContactStatsCte := CTE("ContactStats")
	// ConversationStatsCte := CTE("ConversationStats")
	// MessageStatsCte := CTE("MessageStats")

	aggregateStatsQuery := WITH(
		CampaignStatsCte.AS(
			SELECT(
				COUNT(CASE().WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Draft.String()))).THEN(Bool(true))).AS("draft"),
				COUNT(CASE().WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Cancelled.String()))).THEN(Bool(true))).AS("cancelled"),
				COUNT(CASE().WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Running.String()))).THEN(Bool(true))).AS("running"),
				COUNT(CASE().WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Finished.String()))).THEN(Bool(true))).AS("finished"),
				COUNT(CASE().WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Paused.String()))).THEN(Bool(true))).AS("paused"),
				COUNT(CASE().WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Scheduled.String()))).THEN(Bool(true))).AS("scheduled"),
			).FROM(table.Campaign).WHERE(table.Campaign.OrganizationId.EQ(UUID(orgUuid))),
		),
		ContactStatsCte.AS(
			SELECT(
				COUNT(CASE().WHEN(table.Contact.Status.EQ(utils.EnumExpression(model.ContactStatus_Active.String()))).THEN(Bool(true))).AS("active"),
				COUNT(CASE().WHEN(table.Contact.Status.EQ(utils.EnumExpression(model.ContactStatus_Blocked.String()))).THEN(Bool(true))).AS("blocked"),
			).FROM(table.Contact).WHERE(table.Contact.OrganizationId.EQ(UUID(orgUuid))),
		),
		// ConversationStatsCte.AS(
		// 	SELECT(
		// 		COUNT(CASE(table.Conversation.Status.EQ(String(string(api_types.ConversationStatusEnumActive)))).THEN(Bool(true))).AS("active"),
		// 		COUNT(CASE(table.Conversation.Status.EQ(String(string(api_types.ConversationStatusEnumClosed)))).THEN(Bool(true))).AS("closed"),
		// 		// COUNT(CASE(table.Conversation.Status.EQ(String(string(api_types.conversa)))).THEN(Bool(true))).AS("pending"), // here we need to determine if there are unread incoming messages for a conversation exists then count it in pending
		// 		COUNT(table.Conversation.UniqueId).AS("total"),
		// 	).FROM(table.Conversation).WHERE(table.Conversation.OrganizationId.EQ(UUID(orgUuid))),
		// ),
		// MessageStatsCte.AS(
		// 	SELECT(
		// 		COUNT(CASE(table.Message.Status.EQ(String(string(api_types.MessageStatusEnumDelivered)))).THEN(Bool(true))).AS("delivered"),
		// 		COUNT(CASE(table.Message.Status.EQ(String(string(api_types.MessageStatusEnumFailed)))).THEN(Bool(true))).AS("failed"),
		// 		COUNT(CASE(table.Message.Status.EQ(String(string(api_types.MessageStatusEnumRead)))).THEN(Bool(true))).AS("read"),
		// 		COUNT(CASE(table.Message.Status.EQ(String(string(api_types.MessageStatusEnumSent)))).THEN(Bool(true))).AS("sent"),
		// 		COUNT(CASE(table.Message.Status.EQ(String(string(api_types.MessageStatusEnumUnDelivered)))).THEN(Bool(true))).AS("undelivered"),
		// 		COUNT(CASE(table.Message.Status.EQ(String(string(api_types.MessageStatusEnumRead)))).THEN(Int(0)).ELSE(Bool(true))).AS("unread"),
		// 		COUNT(table.Message.UniqueId).AS("total"),
		// 	).FROM(table.Message).WHERE(table.Message.OrganizationId.EQ(UUID(orgUuid))),
		// ),
	)(SELECT(
		CampaignStatsCte.AllColumns(),
		ContactStatsCte.AllColumns(),
	).FROM(
		CampaignStatsCte,
		ContactStatsCte,
	))

	stringSql := aggregateStatsQuery.DebugSql()

	fmt.Println("sql is", stringSql)

	err = aggregateStatsQuery.QueryContext(context.Request().Context(), context.App.Db, &aggregateAnalyticsData)

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
				Active:  0,
				Blocked: 0,
				Total:   0,
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

	fmt.Println("query results are", aggregateAnalyticsData)

	if err != nil {
		fmt.Println("error is", err.Error())
		if err.Error() == "no rows in result set" {
			return context.JSON(http.StatusOK, responseToReturn)
		}
		return context.JSON(http.StatusInternalServerError, "Error getting aggregate stats")
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
	responseToReturn := api_types.CampaignAnalyticsResponseSchema{
		MessagesDelivered:   0,
		MessagesFailed:      0,
		MessagesRead:        0,
		MessagesSent:        0,
		MessagesUndelivered: 0,
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func handleGetCampaignAnalytics(context interfaces.ContextWithSession) error {
	responseToReturn := api_types.CampaignAnalyticsResponseSchema{
		MessagesDelivered:   0,
		MessagesFailed:      0,
		MessagesRead:        0,
		MessagesSent:        0,
		MessagesUndelivered: 0,
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func _getMessagingStats(organizationId string, from, to time.Time) {

}

func _getConversationAnalytics(organizationId string, from, to time.Time) {

}

func _getAggregateDashboardStats(organizationId string, from, to time.Time) {

}

func _getLinkClicks(organizationId string, from, to time.Time) {

}
