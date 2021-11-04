package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// 設定の構造体
type Config struct {
	// DevMode 開発モードかどうか (default: false)
	Dev_Mode bool `json:"dev_mode,omitempty"`
	// Bot_ID ボットのID (default: "")
	Bot_ID string `json:"bot_id,omitempty"`
	// Verification_Token Botへのリクエストの認証トークン (default: "")
	Verification_Token string `json:"verification_token,omitempty"`
	// Bot_Access_Token Botからのアクセストークン (default: "")
	Bot_Access_Token string `json:"bot_access_token,omitempty"`
}

// 設定を読み込み
func GetConfig() (*Config, error) {
	// デフォルト値の設定
	viper.SetDefault("Dev_Mode", false)
	viper.SetDefault("Bot_ID", "")
	viper.SetDefault("Verification_Token", "")
	viper.SetDefault("Bot_Access_Token", "")

	// 環境変数の読み込み
	viper.AutomaticEnv()

	// .envファイルの読み込み
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Print("Unable to find .env file, default settings or environmental variables are to be used.")
		} else {
			return nil, fmt.Errorf("Error: failed to load .env file - %s ", err)
		}
	}

	var C *Config // 設定格納用変数

	// 設定格納用変数に設定を移す
	err := viper.Unmarshal(&C)
	if err != nil {
		return nil, fmt.Errorf("Error: failed to unmarshal config - %s ", err)
	}

	return C, nil
}
