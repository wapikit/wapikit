//go:build community_edition
// +build community_edition

package ai_controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/wapikit/wapikit/.db-generated/model"
	"github.com/wapikit/wapikit/api/api_types"
	controller "github.com/wapikit/wapikit/api/controllers"
	"github.com/wapikit/wapikit/interfaces"
	"github.com/wapikit/wapikit/utils"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/wapikit/wapikit/.db-generated/table"
)

type AiController struct {
	controller.BaseController `json:"-,inline"`
}

func NewAiController() *AiController {
	return &AiController{
		BaseController: controller.BaseController{
			Name:        "AI Controller",
			RestApiPath: "/api/ai",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/ai/chats",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetChats),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
				// ! TODO: add the route to create a new chat in the future when we will support multiple chat
				{
					Path:                    "/api/ai/chat/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetChatById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
				{
					Path:                    "/api/ai/chat/:id/messages",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetChatMessages),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
				{
					Path:                    "/api/ai/chat/:id/messages",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(handleReplyToChat),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
				{
					Path:                    "/api/ai/chat/:id/vote",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getMessageVotes),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
				{
					Path:                    "/api/ai/chat/:id/vote",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(voteMessage),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
				{
					Path:                    "/api/ai/segment-recommendations",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetSegmentRecommendation),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
				{
					Path:                    "/api/ai/get-chat-summary",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetChatSummary),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
				{
					Path:                    "/api/ai/response-suggestions",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetResponseSuggestions),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
			},
		},
	}
}

func handleGetResponseSuggestions(context interfaces.ContextWithSession) error {

	params := new(api_types.GetConversationResponseSuggestionsParams)

	err := utils.BindQueryParams(context, params)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	conversationId := params.ConversationId

	if conversationId == "" {
		return context.JSON(http.StatusBadRequest, "Invalid conversation Id")
	}

	conversationUuid, err := uuid.Parse(conversationId)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Invalid conversation Id")
	}

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Invalid organization Id")
	}

	// * get the conversation from the db

	var dest struct {
		model.Conversation
		Messages []model.Message
	}

	fetchConversationQuery := SELECT(
		table.Conversation.AllColumns,
		table.Message.AllColumns,
	).FROM(
		table.Conversation.LEFT_JOIN(
			table.Message,
			table.Conversation.UniqueId.EQ(table.Message.ConversationId),
		),
	).WHERE(
		table.Conversation.UniqueId.EQ(UUID(conversationUuid)).AND(
			table.Conversation.OrganizationId.EQ(UUID(orgUuid)),
		),
	)

	err = fetchConversationQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "Conversation not found")
		} else {
			return context.JSON(http.StatusInternalServerError, "Error fetching conversation")
		}
	}

	if len(dest.Messages) == 0 {
		return context.JSON(http.StatusOK, api_types.GetResponseSuggestionsResponse{
			Suggestions: []string{},
		})
	}

	// * get the response suggestions from the AI model
	aiService := context.App.AiService
	suggestions, err := aiService.GetResponseSuggestions(
		context.Request().Context(),
		dest.Messages,
	)

	return context.JSON(http.StatusOK, api_types.GetResponseSuggestionsResponse{
		Suggestions: suggestions,
	})
}

