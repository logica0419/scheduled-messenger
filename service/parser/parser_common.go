package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cosiner/argv"
	"github.com/logica0419/scheduled-messenger-bot/model"
)

// プレーンテキストのメッセージを配列に分解
func argvParse(message string) ([]string, error) {
	// パース用関数を定義
	var identity = func(s string) (string, error) { return s, nil }

	// パース
	parsed, err := argv.Argv(message, identity, identity)
	if err != nil || len(parsed) == 0 {
		return nil, err
	}

	return parsed[0], nil
}

// body から特定のルールにマッチする文字列を変換
func bodyParse(body *string) *string {
	// メンションよけのパース (@.{id} を @{id} に変換)
	replacedBody := strings.Replace(*body, "@.", "@", -1)

	return &replacedBody
}

// 記入された時間を time.Time に変換
func TimeParse(t *string) (*time.Time, error) {
	// 記入フォーマットを定義
	const format = "2006/01/02/15:04"

	// フォーマットに従ってパース
	parsed, err := time.ParseInLocation(format, *t, time.Local)
	if err != nil {
		return nil, err
	}

	// 指定された時間が現在時刻より後か確認する
	if time.Now().After(parsed) {
		return nil, fmt.Errorf("現在時刻より後の時間を指定してください")
	}

	return &parsed, nil
}

// 記入された時間を定期投稿の time に変換
func TimeParsePeriodic(t *string) (*model.PeriodicTime, error) {
	// スラッシュとコロンで区切る
	timeArr := regexp.MustCompile("[/:]").Split(*t, -1)

	// 配列の長さの確認
	if len(timeArr) != 5 && len(timeArr) != 6 {
		return nil, fmt.Errorf("フォーマットが異なります")
	}

	// 年がワイルドカードになっていることの確認
	if timeArr[0] != "*" {
		return nil, fmt.Errorf("定期投稿の場合、年は * のみしか使えません")
	}

	// 年をドロップ
	timeArr = timeArr[1:]

	// 定期投稿の time を作成
	parsedTime := &model.PeriodicTime{}
	for i := range timeArr {
		// ワイルドカードの時は何もしない (その項目は何も代入されないのでポインターが nil になる)
		if timeArr[i] != "*" {
			// 値を数値に変換
			intTime, err := strconv.Atoi(timeArr[i])
			if err != nil {
				return nil, fmt.Errorf("時間の数値変換ができません\n%s", err)
			}

			// 項目ごとに Validation と代入
			switch i {
			case 0: // 月
				if intTime < 1 || intTime > 12 {
					return nil, fmt.Errorf("有効な月ではありません")
				}
				parsedTime.Month = &intTime

			case 1: // 日付
				// 月ごとに上限の日付が違うので月によって Validation を変更
				if parsedTime.Month == nil { // 月が nil だった (月の指定がない) 場合
					if intTime < 1 || intTime > 31 {
						return nil, fmt.Errorf("有効な日付ではありません")
					}
				} else {
					switch *parsedTime.Month {
					case 1, 3, 5, 7, 8, 10, 12: // 31 日まである月
						if intTime < 1 || intTime > 31 {
							return nil, fmt.Errorf("有効な日付ではありません")
						}
					case 4, 6, 9, 11: // 30 日まである月
						if intTime < 1 || intTime > 30 {
							return nil, fmt.Errorf("有効な日付ではありません")
						}
					case 2: // 29 日まである月
						if intTime < 1 || intTime > 29 {
							return nil, fmt.Errorf("有効な日付ではありません")
						}
					}
				}
				parsedTime.Date = &intTime

			case 2: // 時間
				if intTime < 0 || intTime > 23 {
					return nil, fmt.Errorf("有効な時刻ではありません")
				}
				parsedTime.Hour = &intTime

			case 3: // 分
				if intTime < 0 || intTime > 59 {
					return nil, fmt.Errorf("有効な時刻ではありません")
				}
				parsedTime.Minute = &intTime

			case 4: // 曜日 (Optional)
				if intTime < 0 || intTime > 6 {
					return nil, fmt.Errorf("有効な曜日ではありません")
				}
				parsedTime.Day = &intTime
			}
		}
	}

	return parsedTime, nil
}
