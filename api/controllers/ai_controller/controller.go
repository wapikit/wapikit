package ai_controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/wapikit/wapikit/.db-generated/model"
	controller "github.com/wapikit/wapikit/api/controllers"
	"github.com/wapikit/wapikit/internal/api_types"
	"github.com/wapikit/wapikit/internal/core/ai_service"
	"github.com/wapikit/wapikit/internal/core/utils"
	"github.com/wapikit/wapikit/internal/interfaces"

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
							MaxRequests:    10,
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
							MaxRequests:    10,
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
							MaxRequests:    10,
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
							MaxRequests:    10,
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
							MaxRequests:    10,
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
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
			},
		},
	}
}

func handleGetChats(context interfaces.ContextWithSession) error {
	params := new(api_types.GetAiChatsParams)

	err := utils.BindQueryParams(context, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	pageNumber := params.Page
	pageSize := params.PerPage

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)
	if err != nil {
		return context.String(http.StatusInternalServerError, "Error parsing organization UUID")
	}
	userUuid, err := uuid.Parse(context.Session.User.UniqueId)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Error parsing user UUID")
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
			Description: *chat.Description,
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
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Chat Id")
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
			return echo.NewHTTPError(http.StatusNotFound, "Chat not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching chat")
		}
	}

	responseToReturn := api_types.GetAiChatByIdResponseSchema{
		Chat: api_types.AiChatSchema{
			UniqueId:    dest.UniqueId.String(),
			CreatedAt:   dest.CreatedAt,
			Description: *dest.Description,
			Title:       dest.Title,
		},
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func handleGetChatMessages(context interfaces.ContextWithSession) error {
	chatId := context.Param("id")
	if chatId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Chat Id")
	}

	chatUuid, _ := uuid.Parse(chatId)

	var dest []struct {
		model.AiChatMessage
	}

	fetchMessagesQuery := SELECT(
		table.AiChatMessage.AllColumns,
	).FROM(
		table.AiChatMessage,
	).WHERE(
		table.AiChatMessage.AiChatId.EQ(UUID(chatUuid)),
	)

	err := fetchMessagesQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return echo.NewHTTPError(http.StatusNotFound, "Messages not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching messages")
		}
	}

	messagesToReturn := []api_types.AiChatMessageSchema{}

	for _, message := range dest {
		role := api_types.AiChatMessageRoleEnum(message.Role)
		messagesToReturn = append(messagesToReturn, api_types.AiChatMessageSchema{
			UniqueId:  message.UniqueId.String(),
			CreatedAt: message.CreatedAt,
			// Content:   message.Content,
			Role: role,
		})
	}

	return context.JSON(http.StatusOK, api_types.GetAiChatMessagesResponseSchema{
		Messages: messagesToReturn,
	})
}

func voteMessage(context interfaces.ContextWithSession) error {
	return nil
}

func getMessageVotes(context interfaces.ContextWithSession) error {
	params := new(api_types.GetAiChatMessageVotesParams)

	err := utils.BindQueryParams(context, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	pageNumber := params.Page
	pageSize := params.PerPage

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)
	if err != nil {
		return context.String(http.StatusInternalServerError, "Error parsing organization UUID")
	}
	userUuid, err := uuid.Parse(context.Session.User.UniqueId)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Error parsing user UUID")
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

func handleReplyToChat(context interfaces.ContextWithSession) error {

	logger := context.App.Logger

	// * read the users query from here
	chatId := context.Param("id")
	if chatId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Chat Id")
	}

	chatUuid, err := uuid.Parse(chatId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid chat Id")
	}

	userUuid, err := uuid.Parse(context.Session.User.UniqueId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid user id found")
	}

	logger.Info("User UUID: %v", userUuid, nil)

	// * get the chat from the db
	var dest struct {
		model.AiChat
		Organization       model.Organization       `alias:"organization"`
		OrganizationMember model.OrganizationMember `alias:"organization_member"`
	}

	fetchChatQuery := SELECT(
		table.AiChat.AllColumns,
		table.Organization.AllColumns,
		table.OrganizationMember.AllColumns,
	).FROM(
		table.AiChat.LEFT_JOIN(
			table.Organization,
			table.AiChat.OrganizationId.EQ(table.Organization.UniqueId),
		).LEFT_JOIN(
			table.OrganizationMember,
			table.AiChat.OrganizationMemberId.EQ(table.OrganizationMember.UniqueId),
		),
	).WHERE(
		table.AiChat.UniqueId.EQ(UUID(chatUuid)),
	).LIMIT(1)

	err = fetchChatQuery.Query(context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return echo.NewHTTPError(http.StatusNotFound, "Chat not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error while generating response")
		}
	}

	payload := new(api_types.AiChatQuerySchema)
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	logger.Info("User query: %v", payload.Query, nil)

	// * check the limit
	isLimitReached := ai_service.CheckAiRateLimit()

	if isLimitReached {
		return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit reached")
	}

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

	_, err = userAiChatMessageToInsertQuery.ExecContext(context.Request().Context(), context.App.Db)

	// * check the intent of ths user
	// userIntent, err := ai_service.DetectIntent(payload.Query)

	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "Error detecting intent")
	// }

	// * get the corresponding context from the db like campaign right now, we will only support campaign
	// dataContext, err := ai_service.FetchRelevantData(userIntent, dest.OrganizationId, userUuid)

	// logger.Info("Data context: %v", dataContext, nil)

	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching data")
	// }

	// * create the prompt to pass to AI with the context and the user query
	var prompt string

	// * call the AI API with the context and the user query
	streamChannel, err := ai_service.QueryAiModelWithStreaming(context.Request().Context(), api_types.Gpt35Turbo, prompt)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error querying AI model")
	}

	bufferedResponse := ""

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

	_, err = aiChatMessageToInsertQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error inserting message")
	}

	context.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)
	context.Response().WriteHeader(http.StatusOK)
	enc := json.NewEncoder(context.Response())

	for response := range streamChannel {
		delta := map[string]string{
			"type":      "text-delta",
			"textDelta": response,
		}

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
		"content":      "",
	}

	if err := enc.Encode(finishMessage); err != nil {
		return err
	}

	return nil
}
