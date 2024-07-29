package organization_service

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/internal/api_types"
	"github.com/sarthakjdev/wapikit/internal/core/utils"
	"github.com/sarthakjdev/wapikit/internal/interfaces"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/sarthakjdev/wapikit/.db-generated/model"
	table "github.com/sarthakjdev/wapikit/.db-generated/table"
)

type OrganizationService struct {
	services.BaseService `json:"-,inline"`
}

func NewOrganizationService() *OrganizationService {
	return &OrganizationService{
		BaseService: services.BaseService{
			Name:        "Organization Service",
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
							MaxRequests:    10,
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
						PermissionRoleLevel: api_types.Admin,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
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
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
					},
				},
				{
					Path:                    "/api/organization/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getOrganizationById),
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
					Path:                    "/api/organization/settings",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(getOrganizationSettings),
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
					Path:                    "/api/organization/invites",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getOrganizationInvites),
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
					Path:                    "/api/organization/invites",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(createNewOrganizationInvite),
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
					Path:                    "/api/organization/roles",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getOrganizationRoles),
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
					Path:                    "/api/organization/roles/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getRoleById),
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
					Path:                    "/api/organization/roles/:id",
					Method:                  http.MethodDelete,
					Handler:                 interfaces.HandlerWithSession(deleteRoleById),
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
					Path:                    "/api/organization/roles/:id",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(updateRoleById),
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
					Path:                    "/api/organization/members",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getOrganizationMembers),
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
					Path:                    "/api/organization/members",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(createNewOrganizationMember),
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
					Path:                    "/api/organization/members/:id",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(updateOrgMemberById),
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
					Path:                    "/api/organization/members/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getOrgMemberById),
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
					Path:                    "/api/organization/members/:id",
					Method:                  http.MethodDelete,
					Handler:                 interfaces.HandlerWithSession(updateOrgMemberById),
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
					Path:                    "/api/organization/members/:id/role",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(updateOrganizationMemberRoles),
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
					Path:                    "/api/organization/syncTemplates",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(syncTemplates),
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
					Path:                    "/api/organization/syncMobileNumbers",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(syncMobileNumbers),
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
	err = table.Organization.INSERT().
		MODEL(newOrg).
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
	err = table.OrganizationMember.INSERT().MODEL(model.OrganizationMember{
		AccessLevel:    model.UserPermissionLevel_Owner,
		OrganizationId: newOrg.UniqueId,
		UserId:         userUuid,
	}).RETURNING(table.OrganizationMember.AllColumns).QueryContext(context.Request().Context(), tx, &member)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 3. Create API key for the organization
	claims := &interfaces.JwtPayload{
		ContextUser: interfaces.ContextUser{
			Username:       context.Session.User.Username,
			Email:          context.Session.User.Email,
			Role:           api_types.UserPermissionLevel(api_types.Owner),
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

	err = table.ApiKey.INSERT().MODEL(model.ApiKey{
		MemberId:       member.UniqueId,
		OrganizationId: newOrg.UniqueId,
		Key:            token,
	}).RETURNING(table.ApiKey.AllColumns).QueryContext(context.Request().Context(), tx, &apiKey)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.String(http.StatusOK, "OK")
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

	var dest struct {
		TotalOrganizations int `json:"total_organizations"`
		Organizations      []model.Organization
	}

	err = orgQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		context.App.Logger.Info("no rows in result set error occurred")
		if err.Error() == "qrm: no rows in result set" {
			organizations := make([]api_types.OrganizationSchema, 0)
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
	for _, org := range dest.Organizations {
		uniqueId := org.UniqueId.String()
		organization := api_types.OrganizationSchema{
			CreatedAt: org.CreatedAt,
			Name:      org.Name,
			UniqueId:  uniqueId,
		}
		userOrganizations = append(userOrganizations, organization)
	}

	response := api_types.GetOrganizationsResponseSchema{
		Organizations: userOrganizations,
		PaginationMeta: api_types.PaginationMeta{
			Page:    param.Page,
			PerPage: param.PerPage,
			Total:   dest.TotalOrganizations,
		},
	}

	return context.JSON(http.StatusOK, response)
}

func getOrganizationById(context interfaces.ContextWithSession) error {
	organizationId := context.Param("id")
	if organizationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization id")
	}

	orgUuid, _ := uuid.Parse(organizationId)
	hasAccess := verifyAccessToOrganization(context, context.Session.User.UniqueId, organizationId)

	if !hasAccess {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this organization")
	}

	var dest model.Organization
	organizationQuery := SELECT(table.Organization.AllColumns).
		FROM(table.Organization).
		WHERE(table.Organization.UniqueId.EQ(UUID(orgUuid)))
	err := organizationQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	uniqueId := dest.UniqueId.String()
	return context.JSON(http.StatusOK, api_types.GetOrganizationByIdResponseSchema{
		Organization: api_types.OrganizationSchema{
			Name:       dest.Name,
			CreatedAt:  dest.CreatedAt,
			UniqueId:   uniqueId,
			FaviconUrl: &dest.FaviconUrl,
			LogoUrl:    dest.LogoUrl,
			WebsiteUrl: dest.WebsiteUrl,
		},
	})
}

func deleteOrganization(context interfaces.ContextWithSession) error {

	return context.String(http.StatusInternalServerError, "NOT IMPLEMENTED YET")

	organizationId := context.Param("id")
	if organizationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization id")
	}

	hasAccess := verifyAccessToOrganization(context, context.Session.User.UniqueId, organizationId)

	if !hasAccess {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this organization")
	}

	return context.String(http.StatusOK, "OK")
}

