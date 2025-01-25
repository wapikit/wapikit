//go:build managed_cloud
// +build managed_cloud

package interfaces

import (
	"github.com/labstack/echo/v4"

	enterprise_interfaces "github.com/wapikit/wapikit-enterprise/interfaces"
	quota_service "github.com/wapikit/wapikit-enterprise/services/quota"
	subscription_service "github.com/wapikit/wapikit-enterprise/services/subscription"
)

type ContextWithSession struct {
	echo.Context        `json:",inline"`
	App                 App                                                    `json:"app,omitempty"`
	Session             ContextSession                                         `json:"session,omitempty"`
	UserIp              string                                                 `json:"user_ip,omitempty"`
	UserCountry         string                                                 `json:"user_country,omitempty"`
	SubscriptionDetails *enterprise_interfaces.OrganizationSubscriptionDetails `json:"subscription_details,omitempty"`
}

type ContextWithoutSession struct {
	echo.Context `json:",inline"`
	App          App    `json:"app,omitempty"`
	UserIp       string `json:"user_ip,omitempty"`
	UserCountry  string `json:"user_country,omitempty"`
}

func BuildContextWithoutSession(ctx echo.Context, app App, userIp, userCountry string) ContextWithoutSession {
	return ContextWithoutSession{
		Context:     ctx,
		App:         app,
		UserIp:      userIp,
		UserCountry: userCountry,
	}
}

func BuildContextWithSession(ctx echo.Context, app App, session ContextSession, userIp, userCountry string) ContextWithSession {
	context := ContextWithSession{
		Context:     ctx,
		App:         app,
		Session:     session,
		UserIp:      userIp,
		UserCountry: userCountry,
	}
	subscriptionServiceInstance := subscription_service.GetSubscriptionServiceInstance()
	subscriptionDetails := subscriptionServiceInstance.GetOrganizationSubscriptionDetails(context.Request().Context(), context.App.Db, context.Session.User.OrganizationId, context.App.Constants.IsCommunityEdition)
	context.SubscriptionDetails = subscriptionDetails
	return context
}

func (ctx *ContextWithSession) IsAiLimitReached() bool {
	subscriptionService := quota_service.GetQuotaServiceInstance()
	return subscriptionService.IsAiLimitReached(ctx.SubscriptionDetails)
}

func (ctx *ContextWithSession) IsContactCreationLimitReached() bool {
	quotaInstance := quota_service.GetQuotaServiceInstance()
	return quotaInstance.IsContactCreationLimitReached(ctx.Request().Context(), ctx.Session.User.OrganizationId, ctx.App.Db, ctx.SubscriptionDetails)
}

func (ctx *ContextWithSession) IsOrganizationMemberLimitReached() bool {
	quotaInstance := quota_service.GetQuotaServiceInstance()
	return quotaInstance.IsOrganizationMemberLimitReached(ctx.Request().Context(), ctx.Session.User.OrganizationId, ctx.App.Db, ctx.SubscriptionDetails)
}

func (ctx *ContextWithSession) IsActiveCampaignLimitReached() bool {
	quotaInstance := quota_service.GetQuotaServiceInstance()
	return quotaInstance.IsActiveCampaignLimitReached(ctx.Request().Context(), ctx.Session.User.OrganizationId, ctx.App.Db, ctx.SubscriptionDetails)
}

func (ctx *ContextWithSession) CanUseLinkTracking() bool {
	quotaInstance := quota_service.GetQuotaServiceInstance()
	return quotaInstance.CanUseLinkTracking()
}

func (ctx *ContextWithSession) CanAccessApi() bool {
	quotaInstance := quota_service.GetQuotaServiceInstance()
	return quotaInstance.CanAccessApi()
}

func (ctx *ContextWithSession) CanUseLiveTeamInbox() bool {
	quotaInstance := quota_service.GetQuotaServiceInstance()
	return quotaInstance.CanUseLiveTeamInbox()
}
