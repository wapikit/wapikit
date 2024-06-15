package campaign_service

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/database"
	"github.com/sarthakjdev/wapikit/internal/api_types"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
	"github.com/sarthakjdev/wapikit/services"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/sarthakjdev/wapikit/.db-generated/model"
	table "github.com/sarthakjdev/wapikit/.db-generated/table"
)

type CampaignService struct {
	services.BaseService `json:"-,inline"`
}

func NewCampaignService() *CampaignService {
	return &CampaignService{
		BaseService: services.BaseService{
			Name:        "Campaign Service",
			RestApiPath: "/api/campaign",
			Routes: []interfaces.Route{
				{
					Path:                    "/campaigns",
					Method:                  http.MethodGet,
					Handler:                 GetCampaigns,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/campaigns",
					Method:                  http.MethodPost,
					Handler:                 CreateNewCampaign,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/campaigns/:id",
					Method:                  http.MethodGet,
					Handler:                 GetCampaignByID,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/campaigns/:id",
					Method:                  http.MethodPut,
					Handler:                 UpdateCampaignById,
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/campaigns/:id",
					Method:                  http.MethodDelete,
					Handler:                 DeleteCampaignById,
					IsAuthorizationRequired: true,
				},
			},
		},
	}
}

func GetCampaigns(context interfaces.CustomContext) error {
	params := new(api_types.GetCampaignsParams)
	if err := context.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return context.String(http.StatusOK, "OK")
}

func CreateNewCampaign(context interfaces.CustomContext) error {
	payload := new(api_types.CreateCampaignJSONRequestBody)
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return context.String(http.StatusOK, "OK")
}

func GetCampaignByID(context interfaces.CustomContext) error {
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

func UpdateCampaignById(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}

func DeleteCampaignById(context interfaces.CustomContext) error {
	return context.String(http.StatusOK, "OK")
}
