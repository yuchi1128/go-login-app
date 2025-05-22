package main

import (
	"log"
	"os"

	"github.com/yuchi1128/go-login-app/database"
	"github.com/yuchi1128/go-login-app/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load()
		if err != nil {
			log.Println("注意: .envファイルの読み込みに失敗しました (Docker環境では通常問題ありません)")
		}
	}

	// データベース接続
	database.ConnectDB()
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)


	// ルーターのセットアップ
	router := routes.SetupRouter()

	// サーバーの起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("サーバーをポート %s で起動します...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("サーバー起動エラー: %v", err)
	}
}