func handleGetChats(context interfaces.ContextWithSession) error {
	params := new(api_types.GetAiChatsParams)

	err := utils.BindQueryParams(context, params)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	pageNumber := params.Page
	pageSize := params.PerPage

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Error parsing organization UUID")
	}
	userUuid, err := uuid.Parse(context.Session.User.UniqueId)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Error parsing user UUID")
	}

	orgMemberQuery := SELECT(
		table.OrganizationMember.AllColumns,
	).FROM(
		table.OrganizationMember,
	).WHERE(
		table.OrganizationMember.OrganizationId.EQ(UUID(orgUuid)).AND(
			table.OrganizationMember.UserId.EQ(UUID(userUuid)),
		),
	).LIMIT(1)

	var orgMember model.OrganizationMember

	err = orgMemberQuery.Query(context.App.Db, &orgMember)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "Organization member not found")
		} else {
			return context.JSON(http.StatusInternalServerError, "Error fetching organization member")
		}
	}

	// * get the chats from the db
	var dest []struct {
		TotalAiChats int `json:"totalAiChats"`
		model.AiChat
	}

	fetchChatsQuery := SELECT(
		table.AiChat.AllColumns,
		COUNT(table.AiChat.UniqueId).OVER().AS("totalAiChats"),
	).FROM(
		table.AiChat,
	).WHERE(
		table.AiChat.OrganizationId.EQ(UUID(orgUuid)).AND(
			table.AiChat.OrganizationMemberId.EQ(UUID(orgMember.UniqueId)),
		),
	)

	err = fetchChatsQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			total := 0
			chats := make([]api_types.AiChatSchema, 0)
			return context.JSON(http.StatusOK, api_types.GetAiChatsResponseSchema{
				Chats: chats,
				PaginationMeta: api_types.PaginationMeta{
					Page:    *pageNumber,
					PerPage: *pageSize,
					Total:   total,
				},
			})

		} else {
			return context.JSON(http.StatusInternalServerError, "Error fetching chats")
		}
	}

	chatsToReturn := []api_types.AiChatSchema{}

	for _, chat := range dest {
		chatsToReturn = append(chatsToReturn, api_types.AiChatSchema{
			UniqueId:    chat.UniqueId.String(),
			CreatedAt:   chat.CreatedAt,
			Description: chat.Description,
			Title:       chat.Title,
		})
	}

	total := 0
	if len(chatsToReturn) > 0 {
		total = dest[0].TotalAiChats
	}

	return context.JSON(http.StatusOK, api_types.GetAiChatsResponseSchema{
		Chats: chatsToReturn,
		PaginationMeta: api_types.PaginationMeta{
			Page:    *pageNumber,
			PerPage: *pageSize,
			Total:   total,
		},
	})
}

func handleGetChatById(context interfaces.ContextWithSession) error {
	chatId := context.Param("id")
	if chatId == "" {
		return context.JSON(http.StatusBadRequest, "Invalid Chat Id")
	}

	chatUuid, _ := uuid.Parse(chatId)
	var dest model.AiChat
	fetchChatQuery := SELECT(
		table.AiChat.AllColumns,
	).FROM(
		table.AiChat,
	).WHERE(
		table.AiChat.UniqueId.EQ(UUID(chatUuid)),
	).LIMIT(1)

	err := fetchChatQuery.Query(context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "Chat not found")
		} else {
			return context.JSON(http.StatusInternalServerError, "Error fetching chat")
		}
	}

	responseToReturn := api_types.GetAiChatByIdResponseSchema{
		Chat: api_types.AiChatSchema{
			UniqueId:    dest.UniqueId.String(),
			CreatedAt:   dest.CreatedAt,
			Description: dest.Description,
			Title:       dest.Title,
		},
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func handleGetChatMessages(context interfaces.ContextWithSession) error {
	chatId := context.Param("id")
	if chatId == "" {
		return context.JSON(http.StatusBadRequest, "Invalid Chat Id")
	}
	chatUuid, err := uuid.Parse(chatId)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Invalid chat Id")
	}

	params := new(api_types.GetAiChatMessageVotesParams)
	err = utils.BindQueryParams(context, params)

	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	pageNumber := params.Page
	pageSize := params.PerPage

	var dest []struct {
		model.AiChatMessage
	}

	fetchMessagesQuery := SELECT(
		table.AiChatMessage.AllColumns,
	).FROM(
		table.AiChatMessage,
	).WHERE(
		table.AiChatMessage.AiChatId.EQ(UUID(chatUuid)),
	).ORDER_BY(
		table.AiChatMessage.CreatedAt.ASC(),
	).
		LIMIT(pageSize).
		OFFSET((pageNumber - 1) * pageSize)

	err = fetchMessagesQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "Messages not found")
		} else {
			return context.JSON(http.StatusInternalServerError, "Error fetching messages")
		}
	}

	messagesToReturn := []api_types.AiChatMessageSchema{}

	for _, message := range dest {
		role := api_types.AiChatMessageRoleEnum(message.Role)
		messagesToReturn = append(messagesToReturn, api_types.AiChatMessageSchema{
			UniqueId:  message.UniqueId.String(),
			CreatedAt: message.CreatedAt,
			Content:   message.Content,
			Role:      role,
		})
	}

	return context.JSON(http.StatusOK, api_types.GetAiChatMessagesResponseSchema{
		Messages: messagesToReturn,
	})
}

