package contact_service

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/database"
	"github.com/sarthakjdev/wapikit/internal"
	"github.com/sarthakjdev/wapikit/internal/api_types"
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
					Handler:                 GetContacts,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/contacts",
					Method:                  http.MethodPost,
					Handler:                 CreateNewContacts,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/contacts/:id",
					Method:                  http.MethodGet,
					Handler:                 GetContactById,
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Admin,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
					},
				},
				{
					Path:                    "/api/contacts/:id",
					Method:                  http.MethodDelete,
					Handler:                 DeleteContactById,
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Admin,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
					},
				},
			},
		},
	}
}

func GetContacts(context interfaces.CustomContext) error {
	params := new(api_types.GetContactsParams)

	err := internal.BindQueryParams(context, params)
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
	var dest struct {
		TotalContacts int `json:"totalContacts"`
		Contacts      []struct {
			model.Contact
			ContactLists []struct {
				model.ContactList
			}
		}
	}

	orgUuid, _ := uuid.FromBytes([]byte(context.Session.User.OrganizationId))
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
		COUNT(table.Contact.OrganizationId.EQ(UUID(orgUuid))).OVER().AS("totalContacts"),
	).
		FROM(table.Contact.
			LEFT_JOIN(table.ContactListContact, table.ContactListContact.ContactId.EQ(table.Contact.UniqueId)).
			LEFT_JOIN(table.ContactList, table.Contact.UniqueId.EQ(table.ContactListContact.ContactId)),
		).
		WHERE(whereCondition).
		LIMIT(limit).
		OFFSET(page * limit)

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

	if err.Error() == "qrm: no rows in result set" {
		total := 0
		contacts := make([]api_types.ContactSchema, 0)
		return context.JSON(http.StatusOK, api_types.GetContactsResponseSchema{
			Contacts: &contacts,
			PaginationMeta: &api_types.PaginationMeta{
				Page:    &page,
				PerPage: &limit,
				Total:   &total,
			},
		})

	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	contactsToReturn := []api_types.ContactSchema{}
	totalContacts := 0

	if len(dest.Contacts) > 0 {
		for _, contact := range dest.Contacts {
			lists := []api_types.ContactListSchema{}

			for _, contactList := range contact.ContactLists {
				stringUniqueId := contactList.UniqueId.String()
				listToAppend := api_types.ContactListSchema{
					UniqueId: &stringUniqueId,
					Name:     &contactList.Name,
				}
				lists = append(lists, listToAppend)
			}
			contactId := contact.UniqueId.String()
			attr := map[string]interface{}{}
			json.Unmarshal([]byte(*contact.Attributes), &attr)
			cntct := api_types.ContactSchema{
				UniqueId:   &contactId,
				CreatedAt:  &contact.CreatedAt,
				Name:       &contact.Name,
				Lists:      &lists,
				Phone:      &contact.PhoneNumber,
				Attributes: &attr,
			}
			contactsToReturn = append(contactsToReturn, cntct)
		}

		totalContacts = dest.TotalContacts
	}

	return context.JSON(http.StatusOK, api_types.GetContactsResponseSchema{
		Contacts: &contactsToReturn,
		PaginationMeta: &api_types.PaginationMeta{
			Page:    &page,
			PerPage: &limit,
			Total:   &totalContacts,
		},
	})

}

func CreateNewContacts(context interfaces.CustomContext) error {
	payload := new(interface{})
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	database.GetDbInstance()
	return context.String(http.StatusOK, "OK")
}

func GetContactById(context interfaces.CustomContext) error {

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

	contactsQuery := SELECT(
		table.Contact.AllColumns,
		table.ContactListContact.AllColumns,
		table.ContactList.AllColumns,
	).
		FROM(table.Contact.
			LEFT_JOIN(table.ContactListContact, table.ContactListContact.ContactId.EQ(table.Contact.UniqueId)).
			LEFT_JOIN(table.ContactList, table.Contact.UniqueId.EQ(table.ContactListContact.ContactId)),
		).
		WHERE(table.Contact.OrganizationId.EQ(String(context.Session.User.OrganizationId)).AND(table.Contact.UniqueId.EQ(String(contactId))))

	err := contactsQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	lists := []api_types.ContactListSchema{}

	for _, contactList := range dest.ContactLists {
		stringUniqueId := contactList.UniqueId.String()
		listToAppend := api_types.ContactListSchema{
			UniqueId: &stringUniqueId,
			Name:     &contactList.Name,
		}
		lists = append(lists, listToAppend)
	}

	contactIdString := dest.UniqueId.String()
	attr := map[string]interface{}{}
	json.Unmarshal([]byte(*dest.Attributes), &attr)

	return context.JSON(http.StatusOK, api_types.GetContactByIdResponseSchema{
		Contact: &api_types.ContactSchema{
			UniqueId:   &contactIdString,
			CreatedAt:  &dest.CreatedAt,
			Name:       &dest.Name,
			Lists:      &lists,
			Phone:      &dest.PhoneNumber,
			Attributes: &attr,
		},
	})
}

func DeleteContactById(context interfaces.CustomContext) error {
	contactId := context.Param("id")

	if contactId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid contact id")
	}

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
