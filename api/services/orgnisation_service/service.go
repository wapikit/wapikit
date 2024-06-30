package organization_service

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
					Handler:                 GetOrganizations,
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
					Path:                    "/api/organization",
					Method:                  http.MethodPost,
					Handler:                 CreateNewOrganization,
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
					Handler:                 UpdateOrganization,
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
					Handler:                 GetOrganizationSettings,
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
					Handler:                 GetOrganizationRoles,
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
					Handler:                 GetRoleById,
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
					Handler:                 DeleteRoleById,
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
					Handler:                 UpdateRoleById,
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
					Handler:                 GetOrganizationMembers,
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
					Handler:                 CreateNewOrganizationMember,
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
					Handler:                 UpdateOrgMemberById,
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
					Handler:                 GetOrgMemberById,
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
					Handler:                 UpdateOrgMemberById,
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
					Path:                    "/api/organization/members/:id/roles",
					Method:                  http.MethodPost,
					Handler:                 UpdateMemberRoles,
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
					Handler:                 UpdateMemberRoles,
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
					Handler:                 UpdateMemberRoles,
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

func CreateNewOrganization(context interfaces.CustomContext) error {
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

	userUuid, err := uuid.FromBytes([]byte(context.Session.User.UniqueId))
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

	err = tx.Commit()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.String(http.StatusOK, "OK")
}

func GetOrganizations(context interfaces.CustomContext) error {
	organizationId := context.Param("id")
	if organizationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization id")
	}

	orgUuid, _ := uuid.FromBytes([]byte(organizationId))

	hasAccess := VerifyAccessToOrganization(context, context.Session.User.UniqueId, organizationId)

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
	return context.JSON(http.StatusOK, api_types.GetOrganizationResponseSchema{
		Organization: &api_types.OrganizationSchema{
			Name:       &dest.Name,
			CreatedAt:  &dest.CreatedAt,
			UniqueId:   &uniqueId,
			FaviconUrl: &dest.FaviconUrl,
			LogoUrl:    dest.LogoUrl,
			WebsiteUrl: dest.WebsiteUrl,
		},
	})
}

func DeleteOrganization(context interfaces.CustomContext) error {
	organizationId := context.Param("id")
	if organizationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization id")
	}

	hasAccess := VerifyAccessToOrganization(context, context.Session.User.UniqueId, organizationId)

	if !hasAccess {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this organization")
	}

	return context.String(http.StatusOK, "OK")
}

func UpdateOrganization(context interfaces.CustomContext) error {
	organizationId := context.Param("id")
	if organizationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization id")
	}

	hasAccess := VerifyAccessToOrganization(context, context.Session.User.UniqueId, organizationId)

	if !hasAccess {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this organization")
	}
	return context.String(http.StatusOK, "OK")
}

func GetOrganizationRoles(context interfaces.CustomContext) error {
	organizationId := context.Param("id")
	if organizationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization id")
	}

	hasAccess := VerifyAccessToOrganization(context, context.Session.User.UniqueId, organizationId)

	if !hasAccess {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this organization")
	}

	return context.String(http.StatusOK, "OK")
}

func GetOrganizationSettings(context interfaces.CustomContext) error {
	organizationId := context.Param("id")
	if organizationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization id")
	}

	hasAccess := VerifyAccessToOrganization(context, context.Session.User.UniqueId, organizationId)
	if !hasAccess {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this organization")
	}

	return context.String(http.StatusOK, "OK")
}

