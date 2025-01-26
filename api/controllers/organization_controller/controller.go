package organization_controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	wapi "github.com/wapikit/wapi.go/pkg/client"
	"github.com/wapikit/wapikit/api/api_types"
	controller "github.com/wapikit/wapikit/api/controllers"
	"github.com/wapikit/wapikit/interfaces"
	"github.com/wapikit/wapikit/utils"

	"github.com/go-jet/jet/qrm"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
)

type OrganizationController struct {
	controller.BaseController `json:"-,inline"`
}

func NewOrganizationController() *OrganizationController {
	return &OrganizationController{
		BaseController: controller.BaseController{
			Name:        "Organization Controller",
			RestApiPath: "/api/organization",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/organization",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getOrganizations),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
					},
				},
				{
					Path:                    "/api/organization",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(createNewOrganization),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
					},
				},
				{
					Path:                    "/api/organization/ai-configuration",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getFullAiConfiguration),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RequiredPermission:  []api_types.RolePermissionEnum{},
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
					},
				},
				{
					Path:                    "/api/organization/:id",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(updateOrganizationById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Owner,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.UpdateOrganization,
						},
					},
				},
				{
					Path:                    "/api/organization/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getOrganizationById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
					},
				},
				{
					Path:                    "/api/organization/tags",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getOrganizationTags),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetTag,
						},
					},
				},
				{
					Path:                    "/api/organization/:id/transfer",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(transferOwnershipOfOrganization),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Owner,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
					},
				},
				{
					Path:                    "/api/organization/invites",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getOrganizationInvites),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetOrganizationMember,
						},
					},
				},
				{
					Path:                    "/api/organization/invites",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(createNewOrganizationInvite),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.CreateOrganizationMember,
						},
					},
				},
				{
					Path:                    "/api/organization/members",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getOrganizationMembers),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetOrganizationMember,
						},
					},
				},
				{
					Path:                    "/api/organization/members/:id",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(updateOrgMemberById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.UpdateOrganizationMember,
						},
					},
				},
				{
					Path:                    "/api/organization/members/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getOrgMemberById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetOrganizationMember,
						},
					},
				},
				{
					Path:                    "/api/organization/members/:id",
					Method:                  http.MethodDelete,
					Handler:                 interfaces.HandlerWithSession(deleteOrgMemberById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.DeleteOrganizationMember,
						},
					},
				},
				{
					Path:                    "/api/organization/members/:id/role",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(updateOrganizationMemberRoles),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.UpdateOrganizationMember,
						},
					},
				},
				{
					Path:                    "/api/organization/templates",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getAllMessageTemplates),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetMessageTemplates,
						},
					},
				},
				{
					Path:                    "/api/organization/templates/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getMessageTemplateById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetMessageTemplates,
						},
					},
				},
				{
					Path:                    "/api/organization/phone-numbers",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getAllPhoneNumbers),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetPhoneNumbers,
						},
					},
				},
				{
					Path:                    "/api/organization/phone-numbers/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getPhoneNumberById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetPhoneNumbers,
						},
					},
				},
				{
					Path:                    "/api/organization/whatsappBusinessAccount",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(handleUpdateWhatsappBusinessAccountDetails),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Owner,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
					},
				},
			},
		},
	}
}

