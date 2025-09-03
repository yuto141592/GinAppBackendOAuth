package main

import (
	"hello_gin/handlers"
	"hello_gin/middleware"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	// 認証確認用（ログイン確認用）
	r.GET("/protected", middleware.FirebaseAuthMiddleware(authClient), func(c *gin.Context) {
		uid := c.GetString("uid")
		c.JSON(http.StatusOK, gin.H{"message": "Hello " + uid})
	})

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
