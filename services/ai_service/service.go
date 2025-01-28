package ai_service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/mistral"
	"github.com/tmc/langchaingo/llms/openai"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/wapikit/wapikit/.db-generated/model"
	"github.com/wapikit/wapikit/.db-generated/table"
	"github.com/wapikit/wapikit/api/api_types"
	cache_service "github.com/wapikit/wapikit/services/redis_service"
	"github.com/wapikit/wapikit/utils"
)

type UserQueryIntent string

const (
	UserIntentCampaigns     UserQueryIntent = "campaigns"
	UserIntentCampaign      UserQueryIntent = "campaign"
	UserIntentGenerateChats UserQueryIntent = "chats"
	UserIntentGenerateChat  UserQueryIntent = "chat"
)

const (
	SYSTEM_PROMPT_AI_CHAT_BOX_QUERY = `You are a AI assistant for a WhatsApp Business Management tool used for sending our marketing campaigns and customer engagement. You will act as a data analyst to provide insights on the data and helps in decision making. You will be provided with the relevant contextual data from the organization database, you responsibility is to provide insights, without any buzz words or jargons. You must use easy and simple sentences.`

	SYSTEM_PROMPT_INTENT_DETECTION = `
You are an AI assistant for a WhatsApp Business Management tool specializing in sending marketing campaigns and customer engagement. Your primary responsibility is to detect the intent in user queries.
- The intent can take one of the following values: "campaign" (for data related to campaigns, such as insights or performance) or "chats" (for data related to customer conversations or replies).
- If a query indirectly refers to marketing efforts, customer engagement, campaign strategy, or responses from customers, you should infer whether the intent requires "campaign" data, "chat" data, or both.
- For each detected intent, provide the following in a JSON string:
  - The "intent": either "campaigns" or "chats".
  - The "startDate" and "endDate" (in UTC string). If the user query does not provide a specific timeframe, leave the dates as empty strings.
- If the intent cannot be determined, return an empty JSON object: {}.
- Always interpret queries with a "marketing and customer engagement perspective," even when questions are indirect or open-ended.
`

	SYSTEM_PROMPT_CHAT_SUMMARY_GENERATION = "You are an AI assistant for a WhatsApp Business Management tool used for sending marketing campaign and customer engagement. You will be provided with the chat messages between user and a internal organization member, and you responsibility is to generate a summary of the chat. The summary should be not be more than 500 words. It should clear and concise, with easy english and no jargons and try to answer in bullet points. The summary should depict the main points of the chat and the conclusion of the chat. And finally, what actionable steps can be taken from the chat."

	SYSTEM_PROMPT_RESPONSE_SUGGESTION_GENERATION = "You are an AI assistant for a WhatsApp Business Management tool used for sending marketing campaign and customer engagement. You will be provided with the chat messages between user and a internal organization member, and you responsibility is to generate a response to the chat. The response should be not be more than 500 words. It should clear and concise, with easy english and no jargons. The response should be in a way that it should be able to answer the query of the user and also provide a solution to the query from the perspective of the organization. You must response back with a json string with response property where response is an array of strings where each element of array is a response suggestion to the contact message."

	SYSTEM_PROMPT_SEGMENT_SUGGESTIONS = "You are an AI assistant for a WhatsApp Business Management tool used for sending marketing campaign and customer engagement. You will be provided with the conversation or a contact details from the application. Your responsibility is to generate a segment suggestion for the entity provided. You should response back with a json string with tags properties where tags are array of string. Keep them shorter in length, the tag strings will be used to create tags in the backend where your provided string will be label of the tag and will be used as segment identifiers. You should provide a maximum of 5 tags. You should not use any buzz words or jargons. Example are: ['VIP', 'High Potential', 'Low Potential', 'Engaged', 'Not Engaged', 'Potential Customer', 'Existing Customer', 'Needs Support', 'Urgent']"
)

