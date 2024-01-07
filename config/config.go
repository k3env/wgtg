package config

import (
	"errors"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type BotAuth string

const (
	BotAuthStatic   = "static"
	BotAuthDynamic  = "dynamic"
	BotAuthDisabled = "disabled"
	BotAuthBoth     = "both"
)

type Config struct {
	Address       string  `mapstructure:"MTWG_API_ADDR"`
	Port          int     `mapstructure:"MTWG_API_PORT"`
	User          string  `mapstructure:"MTWG_API_USER"`
	Password      string  `mapstructure:"MTWG_API_PASS"`
	BotToken      string  `mapstructure:"MTWG_TG_TOKEN"`
	BotAdminCode  string  `mapstructure:"MTWG_TG_ADMIN_CODE"`
	PublicAddress string  `mapstructure:"MTWG_WG_PUBLIC_IP"`
	BotAdmins     []int64 `mapstructure:"MTWG_TG_ADMINS"`
	AuthType      BotAuth `mapstructure:"MTWG_TG_AUTH"`
}

func Load(path string) (config *Config, err error) {
	if path != "" {
		viper.SetConfigFile(path)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if !errors.Is(err, viper.ConfigFileNotFoundError{}) {
			return nil, err
		}
	}
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	if config.BotAdminCode == "" {
		config.BotAdminCode = uuid.New().String()
	}
	return
}
