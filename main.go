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

// var client *routeros.Client
var err error
var config types.Config

//var res *routeros.Reply
//var wgInterfaces = make(map[string]*types.WGInterface)

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
	//client, err = routeros.Dial(dial, config.User, config.Password)
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	log.Printf("Connected as %s", config.User)
	defer api.Disconnect()

	err = api.Load()
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Close()

	//publicIp, err := getPublicIP(client)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//res, err = client.RunArgs([]string{"/interface/wireguard/print"})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if len(res.Re) == 0 {
	//	log.Fatal("Wireguard interfaces not found!")
	//}
	//for _, cfg := range res.Re {
	//	var wgif types.WGInterface
	//	err = wgif.Parse(*cfg)
	//	if err != nil {
	//		log.Printf("Error while parsing %s interface config: %s", cfg.Map["name"], err)
	//		continue
	//	}
	//	ip, network, err := getInterfaceIp(cfg.Map["name"])
	//	if err != nil {
	//		log.Printf("Error while parsing %s interface config: %s", cfg.Map["name"], err)
	//		continue
	//	}
	//	err = wgif.SetNetworks(ip, network, publicIp.String())
	//	if err != nil {
	//		log.Printf("Error while parsing %s interface config: %s", wgif.Interface, err)
	//		continue
	//	}
	//	wgInterfaces[wgif.Interface] = &wgif
	//}
	//res, err = client.Run("/interface/wireguard/peers/print")
	//if err != nil {
	//	log.Fatalf("Error %s", err)
	//}
	//for _, peer := range res.Re {
	//	data := peer.Map
	//	ifName := data["interface"]
	//
	//	if wgInterfaces[ifName].Interface == ifName {
	//		parts := strings.Split(data["comment"], "|")
	//		user := ""
	//		var ip net.IP
	//		if len(parts) >= 2 {
	//			ip = net.ParseIP(parts[1])
	//		}
	//		if len(parts) >= 1 {
	//			user = parts[0]
	//		}
	//		shared := data["preshared-key"]
	//		key := data["public-key"]
	//		ipList := strings.Split(data["allowed-address"], ",")
	//		wgPeer := &types.WGPeer{
	//			Interface:  wgInterfaces[ifName],
	//			Name:       user,
	//			IP:         ip,
	//			PrivateKey: "",
	//			PublicKey:  key,
	//			SharedKey:  shared,
	//			AllowedIPs: ipList,
	//		}
	//		wgInterfaces[ifName].ImportPeer(*wgPeer)
	//	}
	//}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	tg, err := bot.New(config.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	api.BindBot(tg)
	defer api.Unbind(tg)
	//h := tg.RegisterHandler(bot.HandlerTypeMessageText, "/generate", bot.MatchTypePrefix, addPeerHandler)
	//defer tg.UnregisterHandler(h)
	tg.Start(ctx)
}

//func addPeerHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
//	cmd := strings.Split(update.Message.Text, " ")
//	if len(cmd) != 3 {
//		b.SendMessage(ctx, &bot.SendMessageParams{
//			ChatID: update.Message.Chat.ID,
//			Text:   "Insufficent params",
//		})
//		return
//	}
//	if wgInterfaces[cmd[1]] == nil {
//		b.SendMessage(ctx, &bot.SendMessageParams{
//			ChatID: update.Message.Chat.ID,
//			Text:   "Interface not found",
//		})
//		return
//	}
//	peer := util.Find[types.WGPeer](wgInterfaces[cmd[1]].Peers, func(p *types.WGPeer) bool {
//		return p.Name == cmd[2]
//	})
//	if peer == nil {
//		peer, _ = addNewPeer(*wgInterfaces[cmd[1]], cmd[2], false)
//	}
//	if peer.PrivateKey == "" {
//		b.SendMessage(ctx, &bot.SendMessageParams{
//			ChatID: update.Message.Chat.ID,
//			Text:   "Cant export config, private key not found",
//		})
//		return
//	}
//	t, err := tmplts.ReadFile("tmplts/client.conf")
//	if err != nil {
//		log.Fatal(err)
//	}
//	tpl, err := template.New("client").Parse(string(t))
//	buff := bytes.NewBufferString("")
//	tpl.Execute(buff, peer.ExportConfig())
//	fileData := &bot.SendDocumentParams{
//		ChatID:   update.Message.Chat.ID,
//		Document: &models.InputFileUpload{Filename: fmt.Sprintf("client-%s.conf", strings.Replace(cmd[2], " ", "_", -1)), Data: buff},
//		Caption:  fmt.Sprintf("client-%s.conf", strings.Replace(cmd[2], " ", "_", -1)),
//	}
//	_, err = b.SendDocument(ctx, fileData)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//}
//
//func addNewPeer(wgInterface types.WGInterface, name string, usePsk bool) (*types.WGPeer, error) {
//	peer, _ := wgInterface.AddPeer(name, usePsk)
//	cmd := fmt.Sprintf("/interface/wireguard/peers/add =public-key=%s =interface=%s =allowed-address=%s =comment=%s|%s", peer.PublicKey, wgInterface.Interface, strings.Join(peer.AllowedIPs, ","), peer.Name, peer.IP.String())
//	args := strings.Split(cmd, " ")
//	_, err = client.RunArgs(args)
//	if err != nil {
//		return nil, err
//	}
//	return peer, nil
//}
//
//func getPublicIP(client *routeros.Client) (net.IP, error) {
//	res, err := client.Run("/ip/cloud/print")
//	if err != nil {
//		return nil, err
//	}
//	addr := res.Re[0].Map["public-address"]
//	ip := net.ParseIP(addr)
//	return ip, nil
//}
//
//func getInterfaceIp(ifname string) (net.IP, *net.IPNet, error) {
//	args := []string{"/ip/address/print", fmt.Sprintf("?=interface=%s", ifname)}
//	res, err = client.RunArgs(args)
//	if err != nil {
//		return nil, nil, err
//	}
//	if len(res.Re) == 0 {
//		return nil, nil, errors.IPNotFoundError
//	}
//	addr := res.Re[0].Map["address"]
//	return net.ParseCIDR(addr)
//}
