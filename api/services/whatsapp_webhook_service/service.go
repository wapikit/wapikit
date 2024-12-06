package webhook_service

import (
	"fmt"
	"net/http"

	wapi "github.com/wapikit/wapi.go/pkg/client"
	"github.com/wapikit/wapi.go/pkg/events"
	"github.com/wapikit/wapikit/api/services"
	"github.com/wapikit/wapikit/internal/interfaces"
)

type WebhookService struct {
	services.BaseService `json:"-,inline"`
	wapiClient           *wapi.Client
	handlerMap           map[events.EventType]func(events.BaseEvent, interfaces.App)
}

func NewWhatsappWebhookServiceWebhookService(wapiClient *wapi.Client) *WebhookService {
	service := &WebhookService{
		BaseService: services.BaseService{
			Name:        "Webhook Service",
			RestApiPath: "/api/webhook",
			Routes:      []interfaces.Route{},
		},
		wapiClient: wapiClient,
	}

	service.BaseService.Routes = []interfaces.Route{
		{
			Path:                    "/api/webhook/whatsapp",
			Method:                  http.MethodGet,
			Handler:                 interfaces.HandlerWithoutSession(service.handleWebhookGetRequest), // Using service method here
			IsAuthorizationRequired: false,
		},
		{
			Path:                    "/api/webhook/whatsapp",
			Method:                  http.MethodPost,
			Handler:                 interfaces.HandlerWithoutSession(service.handleWebhookPostRequest), // Using service method here
			IsAuthorizationRequired: false,
		},
	}

	service.handlerMap = map[events.EventType]func(event events.BaseEvent, app interfaces.App){
		events.TextMessageEventType:                  handleTextMessage,
		events.VideoMessageEventType:                 handleVideoMessageEvent,
		events.ImageMessageEventType:                 handleImageMessageEvent,
		events.AccountAlertsEventType:                handleAccountAlerts,
		events.DocumentMessageEventType:              handleDocumentMessageEvent,
		events.AudioMessageEventType:                 handleAudioMessageEvent,
		events.MessageReadEventType:                  handleMessageReadEvent,
		events.SecurityEventType:                     handleSecurityEvent,
		events.ErrorEventType:                        handleErrorEvent,
		events.AdInteractionEventType:                handleAdInteractionEvent,
		events.CustomerNumberChangedEventType:        handlePhoneNumberChangeEvent,
		events.AccountReviewUpdateEventType:          handleAccountReviewUpdateEvent,
		events.AccountUpdateEventType:                handleAccountUpdateEvent,
		events.TemplateMessageEventType:              handleTemplateMessageEvent,
		events.ContactMessageEventType:               handleContactMessageEvent,
		events.ListInteractionMessageEventType:       handleListInteractionMessageEvent,
		events.LocationMessageEventType:              handleLocationMessageEvent,
		events.MessageDeliveredEventType:             handleMessageDeliveredEvent,
		events.MessageFailedEventType:                handleMessageFailedEvent,
		events.QuickReplyMessageEventType:            handleQuickReplyMessageEvent,
		events.ReplyButtonInteractionEventType:       handleReplyButtonInteractionEvent,
		events.ReactionMessageEventType:              handleReactionMessageEvent,
		events.BusinessCapabilityUpdateEventType:     handleBusinessCapabilityUpdateEvent,
		events.ProductInquiryEventType:               handleProductInquiryEvent,
		events.OrderReceivedEventType:                handleOrderReceivedEvent,
		events.StickerMessageEventType:               handleStickerMessageEvent,
		events.MessageUndeliveredEventType:           handleMessageUndeliveredEvent,
		events.CustomerIdentityChangedEventType:      handleCustomerIdentityChangedEvent,
		events.MessageSentEventType:                  handleMessageSentEvent,
		events.UnknownEventType:                      handleUnknownEvent,
		events.WarnEventType:                         handleWarnEvent,
		events.ReadyEventType:                        handleReadyEvent,
		events.MessageTemplateStatusUpdateEventType:  handleMessageTemplateUpdateEvent,
		events.MessageTemplateQualityUpdateEventType: handleMessageTemplateQualityUpdateEvent,
		events.PhoneNumberNameUpdateEventType:        handlePhoneNumberNameUpdateEvent,
		events.PhoneNumberQualityUpdateEventType:     handlePhoneNumberQualityUpdateEvent,
	}

	return service

}

func (service *WebhookService) handleWebhookGetRequest(context interfaces.ContextWithoutSession) error {
	// ! TODO: here we need to call the wapiClient get request handler
	return nil
}

