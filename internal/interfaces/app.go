package interfaces

import (
	"database/sql"
	"log/slog"

	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

type constants struct {
	SiteName      string `koanf:"site_name"`
	RootURL       string `koanf:"root_url"`
	LogoURL       string `koanf:"logo_url"`
	FaviconURL    string `koanf:"favicon_url"`
	AdminUsername []byte `koanf:"admin_username"`
	AdminPassword []byte `koanf:"admin_password"`
	IsDevelopment bool
	IsProduction  bool
}

type App struct {
	Db        *sql.DB
	Logger    slog.Logger
	Koa       *koanf.Koanf
	Fs        stuffbin.FileSystem
	Constants *constants
}
