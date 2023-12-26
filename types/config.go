package types

import (
	"github.com/spf13/viper"
)

type Config struct {
	Address  string `mapstructure:"MTWG_API_ADDR"`
	Port     int    `mapstructure:"MTWG_API_PORT"`
	User     string `mapstructure:"MTWG_API_USER"`
	Password string `mapstructure:"MTWG_API_PASS"`
	BotToken string `mapstructure:"MTWG_TG_TOKEN"`
	BotAdmin int    `mapstructure:"MTWG_TG_ADMIN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
