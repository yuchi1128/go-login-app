# Go言語 ログインAPIサンプル

Go (Gin) 製ログインAPI。ユーザー名/Email対応、JWT認証、PostgreSQL使用。Docker + Dev Containers環境。

## 起動手順

1.  `git clone https://github.com/yuchi1128/go-login-app`
2.  VS Codeで開き「Reopen in Container」
3.  `go run main.go`

## API (抜粋)

* `POST /auth/register` : ユーザー登録
* `POST /auth/login` : ログイン
* `GET  /api/profile` : プロフィール (要認証)
* `GET  /health` : ヘルスチェック
