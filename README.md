# Scheduled-Messenger

指定された時間・指定されたチャンネルに指定されたメッセージを送る traQ Bot

## 環境変数

- `DevMode`  
開発モード (default: false)
- `Bot_ID`  
ボットのID (default: "")
- `Verification_Token`  
Botへのリクエストの認証トークン (default: "")
- `Bot_Access_Token`  
Botからのアクセストークン (default: "")
- `Log_Chan_ID`  
エラーログを送信するチャンネルのID (default: "")
- `MariaDB_Hostname`  
DB のホスト (default: "mariadb")
- `MariaDB_Database`  
DB の DB 名 (default: "SchMes")
- `MariaDB_Username`  
DB のユーザー名 (default: "root")
- `MariaDB_Password`  
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
export LOG_CHAN_ID
export MARIADB_HOSTNAME=
export MARIADB_DATABASE=
export MARIADB_USERNAME=
export MARIADB_PASSWORD=

go run ./*.go
```

MariaDBが`{MONGODB_HOSTNAME}:3306`(デフォルトのポート)で立っていることを確認してください。  
ポート`8080`でサーバーが立つので、`localhost:8080`のエンドポイントにリクエストを送り、レスポンスを確かめてください。
