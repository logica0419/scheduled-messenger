# Scheduled-Messenger

指定された時間・指定されたチャンネルに指定されたメッセージを送る traQ Bot

## 環境変数

- `DEV_MODE`  
開発モード (default: false)
- `BOT_ID`  
ボットのID (default: "")
- `VERIFICATION_TOKEN`  
Botへのリクエストの認証トークン (default: "")
- `BOT_ACCESS_TOKEN`  
Botからのアクセストークン (default: "")
- `LOG_CHAN_ID`  
エラーログを送信するチャンネルのID (default: "")
- `NS_MARIADB_HOSTNAME`  
DB のホスト (default: "mariadb")
- `NS_MARIADB_DATABASE`  
DB の DB 名 (default: "SchMes")
- `NS_MARIADB_USER`  
DB のユーザー名 (default: "root")
- `NS_MARIADB_PASSWORD`  
DB のパスワード (default: "password")

## ローカルで動かすときのサンプル

シェルスクリプトを使いましょう。  
ディレクトリ内に`{任意の名前}.sh`を作り、下のコードをコピーして環境変数を設定した後、`sh {任意の名前}.sh`で実行します。

```sh *.sh
#!/bin/sh

export DEV_MODE=
export BOT_ID=
export VERIFICATION_TOKEN=
export BOT_ACCESS_TOKEN=
export LOG_CHAN_ID=
export NS_MARIADB_HOSTNAME=
export NS_MARIADB_DATABASE=
export NS_MARIADB_USER=
export NS_MARIADB_PASSWORD=

go run ./*.go
```

MariaDBが`{NS_MARIADB_HOSTNAME}:3306`(デフォルトのポート)で立っていることを確認してください。  
ポート`8080`でサーバーが立つので、`localhost:8080`のエンドポイントにリクエストを送り、レスポンスを確かめてください。
