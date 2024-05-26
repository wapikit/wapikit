package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func handleGetCampaigns(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func handleCreateNewCampaign(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func handleGetCampaignByID(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func updateCampaignStatus(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
