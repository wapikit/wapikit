package webhook_service

import (
	wapi "github.com/sarthakjdev/wapi.go/pkg/client"
	"github.com/sarthakjdev/wapi.go/pkg/events"
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

type WebhookService struct {
	services.BaseService `json:"-,inline"`
	client               wapi.Client
}

func NewWhatsappWebhookServiceWebhookService(wapiClient wapi.Client) *WebhookService {

	// wapiClient := internal.GetWapiCloudClient(
	// 	koa.String("PHONE_NUMBER_ID"),
	// 	koa.String("WHATSAPP_BUSINESS_ACCOUNT_ID"),
	// 	koa.String("WHATSAPP_WEBHOOK_SECRET"),
	// 	koa.String("WHATSAPP_API_ACCESS_TOKEN"),
	// )

	return &WebhookService{
		BaseService: services.BaseService{
			Name:        "Webhook Service",
			RestApiPath: "/api/webhook",
			Routes:      []interfaces.Route{},
		},
		client: wapiClient,
	}
}

func (service *WebhookService) ResgiterEventListeners() {
	service.client.On(events.TextMessageEventType, service.handleTextMessage)
	service.client.On(events.VideoMessageEventType, service.handleVideoMessageEventListeners)
	service.client.On(events.ImageMessageEventType, service.handleImageMessageEventListeners)
	service.client.On(events.DocumentMessageEventType, service.handleDocumentMessageEventListeners)
	service.client.On(events.AudioMessageEventType, service.handleAudioMessageEventListeners)
	service.client.On(events.MessageReadEventType, service.handleMessageReadEventListeners)

}

func (service *WebhookService) handleTextMessage(event events.BaseEvent) {
	// Handle text message event

	// ! TODO:

	// update the db
	// send an api_server_event to webhook
	// check for quick action, now feature flag must be checked here
	// if quick action keywords are enabled then send a quick reply

}

func (service *WebhookService) handleVideoMessageEventListeners(event events.BaseEvent) {

}

func (service *WebhookService) handleImageMessageEventListeners(event events.BaseEvent) {

}

func (service *WebhookService) handleDocumentMessageEventListeners(event events.BaseEvent) {

}

func (service *WebhookService) handleAudioMessageEventListeners(event events.BaseEvent) {

}

func (service *WebhookService) handleMessageReadEventListeners(event events.BaseEvent) {
	// send an api_server_event to webhook

}
