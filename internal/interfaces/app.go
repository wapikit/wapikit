package interfaces

import (
	"database/sql"
	"log/slog"

	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
	wapi "github.com/sarthakjdev/wapi.go/pkg/client"
	"github.com/sarthakjdev/wapikit/internal"
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
}

type App struct {
	Db         *sql.DB
	Redis      *internal.RedisClient
	WapiClient *wapi.Client
	Logger     slog.Logger
	Koa        *koanf.Koanf
	Fs         stuffbin.FileSystem
	Constants  *Constants
}
