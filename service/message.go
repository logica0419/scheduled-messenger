package service

import "fmt"

// チャンネルに JOIN した際のメッセージを生成
func CreateJoinMessage(path string) string {
	return fmt.Sprintf("Scheduled Messenher Bot will now check all the messages in `%s`!\nPlease call me whenever you want me!", path)
}
