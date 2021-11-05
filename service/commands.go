package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/logica0419/scheduled-messenger-bot/model/event"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
)

var Commands = map[string]string{
	"schedule": "!schedule",
	"join":     "!join",
	"leave":    "!leave",
}

func ParseScheduleCommand(api *api.API, req *event.MessageEvent) (time.Time, string, string, string, error) {
	// メッセージを配列に
	listedReqMes, err := argvParse(req.GetText())
	if err != nil {
		return time.Now(), "", "", "", fmt.Errorf("failed to parse argv: %s", err)
	}

	// 冒頭にメンションがついていた場合要素をドロップ
	if strings.Contains(listedReqMes[0], "@") {
		listedReqMes = listedReqMes[1:]
	}

	// メッセージをパース
	parsedTime, distChannel, body, err := argparseScheduleCommand(listedReqMes)
	if err != nil {
		_ = api.SendMessage(req.GetChannelID(), err.Error())
		return time.Now(), "", "", "", fmt.Errorf("failed to parse argv: %s", err)
	}

	// チャンネルが指定されてないとき、ID と名前を取得
	distChannelID := ""
	if distChannel == "" {
		distChannel = "このチャンネル"
		distChannelID = req.GetChannelID()
	} else {
		// 指定されている場合 distChannel のIDを embedded から取得
		for _, v := range req.Message.Embedded {
			if v.Raw == distChannel && v.Type == "channel" {
				distChannelID = v.ID
				break
			}
		}
	}

	return parsedTime, distChannel, distChannelID, body, err
}
