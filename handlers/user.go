package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/database"
	"github.com/sarthakjdev/wapikit/internal"
)

func GetUsers(context internal.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func CreateUser(context internal.CustomContext) error {
	payload := new(CreateNewUserHandlerBodySchemaType)
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

func GetUserByID(context internal.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}