var AiModelEnumToLlmModelString = map[api_types.AiModelEnum]string{
	api_types.GPT4Mini:    "gpt-4o-mini",
	api_types.Gpt35Turbo:  "gpt-3.5-turbo",
	api_types.Gpt4o:       "gpt-4o",
	api_types.Mistral:     "mistral",
	api_types.Gemini15Pro: "gemini-15-pro",
	api_types.Claude35:    "claude-3-5-haiku-latest",
}

type AiService struct {
	Logger *slog.Logger
	Redis  *cache_service.RedisClient
	Db     *sql.DB
	ApiKey string
}

func NewAiService(
	logger *slog.Logger,
	redis *cache_service.RedisClient,
	db *sql.DB,
	apiKey string,
) *AiService {
	return &AiService{
		Logger: logger,
		Redis:  redis,
		Db:     db,
		ApiKey: apiKey,
	}
}

func (ai *AiService) FetchRelevantData(organizationId uuid.UUID, intentDetails *DetectIntentResponse, ctx context.Context, db *sql.DB) (string, error) {

	switch intentDetails.Intent {
	case UserIntentCampaign:
		{
			var dest []struct {
				model.Campaign
				MessagesSent           int    `json:"messagesSent"`
				MessagesRead           int    `json:"messagesRead"`
				MessagesReplied        int    `json:"messagesReplied"`
				TemplateUsed           string `json:"templateUsed"`
				MessagesFailedToBeSent int    `json:"messagesFailedToBeSent"`
				Tags                   []struct {
					model.Tag
				}
				Lists []struct {
					model.ContactList
					NumberOfContacts int `json:"numberOfContacts"`
				}
			}

			whereCondition := table.Campaign.OrganizationId.EQ(UUID(organizationId))

			// ! TODO: timestamp impression fix
			// if !intentDetails.StartDate.IsZero() {
			// 	whereCondition = whereCondition.AND(table.Campaign.CreatedAt.GT_EQ(intentDetails.StartDate))
			// }

			// if !intentDetails.EndDate.IsZero() {
			// 	whereCondition = whereCondition.AND(table.Campaign.CreatedAt.LT_EQ(intentDetails.EndDate))
			// }

			campaignQuery := SELECT(
				table.Campaign.AllColumns,
				table.Tag.AllColumns,
				table.CampaignList.AllColumns,
				table.ContactList.AllColumns,
				table.CampaignTag.AllColumns,
				COUNT(table.Campaign.UniqueId).OVER().AS("totalCampaigns"),
			).
				FROM(table.Campaign.
					LEFT_JOIN(table.CampaignTag, table.CampaignTag.CampaignId.EQ(table.Campaign.UniqueId)).
					LEFT_JOIN(table.Tag, table.Tag.UniqueId.EQ(table.CampaignTag.TagId)).
					LEFT_JOIN(table.CampaignList, table.CampaignList.CampaignId.EQ(table.Campaign.UniqueId)).
					LEFT_JOIN(table.ContactList, table.ContactList.UniqueId.EQ(table.CampaignList.ContactListId)),
				).
				WHERE(whereCondition)

			err := campaignQuery.QueryContext(ctx, db, &dest)

			if err != nil {
				return "", err
			}

			campaignJson, _ := json.Marshal(dest)
			return string(campaignJson), nil
		}

	case UserIntentGenerateChats:
		{

			var dest []struct {
				model.Conversation
				Contact    model.Contact `json:"contact"`
				AssignedTo struct {
					Member model.OrganizationMember `json:"member"`
				}
				Messages []model.Message `json:"messages"`
			}

			// ! TODO: implement this

			_ = dest

			return "", nil
		}

	}
	return "", nil
}

type AiQueryResponse struct {
	Content         string
	InputTokenUsed  int
	OutputTokenUsed int
}

