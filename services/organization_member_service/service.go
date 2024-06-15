package organization_member_service

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/database"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
	"github.com/sarthakjdev/wapikit/services"
)

type OrganizationMemberService struct {
	services.BaseService `json:"-,inline"`
}

func NewOrganizationService() *OrganizationMemberService {
	return &OrganizationMemberService{
		BaseService: services.BaseService{
			Name:        "Organization Service",
			RestApiPath: "/api",
			Routes: []interfaces.Route{
				{
					Path:                    "/organization/members",
					Method:                  http.MethodGet,
					Handler:                 GetOrganizationMember,
					IsAuthorizationRequired: true,
				},
			},
		},
	}
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
