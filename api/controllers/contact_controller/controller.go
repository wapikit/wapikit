package contact_controller

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/wapikit/wapikit/api/api_types"
	controller "github.com/wapikit/wapikit/api/controllers"
	"github.com/wapikit/wapikit/interfaces"
	"github.com/wapikit/wapikit/utils"

	"github.com/go-jet/jet/qrm"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
)

type ContactController struct {
	controller.BaseController `json:"-,inline"`
}

func NewContactController() *ContactController {
	return &ContactController{
		BaseController: controller.BaseController{
			Name:        "Contact Controller",
			RestApiPath: "/api/contact",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/contacts",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getContacts),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    600,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetContact,
						},
					},
				},
				{
					Path:                    "/api/contacts",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(createNewContacts),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    600,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.CreateContact,
						},
					},
				},
				{
					Path:                    "/api/contacts/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getContactById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetContact,
						},
					},
				},
				{
					Path:                    "/api/contacts/:id",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(updateContactById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.UpdateContact,
						},
					},
				},
				{
					Path:                    "/api/contacts/:id",
					Method:                  http.MethodDelete,
					Handler:                 interfaces.HandlerWithSession(deleteContactById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.DeleteContact,
						},
					},
				},
				{
					Path:                    "/api/contacts/bulkImport",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(bulkImport),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.BulkImportContacts,
							api_types.CreateContact,
						},
					},
				},
			},
		},
	}
}

func getContacts(context interfaces.ContextWithSession) error {
	logger := context.App.Logger
	params := new(api_types.GetContactsParams)

	err := utils.BindQueryParams(context, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	page := params.Page
	limit := params.PerPage
	listId := params.ListId
	order := params.Order
	status := params.Status

	if page == 0 || limit > 50 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid page or perPage value")
	}

	var dest []struct {
		TotalContacts int `json:"totalContacts"`
		model.Contact
		ContactLists []struct {
			model.ContactList
		}
		Conversations []struct {
			model.Conversation
			Messages []model.Message
		}
	}

	orgUuid, _ := uuid.Parse(context.Session.User.OrganizationId)
	whereCondition := table.Contact.OrganizationId.EQ(UUID(orgUuid))

	if listId != nil {
		logger.Debug("List ID:", *listId, nil)
		listUuid, err := uuid.Parse(*listId)
		if err != nil {
			// * skip the list if it is not a valid UUID
		} else {
			whereCondition = whereCondition.AND(table.ContactListContact.ContactListId.EQ(UUID(listUuid)))
		}
	}

	contactsQuery := SELECT(
		table.Contact.AllColumns,
		table.ContactListContact.AllColumns,
		table.ContactList.AllColumns,
		COUNT(table.Contact.UniqueId).OVER().AS("totalContacts"),
		table.Conversation.AllColumns,
		table.Message.AllColumns,
	).
		FROM(table.Contact.
			LEFT_JOIN(table.ContactListContact, table.ContactListContact.ContactId.EQ(table.Contact.UniqueId)).
			LEFT_JOIN(table.ContactList, table.ContactList.UniqueId.EQ(table.ContactListContact.ContactListId)).
			LEFT_JOIN(table.Conversation, table.Conversation.ContactId.EQ(table.Contact.UniqueId).AND(table.Conversation.OrganizationId.EQ(UUID(orgUuid)))).
			LEFT_JOIN(table.Message, table.Message.ConversationId.EQ(table.Conversation.UniqueId).AND(table.Message.OrganizationId.EQ(UUID(orgUuid)))),
		).
		WHERE(whereCondition).
		LIMIT(limit).
		OFFSET((page - 1) * limit)

	if order != nil {
		if *order == api_types.Asc {
			contactsQuery.ORDER_BY(table.Contact.CreatedAt.ASC())
		} else {
			contactsQuery.ORDER_BY(table.Contact.CreatedAt.DESC())
		}
	}

	if status != nil {
		whereCondition.AND(table.Contact.Status.EQ(String(*status)))
	}

	err = contactsQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			total := 0
			contacts := make([]api_types.ContactSchema, 0)
			return context.JSON(http.StatusOK, api_types.GetContactsResponseSchema{
				Contacts: contacts,
				PaginationMeta: api_types.PaginationMeta{
					Page:    page,
					PerPage: limit,
					Total:   total,
				},
			})
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	contactsToReturn := []api_types.ContactSchema{}
	totalContacts := 0

	if len(dest) > 0 {
		for _, contact := range dest {

			conversations := []api_types.ConversationWithoutContactSchema{}

			for _, conversation := range contact.Conversations {
				messages := []api_types.MessageSchema{}

				for _, message := range conversation.Messages {
					messageData := map[string]interface{}{}
					json.Unmarshal([]byte(*message.MessageData), &messageData)
					messageToAppend := api_types.MessageSchema{
						UniqueId:       message.UniqueId.String(),
						ConversationId: message.ConversationId.String(),
						CreatedAt:      message.CreatedAt,
						Direction:      api_types.MessageDirectionEnum(message.Direction.String()),
						MessageData:    &messageData,
						MessageType:    api_types.MessageTypeEnum(message.MessageType.String()),
						Status:         api_types.MessageStatusEnum(message.Status.String()),
					}
					messages = append(messages, messageToAppend)
				}

				campaignId := ""

				if conversation.InitiatedByCampaignId != nil {
					campaignId = string(conversation.InitiatedByCampaignId.String())
				}

				conversationToAppend := api_types.ConversationWithoutContactSchema{
					UniqueId:       conversation.UniqueId.String(),
					CreatedAt:      conversation.CreatedAt,
					Messages:       messages,
					ContactId:      conversation.ContactId.String(),
					OrganizationId: conversation.OrganizationId.String(),
					InitiatedBy:    api_types.ConversationInitiatedByEnum(conversation.InitiatedBy.String()),
					CampaignId:     &campaignId,
					Status:         api_types.ConversationStatusEnum(conversation.Status.String()),
				}

				conversations = append(conversations, conversationToAppend)
			}

			lists := []api_types.ContactListSchema{}

			for _, contactList := range contact.ContactLists {
				stringUniqueId := contactList.UniqueId.String()
				listToAppend := api_types.ContactListSchema{
					UniqueId: stringUniqueId,
					Name:     contactList.Name,
				}
				lists = append(lists, listToAppend)
			}
			contactId := contact.UniqueId.String()
			attr := map[string]interface{}{}
			json.Unmarshal([]byte(*contact.Attributes), &attr)
			cntct := api_types.ContactSchema{
				UniqueId:      contactId,
				CreatedAt:     contact.CreatedAt,
				Name:          contact.Name,
				Lists:         lists,
				Phone:         contact.PhoneNumber,
				Attributes:    attr,
				Status:        api_types.ContactStatusEnum(contact.Status),
				Conversations: &conversations,
			}
			contactsToReturn = append(contactsToReturn, cntct)
		}

		totalContacts = dest[0].TotalContacts

	}

	return context.JSON(http.StatusOK, api_types.GetContactsResponseSchema{
		Contacts: contactsToReturn,
		PaginationMeta: api_types.PaginationMeta{
			Page:    page,
			PerPage: limit,
			Total:   totalContacts,
		},
	})
}

