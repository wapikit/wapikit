//go:build community_edition
// +build community_edition

package main

import (
	"github.com/wapikit/wapikit/api/api_types"
	"github.com/wapikit/wapikit/interfaces"
	ai_service "github.com/wapikit/wapikit/services/ai_service"
)

func MountServices(app *interfaces.App) {
	logger := app.Logger
	redis := app.Redis
	db := app.Db
	logger.Info("Mounting AI service")
	aiService := ai_service.NewAiService(
		&logger,
		redis,
		db,
		koa.String("ai.api_key"),
		api_types.Gpt4o,
	)

	app.AiService = aiService

}