func (ai *AiService) GetResponseSuggestions(ctx context.Context, messages []model.Message) ([]string, error) {
	var responseSuggestions struct {
		Response []string `json:"response"`
	}

	systemPrompt := llms.TextParts(llms.ChatMessageTypeSystem, SYSTEM_PROMPT_RESPONSE_SUGGESTION_GENERATION)
	userPrompt := llms.TextParts(llms.ChatMessageTypeHuman, "Generate response suggestions for the chat messages")
	inputPrompt := []llms.MessageContent{
		systemPrompt,
	}

	// ! TODO: what if messages has attachments ?
	for _, message := range messages {
		if message.Direction == model.MessageDirectionEnum_InBound {
			inputPrompt = append(inputPrompt, llms.TextParts(llms.ChatMessageTypeHuman, *message.MessageData))
		} else {
			inputPrompt = append(inputPrompt, llms.TextParts(llms.ChatMessageTypeGeneric, *message.MessageData))
		}
	}

	inputPrompt = append(inputPrompt, userPrompt)

	aiResponse, err := ai.QueryAiModel(ctx, api_types.Gpt35Turbo, inputPrompt)

	if err != nil {
		return []string{}, err
	}

	err = json.Unmarshal([]byte(aiResponse.Content), &responseSuggestions)

	if err != nil {
		fmt.Println("Error unmarshalling response suggestions", err)
		return []string{}, err
	}

	return responseSuggestions.Response, nil
}

func (ai *AiService) GetLlmFromModel(ctx context.Context, model api_types.AiModelEnum) (llms.Model, error) {
	var llm llms.Model
	var err error
	err = nil

	if model == api_types.Gemini15Pro {
		llm, err = googleai.New(ctx, googleai.WithAPIKey(ai.ApiKey))
		return llm, err
	} else if model == api_types.Mistral {
		llm, err = mistral.New(mistral.WithAPIKey(ai.ApiKey))
		return llm, err
	} else {
		llm, err = openai.New(openai.WithModel(AiModelEnumToLlmModelString[model]), openai.WithToken(ai.ApiKey))
		return llm, err
	}
}

func (ai *AiService) QueryAiModel(ctx context.Context, model api_types.AiModelEnum, inputPrompt []llms.MessageContent) (*AiQueryResponse, error) {
	var llm llms.Model

	llm, err := ai.GetLlmFromModel(ctx, model)

	if err != nil {
		ai.Logger.Error("Error creating OpenAI model in query AI model function", err)
	}

	completion, err := llm.GenerateContent(ctx,
		inputPrompt,
	)

	if err != nil {
		ai.Logger.Error("Error generating content from AI model", err)
		return nil, err
	}

	rawJson, _ := json.Marshal(completion)
	ai.Logger.Info("AI response", string(rawJson))

	var response AiQueryResponse
	for _, choice := range completion.Choices {
		response.Content = choice.Content
		response.InputTokenUsed = choice.GenerationInfo["PromptTokens"].(int)
		response.OutputTokenUsed = choice.GenerationInfo["CompletionTokens"].(int)
	}

	return &response, nil
}

func (ai *AiService) BuildChatBoxQueryInputPrompt(query string, contextMessages []api_types.AiChatMessageSchema, dataContext *string) []llms.MessageContent {
	systemPrompt := llms.TextParts(llms.ChatMessageTypeSystem, SYSTEM_PROMPT_AI_CHAT_BOX_QUERY)
	userPrompt := llms.TextParts(llms.ChatMessageTypeHuman, query)
	inputPrompt := []llms.MessageContent{
		systemPrompt,
	}

	for _, message := range contextMessages {
		if message.Role == api_types.Assistant {
			inputPrompt = append(inputPrompt, llms.TextParts(llms.ChatMessageTypeAI, message.Content))
		} else if message.Role == api_types.User {
			inputPrompt = append(inputPrompt, llms.TextParts(llms.ChatMessageTypeHuman, message.Content))
		}
	}

	if dataContext != nil || *dataContext != "" {
		fullContextText := strings.Join([]string{"Heres the data you may need:", *dataContext}, " ")
		inputPrompt = append(inputPrompt, llms.TextParts(llms.ChatMessageTypeHuman, fullContextText))
	}

	inputPrompt = append(inputPrompt, userPrompt)

	jsonInputPrompt, _ := json.Marshal(inputPrompt)
	ai.Logger.Info("Input prompt for AI model", string(jsonInputPrompt))

	return inputPrompt
}

type StreamingResult struct {
	StreamChannel    <-chan string
	InputTokensUsed  int
	OutputTokensUsed int
}

