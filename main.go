package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"github.com/go-telegram/bot"
	api2 "github.com/k3env/wgtg/api"
	"github.com/k3env/wgtg/types"
	"log"
	"os"
	"os/signal"
)

var err error
var config types.Config

//go:embed all:tmplts
var tmplts embed.FS

func main() {
	var cfgFlag = flag.String("config", "app.env", "")
	flag.Parse()

	config, err = types.LoadConfig(*cfgFlag)
	if err != nil {
		log.Fatalf("Error on config loading: %s", err)
	}

	dial := fmt.Sprintf("%s:%d", config.Address, config.Port)

	log.Printf("Connecting to: %s...", dial)
	api := api2.NewAPI(dial, config.User, config.Password).WithDefaultLogger().WithTemplateFS(tmplts)
	err = api.Connect()
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	log.Printf("Connected as %s", config.User)
	defer api.Disconnect()

	err = api.Load(&config)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	tg, err := bot.New(config.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	api.BindBot(tg)
	defer api.Unbind(tg)
	tg.Start(ctx)
}
