package webhook_controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	wapi "github.com/wapikit/wapi.go/pkg/client"
	"github.com/wapikit/wapi.go/pkg/events"
	"github.com/wapikit/wapikit/api/api_types"
	controller "github.com/wapikit/wapikit/api/controllers"
	"github.com/wapikit/wapikit/interfaces"
	"github.com/wapikit/wapikit/internal/api_server_events"
	"github.com/wapikit/wapikit/utils"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
)

type WebhookController struct {
	controller.BaseController `json:"-,inline"`
	handlerMap                map[events.EventType]func(events.BaseEvent, interfaces.App)
}

func NewWhatsappWebhookWebhookController(wapiClient *wapi.Client) *WebhookController {
	service := &WebhookController{
		BaseController: controller.BaseController{
			Name:        "Webhook Controller",
			RestApiPath: "/api/webhook",
			Routes:      []interfaces.Route{},
		},
	}

	service.BaseController.Routes = []interfaces.Route{
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

func (service *WebhookController) handleWebhookGetRequest(context interfaces.ContextWithoutSession) error {

	decrypter := context.App.EncryptionService
	logger := context.App.Logger
	webhookVerificationToken := context.QueryParam("hub.verify_token")
	logger.Info("webhook verification token", webhookVerificationToken, nil)

	var decryptedDetails utils.WebhookSecretData

	err := decrypter.DecryptData(webhookVerificationToken, &decryptedDetails)
	logger.Info("decrypted details", decryptedDetails, nil)
	if err != nil {
		logger.Error("error decrypting webhook verification token", err.Error(), nil)
		return context.JSON(http.StatusBadRequest, "Invalid verification token")
	}

	if &decryptedDetails == nil {
		logger.Error("decrypted details are nil", "", nil)
		return context.JSON(http.StatusBadRequest, "Invalid verification token")
	}

	// ! FETCH THE BUSINESS ACCOUNT DETAILS FROM THE DATABASE
	orgUuid, err := uuid.Parse(decryptedDetails.OrganizationId)

	if err != nil {
		logger.Error("error parsing organization id", err.Error(), nil)
		return context.JSON(http.StatusBadRequest, "Invalid organization id")
	}

	whatsappBusinessAccountDetails := SELECT(
		table.WhatsappBusinessAccount.AllColumns,
	).FROM(
		table.WhatsappBusinessAccount,
	).WHERE(
		table.WhatsappBusinessAccount.AccountId.EQ(String(decryptedDetails.WhatsappBusinessAccountId)).AND(
			table.WhatsappBusinessAccount.OrganizationId.EQ(UUID(orgUuid)),
		),
	).LIMIT(1)

	var businessAccount model.WhatsappBusinessAccount

	err = whatsappBusinessAccountDetails.Query(context.App.Db, &businessAccount)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			logger.Error("business account not found", err.Error(), nil)
			return context.JSON(http.StatusNotFound, "Business account not found")
		}

		logger.Error("error fetching business account details", err.Error(), nil)

		return context.JSON(http.StatusInternalServerError, "Internal server error")
	}

	wapiClient := wapi.New(&wapi.ClientConfig{
		BusinessAccountId: businessAccount.AccountId,
		ApiAccessToken:    businessAccount.AccessToken,
		WebhookSecret:     businessAccount.WebhookSecret,
	})

	getHandler := wapiClient.GetWebhookGetRequestHandler()
	getHandler(context)
	return nil
}

type ContactWithAllDetails struct {
	model.Contact
	Organization model.Organization
}

type BusinessAccountDetails struct {
	model.WhatsappBusinessAccount
	Organization model.Organization `json:"organization"`
}

func fetchBusinessAccountDetails(businessAccountId string, app interfaces.App) (*BusinessAccountDetails, error) {
	var businessAccountDetails BusinessAccountDetails

	businessAccountQuery := SELECT(
		table.WhatsappBusinessAccount.AllColumns,
		table.Organization.AllColumns,
	).FROM(
		table.WhatsappBusinessAccount.
			LEFT_JOIN(table.Organization, table.Organization.UniqueId.EQ(table.WhatsappBusinessAccount.OrganizationId)),
	).WHERE(
		table.WhatsappBusinessAccount.AccountId.EQ(String(businessAccountId)),
	).LIMIT(1)

	err := businessAccountQuery.Query(app.Db, &businessAccountDetails)

	if err != nil {
		return nil, err
	}

	// ! TODO: add caching here for 30 minutes

	return &businessAccountDetails, nil

}

