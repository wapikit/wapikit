package user_controller

import (
	"net/http"

	"github.com/go-jet/jet/qrm"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
	controller "github.com/wapikit/wapikit/api/controllers"
	"github.com/wapikit/wapikit/internal/api_types"
	"github.com/wapikit/wapikit/internal/core/utils"
	"github.com/wapikit/wapikit/internal/interfaces"
)

type UserController struct {
	controller.BaseController `json:"-,inline"`
}

func NewUserController() *UserController {
	return &UserController{
		BaseController: controller.BaseController{
			Name:        "User Controller",
			RestApiPath: "/api",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/user",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getUser),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60 * 60,
						},
					},
				},
				{
					Path:                    "/api/user",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(updateUser),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60 * 60,
						},
					},
				},
				{
					Path:                    "/api/user/notifications",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getNotifications),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60 * 60,
						},
					},
				},
			},
		},
	}
}

func getUser(context interfaces.ContextWithSession) error {
	userUuid, err := uuid.Parse(context.Session.User.UniqueId)
	if err != nil {
		return context.String(http.StatusInternalServerError, "Error parsing user UUID")
	}

	orgUuid, err := uuid.Parse(context.Session.User.OrganizationId)

	if err != nil {
		// it might be possible that user have not joined any organization yet
	}

	userQuery := SELECT(
		table.User.AllColumns,
		table.Organization.AllColumns,
		table.WhatsappBusinessAccount.AllColumns,
		table.OrganizationMember.AllColumns,
	).
		FROM(
			table.User.
				LEFT_JOIN(table.OrganizationMember, table.OrganizationMember.OrganizationId.EQ(UUID(orgUuid)).
					AND(table.OrganizationMember.UserId.EQ(table.User.UniqueId))).
				LEFT_JOIN(table.Organization, table.Organization.UniqueId.EQ(table.OrganizationMember.OrganizationId)).
				LEFT_JOIN(table.WhatsappBusinessAccount, table.WhatsappBusinessAccount.OrganizationId.EQ(table.Organization.UniqueId)),
		).
		WHERE(
			table.User.UniqueId.EQ(UUID(userUuid)).
				AND(
					table.User.Email.EQ(String(context.Session.User.Email)),
				),
		).LIMIT(1)

	type UserWithOrgDetails struct {
		User                    model.User
		Organization            model.Organization
		OrganizationMember      model.OrganizationMember
		WhatsappBusinessAccount model.WhatsappBusinessAccount
	}

	user := UserWithOrgDetails{}

	userQuery.Query(context.App.Db, &user)

	isOwner := false

	if user.OrganizationMember.AccessLevel == model.UserPermissionLevelEnum_Owner {
		isOwner = true
	}

	currentPermissionLevel := api_types.UserPermissionLevelEnum(user.OrganizationMember.AccessLevel)

	// find the current logged in organization
	response := api_types.GetUserResponseSchema{
		User: api_types.UserSchema{
			CreatedAt:                      user.User.CreatedAt,
			Name:                           user.User.Name,
			Email:                          user.User.Email,
			Username:                       user.User.Username,
			UniqueId:                       context.Session.User.UniqueId,
			ProfilePicture:                 user.User.ProfilePictureUrl,
			IsOwner:                        isOwner,
			CurrentOrganizationAccessLevel: &currentPermissionLevel,
		},
	}

	if user.Organization.UniqueId.String() != uuid.Nil.String() {
		response.User.Organization = &api_types.OrganizationSchema{
			Name:        user.Organization.Name,
			CreatedAt:   user.Organization.CreatedAt,
			UniqueId:    user.Organization.UniqueId.String(),
			FaviconUrl:  &user.Organization.FaviconUrl,
			LogoUrl:     user.Organization.LogoUrl,
			WebsiteUrl:  user.Organization.WebsiteUrl,
			Description: user.Organization.Description,
		}
	}

	if user.WhatsappBusinessAccount.AccessToken != "" {
		response.User.Organization.WhatsappBusinessAccountDetails = &api_types.WhatsAppBusinessAccountDetailsSchema{
			AccessToken:       user.WhatsappBusinessAccount.AccessToken,
			BusinessAccountId: user.WhatsappBusinessAccount.AccountId,
			WebhookSecret:     user.WhatsappBusinessAccount.WebhookSecret,
		}
	}

	return context.JSON(http.StatusOK, response)
}

func updateUser(context interfaces.ContextWithSession) error {
	userUuid, err := uuid.Parse(context.Session.User.UniqueId)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Error parsing user UUID")
	}

	payload := new(api_types.UpdateUserSchema)
	var user model.User
	updateUserQuery := table.User.
		UPDATE(table.User.Name, table.User.ProfilePictureUrl).
		SET(payload.Name, payload.ProfilePicture).
		WHERE(table.User.UniqueId.EQ(UUID(userUuid)))

	err = updateUserQuery.QueryContext(context.Request().Context(), context.App.Db, &user)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Error updating user")
	}

	isUpdated := false
	responseToReturn := api_types.UpdateUserResponseSchema{
		IsUpdated: isUpdated,
	}

	return context.JSON(http.StatusOK, responseToReturn)
}

func getNotifications(context interfaces.ContextWithSession) error {

	params := new(api_types.GetUserNotificationsParams)

	err := utils.BindQueryParams(context, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	page := params.Page
	limit := params.PerPage

	var dest []struct {
		TotalNotifications int `json:"totalNotifications"`
		model.Notification
	}

	notificationsQuery := SELECT(
		table.Notification.AllColumns,
		COUNT(table.Notification.UniqueId).OVER().AS("totalNotifications"),
	).
		FROM(table.Notification).
		WHERE(table.Notification.UserId.EQ(UUID(uuid.MustParse(context.Session.User.UniqueId)))).
		ORDER_BY(table.Notification.CreatedAt.DESC()).
		LIMIT(limit).
		OFFSET((page - 1) * limit)

	err = notificationsQuery.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			total := 0
			notifications := make([]api_types.NotificationSchema, 0)
			return context.JSON(http.StatusOK, api_types.GetUserNotificationsResponseSchema{
				Notifications: notifications,
				UnreadCount:   0,
				PaginationMeta: api_types.PaginationMeta{
					Page:    page,
					PerPage: limit,
					Total:   total,
				},
			})
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	// ! return the notifications to the user

	totalNotifications := 0

	if len(dest) > 0 {
		totalNotifications = dest[0].TotalNotifications
	}

	response := api_types.GetUserNotificationsResponseSchema{
		Notifications: []api_types.NotificationSchema{},
		PaginationMeta: api_types.PaginationMeta{
			Page:    page,
			PerPage: limit,
			Total:   0,
		},
	}

	for _, notification := range dest {
		response.Notifications = append(response.Notifications, api_types.NotificationSchema{
			UniqueId:    notification.UniqueId.String(),
			CreatedAt:   notification.CreatedAt,
			Title:       notification.Title,
			Description: notification.Description,
			CtaUrl:      notification.CtaUrl,
			Type:        *notification.Type,
		})
	}

	response.PaginationMeta.Total = totalNotifications

	return context.JSON(http.StatusOK, response)
}

func DeleteAccountStepOne(context interfaces.ContextWithSession) error {
	// ! generate a deletion token here
	// ! send the link to delete account with token in it to the user email
	// ! get the user details from the token from the frontend and then check if the account is even deletable or not because if a user owns a organization he/she must need to transfer the ownership to someone else before deleting the account
	return context.String(http.StatusOK, "OK")
}

func DeleteAccountStetTwo(context interfaces.ContextWithSession) error {
	return context.String(http.StatusOK, "OK")
}
