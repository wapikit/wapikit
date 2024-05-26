package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetUsers(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func CreateUser(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func GetUserByID(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
