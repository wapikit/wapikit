//go:build managed_cloud
// +build managed_cloud

package main

import (
	ai_service "github.com/wapikit/wapikit-enterprise/services/ai"
	"github.com/wapikit/wapikit/api/api_types"
	"github.com/wapikit/wapikit/interfaces"
	"github.com/wapikit/wapikit/services/notification_service"
)

func MountServices(app *interfaces.App) {
	logger := app.Logger
	redis := app.Redis
	db := app.Db

	app.AiService = ai_service.NewAiService(
		&logger,
		redis,
		db,
		koa.String("ai.api_key"),
		api_types.Gpt4o,
	)

	app.NotificationService = &notification_service.NotificationService{
		Logger: &app.Logger,
		SlackConfig: &notification_service.SlackConfig{
			SlackWebhookUrl: koa.String("slack.webhook_url"),
			SlackChannel:    koa.String("slack.channel"),
		},
		EmailConfig: &notification_service.EmailConfig{
			Host:     koa.String("email.host"),
			Port:     koa.String("email.port"),
			Password: koa.String("email.password"),
			Username: koa.String("email.username"),
		},
	}

	app.CampaignManager.NotificationService = app.NotificationService
}
