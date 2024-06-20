package main

import (
	"html/template"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/knadh/koanf/parsers/toml"
	file "github.com/knadh/koanf/providers/file"
	posflag "github.com/knadh/koanf/providers/posflag"

	_ "github.com/go-jet/jet/v2/postgres"
	"github.com/knadh/koanf/v2"
	echo "github.com/labstack/echo/v4"
	_ "github.com/sarthakjdev/wapikit/.db-generated/model"
	_ "github.com/sarthakjdev/wapikit/.db-generated/table"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
	flag "github.com/spf13/pflag"
)

type tplRenderer struct {
	templates  *template.Template
	SiteName   string
	RootURL    string
	LogoURL    string
	FaviconURL string
}

type tplData struct {
	SiteName   string
	RootURL    string
	LogoURL    string
	FaviconURL string
	Data       interface{}
}

// Render executes and renders a template for echo.
func (t *tplRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, tplData{
		SiteName:   t.SiteName,
		RootURL:    t.RootURL,
		LogoURL:    t.LogoURL,
		FaviconURL: t.FaviconURL,
		Data:       data,
	})
}

func initConstants() *interfaces.Constants {
	var c interfaces.Constants

	if err := koa.Unmarshal("app", &c); err != nil {
		logger.Error("error loading app config: %v", err)
	}

	if koa.String("env") == "development" {
		c.IsDevelopment = true
		c.IsProduction = false
	} else {
		c.IsProduction = true
		c.IsDevelopment = false
	}

	c.RootURL = strings.TrimRight("http://127.0.0.0.1:5000/", "/")
	c.SiteName = "Wapikit"
	logger.Info("loading app constants %s", c, nil)
	return &c
}

func initSettings(app *interfaces.App) {
	// get the settings from the DB and load it into the koanf

	// var out map[string]interface{}
	// if err := json.Unmarshal(s, &out); err != nil {
	// 	app.Logger.Error("error unmarshalling settings from DB: %v", err)
	// }
	// if err := app.Koa.Load(confmap.Provider(out, "."), nil); err != nil {
	// 	app.Logger.Error("error parsing settings from DB: %v", err)
	// }

}

func initFlags() {
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		// ! TODO: trigger the --help command here
		logger.Info(f.FlagUsages())
		os.Exit(0)
	}

	// Register the command line flags.
	f.StringSlice("config", []string{"config-dev.toml"},
		"path to one or more config files (will be merged in order)")
	f.Bool("install", false, "setup database (first time)")
	f.Bool("idempotent", false, "make --install run only if the database isn't already setup")
	f.Bool("upgrade", false, "upgrade database to the current version")
	f.Bool("version", false, "show current version of the build")
	f.Bool("yes", false, "assume 'yes' to prompts during --install/upgrade")
	f.Bool("db-migrate", false, "apply database migrations")
	f.Bool("db-apply", false, "apply database migrations")

	if err := f.Parse(os.Args[1:]); err != nil {
		logger.Error("error loading flags: %v", err)
	}

	if err := koa.Load(posflag.Provider(f, ".", koa), nil); err != nil {
		logger.Error("error loading config: %v", err)
	}
}

func initDb() {
	// check for the database connection
	// check if the database if the database is already setup
	// if not, then setup the database

	// if the database is already setup, then check for the migrations
	// if the migrations are not applied, then apply the migrations

	// create a default user and organisation
}

func handleMigrationApply() {
	// ask for the confirmation here first
}

func checkDbState() {
}

func initMigrate() {
}

func loadConfigFiles(filePaths []string, koa *koanf.Koanf) {
	for _, filePath := range filePaths {
		logger.Info("reading config: %s", filePath, nil)
		if err := koa.Load(file.Provider(filePath), toml.Parser()); err != nil {
			if os.IsNotExist(err) {
				logger.Error("config file not found. If there isn't one yet, run --config to generate one.")
			}
			logger.Error("error loading config from file: %v.", err)
		}
	}

}

func initTplFuncs(cs *interfaces.Constants) template.FuncMap {
	funcs := template.FuncMap{
		"RootURL": func() string {
			return cs.RootURL
		},
		"LogoURL": func() string {
			return cs.LogoURL
		},
		"Date": func(layout string) string {
			if layout == "" {
				layout = time.ANSIC
			}
			return time.Now().Format(layout)
		},
	}

	for k, v := range sprig.GenericFuncMap() {
		funcs[k] = v
	}

	return funcs
}
