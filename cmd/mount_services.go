//go:build community_edition
// +build community_edition

package main

import (
	"database/sql"

	"github.com/wapikit/wapikit/api/api_types"
	"github.com/wapikit/wapikit/interfaces"
	ai_service "github.com/wapikit/wapikit/services/ai_service"
	cache_service "github.com/wapikit/wapikit/services/redis_service"
)

func MountServices(app interfaces.App, redis *cache_service.RedisClient, db *sql.DB) {
	aiService := ai_service.NewAiService(
		logger,
		redis,
		db,
		koa.String("ai.api_key"),
		api_types.Gpt4o,
	)

	app.AiService = aiService

}
