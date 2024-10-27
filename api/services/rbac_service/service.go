package organization_service

import (
	"net/http"

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

type RoleBasedAccessControlService struct {
	services.BaseService `json:"-,inline"`
}

func NewRoleBasedAccessControlService() *RoleBasedAccessControlService {
	return &RoleBasedAccessControlService{
		BaseService: services.BaseService{
			Name:        "Role Based Access Control Service",
			RestApiPath: "/api/rbac",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/rbac/roles",
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
					Path:                    "/api/rbac/roles/:id",
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
					Path:                    "/api/rbac/roles/:id",
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
					Path:                    "/api/rbac/roles/:id",
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
			},
		},
	}
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
			var roles []api_types.OrganizationRoleSchema
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

	var rolesToReturn []api_types.OrganizationRoleSchema

	if len(dest.roles) > 0 {
		for _, role := range dest.roles {
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

	var permissionToReturn []api_types.RolePermissionEnum

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

	var updatedPermissions []model.OrganizaRolePermissionEnum

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

	var permissionsToReturn []api_types.RolePermissionEnum

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