func createNewContacts(context interfaces.ContextWithSession) error {
	payload := new(api_types.CreateContactsJSONBody)
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// * insert contact into the contact table
	contactsToInsert := []model.Contact{}
	var insertedContacts []model.Contact

	for _, contact := range *payload {
		jsonAttributes, _ := json.Marshal(contact.Attributes)
		stringAttributes := string(jsonAttributes)
		contactToInsert := model.Contact{
			OrganizationId: orgUuid,
			Name:           contact.Name,
			PhoneNumber:    contact.Phone,
			Attributes:     &stringAttributes,
			Status:         model.ContactStatusEnum_Active,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		contactsToInsert = append(contactsToInsert, contactToInsert)
	}

	insertQuery := table.Contact.
		INSERT(table.Contact.MutableColumns).
		MODELS(contactsToInsert).
		ON_CONFLICT(table.Contact.PhoneNumber, table.Contact.OrganizationId).
		DO_NOTHING().
		RETURNING(table.Contact.AllColumns)

	err = insertQuery.QueryContext(context.Request().Context(), context.App.Db, &insertedContacts)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// * insert contact into the contact list contact table

	insertedContactListContact := []model.ContactListContact{}

	for _, contact := range *payload {
		if len(contact.ListsIds) > 0 {
			// * find the inserted db record for this contact to get the uniqueId
			insertedDbRecordOfThisContact := model.Contact{}

			for _, contactRecord := range insertedContacts {
				if contactRecord.PhoneNumber == contact.Phone {
					insertedDbRecordOfThisContact = contactRecord
					break
				}
			}

			if insertedDbRecordOfThisContact.UniqueId == uuid.Nil {
				// * skip this contact if it is not inserted
				continue
			}

			for _, listId := range contact.ListsIds {
				listUuid, err := uuid.Parse(listId)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, "Invalid list ID format")
				}

				contactListContact := model.ContactListContact{
					ContactId:     insertedDbRecordOfThisContact.UniqueId,
					ContactListId: listUuid,
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}

				insertedContactListContact = append(insertedContactListContact, contactListContact)
			}
		}
	}

	insertedContactListContactQuery := table.ContactListContact.
		INSERT().
		MODELS(insertedContactListContact).
		ON_CONFLICT(table.ContactListContact.ContactId, table.ContactListContact.ContactListId).
		DO_NOTHING()

	_, err = insertedContactListContactQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	numberOfRows := len(contactsToInsert)

	response := api_types.CreateNewContactResponseSchema{
		Message: strings.Join([]string{"Successfully created ", string(numberOfRows), " contacts"}, " "),
	}

	return context.JSON(http.StatusOK, response)
}

