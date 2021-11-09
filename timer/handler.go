package timer

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/logica0419/scheduled-messenger-bot/model"
	"github.com/logica0419/scheduled-messenger-bot/service"
)

// 予約投稿ハンドラ
func (t *Timer) schMesHandler() {
	// 現在時刻を取得
	currentTime := time.Now()

	// 現在時刻より前に送信予定のメッセージを DB から取得
	mesList, err := t.repo.GetSchMesByTime(currentTime)
	if err != nil {
		// エラーが発生した時は、ログを config で指定したチャンネルに送信
		_ = t.api.SendMessage(t.c.Log_Chan_ID, fmt.Sprintf("ErrorLog: %s レコードの取得に失敗しました\n```\nError: %s\n```", currentTime.Format("01/02 15:04"), err.Error()))
		return
	}

	// メッセージが無く、Dev Mode でない場合 return
	if len(mesList) == 0 && !t.c.Dev_Mode {
		return
	}

	// ログ用の、実際に送信されたメッセージカウント変数
	sentMes := 0

	// 処理終了待機用ウェイトグループ
	wg := new(sync.WaitGroup)

	// 全メッセージの送信と削除処理を並列で行う
	for _, v := range mesList {
		// ウェイトグループに完了待ちを 1 追加
		wg.Add(1)

		go func(mes *model.SchMes) {
			// 関数の処理が終了したらウェイトグループに完了を送る
			defer wg.Done()

			// メッセージを作成
			sendingMes := service.CreateScheduledMessage(mes.UserID, mes.Body)

			// 指定したチャンネルにメッセージを送信
			err = t.api.SendMessage(mes.ChannelID.String(), sendingMes)
			// エラーが起きたらログを config で指定したチャンネルに送信
			if err != nil {
				_ = t.api.SendMessage(t.c.Log_Chan_ID, fmt.Sprintf("ErrorLog: %s メッセージの送信に失敗しました\n```\nID: %s\nError: %s\n```", currentTime.Format("01/02 15:04"), mes.ID, err.Error()))
			} else {
				// 送った ID のメッセージを DB から削除
				err = t.repo.DeleteSchMesByID(mes.ID)
				// エラーが起きたらログを config で指定したチャンネルに送信
				if err != nil {
					_ = t.api.SendMessage(t.c.Log_Chan_ID, fmt.Sprintf("ErrorLog: %s レコードの削除に失敗しました\n```\nID: %s\nError: %s\n```", currentTime.Format("01/02 15:04"), mes.ID, err.Error()))
				} else {
					// メッセージカウントを 1 足す
					sentMes++
				}
			}
		}(v)
	}

	// 全 routine の終了を待機し、ログを表示
	wg.Wait()
	log.Printf("Log: %s %d個の予約投稿メッセージが正常に送信されました", currentTime.Format("01/02 15:04"), sentMes)
}

// 定期投稿ハンドラ
func (t *Timer) schMesPeriodicHandler() {
	// 現在時刻を取得
	currentTime := time.Now()

	// メッセージを DB から全取得
	mesList, err := t.repo.GetSchMesPeriodicAll()
	if err != nil {
		// エラーが発生した時は、ログを config で指定したチャンネルに送信
		_ = t.api.SendMessage(t.c.Log_Chan_ID, fmt.Sprintf("ErrorLog: %s レコードの取得に失敗しました\n```\nError: %s\n```", currentTime.Format("01/02 15:04"), err.Error()))
		return
	}

	// メッセージを時間でフィルター
	mesList = filterSchMesPeriodicByTime(mesList, currentTime)

	// メッセージが無く、Dev Mode でない場合 return
	if len(mesList) == 0 && !t.c.Dev_Mode {
		return
	}

	// ログ用の、実際に送信されたメッセージカウント変数
	sentMes := 0

	// 処理終了待機用ウェイトグループ
	wg := new(sync.WaitGroup)

	// 全メッセージの送信と削除処理を並列で行う
	for _, v := range mesList {
		// ウェイトグループに完了待ちを 1 追加
		wg.Add(1)

		go func(mes *model.SchMesPeriodic) {
			// 関数の処理が終了したらウェイトグループに完了を送る
			defer wg.Done()

			// メッセージを作成
			sendingMes := service.CreateScheduledMessage(mes.UserID, mes.Body)

			// 指定したチャンネルにメッセージを送信
			err = t.api.SendMessage(mes.ChannelID.String(), sendingMes)
			// エラーが起きたらログを config で指定したチャンネルに送信
			if err != nil {
				_ = t.api.SendMessage(t.c.Log_Chan_ID, fmt.Sprintf("ErrorLog: %s メッセージの送信に失敗しました\n```\nID: %s\nError: %s\n```", currentTime.Format("01/02 15:04"), mes.ID, err.Error()))
				return
			}

			// リピートの項目がある場合、回数を一回減らす
			if mes.Repeat != nil {
				repeat := *mes.Repeat - 1
				err = t.repo.UpdateSchMesPeriodic(&model.SchMesPeriodic{ID: mes.ID, Repeat: &repeat})
				if err != nil {
					_ = t.api.SendMessage(t.c.Log_Chan_ID, fmt.Sprintf("ErrorLog: %s レコードの更新に失敗しました\n```\nID: %s\nError: %s\n```", currentTime.Format("01/02 15:04"), mes.ID, err.Error()))
				}

				// 回数が 0 になったら、DB から削除
				if repeat == 0 {
					err = t.repo.DeleteSchMesPeriodicByID(mes.ID)
					if err != nil {
						_ = t.api.SendMessage(t.c.Log_Chan_ID, fmt.Sprintf("ErrorLog: %s レコードの削除に失敗しました\n```\nID: %s\nError: %s\n```", currentTime.Format("01/02 15:04"), mes.ID, err.Error()))
					}
				}
			}

			// メッセージカウントを 1 足す
			sentMes++
		}(v)
	}

	// 全 routine の終了を待機し、ログを表示
	wg.Wait()
	log.Printf("Log: %s %d個の定期投稿メッセージが正常に送信されました", currentTime.Format("01/02 15:04"), sentMes)
}

func filterSchMesPeriodicByTime(mesList []*model.SchMesPeriodic, currentTime time.Time) []*model.SchMesPeriodic {
	// フィルター済みのメッセージリスト格納用変数
	filtered := []*model.SchMesPeriodic{}

	// 個々のメッセージに対して現在時刻が表現にマッチするか調べ、マッチしたらフィルター済みリストに追加
	for _, v := range mesList {
		matches := v.Time.Matches(currentTime)
		if matches {
			filtered = append(filtered, v)
		}
	}

	return filtered
}
