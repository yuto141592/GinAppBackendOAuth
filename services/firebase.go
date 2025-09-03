package services

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var (
	AuthClient      *auth.Client
	FirestoreClient *firestore.Client
)

// InitFirebase は Firebase を初期化し、AuthClient と FirestoreClient を準備する
func InitFirebase() *auth.Client {
	var app *firebase.App
	var err error

	if os.Getenv("FIREBASE_USE_LOCAL_KEY") == "true" {
		opt := option.WithCredentialsFile("serviceAccountKey.json")
		app, err = firebase.NewApp(context.Background(), nil, opt)
	} else {
		secretJSON := os.Getenv("FIREBASE_SERVICE_ACCOUNT")
		if secretJSON == "" {
			log.Fatal("FIREBASE_SERVICE_ACCOUNT が設定されていません")
		}
		opt := option.WithCredentialsJSON([]byte(secretJSON))
		app, err = firebase.NewApp(context.Background(), nil, opt)
	}

	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
	}

	// Auth Client
	AuthClient, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v", err)
	}

	// Firestore Client
	FirestoreClient, err = app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v", err)
	}

	return AuthClient
}
