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
	. "github.com/sarthakjdev/wapikit/.db-generated/table"
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
					Path:                    "/api/regenerate",
					Method:                  http.MethodPost,
					Handler:                 RegenerateApiKey,
					IsAuthorizationRequired: true,
				},
				// Add more routes as needed
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

	query := SELECT(OrganisationMember.AllColumns).WHERE(
		OR(
			OrganisationMember.Email.EQ(String(payload.Username)),
			OrganisationMember.Username.EQ(String(payload.Username)))).LIMIT(1)

	user := model.OrganisationMember{}
	query.Query(database.GetDbInstance(), &user)

	if user == (model.OrganisationMember{}) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email / password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email / password")
	}

	// generate the JWT token

	claims := &interfaces.JwtPayload{
		ContextUser: interfaces.ContextUser{
			Username: user.Username,
			Email:    user.Email,
			Role:     interfaces.PermissionRole(user.Role.String()),
			UniqueId: user.UniqueId.String(),
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60).Unix(), // 60-day expiration
			Issuer:    "wapikit",
		},
	}
	//Create the token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(context.App.Koa.String("jwt_secret")))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating token")
	}

	return context.JSON(http.StatusOK, map[string]string{
		"token":  token,
		"Status": "OK",
	})
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
