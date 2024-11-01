package contact_service

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/internal/api_types"
	"github.com/sarthakjdev/wapikit/internal/core/utils"
	"github.com/sarthakjdev/wapikit/internal/interfaces"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/sarthakjdev/wapikit/.db-generated/model"
	table "github.com/sarthakjdev/wapikit/.db-generated/table"
)

type ContactService struct {
	services.BaseService `json:"-,inline"`
}

func NewContactService() *ContactService {
	return &ContactService{
		BaseService: services.BaseService{
			Name:        "Contact Service",
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
							MaxRequests:    100,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
					},
				},
				{
					Path:                    "/api/contacts",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(createNewContacts),
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/contacts/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getContactById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60, // 1 minute
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
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
					},
				},
				{
					Path:                    "/api/contacts/bulkImport",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(bulkImport),
					IsAuthorizationRequired: true,
				},
			},
		},
	}
}

func getContacts(context interfaces.ContextWithSession) error {
	params := new(api_types.GetContactsParams)

	err := utils.BindQueryParams(context, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	page := params.Page
	limit := params.PerPage
	listIds := params.ListId
	order := params.Order
	status := params.Status

	if page == 0 || limit > 50 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid page or perPage value")
	}

	// ! TODO: need to have it as slice here
	var dest []struct {
		TotalContacts int `json:"totalContacts"`
		model.Contact
		ContactLists []struct {
			model.ContactList
		}
	}

	orgUuid, _ := uuid.Parse(context.Session.User.OrganizationId)
	whereCondition := table.Contact.OrganizationId.EQ(UUID(orgUuid))

	if listIds != nil {
		listsIdArray := strings.Split(*listIds, ",")
		expressionArr := make([]Expression, len(listsIdArray))
		for _, listId := range listsIdArray {
			expressionArr = append(expressionArr, String(listId))
		}
		// ! TODO: verify there might be bug in the IN expression here
		whereCondition.AND(table.ContactListContact.ContactListId.IN(expressionArr...))
	}

	contactsQuery := SELECT(
		table.Contact.AllColumns,
		table.ContactListContact.AllColumns,
		table.ContactList.AllColumns,
		COUNT(table.Contact.UniqueId).OVER().AS("totalContacts"),
	).
		FROM(table.Contact.
			LEFT_JOIN(table.ContactListContact, table.ContactListContact.ContactId.EQ(table.Contact.UniqueId)).
			LEFT_JOIN(table.ContactList, table.Contact.UniqueId.EQ(table.ContactListContact.ContactId)),
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
		if err.Error() == "qrm: no rows in result set" {
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
				UniqueId:   contactId,
				CreatedAt:  contact.CreatedAt,
				Name:       contact.Name,
				Lists:      lists,
				Phone:      contact.PhoneNumber,
				Attributes: attr,
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

	insertedContact := []model.Contact{}

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for _, contact := range *payload {
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		jsonAttributes, _ := json.Marshal(contact.Attributes)
		stringAttributes := string(jsonAttributes)
		contactToInsert := model.Contact{
			OrganizationId: orgUuid,
			Name:           contact.Name,
			PhoneNumber:    contact.Phone,
			Attributes:     &stringAttributes,
			Status:         model.ContactStatus_Active,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		insertedContact = append(insertedContact, contactToInsert)
	}

	insertQuery := table.Contact.
		INSERT(table.Contact.MutableColumns).
		MODELS(insertedContact).
		ON_CONFLICT(table.Contact.PhoneNumber, table.Contact.OrganizationId).
		DO_NOTHING()

	stringQuery := insertQuery.DebugSql()

	fmt.Println(stringQuery)

	result, err := insertQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	numberOfRows, _ := result.RowsAffected()

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
	}

	orgUuid, _ := uuid.Parse(context.Session.User.OrganizationId)
	contactUuid, _ := uuid.Parse(contactId)

	contactsQuery := SELECT(
		table.Contact.AllColumns,
		table.ContactListContact.AllColumns,
		table.ContactList.AllColumns,
	).
		FROM(table.Contact.
			LEFT_JOIN(table.ContactListContact, table.ContactListContact.ContactId.EQ(table.Contact.UniqueId)).
			LEFT_JOIN(table.ContactList, table.Contact.UniqueId.EQ(table.ContactListContact.ContactId)),
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

	contactIdString := dest.UniqueId.String()
	attr := map[string]interface{}{}
	json.Unmarshal([]byte(*dest.Attributes), &attr)

	return context.JSON(http.StatusOK, api_types.GetContactByIdResponseSchema{
		Contact: api_types.ContactSchema{
			UniqueId:   contactIdString,
			CreatedAt:  dest.CreatedAt,
			Name:       dest.Name,
			Lists:      lists,
			Phone:      dest.PhoneNumber,
			Attributes: attr,
		},
	})
}

func bulkImport(context interfaces.ContextWithSession) error {

	payload := new(api_types.BulkImportSchema)
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Retrieve the CSV file from the request
	file, err := context.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "CSV file is required")
	}

	// Open the CSV file
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	delimeter := ','

	if payload.Delimiter != nil && len(*payload.Delimiter) != 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Delimiter must be a single character")
	}

	if payload.Delimiter != nil {
		delimeter = rune((*payload.Delimiter)[0])
	}

	// Initialize the CSV reader
	reader := csv.NewReader(src)
	reader.Comma = delimeter

	var importedContacts []model.Contact
	orgUuid, _ := uuid.Parse(context.Session.User.OrganizationId)

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
			Status:         model.ContactStatus_Active,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		importedContacts = append(importedContacts, contact)
	}

	// Insert contacts into the database
	insertQuery := table.Contact.
		INSERT(table.Contact.MutableColumns).
		MODELS(importedContacts).
		ON_CONFLICT(table.Contact.PhoneNumber, table.Contact.OrganizationId).
		DO_NOTHING()

	_, err = insertQuery.ExecContext(context.Request().Context(), context.App.Db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to insert contacts")
	}

	listId := *payload.ListId

	if listId == "" {
		return context.JSON(http.StatusOK, api_types.BulkImportResponseSchema{
			Message: strconv.Itoa(len(importedContacts)) + " contacts imported successfully",
		})
	}

	// Parse the List ID into a UUID
	listUUID, err := uuid.Parse(listId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid list ID format")
	}

	// Associate imported contacts with the specified list
	for _, contact := range importedContacts {
		associateQuery := table.ContactListContact.
			INSERT(table.ContactListContact.ContactListId, table.ContactListContact.ContactId).
			VALUES(listUUID, contact.UniqueId).
			ON_CONFLICT(table.ContactListContact.ContactId, table.ContactListContact.ContactListId).
			DO_NOTHING()
		_, err = associateQuery.ExecContext(context.Request().Context(), context.App.Db)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to associate contacts with list")
		}
	}

	// Prepare a success message
	response := api_types.BulkImportResponseSchema{
		Message: strconv.Itoa(len(importedContacts)) + " contacts imported successfully",
	}

	return context.JSON(http.StatusOK, response)
}

func deleteContactById(context interfaces.ContextWithSession) error {
	contactId := context.Param("id")

	if contactId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid contact id")
	}

	// ! TODO: check if there is any conversation associated with this contact
	// ! TODO: also before deleting the contact, remove the contact from all the lists and delete all their messages

	contactQuery := table.Contact.DELETE().WHERE(table.Contact.UniqueId.EQ(String(contactId)))
	result, err := contactQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if res, _ := result.RowsAffected(); res == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Contact not found")
	}

	return context.String(http.StatusOK, "OK")
}
