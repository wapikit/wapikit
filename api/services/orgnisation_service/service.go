package organization_service

import (
	"net/http"

	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

type OrganizationService struct {
	services.BaseService `json:"-,inline"`
}

func NewOrganizationService() *OrganizationService {
	return &OrganizationService{
		BaseService: services.BaseService{
			Name:        "Organization Service",
			RestApiPath: "/api",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/organization/:id",
					Method:                  http.MethodGet,
					Handler:                 GetOrganization,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/organization",
					Method:                  http.MethodPost,
					Handler:                 CreateNewOrganization,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/organization/:id",
					Method:                  http.MethodDelete,
					Handler:                 DeleteOrganization,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/organization/:id",
					Method:                  http.MethodPost,
					Handler:                 UpdateOrganization,
					IsAuthorizationRequired: true,
				},
			},
		},
	}
}

func GetOrganization(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func CreateNewOrganization(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")

}

func DeleteOrganization(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func UpdateOrganization(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}
