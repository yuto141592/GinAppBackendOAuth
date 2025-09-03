package main

import (
	"hello_gin/services"
	"log"
)

func main() {
	// Firebase 初期化
	authClient := services.InitFirebase()
	defer services.FirestoreClient.Close()

	// Router 設定
	r := SetupRouter(authClient)

	log.Println("Server Started!")
	r.Run(":8080")
}
