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

	pii_redactor "github.com/wapikit/wapikit/services/pii_redactor_service"
	cache_service "github.com/wapikit/wapikit/services/redis_service"
	"github.com/wapikit/wapikit/utils"

	"github.com/wapikit/wapikit/.db-generated/model"
	"github.com/wapikit/wapikit/.db-generated/table"
	"github.com/wapikit/wapikit/api/api_types"
)

// Expanded intent types and constants
type UserQueryIntent string

const (
	UserIntentCampaigns    UserQueryIntent = "campaigns"
	UserIntentConversation UserQueryIntent = "conversation"
	UserIntentChat         UserQueryIntent = "chat"
)

// Enhanced system prompts with versioning
const (
	SYSTEM_PROMPT_AI_CHAT_BOX_QUERY = `You are a AI assistant for a WhatsApp Business Management tool used for sending our marketing campaigns and customer engagement. You will act as a data analyst to provide insights on the data and helps in decision making. You will be provided with the relevant contextual data from the organization database, you responsibility is to provide insights, without any buzz words or jargons. You must use easy and simple sentences.`

	SYSTEM_PROMPT_INTENT_DETECTION = `
You are an AI assistant for a WhatsApp Business Management tool specialized in sending marketing campaigns and customer engagement. Your primary task is to analyze user queries and determine their intent. For each query, produce a JSON object with the following keys:
- "primaryIntent": A string representing the main intent. It should be either "campaigns" (for queries related to campaign performance, insights, or strategy) or "chats" (for queries related to customer conversations or replies).
- "confidence": A number between 0 and 1 that reflects your confidence in the detected intent.
- "entities": An array of objects, where each object has a "type" (such as "date", "keyword", or "location") and a "value" that provides the recognized detail from the query.
- "temporalContext": An object with the keys "start", "end", and "timezone". If the query specifies a timeframe, "start" and "end" should be the UTC timestamps; otherwise, they should be empty strings and "timezone" should be "UTC".

If you cannot determine the intent, return an empty JSON object: {}.

Always interpret queries from a marketing and customer engagement perspective, and use clear, simple language in your analysis.
`

	SYSTEM_PROMPT_CHAT_SUMMARY_GENERATION = "You are an AI assistant for a WhatsApp Business Management tool used for sending marketing campaign and customer engagement. You will be provided with the chat messages between user and a internal organization member, and you responsibility is to generate a summary of the chat. The summary should be not be more than 500 words. It should clear and concise, with easy english and no jargons and try to answer in bullet points. The summary should depict the main points of the chat and the conclusion of the chat. And finally, what actionable steps can be taken from the chat."

	SYSTEM_PROMPT_RESPONSE_SUGGESTION_GENERATION = "You are an AI assistant for a WhatsApp Business Management tool used for sending marketing campaign and customer engagement. You will be provided with the chat messages between user and a internal organization member, and you responsibility is to generate a response to the chat. The response should be not be more than 500 words. It should clear and concise, with easy english and no jargons. The response should be in a way that it should be able to answer the query of the user and also provide a solution to the query from the perspective of the organization. You must response back with a json string with response property where response is an array of strings where each element of array is a response suggestion to the contact message."

	SYSTEM_PROMPT_SEGMENT_SUGGESTIONS = "You are an AI assistant for a WhatsApp Business Management tool used for sending marketing campaign and customer engagement. You will be provided with the conversation or a contact details from the application. Your responsibility is to generate a segment suggestion for the entity provided. You should response back with a json string with tags properties where tags are array of string. Keep them shorter in length, the tag strings will be used to create tags in the backend where your provided string will be label of the tag and will be used as segment identifiers. You should provide a maximum of 5 tags. You should not use any buzz words or jargons. Example are: ['VIP', 'High Potential', 'Low Potential', 'Engaged', 'Not Engaged', 'Potential Customer', 'Existing Customer', 'Needs Support', 'Urgent']"
)

// New data structures for enhanced functionality
type Entity struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type DetectIntentResponse struct {
	PrimaryIntent   UserQueryIntent `json:"primaryIntent"`
	Confidence      float64         `json:"confidence"`
	Entities        []Entity        `json:"entities"`
	TemporalContext TemporalRange   `json:"temporalContext"`
}

type TemporalRange struct {
	Start    *time.Time `json:"start"`
	End      *time.Time `json:"end"`
	Timezone string     `json:"timezone"`
}