func getContactById(context interfaces.ContextWithSession) error {
	contactId := context.Param("id")
	if contactId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid contact id")
	}

	var dest struct {
		model.Contact
		ContactLists []struct {
			model.ContactList
		}
		Conversations []struct {
			model.Conversation
			Messages []model.Message
		}
	}

	orgUuid, _ := uuid.Parse(context.Session.User.OrganizationId)
	contactUuid, _ := uuid.Parse(contactId)

	contactsQuery := SELECT(
		table.Contact.AllColumns,
		table.ContactListContact.AllColumns,
		table.ContactList.AllColumns,
		table.Conversation.AllColumns,
	).
		FROM(table.Contact.
			LEFT_JOIN(table.ContactListContact, table.ContactListContact.ContactId.EQ(table.Contact.UniqueId)).
			LEFT_JOIN(table.ContactList, table.ContactList.UniqueId.EQ(table.ContactListContact.ContactListId)).
			LEFT_JOIN(table.Conversation, table.Conversation.ContactId.EQ(table.Contact.UniqueId).AND(table.Conversation.OrganizationId.EQ(UUID(orgUuid)))).
			LEFT_JOIN(table.Message, table.Message.ConversationId.EQ(table.Conversation.UniqueId).AND(table.Message.OrganizationId.EQ(UUID(orgUuid)))),
		).
		WHERE(table.Contact.OrganizationId.EQ(UUID(orgUuid)).AND(table.Contact.UniqueId.EQ(UUID(contactUuid))))

	err := contactsQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	lists := []api_types.ContactListSchema{}

	for _, contactList := range dest.ContactLists {
		stringUniqueId := contactList.UniqueId.String()
		listToAppend := api_types.ContactListSchema{
			UniqueId: stringUniqueId,
			Name:     contactList.Name,
		}
		lists = append(lists, listToAppend)
	}

	conversations := []api_types.ConversationWithoutContactSchema{}

	for _, conversation := range dest.Conversations {
		messages := []api_types.MessageSchema{}

		for _, message := range conversation.Messages {
			messageData := map[string]interface{}{}
			json.Unmarshal([]byte(*message.MessageData), &messageData)
			messageToAppend := api_types.MessageSchema{
				UniqueId:       message.UniqueId.String(),
				ConversationId: message.ConversationId.String(),
				CreatedAt:      message.CreatedAt,
				Direction:      api_types.MessageDirectionEnum(message.Direction.String()),
				MessageData:    &messageData,
				MessageType:    api_types.MessageTypeEnum(message.MessageType.String()),
				Status:         api_types.MessageStatusEnum(message.Status.String()),
			}
			messages = append(messages, messageToAppend)
		}

		campaignId := ""

		if conversation.InitiatedByCampaignId != nil {
			campaignId = string(conversation.InitiatedByCampaignId.String())
		}

		conversationToAppend := api_types.ConversationWithoutContactSchema{
			UniqueId:       conversation.UniqueId.String(),
			CreatedAt:      conversation.CreatedAt,
			Messages:       messages,
			ContactId:      conversation.ContactId.String(),
			OrganizationId: conversation.OrganizationId.String(),
			InitiatedBy:    api_types.ConversationInitiatedByEnum(conversation.InitiatedBy.String()),
			CampaignId:     &campaignId,
			Status:         api_types.ConversationStatusEnum(conversation.Status.String()),
		}

		conversations = append(conversations, conversationToAppend)
	}

	contactIdString := dest.UniqueId.String()
	attr := map[string]interface{}{}
	json.Unmarshal([]byte(*dest.Attributes), &attr)

	return context.JSON(http.StatusOK, api_types.GetContactByIdResponseSchema{
		Contact: api_types.ContactSchema{
			UniqueId:      contactIdString,
			CreatedAt:     dest.CreatedAt,
			Name:          dest.Name,
			Lists:         lists,
			Phone:         dest.PhoneNumber,
			Attributes:    attr,
			Status:        api_types.ContactStatusEnum(dest.Status),
			Conversations: &conversations,
		},
	})
}

