package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/internal"
)

func HandleHealthCheck(context internal.CustomContext) error {

	// get the system metric here
	context.String(http.StatusOK, "OK")
	return nil
}

func UpdateBusinessAccountDetails(context internal.CustomContext) error {
	payload := new(interface{})
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// if payload.BusinessAccountId == "" {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "businessAccountId is required")
	// }

	context.String(http.StatusOK, "OK")
	return nil
}
