package webhook_service

import (
	"github.com/sarthakjdev/wapikit/internal/interfaces"
	"github.com/sarthakjdev/wapikit/services"
)

type WebhookService struct {
	services.BaseService `json:"-,inline"`
}

func NewWebhookService() *WebhookService {

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
