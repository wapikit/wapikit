package system_controller

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/wapikit/wapikit/api/api_types"
	controller "github.com/wapikit/wapikit/api/controllers"
	"github.com/wapikit/wapikit/interfaces"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
)

type SystemController struct {
	controller.BaseController `json:"-,inline"`
}

func NewSystemController() *SystemController {
	return &SystemController{
		BaseController: controller.BaseController{
			Name:        "System Controller",
			RestApiPath: "/api/system",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/health",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleHealthCheck),
					IsAuthorizationRequired: false,
				},
				{
					Path:                    "/api/metadata",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetMetaData),
					IsAuthorizationRequired: false,
				},
				{
					Path:                    "/api/system/feature-flags",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetFeatureFlags),
					IsAuthorizationRequired: true,
				},
			},
		},
	}
}

func handleHealthCheck(context interfaces.ContextWithSession) error {
	// get the system metric here
	context.JSON(http.StatusOK, "OK")
	return nil
}

func handleGetMetaData(context interfaces.ContextWithSession) error {
	// get the system metric here

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Error parsing organization UUID")
	}

	var dest model.Organization

	organizationQuery := SELECT(table.Organization.AllColumns).
		FROM(table.Organization).
		WHERE(table.Organization.UniqueId.EQ(UUID(orgUuid))).LIMIT(1)

	err = organizationQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	metaTitle := strings.Join([]string{"WapiKit", dest.Name}, " | ")

	responseToReturn := api_types.GetMetaDataResponseSchema{
		MetaTitle:       &metaTitle,
		MetaDescription: dest.Description,
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func handleGetFeatureFlags(context interfaces.ContextWithSession) error {
	featureFlags := api_types.FeatureFlags{
		SystemFeatureFlags: api_types.SystemFeatureFlags{
			IsAiIntegrationEnabled:                true,
			IsApiAccessEnabled:                    true,
			IsMultiOrganizationEnabled:            true,
			IsRoleBasedAccessControlEnabled:       true,
			IsPluginIntegrationMarketplaceEnabled: true,
			IsCloudEdition:                        false,
			IsEnterpriseEdition:                   false,
		},
	}

	responseToReturn := api_types.GetFeatureFlagsResponseSchema{
		FeatureFlags: featureFlags,
	}

	return context.JSON(http.StatusOK, responseToReturn)
}
