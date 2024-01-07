package main

import (
	"context"
	"embed"
	"flag"
	"github.com/go-telegram/bot"
	tikapi "github.com/k3env/wgtg/api"
	"github.com/k3env/wgtg/config"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
)

var err error
var cfg *config.Config

//go:embed all:tmplts
var tmplts embed.FS

func main() {
	zo := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "02.01.2006 15:04"}
	zl := zerolog.New(zo).Level(zerolog.InfoLevel).With().Timestamp().Logger()

	var cfgFlag = flag.String("config", "", "")
	flag.Parse()

	cfg, err = config.Load(*cfgFlag)
	if err != nil {
		zl.Fatal().Err(err).Msg("Error on config loading")
	}
	if cfg.AuthType == config.BotAuthDisabled {
		zl.Warn().Msg("Authentication disabled, anyone can create wg profiles")
	}
	if len(cfg.BotAdmins) == 0 && cfg.AuthType != config.BotAuthDisabled {
		if cfg.AuthType == config.BotAuthStatic {

		} else {
			zl.Info().Msg("No admin accounts found\nYou need register one")
		}
	}
	if cfg.AuthType == config.BotAuthDynamic || cfg.AuthType == config.BotAuthBoth {
		zl.Info().Msgf("Admin registration token is: %s", cfg.BotAdminCode)
	}

	api := tikapi.New().WithLogger(&zl).WithTemplateFS(tmplts).WithConfig(cfg)
	err = api.Connect()
	if err != nil {
		zl.Fatal().Err(err).Msg("Connection to Mikrotik error")
	}
	defer api.Disconnect()

	err = api.Load()
	if err != nil {
		zl.Fatal().Err(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	tg, err := bot.New(cfg.BotToken)
	if err != nil {
		zl.Fatal().Err(err)
	}

	api.BindBot(tg)
	defer api.Unbind(tg)
	tg.Start(ctx)
}
