package conversation_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wapikit/wapi.go/pkg/components"
	"github.com/wapikit/wapikit/.db-generated/model"
	"github.com/wapikit/wapikit/.db-generated/table"
	"github.com/wapikit/wapikit/api/api_types"
	controller "github.com/wapikit/wapikit/api/controllers"
	"github.com/wapikit/wapikit/interfaces"
	"github.com/wapikit/wapikit/services/event_service"
	"github.com/wapikit/wapikit/utils"

	"github.com/go-jet/jet/qrm"
	. "github.com/go-jet/jet/v2/postgres"
)

type ConversationController struct {
	controller.BaseController `json:"-,inline"`
}

func NewConversationController() *ConversationController {
	return &ConversationController{
		BaseController: controller.BaseController{
			Name:        "Conversation Controller",
			RestApiPath: "/api/conversation",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/conversations",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetConversations),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    600,
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
							MaxRequests:    600,
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
							MaxRequests:    600,
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
							MaxRequests:    600,
							WindowTimeInMs: time.Hour.Milliseconds(),
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.DeleteConversation,
						},
					},
				},
				{
					Path:                    "/api/conversation/:id/assign",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(handleAssignConversation),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    600,
							WindowTimeInMs: time.Hour.Milliseconds(),
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.AssignConversation,
						},
					},
				},
				{
					Path:                    "/api/conversation/:id/unassign",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(handleUnassignConversation),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    600,
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
							MaxRequests:    600,
							WindowTimeInMs: time.Hour.Milliseconds(),
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetConversation,
						},
					},
				},
				{
					Path:                    "/api/conversation/:id/messages",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(handleSendMessage),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    600,
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

	orgId := context.Session.User.OrganizationId
	orgUuid := uuid.MustParse(orgId)

	fmt.Println("Query Params are:", context.QueryParams())
	queryParams := new(api_types.GetConversationsParams)
	if err := utils.BindQueryParams(context, queryParams); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	page := queryParams.Page
	limit := queryParams.PerPage
	campaignId := queryParams.CampaignId
	status := queryParams.Status
	// listIds := queryParams.ListId
	// order := queryParams.Order

	if page == 0 || limit > 50 {
		return context.JSON(http.StatusBadRequest, "Invalid page or perPage value")
	}

	// ! fetch conversations from the database paginated
	// ! always keep the unresolved conversation with unread messages on top, sorted by the latest messages
	// ! always fetch the last 20 messages from each conversation
	// ! fetch the user assigned to the conversation

	type FetchedConversation struct {
		model.Conversation
		Contact struct {
			model.Contact
			ContactLists []struct {
				model.ContactList
			} `json:"contactLists"`
		} `json:"contact"`
		Tags       []model.Tag     `json:"tags"`
		Messages   []model.Message `json:"messages"`
		AssignedTo struct {
			model.OrganizationMember
			User model.User `json:"user"`
		} `json:"assignedTo"`
		NumberOfUnreadMessages int `json:"numberOfUnreadMessages"`
	}

	var fetchedConversations []FetchedConversation

	conversationWhereQuery := table.Conversation.OrganizationId.EQ(UUID(orgUuid))

	if status != nil {
		conversationWhereQuery = conversationWhereQuery.AND(
			table.Conversation.Status.EQ(utils.EnumExpression(string(*status))),
		)
	} else {
		conversationWhereQuery = conversationWhereQuery.AND(
			table.Conversation.Status.NOT_IN(
				utils.EnumExpression(model.ConversationStatusEnum_Deleted.String()),
				utils.EnumExpression(model.ConversationStatusEnum_Closed.String()),
			),
		)
	}

	if campaignId != nil {
		conversationWhereQuery = conversationWhereQuery.AND(table.Conversation.InitiatedByCampaignId.EQ(UUID(uuid.MustParse(*campaignId))))
	}

	conversationQuery := SELECT(
		table.Conversation.AllColumns,
		table.Contact.AllColumns,
		table.ContactListContact.AllColumns,
		table.ContactList.AllColumns,
		table.ConversationAssignment.AllColumns,
		table.OrganizationMember.AllColumns,
		table.User.AllColumns,
		table.Message.AllColumns,
		table.Tag.AllColumns,
		table.ConversationTag.AllColumns,
	).FROM(table.Conversation.
		LEFT_JOIN(table.Contact, table.Conversation.ContactId.EQ(table.Contact.UniqueId)).
		LEFT_JOIN(table.ContactListContact, table.Contact.UniqueId.EQ(table.ContactListContact.ContactId)).
		LEFT_JOIN(table.ContactList, table.Contact.UniqueId.EQ(table.ContactListContact.ContactId)).
		LEFT_JOIN(table.ConversationAssignment, table.Conversation.UniqueId.EQ(table.ConversationAssignment.ConversationId)).
		LEFT_JOIN(table.OrganizationMember, table.ConversationAssignment.AssignedToOrganizationMemberId.EQ(table.OrganizationMember.UniqueId)).
		LEFT_JOIN(table.User, table.OrganizationMember.UserId.EQ(table.User.UniqueId)).
		LEFT_JOIN(table.Message, table.Conversation.UniqueId.EQ(table.Message.ConversationId)).
		LEFT_JOIN(table.ConversationTag, table.Conversation.UniqueId.EQ(table.ConversationTag.ConversationId)).
		LEFT_JOIN(table.Tag, table.ConversationTag.TagId.EQ(table.Tag.UniqueId)),
	).
		WHERE(conversationWhereQuery).
		ORDER_BY(
			Raw(` MAX("Message"."CreatedAt") OVER (PARTITION BY "Conversation"."UniqueId") DESC,
			     "Message"."CreatedAt" ASC`,
			),
		).
		LIMIT(limit).
		OFFSET((page - 1) * limit)

	err := conversationQuery.QueryContext(context.Request().Context(), context.App.Db, &fetchedConversations)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
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

		lists := []api_types.ContactListSchema{}

		for _, contactList := range conversation.Contact.ContactLists {
			stringUniqueId := contactList.UniqueId.String()
			listToAppend := api_types.ContactListSchema{
				UniqueId: stringUniqueId,
				Name:     contactList.Name,
			}
			lists = append(lists, listToAppend)
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
			NumberOfUnreadMessages: conversation.NumberOfUnreadMessages,
			Contact: api_types.ContactWithoutConversationSchema{
				UniqueId:   conversation.Contact.UniqueId.String(),
				Name:       conversation.Contact.Name,
				Phone:      conversation.Contact.PhoneNumber,
				Attributes: attr,
				CreatedAt:  conversation.Contact.CreatedAt,
				Status:     api_types.ContactStatusEnum(conversation.Contact.Status.String()),
			},
			Tags: []api_types.TagSchema{},
		}

		context.App.Logger.Info("conversation: %v", conversation.AssignedTo)

		if conversation.AssignedTo.UniqueId != uuid.Nil {
			member := conversation.AssignedTo
			accessLevel := api_types.UserPermissionLevelEnum(member.AccessLevel)
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
				Label:    tag.Label,
			}
			conversationToAppend.Tags = append(conversationToAppend.Tags, tagToAppend)
		}

		for _, message := range conversation.Messages {
			messageData := map[string]interface{}{}
			json.Unmarshal([]byte(*message.MessageData), &messageData)
			message := api_types.MessageSchema{
				UniqueId:       message.UniqueId.String(),
				ConversationId: message.ConversationId.String(),
				CreatedAt:      message.CreatedAt,
				Direction:      api_types.MessageDirectionEnum(message.Direction.String()),
				MessageData:    &messageData,
				MessageType:    api_types.MessageTypeEnum(message.MessageType.String()),
				Status:         api_types.MessageStatusEnum(message.Status.String()),
			}
			conversationToAppend.Messages = append(conversationToAppend.Messages, message)
		}

		response.Conversations = append(response.Conversations, conversationToAppend)

	}

	return context.JSON(http.StatusOK, response)
}

