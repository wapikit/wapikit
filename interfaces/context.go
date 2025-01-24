//go:build community_edition
// +build community_edition

package interfaces

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// only for cloud edition compatibility, will never be used for real checks
type SubscriptionDetails struct {
	ValidTill             time.Time `form:"valid_till"`
	PlanType              string    `form:"plan_type"`
	GatewaySubscriptionId string    `form:"gateway_subscription_id"`
	DbSubscriptionId      uuid.UUID `form:"db_subscription_id"`
}

type ContextWithSession struct {
	echo.Context        `json:",inline"`
	App                 App                 `json:"app,omitempty"`
	Session             ContextSession      `json:"session,omitempty"`
	UserIp              string              `json:"user_ip,omitempty"`
	UserCountry         string              `json:"user_country,omitempty"`
	SubscriptionDetails SubscriptionDetails `json:"subscription_details,omitempty"`
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
	return ContextWithSession{
		Context:     ctx,
		App:         app,
		Session:     session,
		UserIp:      userIp,
		UserCountry: userCountry,
	}
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
