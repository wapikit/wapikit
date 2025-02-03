//go:build community_edition
// +build community_edition

package controller

import (
	"github.com/wapikit/wapikit/.db-generated/model"
	"github.com/wapikit/wapikit/api/api_types"
	"github.com/wapikit/wapikit/interfaces"
	"github.com/wapikit/wapikit/services/ai_service"
)

func mountServices(app *interfaces.App, org *model.Organization) {
	if org.IsAiEnabled && !app.Constants.IsCloudEdition {
		aiService := ai_service.NewAiService(&app.Logger, app.Redis, app.Db, org.AiApiKey, api_types.AiModelEnum(*org.AiModel))
		app.AiService = aiService
	}
}