func updateContactById(context interfaces.ContextWithSession) error {
	logger := context.App.Logger
	contactId := context.Param("id")
	if contactId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid contact id")
	}

	orgUuid, _ := uuid.Parse(context.Session.User.OrganizationId)
	contactUuid, _ := uuid.Parse(contactId)

	var existingContact struct {
		model.Contact
		ContactLists []struct {
			model.ContactList
		}
	}

	contactsQuery := SELECT(
		table.Contact.AllColumns,
		table.ContactListContact.AllColumns,
		table.ContactList.AllColumns,
	).
		FROM(table.Contact.
			LEFT_JOIN(table.ContactListContact, table.ContactListContact.ContactId.EQ(table.Contact.UniqueId)).
			LEFT_JOIN(table.ContactList, table.ContactList.UniqueId.EQ(table.ContactListContact.ContactListId)),
		).
		WHERE(table.Contact.OrganizationId.EQ(UUID(orgUuid)).AND(table.Contact.UniqueId.EQ(UUID(contactUuid))))

	err := contactsQuery.QueryContext(context.Request().Context(), context.App.Db, &existingContact)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return echo.NewHTTPError(http.StatusNotFound, "Contact not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	payload := new(api_types.UpdateContactSchema)
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// * updating lists

	// * CTE for both these cases
	// * case 1: if the contact is not in any list and the payload has lists
	// * case 2: if the contact is in some lists and the payload has lists

	oldListsUuids := make([]uuid.UUID, 0)
	newListsUuids := make([]uuid.UUID, 0)

	for _, list := range existingContact.ContactLists {
		oldListsUuids = append(oldListsUuids, list.UniqueId)
	}

	for _, listId := range payload.Lists {
		listUuid, err := uuid.Parse(listId)
		if err != nil {
			continue
		}
		newListsUuids = append(newListsUuids, listUuid)
	}

	listsToBeDeleted := make([]Expression, 0)
	listsToBeInserted := make([]model.ContactListContact, 0)

	commonListIds := make([]uuid.UUID, 0)

	// * the list ids that are in oldListsUuids but not in newListsUuids are needed to be deleted
	for _, oldList := range oldListsUuids {
		found := false
		for _, newList := range newListsUuids {
			if oldList == newList {
				found = true
				commonListIds = append(commonListIds, oldList)
				break
			}
		}
		if !found {
			listsToBeDeleted = append(listsToBeDeleted, UUID(oldList))
		}
	}

	// * the lists ids that are in newListsUuids but not in oldListsUuids are needed to be inserted
	for _, newList := range newListsUuids {
		found := false
		for _, oldList := range oldListsUuids {
			if newList == oldList {
				found = true
				commonListIds = append(commonListIds, newList)
				break
			}

		}

		if !found {
			contactListContact := model.ContactListContact{
				ContactId:     existingContact.UniqueId,
				ContactListId: newList,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			listsToBeInserted = append(listsToBeInserted, contactListContact)
		}
	}

	if len(listsToBeDeleted) > 0 {
		deleteQuery := table.ContactListContact.
			DELETE().
			WHERE(table.ContactListContact.ContactId.EQ(UUID(existingContact.UniqueId)).
				AND(table.ContactListContact.ContactListId.IN(listsToBeDeleted...)))

		_, err = deleteQuery.ExecContext(context.Request().Context(), context.App.Db)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	var insertedLists []model.ContactList

	if len(listsToBeInserted) > 0 {
		listToBeInsertedExpression := make([]Expression, 0)
		for _, list := range listsToBeInserted {
			listToBeInsertedExpression = append(listToBeInsertedExpression, UUID(list.ContactListId))
		}
		listToBeInsertedCte := CTE("lists_to_be_inserted")
		contactListContactInsertQuery := WITH(
			listToBeInsertedCte.AS(
				SELECT(table.ContactList.AllColumns).FROM(
					table.ContactList,
				).WHERE(
					table.ContactList.UniqueId.IN(listToBeInsertedExpression...),
				),
			),
			CTE("insert_list").AS(
				table.ContactListContact.
					INSERT().
					MODELS(listsToBeInserted).
					ON_CONFLICT(table.ContactListContact.ContactId, table.ContactListContact.ContactListId).
					DO_NOTHING(),
			),
		)(
			SELECT(listToBeInsertedCte.AllColumns()).FROM(listToBeInsertedCte),
		)

		err = contactListContactInsertQuery.QueryContext(context.Request().Context(), context.App.Db, &insertedLists)

		if err != nil {
			logger.Error("Error inserting lists:", err.Error(), nil)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	// * update rest of the contact details

	jsonAttributes, _ := json.Marshal(payload.Attributes)
	stringAttributes := string(jsonAttributes)

	contactToUpdate := model.Contact{
		UpdatedAt:   time.Now(),
		Name:        payload.Name,
		PhoneNumber: payload.Phone,
		Attributes:  &stringAttributes,
		Status:      model.ContactStatusEnum(payload.Status),
	}

	var updatedContact model.Contact

	updateQuery := table.Contact.
		UPDATE(table.Contact.Name, table.Contact.PhoneNumber, table.Contact.Attributes, table.Contact.Status, table.Contact.UpdatedAt).
		MODEL(contactToUpdate).
		WHERE(table.Contact.UniqueId.EQ(UUID(existingContact.UniqueId))).
		RETURNING(table.Contact.AllColumns)

	err = updateQuery.QueryContext(context.Request().Context(), context.App.Db, &updatedContact)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	listToReturn := []api_types.ContactListSchema{}

	for _, list := range existingContact.ContactLists {
		// * if the list is not in the commonListIds, then it is not in the payload lists
		found := false
		for _, commonListId := range commonListIds {
			if list.UniqueId == commonListId {
				found = true
				break
			}
		}
		if !found {
			continue
		}

		stringUniqueId := list.UniqueId.String()
		listToAppend := api_types.ContactListSchema{
			UniqueId: stringUniqueId,
			Name:     list.Name,
		}
		listToReturn = append(listToReturn, listToAppend)
	}

	for _, list := range insertedLists {
		stringUniqueId := list.UniqueId.String()
		listToAppend := api_types.ContactListSchema{
			UniqueId: stringUniqueId,
			Name:     list.Name,
		}

		listToReturn = append(listToReturn, listToAppend)
	}

	return context.JSON(http.StatusOK, api_types.UpdateContactByIdResponseSchema{
		Contact: api_types.ContactSchema{
			UniqueId:   updatedContact.UniqueId.String(),
			CreatedAt:  updatedContact.CreatedAt,
			Name:       updatedContact.Name,
			Attributes: payload.Attributes,
			Phone:      updatedContact.PhoneNumber,
			Lists:      listToReturn,
		},
	})
}

// ! TODO: change this to a streaming endpoint
func bulkImport(context interfaces.ContextWithSession) error {
	logger := context.App.Logger

	r := context.Request()

	err := r.ParseMultipartForm(10 << 20) // 10 MB max memory
	if err != nil {
		logger.Error("Error parsing form data:", err.Error(), nil)
	}

	// Get the file
	file, _, err := r.FormFile("file")
	if err != nil {
		logger.Error("Error getting file:", err.Error(), nil)
		return echo.NewHTTPError(http.StatusBadRequest, "Error getting file")
	}
	defer file.Close()

	// Get the delimiter
	payloadDelimiter := r.FormValue("delimiter")

	// Get listIds as a JSON string, then parse it into a slice
	listIdsStr := r.FormValue("listIds")
	var listIds []string

	err = json.Unmarshal([]byte(listIdsStr), &listIds)
	if err != nil {
		logger.Error("Error parsing list IDs:", err.Error(), nil)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid list IDs")
	}

	delimeter := ','

	if payloadDelimiter == "" && len(payloadDelimiter) != 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Delimiter must be a single character")
	}

	if payloadDelimiter != "" {
		delimeter = rune((payloadDelimiter)[0])
	}

	// Initialize the CSV reader
	reader := csv.NewReader(file)
	reader.Comma = delimeter

	var contactToImport []model.Contact
	orgUuid, _ := uuid.Parse(context.Session.User.OrganizationId)

	// Skip the first row (header)
	if _, err := reader.Read(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid CSV format or empty file")
	}

	// Iterate through CSV rows and parse each contact
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid CSV format")
		}

		// Assuming the CSV columns: Name, Phone, Attributes (in JSON format)
		if len(record) < 3 {
			return echo.NewHTTPError(http.StatusBadRequest, "Each row must contain Name, Phone, and Attributes")
		}

		name := record[0]
		phone := record[1]
		attributes := record[2]

		// Convert attributes JSON string to map
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attributes), &attrMap); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format in attributes")
		}

		// Prepare the contact to insert
		jsonAttributes, _ := json.Marshal(attrMap)
		stringAttributes := string(jsonAttributes)

		contact := model.Contact{
			OrganizationId: orgUuid,
			Name:           name,
			PhoneNumber:    phone,
			Attributes:     &stringAttributes,
			Status:         model.ContactStatusEnum_Active,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		contactToImport = append(contactToImport, contact)
	}

	// Insert contacts into the database
	insertQuery := table.Contact.
		INSERT(table.Contact.MutableColumns).
		MODELS(contactToImport).
		ON_CONFLICT(table.Contact.PhoneNumber, table.Contact.OrganizationId).
		DO_NOTHING().
		RETURNING(table.Contact.AllColumns)

	var importedContact []model.Contact

	err = insertQuery.QueryContext(context.Request().Context(), context.App.Db, &importedContact)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to insert contacts")
	}

	if len(listIds) == 0 {
		return context.JSON(http.StatusOK, api_types.BulkImportResponseSchema{
			Message: strconv.Itoa(len(contactToImport)) + " contacts imported successfully",
		})
	}

	listUuids := make([]uuid.UUID, 0)

	for _, listId := range listIds {
		var list model.ContactList
		listUuid, err := uuid.Parse(listId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid list ID format")
		}

		// check if the list exists

		listQuery := table.ContactList.
			SELECT(table.ContactList.UniqueId).
			WHERE(table.ContactList.UniqueId.EQ(UUID(listUuid)))

		err = listQuery.QueryContext(context.Request().Context(), context.App.Db, &list)
		if err != nil {
			// do nothing
		} else {
			listUuids = append(listUuids, listUuid)
		}
	}

	if len(listUuids) == 0 {
		return context.JSON(http.StatusOK, api_types.BulkImportResponseSchema{
			Message: strconv.Itoa(len(contactToImport)) + " contacts imported successfully",
		})
	}

	for _, listId := range listUuids {
		var records []model.ContactListContact
		for _, contact := range importedContact {
			contactListContact := model.ContactListContact{
				ContactId:     contact.UniqueId,
				ContactListId: listId,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			records = append(records, contactListContact)
		}

		if len(records) == 0 {
			return context.JSON(http.StatusOK, api_types.BulkImportResponseSchema{
				Message: strconv.Itoa(len(contactToImport)) + " contacts imported successfully",
			})
		}

		contactInsertionToListQuery := table.ContactListContact.
			INSERT().
			MODELS(records).
			ON_CONFLICT(table.ContactListContact.ContactId, table.ContactListContact.ContactListId).
			DO_NOTHING()

		_, err = contactInsertionToListQuery.ExecContext(context.Request().Context(), context.App.Db)

		if err != nil {
			logger.Error("Error inserting contacts into list:", err.Error(), nil)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to insert contacts into list")
		}
	}

	// Prepare a success message
	response := api_types.BulkImportResponseSchema{
		Message: strconv.Itoa(len(contactToImport)) + " contacts imported successfully",
	}

	return context.JSON(http.StatusOK, response)
}

func deleteContactById(context interfaces.ContextWithSession) error {
	contactId := context.Param("id")

	if contactId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid contact id")
	}

	contactUuid, _ := uuid.Parse(contactId)

	// ! TODO: check if there is any conversation associated with this contact
	// ! TODO: also before deleting the contact, remove the contact from all the lists and delete all their messages

	contactQuery := table.Contact.DELETE().WHERE(table.Contact.UniqueId.EQ(UUID(contactUuid)))
	result, err := contactQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if res, _ := result.RowsAffected(); res == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Contact not found")
	}

	response := api_types.DeleteContactByIdResponseSchema{
		Data: true,
	}

	return context.JSON(http.StatusOK, response)
}
