package contact_list_service

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/api/services"
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
					Handler:                 GetContactLists,
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
					Handler:                 CreateNewContactLists,
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
					Handler:                 GetContactListById,
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
					Handler:                 DeleteContactListById,
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
					Handler:                 UpdateContactListById,
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

func GetContactLists(context interfaces.CustomContext) error {
	params := new(api_types.GetContactListsParams)

	if err := context.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	order := params.Order
	pageNumber := params.Page
	pageSize := params.PerPage

	listsQuery := SELECT(
		table.ContactList.AllColumns,
		table.Tag.AllColumns,
		COUNT(table.ContactList.OrganizationId.EQ(String(context.Session.User.OrganizationId))).OVER().AS("totalLists"),
		COUNT(table.ContactListContact.ContactListId.EQ(table.ContactList.UniqueId)).OVER().AS("totalContacts"),
		COUNT(table.Campaign.OrganizationId.EQ(String(context.Session.User.OrganizationId)).
			AND(table.Campaign.Status.NOT_EQ(String(model.CampaignStatus_Draft.String()))).
			AND(table.CampaignList.ContactListId.EQ(table.ContactList.UniqueId))).
			OVER().
			AS("totalCampaigns"),
	).
		FROM(
			table.ContactList.
				LEFT_JOIN(table.ContactListTag, table.ContactListTag.ContactListId.EQ(table.ContactList.UniqueId)).
				LEFT_JOIN(table.Tag, table.Tag.UniqueId.EQ(table.ContactListTag.TagId))).
		WHERE(table.ContactList.OrganizationId.EQ(String(context.Session.User.OrganizationId))).
		LIMIT(*pageSize).
		OFFSET(*pageNumber * *pageSize)

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

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	jsonLists, _ := json.Marshal(dest)
	context.App.Logger.Info("Lists: %v", jsonLists)

	listsToReturn := []api_types.ContactListSchema{}

	for _, list := range dest.Lists {

		tags := []api_types.TagSchema{}
		if len(list.Tags) > 0 {
			for _, tag := range list.Tags {
				stringUniqueId := tag.UniqueId.String()
				tagToAppend := api_types.TagSchema{
					UniqueId: &stringUniqueId,
					Name:     &tag.Label,
				}

				tags = append(tags, tagToAppend)
			}
		}

		uniqueId := list.UniqueId.String()

		lst := api_types.ContactListSchema{
			CreatedAt:             &list.CreatedAt,
			Name:                  &list.Name,
			Description:           &list.Name,
			NumberOfCampaignsSent: &list.TotalCampaigns,
			NumberOfContacts:      &list.TotalContacts,
			Tags:                  &tags,
			UniqueId:              &uniqueId,
		}
		listsToReturn = append(listsToReturn, lst)
	}

	return context.JSON(http.StatusOK, api_types.GetContactListResponseSchema{
		Lists: &listsToReturn,
		PaginationMeta: &api_types.PaginationMeta{
			Page:    pageNumber,
			PerPage: pageSize,
			Total:   &dest.TotalLists,
		},
	})
}

func CreateNewContactLists(context interfaces.CustomContext) error {
	return nil
}

func GetContactListById(context interfaces.CustomContext) error {

	contactListId := context.Param("id")

	if contactListId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Contact list id is required")
	}

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
			table.ContactList.OrganizationId.EQ(String(context.Session.User.OrganizationId)).
				AND(table.ContactList.UniqueId.EQ(String(contactListId))),
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
				UniqueId: &stringUniqueId,
				Name:     &tag.Label,
			}
			tags = append(tags, tagToAppend)
		}
	}

	uniqueId := dest.UniqueId.String()

	return context.JSON(http.StatusOK, api_types.GetContactListByIdSchema{
		List: &api_types.ContactListSchema{
			CreatedAt:             &dest.CreatedAt,
			Name:                  &dest.Name,
			Description:           &dest.Name,
			NumberOfCampaignsSent: &dest.TotalCampaigns,
			NumberOfContacts:      &dest.TotalContacts,
			Tags:                  &tags,
			UniqueId:              &uniqueId,
		},
	})
}

func DeleteContactListById(context interfaces.CustomContext) error {

	contactListId := context.Param("id")

	if contactListId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Contact list id is required")
	}

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

func UpdateContactListById(context interfaces.CustomContext) error {
	return nil
}
