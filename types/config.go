package types

import (
	"github.com/spf13/viper"
)

type Config struct {
	Address       string `mapstructure:"MTWG_API_ADDR"`
	Port          int    `mapstructure:"MTWG_API_PORT"`
	User          string `mapstructure:"MTWG_API_USER"`
	Password      string `mapstructure:"MTWG_API_PASS"`
	BotToken      string `mapstructure:"MTWG_TG_TOKEN"`
	BotAdmin      int    `mapstructure:"MTWG_TG_ADMIN"`
	PublicAddress string `mapstructure:"MTWG_WG_PUBLIC_IP"`
}

func LoadConfig(path string) (config *Config, err error) {
	if path != "" {
		viper.SetConfigFile(path)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignoring error
		} else {
			return nil, err
		}
	}
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
