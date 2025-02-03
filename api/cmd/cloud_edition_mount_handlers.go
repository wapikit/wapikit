//go:build managed_cloud
// +build managed_cloud

package api

import (
	"github.com/labstack/echo/v4"
	"github.com/wapikit/wapikit-enterprise/api/controllers/ai_controller"
	"github.com/wapikit/wapikit-enterprise/api/controllers/payment_controller"
	"github.com/wapikit/wapikit-enterprise/api/controllers/subscription_controller"
	"github.com/wapikit/wapikit-enterprise/api/controllers/system_controller"
	"github.com/wapikit/wapikit/api/controllers/analytics_controller"
	"github.com/wapikit/wapikit/api/controllers/auth_controller"
	"github.com/wapikit/wapikit/api/controllers/campaign_controller"
	"github.com/wapikit/wapikit/api/controllers/contact_controller"
	"github.com/wapikit/wapikit/api/controllers/contact_list_controller"
	"github.com/wapikit/wapikit/api/controllers/conversation_controller"
	"github.com/wapikit/wapikit/api/controllers/event_controller"
	"github.com/wapikit/wapikit/api/controllers/integration_controller"
	"github.com/wapikit/wapikit/api/controllers/organization_controller"
	"github.com/wapikit/wapikit/api/controllers/rbac_controller"
	"github.com/wapikit/wapikit/api/controllers/user_controller"
	"github.com/wapikit/wapikit/api/controllers/webhook_controller"
	"github.com/wapikit/wapikit/interfaces"
)

// registerHandlers registers HTTP handlers.
func mountHandlerServices(e *echo.Echo, app *interfaces.App) {
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
	paymentController := payment_controller.NewPaymentController()
	subscriptionController := subscription_controller.NewSubscriptionController()
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
		paymentController,
		subscriptionController,
		eventController,
	)

	for _, service := range controllersToRegister {
		service.Register(e)
	}
}