func updateOrganizationById(context interfaces.ContextWithSession) error {
	organizationId := context.Param("id")
	if organizationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization id")
	}

	hasAccess := verifyAccessToOrganization(context, context.Session.User.UniqueId, organizationId)

	if !hasAccess {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this organization")
	}

	payload := new(api_types.UpdateOrganizationSchema)

	if payload.Name == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Name is required")
	}

	orgUuid, _ := uuid.Parse(organizationId)

	updateOrgQuery := table.Organization.
		UPDATE(table.Organization.Name).
		SET(*payload.Name).
		WHERE(table.Organization.UniqueId.EQ(UUID(orgUuid)))

	results, err := updateOrgQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if rows, _ := results.RowsAffected(); rows == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	return context.String(http.StatusOK, "OK")
}

func getOrganizationRoles(context interfaces.ContextWithSession) error {
	params := new(api_types.GetOrganizationRolesParams)
	err := utils.BindQueryParams(context, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var dest struct {
		TotalRoles int `json:"totalRoles"`
		roles      []struct {
			model.OrganizationRole
		}
	}

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	whereCondition := table.OrganizationRole.OrganizationId.EQ(UUID(orgUuid))

	organizationRolesQuery := SELECT(table.OrganizationRole.AllColumns).
		FROM(table.OrganizationRole).
		WHERE(whereCondition).
		LIMIT(params.PerPage).
		OFFSET((params.Page - 1) * params.PerPage)

	if params.SortBy != nil {
		if *params.SortBy == api_types.Asc {
			organizationRolesQuery.ORDER_BY(table.OrganizationRole.CreatedAt.ASC())
		} else {
			organizationRolesQuery.ORDER_BY(table.OrganizationRole.CreatedAt.DESC())
		}
	}

	err = organizationRolesQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == "qrm: no rows in result set" {
			roles := make([]api_types.OrganizationRoleSchema, 0)
			total := 0
			return context.JSON(http.StatusOK, api_types.GetOrganizationRolesResponseSchema{
				Roles: roles,
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

	rolesToReturn := make([]api_types.OrganizationRoleSchema, len(dest.roles))

	if len(dest.roles) > 0 {
		for _, role := range dest.roles {
			permissions := make([]api_types.RolePermissionEnum, len(role.Permissions))
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

			rolesToReturn = append(rolesToReturn, roleToReturn)
		}
	}

	return context.JSON(http.StatusOK, api_types.GetOrganizationRolesResponseSchema{
		Roles: rolesToReturn,
		PaginationMeta: api_types.PaginationMeta{
			Page:    params.Page,
			PerPage: params.PerPage,
			Total:   dest.TotalRoles,
		},
	})
}

func getOrganizationSettings(context interfaces.ContextWithSession) error {
	organizationId := context.Param("id")
	if organizationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization id")
	}

	hasAccess := verifyAccessToOrganization(context, context.Session.User.UniqueId, organizationId)
	if !hasAccess {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this organization")
	}

	return context.String(http.StatusOK, "OK")
}

func getRoleById(context interfaces.ContextWithSession) error {
	roleId := context.Param("id")
	if roleId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid role id")
	}

	roleUuid, _ := uuid.Parse(roleId)
	roleQuery := SELECT(table.OrganizationRole.AllColumns).FROM(table.OrganizationRole).WHERE(table.OrganizationRole.UniqueId.EQ(UUID(roleUuid))).LIMIT(1)

	var dest model.OrganizationRole
	err := roleQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == "qrm: no rows in result set" {
			role := new(api_types.OrganizationRoleSchema)
			return context.JSON(http.StatusOK, api_types.GetRoleByIdResponseSchema{
				Role: *role,
			})
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	if dest.OrganizationId.String() != context.Session.User.OrganizationId {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this resource")
	}

	permissionToReturn := make([]api_types.RolePermissionEnum, len(dest.Permissions))

	for _, perm := range dest.Permissions {
		permissionToReturn = append(permissionToReturn, api_types.RolePermissionEnum(perm))
	}

	role := api_types.OrganizationRoleSchema{
		Description: dest.Description,
		Name:        dest.Name,
		Permissions: permissionToReturn,
		UniqueId:    roleId,
	}

	return context.JSON(http.StatusOK, role)
}

func deleteRoleById(context interfaces.ContextWithSession) error {
	// ! destructive endpoint, we are currently allowing deletion of roles even though its being assigned to user,
	// ! at the frontend, there must be double confirmation before deleting a role

	roleId := context.Param("id")
	if roleId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid role id")
	}

	roleUuid, _ := uuid.Parse(roleId)

	// check if the role exists and belongs to the organization

	var role model.OrganizationRole

	existingRoleQuery := SELECT(table.OrganizationRole.AllColumns).
		WHERE(table.OrganizationRole.UniqueId.EQ(UUID(roleUuid))).
		LIMIT(1)

	err := existingRoleQuery.QueryContext(context.Request().Context(), context.App.Db, &role)

	if err != nil {
		if err.Error() == "qrm: no rows in result set" {
			return echo.NewHTTPError(http.StatusNotFound, "Role not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	if role.OrganizationId.String() != context.Session.User.OrganizationId {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this resource")
	}

	roleAssignmentDeleteQuery := table.RoleAssignment.DELETE().WHERE(table.RoleAssignment.OrganizationRoleId.EQ(UUID(roleUuid)))

	_, err = roleAssignmentDeleteQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// delete the role

	roleQuery := table.OrganizationRole.DELETE().WHERE(table.OrganizationRole.UniqueId.EQ(UUID(roleUuid)))

	_, err = roleQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := api_types.DeleteRoleByIdResponseSchema{
		Data: true,
	}

	return context.JSON(http.StatusOK, response)
}

func updateRoleById(context interfaces.ContextWithSession) error {
	roleId := context.Param("id")

	if roleId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid role id")
	}

	roleUuid, _ := uuid.Parse(roleId)

	payload := new(api_types.RoleUpdateSchema)

	if &payload.Name == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Name is required")
	}

	// check if the role exists and belongs to the organization

	var role model.OrganizationRole

	existingRoleQuery := SELECT(table.OrganizationRole.AllColumns).
		WHERE(table.OrganizationRole.UniqueId.EQ(UUID(roleUuid))).
		LIMIT(1)

	err := existingRoleQuery.QueryContext(context.Request().Context(), context.App.Db, &role)

	if err != nil {
		if err.Error() == "qrm: no rows in result set" {
			return echo.NewHTTPError(http.StatusNotFound, "Role not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	if role.OrganizationId.String() != context.Session.User.OrganizationId {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this resource")
	}

	updatedPermissions := make([]model.OrganizaRolePermissionEnum, len(payload.Permissions))

	for _, perm := range payload.Permissions {
		updatedPermissions = append(updatedPermissions, model.OrganizaRolePermissionEnum(perm))
	}

	var updatedRole model.OrganizationRole

	// update the role
	updateRoleQuery := table.OrganizationRole.
		UPDATE(table.OrganizationRole.Name, table.OrganizationRole.Description, table.OrganizationRole.Permissions).
		SET(payload.Name, *payload.Description, updatedPermissions).
		WHERE(table.OrganizationRole.UniqueId.EQ(UUID(roleUuid))).
		RETURNING(table.OrganizationRole.AllColumns)

	err = updateRoleQuery.QueryContext(context.Request().Context(), context.App.Db, &updatedRole)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	permissionsToReturn := make([]api_types.RolePermissionEnum, len(updatedRole.Permissions))

	for _, perm := range updatedRole.Permissions {
		permissionsToReturn = append(permissionsToReturn, api_types.RolePermissionEnum(perm))
	}

	roleToReturn := api_types.OrganizationRoleSchema{
		Description: updatedRole.Description,
		Name:        updatedRole.Name,
		Permissions: permissionsToReturn,
		UniqueId:    roleId,
	}

	return context.JSON(http.StatusOK, api_types.UpdateRoleByIdResponseSchema{
		Role: roleToReturn,
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

	var dest struct {
		TotalMembers int `json:"totalMembers"`
		Members      []struct {
			model.OrganizationMember
			model.User
			Roles []struct {
				model.OrganizationRole
			}
		}
	}

	err = organizationMembersQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == "qrm: no rows in result set" {
			members := make([]api_types.OrganizationMemberSchema, 0)
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

	membersToReturn := make([]api_types.OrganizationMemberSchema, len(dest.Members))

	if len(dest.Members) > 0 {
		for _, member := range dest.Members {
			memberRoles := make([]api_types.OrganizationRoleSchema, len(member.Roles))
			if len(member.Roles) > 0 {
				for _, role := range member.Roles {
					permissions := make([]api_types.RolePermissionEnum, len(role.Permissions))
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

			accessLevel := api_types.UserPermissionLevel(member.OrganizationMember.AccessLevel)
			memberId := member.User.UniqueId.String()
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

	return context.JSON(http.StatusOK, api_types.GetOrganizationMembersResponseSchema{
		Members: membersToReturn,
		PaginationMeta: api_types.PaginationMeta{
			Page:    pageNumber,
			PerPage: pageSize,
			Total:   dest.TotalMembers,
		}})
}

func createNewOrganizationMember(context interfaces.ContextWithSession) error {
	payload := new(interface{})
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	//     return "", fmt.Errorf("error hashing password: %w", err)
	// }
	// return string(hash), nil

	return context.String(http.StatusOK, "OK")
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
		if err.Error() == "qrm: no rows in result set" {
			member := new(api_types.OrganizationMemberSchema)
			return context.JSON(http.StatusOK, api_types.GetOrganizationMemberByIdResponseSchema{
				Member: *member,
			})
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	memberRoles := make([]api_types.OrganizationRoleSchema, len(dest.member.Roles))
	if len(dest.member.Roles) > 0 {
		for _, role := range dest.member.Roles {
			permissions := make([]api_types.RolePermissionEnum, len(role.Permissions))
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

	accessLevel := api_types.UserPermissionLevel(dest.member.OrganizationMember.AccessLevel)

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
	return context.String(http.StatusOK, "OK")
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

	roleUuid, err := uuid.Parse(*payload.RoleUniqueId)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid member id")
	}

	orgRole := SELECT(table.OrganizationRole.AllColumns, table.RoleAssignment.AllColumns).
		FROM(table.OrganizationRole.
			LEFT_JOIN(table.RoleAssignment, table.RoleAssignment.OrganizationRoleId.EQ(table.OrganizationRole.UniqueId)),
		).WHERE(table.OrganizationRole.UniqueId.EQ(UUID(roleUuid))).
		LIMIT(1)

	var dest struct {
		model.OrganizationRole
		Assignment model.RoleAssignment
	}

	err = orgRole.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == "qrm: no rows in result set" {
			return echo.NewHTTPError(http.StatusNotFound, "Role not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	switch payload.Action {
	case api_types.Add:
		{

			// Check if the role is already assigned to the member
			if dest.Assignment.UniqueId != uuid.Nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Role already assigned to the member")
			} else {
				// Assign the role to the member
				var roleAssignmentDest model.RoleAssignment

				roleAssignment := model.RoleAssignment{
					OrganizationMemberId: memberUuid,
					OrganizationRoleId:   dest.UniqueId,
				}

				err := table.RoleAssignment.INSERT().
					MODEL(roleAssignment).
					RETURNING(table.RoleAssignment.AllColumns).
					QueryContext(context.Request().Context(), context.App.Db, &roleAssignmentDest)

				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}

				if roleAssignmentDest.UniqueId == uuid.Nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "Error assigning role")
				} else {
					return context.String(http.StatusOK, "OK")
				}
			}
		}

	case api_types.Remove:
		{
			if dest.Assignment.UniqueId == uuid.Nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Role not assigned to the member")
			} else {
				_, err := table.RoleAssignment.DELETE().
					WHERE(table.RoleAssignment.UniqueId.EQ(UUID(dest.UniqueId))).
					ExecContext(context.Request().Context(), context.App.Db)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}

				return context.String(http.StatusOK, "OK")
			}
		}
	}

	return context.String(http.StatusOK, "OK")
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
		if err.Error() == "qrm: no rows in result set" {
			invites := make([]api_types.OrganizationMemberInviteSchema, 0)
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

	invitesToReturn := make([]api_types.OrganizationMemberInviteSchema, len(dest.Invites))

	if len(dest.Invites) > 0 {
		for _, invite := range dest.Invites {
			accessLevel := api_types.UserPermissionLevel(invite.AccessLevel)
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
		if err.Error() == "qrm: no rows in result set" {
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
		AccessLevel:     model.UserPermissionLevel(payload.AccessLevel),
		InvitedByUserId: uuid.MustParse(context.Session.User.UniqueId),
		Status:          model.OrganizationInviteStatusEnum_Pending,
	}

	insertQuery := table.OrganizationMemberInvite.INSERT().MODEL(invite).
		RETURNING(table.OrganizationMemberInvite.AllColumns)

	err = insertQuery.QueryContext(context.Request().Context(), context.App.Db, &inviteDest)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// ! TODO: send email to the user for the invite

	response := api_types.CreateInviteResponseSchema{
		Invite: api_types.OrganizationMemberInviteSchema{
			AccessLevel: api_types.UserPermissionLevel(inviteDest.AccessLevel),
			Email:       inviteDest.Email,
			Status:      api_types.InviteStatusEnum(inviteDest.Status),
			CreatedAt:   inviteDest.CreatedAt,
			UniqueId:    inviteDest.UniqueId.String(),
		},
	}

	return context.JSON(http.StatusOK, response)
}

func syncTemplates(context interfaces.ContextWithSession) error {
	return context.String(http.StatusOK, "OK")
}

func syncMobileNumbers(context interfaces.ContextWithSession) error {
	return context.String(http.StatusOK, "OK")
}

func verifyAccessToOrganization(context interfaces.ContextWithSession, userId, organizationId string) bool {
	orgQuery := SELECT(table.OrganizationMember.AllColumns, table.Organization.AllColumns).
		FROM(table.OrganizationMember.
			LEFT_JOIN(table.Organization, table.Organization.UniqueId.EQ(table.OrganizationMember.OrganizationId)),
		).
		WHERE(table.OrganizationMember.UserId.EQ(String(userId)).
			AND(table.OrganizationMember.OrganizationId.EQ(String(organizationId))))

	var dest struct {
		model.OrganizationMember
		Organization model.Organization
	}

	err := orgQuery.Query(context.App.Db, &dest)

	if err != nil {
		return false
	}

	if dest.Organization.UniqueId.String() == "" {
		return false
	}

	return true
}