func createNewOrganization(context interfaces.ContextWithSession) error {
	payload := new(api_types.NewOrganizationSchema)
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var newOrg model.Organization
	var member model.OrganizationMember

	tx, err := context.App.Db.BeginTx(context.Request().Context(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()

	// 1. Insert Organization
	err = table.Organization.INSERT(table.Organization.MutableColumns).
		MODEL(model.Organization{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Name:        payload.Name,
			Description: payload.Description,
		}).
		RETURNING(table.Organization.AllColumns).
		QueryContext(context.Request().Context(), tx, &newOrg)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	userUuid, err := uuid.Parse(context.Session.User.UniqueId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// 2. Insert Organization Member
	err = table.OrganizationMember.INSERT(table.OrganizationMember.MutableColumns).MODEL(model.OrganizationMember{
		AccessLevel:    model.UserPermissionLevelEnum_Owner,
		OrganizationId: newOrg.UniqueId,
		UserId:         userUuid,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}).RETURNING(table.OrganizationMember.AllColumns).QueryContext(context.Request().Context(), tx, &member)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 3. Create API key for the organization
	claims := &interfaces.JwtPayload{
		ContextUser: interfaces.ContextUser{
			Username:       context.Session.User.Username,
			Email:          context.Session.User.Email,
			Role:           api_types.UserPermissionLevelEnum(api_types.Owner),
			UniqueId:       context.Session.User.UniqueId,
			OrganizationId: newOrg.UniqueId.String(),
			Name:           context.Session.User.Name,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 365 * 2).Unix(), // 60-day expiration
			Issuer:    "wapikit",
		},
	}

	//Create the token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(context.App.Koa.String("app.jwt_secret")))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating token")
	}

	var apiKey model.ApiKey

	err = table.ApiKey.INSERT(table.ApiKey.MutableColumns).MODEL(model.ApiKey{
		MemberId:       member.UniqueId,
		OrganizationId: newOrg.UniqueId,
		Key:            token,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}).RETURNING(table.ApiKey.AllColumns).QueryContext(context.Request().Context(), tx, &apiKey)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	responseToReturn := api_types.CreateNewOrganizationResponseSchema{
		Organization: api_types.OrganizationSchema{
			Name:        newOrg.Name,
			CreatedAt:   newOrg.CreatedAt,
			UniqueId:    newOrg.UniqueId.String(),
			Description: newOrg.Description,
			LogoUrl:     newOrg.LogoUrl,
			FaviconUrl:  &newOrg.FaviconUrl,
			WebsiteUrl:  newOrg.WebsiteUrl,
		},
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func getOrganizations(context interfaces.ContextWithSession) error {
	param := new(api_types.GetUserOrganizationsParams)
	if err := utils.BindQueryParams(context, param); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	userUUid, err := uuid.Parse(context.Session.User.UniqueId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	orgMembers := SELECT(table.OrganizationMember.AllColumns).
		FROM(table.OrganizationMember).
		WHERE(table.OrganizationMember.UserId.EQ(UUID(userUUid))).AsTable("Member")

	orgQuery := SELECT(
		orgMembers.AllColumns(),
		table.Organization.AllColumns,
		COUNT(table.Organization.UniqueId).OVER().AS("total_organizations"),
	).FROM(
		orgMembers.INNER_JOIN(
			table.Organization, table.Organization.UniqueId.EQ(table.OrganizationMember.OrganizationId.From(orgMembers)),
		),
	).
		LIMIT(param.PerPage).
		OFFSET((param.Page - 1) * param.PerPage)

	if param.SortBy != nil {
		if *param.SortBy == api_types.Asc {
			orgQuery.ORDER_BY(table.Organization.CreatedAt.ASC())
		} else {
			orgQuery.ORDER_BY(table.Organization.CreatedAt.DESC())
		}
	}

	var dest []struct {
		model.Organization
		TotalOrganizations int `json:"total_organizations"`
	}

	err = orgQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			var organizations []api_types.OrganizationSchema
			total := 0
			return context.JSON(http.StatusOK, api_types.GetOrganizationsResponseSchema{
				Organizations: organizations,
				PaginationMeta: api_types.PaginationMeta{
					Page:    param.Page,
					PerPage: param.PerPage,
					Total:   total,
				},
			})
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	userOrganizations := []api_types.OrganizationSchema{}
	for _, org := range dest {
		uniqueId := org.UniqueId.String()
		organization := api_types.OrganizationSchema{
			Name:       org.Name,
			CreatedAt:  org.CreatedAt,
			UniqueId:   uniqueId,
			FaviconUrl: &org.FaviconUrl,
			LogoUrl:    org.LogoUrl,
			WebsiteUrl: org.WebsiteUrl,
		}

		if org.SlackChannel != nil && org.SlackWebhookUrl != nil {
			organization.SlackNotificationConfiguration = &api_types.SlackNotificationConfigurationSchema{
				SlackChannel:    *org.SlackChannel,
				SlackWebhookUrl: *org.SlackWebhookUrl,
			}
		}

		if org.SmtpClientHost != nil && org.SmtpClientPassword != nil && org.SmtpClientPort != nil && org.SmtpClientUsername != nil {
			organization.EmailNotificationConfiguration = &api_types.EmailNotificationConfigurationSchema{
				SmtpHost:     *org.SmtpClientHost,
				SmtpPassword: *org.SmtpClientPassword,
				SmtpPort:     *org.SmtpClientPort,
				SmtpUsername: *org.SmtpClientUsername,
			}
		}

		if org.IsAiEnabled {
			model := api_types.AiModelEnum(*org.AiModel)
			organization.AiConfiguration = &api_types.AiConfigurationDetailsSchema{
				IsEnabled: &org.IsAiEnabled,
				Model:     model,
			}
		}

		userOrganizations = append(userOrganizations, organization)
	}

	totalOrganizations := 0

	if len(dest) > 0 {
		totalOrganizations = dest[0].TotalOrganizations
	}

	response := api_types.GetOrganizationsResponseSchema{
		Organizations: userOrganizations,
		PaginationMeta: api_types.PaginationMeta{
			Page:    param.Page,
			PerPage: param.PerPage,
			Total:   totalOrganizations,
		},
	}

	return context.JSON(http.StatusOK, response)
}

func getOrganizationById(context interfaces.ContextWithSession) error {
	organizationId := context.Param("id")
	if organizationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization id")
	}

	orgUuid, err := uuid.Parse(organizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid organization Id")
	}

	userUuid, err := uuid.Parse(context.Session.User.UniqueId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid user Id")
	}

	hasAccess := _verifyAccessToOrganization(context, userUuid, orgUuid)

	if !hasAccess {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this organization")
	}

	var dest model.Organization
	organizationQuery := SELECT(table.Organization.AllColumns).
		FROM(table.Organization).
		WHERE(table.Organization.UniqueId.EQ(UUID(orgUuid))).LIMIT(1)
	err = organizationQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	uniqueId := dest.UniqueId.String()
	orgToReturn := api_types.OrganizationSchema{
		Name:       dest.Name,
		CreatedAt:  dest.CreatedAt,
		UniqueId:   uniqueId,
		FaviconUrl: &dest.FaviconUrl,
		LogoUrl:    dest.LogoUrl,
		WebsiteUrl: dest.WebsiteUrl,
	}

	if dest.SlackChannel != nil && dest.SlackWebhookUrl != nil {
		orgToReturn.SlackNotificationConfiguration = &api_types.SlackNotificationConfigurationSchema{
			SlackChannel:    *dest.SlackChannel,
			SlackWebhookUrl: *dest.SlackWebhookUrl,
		}
	}

	if dest.SmtpClientHost != nil && dest.SmtpClientPassword != nil && dest.SmtpClientPort != nil && dest.SmtpClientUsername != nil {
		orgToReturn.EmailNotificationConfiguration = &api_types.EmailNotificationConfigurationSchema{
			SmtpHost:     *dest.SmtpClientHost,
			SmtpPassword: *dest.SmtpClientPassword,
			SmtpPort:     *dest.SmtpClientPort,
			SmtpUsername: *dest.SmtpClientUsername,
		}
	}

	if dest.IsAiEnabled {
		model := api_types.AiModelEnum(*dest.AiModel)
		orgToReturn.AiConfiguration = &api_types.AiConfigurationDetailsSchema{
			IsEnabled: &dest.IsAiEnabled,
			Model:     model,
		}
	}

	return context.JSON(http.StatusOK, api_types.GetOrganizationByIdResponseSchema{
		Organization: orgToReturn,
	})
}

func deleteOrganization(context interfaces.ContextWithSession) error {

	return context.String(http.StatusInternalServerError, "NOT IMPLEMENTED YET")

	organizationId := context.Param("id")
	if organizationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization id")
	}

	orgUuid, err := uuid.Parse(organizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid organization Id")
	}

	userUuid, err := uuid.Parse(context.Session.User.UniqueId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid user Id")
	}

	hasAccess := _verifyAccessToOrganization(context, userUuid, orgUuid)

	if !hasAccess {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this organization")
	}

	return context.String(http.StatusOK, "OK")
}

func updateOrganizationById(context interfaces.ContextWithSession) error {
	logger := context.App.Logger
	organizationId := context.Param("id")
	if organizationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization id")
	}

	orgUuid, err := uuid.Parse(organizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid organization Id")
	}

	userUuid, err := uuid.Parse(context.Session.User.UniqueId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid user Id")
	}

	hasAccess := _verifyAccessToOrganization(context, userUuid, orgUuid)

	if !hasAccess {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this organization")
	}

	payload := new(api_types.UpdateOrganizationSchema)

	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	orgUpdates := model.Organization{
		Name:        payload.Name,
		UpdatedAt:   time.Now(),
		Description: payload.Description,
	}

	if payload.EmailNotificationConfiguration != nil {
		orgUpdates.SmtpClientHost = &payload.EmailNotificationConfiguration.SmtpHost
		orgUpdates.SmtpClientUsername = &payload.EmailNotificationConfiguration.SmtpUsername
		orgUpdates.SmtpClientPassword = &payload.EmailNotificationConfiguration.SmtpPassword
		orgUpdates.SmtpClientPort = &payload.EmailNotificationConfiguration.SmtpPort
	}

	if payload.SlackNotificationConfiguration != nil {
		orgUpdates.SlackWebhookUrl = &payload.SlackNotificationConfiguration.SlackWebhookUrl
		orgUpdates.SlackChannel = &payload.SlackNotificationConfiguration.SlackChannel
	}

	if payload.AiConfiguration != nil {
		orgUpdates.IsAiEnabled = *payload.AiConfiguration.IsEnabled
		orgUpdates.AiModel = (*model.AiModelEnum)(&payload.AiConfiguration.Model)
		orgUpdates.AiApiKey = payload.AiConfiguration.ApiKey
	}

	var updatedOrg model.Organization

	updateOrgQuery := table.Organization.
		UPDATE(table.Organization.MutableColumns).
		MODEL(orgUpdates).
		WHERE(table.Organization.UniqueId.EQ(UUID(orgUuid))).
		RETURNING(table.Organization.AllColumns)

	err = updateOrgQuery.QueryContext(context.Request().Context(), context.App.Db, &updatedOrg)

	if err != nil {
		logger.Error("Error updating organization", err.Error(), nil)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// if AI chat has been enabled, we have to create a default chat for every user in the organization

	if payload.AiConfiguration != nil && *payload.AiConfiguration.IsEnabled {
		allOrgMembersQuery := SELECT(
			table.OrganizationMember.AllColumns,
			table.AiChat.AllColumns,
		).
			FROM(
				table.OrganizationMember.
					LEFT_JOIN(table.AiChat, table.AiChat.OrganizationMemberId.EQ(table.OrganizationMember.UniqueId)),
			).
			WHERE(table.OrganizationMember.OrganizationId.EQ(UUID(orgUuid)))

		var orgMembers []struct {
			model.OrganizationMember
			AiChat model.AiChat `json:"aiChat"`
		}

		err = allOrgMembersQuery.QueryContext(context.Request().Context(), context.App.Db, &orgMembers)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		aiChatsToCreate := []model.AiChat{}

		for _, member := range orgMembers {
			if member.AiChat.UniqueId == uuid.Nil {
				aiChat := model.AiChat{
					UniqueId:             uuid.New(),
					CreatedAt:            time.Now(),
					UpdatedAt:            time.Now(),
					Status:               model.AiChatStatusEnum_Active,
					OrganizationId:       orgUuid,
					OrganizationMemberId: member.UniqueId,
					Title:                "Default Chat",
					Visibility:           model.AiChatVisibilityEnum_Public,
				}
				aiChatsToCreate = append(aiChatsToCreate, aiChat)
			}
		}

		if len(aiChatsToCreate) > 0 {
			_, err = table.AiChat.INSERT(table.AiChat.AllColumns).
				MODELS(aiChatsToCreate).
				ExecContext(context.Request().Context(), context.App.Db)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong while updating your organization!!")
	}

	return context.JSON(http.StatusOK, api_types.UpdateOrganizationByIdResponseSchema{
		IsUpdated: true,
	})
}

func getOrganizationTags(context interfaces.ContextWithSession) error {
	params := new(api_types.GetOrganizationTagsParams)
	err := utils.BindQueryParams(context, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid query params")
	}

	var dest []struct {
		model.Tag
		TotalTags int `json:"totalTags"`
	}

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error parsing organization UUID")
	}

	whereCondition := table.Tag.OrganizationId.EQ(UUID(orgUuid))

	organizationTagsQuery := SELECT(table.Tag.AllColumns,
		COUNT(table.Tag.UniqueId).OVER().AS("totalTags"),
	).
		FROM(table.Tag).
		WHERE(whereCondition).
		LIMIT(params.PerPage).
		OFFSET((params.Page - 1) * params.PerPage)

	if params.SortBy != nil {
		if *params.SortBy == api_types.Asc {
			organizationTagsQuery.ORDER_BY(table.Tag.CreatedAt.ASC())
		} else {
			organizationTagsQuery.ORDER_BY(table.Tag.CreatedAt.DESC())
		}
	}

	err = organizationTagsQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			var tags []api_types.TagSchema
			total := 0
			return context.JSON(http.StatusOK, api_types.GetOrganizationTagsResponseSchema{
				Tags: tags,
				PaginationMeta: api_types.PaginationMeta{
					Page:    params.Page,
					PerPage: params.PerPage,
					Total:   total,
				},
			})
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	tagsToReturn := []api_types.TagSchema{}

	numberOfTotalTag := 0

	if len(dest) > 0 {
		numberOfTotalTag = dest[0].TotalTags
	}

	if len(dest) > 0 {
		for _, tag := range dest {
			tagId := tag.UniqueId.String()
			tagToReturn := api_types.TagSchema{
				Name:     tag.Label,
				UniqueId: tagId,
			}

			tagsToReturn = append(tagsToReturn, tagToReturn)
		}
	}

	return context.JSON(http.StatusOK, api_types.GetOrganizationTagsResponseSchema{
		Tags: tagsToReturn,
		PaginationMeta: api_types.PaginationMeta{
			Page:    params.Page,
			PerPage: params.PerPage,
			Total:   numberOfTotalTag,
		},
	})
}

func getOrganizationMembers(context interfaces.ContextWithSession) error {
	params := new(api_types.GetOrganizationMembersParams)

	if err := utils.BindQueryParams(context, params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	pageNumber := params.Page
	pageSize := params.PerPage
	sortBy := params.SortBy

	organizationUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	organizationMembersQuery := SELECT(table.OrganizationMember.AllColumns,
		table.User.Username,
		table.User.Name,
		table.User.Email,
		table.RoleAssignment.AllColumns,
		table.OrganizationRole.AllColumns,
		COUNT(table.OrganizationMember.UniqueId).OVER().AS("totalMembers"),
	).
		FROM(table.OrganizationMember.
			LEFT_JOIN(table.User, table.User.UniqueId.EQ(table.OrganizationMember.UserId)).
			LEFT_JOIN(table.RoleAssignment, table.RoleAssignment.OrganizationMemberId.EQ(table.OrganizationMember.UniqueId)).
			LEFT_JOIN(table.OrganizationRole, table.OrganizationRole.UniqueId.EQ(table.RoleAssignment.OrganizationRoleId))).
		WHERE(table.OrganizationMember.OrganizationId.EQ(UUID(organizationUuid))).
		GROUP_BY(
			table.OrganizationMember.UniqueId,
			table.User.Username,
			table.User.Name,
			table.User.Email,
			table.RoleAssignment.UniqueId,
			table.OrganizationRole.UniqueId,
		).
		LIMIT(pageSize).
		OFFSET((pageNumber - 1) * pageSize)

	if sortBy != nil {
		if *sortBy == api_types.Asc {
			organizationMembersQuery.ORDER_BY(table.OrganizationMember.CreatedAt.ASC())
		} else {
			organizationMembersQuery.ORDER_BY(table.OrganizationMember.CreatedAt.ASC())
		}
	}

	var dest []struct {
		model.OrganizationMember
		model.User
		Roles []struct {
			model.OrganizationRole
		}
		TotalMembers int `json:"totalMembers"`
	}

	err = organizationMembersQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			var members []api_types.OrganizationMemberSchema
			total := 0
			return context.JSON(http.StatusOK, api_types.GetOrganizationMembersResponseSchema{
				Members: members,
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

	var membersToReturn []api_types.OrganizationMemberSchema

	if len(dest) > 0 {
		for _, member := range dest {
			var memberRoles []api_types.OrganizationRoleSchema
			if len(member.Roles) > 0 {
				for _, role := range member.Roles {
					var permissions []api_types.RolePermissionEnum
					permissionArray := strings.Split(role.Permissions, ",")

					for _, perm := range permissionArray {
						if perm == "" {
							continue
						}
						permissions = append(permissions, api_types.RolePermissionEnum(perm))
					}

					roleId := role.UniqueId.String()
					roleToReturn := api_types.OrganizationRoleSchema{
						Description: role.Description,
						Name:        role.Name,
						Permissions: permissions,
						UniqueId:    roleId,
					}

					memberRoles = append(memberRoles, roleToReturn)
				}
			}

			accessLevel := api_types.UserPermissionLevelEnum(member.OrganizationMember.AccessLevel)
			memberId := member.OrganizationMember.UniqueId.String()
			mmbr := api_types.OrganizationMemberSchema{
				CreatedAt:   member.OrganizationMember.CreatedAt,
				AccessLevel: accessLevel,
				UniqueId:    memberId,
				Email:       member.User.Email,
				Name:        member.User.Name,
				Roles:       memberRoles,
			}

			membersToReturn = append(membersToReturn, mmbr)
		}
	}

	totalMembers := 0

	if len(dest) > 0 {
		totalMembers = dest[0].TotalMembers
	}

	return context.JSON(http.StatusOK, api_types.GetOrganizationMembersResponseSchema{
		Members: membersToReturn,
		PaginationMeta: api_types.PaginationMeta{
			Page:    pageNumber,
			PerPage: pageSize,
			Total:   totalMembers,
		}})
}

func getOrgMemberById(context interfaces.ContextWithSession) error {
	memberId := context.Param("id")

	if memberId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid member id")
	}

	memberUuid, _ := uuid.Parse(memberId)
	memberQuery := SELECT(table.OrganizationMember.AllColumns,
		table.User.Username,
		table.User.Name,
		table.User.Email,
		table.RoleAssignment.AllColumns,
		table.OrganizationRole.AllColumns,
	).
		FROM(table.OrganizationMember.
			LEFT_JOIN(table.User, table.User.UniqueId.EQ(table.OrganizationMember.UserId)).
			LEFT_JOIN(table.RoleAssignment, table.RoleAssignment.OrganizationMemberId.EQ(table.OrganizationMember.UniqueId)).
			LEFT_JOIN(table.OrganizationRole, table.OrganizationRole.UniqueId.EQ(table.RoleAssignment.OrganizationRoleId))).
		WHERE(table.OrganizationMember.UniqueId.EQ(UUID(memberUuid))).
		LIMIT(1)

	var dest struct {
		member struct {
			model.OrganizationMember
			model.User
			Roles []struct {
				model.OrganizationRole
			}
		}
	}

	err := memberQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			member := new(api_types.OrganizationMemberSchema)
			return context.JSON(http.StatusOK, api_types.GetOrganizationMemberByIdResponseSchema{
				Member: *member,
			})
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	var memberRoles []api_types.OrganizationRoleSchema
	if len(dest.member.Roles) > 0 {
		for _, role := range dest.member.Roles {
			var permissions []api_types.RolePermissionEnum
			for _, perm := range role.Permissions {
				permissions = append(permissions, api_types.RolePermissionEnum(perm))
			}
			roleId := role.UniqueId.String()
			roleToReturn := api_types.OrganizationRoleSchema{
				Description: role.Description,
				Name:        role.Name,
				Permissions: permissions,
				UniqueId:    roleId,
			}
			memberRoles = append(memberRoles, roleToReturn)
		}
	}

	accessLevel := api_types.UserPermissionLevelEnum(dest.member.OrganizationMember.AccessLevel)

	member := api_types.OrganizationMemberSchema{
		CreatedAt:   dest.member.OrganizationMember.CreatedAt,
		AccessLevel: accessLevel,
		UniqueId:    memberId,
		Email:       dest.member.User.Email,
		Name:        dest.member.User.Name,
		Roles:       memberRoles,
	}

	return context.JSON(http.StatusOK, api_types.GetOrganizationMemberByIdResponseSchema{
		Member: member,
	})
}

func deleteOrgMemberById(context interfaces.ContextWithSession) error {
	memberId := context.Param("id")

	if memberId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid member id")
	}

	memberUuid, _ := uuid.Parse(memberId)

	// * delete all role assignments first
	deleteRoleAssignmentQuery := table.RoleAssignment.DELETE().
		WHERE(table.RoleAssignment.OrganizationMemberId.EQ(UUID(memberUuid)))

	_, err := deleteRoleAssignmentQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// * delete the member
	deleteMemberQuery := table.OrganizationMember.DELETE().
		WHERE(table.OrganizationMember.UniqueId.EQ(UUID(memberUuid))).
		RETURNING(table.OrganizationMember.AllColumns)

	_, err = deleteMemberQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := api_types.DeleteOrganizationMemberByIdResponseSchema{
		Data: true,
	}

	return context.JSON(http.StatusOK, response)
}

func updateOrgMemberById(context interfaces.ContextWithSession) error {
	memberId := context.Param("id")
	if memberId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid member id")
	}

	memberUuid, _ := uuid.Parse(memberId)

	payload := new(api_types.UpdateOrganizationMemberSchema)
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	updateMemberQuery := table.OrganizationMember.
		UPDATE(table.OrganizationMember.AccessLevel).
		SET(payload.AccessLevel).
		WHERE(table.OrganizationMember.UniqueId.EQ(UUID(memberUuid))).
		RETURNING(table.OrganizationMember.AllColumns)

	_, err := updateMemberQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.String(http.StatusOK, "OK")
}

func updateOrganizationMemberRoles(context interfaces.ContextWithSession) error {
	memberId := context.Param("id")
	if memberId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid member id")
	}

	memberUuid, _ := uuid.Parse(memberId)
	payload := new(api_types.UpdateOrganizationMemberRoleSchema)
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var orgMember model.OrganizationMember

	memberQuery := SELECT(table.OrganizationMember.AllColumns).
		FROM(table.OrganizationMember).
		WHERE(table.OrganizationMember.UniqueId.EQ(UUID(memberUuid))).
		LIMIT(1)

	err := memberQuery.QueryContext(context.Request().Context(), context.App.Db, &orgMember)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return echo.NewHTTPError(http.StatusNotFound, "Member not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	roleIdExpressions := make([]Expression, 0)

	for _, role := range payload.UpdatedRoleIds {
		roleUuid, err := uuid.Parse(role)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid role id")
		}
		roleIdExpressions = append(roleIdExpressions, UUID(roleUuid))
	}

	var roles []struct {
		model.OrganizationRole
		Assignment model.RoleAssignment
	}

	if len(roleIdExpressions) > 0 {
		orgRoleQuery := SELECT(table.OrganizationRole.AllColumns, table.RoleAssignment.AllColumns).
			FROM(table.OrganizationRole.
				LEFT_JOIN(table.RoleAssignment, table.RoleAssignment.OrganizationRoleId.EQ(table.OrganizationRole.UniqueId)),
			).WHERE(table.OrganizationRole.UniqueId.IN(roleIdExpressions...))

		err := orgRoleQuery.QueryContext(context.Request().Context(), context.App.Db, &roles)

		if err != nil {
			if err.Error() == qrm.ErrNoRows.Error() {
				return echo.NewHTTPError(http.StatusNotFound, "Role not found")
			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}
	}

	// ! delete all roles assignment for which the member has removed the role for

	var removedRolesQuery DeleteStatement

	if len(roleIdExpressions) == 0 {
		removedRolesQuery = table.RoleAssignment.DELETE().WHERE(table.RoleAssignment.OrganizationMemberId.EQ(UUID(memberUuid)))
	} else {
		removedRolesQuery = table.RoleAssignment.DELETE().WHERE(table.RoleAssignment.OrganizationMemberId.EQ(UUID(memberUuid)).AND(
			table.RoleAssignment.OrganizationRoleId.NOT_IN(roleIdExpressions...),
		))
	}

	_, err = removedRolesQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error removing roles")
	}

	// if all roles are removed then return
	if len(payload.UpdatedRoleIds) == 0 {
		responseToReturn := api_types.UpdateOrganizationMemberRoleByIdResponseSchema{
			IsRoleUpdated: true,
		}

		return context.JSON(http.StatusOK, responseToReturn)
	}

	// ! run a up-sert query

	rolesToUpsert := []model.RoleAssignment{}

	for _, role := range roles {
		roleAssignment := model.RoleAssignment{
			OrganizationMemberId: memberUuid,
			OrganizationRoleId:   role.OrganizationRole.UniqueId,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		rolesToUpsert = append(rolesToUpsert, roleAssignment)
	}

	_, err = table.RoleAssignment.INSERT(table.RoleAssignment.MutableColumns).
		MODELS(rolesToUpsert).
		RETURNING(table.RoleAssignment.AllColumns).
		ON_CONFLICT(table.RoleAssignment.OrganizationMemberId, table.RoleAssignment.OrganizationRoleId).
		DO_NOTHING().ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	responseToReturn := api_types.UpdateOrganizationMemberRoleByIdResponseSchema{
		IsRoleUpdated: true,
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func getOrganizationInvites(context interfaces.ContextWithSession) error {
	params := new(api_types.GetOrganizationInvitesParams)
	err := utils.BindQueryParams(context, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var dest struct {
		TotalRoles int `json:"totalRoles"`
		Invites    []struct {
			model.OrganizationMemberInvite
		}
	}

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	organizationInvitesQuery := SELECT(table.OrganizationMemberInvite.AllColumns).
		FROM(table.OrganizationMemberInvite).
		WHERE(table.OrganizationMemberInvite.OrganizationId.EQ(UUID(orgUuid))).
		LIMIT(params.PerPage).
		OFFSET((params.Page - 1) * params.PerPage)

	if params.SortBy != nil {
		if *params.SortBy == api_types.Asc {
			organizationInvitesQuery.ORDER_BY(table.OrganizationMemberInvite.CreatedAt.ASC())
		} else {
			organizationInvitesQuery.ORDER_BY(table.OrganizationMemberInvite.CreatedAt.DESC())
		}
	}

	err = organizationInvitesQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			var invites []api_types.OrganizationMemberInviteSchema
			total := 0
			return context.JSON(http.StatusOK, api_types.GetOrganizationMemberInvitesResponseSchema{
				Invites: invites,
				PaginationMeta: api_types.PaginationMeta{
					Page:    params.Page,
					PerPage: params.PerPage,
					Total:   total,
				},
			})
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	var invitesToReturn []api_types.OrganizationMemberInviteSchema

	if len(dest.Invites) > 0 {
		for _, invite := range dest.Invites {
			accessLevel := api_types.UserPermissionLevelEnum(invite.AccessLevel)
			inviteId := invite.UniqueId.String()
			inv := api_types.OrganizationMemberInviteSchema{
				CreatedAt:   invite.CreatedAt,
				AccessLevel: accessLevel,
				Email:       invite.Email,
				Status:      api_types.InviteStatusEnum(invite.Status),
				UniqueId:    inviteId,
			}
			invitesToReturn = append(invitesToReturn, inv)
		}
	}

	return context.JSON(http.StatusOK, api_types.GetOrganizationMemberInvitesResponseSchema{
		Invites: invitesToReturn,
		PaginationMeta: api_types.PaginationMeta{
			Page:    params.Page,
			PerPage: params.PerPage,
			Total:   dest.TotalRoles,
		},
	})
}

func createNewOrganizationInvite(context interfaces.ContextWithSession) error {
	payload := new(api_types.CreateNewOrganizationInviteSchema)

	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if ok := utils.IsValidEmail(payload.Email); ok == false {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid email address")
	}

	organizationUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var userDest struct {
		model.User
		OrganizationMember       model.OrganizationMember
		OrganizationMemberInvite model.OrganizationMemberInvite
	}

	// check if user exists and is a member of this organization or may have been already sent an invite
	userQuery := SELECT(table.User.AllColumns, table.OrganizationMember.AllColumns, table.OrganizationMemberInvite.AllColumns).
		FROM(table.User.
			LEFT_JOIN(table.OrganizationMember, table.OrganizationMember.UserId.EQ(table.User.UniqueId)).
			LEFT_JOIN(table.OrganizationMemberInvite, table.OrganizationMemberInvite.Email.EQ(String(payload.Email)).
				AND(table.OrganizationMemberInvite.OrganizationId.EQ(UUID(organizationUuid))),
			)).
		WHERE(table.User.Email.EQ(String(payload.Email)).AND(table.OrganizationMember.OrganizationId.EQ(UUID(organizationUuid)))).
		LIMIT(1)

	err = userQuery.QueryContext(context.Request().Context(), context.App.Db, &userDest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			// * user not found create a invite in the db table and also send email to the  user for the invite
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	// * user found, check if the user is already a member of the organization
	if userDest.OrganizationMember.OrganizationId != uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User already a member of the organization")
	}

	// * user found, check if the user has already been sent an invite
	if userDest.OrganizationMemberInvite.OrganizationId != uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User already sent an invite")
	}

	// * user not found create a invite in the db table and also send email to the  user for the invite

	var inviteDest model.OrganizationMemberInvite

	inviteSlug, err := gonanoid.New()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	invite := model.OrganizationMemberInvite{
		Email:           payload.Email,
		OrganizationId:  organizationUuid,
		Slug:            inviteSlug,
		AccessLevel:     model.UserPermissionLevelEnum(payload.AccessLevel),
		InvitedByUserId: uuid.MustParse(context.Session.User.UniqueId),
		Status:          model.OrganizationInviteStatusEnum_Pending,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	insertQuery := table.OrganizationMemberInvite.INSERT(table.OrganizationMemberInvite.MutableColumns).MODEL(invite).
		RETURNING(table.OrganizationMemberInvite.AllColumns)

	err = insertQuery.QueryContext(context.Request().Context(), context.App.Db, &inviteDest)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// ! TODO: send email to the user for the invite

	response := api_types.CreateInviteResponseSchema{
		Invite: api_types.OrganizationMemberInviteSchema{
			AccessLevel: api_types.UserPermissionLevelEnum(inviteDest.AccessLevel),
			Email:       inviteDest.Email,
			Status:      api_types.InviteStatusEnum(inviteDest.Status),
			CreatedAt:   inviteDest.CreatedAt,
			UniqueId:    inviteDest.UniqueId.String(),
		},
	}

	return context.JSON(http.StatusOK, response)
}

// ! cache the response here and return the cached response, revalidate the cache for this endpoint on template update webhook
func getMessageTemplateById(context interfaces.ContextWithSession) error {
	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid organization id")
	}

	templateId := context.Param("id")

	if templateId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid template id")
	}

	businessAccountDetails := SELECT(table.WhatsappBusinessAccount.AllColumns).
		FROM(table.WhatsappBusinessAccount).
		WHERE(table.WhatsappBusinessAccount.OrganizationId.EQ(UUID(orgUuid))).
		LIMIT(1)

	var businessAccount model.WhatsappBusinessAccount

	err = businessAccountDetails.QueryContext(context.Request().Context(), context.App.Db, &businessAccount)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching business account details")
	}

	if businessAccount.UniqueId.String() == "" || businessAccount.AccessToken == "" || businessAccount.AccountId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Please update your business account details in the settings first.")
	}

	// initialize a wapi client and fetch the templates

	wapiClient := wapi.New(&wapi.ClientConfig{
		BusinessAccountId: businessAccount.AccountId,
		ApiAccessToken:    businessAccount.AccessToken,
		WebhookSecret:     businessAccount.WebhookSecret,
	})

	templateResponse, err := wapiClient.Business.Template.Fetch(templateId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, templateResponse)
}

func getAllMessageTemplates(context interfaces.ContextWithSession) error {
	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid organization id")
	}

	businessAccountDetails := SELECT(table.WhatsappBusinessAccount.AllColumns).
		FROM(table.WhatsappBusinessAccount).
		WHERE(table.WhatsappBusinessAccount.OrganizationId.EQ(UUID(orgUuid))).
		LIMIT(1)

	var businessAccount model.WhatsappBusinessAccount

	err = businessAccountDetails.QueryContext(context.Request().Context(), context.App.Db, &businessAccount)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusOK, []api_types.TemplateSchema{})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching business account details")
	}

	if businessAccount.UniqueId.String() == "" || businessAccount.AccessToken == "" || businessAccount.AccountId == "" {
		// return empty response
		return context.JSON(http.StatusOK, []api_types.TemplateSchema{})
	}

	// initialize a wapi client and fetch the templates

	wapiClient := wapi.New(&wapi.ClientConfig{
		BusinessAccountId: businessAccount.AccountId,
		ApiAccessToken:    businessAccount.AccessToken,
		WebhookSecret:     businessAccount.WebhookSecret,
	})

	templateResponse, err := wapiClient.Business.Template.FetchAll()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// return the templates to the user

	return context.JSON(http.StatusOK, templateResponse.Data)
}