func (ai *AiService) QueryAiModelWithStreaming(ctx context.Context, model api_types.AiModelEnum, inputPrompt []llms.MessageContent) (*StreamingResult, error) {
	streamChannel := make(chan string)
	tokenChannel := make(chan struct {
		inputTokens  int
		outputTokens int
	})

	model = api_types.Gpt35Turbo
	llm, err := ai.GetLlmFromModel(ctx, model)

	if err != nil {
		return nil, err
	}

	go func() {
		defer close(streamChannel)
		defer close(tokenChannel)

		resp, err := llm.GenerateContent(ctx,
			inputPrompt,
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				streamChannel <- string(chunk)
				return nil
			}),
		)

		if err != nil {
			log.Fatal(err)
		}

		inputTokenUsed := 0
		outputTokenUsed := 0

		for _, choice := range resp.Choices {
			inputTokenUsed = choice.GenerationInfo["PromptTokens"].(int)
			outputTokenUsed = choice.GenerationInfo["CompletionTokens"].(int)
		}

		tokenChannel <- struct {
			inputTokens  int
			outputTokens int
		}{inputTokens: inputTokenUsed, outputTokens: outputTokenUsed}
	}()

	result := &StreamingResult{
		StreamChannel: streamChannel,
	}

	go func() {
		tokens := <-tokenChannel
		result.InputTokensUsed = tokens.inputTokens
		result.OutputTokensUsed = tokens.outputTokens
	}()

	return result, nil
}

type DetectIntentResponse struct {
	Intent    UserQueryIntent `json:"intent"`
	StartDate *time.Time      `json:"startDate"`
	EndDate   *time.Time      `json:"endDate"`
}

func (ai *AiService) DetectIntent(query string, organizationId string) (*DetectIntentResponse, error) {
	systemPromptContent := strings.Join([]string{SYSTEM_PROMPT_INTENT_DETECTION, "Current time and date is", utils.GetCurrentTimeAndDateInUTCString()}, " ")
	systemPrompt := llms.TextParts(llms.ChatMessageTypeSystem, systemPromptContent)
	userPrompt := llms.TextParts(llms.ChatMessageTypeHuman, query)

	inputPrompt := []llms.MessageContent{
		systemPrompt,
		userPrompt,
	}

	intentResponse, err := ai.QueryAiModel(context.Background(), api_types.Gpt35Turbo, inputPrompt)

	fmt.Println("Intent response", intentResponse.Content)

	if err != nil {
		return nil, err
	}

	var detectIntentResponse DetectIntentResponse
	err = json.Unmarshal([]byte(intentResponse.Content), &detectIntentResponse)
	if err != nil {
		fmt.Println("Error unmarshalling intent response", err)
		return nil, err
	}

	// * log the API call
	ai.LogApiCall(uuid.MustParse(organizationId), ai.Db, query, intentResponse.Content, intentResponse.InputTokenUsed, intentResponse.OutputTokenUsed)
	return &detectIntentResponse, nil
}

func (ai *AiService) LogApiCall(organizationId uuid.UUID, db *sql.DB, request, response string, inputTokenUsed, outputTokenUsed int) error {

	fmt.Println("Logging API call")

	apiLogToInsert := model.AiApiCallLogs{
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		OrganizationId:  organizationId,
		Request:         request,
		Response:        response,
		InputTokenUsed:  int32(inputTokenUsed),
		OutputTokenUsed: int32(outputTokenUsed),
	}

	insertQuery := table.AiApiCallLogs.INSERT(
		table.AiApiCallLogs.MutableColumns,
	).MODEL(
		apiLogToInsert,
	).RETURNING(
		table.AiApiCallLogs.AllColumns,
	)

	_, err := insertQuery.Exec(db)

	if err != nil {
		fmt.Println("Error inserting API log: %v", err)
		return err
	}

	return nil
}

func (ai *AiService) CacheAiResponse(query string, response string) error {
	err := ai.Redis.CacheData(query, response, 0)
	if err != nil {
		return err
	}
	return nil
}
