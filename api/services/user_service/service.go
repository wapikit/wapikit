package organization_service

import (
	"net/http"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/sarthakjdev/wapikit/.db-generated/model"
	table "github.com/sarthakjdev/wapikit/.db-generated/table"
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/internal/api_types"
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
					Path:                    "/api/user",
					Method:                  http.MethodGet,
					Handler:                 GetUser,
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60 * 60, // 1 hour
						},
					},
				},
			},
		},
	}
}

func GetUser(context interfaces.CustomContext) error {
	userQuery := SELECT(
		table.User.AllColumns,
		table.Organization.AllColumns,
		table.OrganizationMember.AllColumns,
		table.RoleAssignment.AllColumns,
	).
		FROM(
			table.User.
				LEFT_JOIN(table.Organization, table.User.UniqueId.EQ(table.OrganizationMember.UserId).AND(table.Organization.UniqueId.EQ(table.OrganizationMember.OrganizationId))).
				LEFT_JOIN(table.OrganizationMember, table.OrganizationMember.OrganizationId.EQ(table.Organization.UniqueId).AND(table.OrganizationMember.UserId.EQ(table.User.UniqueId))).
				LEFT_JOIN(table.RoleAssignment, table.RoleAssignment.OrganizationMemberId.EQ(table.OrganizationMember.UniqueId)),
		).
		WHERE(
			table.User.UniqueId.EQ(String(context.Session.User.UniqueId)).
				AND(
					table.User.Email.EQ(String(context.Session.User.Email)),
				),
		).LIMIT(1)

	type UserWithOrgDetails struct {
		User          model.User `json:"-,inline"`
		Organizations []struct {
			Organization struct {
				model.Organization `json:"-,inline"`
				MemberDetails      model.OrganizationMember `json:"member_details"`
			}
			AssignedRoles []model.RoleAssignment `json:"assigned_roles"`
		} `json:"organizations"`
	}

	user := UserWithOrgDetails{}

	userQuery.Query(context.App.Db, &user)
	role := string(context.Session.User.Role)

	userOrganizations := []api_types.OrganizationSchema{}
	for _, org := range user.Organizations {
		uniqueId := org.Organization.UniqueId.String()
		organization := api_types.OrganizationSchema{
			CreatedAt: &org.Organization.CreatedAt,
			Name:      &org.Organization.Name,
			UniqueId:  &uniqueId,
		}
		userOrganizations = append(userOrganizations, organization)
	}

	response := api_types.GetUserResponseSchema{
		User: &api_types.UserSchema{
			CreatedAt:               &user.User.CreatedAt,
			Name:                    &user.User.Name,
			Email:                   &user.User.Email,
			Username:                &user.User.Username,
			UniqueId:                &context.Session.User.UniqueId,
			CurrentOrganizationRole: &role,
			ProfilePicture:          user.User.ProfilePictureUrl,
			Organizations:           &userOrganizations,
		},
	}

	return context.JSON(http.StatusOK, response)
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
