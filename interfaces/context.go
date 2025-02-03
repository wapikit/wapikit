//go:build community_edition
// +build community_edition

package interfaces

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
	"github.com/labstack/echo/v4"
	wapi "github.com/wapikit/wapi.go/pkg/client"
	"github.com/wapikit/wapikit/api/api_types"
	"github.com/wapikit/wapikit/internal/campaign_manager"
	ai_service "github.com/wapikit/wapikit/services/ai_service"
	"github.com/wapikit/wapikit/services/encryption_service"
	"github.com/wapikit/wapikit/services/event_service"
	"github.com/wapikit/wapikit/services/notification_service"
	cache_service "github.com/wapikit/wapikit/services/redis_service"
)

type App struct {
	Db                  *sql.DB
	Redis               *cache_service.RedisClient
	WapiClient          *wapi.Client
	Logger              slog.Logger
	Koa                 *koanf.Koanf
	Fs                  stuffbin.FileSystem
	Constants           *Constants
	CampaignManager     *campaign_manager.CampaignManager
	AiService           *ai_service.AiService
	EncryptionService   *encryption_service.EncryptionService
	NotificationService *notification_service.NotificationService
	EventService        *event_service.EventService
}

type RateLimitConfig struct {
	MaxRequests    int   `json:"maxRequests"`
	WindowTimeInMs int64 `json:"windowTime"`
}

type RouteMetaData struct {
	PermissionRoleLevel  api_types.UserPermissionLevelEnum `json:"permissionRoleLevel"`
	RequiredPermission   []api_types.RolePermissionEnum    `json:"requiredPermission"`
	RateLimitConfig      RateLimitConfig                   `json:"rateLimitConfig"`
	RequiredFeatureFlags []string                          `json:"requiredFeatureFlags"`
}

type Route struct {
	Path                    string `json:"path"`
	Method                  string `json:"method"`
	Handler                 Handler
	IsAuthorizationRequired bool
	MetaData                RouteMetaData `json:"metaData"`
}

type ApiController interface {
	Register(server *echo.Echo)
	GetControllerName() string
}

type Handler interface {
	Handle(context echo.Context) error
}

type HandlerWithoutSession func(context ContextWithoutSession) error

func (eh HandlerWithoutSession) Handle(context echo.Context) error {
	return eh(context.(ContextWithoutSession))
}

type HandlerWithSession func(context ContextWithSession) error

func (ch HandlerWithSession) Handle(context echo.Context) error {
	return ch(context.(ContextWithSession))
}

// only for cloud edition compatibility, will never be used for real checks
type SubscriptionDetails struct {
	ValidTill             time.Time `form:"valid_till"`
	PlanType              string    `form:"plan_type"`
	LastPurchaseOn        time.Time `form:"last_purchase_on"`
	Validity              string    `form:"validity"`
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

func (ctx *ContextWithoutSession) GetSessionDetailsIfAuthenticated() (ContextSession, bool) {
	return ContextSession{}, false
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

func (ctx *ContextWithSession) IsAiLimitReached() bool {
	return true
}

func (ctx *ContextWithSession) IsContactCreationLimitReached() bool {
	return true
}

func (ctx *ContextWithSession) IsOrganizationMemberLimitReached() bool {
	return true
}

func (ctx *ContextWithSession) IsActiveCampaignLimitReached() bool {
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