func (service *WebhookService) handleWebhookPostRequest(context interfaces.ContextWithoutSession) error {
	for eventType, handler := range service.handlerMap {
		//  ! TODO: a middleware here which parses the required event handler parameter and type cast it to the corresponding type
		service.wapiClient.On(eventType, func(event events.BaseEvent) {
			handler(event, context.App)
		})
	}
	// app := context.App
	// app.SendApiServerEvent()
	return nil
}

func handleTextMessage(event events.BaseEvent, app interfaces.App) {
	// Handle text message event
	// ! TODO: type cast this base event to textMessage event
	fmt.Println(event)
	// ! TODO:
	// ! update the db
	// ! send an api_server_event to websocket server which will in return check if the user with the id exists if exists then it broadcasts message to frontend
	// ! check for quick action, now feature flag must be checked here
	// ! if quick action keywords are enabled then send a quick reply
}

func handleVideoMessageEvent(event events.BaseEvent, app interfaces.App) {
}

func handleImageMessageEvent(event events.BaseEvent, app interfaces.App) {

}

func handleDocumentMessageEvent(event events.BaseEvent, app interfaces.App) {

}

func handleAudioMessageEvent(event events.BaseEvent, app interfaces.App) {

}

func handleMessageReadEvent(event events.BaseEvent, app interfaces.App) {

	// ! TODO: mark the message as read in the database

	// ! send an api_server_event to webhook

}

func handlePhoneNumberChangeEvent(event events.BaseEvent, app interfaces.App) {

	// ! check for the contact in the database

	// ! change the phone number
	// send an api_server_event to webhook

}

func handleSecurityEvent(event events.BaseEvent, app interfaces.App) {
	// send an api_server_event to webhook

}

func handleAccountAlerts(event events.BaseEvent, app interfaces.App) {
	// send an api_server_event to webhook

}

func handleAdInteractionEvent(event events.BaseEvent, app interfaces.App) {
	// send an api_server_event to webhook

}

func handleErrorEvent(event events.BaseEvent, app interfaces.App) {
	// send an api_server_event to webhook

}

func handleAccountReviewUpdateEvent(event events.BaseEvent, app interfaces.App) {

}

func handleAccountUpdateEvent(event events.BaseEvent, app interfaces.App) {

}

func handleTemplateMessageEvent(event events.BaseEvent, app interfaces.App) {

}

func handleContactMessageEvent(event events.BaseEvent, app interfaces.App) {

}

func handleListInteractionMessageEvent(event events.BaseEvent, app interfaces.App) {

}

func handleLocationMessageEvent(event events.BaseEvent, app interfaces.App) {

}

func handleMessageDeliveredEvent(event events.BaseEvent, app interfaces.App) {

}

func handleMessageFailedEvent(event events.BaseEvent, app interfaces.App) {

}

func handleQuickReplyMessageEvent(event events.BaseEvent, app interfaces.App) {

}

func handleReplyButtonInteractionEvent(event events.BaseEvent, app interfaces.App) {

}

func handleReactionMessageEvent(event events.BaseEvent, app interfaces.App) {

}

func handleBusinessCapabilityUpdateEvent(event events.BaseEvent, app interfaces.App) {

}

func handleProductInquiryEvent(event events.BaseEvent, app interfaces.App) {

}

func handleOrderReceivedEvent(event events.BaseEvent, app interfaces.App) {

}

func handleStickerMessageEvent(event events.BaseEvent, app interfaces.App) {

}

func handleMessageUndeliveredEvent(event events.BaseEvent, app interfaces.App) {

}

func handleCustomerIdentityChangedEvent(event events.BaseEvent, app interfaces.App) {
	// ! TODO:
	// ! 1. check if the customer exists in the database
	// ! 2. update the customer identity
	// ! 3. send an api_server_event to websocket server to logout the user if connected, else send a notification to the user to login again
}

func handleMessageSentEvent(event events.BaseEvent, app interfaces.App) {

}

func handleUnknownEvent(event events.BaseEvent, app interfaces.App) {
	// ! TODO: in this handle we need to save the log in the database
}

func handleWarnEvent(event events.BaseEvent, app interfaces.App) {
}

func handleReadyEvent(event events.BaseEvent, app interfaces.App) {
}

func handleMessageTemplateUpdateEvent(event events.BaseEvent, app interfaces.App) {
}

func handleMessageTemplateQualityUpdateEvent(event events.BaseEvent, app interfaces.App) {
}

func handlePhoneNumberNameUpdateEvent(event events.BaseEvent, app interfaces.App) {
}

func handlePhoneNumberQualityUpdateEvent(event events.BaseEvent, app interfaces.App) {
}