func getAllPhoneNumbers(context interfaces.ContextWithSession) error {

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid organization id")
	}

	businessAccountDetails := SELECT(table.WhatsappBusinessAccount.AllColumns).
		FROM(table.WhatsappBusinessAccount).
		WHERE(table.WhatsappBusinessAccount.OrganizationId.EQ(UUID(orgUuid))).
		LIMIT(1)

	var businessAccount model.WhatsappBusinessAccount

	err = businessAccountDetails.QueryContext(context.Request().Context(), context.App.Db, &businessAccount)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusOK, []api_types.PhoneNumberSchema{})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching business account details")
	}

	if businessAccount.UniqueId.String() == "" || businessAccount.AccessToken == "" || businessAccount.AccountId == "" {
		return context.JSON(http.StatusOK, []api_types.PhoneNumberSchema{})
	}

	// initialize a wapi client and fetch the templates

	wapiClient := wapi.New(&wapi.ClientConfig{
		BusinessAccountId: businessAccount.AccountId,
		ApiAccessToken:    businessAccount.AccessToken,
		WebhookSecret:     businessAccount.WebhookSecret,
	})

	phoneNumbersResponse, err := wapiClient.Business.PhoneNumber.FetchAll(true)

	fmt.Println("phoneNumbersResponse", phoneNumbersResponse)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, phoneNumbersResponse.Data)
}

