package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/database"
	"github.com/sarthakjdev/wapikit/internal"
	"golang.org/x/crypto/bcrypt"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/sarthakjdev/wapikit/.db-generated/wapikit/public/model"
	. "github.com/sarthakjdev/wapikit/.db-generated/wapikit/public/table"
)

func HandleSignIn(context internal.CustomContext) error {
	payload := new(AuthHandlerBodySchemaType)

	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if payload.Email == "" || payload.Password == "" {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "email / password is required")
	}

	// get the user from the database

	query := SELECT(OrganisationMember.AllColumns).WHERE(
		OR(
			OrganisationMember.Email.EQ(String(payload.Email)),
			OrganisationMember.Username.EQ(String(payload.Email)))).LIMIT(1)

	user := model.OrganisationMember{}
	query.Query(database.GetDbInstance(), &user)

	if user == (model.OrganisationMember{}) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email / password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email / password")
	}

	// generate the JWT token

	claims := &JwtPayload{
		internal.ContextUser{
			Username: user.Username,
			Email:    user.Email,
			Role:     internal.PermissionRole(user.Role.String()),
			UniqueId: user.UniqueId.String(),
		},
		jwt.StandardClaims{
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
