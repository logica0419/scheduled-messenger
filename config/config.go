package config

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// 設定の構造体
type Config struct {
	Dev_Mode            bool   `json:"dev_mode,omitempty"`            // 開発モード (default: false)
	Bot_ID              string `json:"bot_id,omitempty"`              // ボットのID (default: "")
	Verification_Token  string `json:"verification_token,omitempty"`  // Bot へのリクエストの認証トークン (default: "")
	Bot_Access_Token    string `json:"bot_access_token,omitempty"`    // Bot からのアクセストークン (default: "")
	Log_Chan_ID         string `json:"log_chan_id,omitempty"`         // エラーログを送信するチャンネルの ID (default: "")
	NS_MariaDB_Hostname string `json:"ns_mariadb_hostname,omitempty"` // DB のホスト (default: "mariadb")
	NS_MariaDB_Database string `json:"ns_mariadb_database,omitempty"` // DB の DB 名 (default: "SchMes")
	NS_MariaDB_User     string `json:"ns_mariadb_user,omitempty"`     // DB のユーザー名 (default: "root")
	NS_MariaDB_Password string `json:"ns_mariadb_password,omitempty"` // DB のパスワード (default: "password")
}

// 設定を読み込み
func GetConfig() (*Config, error) {
	// デフォルト値の設定
	viper.SetDefault("Dev_Mode", false)
	viper.SetDefault("NS_MariaDB_Hostname", "localhost")
	viper.SetDefault("NS_MariaDB_Database", "SchMes")
	viper.SetDefault("NS_MariaDB_User", "root")
	viper.SetDefault("NS_MariaDB_Password", "password")

	// 環境変数の読み込み
	viper.AutomaticEnv()

	// config.json ファイルの読み込み
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			log.Print("Unable to find config.json, default settings or environmental variables are to be used.")
		} else {
			return nil, fmt.Errorf("Error: failed to load config.json - %s ", err)
		}
	}

	// 設定格納用変数
	var c Config

	// 設定格納用変数に設定を移す
	err := viper.Unmarshal(&c)
	if err != nil {
		return nil, fmt.Errorf("Error: failed to parse configs - %s ", err)
	}

	log.Printf("%+v", c)
	return &c, nil
}
