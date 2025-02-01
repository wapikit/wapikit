package campaign_manager

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/wapikit/wapi.go/manager"
	wapiComponents "github.com/wapikit/wapi.go/pkg/components"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
	"github.com/wapikit/wapikit/services/notification_service"
)

// templateComponentParameters holds the parameters for each component type
type templateComponentParameters struct {
	Header  []string `json:"header"`
	Body    []string `json:"body"`
	Buttons []string `json:"buttons"`
}

// For the purposes of these helper functions, we assume that the fetched template’s
// component has these fields. (Adjust as needed.)
type BusinessTemplateComponent struct {
	Type    string
	Format  string
	Example struct {
		BodyText     [][]string
		HeaderText   []string
		HeaderHandle []string
	}
	Buttons []struct {
		Type    string
		Example []string
	}
}

// buildTemplateMessage creates a new template message and adds all components.
func (cm *CampaignManager) buildTemplateMessage(templateInUse *manager.WhatsAppBusinessMessageTemplateNode, params templateComponentParameters) (*wapiComponents.TemplateMessage, error) {

	cm.Logger.Info("Building template message", nil)
	cm.Logger.Info("Template in use", templateInUse)
	cm.Logger.Info("Template parameters", params)

	templateMessage, err := wapiComponents.NewTemplateMessage(&wapiComponents.TemplateMessageConfigs{
		Name:     templateInUse.Name,
		Language: templateInUse.Language,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating new template message: %v", err)
	}

	// Loop through each component from the fetched template and add to our message.
	for _, comp := range templateInUse.Components {
		cm.Logger.Info("Component", comp)
		switch comp.Type {
		case manager.MessageTemplateComponentTypeBody:
			if err := cm.addBodyComponent(templateMessage, comp, params); err != nil {
				return nil, err
			}
		case manager.MessageTemplateComponentTypeHeader:
			if err := cm.addHeaderComponent(templateMessage, comp, params); err != nil {
				return nil, err
			}
		case manager.MessageTemplateComponentTypeButtons:
			if err := cm.addButtonComponents(templateMessage, comp, params); err != nil {
				return nil, err
			}
		default:
		}
	}

	cm.Logger.Info("Template message", templateMessage)

	return templateMessage, nil
}

// addBodyComponent builds and adds the BODY component.
func (cm *CampaignManager) addBodyComponent(templateMessage *wapiComponents.TemplateMessage, comp manager.WhatsAppBusinessHSMWhatsAppHSMComponent, params templateComponentParameters) error {
	var bodyParameters []wapiComponents.TemplateMessageParameter
	if len(comp.Example.BodyText) > 0 {
		// For each stored body parameter, create a text parameter.
		for _, bodyText := range params.Body {
			bodyParameters = append(bodyParameters, wapiComponents.TemplateMessageBodyAndHeaderParameter{
				Type: wapiComponents.TemplateMessageParameterTypeText,
				Text: &bodyText,
			})
		}
	}
	templateMessage.AddBody(wapiComponents.TemplateMessageComponentBodyType{
		Type:       wapiComponents.TemplateMessageComponentTypeBody,
		Parameters: bodyParameters,
	})
	return nil
}

// addHeaderComponent builds and adds the HEADER component.
// It supports TEXT, IMAGE, VIDEO, DOCUMENT, and LOCATION formats.
func (cm *CampaignManager) addHeaderComponent(templateMessage *wapiComponents.TemplateMessage, comp manager.WhatsAppBusinessHSMWhatsAppHSMComponent, params templateComponentParameters) error {
	// If no header examples exist, simply add an empty header.
	if len(comp.Example.HeaderText) == 0 && len(comp.Example.HeaderHandle) == 0 {
		templateMessage.AddHeader(wapiComponents.TemplateMessageComponentHeaderType{
			Type: wapiComponents.TemplateMessageComponentTypeHeader,
		})
		return nil
	}

	switch comp.Format {
	case "TEXT":
		var headerParameters []wapiComponents.TemplateMessageParameter
		for _, headerText := range params.Header {
			headerParameters = append(headerParameters, wapiComponents.TemplateMessageBodyAndHeaderParameter{
				Type: wapiComponents.TemplateMessageParameterTypeText,
				Text: &headerText,
			})
		}
		templateMessage.AddHeader(wapiComponents.TemplateMessageComponentHeaderType{
			Type:       wapiComponents.TemplateMessageComponentTypeHeader,
			Parameters: &headerParameters,
		})
	case "IMAGE":
		var headerParameters []wapiComponents.TemplateMessageParameter
		for _, mediaUrl := range params.Header {
			headerParameters = append(headerParameters, wapiComponents.TemplateMessageBodyAndHeaderParameter{
				Type: wapiComponents.TemplateMessageParameterTypeImage,
				Image: &wapiComponents.TemplateMessageParameterMedia{
					Link: mediaUrl,
				},
			})
		}
		templateMessage.AddHeader(wapiComponents.TemplateMessageComponentHeaderType{
			Type:       wapiComponents.TemplateMessageComponentTypeHeader,
			Parameters: &headerParameters,
		})
	case "VIDEO":
		var headerParameters []wapiComponents.TemplateMessageParameter
		for _, mediaUrl := range params.Header {
			headerParameters = append(headerParameters, wapiComponents.TemplateMessageBodyAndHeaderParameter{
				Type: wapiComponents.TemplateMessageParameterTypeVideo,
				Video: &wapiComponents.TemplateMessageParameterMedia{
					Link: mediaUrl,
				},
			})
		}
		templateMessage.AddHeader(wapiComponents.TemplateMessageComponentHeaderType{
			Type:       wapiComponents.TemplateMessageComponentTypeHeader,
			Parameters: &headerParameters,
		})
	case "DOCUMENT":
		var headerParameters []wapiComponents.TemplateMessageParameter
		for _, mediaUrl := range params.Header {
			headerParameters = append(headerParameters, wapiComponents.TemplateMessageBodyAndHeaderParameter{
				Type: wapiComponents.TemplateMessageParameterTypeDocument,
				Document: &wapiComponents.TemplateMessageParameterMedia{
					Link: mediaUrl,
				},
			})
		}
		templateMessage.AddHeader(wapiComponents.TemplateMessageComponentHeaderType{
			Type:       wapiComponents.TemplateMessageComponentTypeHeader,
			Parameters: &headerParameters,
		})
	case "LOCATION":
		// For a location header, we expect a JSON string in the first header parameter.
		if len(params.Header) == 0 {
			return fmt.Errorf("no location parameter provided in header")
		}
		var loc wapiComponents.TemplateMessageParameterLocation
		if err := json.Unmarshal([]byte(params.Header[0]), &loc); err != nil {
			return fmt.Errorf("error unmarshalling location header parameter: %v", err)
		}
		headerParameters := []wapiComponents.TemplateMessageParameter{
			wapiComponents.TemplateMessageBodyAndHeaderParameter{
				Type:     wapiComponents.TemplateMessageParameterTypeLocation,
				Location: &loc,
			},
		}
		templateMessage.AddHeader(wapiComponents.TemplateMessageComponentHeaderType{
			Type:       wapiComponents.TemplateMessageComponentTypeHeader,
			Parameters: &headerParameters,
		})
	default:
		return fmt.Errorf("unsupported header format: %s", comp.Format)
	}
	return nil
}

// addButtonComponents builds and adds all BUTTON components.
func (cm *CampaignManager) addButtonComponents(templateMessage *wapiComponents.TemplateMessage, comp manager.WhatsAppBusinessHSMWhatsAppHSMComponent, params templateComponentParameters) error {
	// Loop over each button in the fetched component.
	for index, button := range comp.Buttons {
		cm.Logger.Info("Button type", button.Type)

		switch button.Type {
		case manager.TemplateMessageButtonTypeUrl:
			{
				if len(params.Buttons) > index {
					templateMessage.AddButton(wapiComponents.TemplateMessageComponentButtonType{
						Type:    wapiComponents.TemplateMessageComponentTypeButton,
						SubType: wapiComponents.TemplateMessageButtonComponentTypeUrl,
						Index:   index,
						Parameters: &[]wapiComponents.TemplateMessageParameter{
							wapiComponents.TemplateMessageButtonParameter{
								Type: wapiComponents.TemplateMessageButtonParameterTypeText,
								Text: params.Buttons[index],
							},
						},
					})
				}
			}
		case manager.TemplateMessageButtonTypeQuickReply:
			{

				if len(params.Buttons) > index {
					templateMessage.AddButton(wapiComponents.TemplateMessageComponentButtonType{
						Type:    wapiComponents.TemplateMessageComponentTypeButton,
						SubType: wapiComponents.TemplateMessageButtonComponentTypeQuickReply,
						Index:   index,
						Parameters: &[]wapiComponents.TemplateMessageParameter{
							wapiComponents.TemplateMessageButtonParameter{
								Type:    wapiComponents.TemplateMessageButtonParameterTypePayload,
								Payload: params.Buttons[index],
							},
						},
					})
				} else {
					// templateMessage.AddButton(wapiComponents.TemplateMessageComponentButtonType{
					// 	Type:    wapiComponents.TemplateMessageComponentTypeButton,
					// 	SubType: wapiComponents.TemplateMessageButtonComponentTypeQuickReply,
					// 	Index:   index,
					// })
				}
			}
		// case "COPY_CODE", "copy_code":
		// 	{
		// 		// Implement the copy code button.
		// 		// We assume that wapiComponents now defines:
		// 		//   TemplateMessageButtonComponentTypeCopyCode
		// 		// and that for copy code buttons the parameter is a text string.
		// 		if len(params.Buttons) > index {
		// 			templateMessage.AddButton(wapiComponents.TemplateMessageComponentButtonType{
		// 				Type:    wapiComponents.TemplateMessageComponentTypeButton,
		// 				SubType: wapiComponents.TemplateMessageButtonComponentTypeCopyCode,
		// 				Index:   index,
		// 				Parameters: []wapiComponents.TemplateMessageParameter{
		// 					wapiComponents.TemplateMessageButtonParameter{
		// 						Type: wapiComponents.TemplateMessageButtonParameterTypeText,
		// 						Text: params.Buttons[index],
		// 					},
		// 				},
		// 			})
		// 		} else {
		// 			templateMessage.AddButton(wapiComponents.TemplateMessageComponentButtonType{
		// 				Type:       wapiComponents.TemplateMessageComponentTypeButton,
		// 				SubType:    wapiComponents.TemplateMessageButtonComponentTypeCopyCode,
		// 				Index:      index,
		// 				Parameters: []wapiComponents.TemplateMessageParameter{},
		// 			})
		// 		}
		// 	}
		default:
			// * DO NOTHING *
		}
	}
	return nil
}

// --- Main sendMessage function ---

func (cm *CampaignManager) sendMessage(message *CampaignMessage) error {
	// Ensure that the campaign wait group is decremented and update the last contact ID,
	// irrespective of whether sending succeeds.
	defer func() {
		message.Campaign.wg.Done()
		if err := cm.UpdateLastContactId(message.Campaign.UniqueId, message.Contact.UniqueId); err != nil {
			cm.Logger.Error("error updating last contact id", err.Error())
		}
	}()

	// Retrieve the business worker.
	cm.businessWorkersMutex.RLock()
	worker, exists := cm.businessWorkers[message.Campaign.BusinessAccountId]
	cm.businessWorkersMutex.RUnlock()
	if !exists {
		cm.Logger.Error("Business worker not found", nil)
		cm.NotificationService.SendSlackNotification(notification_service.SlackNotificationParams{
			Title:   "🚨🚨 Business worker not found in send message 🚨🚨",
			Message: "Business worker not found for business account ID: " + message.Campaign.BusinessAccountId,
		})
		return fmt.Errorf("business worker not found for business account ID: %s", message.Campaign.BusinessAccountId)
	}

	// Fetch the template details.
	client := message.Campaign.WapiClient
	templateInUse, err := client.Business.Template.Fetch(*message.Campaign.MessageTemplateId)
	if err != nil {
		message.Campaign.ErrorCount.Add(1)
		return fmt.Errorf("error fetching template: %v", err)
	}

	// Determine if the template requires parameters.
	doTemplateRequireParameter := false
	for _, comp := range templateInUse.Components {
		if len(comp.Example.BodyText) > 0 || len(comp.Example.HeaderText) > 0 || len(comp.Example.HeaderHandle) > 0 {
			doTemplateRequireParameter = true
		}
		if len(comp.Buttons) > 0 {
			for _, btn := range comp.Buttons {
				if len(btn.Example) > 0 {
					doTemplateRequireParameter = true
				}
			}
		}
	}

	// Unmarshal stored parameters from the database.
	var params templateComponentParameters
	if err = json.Unmarshal([]byte(*message.Campaign.TemplateMessageComponentParameters), &params); err != nil {
		return fmt.Errorf("error unmarshalling template parameters: %v", err)
	}

	if doTemplateRequireParameter && reflect.DeepEqual(params, templateComponentParameters{}) {
		cm.StopCampaign(message.Campaign.UniqueId.String())
		return fmt.Errorf("template requires parameters, but no parameter found in the database")
	}

	// Build the template message using our helper functions.
	templateMessage, err := cm.buildTemplateMessage(templateInUse, params)
	if err != nil {
		return fmt.Errorf("error building template message: %v", err)
	}

	// Send the message using the messaging client.
	messagingClient := client.NewMessagingClient(message.Campaign.PhoneNumberToUse)
	response, err := messagingClient.Message.Send(templateMessage, message.Contact.PhoneNumber)

	cm.Logger.Info("Message send response", response)

	messageStatus := model.MessageStatusEnum_Sent
	if err != nil {
		cm.Logger.Error("error sending message to user", err.Error())
		message.Campaign.ErrorCount.Add(1)
		messageStatus = model.MessageStatusEnum_Failed
		return err
	}

	// Convert the sent message to JSON for record keeping.
	jsonMessage, err := templateMessage.ToJson(wapiComponents.ApiCompatibleJsonConverterConfigs{
		SendToPhoneNumber: message.Contact.PhoneNumber,
	})
	if err != nil {
		return err
	}
	stringifiedJsonMessage := string(jsonMessage)

	// Save a record of the sent message to the database.
	messageSent := model.Message{
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		CampaignId:      &message.Campaign.UniqueId,
		Direction:       model.MessageDirectionEnum_OutBound,
		ContactId:       message.Contact.UniqueId,
		PhoneNumberUsed: message.Campaign.PhoneNumberToUse,
		OrganizationId:  message.Campaign.OrganizationId,
		MessageData:     &stringifiedJsonMessage,
		MessageType:     model.MessageTypeEnum_Template,
		Status:          messageStatus,
	}

	messageSentRecordQuery := table.Message.
		INSERT(table.Message.MutableColumns).
		MODEL(messageSent).
		RETURNING(table.Message.AllColumns)

	if err = messageSentRecordQuery.Query(cm.Db, &messageSent); err != nil {
		cm.Logger.Error("error saving message record to the database", err.Error())
	}

	// Update rate limiter and campaign counters.
	worker.rateLimiter.Incr(1)
	message.Campaign.Sent.Add(1)
	return nil
}
