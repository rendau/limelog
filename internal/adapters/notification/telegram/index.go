package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mechta-market/limelog/internal/cns"
	"github.com/mechta-market/limelog/internal/interfaces"
)

type St struct {
	lg       interfaces.Logger
	botToken string
	chatId   int64

	botApi *tgbotapi.BotAPI
}

func New(lg interfaces.Logger, botToken string, chatId int64) (*St, error) {
	var err error

	res := &St{
		lg:       lg,
		botToken: botToken,
		chatId:   chatId,
	}

	res.botApi, err = tgbotapi.NewBotAPI(res.botToken)
	if err != nil {
		lg.Errorw("Fail to create telegram-bot", err)
		return nil, err
	}

	return res, nil
}

func (o *St) Send(msg map[string]interface{}) {
	var err error

	const maxMsgFieldValueSize = 120

	filteredFields := map[string]interface{}{}

	for k, v := range msg {
		if strings.HasPrefix(k, cns.SystemFieldPrefix) {
			if k == cns.SfTsFieldName {
				switch val := v.(type) {
				case int64:
					v = time.UnixMilli(val).In(cns.AppTimeLocation).Format(time.Stamp)
				case float64:
					v = time.UnixMilli(int64(val)).In(cns.AppTimeLocation).Format(time.Stamp)
				default:
					continue
				}
			} else { // ignore system fields
				continue
			}
		}

		vStr, ok := v.(string)
		if ok && len(vStr) > maxMsgFieldValueSize {
			v = vStr[:maxMsgFieldValueSize] + "..."
		}

		filteredFields[k] = v
	}

	filteredFieldsRaw, err := json.MarshalIndent(&filteredFields, "", "   ")
	if err != nil {
		log.Println("Fail ot marshal json", err)
		return
	}

	tag, ok := (msg[cns.SfTagFieldName]).(string)
	if !ok {
		tag = ""
	}

	msgContent := "-------  " + tag + ": \n\n```\n" + string(filteredFieldsRaw) + "\n```"

	fmt.Println(msgContent)

	tgMsg := tgbotapi.NewMessage(o.chatId, msgContent)
	tgMsg.ParseMode = "Markdown"

	_, err = o.botApi.Send(tgMsg)
	if err != nil {
		o.lg.Errorw("Fail to send telegram-msg", err)
		return
	}
}
