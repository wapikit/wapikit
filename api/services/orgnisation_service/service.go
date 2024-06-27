package organization_service

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/database"
	"github.com/sarthakjdev/wapikit/internal/api_types"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
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

	

	return context.String(http.StatusOK, "OK")
}

func GetOrganization(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func DeleteOrganization(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func UpdateOrganization(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func GetOrganizationRoles(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func GetOrganizationSettings(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func GetRoleById(context interfaces.CustomContext) error {
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