func voteMessage(context interfaces.ContextWithSession) error {

	chatId := context.Param("id")
	if chatId == "" {
		return context.JSON(http.StatusBadRequest, "Invalid Chat Id")
	}

	chatUuid, err := uuid.Parse(chatId)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Invalid chat Id")
	}

	payload := new(api_types.AiChatMessageVoteCreateSchema)
	if err := context.Bind(payload); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	// * check if the chat and message exists

	var dest struct {
		model.AiChat
		AiChatMessage model.AiChatMessage `alias:"ai_chat_message"`
	}

	messageUuid := uuid.MustParse(payload.MessageId)

	fetchChatQuery := SELECT(
		table.AiChat.AllColumns,
		table.AiChatMessage.AllColumns,
	).FROM(
		table.AiChat.LEFT_JOIN(
			table.AiChatMessage,
			table.AiChat.UniqueId.EQ(table.AiChatMessage.AiChatId),
		),
	).WHERE(
		table.AiChat.UniqueId.EQ(UUID(chatUuid)).AND(
			table.AiChatMessage.UniqueId.EQ(UUID(messageUuid)),
		),
	).LIMIT(1)

	err = fetchChatQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "Chat or message not found")
		} else {
			return context.JSON(http.StatusInternalServerError, "Error fetching chat or message")
		}
	}

	return nil
}

