package contact_list_controller

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wapikit/wapikit/api/api_types"
	controller "github.com/wapikit/wapikit/api/controllers"
	"github.com/wapikit/wapikit/interfaces"
	"github.com/wapikit/wapikit/utils"

	"github.com/go-jet/jet/qrm"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
)

type ContactListController struct {
	controller.BaseController `json:"-,inline"`
}

func NewContactListController() *ContactListController {
	return &ContactListController{
		BaseController: controller.BaseController{
			Name:        "Contact List Controller",
			RestApiPath: "/api/contact-list",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/lists",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(GetContactLists),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetList,
						},
					},
				},
				{
					Path:                    "/api/lists",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(CreateNewContactLists),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.CreateList,
						},
					},
				},
				{
					Path:                    "/api/lists/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(GetContactListById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetList,
						},
					},
				},
				{
					Path:                    "/api/contacts/:id",
					Method:                  http.MethodDelete,
					Handler:                 interfaces.HandlerWithSession(DeleteContactListById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.DeleteList,
						},
					},
				},
				{
					Path:                    "/api/contacts/:id",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(UpdateContactListById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.UpdateList,
						},
					},
				},
			},
		},
	}
}

func GetContactLists(context interfaces.ContextWithSession) error {
	params := new(api_types.GetContactListsParams)

	if err := utils.BindQueryParams(context, params); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
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
		COUNT(DISTINCT(table.ContactListContact.ContactId)).AS("totalContacts"),
		COUNT(DISTINCT(table.CampaignList.CampaignId)).
			AS("totalCampaigns"),
	).
		FROM(
			table.ContactList.
				LEFT_JOIN(table.ContactListTag, table.ContactListTag.ContactListId.EQ(table.ContactList.UniqueId)).
				LEFT_JOIN(table.Tag, table.Tag.UniqueId.EQ(table.ContactListTag.TagId)).
				LEFT_JOIN(table.ContactListContact, table.ContactListContact.ContactListId.EQ(table.ContactList.UniqueId)).
				LEFT_JOIN(table.CampaignList, table.CampaignList.ContactListId.EQ(table.ContactList.UniqueId)),
		).
		WHERE(whereCondition).
		GROUP_BY(
			table.ContactList.UniqueId,
			table.Tag.UniqueId,
		).
		LIMIT(pageSize).
		OFFSET((pageNumber - 1) * pageSize)

	debugSql := listsQuery.DebugSql()
	context.App.Logger.Debug(debugSql)

	if order != nil {
		if *order == api_types.Asc {
			listsQuery.ORDER_BY(table.ContactList.CreatedAt.ASC())
		} else {
			listsQuery.ORDER_BY(table.ContactList.CreatedAt.DESC())
		}
	}

	var dest []struct {
		model.ContactList
		TotalContacts  int `json:"totalContacts"`
		TotalLists     int `json:"totalLists"`
		TotalCampaigns int `json:"totalCampaigns"`
		Tags           []struct {
			model.Tag
		}
	}

	err := listsQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
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
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	listsToReturn := []api_types.ContactListSchema{}

	numberOfTotalLists := 0
	if len(dest) > 0 {
		numberOfTotalLists = dest[0].TotalLists
	}

	if len(dest) > 0 {
		for _, list := range dest {
			tags := []api_types.TagSchema{}
			if len(list.Tags) > 0 {
				for _, tag := range list.Tags {
					stringUniqueId := tag.UniqueId.String()
					tagToAppend := api_types.TagSchema{
						UniqueId: stringUniqueId,
						Label:    tag.Label,
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
			Total:   numberOfTotalLists,
		},
	})
}

func CreateNewContactLists(context interfaces.ContextWithSession) error {

	payload := new(api_types.NewContactListSchema)

	if err := context.Bind(payload); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	orgUuid, _ := uuid.Parse(context.Session.User.OrganizationId)

	var contactList = model.ContactList{
		Name: payload.Name,
		// Description:    payload.Description,
		OrganizationId: orgUuid,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	insertQuery := table.ContactList.
		INSERT(table.ContactList.MutableColumns).MODEL(contactList).
		RETURNING(table.ContactList.AllColumns)

	var dest model.ContactList

	err := insertQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	uniqueId := dest.UniqueId.String()

	return context.JSON(http.StatusCreated, api_types.CreateNewListResponseSchema{
		List: api_types.ContactListSchema{
			CreatedAt:   dest.CreatedAt,
			Name:        dest.Name,
			Description: dest.Name,
			UniqueId:    uniqueId,
		},
	})
}

func GetContactListById(context interfaces.ContextWithSession) error {
	contactListId := context.Param("id")
	if contactListId == "" {
		return context.JSON(http.StatusBadRequest, "Contact list id is required")
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
		return context.JSON(http.StatusInternalServerError, "Error fetching list")
	}

	tags := []api_types.TagSchema{}

	if len(dest.Tags) > 0 {
		for _, tag := range dest.Tags {
			stringUniqueId := tag.UniqueId.String()
			tagToAppend := api_types.TagSchema{
				UniqueId: stringUniqueId,
				Label:    tag.Label,
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
		return context.JSON(http.StatusBadRequest, "Contact list id is required")
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
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	if res, _ := result.RowsAffected(); res == 0 {
		return context.JSON(http.StatusNotFound, "List not found")
	}

	return context.JSON(http.StatusOK, "OK")
}

func UpdateContactListById(context interfaces.ContextWithSession) error {
	contactListId := context.Param("id")
	contactListUuid, err := uuid.Parse(contactListId)

	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid contact list id")
	}

	orgUUid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid contact list id")
	}

	payload := new(api_types.UpdateContactListSchema)

	if err := context.Bind(payload); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	var contactList model.ContactList

	contactListQuery := SELECT(table.ContactList.AllColumns).
		FROM(table.ContactList).
		WHERE(table.ContactList.UniqueId.EQ(UUID(contactListUuid)).
			AND(table.ContactList.OrganizationId.EQ(UUID(orgUUid))))

	err = contactListQuery.QueryContext(context.Request().Context(), context.App.Db, &contactList)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "Contact list not found")
		} else {
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	if contactList.UniqueId.String() == "" {
		return context.JSON(http.StatusNotFound, "Contact list not found")
	}

	updateQuery := table.ContactList.
		UPDATE().
		SET(table.ContactList.Name.SET(String(payload.Name))).
		WHERE(
			table.ContactList.UniqueId.EQ(UUID(contactListUuid)).
				AND(table.ContactList.OrganizationId.EQ(UUID(orgUUid))),
		).RETURNING(table.ContactList.AllColumns)

	_, err = updateQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	response := api_types.UpdateListByIdResponseSchema{
		List: api_types.ContactListSchema{
			CreatedAt:             contactList.CreatedAt,
			Name:                  payload.Name,
			Description:           payload.Name,
			NumberOfCampaignsSent: 0,
			NumberOfContacts:      0,
			Tags:                  []api_types.TagSchema{},
			UniqueId:              contactList.UniqueId.String(),
		},
	}

	// ! TODO: add tags here to

	return context.JSON(http.StatusOK, response)
}
