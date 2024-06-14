package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/database"
	"github.com/sarthakjdev/wapikit/internal"

	. "github.com/go-jet/jet/v2/postgres"
	model "github.com/sarthakjdev/wapikit/.db-generated/model"
	table "github.com/sarthakjdev/wapikit/.db-generated/table"
)

func GetCampaigns(context internal.CustomContext) error {
	params := new(GetCampaignsParams)
	if err := context.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return context.String(http.StatusOK, "OK")
}

func CreateNewCampaign(context internal.CustomContext) error {
	payload := new(CreateCampaignJSONRequestBody)
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return context.String(http.StatusOK, "OK")
}

func GetCampaignByID(context internal.CustomContext) error {
	campaignId := context.Param("id")
	if campaignId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Campaign Id")
	}

	sqlStatement := SELECT(table.Campaign.AllColumns, table.Tag.AllColumns).
		// FROM(table.Campaign.INNER_JOIN(

		// )).
		WHERE(AND(
			table.Campaign.OrganisationId.EQ(String(context.Session.User.OrganizationId)),
			table.Campaign.UniqueId.EQ(String(campaignId)),
		)).LIMIT(1)

	campaignResponse := model.Campaign{}
	sqlStatement.Query(database.GetDbInstance(), &campaignResponse)

	if campaignResponse.UniqueId.String() == "" {
		return echo.NewHTTPError(http.StatusNotFound, "Campaign not found")
	}

	// return context.JSON(http.StatusOK, CampaignSchema{
	// 	CreatedAt:             &campaignResponse.CreatedAt,
	// 	Name:                  &campaignResponse.Name,
	// 	Description:           &campaignResponse.Name,
	// 	IsLinkTrackingEnabled: &campaignResponse.,
	// 	TemplateMessageId:    &campaignResponse.MessageTemplateId,
	// 	Status: 			  &campaignResponse.Status,
	// 	ListId:                &campaignResponse.ListId,
	// 	SentAt:                nil,
	// })

	return context.String(http.StatusOK, "OK")

}

func UpdateCampaignStatus(context internal.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}
