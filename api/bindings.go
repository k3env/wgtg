package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/k3env/wgtg/config"
	"github.com/k3env/wgtg/util"
	"github.com/k3env/wgtg/wg"
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
	peer := util.Find[wg.WGPeer](wgif.Peers, func(p *wg.WGPeer) bool {
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
	if err != nil && api.logger != nil {
		api.logger.Fatal().Err(err)
	}
	if api.logger != nil {
		api.logger.Info().
			Str("config", cmd[2]).
			Int64("sender", update.Message.From.ID).
			Str("name", fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)).
			Str("user", update.Message.From.Username).
			Msg("Config issued")
	}
}
func (api *MikrotikAPI) bindingListPeers(ctx context.Context, b *bot.Bot, update *models.Update) {

}
func (api *MikrotikAPI) bindingRemovePeer(ctx context.Context, b *bot.Bot, update *models.Update) {

}

func (api *MikrotikAPI) bindingsLogin(ctx context.Context, b *bot.Bot, update *models.Update) {
	cmd := strings.Split(update.Message.Text, " ")
	if api.bot.authType == config.BotAuthStatic || api.bot.authType == config.BotAuthDisabled {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Authentication by code disabled",
		})
		return
	}
	if len(cmd) != 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Insufficent params",
		})
		return
	}
	code := cmd[1]
	if api.bot.authCode != code {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Wrong key",
		})
		if api.logger != nil {
			api.logger.Warn().
				Int64("sender", update.Message.From.ID).
				Str("name", fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)).
				Str("user", update.Message.From.Username).
				Msg("Wrong key")
		}
		return
	}
	if util.Has(api.bot.admins, func(i int64) bool {
		return update.Message.From.ID == i
	}) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "You're already admin",
		})
		return
	}
	api.bot.admins = append(api.bot.admins, update.Message.From.ID)
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Now you're admin"})
	if api.logger != nil {
		api.logger.Info().Int64("sender", update.Message.From.ID).
			Str("name", fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)).
			Str("user", update.Message.From.Username).
			Msg("New admin")
	}
}

func (api *MikrotikAPI) middleCheckAuth(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if api.bot.authType == config.BotAuthDisabled {
			if api.logger != nil {
				api.logger.Info().Msg("Auth disabled, passing")
			}
			next(ctx, b, update)
		}
		if util.Has(api.bot.admins, func(id int64) bool {
			return id == update.Message.From.ID
		}) {
			next(ctx, b, update)
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.From.ID,
				Text:   "Forbidden, contact administrator",
			})
			if api.logger != nil {
				api.logger.Warn().
					Int64("sender", update.Message.From.ID).
					Str("name", fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)).
					Str("user", update.Message.From.Username).
					Msg("Forbidden")
			}
		}
	}
}
