package conversation_service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/wapikit/wapikit/.db-generated/model"
	"github.com/wapikit/wapikit/.db-generated/table"
	"github.com/wapikit/wapikit/api/services"
	"github.com/wapikit/wapikit/internal/api_types"
	"github.com/wapikit/wapikit/internal/core/api_server_events"
	"github.com/wapikit/wapikit/internal/core/utils"
	"github.com/wapikit/wapikit/internal/interfaces"

	"github.com/go-jet/jet/qrm"
	. "github.com/go-jet/jet/v2/postgres"
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
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetConversation,
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
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetConversation,
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
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.UpdateConversation,
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
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.DeleteConversation,
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
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.AssignConversation,
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
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.UnassignConversation,
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
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetConversation,
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

	page := queryParams.Page
	limit := queryParams.PerPage
	campaignId := queryParams.CampaignId
	status := queryParams.Status
	// listIds := queryParams.ListId
	// order := queryParams.Order

	if page == 0 || limit > 50 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid page or perPage value")
	}

	// ! fetch conversations from the database paginated
	// ! always keep the unresolved conversation with unread messages on top, sorted by the latest messages
	// ! always fetch the last 20 messages from each conversation
	// ! fetch the user assigned to the conversation

	type FetchedConversation struct {
		model.Conversation
		Contact    model.Contact   `json:"contact"`
		Tags       []model.Tag     `json:"tags"`
		Messages   []model.Message `json:"messages"`
		AssignedTo struct {
			model.OrganizationMember
			User model.User `json:"user"`
		} `json:"assignedTo"`
		NumberOfUnreadMessages int `json:"numberOfUnreadMessages"`
	}

	var fetchedConversations []FetchedConversation

	var conversationWhereQuery BoolExpression

	if *status != "" {
		conversationWhereQuery = table.Conversation.Status.EQ(utils.EnumExpression(string(*status)))
	} else {
		conversationWhereQuery = table.Conversation.Status.NOT_IN(
			utils.EnumExpression(model.ConversationStatusEnum_Deleted.String()),
			utils.EnumExpression(model.ConversationStatusEnum_Closed.String()),
		)
	}

	if *campaignId != "" {
		conversationWhereQuery = conversationWhereQuery.AND(table.Conversation.InitiatedByCampaignId.EQ(UUID(uuid.MustParse(*campaignId))))
	}

	conversationQuery := SELECT(
		table.Conversation.AllColumns,
		table.Contact.AllColumns,
		table.ConversationAssignment.AllColumns,
		table.Message.AllColumns,
		table.Tag.AllColumns,
		table.ConversationTag.AllColumns,
	).FROM(table.Conversation.
		LEFT_JOIN(table.Contact, table.Conversation.ContactId.EQ(table.Contact.UniqueId)).
		LEFT_JOIN(table.ConversationAssignment, table.Conversation.UniqueId.EQ(table.ConversationAssignment.ConversationId)).
		LEFT_JOIN(table.Message, table.Conversation.UniqueId.EQ(table.Message.ConversationId)).
		LEFT_JOIN(table.ConversationTag, table.Conversation.UniqueId.EQ(table.ConversationTag.ConversationId)).
		LEFT_JOIN(table.Tag, table.ConversationTag.TagId.EQ(table.Tag.UniqueId)),
	).
		WHERE(conversationWhereQuery).
		ORDER_BY(
			// 1. Prioritize conversations with unread messages
			CASE().
				WHEN(table.Message.Status.EQ(utils.EnumExpression(model.MessageStatus_Delivered.String()))).
				THEN(CAST(Int(1)).AS_INTEGER()).
				ELSE(CAST(Int(2)).AS_INTEGER()),
			// 2. Sort by most recent activity
			table.Conversation.UpdatedAt.DESC(),
			// 3. Active conversations on top
			CASE().
				WHEN(table.Conversation.Status.EQ(utils.EnumExpression(model.ConversationStatusEnum_Active.String()))).
				THEN(CAST(Int(1)).AS_INTEGER()).
				ELSE(CAST(Int(2)).AS_INTEGER()),
		).
		LIMIT(limit).
		OFFSET((page - 1) * limit)

	err := conversationQuery.QueryContext(context.Request().Context(), context.App.Db, &fetchedConversations)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := api_types.GetConversationsResponseSchema{
		Conversations: make([]api_types.ConversationSchema, 0),
		PaginationMeta: api_types.PaginationMeta{
			Page:    page,
			PerPage: limit,
		},
	}

	for _, conversation := range fetchedConversations {

		attr := map[string]interface{}{}
		json.Unmarshal([]byte(*conversation.Contact.Attributes), &attr)

		campaignId := ""

		if conversation.InitiatedByCampaignId != nil {
			campaignId = string(conversation.InitiatedByCampaignId.String())
		}

		conversationToAppend := api_types.ConversationSchema{
			UniqueId:               conversation.UniqueId.String(),
			ContactId:              conversation.ContactId.String(),
			OrganizationId:         conversation.OrganizationId.String(),
			InitiatedBy:            api_types.ConversationInitiatedByEnum(conversation.InitiatedBy.String()),
			CampaignId:             &campaignId,
			CreatedAt:              conversation.CreatedAt,
			Status:                 api_types.ConversationStatusEnum(conversation.Status.String()),
			Messages:               []api_types.MessageSchema{},
			NumberOfUnreadMessages: &conversation.NumberOfUnreadMessages,
			Contact: &api_types.ContactSchema{
				UniqueId:   conversation.Contact.UniqueId.String(),
				Name:       conversation.Contact.Name,
				Phone:      conversation.Contact.PhoneNumber,
				Attributes: attr,
				CreatedAt:  conversation.Contact.CreatedAt,
			},
			Tags: []api_types.TagSchema{},
		}

		if conversation.AssignedTo.UniqueId != uuid.Nil {
			member := conversation.AssignedTo
			accessLevel := api_types.UserPermissionLevel(member.AccessLevel)
			assignedToOrgMember := api_types.OrganizationMemberSchema{
				CreatedAt:   conversation.AssignedTo.CreatedAt,
				AccessLevel: accessLevel,
				UniqueId:    member.UniqueId.String(),
				Email:       member.User.Email,
				Name:        member.User.Name,
				Roles:       []api_types.OrganizationRoleSchema{},
			}

			conversationToAppend.AssignedTo = &assignedToOrgMember

		}

		for _, tag := range conversation.Tags {
			tagToAppend := api_types.TagSchema{
				UniqueId: tag.UniqueId.String(),
				Name:     tag.Label,
			}

			conversationToAppend.Tags = append(conversationToAppend.Tags, tagToAppend)
		}

		for _, message := range conversation.Messages {
			message := api_types.MessageSchema{
				UniqueId: message.UniqueId.String(),
			}
			conversationToAppend.Messages = append(conversationToAppend.Messages, message)
		}

		response.Conversations = append(response.Conversations, conversationToAppend)
	}

	return context.JSON(http.StatusOK, response)

	// WITH latest_messages AS (
	// 	SELECT
	// 		"Message"."ConversationId",
	// 		ARRAY_AGG("Message" ORDER BY "Message"."CreatedAt" DESC LIMIT 20) AS "Messages",
	// 		MAX("Message"."CreatedAt") AS "LastMessageAt",
	// 		COUNT(*) FILTER (WHERE "Message"."Status" = 'UNREAD') AS "UnreadMessageCount"
	// 	FROM
	// 		public."Message"
	// 	GROUP BY
	// 		"Message"."ConversationId"
	// ),
	// conversation_data AS (
	// 	SELECT
	// 		"Conversation"."UniqueId" AS "ConversationUniqueId",
	// 		"Conversation"."ContactId",
	// 		"Conversation"."OrganizationId",
	// 		"Conversation"."InitiatedBy",
	// 		"Conversation"."InitiatedByCampaignId",
	// 		"Conversation"."CreatedAt",
	// 		"Conversation"."UpdatedAt",
	// 		"Conversation"."Status",
	// 		"Contact"."UniqueId" AS "ContactUniqueId",
	// 		"Contact"."Name" AS "ContactName",
	// 		"Contact"."PhoneNumber" AS "ContactPhone",
	// 		"Contact"."Attributes" AS "ContactAttributes",
	// 		"Contact"."CreatedAt" AS "ContactCreatedAt",
	// 		"latest_messages"."Messages" AS "Messages",
	// 		"latest_messages"."UnreadMessageCount",
	// 		"latest_messages"."LastMessageAt"
	// 	FROM
	// 		public."Conversation"
	// 	LEFT JOIN
	// 		public."Contact" ON "Conversation"."ContactId" = "Contact"."UniqueId"
	// 	LEFT JOIN
	// 		latest_messages ON "Conversation"."UniqueId" = latest_messages."ConversationId"
	// ),
	// tag_data AS (
	// 	SELECT
	// 		"ConversationTag"."ConversationId",
	// 		JSON_AGG("Tag") AS "Tags"
	// 	FROM
	// 		public."ConversationTag"
	// 	LEFT JOIN
	// 		public."Tag" ON "ConversationTag"."TagId" = "Tag"."UniqueId"
	// 	GROUP BY
	// 		"ConversationTag"."ConversationId"
	// ),
	// assigned_users AS (
	// 	SELECT
	// 		"Conversation"."UniqueId" AS "ConversationUniqueId",
	// 		"User"."UniqueId" AS "UserUniqueId",
	// 		"User"."Name" AS "UserName"
	// 	FROM
	// 		public."Conversation"
	// 	LEFT JOIN
	// 		public."User" ON "Conversation"."AssignedTo" = "User"."UniqueId"
	// )
	// SELECT
	// 	conversation_data."ConversationUniqueId" AS "uniqueId",
	// 	conversation_data."ContactId" AS "contactId",
	// 	conversation_data."OrganizationId" AS "organizationId",
	// 	conversation_data."InitiatedBy" AS "initiatedBy",
	// 	conversation_data."InitiatedByCampaignId" AS "campaignId",
	// 	conversation_data."CreatedAt" AS "createdAt",
	// 	conversation_data."Status" AS "status",
	// 	conversation_data."Messages" AS "messages",
	// 	conversation_data."UnreadMessageCount" AS "numberOfUnreadMessages",
	// 	JSON_BUILD_OBJECT(
	// 		'uniqueId', conversation_data."ContactUniqueId",
	// 		'name', conversation_data."ContactName",
	// 		'phone', conversation_data."ContactPhone",
	// 		'attributes', conversation_data."ContactAttributes",
	// 		'createdAt', conversation_data."ContactCreatedAt"
	// 	) AS "contact",
	// 	tag_data."Tags" AS "tags",
	// 	JSON_BUILD_OBJECT(
	// 		'uniqueId', assigned_users."UserUniqueId",
	// 		'name', assigned_users."UserName"
	// 	) AS "assignedTo"
	// FROM
	// 	conversation_data
	// LEFT JOIN
	// 	tag_data ON conversation_data."ConversationUniqueId" = tag_data."ConversationId"
	// LEFT JOIN
	// 	assigned_users ON conversation_data."ConversationUniqueId" = assigned_users."ConversationUniqueId"
	// WHERE
	// 	conversation_data."Status" != 'CLOSED'
	// ORDER BY
	// 	conversation_data."UnreadMessageCount" DESC,
	// 	conversation_data."LastMessageAt" DESC
	// LIMIT 20;

	// conversationQuery := WITH(
	// 	conversationCte.AS(
	// 		SELECT(
	// 			table.Conversation.AllColumns,
	// 		).
	// 			FROM(table.Conversation).
	// 			WHERE(
	// 				table.Conversation.Status.NOT_IN(utils.EnumExpression(
	// 					model.ConversationStatusEnum_Deleted.String()),
	// 					utils.EnumExpression(model.ConversationStatusEnum_Closed.String()
	// 				),
	// 			).
	// 			ORDER_BY(
	// 				table.Conversation.LastMessageAt.DESC(),
	// 			).
	// 			LIMIT(queryParams.Limit).
	// 			OFFSET(queryParams.Offset),
	// 	)),
	// 	contactCte.AS(
	// 		SELECT(
	// 		table.Contact.AllColumns,
	// 		table.ContactListContact.AllColumns,
	// 		table.ContactList.AllColumns,
	// 		).
	// 		FROM(table.Contact.
	// 			LEFT_JOIN(table.ContactListContact, table.ContactListContact.ContactId.EQ(table.Contact.UniqueId)).
	// 			LEFT_JOIN(table.ContactList, table.Contact.UniqueId.EQ(table.ContactListContact.ContactId)),
	// 		).
	// 		WHERE(table.Contact.OrganizationId.EQ(UUID(orgUuid)).AND(Conversa)).
	// 		LIMIT(1)
	// 	),
	// 	messagesCte.AS(
	// 		SELECT(
	// 			table.Message.AllColumns,
	// 		).
	// 			FROM(table.Message).
	// 			WHERE(
	// 				table.Message.ConversationId.IN(
	// 					SELECT(conversationCte.Field("UniqueId")).FROM(conversationCte),
	// 				),
	// 			).
	// 			ORDER_BY(
	// 				table.Message.CreatedAt.DESC(),
	// 			).
	// 			LIMIT(20),
	// 	),
	// 	assignmentCte.AS(
	// 		SELECT(
	// 			table.Assignment.AllColumns,
	// 		).
	// 			FROM(table.Assignment).
	// 			WHERE(
	// 				table.Assignment.ConversationId.IN(
	// 					SELECT(conversationCte.Field("UniqueId")).FROM(conversationCte),
	// 				),
	// 			),
	// 	),
	// )(
	// 	SELECT(
	// 		conversationCte.AllColumns(),
	// 		messagesCte.AllColumns(),
	// 		assignmentCte.AllColumns(),
	// 	).
	// 		FROM(
	// 			conversationCte.
	// 				LEFT_JOIN(messagesCte, messagesCte.Field("ConversationId").EQ(conversationCte.Field("UniqueId"))).
	// 				LEFT_JOIN(assignmentCte, assignmentCte.Field("ConversationId").EQ(conversationCte.Field("UniqueId"))),
	// 		).
	// 		ORDER_BY(
	// 			conversationCte.Field("LastMessageAt").DESC(),
	// 		),
	// )

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

	return context.JSON(http.StatusBadRequest, nil)
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
	conversationId := context.Param("id")
	if conversationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "conversation id is required")
	}
	conversationUuid, err := uuid.Parse(conversationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid conversation id")
	}

	payload := new(api_types.AssignConversationSchema)
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	orgMemberUuid, err := uuid.Parse(payload.OrganizationMemberId)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid organization member id")
	}

	var conversation struct {
		model.Conversation
		Assignment model.ConversationAssignment
	}

	var organizationMember model.OrganizationMember

	conversationFetchQuery := SELECT(
		table.Conversation.AllColumns,
		table.ConversationAssignment.AllColumns,
	).FROM(
		table.Conversation.
			LEFT_JOIN(table.ConversationAssignment, table.Conversation.UniqueId.EQ(table.ConversationAssignment.ConversationId).AND(
				table.ConversationAssignment.Status.EQ(utils.EnumExpression(model.ConversationAssignmentStatus_Assigned.String())),
			)),
	).
		WHERE(
			table.Conversation.UniqueId.EQ(UUID(conversationUuid)),
		).LIMIT(1)

	organizationMemberQuery := SELECT(
		table.OrganizationMember.AllColumns,
	).FROM(
		table.OrganizationMember,
	).WHERE(
		table.OrganizationMember.UniqueId.EQ(UUID(orgMemberUuid)),
	).LIMIT(1)

	err = organizationMemberQuery.QueryContext(context.Request().Context(), context.App.Db, &organizationMember)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return echo.NewHTTPError(http.StatusNotFound, "organization member not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = conversationFetchQuery.QueryContext(context.Request().Context(), context.App.Db, &conversation)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return echo.NewHTTPError(http.StatusNotFound, "conversation not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	assignmentToInsert := model.ConversationAssignment{
		ConversationId:                 conversationUuid,
		Status:                         model.ConversationAssignmentStatus_Assigned,
		CreatedAt:                      time.Now(),
		UpdatedAt:                      time.Now(),
		AssignedToOrganizationMemberId: orgMemberUuid,
	}

	// ! update the conversation assignment record and create a new record for new assignment.

	if conversation.Assignment.ConversationId == uuid.Nil {
		unassignFromPreviousMemberCte := CTE("unassign_from_previous_member_cte")
		assignmentUpdateQuery := WITH(
			unassignFromPreviousMemberCte.AS(
				table.ConversationAssignment.UPDATE(table.ConversationAssignment.Status).
					SET(
						table.ConversationAssignment.Status.SET(utils.EnumExpression(model.ConversationAssignmentStatus_Unassigned.String())),
					).
					WHERE(
						table.ConversationAssignment.ConversationId.EQ(UUID(conversationUuid)),
					).
					RETURNING(table.ConversationAssignment.AllColumns),
			),
		)(
			table.ConversationAssignment.
				INSERT().
				MODEL(assignmentToInsert).
				RETURNING(table.ConversationAssignment.AllColumns),
		)

		err = assignmentUpdateQuery.QueryContext(context.Request().Context(), context.App.Db, &assignmentToInsert)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	} else {
		insertQuery := table.ConversationAssignment.
			INSERT().
			MODEL(assignmentToInsert).
			RETURNING(table.ConversationAssignment.AllColumns)

		err = insertQuery.QueryContext(context.Request().Context(), context.App.Db, &assignmentToInsert)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	// ! send assignment notification to the user
	event := api_server_events.BaseApiServerEvent{
		EventType: api_server_events.ApiServerChatAssignmentEvent,
	}
	context.App.Redis.PublishMessageToRedisChannel(context.App.Constants.RedisEventChannelName, string(event.ToJson()))

	responseToReturn := api_types.AssignConversationResponseSchema{
		Data: true,
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func handleUnassignConversation(context interfaces.ContextWithSession) error {
	conversationId := context.Param("id")
	if conversationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "conversation id is required")
	}
	conversationUuid, err := uuid.Parse(conversationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid conversation id")
	}

	var conversation struct {
		model.Conversation
		Assignment model.ConversationAssignment
	}

	conversationFetchQuery := SELECT(
		table.Conversation.AllColumns,
		table.ConversationAssignment.AllColumns,
	).FROM(
		table.Conversation.
			LEFT_JOIN(table.ConversationAssignment, table.Conversation.UniqueId.EQ(table.ConversationAssignment.ConversationId).AND(
				table.ConversationAssignment.Status.EQ(utils.EnumExpression(model.ConversationAssignmentStatus_Assigned.String())),
			)),
	).
		WHERE(
			table.Conversation.UniqueId.EQ(UUID(conversationUuid)),
		).LIMIT(1)

	err = conversationFetchQuery.QueryContext(context.Request().Context(), context.App.Db, &conversation)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return echo.NewHTTPError(http.StatusNotFound, "conversation not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	updateAssignmentQuery := table.ConversationAssignment.
		UPDATE(table.ConversationAssignment.Status).
		SET(
			table.ConversationAssignment.Status.SET(utils.EnumExpression(model.ConversationAssignmentStatus_Unassigned.String())),
		).
		WHERE(
			table.ConversationAssignment.ConversationId.EQ(UUID(conversationUuid)),
		).
		RETURNING(table.ConversationAssignment.AllColumns)

	err = updateAssignmentQuery.QueryContext(context.Request().Context(), context.App.Db, &conversation.Assignment)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// ! send un-assignment notification to the user
	redis := context.App.Redis

	event := api_server_events.BaseApiServerEvent{
		EventType: api_server_events.ApiServerChatUnAssignmentEvent,
	}

	redis.PublishMessageToRedisChannel(context.App.Constants.RedisEventChannelName, string(event.ToJson()))

	responseToReturn := api_types.UnassignConversationResponseSchema{
		Data: true,
	}

	return context.JSON(http.StatusOK, responseToReturn)
}
