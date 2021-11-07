package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// 設定の構造体
type Config struct {
	// 開発モード (default: false)
	Dev_Mode bool `json:"dev_mode,omitempty"`
	// ボットのID (default: "")
	Bot_ID string `json:"bot_id,omitempty"`
	// Bot へのリクエストの認証トークン (default: "")
	Verification_Token string `json:"verification_token,omitempty"`
	// Bot からのアクセストークン (default: "")
	Bot_Access_Token string `json:"bot_access_token,omitempty"`
	// エラーログを送信するチャンネルの ID (default: "")
	Log_Chan_ID string `json:"log_chan_id,omitempty"`
	// DB のホスト (default: "mariadb")
	MariaDB_Hostname string `json:"mariadb_hostname,omitempty"`
	// DB の DB 名 (default: "SchMes")
	MariaDB_Database string `json:"mariadb_database,omitempty"`
	// DB のユーザー名 (default: "root")
	MariaDB_Username string `json:"mariadb_username,omitempty"`
	// DB のパスワード (default: "password")
	MariaDB_Password string `json:"mariadb_password,omitempty"`
}

// 設定を読み込み
func GetConfig() (*Config, error) {
	// デフォルト値の設定
	viper.SetDefault("Dev_Mode", false)
	viper.SetDefault("Bot_ID", "")
	viper.SetDefault("Verification_Token", "")
	viper.SetDefault("Bot_Access_Token", "")
	viper.SetDefault("Log_Chan_ID", "")
	viper.SetDefault("MariaDB_Hostname", "mariadb")
	viper.SetDefault("MariaDB_Database", "SchMes")
	viper.SetDefault("MariaDB_Username", "root")
	viper.SetDefault("MariaDB_Password", "password")

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

	// 設定格納用変数
	var C *Config

	// 設定格納用変数に設定を移す
	err := viper.Unmarshal(&C)
	if err != nil {
		return nil, fmt.Errorf("Error: failed to unmarshal config - %s ", err)
	}

	return C, nil
}