func handleGetConversationById(context interfaces.ContextWithSession) error {
	conversationId := context.Param("id")

	if conversationId == "" {
		return context.JSON(http.StatusBadRequest, "conversation id is required")
	}
	conversationUuid, err := uuid.Parse(conversationId)

	if err != nil {
		return context.JSON(http.StatusBadRequest, "invalid conversation id")
	}

	type FetchedConversation struct {
		model.Conversation
		Contact struct {
			model.Contact
			ContactLists []struct {
				model.ContactList
			} `json:"contactLists"`
		} `json:"contact"`
		Tags       []model.Tag     `json:"tags"`
		Messages   []model.Message `json:"messages"`
		AssignedTo struct {
			model.OrganizationMember
			User model.User `json:"user"`
		} `json:"assignedTo"`
		NumberOfUnreadMessages int `json:"numberOfUnreadMessages"`
	}

	var conversation FetchedConversation

	conversationQuery := SELECT(
		table.Conversation.AllColumns,
		table.Contact.AllColumns,
		table.ContactListContact.AllColumns,
		table.ContactList.AllColumns,
		table.ConversationAssignment.AllColumns,
		table.Message.AllColumns,
		table.Tag.AllColumns,
		table.ConversationTag.AllColumns,
	).FROM(table.Conversation.
		LEFT_JOIN(table.Contact, table.Conversation.ContactId.EQ(table.Contact.UniqueId)).
		LEFT_JOIN(table.ContactListContact, table.Contact.UniqueId.EQ(table.ContactListContact.ContactId)).
		LEFT_JOIN(table.ContactList, table.Contact.UniqueId.EQ(table.ContactListContact.ContactId)).
		LEFT_JOIN(table.ConversationAssignment, table.Conversation.UniqueId.EQ(table.ConversationAssignment.ConversationId)).
		LEFT_JOIN(table.Message, table.Conversation.UniqueId.EQ(table.Message.ConversationId)).
		LEFT_JOIN(table.ConversationTag, table.Conversation.UniqueId.EQ(table.ConversationTag.ConversationId)).
		LEFT_JOIN(table.Tag, table.ConversationTag.TagId.EQ(table.Tag.UniqueId)),
	).
		WHERE(
			table.Conversation.UniqueId.EQ(UUID(conversationUuid)),
		).
		ORDER_BY(
			Raw(` MAX("Message"."CreatedAt") OVER (PARTITION BY "Conversation"."UniqueId") DESC,
			     "Message"."CreatedAt" ASC`,
			),
		)

	err = conversationQuery.QueryContext(context.Request().Context(), context.App.Db, &conversation)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "conversation not found")
		}
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	response := api_types.GetConversationByIdResponseSchema{
		Conversation: api_types.ConversationSchema{},
	}

	attr := map[string]interface{}{}
	json.Unmarshal([]byte(*conversation.Contact.Attributes), &attr)
	campaignId := ""

	if conversation.InitiatedByCampaignId != nil {
		campaignId = string(conversation.InitiatedByCampaignId.String())
	}

	lists := []api_types.ContactListSchema{}

	for _, contactList := range conversation.Contact.ContactLists {
		stringUniqueId := contactList.UniqueId.String()
		listToAppend := api_types.ContactListSchema{
			UniqueId: stringUniqueId,
			Name:     contactList.Name,
		}
		lists = append(lists, listToAppend)
	}

	response.Conversation = api_types.ConversationSchema{
		UniqueId:               conversation.UniqueId.String(),
		ContactId:              conversation.ContactId.String(),
		OrganizationId:         conversation.OrganizationId.String(),
		InitiatedBy:            api_types.ConversationInitiatedByEnum(conversation.InitiatedBy.String()),
		CampaignId:             &campaignId,
		CreatedAt:              conversation.CreatedAt,
		Status:                 api_types.ConversationStatusEnum(conversation.Status.String()),
		Messages:               []api_types.MessageSchema{},
		NumberOfUnreadMessages: conversation.NumberOfUnreadMessages,
		Contact: api_types.ContactWithoutConversationSchema{
			UniqueId:   conversation.Contact.UniqueId.String(),
			Name:       conversation.Contact.Name,
			Phone:      conversation.Contact.PhoneNumber,
			Attributes: attr,
			CreatedAt:  conversation.Contact.CreatedAt,
			Status:     api_types.ContactStatusEnum(conversation.Contact.Status.String()),
		},
		Tags: []api_types.TagSchema{},
	}

	if conversation.AssignedTo.UniqueId != uuid.Nil {
		member := conversation.AssignedTo
		accessLevel := api_types.UserPermissionLevelEnum(member.AccessLevel)
		assignedToOrgMember := api_types.OrganizationMemberSchema{
			CreatedAt:   conversation.AssignedTo.CreatedAt,
			AccessLevel: accessLevel,
			UniqueId:    member.UniqueId.String(),
			Email:       member.User.Email,
			Name:        member.User.Name,
			Roles:       []api_types.OrganizationRoleSchema{},
		}

		response.Conversation.AssignedTo = &assignedToOrgMember
	}

	for _, tag := range conversation.Tags {
		tagToAppend := api_types.TagSchema{
			UniqueId: tag.UniqueId.String(),
			Label:    tag.Label,
		}
		response.Conversation.Tags = append(response.Conversation.Tags, tagToAppend)
	}

	for _, message := range conversation.Messages {
		messageData := map[string]interface{}{}
		json.Unmarshal([]byte(*message.MessageData), &messageData)
		message := api_types.MessageSchema{
			UniqueId:       message.UniqueId.String(),
			ConversationId: message.ConversationId.String(),
			CreatedAt:      message.CreatedAt,
			Direction:      api_types.MessageDirectionEnum(message.Direction.String()),
			MessageData:    &messageData,
			MessageType:    api_types.MessageTypeEnum(message.MessageType.String()),
			Status:         api_types.MessageStatusEnum(message.Status.String()),
		}
		response.Conversation.Messages = append(response.Conversation.Messages, message)
	}

	return context.JSON(http.StatusOK, response)
}

