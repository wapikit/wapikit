package admin_service

import (
	"net/http"

	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/internal/api_types"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

type AdminService struct {
	services.BaseService `json:"-,inline"`
}

func NewAdminService() *AdminService {
	return &AdminService{
		BaseService: services.BaseService{
			Name:        "Admin Service",
			RestApiPath: "/admin",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/admin/create-role",
					Handler:                 HandleCreateRole,
					Method:                  http.MethodPost,
					PermissionRoleLevel:     api_types.Admin,
					IsAuthorizationRequired: true,
				},
			},
		},
	}
}

func HandleCreateRole(context interfaces.CustomContext) error {
	return nil
}
