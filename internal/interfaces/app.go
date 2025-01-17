package interfaces

import (
	"database/sql"
	"log/slog"

	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
	wapi "github.com/wapikit/wapi.go/pkg/client"
	"github.com/wapikit/wapikit/internal/campaign_manager"
	"github.com/wapikit/wapikit/internal/services/ai_service"
	cache_service "github.com/wapikit/wapikit/internal/services/redis_service"
)

type Constants struct {
	SiteName              string `koanf:"site_name"`
	RootURL               string `koanf:"root_url"`
	LogoURL               string `koanf:"logo_url"`
	FaviconURL            string `koanf:"favicon_url"`
	AdminUsername         []byte `koanf:"admin_username"`
	AdminPassword         []byte `koanf:"admin_password"`
	IsDevelopment         bool
	IsProduction          bool
	RedisEventChannelName string `koanf:"redis_event_channel_name"`
	IsDebugModeEnabled    bool
}

type App struct {
	Db              *sql.DB
	Redis           *cache_service.RedisClient
	WapiClient      *wapi.Client
	Logger          slog.Logger
	Koa             *koanf.Koanf
	Fs              stuffbin.FileSystem
	Constants       *Constants
	CampaignManager *campaign_manager.CampaignManager
	AiService       *ai_service.AiService
	// ! TODO: add some api server event utility so anybody api server event can be published easily.
}
