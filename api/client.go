package api

import (
	"embed"
	"github.com/go-routeros/routeros"
	"github.com/go-telegram/bot"
	"github.com/k3env/wgtg/config"
	"github.com/k3env/wgtg/wg"
	"github.com/rs/zerolog"
)

type MikrotikAPI struct {
	client      *routeros.Client
	isConnected bool
	logger      *zerolog.Logger
	Interfaces  map[string]*wg.WGInterface
	fs          embed.FS
	bot         *botConfig
	mikrotik    *apiConfig
}

type botConfig struct {
	bindings map[string]string
	admins   []int64
	authCode string
	authType config.BotAuth
}

type apiConfig struct {
	endpoint   string
	user       string
	password   string
	publicAddr string
}

func New() *MikrotikAPI {
	return &MikrotikAPI{
		client:      nil,
		isConnected: false,
		Interfaces:  make(map[string]*wg.WGInterface),
	}
}

func (api *MikrotikAPI) Connect() error {
	if api.logger != nil {
		api.logger.Info().Msgf("Connecting to %s", api.mikrotik.endpoint)
	}
	c, err := routeros.Dial(api.mikrotik.endpoint, api.mikrotik.user, api.mikrotik.password)
	if err != nil {
		return err
	}
	api.client = c
	api.isConnected = true
	if api.logger != nil {
		api.logger.Info().Msgf("Connected as %s", api.mikrotik.user)
	}
	return nil
}

func (api *MikrotikAPI) Disconnect() {
	api.client.Close()
	api.isConnected = false
}

func (api *MikrotikAPI) Load() (err error) {
	err = api.loadInterfaces()
	if err != nil {
		return err
	}
	err = api.loadPeers()
	if err != nil {
		return err
	}
	return nil
}

func (api *MikrotikAPI) BindBot(b *bot.Bot) {
	api.bot.bindings["generate"] = b.RegisterHandler(bot.HandlerTypeMessageText, "/generate", bot.MatchTypePrefix, api.middleCheckAuth(api.bindingNewPeer))
	api.bot.bindings["login"] = b.RegisterHandler(bot.HandlerTypeMessageText, "/login", bot.MatchTypePrefix, api.bindingsLogin)
}

func (api *MikrotikAPI) Unbind(b *bot.Bot) {
	for _, v := range api.bot.bindings {
		b.UnregisterHandler(v)
	}
}

/*
	MikrotikAPI extensions
	Moved to api/extenstions.go
*/

/*
	TG bindings
	Moved to api/bindings.go
*/
