package campaign_controller

import (
	"encoding/json"
	"fmt"
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

type CampaignController struct {
	controller.BaseController `json:"-,inline"`
}

func NewCampaignController() *CampaignController {
	return &CampaignController{
		BaseController: controller.BaseController{
			Name:        "Campaign Controller",
			RestApiPath: "/api/campaign",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/campaigns",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getCampaigns),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60 * 60, // 1 hour
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetCampaign,
						},
					},
				},
				{
					Path:                    "/api/campaigns",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(createNewCampaign),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60 * 60, // 1 hour
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.CreateCampaign,
						},
					},
				},
				{
					Path:                    "/api/campaigns/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getCampaignById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60 * 60, // 1 hour
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetCampaign,
						},
					},
				},
				{
					Path:                    "/api/campaigns/:id",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(updateCampaignById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60 * 60, // 1 hour
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.UpdateCampaign,
						},
					},
				},
				{
					Path:                    "/api/campaigns/:id",
					Method:                  http.MethodDelete,
					Handler:                 interfaces.HandlerWithSession(deleteCampaignById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60 * 60, // 1 hour
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.DeleteCampaign,
						},
					},
				},
			},
		},
	}
}

func getCampaigns(context interfaces.ContextWithSession) error {

	params := new(api_types.GetCampaignsParams)

	err := utils.BindQueryParams(context, params)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	pageNumber := params.Page
	pageSize := params.PerPage
	order := params.Order
	status := params.Status

	var dest []struct {
		TotalCampaigns int `json:"totalCampaigns"`
		model.Campaign
		Tags []struct {
			model.Tag
		}
		Lists []struct {
			model.ContactList
		}
	}

	orgUuid, _ := uuid.Parse(context.Session.User.OrganizationId)

	whereCondition := table.Campaign.OrganizationId.EQ(UUID(orgUuid))

	campaignQuery := SELECT(
		table.Campaign.AllColumns,
		table.Tag.AllColumns,
		table.CampaignList.AllColumns,
		table.ContactList.AllColumns,
		table.CampaignTag.AllColumns,
		COUNT(table.Campaign.UniqueId).OVER().AS("totalCampaigns"),
	).
		FROM(table.Campaign.
			LEFT_JOIN(table.CampaignTag, table.CampaignTag.CampaignId.EQ(table.Campaign.UniqueId)).
			LEFT_JOIN(table.Tag, table.Tag.UniqueId.EQ(table.CampaignTag.TagId)).
			LEFT_JOIN(table.CampaignList, table.CampaignList.CampaignId.EQ(table.Campaign.UniqueId)).
			LEFT_JOIN(table.ContactList, table.ContactList.UniqueId.EQ(table.CampaignList.ContactListId)),
		).
		WHERE(whereCondition).
		LIMIT(pageSize).
		OFFSET((pageNumber - 1) * pageSize)

	if order != nil {
		if *order == api_types.OrderEnum(api_types.Asc) {
			campaignQuery.ORDER_BY(table.Campaign.CreatedAt.ASC())
		} else {
			campaignQuery.ORDER_BY(table.Campaign.CreatedAt.DESC())
		}
	}

	if status != nil {
		statusToFilterWith := model.CampaignStatusEnum(*status)
		whereCondition.AND(table.Campaign.Status.EQ(String(statusToFilterWith.String())))
	}

	err = campaignQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			total := 0
			campaigns := make([]api_types.CampaignSchema, 0)
			return context.JSON(http.StatusOK, api_types.GetCampaignResponseSchema{
				Campaigns: campaigns,
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

	fmt.Println("Campaigns: ", dest)

	campaignsToReturn := []api_types.CampaignSchema{}

	if len(dest) > 0 {
		for _, campaign := range dest {
			tags := []api_types.TagSchema{}
			lists := []api_types.ContactListSchema{}
			status := api_types.CampaignStatusEnum(campaign.Status)
			isLinkTrackingEnabled := campaign.IsLinkTrackingEnabled

			if len(campaign.Tags) > 0 {
				for _, tag := range campaign.Tags {
					stringUniqueId := tag.UniqueId.String()
					tagToAppend := api_types.TagSchema{
						UniqueId: stringUniqueId,
						Label:    tag.Label,
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

			// convert string to *map[string]interface{} for template component parameters
			var templateComponentParameters *map[string]interface{}
			if campaign.TemplateMessageComponentParameters != nil {
				var unmarshalled map[string]interface{}
				err := json.Unmarshal([]byte(*campaign.TemplateMessageComponentParameters), &unmarshalled)
				if err != nil {
					context.App.Logger.Error("error unmarshalling template component parameters: %v", err.Error())
				}
				templateComponentParameters = &unmarshalled
			}

			cmpgn := api_types.CampaignSchema{
				CreatedAt:                   campaign.CreatedAt,
				Name:                        campaign.Name,
				Description:                 campaign.Description,
				IsLinkTrackingEnabled:       isLinkTrackingEnabled,
				TemplateMessageId:           campaign.MessageTemplateId,
				Status:                      status,
				Lists:                       lists,
				Tags:                        tags,
				SentAt:                      nil,
				UniqueId:                    campaign.UniqueId.String(),
				PhoneNumberInUse:            &campaign.PhoneNumber,
				TemplateComponentParameters: templateComponentParameters,
				ScheduledAt:                 campaign.ScheduledAt,
			}
			campaignsToReturn = append(campaignsToReturn, cmpgn)
		}
	}

	totalCampaigns := 0

	if len(dest) > 0 {
		totalCampaigns = dest[0].TotalCampaigns
	}

	return context.JSON(http.StatusOK, api_types.GetCampaignResponseSchema{
		Campaigns: campaignsToReturn,
		PaginationMeta: api_types.PaginationMeta{
			Page:    pageNumber,
			PerPage: pageSize,
			Total:   totalCampaigns,
		},
	})
}

func createNewCampaign(context interfaces.ContextWithSession) error {
	payload := new(api_types.CreateCampaignJSONRequestBody)
	if err := context.Bind(payload); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	// create new campaign
	organizationUuid, err := uuid.Parse(context.Session.User.OrganizationId)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	userUuid, err := uuid.Parse(context.Session.User.UniqueId)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	orgMemberQuery := SELECT(table.OrganizationMember.AllColumns).
		FROM(table.OrganizationMember).
		WHERE(table.OrganizationMember.UserId.EQ(UUID(userUuid)).AND(
			table.OrganizationMember.OrganizationId.EQ(UUID(organizationUuid)),
		)).LIMIT(1)

	var orgMember model.OrganizationMember

	err = orgMemberQuery.QueryContext(context.Request().Context(), context.App.Db, &orgMember)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	var newCampaign model.Campaign
	tx, err := context.App.Db.BeginTx(context.Request().Context(), nil)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()
	// 1. Insert Campaign
	insertQuery := table.Campaign.INSERT(table.Campaign.MutableColumns).
		MODEL(model.Campaign{
			Name:                          payload.Name,
			Description:                   payload.Description,
			Status:                        model.CampaignStatusEnum_Draft,
			OrganizationId:                organizationUuid,
			MessageTemplateId:             &payload.TemplateMessageId,
			PhoneNumber:                   payload.PhoneNumberToUse,
			IsLinkTrackingEnabled:         payload.IsLinkTrackingEnabled,
			CreatedByOrganizationMemberId: orgMember.UniqueId,
			CreatedAt:                     time.Now(),
			UpdatedAt:                     time.Now(),
			// ScheduledAt:                   payload.ScheduledAt,
		}).RETURNING(table.Campaign.AllColumns)

	debugSql := insertQuery.DebugSql()
	context.App.Logger.Debug("Debug SQL: %v", debugSql)

	err = insertQuery.QueryContext(context.Request().Context(), tx, &newCampaign)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	// 2. Insert Campaign Tags (if any)
	if len(payload.Tags) > 0 {
		campaignTags := make([]model.CampaignTag, 0)
		for _, payloadTag := range payload.Tags {
			tagUUID, err := uuid.Parse(payloadTag)
			if err != nil {
				context.App.Logger.Error("Error converting tag unique id to uuid: %v", err)
				continue
			}
			campaignTags = append(campaignTags, model.CampaignTag{
				CampaignId: newCampaign.UniqueId, // Use the inserted campaign ID
				TagId:      tagUUID,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			})
		}

		_, err := table.CampaignTag.INSERT().
			MODELS(campaignTags).ExecContext(context.Request().Context(), tx)

		if err != nil {
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	fmt.Println("Payload List Ids: ", payload.ListIds)

	var campaignList []model.CampaignList

	// 3. Insert Campaign Lists (if any)
	if len(payload.ListIds) > 0 {
		campaignLists := make([]model.CampaignList, 0)
		for _, listId := range payload.ListIds {
			listUUID, err := uuid.Parse(listId)
			fmt.Println("List UUID: ", listUUID)
			if err != nil {
				context.App.Logger.Error("Error converting list unique id to uuid: %v", err)
				continue
			}
			campaignLists = append(campaignLists, model.CampaignList{
				CampaignId:    newCampaign.UniqueId, // Use the inserted campaign ID
				ContactListId: listUUID,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			})
		}

		campaignListQuery := table.CampaignList.INSERT().
			MODELS(campaignLists).
			RETURNING(table.CampaignList.AllColumns)

		err = campaignListQuery.QueryContext(context.Request().Context(), tx, &campaignList)

		if err != nil {
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	err = tx.Commit()

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	responseToReturn := api_types.CreateNewCampaignResponseSchema{
		Campaign: api_types.CampaignSchema{
			CreatedAt:             newCampaign.CreatedAt,
			UniqueId:              newCampaign.UniqueId.String(),
			Name:                  newCampaign.Name,
			Description:           newCampaign.Description,
			IsLinkTrackingEnabled: newCampaign.IsLinkTrackingEnabled,
			TemplateMessageId:     newCampaign.MessageTemplateId,
			Status:                api_types.CampaignStatusEnum(newCampaign.Status),
			Lists:                 []api_types.ContactListSchema{},
			Tags:                  []api_types.TagSchema{},
			SentAt:                nil,
		},
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func getCampaignById(context interfaces.ContextWithSession) error {
	campaignId := context.Param("id")
	if campaignId == "" {
		return context.JSON(http.StatusBadRequest, "Invalid Campaign Id")
	}

	campaignUuid, _ := uuid.Parse(campaignId)

	sqlStatement := SELECT(
		table.Campaign.AllColumns,
		table.Tag.AllColumns,
		table.CampaignList.AllColumns,
		table.ContactList.AllColumns,
		table.CampaignTag.AllColumns,
	).
		FROM(table.Campaign.
			LEFT_JOIN(table.CampaignTag, table.CampaignTag.CampaignId.EQ(UUID(campaignUuid))).
			LEFT_JOIN(table.Tag, table.Tag.UniqueId.EQ(table.CampaignTag.TagId)).
			LEFT_JOIN(table.CampaignList, table.CampaignList.CampaignId.EQ(UUID(campaignUuid))).
			LEFT_JOIN(table.ContactList, table.ContactList.UniqueId.EQ(table.CampaignList.ContactListId))).
		WHERE(
			table.Campaign.UniqueId.EQ(UUID(campaignUuid)),
		)

	var campaignResponse struct {
		model.Campaign
		Tags  []model.Tag
		Lists []model.ContactList
	}

	sqlStatement.Query(context.App.Db, &campaignResponse)

	if campaignResponse.UniqueId.String() == "" {
		return context.JSON(http.StatusNotFound, "Campaign not found")
	}

	status := api_types.CampaignStatusEnum(campaignResponse.Status)
	isLinkTrackingEnabled := campaignResponse.IsLinkTrackingEnabled

	stringUniqueId := campaignResponse.UniqueId.String()

	// convert string to *map[string]interface{} for template component parameters
	var templateComponentParameters *map[string]interface{}
	if campaignResponse.TemplateMessageComponentParameters != nil {
		var unmarshalled map[string]interface{}
		err := json.Unmarshal([]byte(*campaignResponse.TemplateMessageComponentParameters), &unmarshalled)
		if err != nil {
			context.App.Logger.Error("error unmarshalling template component parameters: %v", err.Error())
		}
		templateComponentParameters = &unmarshalled
	}

	tags := []api_types.TagSchema{}
	lists := []api_types.ContactListSchema{}

	if len(campaignResponse.Tags) > 0 {
		for _, tag := range campaignResponse.Tags {
			stringUniqueId := tag.UniqueId.String()
			tagToAppend := api_types.TagSchema{
				UniqueId: stringUniqueId,
				Label:    tag.Label,
			}

			tags = append(tags, tagToAppend)
		}
	}

	if len(campaignResponse.Lists) > 0 {
		for _, list := range campaignResponse.Lists {
			stringUniqueId := list.UniqueId.String()
			listToAppend := api_types.ContactListSchema{
				UniqueId: stringUniqueId,
				Name:     list.Name,
			}

			lists = append(lists, listToAppend)
		}
	}

	return context.JSON(http.StatusOK, api_types.GetCampaignByIdResponseSchema{
		Campaign: api_types.CampaignSchema{
			CreatedAt:                   campaignResponse.CreatedAt,
			UniqueId:                    stringUniqueId,
			Name:                        campaignResponse.Name,
			Description:                 campaignResponse.Description,
			IsLinkTrackingEnabled:       isLinkTrackingEnabled,
			TemplateMessageId:           campaignResponse.MessageTemplateId,
			PhoneNumberInUse:            &campaignResponse.PhoneNumber,
			Status:                      status,
			Lists:                       lists,
			Tags:                        tags,
			SentAt:                      nil,
			TemplateComponentParameters: templateComponentParameters,
		},
	})
}

func updateCampaignById(context interfaces.ContextWithSession) error {
	logger := context.App.Logger
	campaignId := context.Param("id")
	if campaignId == "" {
		return context.JSON(http.StatusBadRequest, "Invalid Campaign Id")
	}
	payload := new(api_types.UpdateCampaignByIdJSONRequestBody)
	if err := context.Bind(payload); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	fmt.Println("Payload: ", payload)

	orgUuid, _ := uuid.Parse(context.Session.User.OrganizationId)
	campaignUuid, _ := uuid.Parse(campaignId)

	var campaign struct {
		model.Campaign
		Tags  []model.Tag
		Lists []model.ContactList
	}

	campaignQuery := SELECT(
		table.Campaign.AllColumns,
		table.Tag.AllColumns,
		table.ContactList.AllColumns,
		table.CampaignTag.AllColumns,
		table.CampaignList.AllColumns).
		FROM(table.Campaign.
			LEFT_JOIN(table.CampaignTag, table.CampaignTag.CampaignId.EQ(UUID(campaignUuid))).
			LEFT_JOIN(table.Tag, table.Tag.UniqueId.EQ(table.CampaignTag.TagId)).
			LEFT_JOIN(table.CampaignList, table.CampaignList.CampaignId.EQ(UUID(campaignUuid))).
			LEFT_JOIN(table.ContactList, table.ContactList.UniqueId.EQ(table.CampaignList.ContactListId))).
		WHERE(
			table.Campaign.OrganizationId.EQ(UUID(orgUuid)).AND(
				table.Campaign.UniqueId.EQ(UUID(campaignUuid)),
			),
		)

	err := campaignQuery.QueryContext(context.Request().Context(), context.App.Db, &campaign)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "Campaign not found")
		}
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	// ! if this is a status update, handle it first and return
	if campaign.Status != model.CampaignStatusEnum(*payload.Status) {
		// * this is a status update

		updateStatusQuery :=
			table.Campaign.UPDATE(table.Campaign.Status).
				WHERE(table.Campaign.UniqueId.EQ(UUID(campaignUuid)))

		if *payload.Status == api_types.Finished {
			return context.JSON(http.StatusBadRequest, "user can not finish a campaign, but can cancel it.")
		}

		if *payload.Status == api_types.Running {
			isLimitReachedForActiveCampaigns := context.IsActiveCampaignLimitReached()
			if isLimitReachedForActiveCampaigns {
				return context.JSON(http.StatusBadRequest, "Upgrade to run more campaigns concurrently")
			}

			updateStatusQuery.SET(table.Campaign.Status.SET(utils.EnumExpression(model.CampaignStatusEnum_Running.String())))
			_, err := updateStatusQuery.ExecContext(context.Request().Context(), context.App.Db)
			if err != nil {
				return context.JSON(http.StatusInternalServerError, err.Error())
			}

		} else if *payload.Status == api_types.Paused || *payload.Status == api_types.Cancelled {
			if campaign.Status != model.CampaignStatusEnum_Running {
				return context.JSON(http.StatusBadRequest, "Cannot pause a campaign that is not running")
			}

			updateStatusQuery.SET(table.Campaign.Status.SET(utils.EnumExpression(model.CampaignStatusEnum_Paused.String())))
			_, err := updateStatusQuery.ExecContext(context.Request().Context(), context.App.Db)
			if err != nil {
				return context.JSON(http.StatusInternalServerError, err.Error())
			}

			context.App.CampaignManager.StopCampaign(campaign.UniqueId.String())
		}

		return context.JSON(http.StatusOK, api_types.UpdateCampaignByIdResponseSchema{
			IsUpdated: true,
		})
	}

	if campaign.Status == model.CampaignStatusEnum_Finished || campaign.Status == model.CampaignStatusEnum_Cancelled {
		return context.JSON(http.StatusBadRequest, "Cannot update a finished campaign")
	}

	if campaign.Status == model.CampaignStatusEnum_Running {
		return context.JSON(http.StatusBadRequest, "Cannot update a running campaign, pause the campaign first to update")
	}

	// * ====== SYNC TAGS FOR THIS CAMPAIGN ======

	oldTagsUuids := make([]uuid.UUID, 0)
	newTagsUuids := make([]uuid.UUID, 0)

	for _, tag := range campaign.Tags {
		oldTagsUuids = append(oldTagsUuids, tag.UniqueId)
	}

	for _, tagId := range payload.Tags {
		tagUuid, err := uuid.Parse(tagId)
		if err != nil {
			continue
		}
		newTagsUuids = append(newTagsUuids, tagUuid)
	}

	tagsToBeDeleted := make([]Expression, 0)
	tagsToBeInserted := make([]model.CampaignTag, 0)

	commonTagIds := make([]uuid.UUID, 0)

	// * the tags ids that are in oldTagsUuids but not in newTagsUuids are needed to be deleted
	for _, oldList := range oldTagsUuids {
		found := false
		for _, newList := range newTagsUuids {
			if oldList == newList {
				found = true
				commonTagIds = append(commonTagIds, oldList)
				break
			}
		}
		if !found {
			tagsToBeDeleted = append(tagsToBeDeleted, UUID(oldList))
		}
	}

	// * the tag ids that are in newTagsUuids but not in oldTagsUuids are needed to be inserted
	for _, newTag := range newTagsUuids {
		found := false
		for _, oldList := range oldTagsUuids {
			if newTag == oldList {
				found = true
				commonTagIds = append(commonTagIds, newTag)
				break
			}
		}

		if !found {
			campaignList := model.CampaignTag{
				CampaignId: campaignUuid,
				TagId:      newTag,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			tagsToBeInserted = append(tagsToBeInserted, campaignList)
		}
	}

	if len(tagsToBeDeleted) > 0 {
		deleteQuery := table.CampaignTag.
			DELETE().
			WHERE(table.CampaignTag.CampaignId.EQ(UUID(campaignUuid)).
				AND(table.CampaignTag.TagId.IN(tagsToBeDeleted...)))

		_, err = deleteQuery.ExecContext(context.Request().Context(), context.App.Db)

		if err != nil {
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	var insertedTags []model.CampaignTag

	if len(tagsToBeInserted) > 0 {
		tagToBeInsertedExpression := make([]Expression, 0)
		for _, tag := range tagsToBeInserted {
			tagToBeInsertedExpression = append(tagToBeInsertedExpression, UUID(tag.TagId))
		}

		tagToBeInsertedCte := CTE("tags_to_be_inserted")

		campaignTagInsertQuery := WITH(
			tagToBeInsertedCte.AS(
				SELECT(table.Tag.AllColumns).FROM(
					table.Tag,
				).WHERE(
					table.Tag.UniqueId.IN(tagToBeInsertedExpression...),
				),
			),
			CTE("insert_tag").AS(
				table.CampaignTag.
					INSERT(table.CampaignTag.MutableColumns).
					MODELS(tagsToBeInserted).
					ON_CONFLICT(table.CampaignTag.CampaignId, table.CampaignTag.TagId).
					DO_NOTHING(),
			),
		)(
			SELECT(tagToBeInsertedCte.AllColumns()).FROM(tagToBeInsertedCte),
		)

		err = campaignTagInsertQuery.QueryContext(context.Request().Context(), context.App.Db, &insertedTags)

		if err != nil {
			logger.Error("Error inserting tags:", err.Error(), nil)
			return context.JSON(http.StatusInternalServerError, err.Error())
		}

	}

	// * ====== SYNC LISTS FOR THIS CAMPAIGN ======

	oldListsUuids := make([]uuid.UUID, 0)
	newListsUuids := make([]uuid.UUID, 0)

	for _, list := range campaign.Lists {
		oldListsUuids = append(oldListsUuids, list.UniqueId)
	}

	for _, listId := range payload.ListIds {
		listUuid, err := uuid.Parse(listId)
		if err != nil {
			continue
		}
		newListsUuids = append(newListsUuids, listUuid)
	}

	listsToBeDeleted := make([]Expression, 0)
	listsToBeInserted := make([]model.CampaignList, 0)

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
			campaignList := model.CampaignList{
				CampaignId:    campaignUuid,
				ContactListId: newList,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			listsToBeInserted = append(listsToBeInserted, campaignList)
		}
	}

	if len(listsToBeDeleted) > 0 {
		deleteQuery := table.CampaignList.
			DELETE().
			WHERE(table.CampaignList.CampaignId.EQ(UUID(campaignUuid)).
				AND(table.CampaignList.ContactListId.IN(listsToBeDeleted...)))

		_, err = deleteQuery.ExecContext(context.Request().Context(), context.App.Db)

		if err != nil {
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	var insertedLists []model.CampaignList

	if len(listsToBeInserted) > 0 {
		listToBeInsertedExpression := make([]Expression, 0)
		for _, list := range listsToBeInserted {
			listToBeInsertedExpression = append(listToBeInsertedExpression, UUID(list.ContactListId))
		}
		listToBeInsertedCte := CTE("lists_to_be_inserted")
		campaignListInsertQuery := WITH(
			listToBeInsertedCte.AS(
				SELECT(table.ContactList.AllColumns).FROM(
					table.ContactList,
				).WHERE(
					table.ContactList.UniqueId.IN(listToBeInsertedExpression...),
				),
			),
			CTE("insert_list").AS(
				table.CampaignList.
					INSERT().
					MODELS(listsToBeInserted).
					ON_CONFLICT(table.CampaignList.CampaignId, table.CampaignList.ContactListId).
					DO_NOTHING(),
			),
		)(
			SELECT(listToBeInsertedCte.AllColumns()).FROM(listToBeInsertedCte),
		)

		err = campaignListInsertQuery.QueryContext(context.Request().Context(), context.App.Db, &insertedLists)

		if err != nil {
			logger.Error("Error inserting lists:", err.Error(), nil)
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	// * use default = {} if no parameters are provided
	var stringifiedParameters []byte
	stringifiedParameters, err = json.Marshal(payload.TemplateComponentParameters)
	if err != nil {
		context.App.Logger.Error("Error marshalling template component parameters: %v", err.Error())
	}

	// pitch in default if no parameters are provided
	if stringifiedParameters == nil {
		stringifiedParameters = []byte("{}")
	}

	finalParameters := string(stringifiedParameters)

	campaignUpdateQuery := table.Campaign.UPDATE(table.Campaign.MutableColumns).
		MODEL(model.Campaign{
			Name:                               payload.Name,
			Description:                        payload.Description,
			MessageTemplateId:                  payload.TemplateMessageId,
			PhoneNumber:                        *payload.PhoneNumber,
			IsLinkTrackingEnabled:              payload.EnableLinkTracking,
			UpdatedAt:                          time.Now(),
			CreatedAt:                          campaign.CreatedAt,
			Status:                             model.CampaignStatusEnum(*payload.Status),
			OrganizationId:                     orgUuid,
			CreatedByOrganizationMemberId:      campaign.CreatedByOrganizationMemberId,
			TemplateMessageComponentParameters: &finalParameters,
			ScheduledAt:                        payload.ScheduledAt,
		}).
		WHERE(table.Campaign.UniqueId.EQ(UUID(campaignUuid))).
		RETURNING(table.Campaign.AllColumns)

	var updatedCampaign model.Campaign

	err = campaignUpdateQuery.QueryContext(context.Request().Context(), context.App.Db, &updatedCampaign)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	response := api_types.UpdateCampaignByIdResponseSchema{
		IsUpdated: true,
	}

	return context.JSON(http.StatusOK, response)
}

func deleteCampaignById(context interfaces.ContextWithSession) error {
	campaignId := context.Param("id")
	if campaignId == "" {
		return context.JSON(http.StatusBadRequest, "Invalid Campaign Id")
	}

	orgUuid, _ := uuid.Parse(context.Session.User.OrganizationId)
	campaignUuid, _ := uuid.Parse(campaignId)
	var campaign model.Campaign
	campaignQuery := SELECT(table.Campaign.AllColumns).FROM(table.Campaign).
		WHERE(
			table.Campaign.UniqueId.EQ(UUID(campaignUuid)).
				AND(table.Campaign.OrganizationId.EQ(UUID(orgUuid))))
	err := campaignQuery.QueryContext(context.Request().Context(), context.App.Db, &campaign)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "Campaign not found")
		}
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	if campaign.Status == model.CampaignStatusEnum_Running {
		return context.JSON(http.StatusBadRequest, "Cannot delete a running campaign, pause the campaign first to delete")
	}

	result, err := table.Campaign.DELETE().WHERE(table.Campaign.UniqueId.EQ(String(campaignId))).ExecContext(context.Request().Context(), context.App.Db)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	if res, _ := result.RowsAffected(); res == 0 {
		return context.JSON(http.StatusNotFound, "Campaign not found")
	}

	return context.JSON(http.StatusOK, "OK")
}
