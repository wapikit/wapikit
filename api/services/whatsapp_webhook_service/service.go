package webhook_service

import (
	"encoding/json"
	"fmt"
	"net/http"

	wapi "github.com/wapikit/wapi.go/pkg/client"
	"github.com/wapikit/wapi.go/pkg/events"
	"github.com/wapikit/wapikit/api/services"
	"github.com/wapikit/wapikit/internal/api_types"
	"github.com/wapikit/wapikit/internal/core/api_server_events"
	"github.com/wapikit/wapikit/internal/interfaces"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
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
	fmt.Println("get request received", context.QueryParams())
	getHandler := service.wapiClient.GetWebhookGetRequestHandler()
	getHandler(context)
	return nil
}

func fetchConversation(businessAccountId, phoneNumberId string, app interfaces.App) (*api_server_events.ConversationWithAllDetails, error) {
	var dest *api_server_events.ConversationWithAllDetails

	conversationQuery := SELECT(
		table.Conversation.AllColumns,
		table.WhatsappBusinessAccount.AllColumns,
		table.Organization.AllColumns,
		table.Contact.AllColumns,
		table.ConversationAssignment.AllColumns,
		table.OrganizationMember.AllColumns,
		table.User.AllColumns,
	).FROM(
		table.Conversation.
			LEFT_JOIN(table.Organization, table.Organization.UniqueId.EQ(table.Conversation.OrganizationId)).
			LEFT_JOIN(table.WhatsappBusinessAccount, table.WhatsappBusinessAccount.OrganizationId.EQ(table.Organization.UniqueId)).
			LEFT_JOIN(table.Contact, table.Contact.UniqueId.EQ(table.Conversation.ContactId)).
			LEFT_JOIN(table.ConversationAssignment, table.ConversationAssignment.ConversationId.EQ(table.Conversation.UniqueId)).
			LEFT_JOIN(table.OrganizationMember, table.OrganizationMember.UniqueId.EQ(table.ConversationAssignment.AssignedToOrganizationMemberId)).
			LEFT_JOIN(table.User, table.User.UniqueId.EQ(table.OrganizationMember.UserId)),
	).WHERE(
		table.WhatsappBusinessAccount.AccountId.EQ(String(businessAccountId)).
			AND(
				table.Contact.PhoneNumber.EQ(String(phoneNumberId)),
			),
	).LIMIT(1)

	err := conversationQuery.Query(app.Db, dest)

	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (service *WebhookService) handleWebhookPostRequest(context interfaces.ContextWithoutSession) error {
	for eventType, handler := range service.handlerMap {
		//  ! TODO: a middleware here which parses the required event handler parameter and type cast it to the corresponding type
		service.wapiClient.On(eventType, func(event events.BaseEvent) {
			handler(event, context.App)
		})
	}
	postHandler := service.wapiClient.GetWebhookPostRequestHandler()
	postHandler(context)
	return nil
}

func handleTextMessage(event events.BaseEvent, app interfaces.App) {
	textMessageEvent := event.(*events.TextMessageEvent)

	businessAccountId := textMessageEvent.BusinessAccountId

	// ! TODO: this phone number Id is the phone number id of business to which user has sent the message on whatsapp
	phoneNumber := textMessageEvent.PhoneNumber

	// sentByContactWhatsappAccountId := textMessageEvent.BaseMessageEvent.Context.From

	// ! TODO: need to update wapi.go to return the contact number of user to whom this message has been sent by

	conversation, err := fetchConversation(businessAccountId, phoneNumber, app)
	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			app.Logger.Error("organization not found", err.Error(), nil)
		}
		app.Logger.Error("error fetching organization details", err.Error(), nil)
	}

	messageData := map[string]interface{}{
		"text": textMessageEvent.Text,
	}

	jsonMessageData, _ := json.Marshal(messageData)
	stringMessageData := string(jsonMessageData)

	messageToInsert := model.Message{
		WhatsAppMessageId:         &textMessageEvent.MessageId,
		WhatsappBusinessAccountId: &businessAccountId,
		ConversationId:            &conversation.UniqueId,
		CampaignId:                conversation.InitiatedByCampaignId,
		ContactId:                 conversation.ContactId,
		MessageType:               model.MessageTypeEnum_Text,
		Status:                    model.MessageStatusEnum_Sent,
		Direction:                 model.MessageDirectionEnum_InBound,
		MessageData:               &stringMessageData,
		OrganizationId:            conversation.OrganizationId,
	}

	// * insert this message in DB and get the unique id, then send it to the websocket server, so it can broadcast it to the frontend
	insertQuery := table.Message.INSERT().
		MODEL(messageToInsert).
		RETURNING(table.Message.UniqueId)

	err = insertQuery.Query(app.Db, &messageToInsert)

	if err != nil {
		app.Logger.Error("error inserting message in the database", err.Error(), nil)
	}

	message := api_types.MessageSchema{
		ConversationId: conversation.UniqueId.String(),
		Direction:      api_types.InBound,
		MessageType:    api_types.Text,
		Status:         api_types.MessageStatusEnumSent,
		MessageData:    &messageData,
		UniqueId:       messageToInsert.UniqueId.String(),
	}

	apiServerEvent := api_server_events.NewMessageEvent{
		BaseApiServerEvent: api_server_events.BaseApiServerEvent{
			EventType:    api_server_events.ApiServerNewMessageEvent,
			Conversation: *conversation,
		},
		EventType: api_server_events.ApiServerNewMessageEvent,
		Message:   message,
	}

	fmt.Println("apiServerEvent is", string(apiServerEvent.ToJson()))
	err = app.Redis.PublishMessageToRedisChannel(app.Constants.RedisEventChannelName, apiServerEvent.ToJson())

	if err != nil {
		fmt.Println("error sending api server event", err)
	}

	// ! TODO: quick actions, AI automation replies and other stuff will be added in the future version here
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
