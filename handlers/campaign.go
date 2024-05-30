package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetCampaigns(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func CreateNewCampaign(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func GetCampaignByID(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func UpdateCampaignStatus(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
