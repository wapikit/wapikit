package campaign_service

import (
	"net/http"

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

type CampaignService struct {
	services.BaseService `json:"-,inline"`
}

func NewCampaignService() *CampaignService {
	return &CampaignService{
		BaseService: services.BaseService{
			Name:        "Campaign Service",
			RestApiPath: "/api/campaign",
			Routes: []interfaces.Route{
				{
					Path:                    "/campaigns",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(GetCampaigns),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Admin,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60 * 60, // 1 hour
						},
					},
				},
				{
					Path:                    "/campaigns",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(CreateNewCampaign),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Admin,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60 * 60, // 1 hour
						},
					},
				},
				{
					Path:                    "/campaigns/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(GetCampaignById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Admin,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60 * 60, // 1 hour
						},
					},
				},
				{
					Path:                    "/campaigns/:id",
					Method:                  http.MethodPut,
					Handler:                 interfaces.HandlerWithSession(UpdateCampaignById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Admin,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60 * 60, // 1 hour
						},
					},
				},
				{
					Path:                    "/campaigns/:id",
					Method:                  http.MethodDelete,
					Handler:                 interfaces.HandlerWithSession(DeleteCampaignById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Admin,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60 * 60, // 1 hour
						},
					},
				},
			},
		},
	}
}

func GetCampaigns(context interfaces.ContextWithSession) error {

	params := new(api_types.GetCampaignsParams)

	err := internal.BindQueryParams(context, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	pageNumber := params.Page
	pageSize := params.PerPage
	order := params.Order
	status := params.Status

	var dest struct {
		TotalCampaigns int `json:"totalCampaigns"`
		Campaigns      []struct {
			model.Campaign
			Tags []struct {
				model.Tag
			}
			Lists []struct {
				model.ContactList
			}
		}
	}

	orgUuid, _ := uuid.Parse(context.Session.User.OrganizationId)

	whereCondition := table.Campaign.OrganizationId.EQ(UUID(orgUuid))

	campaignQuery := SELECT(
		table.Campaign.AllColumns,
		table.Tag.AllColumns,
		table.CampaignList.AllColumns,
		table.CampaignList.AllColumns,
		table.CampaignTag.AllColumns,
		COUNT(table.Campaign.OrganizationId.EQ(String(context.Session.User.OrganizationId))).OVER().AS("totalCampaigns"),
	).
		FROM(table.Campaign.
			LEFT_JOIN(table.CampaignTag, table.CampaignTag.CampaignId.EQ(table.Campaign.UniqueId)).
			LEFT_JOIN(table.Tag, table.Tag.UniqueId.EQ(table.CampaignTag.TagId)).
			LEFT_JOIN(table.CampaignList, table.CampaignList.CampaignId.EQ(table.Campaign.UniqueId)).
			LEFT_JOIN(table.ContactList, table.ContactList.UniqueId.EQ(table.CampaignList.ContactListId)),
		).
		WHERE(whereCondition).
		LIMIT(*pageSize).
		OFFSET(*pageNumber * *pageSize)

	if order != nil {
		if *order == api_types.OrderEnum(api_types.Asc) {
			campaignQuery.ORDER_BY(table.Campaign.CreatedAt.ASC())
		} else {
			campaignQuery.ORDER_BY(table.Campaign.CreatedAt.DESC())
		}
	}

	if status != nil {
		statusToFilterWith := model.CampaignStatus(*status)
		whereCondition.AND(table.Campaign.Status.EQ(String(statusToFilterWith.String())))
	}

	err = campaignQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == "qrm: no rows in result set" {
			total := 0
			campaigns := make([]api_types.CampaignSchema, 0)
			return context.JSON(http.StatusOK, api_types.GetCampaignResponseSchema{
				Campaigns: campaigns,
				PaginationMeta: api_types.PaginationMeta{
					Page:    *pageNumber,
					PerPage: *pageSize,
					Total:   total,
				},
			})
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())

		}
	}

	campaignsToReturn := []api_types.CampaignSchema{}

	if len(dest.Campaigns) > 0 {
		for _, campaign := range dest.Campaigns {
			tags := []api_types.TagSchema{}
			lists := []api_types.ContactListSchema{}
			status := api_types.CampaignStatusEnum(campaign.Status)
			var isLinkTrackingEnabled bool

			if len(campaign.Tags) > 0 {
				for _, tag := range campaign.Tags {
					stringUniqueId := tag.UniqueId.String()
					tagToAppend := api_types.TagSchema{
						UniqueId: stringUniqueId,
						Name:     tag.Label,
					}

					tags = append(tags, tagToAppend)
				}
			}

			if len(campaign.Lists) > 0 {
				for _, list := range campaign.Lists {
					stringUniqueId := list.UniqueId.String()
					listToAppend := api_types.ContactListSchema{
						UniqueId: stringUniqueId,
						Name:     list.Name,
					}

					lists = append(lists, listToAppend)
				}
			}

			cmpgn := api_types.CampaignSchema{
				CreatedAt:             campaign.CreatedAt,
				Name:                  campaign.Name,
				Description:           &campaign.Name,
				IsLinkTrackingEnabled: isLinkTrackingEnabled, // ! TODO: db field check
				TemplateMessageId:     &campaign.MessageTemplateId,
				Status:                status,
				Lists:                 lists,
				Tags:                  tags,
				SentAt:                nil,
			}
			campaignsToReturn = append(campaignsToReturn, cmpgn)
		}
	}

	return context.JSON(http.StatusOK, api_types.GetCampaignResponseSchema{
		Campaigns: campaignsToReturn,
		PaginationMeta: api_types.PaginationMeta{
			Page:    *pageNumber,
			PerPage: *pageSize,
			Total:   dest.TotalCampaigns,
		},
	})
}