func getPhoneNumberById(context interfaces.ContextWithSession) error {
	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid organization id")
	}

	businessAccountDetails := SELECT(table.WhatsappBusinessAccount.AllColumns).
		FROM(table.WhatsappBusinessAccount).
		WHERE(table.WhatsappBusinessAccount.OrganizationId.EQ(UUID(orgUuid))).
		LIMIT(1)

	var businessAccount model.WhatsappBusinessAccount

	err = businessAccountDetails.QueryContext(context.Request().Context(), context.App.Db, &businessAccount)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching business account details")
	}

	if businessAccount.UniqueId.String() == "" || businessAccount.AccessToken == "" || businessAccount.AccountId == "" {
		return context.JSON(http.StatusOK, nil)
	}

	phoneNumberId := context.Param("id")

	if phoneNumberId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid phone number id")
	}

	// initialize a wapi client and fetch the templates

	wapiClient := wapi.New(&wapi.ClientConfig{
		BusinessAccountId: businessAccount.AccountId,
		ApiAccessToken:    businessAccount.AccessToken,
		WebhookSecret:     businessAccount.WebhookSecret,
	})

	phoneNumberResponse, err := wapiClient.Business.PhoneNumber.Fetch(phoneNumberId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, phoneNumberResponse)
}

func transferOwnershipOfOrganization(context interfaces.ContextWithSession) error {

	payload := new(api_types.TransferOrganizationOwnershipSchema)

	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	currentUserUuid, err := uuid.Parse(context.Session.User.UniqueId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid user id")
	}

	organizationUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid organization id")
	}

	newOwnerUuid, err := uuid.Parse(payload.NewOwnerId)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid new owner id")
	}

	organizationQuery := SELECT(table.Organization.AllColumns).
		FROM(table.Organization).
		WHERE(table.Organization.UniqueId.EQ(UUID(organizationUuid))).
		LIMIT(1)

	var organization model.Organization
	err = organizationQuery.QueryContext(context.Request().Context(), context.App.Db, &organization)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Organization not found")
	}

	var newOwnerUser model.User
	newUserQuery := SELECT(table.User.AllColumns).
		FROM(table.User).
		WHERE(table.User.UniqueId.EQ(UUID(newOwnerUuid))).
		LIMIT(1)
	err = newUserQuery.QueryContext(context.Request().Context(), context.App.Db, &newOwnerUser)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "New owner not found")
	}

	var newOwnerOrganizationMemberRecord model.OrganizationMember
	newOwnerOrganizationMemberRecordQuery := SELECT(table.OrganizationMember.AllColumns).
		FROM(table.OrganizationMember).
		WHERE(table.OrganizationMember.UniqueId.EQ(UUID(organizationUuid)).AND(table.OrganizationMember.UserId.EQ(UUID(newOwnerUuid)))).
		LIMIT(1)

	err = newOwnerOrganizationMemberRecordQuery.QueryContext(context.Request().Context(), context.App.Db, &newOwnerOrganizationMemberRecord)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching new owner organization member record")
	}

	if newOwnerOrganizationMemberRecord.UniqueId.String() == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "New owner is not a member of this organization")
	}

	// * a cte to swap the roles of the current owner and the new owner in the organization member record table

	updatedOrganizationOriginalOwner := CTE("updated_organization_original_owner")
	updatedOrganizationNewOwner := CTE("updated_organization_new_owner")

	var resp model.OrganizationMember

	stmt := WITH(updatedOrganizationOriginalOwner.AS(
		table.OrganizationMember.UPDATE().
			WHERE(table.OrganizationMember.UniqueId.EQ(UUID(organizationUuid)).
				AND(table.OrganizationMember.UserId.EQ(UUID(currentUserUuid)))).
			SET(table.OrganizationMember.AccessLevel.SET(String(model.UserPermissionLevelEnum_Member.String()))).
			RETURNING(table.OrganizationMember.AllColumns),
	),
		updatedOrganizationNewOwner.AS(
			table.OrganizationMember.UPDATE().
				WHERE(table.OrganizationMember.UniqueId.EQ(UUID(organizationUuid)).
					AND(table.OrganizationMember.UserId.EQ(UUID(newOwnerUuid)))).
				SET(table.OrganizationMember.AccessLevel.SET(String(model.UserPermissionLevelEnum_Owner.String()))).
				RETURNING(table.OrganizationMember.AllColumns),
		),
	)(SELECT(updatedOrganizationOriginalOwner.AllColumns()).FROM(updatedOrganizationOriginalOwner))

	err = stmt.Query(context.App.Db, &resp)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error transferring ownership")
	}

	responseToReturn := api_types.TransferOrganizationOwnershipResponseSchema{
		IsTransferred: true,
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func handleUpdateWhatsappBusinessAccountDetails(context interfaces.ContextWithSession) error {

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid organization id")
	}

	payload := new(api_types.UpdateWhatsAppBusinessAccountDetailsSchema)

	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// ! TODO: sanity check if the details are valid
	businessAccountRecordQuery := SELECT(table.WhatsappBusinessAccount.AllColumns).
		FROM(table.WhatsappBusinessAccount).
		WHERE(table.WhatsappBusinessAccount.OrganizationId.EQ(UUID(orgUuid))).
		LIMIT(1)

	var businessAccount model.WhatsappBusinessAccount

	err = businessAccountRecordQuery.QueryContext(context.Request().Context(), context.App.Db, &businessAccount)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			// create a new record the user is updating its details for the first time

			secretData := utils.WebhookSecretData{
				WhatsappBusinessAccountId: payload.BusinessAccountId,
				OrganizationId:            orgUuid.String(),
			}

			webhookSecret, err := context.App.EncryptionService.EncryptData(secretData)

			insertQuery := table.WhatsappBusinessAccount.
				INSERT(table.WhatsappBusinessAccount.MutableColumns).
				MODEL(model.WhatsappBusinessAccount{
					OrganizationId: orgUuid,
					AccountId:      payload.BusinessAccountId,
					AccessToken:    payload.AccessToken,
					WebhookSecret:  webhookSecret,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}).
				RETURNING(table.WhatsappBusinessAccount.AllColumns)

			err = insertQuery.QueryContext(context.Request().Context(), context.App.Db, &businessAccount)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			responseToReturn := api_types.WhatsAppBusinessAccountDetailsSchema{
				BusinessAccountId: businessAccount.AccountId,
				AccessToken:       businessAccount.AccessToken,
				WebhookSecret:     webhookSecret,
			}

			return context.JSON(http.StatusOK, responseToReturn)

		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	var updatedBusinessAccount model.WhatsappBusinessAccount

	webhookSecret := businessAccount.WebhookSecret

	if businessAccount.AccountId != payload.BusinessAccountId {
		// changing the business account id, let change the webhook secret too

		secretData := utils.WebhookSecretData{
			WhatsappBusinessAccountId: payload.BusinessAccountId,
			OrganizationId:            orgUuid.String(),
		}

		webhookSecret, _ = context.App.EncryptionService.EncryptData(secretData)
	}

	// update the record
	updateQuery := table.WhatsappBusinessAccount.UPDATE(
		table.WhatsappBusinessAccount.AccessToken,
		table.WhatsappBusinessAccount.AccountId,
		table.WhatsappBusinessAccount.WebhookSecret,
		table.WhatsappBusinessAccount.OrganizationId,
		table.WhatsappBusinessAccount.UniqueId,
	).
		MODEL(model.WhatsappBusinessAccount{
			AccountId:      payload.BusinessAccountId,
			AccessToken:    payload.AccessToken,
			OrganizationId: orgUuid,
			UniqueId:       businessAccount.UniqueId,
			WebhookSecret:  webhookSecret,
		}).
		WHERE(table.WhatsappBusinessAccount.UniqueId.EQ(UUID(businessAccount.UniqueId))).
		RETURNING(table.WhatsappBusinessAccount.AllColumns)

	err = updateQuery.QueryContext(context.Request().Context(), context.App.Db, &updatedBusinessAccount)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	responseToReturn := api_types.WhatsAppBusinessAccountDetailsSchema{
		BusinessAccountId: updatedBusinessAccount.AccountId,
		AccessToken:       updatedBusinessAccount.AccessToken,
		WebhookSecret:     updatedBusinessAccount.WebhookSecret,
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func getFullAiConfiguration(context interfaces.ContextWithSession) error {
	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid organization id")
	}

	organizationQuery := SELECT(table.Organization.AllColumns).
		FROM(table.Organization).
		WHERE(table.Organization.UniqueId.EQ(UUID(orgUuid))).
		LIMIT(1)

	var organization model.Organization

	err = organizationQuery.QueryContext(context.Request().Context(), context.App.Db, &organization)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusOK, nil)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	model := api_types.AiModelEnum(*organization.AiModel)

	responseToReturn := api_types.GetAiConfigurationResponseSchema{
		AiConfiguration: api_types.FullAiConfiguration{
			IsEnabled: organization.IsAiEnabled,
			Model:     model,
			ApiKey:    organization.AiApiKey,
		},
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func _verifyAccessToOrganization(context interfaces.ContextWithSession, userId, organizationId uuid.UUID) bool {

	orgQuery := SELECT(table.OrganizationMember.AllColumns, table.Organization.AllColumns).
		FROM(table.OrganizationMember.
			LEFT_JOIN(table.Organization, table.Organization.UniqueId.EQ(table.OrganizationMember.OrganizationId)),
		).
		WHERE(table.OrganizationMember.UserId.EQ(UUID(userId)).
			AND(table.OrganizationMember.OrganizationId.EQ(UUID(organizationId))))

	var dest struct {
		model.OrganizationMember
		Organization model.Organization
	}

	err := orgQuery.Query(context.App.Db, &dest)

	context.App.Logger.Info("dest", dest)

	if err != nil {
		return false
	}

	if dest.Organization.UniqueId.String() == "" {
		return false
	}

	return true
}
