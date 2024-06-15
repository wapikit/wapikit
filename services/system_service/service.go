package system_service

import (
	"net/http"

	"github.com/sarthakjdev/wapikit/internal/interfaces"
	"github.com/sarthakjdev/wapikit/services"
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
					Path:    "/health",
					Method:  http.MethodGet,
					Handler: HandleHealthCheck,
				},
			},
		},
	}
}

func HandleHealthCheck(context interfaces.CustomContext) error {
	// get the system metric here
	context.String(http.StatusOK, "OK")
	return nil
}