func handleUpdateConversationById(context interfaces.ContextWithSession) error {
	conversationId := context.Param("id")
	if conversationId == "" {
		return context.JSON(http.StatusBadRequest, "conversation id is required")
	}
	conversationUuid, err := uuid.Parse(conversationId)
	if err != nil {
		return context.JSON(http.StatusBadRequest, "invalid conversation id")
	}

	type FetchedConversation struct {
		model.Conversation
		Contact struct {
			model.Contact
			ContactLists []struct {
				model.ContactList
			} `json:"contactLists"`
		} `json:"contact"`
		Tags       []model.Tag     `json:"tags"`
		Messages   []model.Message `json:"messages"`
		AssignedTo struct {
			model.OrganizationMember
			User model.User `json:"user"`
		} `json:"assignedTo"`
		NumberOfUnreadMessages int `json:"numberOfUnreadMessages"`
	}

	context.App.Logger.Info("conversation id: %v", conversationUuid)

	return nil
}

func handleDeleteConversationById(context interfaces.ContextWithSession) error {
	conversationId := context.Param("id")
	if conversationId == "" {
		return context.JSON(http.StatusBadRequest, "conversation id is required")
	}
	conversationUuid, err := uuid.Parse(conversationId)
	if err != nil {
		return context.JSON(http.StatusBadRequest, "invalid conversation id")
	}

	context.App.Logger.Info("conversation id: %v", conversationUuid)

	return context.JSON(http.StatusBadRequest, nil)
}

