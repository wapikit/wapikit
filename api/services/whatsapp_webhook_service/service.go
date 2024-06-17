package webhook_service

import (
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

type WebhookService struct {
	services.BaseService `json:"-,inline"`
}

func NewWhatsappWebhookServiceWebhookService() *WebhookService {

	// wapiClient := internal.GetWapiCloudClient(
	// 	koa.String("PHONE_NUMBER_ID"),
	// 	koa.String("WHATSAPP_BUSINESS_ACCOUNT_ID"),
	// 	koa.String("WHATSAPP_WEBHOOK_SECRET"),
	// 	koa.String("WHATSAPP_API_ACCESS_TOKEN"),
	// )

	return &WebhookService{
		BaseService: services.BaseService{
			Name:        "Webhook Service",
			RestApiPath: "/api",
			Routes:      []interfaces.Route{},
		},
	}
}
