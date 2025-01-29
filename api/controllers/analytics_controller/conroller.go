package analytics_controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/wapikit/wapikit/.db-generated/model"
	"github.com/wapikit/wapikit/.db-generated/table"
	"github.com/wapikit/wapikit/api/api_types"
	controller "github.com/wapikit/wapikit/api/controllers"
	"github.com/wapikit/wapikit/interfaces"
	"github.com/wapikit/wapikit/utils"

	"github.com/go-jet/jet/qrm"
	. "github.com/go-jet/jet/v2/postgres"
)

type AnalyticsController struct {
	controller.BaseController `json:"-,inline"`
}

func NewAnalyticsController() *AnalyticsController {
	return &AnalyticsController{
		BaseController: controller.BaseController{
			Name:        "Analytics Controller",
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
					Path:                    "/api/analytics/campaign/:campaignId",
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
			},
		},
	}
}

func handlePrimaryAnalyticsDashboardData(context interfaces.ContextWithSession) error {
	params := new(api_types.GetPrimaryAnalyticsParams)
	err := utils.BindQueryParams(context, params)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
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
				CampaignsCancelled: 0,
				CampaignsRunning:   0,
				CampaignsDraft:     0,
				CampaignsFinished:  0,
				CampaignsPaused:    0,
				CampaignsScheduled: 0,
			},
			ContactStats: api_types.AggregateContactStatsDataPointsSchema{
				ContactsActive:  0,
				ContactsBlocked: 0,
				TotalContacts:   0,
			},
			ConversationStats: api_types.AggregateConversationStatsDataPointsSchema{
				ConversationsActive:  0,
				ConversationsClosed:  0,
				ConversationsPending: 0,
				TotalConversations:   0,
			},
			MessageStats: api_types.AggregateMessageStatsDataPointsSchema{
				MessagesDelivered:   0,
				MessagesFailed:      0,
				MessagesRead:        0,
				MessagesSent:        0,
				TotalMessages:       0,
				MessagesUndelivered: 0,
				MessagesUnread:      0,
			},
		},
		LinkClickAnalytics: []api_types.LinkClicksGraphDataPointSchema{},
		MessageAnalytics:   []api_types.MessageAnalyticGraphDataPointSchema{},
	}

	var primaryAnalyticsData struct {
		CampaignsCancelled   int `json:"campaignsCancelled"`
		CampaignsDraft       int `json:"campaignsDraft"`
		CampaignsFinished    int `json:"campaignsFinished"`
		CampaignsPaused      int `json:"campaignsPaused"`
		CampaignsRunning     int `json:"campaignsRunning"`
		CampaignsScheduled   int `json:"campaignsScheduled"`
		TotalCampaigns       int `json:"totalCampaigns"`
		ContactsActive       int `json:"contactsActive"`
		ContactsBlocked      int `json:"contactsBlocked"`
		TotalContacts        int `json:"totalContacts"`
		ConversationsActive  int `json:"conversationsActive"`
		ConversationsClosed  int `json:"conversationsClosed"`
		ConversationsPending int `json:"conversationsPending"`
		TotalConversations   int `json:"totalConversations"`
		MessagesDelivered    int `json:"messagesDelivered"`
		MessagesFailed       int `json:"messagesFailed"`
		MessagesRead         int `json:"messagesRead"`
		MessagesSent         int `json:"messagesSent"`
		MessagesUndelivered  int `json:"messagesUndelivered"`
		MessagesUnread       int `json:"messagesUnread"`
		TotalMessages        int `json:"totalMessages"`
	}

	CampaignStatsCte := CTE("campaignStats")
	ContactStatsCte := CTE("contactStats")
	ConversationStatsCte := CTE("conversationStats")
	MessageStatsCte := CTE("messageStats")

	aggregateStatsQuery := WITH(
		CampaignStatsCte.AS(
			SELECT(
				COALESCE(
					SUM(CASE().
						WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatusEnum_Draft.String()))).
						THEN(CAST(CAST(Int(1)).AS_INTEGER()).AS_INTEGER()).
						ELSE(CAST(CAST(Int(0)).AS_INTEGER()).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("campaignsDraft"),
				COALESCE(
					SUM(CASE().
						WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatusEnum_Cancelled.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("campaignsCancelled"),
				COALESCE(
					SUM(CASE().
						WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatusEnum_Running.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("campaignsRunning"),
				COALESCE(
					SUM(CASE().
						WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatusEnum_Finished.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("campaignsFinished"),
				COALESCE(
					SUM(CASE().
						WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatusEnum_Paused.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("campaignsPaused"),
				COALESCE(
					SUM(CASE().
						WHEN(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatusEnum_Scheduled.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("campaignsScheduled"),
				COUNT(table.Campaign.UniqueId).AS("totalCampaigns"),
			).FROM(table.Campaign).
				WHERE(table.Campaign.OrganizationId.EQ(UUID(orgUuid))),
		),
		ContactStatsCte.AS(
			SELECT(
				COALESCE(
					SUM(CASE().
						WHEN(table.Contact.Status.EQ(utils.EnumExpression(model.ContactStatusEnum_Active.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("contactsActive"),
				COALESCE(
					SUM(CASE().
						WHEN(table.Contact.Status.EQ(utils.EnumExpression(model.ContactStatusEnum_Blocked.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("contactsBlocked"),
				COUNT(table.Contact.UniqueId).AS("totalContacts"),
			).FROM(table.Contact).
				WHERE(table.Contact.OrganizationId.EQ(UUID(orgUuid))),
		),
		ConversationStatsCte.AS(
			SELECT(
				COALESCE(
					SUM(CASE().
						WHEN(table.Conversation.Status.EQ(utils.EnumExpression(model.ConversationStatusEnum_Active.String()))).THEN(CAST(Int(1)).AS_INTEGER()).ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("conversationsActive"),
				COALESCE(
					SUM(CASE().
						WHEN(table.Conversation.Status.EQ(utils.EnumExpression(model.ConversationStatusEnum_Closed.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("conversationsClosed"),
				// SUM(CASE(table.Conversation.Status.EQ(String(string(api_types.conversa)))).THEN(CAST(Int(1)).AS_INTEGER()).ELSE(CAST(Int(0)).AS_INTEGER())).AS("pending"), // here we need to determine if there are unread incoming messages for a conversation exists then SUM it in pending
				COUNT(table.Conversation.UniqueId).AS("totalConversations"),
			).FROM(table.Conversation).
				WHERE(table.Conversation.OrganizationId.EQ(UUID(orgUuid))),
		),
		MessageStatsCte.AS(
			SELECT(
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatusEnum_Delivered.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("messagesDelivered"),
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatusEnum_Failed.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("messagesFailed"),
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatusEnum_Read.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("messagesRead"),
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatusEnum_Sent.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("messagesSent"),
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatusEnum_UnDelivered.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("messagesUndelivered"),
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatusEnum_Read.String()))).
						THEN(CAST(Int(0)).AS_INTEGER()).
						ELSE(CAST(Int(1)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("messagesUnread"),
				COUNT(table.Message.UniqueId).AS("totalMessages"),
			).FROM(table.Message).
				WHERE(table.Message.OrganizationId.EQ(UUID(orgUuid))),
		),
	)(SELECT(
		CampaignStatsCte.AllColumns(),
		ContactStatsCte.AllColumns(),
		ConversationStatsCte.AllColumns(),
		MessageStatsCte.AllColumns(),
	).FROM(
		CampaignStatsCte,
		ContactStatsCte,
		ConversationStatsCte,
		MessageStatsCte,
	))

	linkClickDataQuery := SELECT(
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

	messageDataSeriesQuery := SELECT(
		table.Message.CreatedAt.AS("date"),
		COALESCE(
			SUM(CASE().WHEN(table.Message.Direction.EQ(utils.EnumExpression(model.MessageDirectionEnum_OutBound.String()))).
				THEN(CAST(Int(1)).AS_INTEGER()).
				ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("sent"),
		COALESCE(
			SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatusEnum_Read.String()))).
				THEN(CAST(Int(1)).AS_INTEGER()).
				ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("read"),
		COALESCE(
			SUM(
				CASE().
					WHEN(
						EXISTS(
							SELECT(table.Message.UniqueId).
								FROM(table.Message).
								WHERE(table.Message.RepliedTo.EQ(table.Message.UniqueId)),
						),
					).
					THEN(CAST(Int(1)).AS_INTEGER()).
					ELSE(CAST(Int(0)).AS_INTEGER()),
			),
			Int(0),
		).AS("replied"),
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

	err = aggregateStatsQuery.QueryContext(context.Request().Context(), context.App.Db, &primaryAnalyticsData)

	if err != nil {
		fmt.Println("error is", err.Error())
		if err.Error() == qrm.ErrNoRows.Error() {
			fmt.Println("No aggregate stats found")
			// do nothing keep the empty response as defined above in the controller
		} else {
			return context.JSON(http.StatusInternalServerError, "Error getting aggregate stats")
		}
	}

	linkDataSeries := []api_types.LinkClicksGraphDataPointSchema{}
	err = linkClickDataQuery.QueryContext(context.Request().Context(), context.App.Db, &linkDataSeries)

	if err != nil {
		fmt.Println("error is", err.Error())
		if err.Error() == qrm.ErrNoRows.Error() {
			fmt.Println("No link click data found")
			// do nothing keep the empty response as defined above in the controller
		} else {
			return context.JSON(http.StatusInternalServerError, "Error getting link click data")
		}
	}

	messageDataSeries := []api_types.MessageAnalyticGraphDataPointSchema{}
	err = messageDataSeriesQuery.QueryContext(context.Request().Context(), context.App.Db, &messageDataSeries)

	if err != nil {
		fmt.Println("error is", err.Error())
		if err.Error() == qrm.ErrNoRows.Error() {
			fmt.Println("No message data found")
			// do nothing keep the empty response as defined above in the controller
		} else {
			return context.JSON(http.StatusInternalServerError, "Error getting message data")
		}
	}

	responseToReturn.AggregateAnalytics = api_types.AggregateAnalyticsSchema{
		CampaignStats: api_types.AggregateCampaignStatsDataPointsSchema{
			CampaignsCancelled: primaryAnalyticsData.CampaignsCancelled,
			CampaignsRunning:   primaryAnalyticsData.CampaignsRunning,
			CampaignsDraft:     primaryAnalyticsData.CampaignsDraft,
			CampaignsFinished:  primaryAnalyticsData.CampaignsFinished,
			CampaignsPaused:    primaryAnalyticsData.CampaignsPaused,
			CampaignsScheduled: primaryAnalyticsData.CampaignsScheduled,
			TotalCampaigns:     primaryAnalyticsData.TotalCampaigns,
		},
		ContactStats: api_types.AggregateContactStatsDataPointsSchema{
			ContactsActive:  primaryAnalyticsData.ContactsActive,
			ContactsBlocked: primaryAnalyticsData.ContactsBlocked,
			TotalContacts:   primaryAnalyticsData.TotalContacts,
		},
		ConversationStats: api_types.AggregateConversationStatsDataPointsSchema{
			ConversationsActive:  primaryAnalyticsData.ConversationsActive,
			ConversationsClosed:  primaryAnalyticsData.ConversationsClosed,
			ConversationsPending: 0,
			TotalConversations:   primaryAnalyticsData.TotalConversations,
		},
		MessageStats: api_types.AggregateMessageStatsDataPointsSchema{
			MessagesDelivered:   primaryAnalyticsData.MessagesDelivered,
			MessagesFailed:      primaryAnalyticsData.MessagesFailed,
			MessagesRead:        primaryAnalyticsData.MessagesRead,
			MessagesSent:        primaryAnalyticsData.MessagesSent,
			TotalMessages:       primaryAnalyticsData.TotalMessages,
			MessagesUndelivered: primaryAnalyticsData.MessagesUndelivered,
			MessagesUnread:      primaryAnalyticsData.MessagesUnread,
		},
	}

	if len(linkDataSeries) != 0 {
		responseToReturn.LinkClickAnalytics = linkDataSeries
	}

	if len(messageDataSeries) != 0 {
		responseToReturn.MessageAnalytics = messageDataSeries
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func handleSecondaryAnalyticsDashboardData(context interfaces.ContextWithSession) error {
	// ! TODO: these analytics we will need once the live team inbox will be implemented

	responseToReturn := api_types.SecondaryAnalyticsDashboardResponseSchema{
		ConversationsAnalytics:                  []api_types.ConversationAnalyticsDataPointSchema{},
		MessageTypeTrafficDistributionAnalytics: []api_types.MessageTypeDistributionGraphDataPointSchema{},
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func handleGetCampaignAnalyticsById(context interfaces.ContextWithSession) error {
	var campaignAnalyticsData struct {
		MessagesDelivered     int                                        `json:"messagesDelivered"`
		MessagesFailed        int                                        `json:"messagesFailed"`
		MessagesRead          int                                        `json:"messagesRead"`
		MessagesSent          int                                        `json:"messagesSent"`
		MessagesUndelivered   int                                        `json:"messagesUndelivered"`
		TotalMessages         int                                        `json:"totalMessages"`
		TotalLinkClicks       int                                        `json:"totalLinkClicks"`
		ConversationInitiated int                                        `json:"conversationInitiated"`
		LinkClicksData        []api_types.LinkClicksGraphDataPointSchema `json:"linkClicksData"`
	}

	campaignMessageDataCte := CTE("messageStats")
	linkClicksCountsCte := CTE("linkStats")
	conversationCountCte := CTE("conversationInitiated")
	linkClicksDataCte := CTE("linkClicksData")

	campaignAnalyticsQuery := WITH(
		campaignMessageDataCte.AS(
			SELECT(
				COUNT(table.Message.UniqueId).AS("totalMessages"),
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatusEnum_Delivered.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("messagesDelivered"),
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatusEnum_Failed.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("messagesFailed"),
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatusEnum_Read.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("messagesRead"),
				COALESCE(

					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatusEnum_Sent.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("messagesSent"),
				COALESCE(
					SUM(CASE().WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatusEnum_UnDelivered.String()))).
						THEN(CAST(Int(1)).AS_INTEGER()).
						ELSE(CAST(Int(0)).AS_INTEGER())), CAST(Int(0)).AS_INTEGER()).AS("messagesUndelivered"),
			).FROM(
				table.Message,
			).WHERE(
				table.Message.CampaignId.EQ(UUID(uuid.MustParse(context.Param("campaignId")))),
			),
		),
		linkClicksCountsCte.AS(
			SELECT(
				COALESCE(
					COUNT(table.TrackLinkClick.UniqueId), CAST(Int(0)).AS_INTEGER()).AS("totalLinkClicks"),
			).FROM(
				table.TrackLinkClick.LEFT_JOIN(table.TrackLink, table.TrackLink.UniqueId.EQ(table.TrackLinkClick.TrackLinkId)),
			).WHERE(
				table.TrackLink.CampaignId.EQ(UUID(uuid.MustParse(context.Param("campaignId")))),
			),
		),
		conversationCountCte.AS(
			SELECT(
				COALESCE(
					COUNT(table.Conversation.UniqueId), CAST(Int(0)).AS_INTEGER()).AS("conversationInitiated"),
			).FROM(
				table.Conversation,
			).WHERE(
				table.Conversation.InitiatedByCampaignId.EQ(
					UUID(uuid.MustParse(context.Param("campaignId"))),
				).AND(
					table.Conversation.InitiatedBy.EQ(utils.EnumExpression(model.ConversationInitiatedEnum_Campaign.String())),
				),
			)),
		linkClicksDataCte.AS(
			SELECT(
				table.TrackLinkClick.CreatedAt.AS("date"),
				COALESCE(COUNT(table.TrackLinkClick.UniqueId), Int(0)).AS("count"),
				TO_CHAR(table.TrackLinkClick.CreatedAt, String("DD-MM-YYYY")).AS("label"),
			).FROM(
				table.TrackLinkClick.
					LEFT_JOIN(table.TrackLink, table.TrackLink.UniqueId.EQ(table.TrackLinkClick.TrackLinkId)),
			).WHERE(
				table.TrackLink.CampaignId.EQ(UUID(uuid.MustParse(context.Param("campaignId")))),
			).GROUP_BY(
				table.TrackLinkClick.CreatedAt,
			).ORDER_BY(
				table.TrackLinkClick.CreatedAt,
			),
		),
	)(SELECT(
		campaignMessageDataCte.AllColumns(),
		linkClicksCountsCte.AllColumns(),
		linkClicksDataCte.AllColumns(),
		conversationCountCte.AllColumns(),
	).FROM(
		campaignMessageDataCte,
		linkClicksCountsCte,
		linkClicksDataCte,
		conversationCountCte,
	))

	sql := campaignAnalyticsQuery.DebugSql()

	fmt.Println("SQL", sql)

	err := campaignAnalyticsQuery.QueryContext(context.Request().Context(), context.App.Db, &campaignAnalyticsData)

	if err != nil {
		fmt.Println("error is", err.Error())
		if err.Error() == qrm.ErrNoRows.Error() {
			context.App.Logger.Info("No campaign analytics found")
			responseToReturn := api_types.CampaignAnalyticsResponseSchema{
				MessagesDelivered:     0,
				MessagesFailed:        0,
				MessagesRead:          0,
				MessagesSent:          0,
				MessagesUndelivered:   0,
				TotalMessages:         0,
				ConversationInitiated: 0,
				TotalLinkClicks:       0,
				LinkClicksData:        []api_types.LinkClicksGraphDataPointSchema{},
			}
			return context.JSON(http.StatusOK, responseToReturn)
		} else {
			return context.JSON(http.StatusInternalServerError, "Error getting campaign analytics")
		}
	}

	responseToReturn := api_types.CampaignAnalyticsResponseSchema{
		MessagesDelivered:     campaignAnalyticsData.MessagesDelivered,
		MessagesFailed:        campaignAnalyticsData.MessagesFailed,
		MessagesRead:          campaignAnalyticsData.MessagesRead,
		MessagesSent:          campaignAnalyticsData.MessagesSent,
		MessagesUndelivered:   campaignAnalyticsData.MessagesUndelivered,
		TotalMessages:         campaignAnalyticsData.TotalMessages,
		ConversationInitiated: campaignAnalyticsData.ConversationInitiated,
		TotalLinkClicks:       campaignAnalyticsData.TotalLinkClicks,
		LinkClicksData:        campaignAnalyticsData.LinkClicksData,
	}

	return context.JSON(http.StatusOK, responseToReturn)
}
