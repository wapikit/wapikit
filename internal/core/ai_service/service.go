package ai_service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"

	"github.com/wapikit/wapikit/.db-generated/model"
	"github.com/wapikit/wapikit/.db-generated/table"
	"github.com/wapikit/wapikit/internal/api_types"
)

type UserQueryIntent string

const (
	UserIntentCampaignInsights    UserQueryIntent = "campaign_insights"
	UserIntentGenerateChatSummary UserQueryIntent = "generate_summary"
)

func FetchRelevantData(intent UserQueryIntent, orgID uuid.UUID, userID uuid.UUID) (json.RawMessage, error) {
	// Query embeddings or database based on intent
	// Example: Fetch last 30 days of campaign insights

	return nil, nil
}

func QueryOpenAi() {
}

func QueryAiModelWithStreaming(ctx context.Context, model api_types.AiModelEnum, input string) (<-chan string, error) {
	streamChannel := make(chan string)
	go func() {
		defer close(streamChannel)

		model := api_types.Gpt35Turbo

		messageContent := []llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, "You are a AI assistant for a WhatsApp Business Management tool used for sending our marketing campaigns. You will act as a data analyst to provide insights on the data and helps in decision making."),
			llms.TextParts(llms.ChatMessageTypeHuman, input),
		}

		switch model {
		case api_types.Gpt35Turbo:
			{
				llm, err := openai.New(
					openai.WithModel("gpt-3.5-turbo"),
					openai.WithToken(""),
				)
				if err != nil {
					log.Fatal(err)
				}

				resp, err := llm.GenerateContent(ctx,
					messageContent,
					llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
						fmt.Printf("Received chunk: %s\n", chunk)

						// Received chunk: Hello
						// Received chunk: !
						// Received chunk:  How
						// Received chunk:  can
						// Received chunk:  I
						// Received chunk:  assist
						// Received chunk:  you
						// Received chunk:  today
						// Received chunk:  with
						// Received chunk:  the
						// Received chunk:  data
						// Received chunk:  analysis
						// Received chunk:  for
						// Received chunk:  your
						// Received chunk:  marketing
						// Received chunk:  campaigns
						// Received chunk: ?
						// add this chunk to the channel
						streamChannel <- string(chunk)
						return nil
					}),
					// ! TODO: may be add a tool which can get the data from database
					// llms.WithTools(tools)
				)

				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Response:", *resp)
			}

		case api_types.Mistral:
			{
				llm, err := ollama.New(ollama.WithModel("mistral"))
				if err != nil {
					log.Fatal(err)
				}

				completion, err := llm.GenerateContent(
					ctx, messageContent,
					llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
						fmt.Print(string(chunk))
						streamChannel <- string(chunk)
						return nil
					}))

				if err != nil {
					log.Fatal(err)
				}
				_ = completion
			}

		case api_types.Gpt4o:
			{

			}

		case api_types.GPT4Mini:
			{

			}

		case api_types.Gemini15Pro:
			{

			}

		default:
			log.Println("Unsupported model")
		}
	}()

	return streamChannel, nil
}

func DetectIntent(query string) (UserQueryIntent, error) {
	// keywords := map[UserQueryIntent][]string{
	// 	UserIntentCampaignInsights:    {"campaign", "insights", "last 30 days"},
	// 	UserIntentGenerateChatSummary: {"summary", "chat", "conversation"},
	// }

	return UserIntentCampaignInsights, nil

	// for intent, words := range keywords {
	// 	for _, word := range words {
	// 		if strings.Contains(strings.ToLower(query), word) {
	// 			return intent, nil
	// 		}
	// 	}
	// }
	// return "", fmt.Errorf("intent not detected")
}

func GenerateEmbedding(content string, model api_types.AiModelEnum) ([]float64, error) {
	// Call embedding model and return vector representation
	return nil, nil
}

func LogApiCall(aiChatId uuid.UUID, db *sql.DB, request, response string) error {
	apiLogToInsert := model.AiApiCallLogs{
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		AiChatId:        aiChatId,
		Request:         request,
		Response:        response,
		InputTokenUsed:  0,
		OutputTokenUsed: 0,
	}

	insertQuery := table.AiApiCallLogs.INSERT(
		table.AiApiCallLogs.MutableColumns,
	).MODEL(
		&apiLogToInsert,
	).RETURNING(
		table.AiApiCallLogs.AllColumns,
	)

	_, err := insertQuery.Exec(db)

	if err != nil {
		return err
	}

	return nil
}

func CheckAiRateLimit() bool {
	return false
}

func GetTotalAiTokenConsumedByOrganization(orgUuid uuid.UUID) int {
	return 0
}

func GetTotalAiTokenConsumedByUser(memberUuid uuid.UUID) int {
	return 0
}
