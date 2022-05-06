package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rendau/limelog/internal/cns"
	"github.com/rendau/limelog/internal/interfaces"
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
	var bytes []byte

	const maxMsgFieldValueSize = 500

	filteredFields := map[string]interface{}{}

	for k, v := range msg {
		if k == "msg" {
			continue
		}

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
		} else {
			switch val := v.(type) {
			case string:
				if len([]rune(val)) > maxMsgFieldValueSize {
					v = string(([]rune(val))[:maxMsgFieldValueSize]) + "..."
				}
			case int64, int32, int16, int8, int, float64, float32:
			default: // try to json-marshal for determine length
				bytes, err = json.Marshal(&v)
				if err != nil {
					log.Println("Fail ot marshal json", err)
					continue
				}
				vStr := string(bytes)
				if len([]rune(vStr)) > maxMsgFieldValueSize {
					vStr = string(([]rune(vStr))[:maxMsgFieldValueSize]) + "..."
					v = vStr
				}
			}
		}

		filteredFields[k] = v
	}

	// filteredFieldsRaw, err := json.MarshalIndent(&filteredFields, "", "   ")
	// if err != nil {
	// 	log.Println("Fail ot marshal json", err)
	// 	return
	// }

	tag, ok := (msg[cns.SfTagFieldName]).(string)
	if !ok {
		tag = ""
	}

	msgContent := " `" + tag + "`: \n\n```\n"

	for k, v := range filteredFields {
		msgContent += fmt.Sprintf("      %s: \t\t%v\n", k, v)
	}

	msgContent += "```"

	tgMsg := tgbotapi.NewMessage(o.chatId, msgContent)
	tgMsg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err = o.botApi.Send(tgMsg)
	if err != nil {
		o.lg.Errorw("Fail to send telegram-msg", err)
		return
	}
}
