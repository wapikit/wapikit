package system_service

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/wapikit/wapikit/api/services"
	"github.com/wapikit/wapikit/internal/api_types"
	"github.com/wapikit/wapikit/internal/interfaces"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
)

type SystemService struct {
	services.BaseService `json:"-,inline"`
}

func NewSystemService() *SystemService {
	return &SystemService{
		BaseService: services.BaseService{
			Name:        "System Service",
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
					Path:                    "/api/feature-flags",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetFeatureFlags),
					IsAuthorizationRequired: false,
				},
				{
					Path:                    "/api/feature-flags",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(handleGetFeatureFlags),
					IsAuthorizationRequired: false,
				},
			},
		},
	}
}

func handleHealthCheck(context interfaces.ContextWithSession) error {
	// get the system metric here
	context.String(http.StatusOK, "OK")
	return nil
}

func handleGetMetaData(context interfaces.ContextWithSession) error {
	// get the system metric here

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Error parsing organization UUID")
	}

	var dest model.Organization

	organizationQuery := SELECT(table.Organization.AllColumns).
		FROM(table.Organization).
		WHERE(table.Organization.UniqueId.EQ(UUID(orgUuid))).LIMIT(1)

	err = organizationQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	metaTitle := strings.Join([]string{"WapiKit", dest.Name}, " | ")

	responseToReturn := api_types.GetMetaDataResponseSchema{
		MetaTitle:       &metaTitle,
		MetaDescription: dest.Description,
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func handleGetFeatureFlags(context interfaces.ContextWithSession) error {
	userUuid, err := uuid.Parse(context.Session.User.UniqueId)
	if err != nil {
		return context.String(http.StatusInternalServerError, "Error parsing user UUID")
	}
	organizationUuid, err := uuid.Parse(context.Session.User.OrganizationId)
	if err != nil {
		return context.String(http.StatusInternalServerError, "Error parsing organization UUID")
	}

	context.App.Logger.Info("userUuid: %v, organizationUuid: %v", userUuid, organizationUuid)

	// ! TODO: get the integration from backend
	response := api_types.GetFeatureFlagsResponseSchema{
		FeatureFlags: &api_types.FeatureFlags{
			SystemFeatureFlags: &api_types.SystemFeatureFlags{
				IsApiAccessEnabled:              true,
				IsMultiOrganizationEnabled:      true,
				IsRoleBasedAccessControlEnabled: true,
			},
			IntegrationFeatureFlags: &api_types.IntegrationFeatureFlags{
				IsCustomChatBoxIntegrationEnabled: true,
				IsOpenAiIntegrationEnabled:        true,
				IsSlackIntegrationEnabled:         true,
			},
		},
	}

	return context.JSON(http.StatusOK, response)
}
