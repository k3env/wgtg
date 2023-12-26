package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/k3env/wgtg/types"
	"github.com/k3env/wgtg/util"
	"log"
	"strings"
	"text/template"
)

func (api *MikrotikAPI) bindingNewPeer(ctx context.Context, b *bot.Bot, update *models.Update) {
	cmd := strings.Split(update.Message.Text, " ")
	if len(cmd) != 3 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Insufficent params",
		})
		return
	}
	wgif := api.Interfaces[cmd[1]]
	if wgif == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Interface not found",
		})
		return
	}
	peer := util.Find[types.WGPeer](wgif.Peers, func(p *types.WGPeer) bool {
		return p.Name == cmd[2]
	})
	if peer == nil {
		peer, _ = api.addNewPeer(wgif, cmd[2], false)
	}
	if peer.PrivateKey == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Cant export config, private key not found",
		})
		return
	}
	t, err := api.fs.ReadFile("tmplts/client.conf")
	if err != nil {
		log.Fatal(err)
	}
	tpl, err := template.New("client").Parse(string(t))
	buff := bytes.NewBufferString("")
	tpl.Execute(buff, peer.ExportConfig())
	fileData := &bot.SendDocumentParams{
		ChatID:   update.Message.Chat.ID,
		Document: &models.InputFileUpload{Filename: fmt.Sprintf("client-%s.conf", strings.Replace(cmd[2], " ", "_", -1)), Data: buff},
		Caption:  fmt.Sprintf("client-%s.conf", strings.Replace(cmd[2], " ", "_", -1)),
	}
	_, err = b.SendDocument(ctx, fileData)
	if err != nil {
		log.Fatal(err)
	}
}
func (api *MikrotikAPI) bindingListPeers(ctx context.Context, b *bot.Bot, update *models.Update) {

}
func (api *MikrotikAPI) bindingRemovePeer(ctx context.Context, b *bot.Bot, update *models.Update) {

}