func fetchContact(sentByContactNumber, businessAccountId string, app interfaces.App) (*ContactWithAllDetails, error) {
	var contact ContactWithAllDetails

	contactQuery := SELECT(
		table.Contact.AllColumns,
		table.Organization.AllColumns,
		table.WhatsappBusinessAccount.AllColumns).
		FROM(table.Contact.
			LEFT_JOIN(table.Organization, table.Organization.UniqueId.EQ(table.Contact.OrganizationId)).
			LEFT_JOIN(table.WhatsappBusinessAccount, table.WhatsappBusinessAccount.OrganizationId.EQ(table.Organization.UniqueId).
				AND(table.WhatsappBusinessAccount.AccountId.EQ(String(businessAccountId)))),
		).
		WHERE(
			table.Contact.PhoneNumber.EQ(String(sentByContactNumber)).AND(
				table.WhatsappBusinessAccount.AccountId.EQ(String(businessAccountId)),
			),
		).LIMIT(1)

	err := contactQuery.Query(app.Db, &contact)

	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func fetchConversation(businessAccountId, sentByContactNumber string, app interfaces.App) (*api_server_events.ConversationWithAllDetails, error) {
	var dest api_server_events.ConversationWithAllDetails

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
		table.Conversation.Status.EQ(utils.EnumExpression(model.ConversationStatusEnum_Active.String())).
			AND(table.WhatsappBusinessAccount.AccountId.EQ(String(businessAccountId))).
			AND(table.Contact.PhoneNumber.EQ(String(sentByContactNumber))),
	).LIMIT(1)

	err := conversationQuery.Query(app.Db, &dest)

	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func (service *WebhookController) handleWebhookPostRequest(context interfaces.ContextWithoutSession) error {

	// ! GET THE BUSINESS ACCOUNT ID HERE

	logger := context.App.Logger

	// * Read the request body so we can parse out the businessAccountId.
	bodyBytes, err := io.ReadAll(context.Request().Body)
	if err != nil {
		logger.Error("Error reading request body: %v", err.Error(), nil)
		return context.JSON(http.StatusInternalServerError, "Error reading request body")
	}

	// * Parse JSON to find the businessAccountId.
	var raw map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &raw); err != nil {
		logger.Error("Error unmarshaling JSON: %v", err.Error(), nil)
		return context.JSON(http.StatusBadRequest, "Invalid JSON")
	}

	var businessAccountId string
	if entryList, ok := raw["entry"].([]interface{}); ok && len(entryList) > 0 {
		if firstEntry, ok := entryList[0].(map[string]interface{}); ok {
			if id, ok := firstEntry["id"].(string); ok {
				businessAccountId = id
			}
		}
	}

	whatsappBusinessAccountDetails := SELECT(
		table.WhatsappBusinessAccount.AllColumns,
	).FROM(
		table.WhatsappBusinessAccount,
	).WHERE(
		table.WhatsappBusinessAccount.AccountId.EQ(String(businessAccountId)),
	).
		LIMIT(1)

	var businessAccount model.WhatsappBusinessAccount
	err = whatsappBusinessAccountDetails.Query(context.App.Db, &businessAccount)
	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "Business account not found")
		}
		return context.JSON(http.StatusInternalServerError, "Internal server error")
	}

	// 3) Reset the body so the wapiClient can read it again.
	context.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 4) Create the wapiClient with the discovered businessAccountId
	wapiClient := wapi.New(&wapi.ClientConfig{
		BusinessAccountId: businessAccountId,
		ApiAccessToken:    businessAccount.AccessToken,
		WebhookSecret:     businessAccount.WebhookSecret,
	})

	for eventType, handler := range service.handlerMap {
		//  ! TODO: a middleware here which parses the required event handler parameter and type cast it to the corresponding type
		wapiClient.On(eventType, func(event events.BaseEvent) {
			handler(event, context.App)
		})
	}

	postHandler := wapiClient.GetWebhookPostRequestHandler()
	err = postHandler(context)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return context.JSON(http.StatusOK, "Success")
}

