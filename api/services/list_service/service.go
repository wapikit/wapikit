package contact_list_service

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/internal"
	"github.com/sarthakjdev/wapikit/internal/api_types"
	"github.com/sarthakjdev/wapikit/internal/interfaces"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/sarthakjdev/wapikit/.db-generated/model"
	table "github.com/sarthakjdev/wapikit/.db-generated/table"
)

type ContactListService struct {
	services.BaseService `json:"-,inline"`
}

func NewContactListService() *ContactListService {
	return &ContactListService{
		BaseService: services.BaseService{
			Name:        "Contact List Service",
			RestApiPath: "/api/contact-list",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/lists",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(GetContactLists),
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
					Path:                    "/api/lists",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(CreateNewContactLists),
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
					Path:                    "/api/lists/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(GetContactListById),
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
					Handler:                 interfaces.HandlerWithSession(DeleteContactListById),
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
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(UpdateContactListById),
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

func GetContactLists(context interfaces.ContextWithSession) error {
	params := new(api_types.GetContactListsParams)

	if err := internal.BindQueryParams(context, params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	order := params.Order
	pageNumber := params.Page
	pageSize := params.PerPage

	orgUuid, _ := uuid.Parse(context.Session.User.OrganizationId)
	whereCondition := table.ContactList.OrganizationId.EQ(UUID(orgUuid))

	listsQuery := SELECT(
		table.ContactList.AllColumns,
		table.Tag.AllColumns,
		COUNT(table.ContactList.UniqueId).OVER().AS("totalLists"),
		// COUNT(table.Contact.UniqueId).OVER().AS("totalContacts"),
		// COUNT(table.Campaign.UniqueId).
		// 	OVER().
		// 	AS("totalCampaigns"),
	).
		FROM(
			table.ContactList.
				LEFT_JOIN(table.ContactListTag, table.ContactListTag.ContactListId.EQ(table.ContactList.UniqueId)).
				LEFT_JOIN(table.Tag, table.Tag.UniqueId.EQ(table.ContactListTag.TagId))).
		WHERE(whereCondition).
		LIMIT(pageSize).
		OFFSET((pageNumber - 1) * pageSize)

	if order != nil {
		if *order == api_types.Asc {
			listsQuery.ORDER_BY(table.ContactList.CreatedAt.ASC())
		} else {
			listsQuery.ORDER_BY(table.ContactList.CreatedAt.DESC())
		}
	}

	var dest struct {
		TotalLists int `json:"totalLists"`
		Lists      []struct {
			model.ContactList
			TotalContacts  int `json:"totalContacts"`
			TotalCampaigns int `json:"totalCampaigns"`
			Tags           []struct {
				model.Tag
			}
		}
	}

	err := listsQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	jsonLists, _ := json.Marshal(dest)
	context.App.Logger.Info("Lists: %v", jsonLists)

	if err != nil {
		if err.Error() == "qrm: no rows in result set" {
			total := 0
			lists := make([]api_types.ContactListSchema, 0)
			return context.JSON(http.StatusOK, api_types.GetContactListResponseSchema{
				Lists: lists,
				PaginationMeta: api_types.PaginationMeta{
					Page:    pageNumber,
					PerPage: pageSize,
					Total:   total,
				},
			})
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	listsToReturn := []api_types.ContactListSchema{}

	if len(dest.Lists) > 0 {
		for _, list := range dest.Lists {
			tags := []api_types.TagSchema{}
			if len(list.Tags) > 0 {
				for _, tag := range list.Tags {
					stringUniqueId := tag.UniqueId.String()
					tagToAppend := api_types.TagSchema{
						UniqueId: stringUniqueId,
						Name:     tag.Label,
					}

					tags = append(tags, tagToAppend)
				}
			}

			uniqueId := list.UniqueId.String()

			lst := api_types.ContactListSchema{
				CreatedAt:             list.CreatedAt,
				Name:                  list.Name,
				Description:           list.Name,
				NumberOfCampaignsSent: list.TotalCampaigns,
				NumberOfContacts:      list.TotalContacts,
				Tags:                  tags,
				UniqueId:              uniqueId,
			}
			listsToReturn = append(listsToReturn, lst)
		}

	}

	return context.JSON(http.StatusOK, api_types.GetContactListResponseSchema{
		Lists: listsToReturn,
		PaginationMeta: api_types.PaginationMeta{
			Page:    pageNumber,
			PerPage: pageSize,
			Total:   dest.TotalLists,
		},
	})
}

func CreateNewContactLists(context interfaces.ContextWithSession) error {
	return nil
}

func GetContactListById(context interfaces.ContextWithSession) error {

	contactListId := context.Param("id")
	if contactListId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Contact list id is required")
	}

	listUuid, _ := uuid.Parse(contactListId)
	orgId, _ := uuid.Parse(context.Session.User.OrganizationId)

	listQuery := SELECT(
		table.ContactList.AllColumns,
		table.Tag.AllColumns,
	).
		FROM(
			table.ContactList.
				LEFT_JOIN(table.ContactListTag, table.ContactListTag.ContactListId.EQ(table.ContactList.UniqueId)).
				LEFT_JOIN(table.Tag, table.Tag.UniqueId.EQ(table.ContactListTag.TagId)),
		).
		WHERE(
			table.ContactList.OrganizationId.EQ(UUID(orgId)).
				AND(table.ContactList.UniqueId.EQ(UUID(listUuid))),
		)

	var dest struct {
		TotalContacts  int `json:"totalContacts"`
		TotalCampaigns int `json:"totalCampaigns"`
		model.ContactList
		Tags []struct {
			model.Tag
		}
	}

	err := listQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tags := []api_types.TagSchema{}

	if len(dest.Tags) > 0 {
		for _, tag := range dest.Tags {
			stringUniqueId := tag.UniqueId.String()
			tagToAppend := api_types.TagSchema{
				UniqueId: stringUniqueId,
				Name:     tag.Label,
			}
			tags = append(tags, tagToAppend)
		}
	}

	uniqueId := dest.UniqueId.String()

	return context.JSON(http.StatusOK, api_types.GetContactListByIdSchema{
		List: api_types.ContactListSchema{
			CreatedAt:             dest.CreatedAt,
			Name:                  dest.Name,
			Description:           dest.Name,
			NumberOfCampaignsSent: dest.TotalCampaigns,
			NumberOfContacts:      dest.TotalContacts,
			Tags:                  tags,
			UniqueId:              uniqueId,
		},
	})
}

func DeleteContactListById(context interfaces.ContextWithSession) error {

	contactListId := context.Param("id")

	if contactListId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Contact list id is required")
	}

	// ! TODO: check for the running campaigns associated with this list, if there's any do not allow deleting the list

	deleteQuery := table.ContactList.
		DELETE().
		WHERE(
			table.ContactList.OrganizationId.EQ(String(context.Session.User.OrganizationId)).
				AND(table.ContactList.UniqueId.EQ(String(contactListId))),
		)

	result, err := deleteQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if res, _ := result.RowsAffected(); res == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "List not found")
	}

	return context.String(http.StatusOK, "OK")
}

func UpdateContactListById(context interfaces.ContextWithSession) error {
	return nil
}
