package main

import (
	"fmt"
	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
	"github.com/k3env/wgtg/util"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Address  string `mapstructure:"MTWG_API_ADDR"`
	Port     int    `mapstructure:"MTWG_API_PORT"`
	User     string `mapstructure:"MTWG_API_USER"`
	Password string `mapstructure:"MTWG_API_PASS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func main() {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatalf("Error on config loading: %s", err)
	}

	dial := fmt.Sprintf("%s:%d", config.Address, config.Port)
	log.Printf("Connecting to: %s...", dial)
	client, err := routeros.Dial(dial, config.User, config.Password)
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	log.Printf("Connected as %s", config.User)
	defer client.Close()

	r, err := client.Run("/interface/wireguard/print")
	if err != nil {
		log.Fatal(err)
	}
	conf, err := util.Find(r.Re, isManagedConnection)
	if err != nil {
		log.Fatal(err)
	}
	serverKey := conf.Map["public-key"]
	serverPort := conf.Map["listen-port"]

	res, err := client.Run("/ip/cloud/print")
	if err != nil {
		log.Fatal(err)
	}
	cloud := res.Re[0]
	endpoint := fmt.Sprintf("%s:%s", cloud.Map["public-address"], serverPort)

	log.Print(endpoint)
	log.Print(serverKey)
}

func isManagedConnection(p *proto.Sentence) bool {
	return p.Map["comment"] == "wegotik"
}