func preHandlerHook(app interfaces.App, businessAccountId string, phoneNumber events.BusinessPhoneNumber, sentByContactNumber string) (*api_server_events.ConversationWithAllDetails, error) {
	conversationDetailsToReturn := &api_server_events.ConversationWithAllDetails{}
	businessAccount, err := fetchBusinessAccountDetails(businessAccountId, app)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			// ! it must be in rare case, because th webhook should not be get to the application, if somebody has a account or running instance and added the webhook details in the whatsapp business account, then it should be in the database
			app.Logger.Error("business account not found", err.Error(), nil)
		}

		app.Logger.Error("error fetching business account details", err.Error(), nil)
		// ! TODO: send notification to the team
		return nil, fmt.Errorf("error fetching business account details")

	} else {
		// * business account found, add it to the response object
		conversationDetailsToReturn.WhatsappBusinessAccount = *businessAccount
	}

	contact, err := fetchContact(sentByContactNumber, businessAccountId, app)

	var contactId uuid.UUID

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			app.Logger.Info("No contact found adding this person to the contacts", err.Error(), nil)

			emptyAttributes := "{}"

			contactToAdd := model.Contact{
				PhoneNumber:    sentByContactNumber,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
				OrganizationId: businessAccount.OrganizationId,
				Status:         model.ContactStatusEnum_Active,
				// ! TODO: add this to wapi.go and then use here
				Name:       "",
				Attributes: &emptyAttributes,
			}

			var insertedContact model.Contact

			insertQuery := table.Contact.INSERT(table.Contact.MutableColumns).
				MODEL(contactToAdd).
				RETURNING(table.Contact.AllColumns)

			err = insertQuery.Query(app.Db, &insertedContact)

			if err != nil {
				return nil, fmt.Errorf("error inserting contact in the database")
			}

			contactId = insertedContact.UniqueId
			conversationDetailsToReturn.Contact = insertedContact
		} else {
			// ! TODO: send notification to the team
			return nil, fmt.Errorf("error inserting contact in the database")
		}
	} else {
		// * contact found, add it to the response object
		conversationDetailsToReturn.Contact = model.Contact{
			UniqueId:       contact.UniqueId,
			PhoneNumber:    contact.PhoneNumber,
			Name:           contact.Name,
			OrganizationId: contact.OrganizationId,
			Status:         contact.Status,
			CreatedAt:      contact.CreatedAt,
			UpdatedAt:      contact.UpdatedAt,
			Attributes:     contact.Attributes,
		}

		contactId = contact.UniqueId
	}

	fetchedConversation, err := fetchConversation(businessAccountId, sentByContactNumber, app)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			// * this is a new message from the user, so we need to create a new conversation
			conversationToInsert := model.Conversation{
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				ContactId:       contactId,
				OrganizationId:  businessAccount.OrganizationId,
				PhoneNumberUsed: phoneNumber.Id,
				InitiatedBy:     model.ConversationInitiatedEnum_Contact,
				Status:          model.ConversationStatusEnum_Active,
			}

			var insertedConversation model.Conversation

			insertQuery := table.Conversation.INSERT(table.Conversation.MutableColumns).
				MODEL(conversationToInsert).
				RETURNING(table.Conversation.AllColumns)

			err = insertQuery.Query(app.Db, &insertedConversation)

			if err != nil {
				return nil, fmt.Errorf("error inserting conversation in the database")
			}

			conversationDetailsToReturn.Conversation = model.Conversation{
				UniqueId:              insertedConversation.UniqueId,
				CreatedAt:             insertedConversation.CreatedAt,
				UpdatedAt:             insertedConversation.UpdatedAt,
				ContactId:             insertedConversation.ContactId,
				OrganizationId:        insertedConversation.OrganizationId,
				Status:                insertedConversation.Status,
				PhoneNumberUsed:       insertedConversation.PhoneNumberUsed,
				InitiatedBy:           insertedConversation.InitiatedBy,
				InitiatedByCampaignId: insertedConversation.InitiatedByCampaignId,
			}

		} else {
			return nil, fmt.Errorf("error fetching conversation from the database")
		}
	} else {
		// * conversation found, add it to the response object
		conversationDetailsToReturn.Conversation = model.Conversation{
			UniqueId:              fetchedConversation.UniqueId,
			CreatedAt:             fetchedConversation.CreatedAt,
			UpdatedAt:             fetchedConversation.UpdatedAt,
			ContactId:             fetchedConversation.ContactId,
			OrganizationId:        fetchedConversation.OrganizationId,
			Status:                fetchedConversation.Status,
			PhoneNumberUsed:       fetchedConversation.PhoneNumberUsed,
			InitiatedBy:           fetchedConversation.InitiatedBy,
			InitiatedByCampaignId: fetchedConversation.InitiatedByCampaignId,
		}

		// ! TODO: handle other properties like assigned to etc etc
	}

	return conversationDetailsToReturn, nil
}

