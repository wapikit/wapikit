//go:build community_edition
// +build community_edition

package api

import (
	"github.com/labstack/echo/v4"
	"github.com/wapikit/wapikit/api/controllers/ai_controller"
	"github.com/wapikit/wapikit/api/controllers/analytics_controller"
	"github.com/wapikit/wapikit/api/controllers/auth_controller"
	"github.com/wapikit/wapikit/api/controllers/campaign_controller"
	"github.com/wapikit/wapikit/api/controllers/contact_controller"
	"github.com/wapikit/wapikit/api/controllers/contact_list_controller"
	"github.com/wapikit/wapikit/api/controllers/conversation_controller"
	"github.com/wapikit/wapikit/api/controllers/event_controller"
	"github.com/wapikit/wapikit/api/controllers/integration_controller"
	"github.com/wapikit/wapikit/api/controllers/next_files_controller"
	"github.com/wapikit/wapikit/api/controllers/organization_controller"
	"github.com/wapikit/wapikit/api/controllers/rbac_controller"
	"github.com/wapikit/wapikit/api/controllers/system_controller"
	"github.com/wapikit/wapikit/api/controllers/user_controller"
	"github.com/wapikit/wapikit/api/controllers/webhook_controller"
	"github.com/wapikit/wapikit/interfaces"
)

// registerHandlers registers HTTP handlers.
func mountHandlerServices(e *echo.Echo, app *interfaces.App) {
	logger := app.Logger
	koa := app.Koa

	isFrontendHostedSeparately := koa.Bool("is_frontend_separately_hosted")

	controllersToRegister := []interfaces.ApiController{}
	userController := user_controller.NewUserController()
	authController := auth_controller.NewAuthController()
	organizationController := organization_controller.NewOrganizationController()
	campaignController := campaign_controller.NewCampaignController()
	analyticsController := analytics_controller.NewAnalyticsController()
	contactsController := contact_controller.NewContactController()
	conversationController := conversation_controller.NewConversationController()
	contactListController := contact_list_controller.NewContactListController()
	systemController := system_controller.NewSystemController()
	integrationController := integration_controller.NewIntegrationController()
	roleBasedAccessControlController := rbac_controller.NewRoleBasedAccessControlController()
	whatsappWebhookController := webhook_controller.NewWhatsappWebhookWebhookController(app.WapiClient)
	aiController := ai_controller.NewAiController()
	eventController := event_controller.NewEventController()

	controllersToRegister = append(
		controllersToRegister,
		userController,
		authController,
		campaignController,
		contactListController,
		contactsController,
		conversationController,
		systemController,
		analyticsController,
		organizationController,
		integrationController,
		roleBasedAccessControlController,
		whatsappWebhookController,
		aiController,
		eventController,
	)

	if !isFrontendHostedSeparately {
		logger.Info("Frontend is not hosted separately")
		nextFileServerService := next_files_controller.NewNextFileServerController()
		controllersToRegister = append(controllersToRegister, nextFileServerService)
	}

	for _, service := range controllersToRegister {
		service.Register(e)
	}
}
