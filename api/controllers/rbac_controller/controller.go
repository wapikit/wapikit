package rbac_controller

import (
	"net/http"
	"strings"
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

type RoleBasedAccessControlController struct {
	controller.BaseController `json:"-,inline"`
}

func NewRoleBasedAccessControlController() *RoleBasedAccessControlController {
	return &RoleBasedAccessControlController{
		BaseController: controller.BaseController{
			Name:        "Role Based Access Control Controller",
			RestApiPath: "/api/rbac",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/rbac/roles",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getOrganizationRoles),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetOrganizationRole,
						},
					},
				},
				{
					Path:                    "/api/rbac/roles",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(createRole),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.CreateOrganizationRole,
						},
					},
				},
				{
					Path:                    "/api/rbac/roles/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getRoleById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetOrganizationRole,
						},
					},
				},
				{
					Path:                    "/api/rbac/roles/:id",
					Method:                  http.MethodDelete,
					Handler:                 interfaces.HandlerWithSession(deleteRoleById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.DeleteOrganizationRole,
						},
					},
				},
				{
					Path:                    "/api/rbac/roles/:id",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(updateRoleById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60, // 1 minute
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.UpdateOrganizationRole,
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
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	var dest []struct {
		TotalRoles int `json:"totalRoles"`
		model.OrganizationRole
	}

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
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
		if err.Error() == qrm.ErrNoRows.Error() {
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
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	var rolesToReturn []api_types.OrganizationRoleSchema

	if len(dest) > 0 {
		for _, role := range dest {
			var permissions []api_types.RolePermissionEnum

			// ! convert the permissions string to an array of RolePermissionEnum
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

			rolesToReturn = append(rolesToReturn, roleToReturn)
		}
	}

	totalRoles := 0

	if len(dest) > 0 {
		totalRoles = dest[0].TotalRoles
	}

	return context.JSON(http.StatusOK, api_types.GetOrganizationRolesResponseSchema{
		Roles: rolesToReturn,
		PaginationMeta: api_types.PaginationMeta{
			Page:    params.Page,
			PerPage: params.PerPage,
			Total:   totalRoles,
		},
	})
}

func getRoleById(context interfaces.ContextWithSession) error {
	roleId := context.Param("id")
	if roleId == "" {
		return context.JSON(http.StatusBadRequest, "Invalid role id")
	}

	roleUuid, _ := uuid.Parse(roleId)
	roleQuery := SELECT(table.OrganizationRole.AllColumns).FROM(table.OrganizationRole).WHERE(table.OrganizationRole.UniqueId.EQ(UUID(roleUuid))).LIMIT(1)

	var dest model.OrganizationRole
	err := roleQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			role := new(api_types.OrganizationRoleSchema)
			return context.JSON(http.StatusOK, api_types.GetRoleByIdResponseSchema{
				Role: *role,
			})
		} else {
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	if dest.OrganizationId.String() != context.Session.User.OrganizationId {
		return context.JSON(http.StatusForbidden, "You do not have access to this resource")
	}

	var permissionToReturn []api_types.RolePermissionEnum

	permissionArray := strings.Split(dest.Permissions, ",")

	for _, perm := range permissionArray {
		if perm == "" {
			continue
		}
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

func createRole(context interfaces.ContextWithSession) error {
	payload := new(api_types.NewOrganizationRoleSchema)
	if err := context.Bind(payload); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	var permissions string

	// create a command separated string of permissions

	for _, perm := range payload.Permissions {
		permissions += string(perm) + ","
	}

	var insertedRole model.OrganizationRole

	err = table.OrganizationRole.INSERT(table.OrganizationRole.MutableColumns).
		MODEL(model.OrganizationRole{
			Name:           payload.Name,
			Description:    payload.Description,
			Permissions:    permissions,
			OrganizationId: orgUuid,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}).
		RETURNING(table.OrganizationRole.AllColumns).
		QueryContext(context.Request().Context(), context.App.Db, &insertedRole)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	permissionsToReturn := []api_types.RolePermissionEnum{}

	permissionArray := strings.Split(insertedRole.Permissions, ",")

	for _, perm := range permissionArray {
		if perm == "" {
			continue
		}
		permissionsToReturn = append(permissionsToReturn, api_types.RolePermissionEnum(perm))
	}

	roleToReturn := api_types.OrganizationRoleSchema{
		Description: insertedRole.Description,
		Name:        insertedRole.Name,
		Permissions: permissionsToReturn,
		UniqueId:    insertedRole.UniqueId.String(),
	}

	return context.JSON(http.StatusCreated, api_types.CreateNewRoleResponseSchema{
		Role: roleToReturn,
	})
}

func deleteRoleById(context interfaces.ContextWithSession) error {
	// ! destructive endpoint, we are currently allowing deletion of roles even though its being assigned to user,
	// ! at the frontend, there must be double confirmation before deleting a role

	roleId := context.Param("id")
	if roleId == "" {
		return context.JSON(http.StatusBadRequest, "Invalid role id")
	}

	roleUuid, _ := uuid.Parse(roleId)

	// check if the role exists and belongs to the organization

	var role model.OrganizationRole

	existingRoleQuery := SELECT(table.OrganizationRole.AllColumns).
		WHERE(table.OrganizationRole.UniqueId.EQ(UUID(roleUuid))).
		LIMIT(1)

	err := existingRoleQuery.QueryContext(context.Request().Context(), context.App.Db, &role)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "Role not found")
		} else {
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	if role.OrganizationId.String() != context.Session.User.OrganizationId {
		return context.JSON(http.StatusForbidden, "You do not have access to this resource")
	}

	roleAssignmentDeleteQuery := table.RoleAssignment.DELETE().WHERE(table.RoleAssignment.OrganizationRoleId.EQ(UUID(roleUuid)))

	_, err = roleAssignmentDeleteQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	// delete the role

	roleQuery := table.OrganizationRole.DELETE().WHERE(table.OrganizationRole.UniqueId.EQ(UUID(roleUuid)))

	_, err = roleQuery.ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	response := api_types.DeleteRoleByIdResponseSchema{
		Data: true,
	}

	return context.JSON(http.StatusOK, response)
}

func updateRoleById(context interfaces.ContextWithSession) error {
	roleId := context.Param("id")

	if roleId == "" {
		return context.JSON(http.StatusBadRequest, "Invalid role id")
	}

	roleUuid, _ := uuid.Parse(roleId)

	payload := new(api_types.RoleUpdateSchema)

	// check if the role exists and belongs to the organization

	var role model.OrganizationRole

	existingRoleQuery := SELECT(table.OrganizationRole.AllColumns).
		WHERE(table.OrganizationRole.UniqueId.EQ(UUID(roleUuid))).
		LIMIT(1)

	err := existingRoleQuery.QueryContext(context.Request().Context(), context.App.Db, &role)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "Role not found")
		} else {
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	if role.OrganizationId.String() != context.Session.User.OrganizationId {
		return context.JSON(http.StatusForbidden, "You do not have access to this resource")
	}

	var updatedPermissions string

	for _, perm := range payload.Permissions {
		updatedPermissions += string(perm) + ","
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
		return context.JSON(http.StatusInternalServerError, err.Error())
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
