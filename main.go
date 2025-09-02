package main

import (
	"context"
	"hello_gin/middleware"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"

	"github.com/gin-contrib/cors"
)

func main() {
	// Firebase 初期化
	var app *firebase.App
	var err error

	// ローカルでは serviceAccountKey.json を使用
	if os.Getenv("FIREBASE_USE_LOCAL_KEY") == "true" {
		opt := option.WithCredentialsFile("serviceAccountKey.json")
		app, err = firebase.NewApp(context.Background(), nil, opt)
	} else {
		// Cloud Run や環境変数から SecretKey を取得
		secretJSON := os.Getenv("FIREBASE_SERVICE_ACCOUNT")
		if secretJSON == "" {
			log.Fatal("FIREBASE_SERVICE_ACCOUNT が設定されていません")
		}
		opt := option.WithCredentialsJSON([]byte(secretJSON))
		app, err = firebase.NewApp(context.Background(), nil, opt)
	}

	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// 認証クライアント作成
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	// Gin 初期化
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://gin-app.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// 認証が必要なルート
	r.GET("/protected", middleware.FirebaseAuthMiddleware(authClient), func(c *gin.Context) {
		uid := c.GetString("uid")
		c.JSON(200, gin.H{"message": "Hello " + uid})
	})

	// ログ出力用ミドルウェア
	r.Use(gin.Logger())

	// テンプレート読み込み
	r.LoadHTMLGlob("templates/views/*")

	// 公開ルート
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Hello, Gin!",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// サーバーを起動
	log.Println("Server Started!")
	r.Run(":8080")
}
