package integration_controller

import (
	"net/http"

	"github.com/wapikit/wapikit/api/api_types"
	controller "github.com/wapikit/wapikit/api/controllers"
	"github.com/wapikit/wapikit/interfaces"
	"github.com/wapikit/wapikit/utils"
)

type IntegrationController struct {
	controller.BaseController `json:"-,inline"`
}

func NewIntegrationController() *IntegrationController {
	return &IntegrationController{
		BaseController: controller.BaseController{
			Name:        "Integration Controller",
			RestApiPath: "/api/integration",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/integrations",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetIntegrations),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
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
							MaxRequests:    60,
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
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60,
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.UpdateIntegrationSettings,
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
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60,
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.UpdateIntegrationSettings,
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
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	responseToReturn := api_types.GetIntegrationResponseSchema{
		Integrations: []api_types.IntegrationSchema{},
		PaginationMeta: api_types.PaginationMeta{
			Total:   0,
			Page:    1,
			PerPage: 10,
		},
	}

	return context.JSON(http.StatusOK, responseToReturn)

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
