package ai_service

import (
	"context"
	"errors"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/mistral"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/wapikit/wapikit/api/api_types"
)

// Enhanced model mapping with performance metrics
var AiModelEnumToLlmModelConfig = map[api_types.AiModelEnum]ModelConfig{
	api_types.GPT4Mini:    {Name: "gpt-4o-mini", RequiresAPIKey: true},
	api_types.Gpt35Turbo:  {Name: "gpt-3.5-turbo", RequiresAPIKey: true},
	api_types.Gpt4o:       {Name: "gpt-4o", RequiresAPIKey: true},
	api_types.Mistral:     {Name: "mistral", RequiresAPIKey: true},
	api_types.Gemini15Pro: {Name: "gemini-15-pro", RequiresAPIKey: true},
	api_types.Claude35:    {Name: "claude-3-5-haiku-latest", RequiresAPIKey: true},
}

// ModelConfig holds configuration details for an AI model.
type ModelConfig struct {
	Name           string
	RequiresAPIKey bool
}

// ModelRouter manages the selection of AI models.
type ModelRouter struct {
	models map[api_types.AiModelEnum]ModelConfig
}

// NewModelRouter initializes and returns a ModelRouter with predefined models.
func NewModelRouter() *ModelRouter {
	return &ModelRouter{
		models: AiModelEnumToLlmModelConfig,
	}
}

// SelectModel returns the appropriate llms.Model based on the provided AiModelEnum and API key.
// ! TODO: in future for cloud edition we will add a logic to use multiple modal for multiple tasks
func (mr *ModelRouter) SelectModel(ctx context.Context, model api_types.AiModelEnum, apiKey string) (llms.Model, error) {
	config, exists := mr.models[model]
	if !exists {
		return nil, errors.New("model configuration not found")
	}

	if config.RequiresAPIKey && apiKey == "" {
		return nil, errors.New("API key is required for the selected model")
	}

	switch model {
	case api_types.Gemini15Pro:
		return googleai.New(ctx, googleai.WithAPIKey(apiKey))
	case api_types.Mistral:
		return mistral.New(mistral.WithAPIKey(apiKey))
	default:
		return openai.New(openai.WithModel(config.Name), openai.WithToken(apiKey))
	}
}

// SelectModelsForTask returns a list of models to try for a given task.
// For simplicity, we return all available models here.
func (mr *ModelRouter) SelectModelsForTask(taskType string) []api_types.AiModelEnum {
	var models []api_types.AiModelEnum
	for model := range mr.models {
		models = append(models, model)
	}
	return models
}