func handleTextMessage(event events.BaseEvent, app interfaces.App) {
	textMessageEvent := event.(*events.TextMessageEvent)
	businessAccountId := textMessageEvent.BusinessAccountId
	phoneNumber := textMessageEvent.PhoneNumber
	sentAt := textMessageEvent.BaseMessageEvent.Timestamp // this is unix timestamp in string, convert this to time.Time

	unixTimestamp, err := strconv.ParseInt(sentAt, 10, 64)
	if err != nil {
		app.Logger.Error("error parsing timestamp", err.Error(), nil)
		return
	}

	sentAtTime := time.Unix(unixTimestamp, 0)
	sentByContactNumber := textMessageEvent.BaseMessageEvent.From

	app.Logger.Debug("details", "businessAccountId", businessAccountId, "phoneNumber", phoneNumber, "sentByContactNumber", sentByContactNumber)

	conversationDetails, err := preHandlerHook(app, businessAccountId, phoneNumber, sentByContactNumber)

	if err != nil {
		app.Logger.Error("error fetching conversation details", err.Error(), nil)
		return
	}

	conversationDetailsString, _ := json.Marshal(conversationDetails)
	app.Logger.Info("conversation details", string(conversationDetailsString), nil)

	messageData := map[string]interface{}{
		"text": textMessageEvent.Text,
	}

	jsonMessageData, _ := json.Marshal(messageData)
	stringMessageData := string(jsonMessageData)

	var insertedMessage model.Message

	messageToInsert := model.Message{
		WhatsAppMessageId:         &textMessageEvent.MessageId,
		WhatsappBusinessAccountId: &businessAccountId,
		ConversationId:            &conversationDetails.UniqueId,
		CampaignId:                conversationDetails.InitiatedByCampaignId,
		ContactId:                 conversationDetails.ContactId,
		MessageType:               model.MessageTypeEnum_Text,
		Status:                    model.MessageStatusEnum_Sent,
		Direction:                 model.MessageDirectionEnum_InBound,
		MessageData:               &stringMessageData,
		OrganizationId:            conversationDetails.OrganizationId,
		CreatedAt:                 sentAtTime,
		UpdatedAt:                 time.Now(),
	}

	// * insert this message in DB and get the unique id, then send it to the websocket server, so it can broadcast it to the frontend
	insertQuery := table.Message.
		INSERT(table.Message.MutableColumns).
		MODEL(messageToInsert).
		RETURNING(table.Message.UniqueId)

	err = insertQuery.Query(app.Db, &insertedMessage)

	if err != nil {
		app.Logger.Error("error inserting message in the database", err.Error(), nil)
	}

	message := api_types.MessageSchema{
		ConversationId: conversationDetails.UniqueId.String(),
		Direction:      api_types.InBound,
		MessageType:    api_types.Text,
		Status:         api_types.MessageStatusEnumSent,
		MessageData:    &messageData,
		UniqueId:       insertedMessage.UniqueId.String(),
		CreatedAt:      sentAtTime,
	}

	apiServerEvent := api_server_events.NewMessageEvent{
		BaseApiServerEvent: api_server_events.BaseApiServerEvent{
			EventType:    api_server_events.ApiServerNewMessageEvent,
			Conversation: *conversationDetails,
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
	videoMessageEvent := event.(*events.VideoMessageEvent)
	businessAccountId := videoMessageEvent.BusinessAccountId
	phoneNumber := videoMessageEvent.PhoneNumber
	sentByContactNumber := videoMessageEvent.BaseMessageEvent.From

	conversationDetails, err := preHandlerHook(app, businessAccountId, phoneNumber, sentByContactNumber)

	if err != nil {
		app.Logger.Error("error fetching conversation details", err.Error(), nil)
		return
	}

	conversationDetailsString, _ := json.Marshal(conversationDetails)
	app.Logger.Info("conversation details", string(conversationDetailsString), nil)

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
