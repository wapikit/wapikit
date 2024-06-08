package handlers

import (
	"net/http"

	"github.com/sarthakjdev/wapikit/internal"
)

func GetCampaigns(context internal.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func CreateNewCampaign(context internal.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func GetCampaignByID(context internal.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func UpdateCampaignStatus(context internal.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}
