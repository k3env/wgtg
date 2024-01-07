package api

import (
	"embed"
	"fmt"
	"github.com/k3env/wgtg/config"
	"github.com/rs/zerolog"
	"os"
)

func (api *MikrotikAPI) WithLogger(logger *zerolog.Logger) *MikrotikAPI {
	api.logger = logger
	return api
}

func (api *MikrotikAPI) WithTemplateFS(fs embed.FS) *MikrotikAPI {
	api.fs = fs
	return api
}

func (api *MikrotikAPI) WithDefaultLogger() *MikrotikAPI {
	zlOut := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(zlOut)
	api.logger = &logger
	return api
}

func (api *MikrotikAPI) WithConfig(config *config.Config) *MikrotikAPI {
	api.mikrotik = &apiConfig{
		endpoint:   fmt.Sprintf("%s:%d", config.Address, config.Port),
		user:       config.User,
		password:   config.Password,
		publicAddr: config.PublicAddress,
	}
	api.bot = &botConfig{
		bindings: make(map[string]string),
		admins:   config.BotAdmins,
		authCode: config.BotAdminCode,
		authType: config.AuthType,
	}
	return api
}
