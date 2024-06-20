package organization_service

import (
	"net/http"

	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

type UserService struct {
	services.BaseService `json:"-,inline"`
}

func NewUserService() *UserService {
	return &UserService{
		BaseService: services.BaseService{
			Name:        "User Service",
			RestApiPath: "/api",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/user/:id",
					Method:                  http.MethodGet,
					Handler:                 GetUser,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/user/:id/stepOne",
					Method:                  http.MethodDelete,
					Handler:                 DeleteAccountStepOne,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/user/:id",
					Method:                  http.MethodPost,
					Handler:                 UpdateUser,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/user/:id/step-one",
					Method:                  http.MethodDelete,
					Handler:                 DeleteAccountStepOne,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/user/:id/step-two",
					Method:                  http.MethodPost,
					Handler:                 DeleteAccountStetTwo,
					IsAuthorizationRequired: true,
				},
			},
		},
	}
}

func GetUser(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func UpdateUser(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func DeleteAccountStepOne(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func DeleteAccountStetTwo(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}
