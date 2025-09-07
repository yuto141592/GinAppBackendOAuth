package main

import (
	"fmt"
	"hello_gin/services"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// .env を読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("GOOGLE_CLIENT_ID:", os.Getenv("GOOGLE_CLIENT_ID"))

	// Firebase 初期化
	authClient := services.InitFirebase()
	defer services.FirestoreClient.Close()

	// Router 設定
	r := SetupRouter(authClient)

	log.Println("Server Started!")
	fmt.Println("GOOGLE_CLIENT_ID:", os.Getenv("GOOGLE_CLIENT_ID"))

	r.Run(":8080")
}
