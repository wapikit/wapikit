package analytics_service

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/internal/api_types"
	"github.com/sarthakjdev/wapikit/internal/core/utils"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

type IntegrationService struct {
	services.BaseService `json:"-,inline"`
}

func NewIntegrationService() *IntegrationService {
	return &IntegrationService{
		BaseService: services.BaseService{
			Name:        "Integration Service",
			RestApiPath: "/api/integration",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/integration",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetIntegrations),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
				{
					Path:                    "/api/analytics/:id",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetIntegrationById),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
				{
					Path:                    "/api/analytics/:id/enable",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(handleEnableIntegration),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
				{
					Path:                    "/api/analytics/:id/disable",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(handleDisableIntegration),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60,
						},
					},
				},
			},
		},
	}
}

func handleGetIntegrations(context interfaces.ContextWithSession) error {
	params := new(api_types.GetIntegrationsParams)
	if err := utils.BindQueryParams(context, params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func handleGetIntegrationById(context interfaces.ContextWithSession) error {
	return nil
}

func handleEnableIntegration(context interfaces.ContextWithSession) error {
	return nil
}

func handleDisableIntegration(context interfaces.ContextWithSession) error {
	return nil
}