func handleGetConversationMessages(context interfaces.ContextWithSession) error {
	conversationId := context.Param("id")
	if conversationId == "" {
		return context.JSON(http.StatusBadRequest, "conversation id is required")
	}
	conversationUuid, err := uuid.Parse(conversationId)
	if err != nil {
		return context.JSON(http.StatusBadRequest, "invalid conversation id")
	}

	queryParams := new(api_types.GetConversationMessagesParams)
	if err := utils.BindQueryParams(context, queryParams); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	page := queryParams.Page
	limit := queryParams.PerPage

	if page == 0 || limit > 50 {
		return context.JSON(http.StatusBadRequest, "Invalid page or perPage value")
	}

	type FetchedMessage struct {
		model.Message
	}

	var dest []struct {
		TotalMessages int `json:"totalMessages"`
		model.Message
	}

	messageQuery := SELECT(
		table.Message.AllColumns,
		COUNT(table.Contact.UniqueId).OVER().AS("totalMessages"),
	).FROM(table.Message).
		WHERE(
			table.Message.ConversationId.EQ(UUID(conversationUuid)),
		).
		ORDER_BY(
			table.Message.CreatedAt.ASC(),
		).
		LIMIT(limit).
		OFFSET((page - 1) * limit)

	err = messageQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			total := 0
			messages := make([]api_types.MessageSchema, 0)
			return context.JSON(http.StatusOK, api_types.GetConversationMessagesResponseSchema{
				Messages: messages,
				PaginationMeta: api_types.PaginationMeta{
					Page:    page,
					PerPage: limit,
					Total:   total,
				},
			})
		}

		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	messagesToReturn := []api_types.MessageSchema{}
	totalMessages := 0

	if len(dest) > 0 {
		for _, message := range dest {
			messageData := map[string]interface{}{}
			json.Unmarshal([]byte(*message.MessageData), &messageData)
			message := api_types.MessageSchema{
				UniqueId:       message.UniqueId.String(),
				ConversationId: message.ConversationId.String(),
				CreatedAt:      message.CreatedAt,
				Direction:      api_types.MessageDirectionEnum(message.Direction.String()),
				MessageData:    &messageData,
				MessageType:    api_types.MessageTypeEnum(message.MessageType.String()),
				Status:         api_types.MessageStatusEnum(message.Status.String()),
			}
			messagesToReturn = append(messagesToReturn, message)
		}

		totalMessages = dest[0].TotalMessages
	}

	response := api_types.GetConversationMessagesResponseSchema{
		Messages: messagesToReturn,
		PaginationMeta: api_types.PaginationMeta{
			Page:    page,
			PerPage: limit,
			Total:   totalMessages,
		},
	}

	return context.JSON(http.StatusOK, response)
}