func getMessageVotes(context interfaces.ContextWithSession) error {
	params := new(api_types.GetAiChatMessageVotesParams)

	err := utils.BindQueryParams(context, params)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	pageNumber := params.Page
	pageSize := params.PerPage

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Error parsing organization UUID")
	}
	userUuid, err := uuid.Parse(context.Session.User.UniqueId)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Error parsing user UUID")
	}

	votesQuery := SELECT(
		table.AiChatMessageVote.AllColumns,
		table.AiChatMessage.AllColumns,
		table.OrganizationMember.AllColumns,
		COUNT(table.AiChatMessageVote.UniqueId).OVER().AS("totalAiChats"),
	).FROM(
		table.AiChatMessageVote.LEFT_JOIN(
			table.AiChatMessage,
			table.AiChatMessageVote.AiChatMessageId.EQ(table.AiChatMessage.UniqueId),
		).LEFT_JOIN(
			table.OrganizationMember,
			table.AiChatMessage.OrganizationMemberId.EQ(table.OrganizationMember.UniqueId),
		),
	).WHERE(
		table.AiChatMessage.OrganizationId.EQ(UUID(orgUuid)).AND(
			table.AiChatMessage.OrganizationMemberId.EQ(UUID(userUuid)),
		),
	)

	var dest []struct {
		TotalVotes int `json:"totalVotes"`
		model.AiChatMessageVote
	}

	err = votesQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			total := 0
			votes := make([]api_types.AiChatMessageVoteSchema, 0)
			return context.JSON(http.StatusOK, api_types.GetAiChatVotesResponseSchema{
				Votes: votes,
				PaginationMeta: api_types.PaginationMeta{
					Page:    pageNumber,
					PerPage: pageSize,
					Total:   total,
				},
			})
		}

		return context.JSON(http.StatusInternalServerError, "Error fetching votes")
	}

	votesToReturn := []api_types.AiChatMessageVoteSchema{}

	for _, vote := range dest {
		votesToReturn = append(votesToReturn, api_types.AiChatMessageVoteSchema{
			UniqueId:  vote.UniqueId.String(),
			CreatedAt: vote.CreatedAt,
			MessageId: vote.AiChatMessageId.String(),
			Vote:      api_types.AiChatMessageVoteEnum(vote.Vote),
		})
	}

	total := 0
	if len(votesToReturn) > 0 {
		total = dest[0].TotalVotes
	}

	responseToReturn := api_types.GetAiChatVotesResponseSchema{
		Votes: votesToReturn,
		PaginationMeta: api_types.PaginationMeta{
			Page:    pageNumber,
			PerPage: pageSize,
			Total:   total,
		},
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

// ! this handle streams response to the client
func handleReplyToChat(context interfaces.ContextWithSession) error {
	isCloudEdition := context.App.Constants.IsCloudEdition
	if isCloudEdition {
		isLimitReached := context.IsAiLimitReached()
		if isLimitReached {
			return context.JSON(http.StatusPaymentRequired, "You need to upgrade your plan to use more AI features")
		}
	}

	logger := context.App.Logger
	aiService := context.App.AiService
	// * read the users query from here
	chatId := context.Param("id")
	if chatId == "" {
		return context.JSON(http.StatusBadRequest, "Invalid Chat Id")
	}
	chatUuid, err := uuid.Parse(chatId)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Invalid chat Id")
	}

	userUuid, err := uuid.Parse(context.Session.User.UniqueId)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Invalid user id found")
	}

	logger.Info("User UUID: %v", userUuid, nil)

	// * get the chat from the db
	var dest struct {
		model.AiChat
		Organization       model.Organization
		OrganizationMember model.OrganizationMember
		AiChatMessages     []model.AiChatMessage
	}

	fetchChatQuery := SELECT(
		table.AiChat.AllColumns,
		table.AiChatMessage.AllColumns,
		table.Organization.AllColumns,
		table.OrganizationMember.AllColumns,
	).FROM(
		table.AiChat.LEFT_JOIN(
			table.Organization,
			table.AiChat.OrganizationId.EQ(table.Organization.UniqueId),
		).LEFT_JOIN(
			table.OrganizationMember,
			table.AiChat.OrganizationMemberId.EQ(table.OrganizationMember.UniqueId),
		).LEFT_JOIN(
			table.AiChatMessage,
			table.AiChat.UniqueId.EQ(table.AiChatMessage.AiChatId),
		),
	).WHERE(
		table.AiChat.UniqueId.EQ(UUID(chatUuid)),
	).ORDER_BY(
		table.AiChatMessage.CreatedAt.ASC(),
	).
		LIMIT(15)

	err = fetchChatQuery.Query(context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "Chat not found")
		} else {
			return context.JSON(http.StatusInternalServerError, "Error while generating response")
		}
	}

	payload := new(api_types.AiChatQuerySchema)
	if err := context.Bind(payload); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	logger.Info("User query: %v", payload.Query, nil)

	var insertedUserMessage model.AiChatMessage

	// create the user message record in the db
	userAiChatMessageToInsert := model.AiChatMessage{
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		AiChatId:             dest.UniqueId,
		OrganizationId:       dest.OrganizationId,
		OrganizationMemberId: dest.OrganizationMemberId,
		Content:              payload.Query,
		Role:                 model.AiChatMessageRoleEnum_User,
	}

	userAiChatMessageToInsertQuery := table.AiChatMessage.INSERT(
		table.AiChatMessage.MutableColumns,
	).MODEL(
		userAiChatMessageToInsert,
	).RETURNING(
		table.AiChatMessage.AllColumns,
	)

	err = userAiChatMessageToInsertQuery.QueryContext(context.Request().Context(), context.App.Db, &insertedUserMessage)

	if err != nil {
		logger.Error("Error inserting user message: %v", err.Error(), nil)
		return context.JSON(http.StatusInternalServerError, "Error inserting user message")
	}

	contextMessages := make([]api_types.AiChatMessageSchema, 0)
	for _, message := range dest.AiChatMessages {
		role := api_types.AiChatMessageRoleEnum(message.Role)
		contextMessages = append(contextMessages, api_types.AiChatMessageSchema{
			UniqueId:  message.UniqueId.String(),
			CreatedAt: message.CreatedAt,
			Content:   message.Content,
			Role:      role,
		})
	}

	// * query the AI model
	inputPrompt := aiService.BuildChatBoxQueryInputPrompt(
		payload.Query,
		contextMessages,
		dest.OrganizationId,
	)

	streamingResponse, err := aiService.QueryAiModelWithStreaming(context.Request().Context(), inputPrompt)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Error querying AI model")
	}

	bufferedResponse := ""

	var insertedAiMessage model.AiChatMessage

	aiChatMessageToInsert := model.AiChatMessage{
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		AiChatId:             dest.UniqueId,
		OrganizationId:       dest.OrganizationId,
		OrganizationMemberId: dest.OrganizationMemberId,
		Content:              bufferedResponse,
		Role:                 model.AiChatMessageRoleEnum_Assistant,
	}

	// * insert the message get the id and sent to user
	aiChatMessageToInsertQuery := table.AiChatMessage.INSERT(
		table.AiChatMessage.MutableColumns,
	).MODEL(
		aiChatMessageToInsert,
	).RETURNING(
		table.AiChatMessage.AllColumns,
	)

	err = aiChatMessageToInsertQuery.QueryContext(context.Request().Context(), context.App.Db, &insertedAiMessage)

	if err != nil {
		logger.Error("Error inserting ai message: %v", err.Error(), nil)
		return context.JSON(http.StatusInternalServerError, "Error inserting ai message to database")
	}

	context.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)
	context.Response().WriteHeader(http.StatusOK)
	enc := json.NewEncoder(context.Response())

	// * send user user message and ai message inserted to frontend
	userMessage := api_types.AiChatMessageSchema{
		UniqueId:  insertedUserMessage.UniqueId.String(),
		Content:   insertedUserMessage.Content,
		CreatedAt: insertedUserMessage.CreatedAt,
		Role:      api_types.AiChatMessageRoleEnum(insertedUserMessage.Role),
	}

	aiMessage := api_types.AiChatMessageSchema{
		UniqueId:  insertedAiMessage.UniqueId.String(),
		Content:   insertedAiMessage.Content,
		CreatedAt: insertedAiMessage.CreatedAt,
		Role:      api_types.AiChatMessageRoleEnum(insertedAiMessage.Role),
	}

	messageDetailsToSend := map[string]interface{}{
		"type":        "messageDetails",
		"userMessage": userMessage,
		"aiMessage":   aiMessage,
	}

	if err := enc.Encode(messageDetailsToSend); err != nil {
		context.Response().Flush()
	}

	for response := range streamingResponse.StreamChannel {
		delta := map[string]string{
			"type":    "text-delta",
			"content": response,
		}
		bufferedResponse += response
		if err := enc.Encode(delta); err != nil {
			return err
		}
		context.Response().Flush()
	}

	// we need to get the usage.promptTokens and usage.completionTokens from the response and send in the finished response too

	// After the loop, send a `finish` signal
	finishMessage := map[string]string{
		"type":         "finish",
		"finishReason": "done",
	}

	// * update the message content in the db

	if err := enc.Encode(finishMessage); err != nil {
		return err
	}

	updateQuery := table.AiChatMessage.UPDATE(
		table.AiChatMessage.Content,
	).SET(
		bufferedResponse,
	).WHERE(
		table.AiChatMessage.UniqueId.EQ(UUID(insertedAiMessage.UniqueId)),
	)
	_, err = updateQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		logger.Error("Error updating ai message: %v", err.Error(), nil)
		return context.JSON(http.StatusInternalServerError, "Error updating ai message")
	}

	aiService.LogApiCall(
		dest.OrganizationId,
		context.App.Db,
		payload.Query,
		bufferedResponse,
		model.AiModelEnum(streamingResponse.ModelUsed),
		streamingResponse.InputTokensUsed,
		streamingResponse.OutputTokensUsed,
	)

	return nil
}

func handleGetSegmentRecommendation(context interfaces.ContextWithSession) error {
	return nil
}

// ! this streams response to the client
func handleGetChatSummary(context interfaces.ContextWithSession) error {
	return nil
}
