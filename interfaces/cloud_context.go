//go:build managed_cloud
// +build managed_cloud

package interfaces

import (
	"github.com/labstack/echo/v4"
	enterprise_utils "github.com/wapikit/wapikit-enterprise/utils"
)

type ContextWithSession struct {
	echo.Context        `json:",inline"`
	App                 App                                              `json:"app,omitempty"`
	Session             ContextSession                                   `json:"session,omitempty"`
	UserIp              string                                           `json:"user_ip,omitempty"`
	UserCountry         string                                           `json:"user_country,omitempty"`
	SubscriptionDetails enterprise_utils.OrganizationSubscriptionDetails `json:"subscription_details,omitempty"`
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
	subscriptionDetails := enterprise_utils.GetOrganizationSubscriptionDetails(context, session.OrganizationId)
	context.SubscriptionDetails = *subscriptionDetails
	return context
}

func (ctx *ContextWithSession) CanUseAiMore() bool {
	return true
}

func (ctx *ContextWithSession) CanCreateMoreContact() bool {
	return true
}

func (ctx *ContextWithSession) CanInviteMoreOrganizationMembers() bool {
	return true
}

func (ctx *ContextWithSession) CanCreateMoreCampaigns() bool {
	return true
}

func (ctx *ContextWithSession) CanUseLinkTracking() bool {
	return true
}

func (ctx *ContextWithSession) CanAccessApi() bool {
	return true
}

func (ctx *ContextWithSession) CanUseLiveTeamInbox() bool {
	return true
}
