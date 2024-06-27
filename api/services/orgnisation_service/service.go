package organization_service

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/database"
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
					Path:                    "/api/organization/:id",
					Method:                  http.MethodGet,
					Handler:                 GetOrganization,
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
					Path:                    "/api/organization/:id/settings",
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
					Path:                    "/api/organization/:id/roles",
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
					Path:                    "/api/organization/:id/members",
					Method:                  http.MethodGet,
					Handler:                 GetOrganizationMember,
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
					Path:                    "/api/organization/:id/members",
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

func GetOrganization(context interfaces.CustomContext) error {
	organizationId := context.Param("id")
	if organizationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization id")
	}

	hasAccess := VerifyAccessToOrganization(context, context.Session.User.UniqueId, organizationId)

	if !hasAccess {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this organization")
	}

	var dest model.Organization
	organizationQuery := SELECT(table.Organization.AllColumns).
		FROM(table.Organization).
		WHERE(table.Organization.UniqueId.EQ(String(organizationId)))
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

	// ! verify access to the role organization
	// ! sanity check if the role.orgId == context.Session.User.OrganizationId

	return context.String(http.StatusOK, "OK")
}

func DeleteRoleById(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func UpdateRoleById(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func GetOrganizationMember(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
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
