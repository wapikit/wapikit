package main

import (
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	file "github.com/knadh/koanf/providers/file"
	posflag "github.com/knadh/koanf/providers/posflag"

	_ "github.com/go-jet/jet/v2/postgres"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
	_ "github.com/wapikit/wapikit/.db-generated/model"
	_ "github.com/wapikit/wapikit/.db-generated/table"
	"github.com/wapikit/wapikit/internal/interfaces"
)

func initConstants() *interfaces.Constants {
	var c interfaces.Constants

	if err := koa.Unmarshal("app", &c); err != nil {
		logger.Error("error loading app config: %v", err)
	}

	if koa.String("environment") == "development" {
		c.IsDevelopment = true
		c.IsProduction = false
	} else {
		c.IsProduction = true
		c.IsDevelopment = false
	}

	c.RootURL = strings.TrimRight("http://127.0.0.0.1:5000/", "/")
	c.SiteName = "Wapikit"
	c.RedisEventChannelName = "ApiServerEvents"
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

func loadConfigFiles(filePaths []string, koa *koanf.Koanf) {
	for _, filePath := range filePaths {
		if err := koa.Load(file.Provider(filePath), toml.Parser()); err != nil {
			if os.IsNotExist(err) {
				logger.Error("config file not found. If there isn't one yet, run --config to generate one.")
			}
			logger.Error("error loading config from file: %v.", err)
		}
		logger.Info("loaded config file %s", koa.All(), nil)
	}
}