// UnmarshalJSON implements a custom JSON unmarshaler for TemporalRange.
// It treats empty strings for "start" or "end" as nil.
func (tr *TemporalRange) UnmarshalJSON(data []byte) error {
	// Create a temporary alias type with Start and End as strings.
	var aux struct {
		Start    string `json:"start"`
		End      string `json:"end"`
		Timezone string `json:"timezone"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse start time if provided.
	if aux.Start != "" {
		t, err := time.Parse(time.RFC3339, aux.Start)
		if err != nil {
			return err
		}
		tr.Start = &t
	} else {
		tr.Start = nil
	}

	// Parse end time if provided.
	if aux.End != "" {
		t, err := time.Parse(time.RFC3339, aux.End)
		if err != nil {
			return err
		}
		tr.End = &t
	} else {
		tr.End = nil
	}

	tr.Timezone = aux.Timezone
	return nil
}

type AiService struct {
	Logger         *slog.Logger
	Redis          *cache_service.RedisClient
	Db             *sql.DB
	ApiKey         string
	DefaultAiModel api_types.AiModelEnum
	ModelRouter    *ModelRouter
	ContextBuilder *ContextBuilder
	Redactor       *pii_redactor.PIIRedactor
	VectorCache    *RedisVectorCache
}

// Enhanced initialization
func NewAiService(
	logger *slog.Logger,
	redis *cache_service.RedisClient,
	db *sql.DB,
	apiKey string,
	defaultModel api_types.AiModelEnum,
) *AiService {
	return &AiService{
		Logger:         logger,
		Redis:          redis,
		Db:             db,
		ApiKey:         apiKey,
		DefaultAiModel: defaultModel,
		ModelRouter:    NewModelRouter(),
		VectorCache:    NewRedisVectorCache(redis),
		Redactor:       pii_redactor.NewPIIRedactor(),
		ContextBuilder: NewContextBuilder(db, logger),
	}
}

func (ai *AiService) DetectIntent(query string, organizationId uuid.UUID) (*DetectIntentResponse, error) {
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
	ai.LogApiCall(organizationId, ai.Db, query, intentResponse.Content, model.AiModelEnum(api_types.Gpt35Turbo), intentResponse.InputTokenUsed, intentResponse.OutputTokenUsed)
	return &detectIntentResponse, nil
}

func (ai *AiService) buildIntentDetectionPrompt(query string) []llms.MessageContent {
	systemPrompt := llms.TextParts(llms.ChatMessageTypeSystem, SYSTEM_PROMPT_INTENT_DETECTION)
	userPrompt := llms.TextParts(llms.ChatMessageTypeHuman, query)

	return []llms.MessageContent{
		systemPrompt,
		userPrompt,
	}
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

type AiQueryResponse struct {
	Content         string
	InputTokenUsed  int
	OutputTokenUsed int
}

func (ai *AiService) QueryAiModel(ctx context.Context, model api_types.AiModelEnum, inputPrompt []llms.MessageContent) (*AiQueryResponse, error) {
	var llm llms.Model

	llm, err := ai.ModelRouter.SelectModel(ctx, model, ai.ApiKey)

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

// Security enhancements
func (ai *AiService) SanitizeInput(input string) string {
	return ai.Redactor.Redact(input)
}

// BuildChatBoxQueryInputPrompt builds an input prompt for the chat box query.
func (ai *AiService) BuildChatBoxQueryInputPrompt(query string, contextMessages []api_types.AiChatMessageSchema, orgId uuid.UUID) []llms.MessageContent {

	intent, err := ai.DetectIntent(query, orgId)
	context := ""

	if err != nil {
		ai.Logger.Error("Error detecting intent", err)
		// * continue without context
	} else {
		ai.Logger.Info("Detected intent", intent)
		context = ai.ContextBuilder.fetchRelevantContext(orgId, *intent)
	}

	systemPrompt := llms.TextParts(llms.ChatMessageTypeSystem, SYSTEM_PROMPT_AI_CHAT_BOX_QUERY)
	userPrompt := llms.TextParts(llms.ChatMessageTypeHuman, query)
	inputPrompt := []llms.MessageContent{systemPrompt}

	for _, message := range contextMessages {
		if message.Role == api_types.Assistant {
			inputPrompt = append(inputPrompt, llms.TextParts(llms.ChatMessageTypeAI, message.Content))
		} else if message.Role == api_types.User {
			inputPrompt = append(inputPrompt, llms.TextParts(llms.ChatMessageTypeHuman, message.Content))
		}
	}
	if context != "" {
		fullContextText := "Here's the data you may need: " + context
		inputPrompt = append(inputPrompt, llms.TextParts(llms.ChatMessageTypeHuman, fullContextText))
	}

	inputPrompt = append(inputPrompt, userPrompt)
	jsonInputPrompt, _ := json.Marshal(inputPrompt)
	ai.Logger.Info("Input prompt for AI model", string(jsonInputPrompt))
	return inputPrompt
}

type StreamingResult struct {
	StreamChannel    <-chan string
	ModelUsed        api_types.AiModelEnum
	InputTokensUsed  int
	OutputTokensUsed int
}

func (ai *AiService) QueryAiModelWithStreaming(ctx context.Context, inputPrompt []llms.MessageContent) (*StreamingResult, error) {
	streamChannel := make(chan string)
	tokenChannel := make(chan struct {
		inputTokens  int
		outputTokens int
	})

	model := ai.DefaultAiModel
	llm, err := ai.ModelRouter.SelectModel(ctx, model, ai.ApiKey)

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
		ModelUsed:     model,
	}

	go func() {
		tokens := <-tokenChannel
		result.InputTokensUsed = tokens.inputTokens
		result.OutputTokensUsed = tokens.outputTokens
	}()

	return result, nil
}

func (ai *AiService) LogApiCall(organizationId uuid.UUID, db *sql.DB, request, response string, aiModel model.AiModelEnum, inputTokenUsed, outputTokenUsed int) error {
	fmt.Println("Logging API call")

	apiLogToInsert := model.AiApiCallLogs{
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		OrganizationId:  organizationId,
		Request:         request,
		Response:        response,
		Model:           model.AiModelEnum(ai.DefaultAiModel),
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