func handleSendMessage(context interfaces.ContextWithSession) error {
	conversationId := context.Param("id")
	if conversationId == "" {
		return context.JSON(http.StatusBadRequest, "conversation id is required")
	}
	conversationUuid, err := uuid.Parse(conversationId)
	if err != nil {
		return context.JSON(http.StatusBadRequest, "invalid conversation id")
	}

	payload := new(api_types.NewMessageSchema)

	if err := context.Bind(payload); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	var conversationWithContact struct {
		model.Conversation
		Contact                 model.Contact                 `json:"contact"`
		WhatsappBusinessAccount model.WhatsappBusinessAccount `json:"whatsappBusinessAccount"`
	}

	conversationFetchQuery := SELECT(
		table.Conversation.AllColumns,
		table.Contact.AllColumns,
		table.WhatsappBusinessAccount.AllColumns,
	).FROM(
		table.Conversation.LEFT_JOIN(
			table.Contact, table.Conversation.ContactId.EQ(table.Contact.UniqueId),
		).LEFT_JOIN(
			table.WhatsappBusinessAccount, table.WhatsappBusinessAccount.OrganizationId.EQ(table.Conversation.OrganizationId),
		),
	).WHERE(
		table.Conversation.UniqueId.EQ(UUID(conversationUuid)),
	).LIMIT(1)

	err = conversationFetchQuery.QueryContext(context.Request().Context(), context.App.Db, &conversationWithContact)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "conversation not found")
		}
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	messageData, err := json.Marshal(payload.MessageData)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	messagingClient := context.App.WapiClient.NewMessagingClient(
		conversationWithContact.PhoneNumberUsed,
	)

	var whatsappMessageId string

	//  ! handle all the message type to send here
	switch *payload.MessageType {
	case api_types.Text:
		if payload.MessageData != nil {
			messageData, err = json.Marshal(payload.MessageData)
			if err != nil {
				return context.JSON(http.StatusInternalServerError, err.Error())
			}
		}

		textMessageData := *payload.MessageData

		{
			textMessage, err := components.NewTextMessage(components.TextMessageConfigs{
				Text: textMessageData["text"].(string),
			})

			if err != nil {
				return context.JSON(http.StatusInternalServerError, err.Error())
			}

			response, err := messagingClient.Message.Send(textMessage, conversationWithContact.Contact.PhoneNumber)
			if err != nil {
				return context.JSON(http.StatusInternalServerError, err.Error())
			}

			var jsonResponse map[string]interface{}
			err = json.Unmarshal([]byte(response), &jsonResponse)

			fmt.Println("response: %v", jsonResponse)

			if err != nil {
				return context.JSON(http.StatusInternalServerError, err.Error())
			}

			whatsappMessageId = jsonResponse["messages"].([]interface{})[0].(map[string]interface{})["id"].(string)

			context.App.Logger.Info("response: %v", response, nil)
		}
	}

	stringMessageData := string(messageData)

	messageToInsert := model.Message{
		ConversationId:            &conversationWithContact.UniqueId,
		Direction:                 model.MessageDirectionEnum_OutBound,
		WhatsAppMessageId:         &whatsappMessageId,
		WhatsappBusinessAccountId: &conversationWithContact.WhatsappBusinessAccount.AccountId,
		CampaignId:                nil,
		ContactId:                 conversationWithContact.ContactId,
		MessageType:               model.MessageTypeEnum_Text,
		Status:                    model.MessageStatusEnum_Sent,
		MessageData:               &stringMessageData,
		OrganizationId:            conversationWithContact.OrganizationId,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
		PhoneNumberUsed:           conversationWithContact.PhoneNumberUsed,
	}

	var insertedMessage model.Message

	insertQuery := table.Message.
		INSERT(table.Message.MutableColumns).
		MODEL(messageToInsert).
		RETURNING(table.Message.AllColumns)

	err = insertQuery.QueryContext(context.Request().Context(), context.App.Db, &insertedMessage)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	responseToReturn := api_types.SendMessageInConversationResponseSchema{
		Message: api_types.MessageSchema{
			UniqueId:       insertedMessage.UniqueId.String(),
			ConversationId: insertedMessage.ConversationId.String(),
			CreatedAt:      insertedMessage.CreatedAt,
			Direction:      api_types.MessageDirectionEnum(insertedMessage.Direction.String()),
			MessageData:    payload.MessageData,
			MessageType:    api_types.MessageTypeEnum(insertedMessage.MessageType.String()),
			Status:         api_types.MessageStatusEnum(insertedMessage.Status.String()),
		},
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func handleAssignConversation(context interfaces.ContextWithSession) error {
	conversationId := context.Param("id")
	if conversationId == "" {
		return context.JSON(http.StatusBadRequest, "conversation id is required")
	}
	conversationUuid, err := uuid.Parse(conversationId)
	if err != nil {
		return context.JSON(http.StatusBadRequest, "invalid conversation id")
	}

	payload := new(api_types.AssignConversationSchema)
	if err := context.Bind(payload); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	orgMemberUuid, err := uuid.Parse(payload.OrganizationMemberId)

	if err != nil {
		return context.JSON(http.StatusBadRequest, "invalid organization member id")
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
			return context.JSON(http.StatusNotFound, "organization member not found")
		}
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	err = conversationFetchQuery.QueryContext(context.Request().Context(), context.App.Db, &conversation)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "conversation not found")
		}
		return context.JSON(http.StatusInternalServerError, err.Error())
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
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		insertQuery := table.ConversationAssignment.
			INSERT().
			MODEL(assignmentToInsert).
			RETURNING(table.ConversationAssignment.AllColumns)

		err = insertQuery.QueryContext(context.Request().Context(), context.App.Db, &assignmentToInsert)

		if err != nil {
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	// ! send assignment notification to the user

	event := event_service.ChatAssignmentEvent{}

	context.App.Redis.PublishMessageToRedisChannel(context.App.Constants.RedisEventChannelName, event.ToJson())

	responseToReturn := api_types.AssignConversationResponseSchema{
		Data: true,
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func handleUnassignConversation(context interfaces.ContextWithSession) error {
	conversationId := context.Param("id")
	if conversationId == "" {
		return context.JSON(http.StatusBadRequest, "conversation id is required")
	}
	conversationUuid, err := uuid.Parse(conversationId)
	if err != nil {
		return context.JSON(http.StatusBadRequest, "invalid conversation id")
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
			return context.JSON(http.StatusNotFound, "conversation not found")
		}
		return context.JSON(http.StatusInternalServerError, err.Error())
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
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	// ! send un-assignment notification to the user
	redis := context.App.Redis

	event := event_service.ChatUnAssignmentEvent{
		BaseApiServerEvent: event_service.BaseApiServerEvent{
			EventType: event_service.ApiServerChatUnAssignmentEvent,
		},
	}

	redis.PublishMessageToRedisChannel(context.App.Constants.RedisEventChannelName, event.ToJson())

	responseToReturn := api_types.UnassignConversationResponseSchema{
		Data: true,
	}

	return context.JSON(http.StatusOK, responseToReturn)
}