func CreateNewCampaign(context interfaces.ContextWithSession) error {
	payload := new(api_types.CreateCampaignJSONRequestBody)
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// create new campaign
	organizationId, err := uuid.Parse(context.Session.User.OrganizationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	userId, err := uuid.Parse(context.Session.User.UniqueId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	var newCampaign model.Campaign
	tx, err := context.App.Db.BeginTx(context.Request().Context(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()
	// 1. Insert Campaign
	err = table.Campaign.INSERT().MODEL(model.Campaign{
		Name:                          payload.Name,
		Status:                        model.CampaignStatus_Draft,
		OrganizationId:                organizationId,
		MessageTemplateId:             *payload.TemplateMessageId,
		CreatedByOrganizationMemberId: userId,
	}).RETURNING(table.Campaign.AllColumns).QueryContext(context.Request().Context(), tx, &newCampaign)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 2. Insert Campaign Tags (if any)
	if len(payload.Tags) > 0 {
		var campaignTags []model.CampaignTag
		for _, payloadTag := range payload.Tags {
			tagUUID, err := uuid.Parse(payloadTag)
			if err != nil {
				context.App.Logger.Error("Error converting tag unique id to uuid: %v", err)
				continue
			}
			campaignTags = append(campaignTags, model.CampaignTag{
				CampaignId: newCampaign.UniqueId, // Use the inserted campaign ID
				TagId:      tagUUID,
			})
		}

		_, err := table.CampaignTag.INSERT().
			MODELS(campaignTags).ExecContext(context.Request().Context(), tx)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	// 3. Insert Campaign Lists (if any)
	if len(payload.ListIds) > 0 {
		var campaignLists []model.CampaignList
		for _, listId := range payload.ListIds {
			listUUID, err := uuid.Parse(listId)
			if err != nil {
				context.App.Logger.Error("Error converting list unique id to uuid: %v", err)
				continue
			}
			campaignLists = append(campaignLists, model.CampaignList{
				CampaignId:    newCampaign.UniqueId, // Use the inserted campaign ID
				ContactListId: listUUID,
			})
		}

		_, err = table.CampaignList.INSERT().
			MODELS(campaignLists).
			ExecContext(context.Request().Context(), tx)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	err = tx.Commit()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.String(http.StatusOK, "OK")
}

func GetCampaignById(context interfaces.ContextWithSession) error {
	campaignId := context.Param("id")
	if campaignId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Campaign Id")
	}

	orgUUid, _ := uuid.Parse(context.Session.User.OrganizationId)
	campaignUuid, _ := uuid.Parse(campaignId)

	sqlStatement := SELECT(table.Campaign.AllColumns, table.Tag.AllColumns).
		FROM(table.Campaign.
			LEFT_JOIN(table.CampaignTag, table.CampaignTag.CampaignId.EQ(String(campaignId))).
			LEFT_JOIN(table.Tag, table.Tag.UniqueId.EQ(table.CampaignTag.TagId))).
		WHERE(AND(
			table.Campaign.OrganizationId.EQ(UUID(orgUUid)),
			table.Campaign.UniqueId.EQ(UUID(campaignUuid)),
		)).LIMIT(1)

	var campaignResponse struct {
		model.Campaign
		Tags []model.Tag
	}
	sqlStatement.Query(database.GetDbInstance(), &campaignResponse)

	if campaignResponse.UniqueId.String() == "" {
		return echo.NewHTTPError(http.StatusNotFound, "Campaign not found")
	}

	status := api_types.CampaignStatusEnum(campaignResponse.Status)
	isLinkTrackingEnabled := false // ! TODO: db field check

	stringUniqueId := campaignResponse.UniqueId.String()

	return context.JSON(http.StatusOK, api_types.GetCampaignByIdResponseSchema{
		Campaign: api_types.CampaignSchema{
			CreatedAt:             campaignResponse.CreatedAt,
			UniqueId:              stringUniqueId,
			Name:                  campaignResponse.Name,
			Description:           &campaignResponse.Name,
			IsLinkTrackingEnabled: isLinkTrackingEnabled, // ! TODO: db field check
			TemplateMessageId:     &campaignResponse.MessageTemplateId,
			Status:                status,
			Lists:                 []api_types.ContactListSchema{},
			Tags:                  []api_types.TagSchema{},
			SentAt:                nil,
		},
	})
}

func UpdateCampaignById(context interfaces.ContextWithSession) error {
	return context.String(http.StatusOK, "OK")
}

func DeleteCampaignById(context interfaces.ContextWithSession) error {
	campaignId := context.Param("id")
	if campaignId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Campaign Id")
	}

	result, err := table.Campaign.DELETE().WHERE(table.Campaign.UniqueId.EQ(String(campaignId))).ExecContext(context.Request().Context(), context.App.Db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if res, _ := result.RowsAffected(); res == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Campaign not found")
	}

	return context.String(http.StatusOK, "OK")
}
