package analytics_service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/wapikit/wapikit/.db-generated/model"
	"github.com/wapikit/wapikit/.db-generated/table"
	"github.com/wapikit/wapikit/api/services"
	"github.com/wapikit/wapikit/internal/api_types"
	"github.com/wapikit/wapikit/internal/core/utils"
	"github.com/wapikit/wapikit/internal/interfaces"

	. "github.com/go-jet/jet/v2/postgres"
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
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetPrimaryAnalytics,
						},
					},
				},
				{
					Path:                    "/api/analytics/secondary",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleSecondaryAnalyticsDashboardData),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetSecondaryAnalytics,
						},
					},
				},
				{
					Path:                    "/api/analytics/campaigns/:campaignId",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetCampaignAnalyticsById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetCampaignAnalytics,
						},
					},
				},
				{
					Path:                    "/api/analytics/campaigns",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetCampaignAnalytics),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetCampaignAnalytics,
						},
					},
				},
			},
		},
	}
}

func handlePrimaryAnalyticsDashboardData(context interfaces.ContextWithSession) error {
	params := new(api_types.GetPrimaryAnalyticsParams)
	err := utils.BindQueryParams(context, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	minDateRange := params.From
	maxDateRange := params.To

	if minDateRange.IsZero() || maxDateRange.IsZero() {
		return context.JSON(http.StatusBadRequest, "Invalid date range")
	}

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Invalid organization id")
	}

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

	var aggregateAnalyticsData api_types.AggregateAnalyticsSchema
	var linkDataSeries []api_types.LinkClicksGraphDataPointSchema
	var messageDataSeries []api_types.MessageAnalyticGraphDataPointSchema

	CampaignStatsCte := CTE("campaignStats")
	ContactStatsCte := CTE("contactStats")
	ConversationStatsCte := CTE("ConversationStats")
	MessageStatsCte := CTE("MessageStats")

	aggregateStatsQuery := WITH(
		CampaignStatsCte.AS(
			SELECT(
				COALESCE(
					SUM(CASE().
						WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Draft.String()))).
						THEN(CAST(CAST(Int(1)).AS_INTEGER()).AS_INTEGER()).
						ELSE(CAST(CAST(Int(0)).AS_INTEGER()).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("draft"),
				COALESCE(
					SUM(CASE().
						WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Cancelled.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("cancelled"),
				COALESCE(
					SUM(CASE().
						WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Running.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("running"),
				COALESCE(
					SUM(CASE().
						WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Finished.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("finished"),
				COALESCE(
					SUM(CASE().
						WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Paused.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("paused"),
				COALESCE(
					SUM(CASE().
						WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Scheduled.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("scheduled"),
				COUNT(table.Campaign.UniqueId).AS("total"),
			).FROM(table.Campaign).
				WHERE(table.Campaign.OrganizationId.EQ(UUID(orgUuid))),
		),
		ContactStatsCte.AS(
			SELECT(
				COALESCE(
					SUM(CASE().
						WHEN(table.Contact.Status.EQ(utils.EnumExpression(model.ContactStatus_Active.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("contactStats.active"),
				COALESCE(
					SUM(CASE().
						WHEN(table.Contact.Status.EQ(utils.EnumExpression(model.ContactStatus_Blocked.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("blocked"),
				COUNT(table.Contact.UniqueId).AS("total"),
			).FROM(table.Contact).
				WHERE(table.Contact.OrganizationId.EQ(UUID(orgUuid))),
		),
		ConversationStatsCte.AS(
			SELECT(
				COALESCE(
					SUM(CASE().
						WHEN(table.Conversation.Status.EQ(utils.EnumExpression(model.ConversationStatus_Active.String()))).THEN(CAST(Int(1)).AS_INTEGER()).ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("active"),
				COALESCE(
					SUM(CASE().
						WHEN(table.Conversation.Status.EQ(utils.EnumExpression(model.ConversationStatus_Closed.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("closed"),
				// SUM(CASE(table.Conversation.Status.EQ(String(string(api_types.conversa)))).THEN(CAST(Int(1)).AS_INTEGER()).ELSE(CAST(Int(0)).AS_INTEGER())).AS("pending"), // here we need to determine if there are unread incoming messages for a conversation exists then SUM it in pending
				COUNT(table.Conversation.UniqueId).AS("total"),
			).FROM(table.Conversation).
				WHERE(table.Conversation.OrganizationId.EQ(UUID(orgUuid))),
		),
		MessageStatsCte.AS(
			SELECT(
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatus_Delivered.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("delivered"),
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatus_Failed.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("failed"),
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatus_Read.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("read"),
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatus_Sent.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("sent"),
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatus_UnDelivered.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("undelivered"),
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatus_Read.String()))).
						THEN(CAST(Int(0)).AS_INTEGER()).
						ELSE(CAST(Int(1)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("unread"),
				COUNT(table.Message.UniqueId).AS("total"),
			).FROM(table.Message).
				WHERE(table.Message.OrganizationId.EQ(UUID(orgUuid))),
		),
	)(SELECT(
		CampaignStatsCte.AllColumns().As("aggregate_analytics_schema.campaignStats"),
		ContactStatsCte.AllColumns().As("aggregate_analytics_schema.contactStats"),
		ConversationStatsCte.AllColumns().As("aggregate_analytics_schema.conversationStats"),
		MessageStatsCte.AllColumns().As("aggregate_analytics_schema.messageStats"),
	).FROM(
		CampaignStatsCte,
		ContactStatsCte,
		ConversationStatsCte,
		MessageStatsCte,
	))

	// ! TODO: debug the above sql query, it returns 0 only, even the db has data
	// stringDebuggingSql := aggregateStatsQuery.DebugSql()
	// stringSql, _ := aggregateStatsQuery.Sql()
	// fmt.Println("stringDebuggingSql is", stringDebuggingSql)
	// fmt.Println("stringSql is", stringSql)

	linkDataQuery := SELECT(
		table.TrackLinkClick.CreatedAt.AS("date"),
		COALESCE(COUNT(table.TrackLinkClick.UniqueId), Int(0)).AS("count"),
		TO_CHAR(table.TrackLinkClick.CreatedAt, String("DD-MM-YYYY")).AS("label"),
	).
		FROM(table.TrackLinkClick.
			LEFT_JOIN(table.TrackLink, table.TrackLinkClick.TrackLinkId.EQ(table.TrackLink.UniqueId))).
		WHERE(table.TrackLink.OrganizationId.EQ(UUID(orgUuid)).
			AND(table.TrackLinkClick.CreatedAt.
				BETWEEN(
					TimestampzExp(Timestamp(minDateRange.Year(), minDateRange.Month(), minDateRange.Day(), minDateRange.Hour(), minDateRange.Minute(), minDateRange.Second())),
					TimestampzExp(Timestamp(maxDateRange.Year(), maxDateRange.Month(), maxDateRange.Day(), maxDateRange.Hour(), maxDateRange.Minute(), maxDateRange.Second())),
				),
			)).
		GROUP_BY(table.TrackLinkClick.CreatedAt).
		ORDER_BY(table.TrackLinkClick.CreatedAt)

	messageDataQuery := SELECT(
		table.Message.CreatedAt.AS("date"),
		COALESCE(
			SUM(CASE().WHEN(table.Message.Direction.EQ(utils.EnumExpression(model.MessageDirection_OutBound.String()))).
				THEN(CAST(Int(1)).AS_INTEGER()).
				ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("sent"),
		COALESCE(
			SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatus_Read.String()))).
				THEN(CAST(Int(1)).AS_INTEGER()).
				ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("read"),
		// ! TODO: for replied get messages for which a reply message exists, means count a record which has a related record in the same table where the relatedRecord.RepliedToId is the current record's UniqueId
		// COALESCE(
		// 	SUM(CASE().
		// 	WHEN(SELECT()).
		// 	THEN(CAST(Int(1)).AS_INTEGER()).ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER(), Int(0))
		// 	.AS("replied"),
		// ),
		TO_CHAR(table.Message.CreatedAt, String("DD-MM-YYYY")).AS("label"),
	).FROM(
		table.Message,
	).WHERE(
		table.Message.OrganizationId.EQ(UUID(orgUuid)).
			AND(table.Message.CreatedAt.
				BETWEEN(
					TimestampzExp(Timestamp(minDateRange.Year(), minDateRange.Month(), minDateRange.Day(), minDateRange.Hour(), minDateRange.Minute(), minDateRange.Second())),
					TimestampzExp(Timestamp(maxDateRange.Year(), maxDateRange.Month(), maxDateRange.Day(), maxDateRange.Hour(), maxDateRange.Minute(), maxDateRange.Second())),
				),
			),
	).GROUP_BY(
		table.Message.CreatedAt,
	).ORDER_BY(
		table.Message.CreatedAt,
	)

	err = aggregateStatsQuery.QueryContext(context.Request().Context(), context.App.Db, &aggregateAnalyticsData)

	if err != nil {
		fmt.Println("error is", err.Error())
		if err.Error() == "no rows in result set" {
			// do nothing keep the empty response as defined above in the controller
		} else {
			return context.JSON(http.StatusInternalServerError, "Error getting aggregate stats")
		}
	}

	responseToReturn.AggregateAnalytics = aggregateAnalyticsData

	fmt.Println("aggregateAnalyticsData query results are", aggregateAnalyticsData)

	err = linkDataQuery.QueryContext(context.Request().Context(), context.App.Db, &linkDataSeries)

	if err != nil {
		fmt.Println("error is", err.Error())
		if err.Error() == "no rows in result set" {
			// do nothing keep the empty response as defined above in the controller
		} else {
			return context.JSON(http.StatusInternalServerError, "Error getting link clicks")
		}
	}

	if len(linkDataSeries) != 0 {
		responseToReturn.LinkClickAnalytics = linkDataSeries
	}

	err = messageDataQuery.QueryContext(context.Request().Context(), context.App.Db, &messageDataSeries)

	if err != nil {
		fmt.Println("error is", err.Error())
		if err.Error() == "no rows in result set" {
			// do nothing keep the empty response as defined above in the controller
		} else {
			return context.JSON(http.StatusInternalServerError, "Error getting message analytics")
		}
	}

	if len(messageDataSeries) != 0 {
		responseToReturn.MessageAnalytics = messageDataSeries
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

func _getAggregateDashboardStats(organizationId string, from, to time.Time) {

}

func _getLinkClicks(organizationId string, from, to time.Time) {

}
