package service

import "fmt"

// チャンネルに JOINED した際のメッセージを生成
func CreateJoinedMessage(path string) string {
	return fmt.Sprintf("これから Scheduled Messenher Bot は`%s`に投稿されるメッセージをチェックします!\nいつでも呼んで下さい!", path)
}

// チャンネルから LEFT した際のメッセージを生成
func CreateLeftMessage() string {
	return "寂しいですがお別れです...\nScheduled Messenher Bot のご利用、ありがとうございました!"
}
