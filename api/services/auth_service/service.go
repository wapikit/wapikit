package auth_service

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/database"
	"github.com/sarthakjdev/wapikit/internal/api_types"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
	"golang.org/x/crypto/bcrypt"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/sarthakjdev/wapikit/.db-generated/model"
	table "github.com/sarthakjdev/wapikit/.db-generated/table"
)

type AuthService struct {
	services.BaseService `json:"-,inline"`
}

func NewAuthService() *AuthService {
	return &AuthService{
		BaseService: services.BaseService{
			Name:        "Auth Service",
			RestApiPath: "/api/auth",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/login",
					Method:                  http.MethodPost,
					Handler:                 HandleSignIn,
					IsAuthorizationRequired: false,
				},
				{
					Path:                    "/api/register",
					Method:                  http.MethodPost,
					Handler:                 HandleUserRegistration,
					IsAuthorizationRequired: false,
				},
				{
					Path:                    "/api/api-keys",
					Method:                  http.MethodGet,
					Handler:                 GetApiKeys,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/api-keys/regenerate",
					Method:                  http.MethodPost,
					Handler:                 RegenerateApiKey,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/oauth",
					Method:                  http.MethodPost,
					Handler:                 HandleLoginWithOAuth,
					IsAuthorizationRequired: true,
				},
			},
		},
	}
}

func HandleSignIn(context interfaces.CustomContext) error {
	payload := new(api_types.LoginRequestBodySchema)

	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if payload.Username == "" || payload.Password == "" {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "Username / password is required")
	}

	stmt := SELECT(
		table.User.AllColumns,
		table.Organization.AllColumns,
		table.OrganizationMember.AllColumns,
		table.RoleAssignment.AllColumns,
	).FROM(
		table.User.
			LEFT_JOIN(table.OrganizationMember, table.User.UniqueId.EQ(table.OrganizationMember.UserId)).
			LEFT_JOIN(table.Organization, table.OrganizationMember.OrganizationId.EQ(table.Organization.UniqueId)).
			LEFT_JOIN(table.RoleAssignment, table.OrganizationMember.UniqueId.EQ(table.RoleAssignment.OrganizationMemberId)),
	).WHERE(
		table.User.Username.EQ(String(payload.Username)).
			OR(table.User.Email.EQ(String(payload.Username))),
	)

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
	stmt.Query(database.GetDbInstance(), &user)

	// if no user found then return 404
	if user.User.UniqueId.String() == "" {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid email / password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.User.Password), []byte(payload.Password)); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid email / password")
	}

	var isOnboardingCompleted bool
	var organizationIdToLoginWith string
	var roleToLoginWith interfaces.PermissionRole
	var claims *interfaces.JwtPayload

	// if no organization found, then simply return the user with a flag saying isOnboardingCompleted
	if len(user.Organizations) == 0 {
		isOnboardingCompleted = false
		// onboarding to be completed by the user
		// return the user with onboarding flag

		claims = &interfaces.JwtPayload{
			ContextUser: interfaces.ContextUser{
				Username: user.User.Username,
				Email:    user.User.Email,
				UniqueId: user.User.UniqueId.String(),
				Name:     user.User.Name,
			},
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24 * 60).Unix(), // 60-day expiration
				Issuer:    "wapikit",
			},
		}

	} else {
		isOnboardingCompleted = true

		// check for the owner org
		for _, org := range user.Organizations {
			if org.Organization.MemberDetails.Role == model.OrganizationMemberRole_Owner {
				organizationIdToLoginWith = org.Organization.UniqueId.String()
				roleToLoginWith = interfaces.OwnerRole
				break
			}
		}

		// check for the admin org
		if organizationIdToLoginWith == "" {
			// no owner org found, login with the org having the highest role
			// here if no owner org found then look for the lower roles too
			for _, org := range user.Organizations {
				if org.Organization.MemberDetails.Role == model.OrganizationMemberRole_Admin {
					organizationIdToLoginWith = org.Organization.UniqueId.String()
					roleToLoginWith = interfaces.AdminRole
					break
				}
			}
		}

		// else login with the first org
		if organizationIdToLoginWith == "" {
			organizationIdToLoginWith = user.Organizations[0].Organization.UniqueId.String()
			roleToLoginWith = interfaces.MemberRole
		}

		claims = &interfaces.JwtPayload{
			ContextUser: interfaces.ContextUser{
				Username:       user.User.Username,
				Email:          user.User.Email,
				Role:           interfaces.PermissionRole(roleToLoginWith),
				UniqueId:       user.User.UniqueId.String(),
				OrganizationId: organizationIdToLoginWith,
				Name:           user.User.Name,
			},
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24 * 60).Unix(), // 60-day expiration
				Issuer:    "wapikit",
			},
		}

	}

	//Create the token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(context.App.Koa.String("jwt_secret")))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating token")
	}

	return context.JSON(http.StatusOK, api_types.LoginResponseBodySchema{
		IsOnboardingCompleted: &isOnboardingCompleted,
		Token:                 &token,
	})
}

func HandleLoginWithOAuth(context interfaces.CustomContext) error {
	return nil
}

func HandleUserRegistration(context interfaces.CustomContext) error {

	return nil
}

func RegenerateApiKey(context interfaces.CustomContext) error {
	return nil
}

func GetApiKeys(context interfaces.CustomContext) error {
	return nil
}

func SwitchOrganization(context interfaces.CustomContext) error {
	// organization id

	// check if the user even a member of the organization

	return nil
	// this function verifies the access to the organization and returns back a new auth token with the access to the new organization with new data
}

func GetUserRoles(context interfaces.CustomContext) error {
	return nil
}
