package main

import (
	"fmt"
	"hello_gin/handlers"
	"hello_gin/middleware"
	"net/http"
	"os"

	"firebase.google.com/go/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func SetupRouter(authClient *auth.Client) *gin.Engine {
	r := gin.Default()

	// CORS設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://gin-app.vercel.app", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.POST("/logout", func(c *gin.Context) {
		// サーバー側のセッション Cookie を削除
		c.SetCookie("session", "", -1, "/", "", false, true)
		c.SetCookie("oauth_state", "", -1, "/", "", false, true)

		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	})

	// 認証確認用（ログイン確認用）
	r.GET("/protected2", middleware.FirebaseAuthMiddleware(authClient), func(c *gin.Context) {
		uid := c.GetString("uid")
		c.JSON(http.StatusOK, gin.H{"message": "Hello " + uid})
	})

	r.GET("/protected3", middleware.FirebaseAuthMiddleware(authClient), func(c *gin.Context) {
		uid := c.GetString("uid")

		token, err := handlers.GetUserToken(uid)
		if err != nil || token == nil || !token.Valid() {
			// OAuth が必要
			c.JSON(http.StatusOK, gin.H{"oauthRequired": true})
			return
		}

		// OAuth トークン有効 → 通常処理
		c.JSON(http.StatusOK, gin.H{"message": "Hello " + uid})
	})

	// OAuth 開始
	// r.GET("/start-oauth", middleware.FirebaseAuthMiddleware(authClient), handlers.StartOAuth)

	// 認証確認用（ログイン確認＋Drive OAuth判定）
	r.GET("/protected", middleware.FirebaseAuthMiddleware(authClient), func(c *gin.Context) {
		uid := c.GetString("uid")
		fmt.Println("GOOGLE_CLIENT_ID:", os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_ID"))

		token, err := handlers.GetUserToken(uid)
		if err != nil || token == nil || !token.Valid() {
			// 1. state を生成
			state := handlers.GenerateRandomState()

			// 2. state を Cookie に保存（1時間有効、httpOnly）
			c.SetCookie("oauth_state", state, 3600, "/", "", false, true)
			c.SetCookie("oauth_uid", uid, 3600, "/", "", false, true)

			// 3. 認可URLを生成
			url := handlers.GetGoogleOauthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
			// ここでバックエンドに出力
			fmt.Println("OAuth URL:", url)

			// 4. フロントに返す
			c.JSON(http.StatusOK, gin.H{"oauthUrl": url})
			return
		}

		// トークン有効 → そのままホームへ
		c.JSON(http.StatusOK, gin.H{"message": "Hello " + uid})
	})

	r.POST("/upload", middleware.FirebaseAuthMiddleware(authClient), handlers.UploadImage)
	r.GET("/images", middleware.FirebaseAuthMiddleware(authClient), handlers.ListImages)

	// GET /oauth/callback
	r.GET("/oauth/callback", handlers.OAuthCallback)
	r.GET("/start-oauth", handlers.StartOAuth)

	// 認証必須ルート
	auth := r.Group("/")
	auth.Use(middleware.FirebaseAuthMiddleware(authClient))
	{
		auth.GET("/items", handlers.GetItems)
		auth.POST("/items", handlers.CreateItem)
		auth.PUT("/items/:id", handlers.UpdateItem)
		auth.DELETE("/items/:id", handlers.DeleteItem)
	}

	return r
}
