package contact_service

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/database"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

type ContactService struct {
	services.BaseService `json:"-,inline"`
}

func NewContactService() *ContactService {
	return &ContactService{
		BaseService: services.BaseService{
			Name:        "Contact Service",
			RestApiPath: "/api/contact",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/contacts",
					Method:                  http.MethodGet,
					Handler:                 GetContacts,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/contacts",
					Method:                  http.MethodPost,
					Handler:                 CreateNewContacts,
					IsAuthorizationRequired: true,
				},
			},
		},
	}
}

func GetContacts(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func CreateNewContacts(context interfaces.CustomContext) error {
	payload := new(interface{})
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	database.GetDbInstance()
	return context.String(http.StatusOK, "OK")
}
