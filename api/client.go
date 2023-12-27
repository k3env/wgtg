package api

import (
	"embed"
	"github.com/go-routeros/routeros"
	"github.com/go-telegram/bot"
	"github.com/k3env/wgtg/types"
	"log"
)

type MikrotikAPI struct {
	apiEndpoint string
	apiUser     string
	apiPass     string
	client      *routeros.Client
	isConnected bool
	logger      *log.Logger
	Interfaces  map[string]*types.WGInterface
	tgBindings  map[string]string
	fs          embed.FS
}

func NewAPI(endpoint string, user string, password string) *MikrotikAPI {
	return &MikrotikAPI{
		apiEndpoint: endpoint,
		apiUser:     user,
		apiPass:     password,
		client:      nil,
		isConnected: false,
		tgBindings:  make(map[string]string),
		Interfaces:  make(map[string]*types.WGInterface),
	}
}

func (api *MikrotikAPI) Connect() error {
	c, err := routeros.Dial(api.apiEndpoint, api.apiUser, api.apiPass)
	if err != nil {
		return err
	}
	api.client = c
	api.isConnected = true
	return nil
}

func (api *MikrotikAPI) Disconnect() {
	api.client.Close()
	api.isConnected = false
}

func (api *MikrotikAPI) Load(conf *types.Config) (err error) {
	err = api.loadInterfaces(conf)
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
	api.tgBindings["generate"] = b.RegisterHandler(bot.HandlerTypeMessageText, "/generate", bot.MatchTypePrefix, api.bindingNewPeer)
}

func (api *MikrotikAPI) Unbind(b *bot.Bot) {
	for _, v := range api.tgBindings {
		b.UnregisterHandler(v)
	}
}

func (api *MikrotikAPI) WithLogger(logger *log.Logger) *MikrotikAPI {
	api.logger = logger
	return api
}

func (api *MikrotikAPI) WithTemplateFS(fs embed.FS) *MikrotikAPI {
	api.fs = fs
	return api
}

func (api *MikrotikAPI) WithDefaultLogger() *MikrotikAPI {
	api.logger = log.Default()
	return api
}

/*
	TG bindings
	Moved to api/bindings.go
*/