func GetRoleById(context interfaces.CustomContext) error {
	roleId := context.Param("id")
	if roleId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid role id")
	}

	roleUuid, _ := uuid.FromBytes([]byte(roleId))
	roleQuery := SELECT(table.OrganizationRole.AllColumns).FROM(table.OrganizationRole).WHERE(table.OrganizationRole.UniqueId.EQ(UUID(roleUuid))).LIMIT(1)

	var dest model.OrganizationRole
	err := roleQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == "qrm: no rows in result set" {
			role := new(api_types.OrganizationRoleSchema)
			return context.JSON(http.StatusOK, api_types.GetRoleByIdResponseSchema{
				Role: role,
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
		Description: &dest.Description,
		Name:        &dest.Name,
		Permissions: &permissionToReturn,
		UniqueId:    &roleId,
	}

	return context.JSON(http.StatusOK, role)
}

func DeleteRoleById(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func UpdateRoleById(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func GetOrganizationMembers(context interfaces.CustomContext) error {

	params := new(api_types.GetOrganizationMembersParams)

	if err := internal.BindQueryParams(context, &params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	pageNumber := params.Page
	pageSize := params.PerPage
	sortBy := params.SortBy

	organizationUuid, err := uuid.FromBytes([]byte(context.Session.User.OrganizationId))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	organizationMembersQuery := SELECT(table.OrganizationMember.AllColumns,
		table.User.Username,
		table.User.Name,
		table.User.Email,
		table.RoleAssignment.AllColumns,
		table.OrganizationRole.AllColumns,
		COUNT(table.OrganizationMember.OrganizationId.EQ(UUID(organizationUuid))).OVER().AS("totalMembers"),
	).
		FROM(table.OrganizationMember.
			LEFT_JOIN(table.User, table.User.UniqueId.EQ(table.OrganizationMember.UserId)).
			LEFT_JOIN(table.RoleAssignment, table.RoleAssignment.OrganizationMemberId.EQ(table.OrganizationMember.UniqueId)).
			LEFT_JOIN(table.OrganizationRole, table.OrganizationRole.UniqueId.EQ(table.RoleAssignment.OrganizationRoleId))).
		WHERE(table.OrganizationMember.OrganizationId.EQ(UUID(organizationUuid))).
		LIMIT(pageSize).
		OFFSET(pageNumber * pageSize)

	if sortBy != nil {
		if *sortBy == api_types.Asc {
			organizationMembersQuery.ORDER_BY(table.OrganizationMember.CreatedAt.ASC())
		} else {
			organizationMembersQuery.ORDER_BY(table.OrganizationMember.CreatedAt.ASC())
		}
	}

	var dest struct {
		TotalMembers int `json:"totalMembers"`
		members      []struct {
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
				Members: &members,
				PaginationMeta: &api_types.PaginationMeta{
					Page:    &pageNumber,
					PerPage: &pageSize,
					Total:   &total,
				},
			})
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	membersToReturn := make([]api_types.OrganizationMemberSchema, len(dest.members))

	if len(dest.members) > 0 {
		for _, member := range dest.members {
			memberRoles := make([]api_types.OrganizationRoleSchema, len(member.Roles))
			if len(member.Roles) > 0 {
				for _, role := range member.Roles {
					permissions := make([]api_types.RolePermissionEnum, len(role.Permissions))
					for _, perm := range role.Permissions {
						permissions = append(permissions, api_types.RolePermissionEnum(perm))
					}

					roleId := role.UniqueId.String()

					roleToReturn := api_types.OrganizationRoleSchema{
						Description: &role.Description,
						Name:        &role.Name,
						Permissions: &permissions,
						UniqueId:    &roleId,
					}

					memberRoles = append(memberRoles, roleToReturn)
				}
			}

			accessLevel := api_types.UserRoleEnum(member.OrganizationMember.AccessLevel)
			memberId := member.User.UniqueId.String()
			mmbr := api_types.OrganizationMemberSchema{
				CreatedAt:   &member.OrganizationMember.CreatedAt,
				AccessLevel: &accessLevel,
				UniqueId:    &memberId,
				Email:       &member.User.Email,
				Name:        &member.User.Name,
				Roles:       &memberRoles,
			}
			membersToReturn = append(membersToReturn, mmbr)
		}

	}

	return context.JSON(http.StatusOK, api_types.GetOrganizationMembersResponseSchema{
		Members: &membersToReturn,
		PaginationMeta: &api_types.PaginationMeta{
			Page:    &pageNumber,
			PerPage: &pageSize,
			Total:   &dest.TotalMembers,
		}})
}

func CreateNewOrganizationMember(context interfaces.CustomContext) error {
	payload := new(interface{})
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	//     return "", fmt.Errorf("error hashing password: %w", err)
	// }
	// return string(hash), nil

	database.GetDbInstance()
	return context.String(http.StatusOK, "OK")
}

func GetOrgMemberById(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func DeleteOrgMemberById(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func UpdateOrgMemberById(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func UpdateMemberRoles(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func VerifyAccessToOrganization(context interfaces.CustomContext, userId, organizationId string) bool {
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
