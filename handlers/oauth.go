package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hello_gin/services"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// handlers/oauth.go
func GetGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes:       []string{"https://www.googleapis.com/auth/drive.file"},
		Endpoint:     google.Endpoint,
	}
}

func getKmsKeyName() string {
	key := os.Getenv("OAUTH_TOKEN_KEY")
	fmt.Println("kmsKeyName:", key)

	if key == "" {
		panic("OAUTH_TOKEN_KEY is not set")
	}
	return key
}

func GenerateRandomState() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

// // Firestoreにトークン保存（暗号化付き）
// func SaveUserToken(userID string, token *oauth2.Token) error {
// 	ctx := context.Background()

// 	encAccess, err := encryptWithKMS(ctx, getKmsKeyName(), token.AccessToken)
// 	fmt.Println("encAccess:", encAccess)

// 	if err != nil {
// 		return err
// 	}

// 	encRefresh, err := encryptWithKMS(ctx, getKmsKeyName(), token.RefreshToken)
// 	fmt.Println("encRefresh:", encRefresh)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = services.FirestoreClient.
// 		Collection("users").
// 		Doc(userID).
// 		Collection("tokens").
// 		Doc("google").
// 		Set(ctx, map[string]interface{}{
// 			"access_token":  encAccess,
// 			"refresh_token": encRefresh,
// 			"expiry":        token.Expiry,
// 			"token_type":    token.TokenType,
// 			"updated_at":    time.Now(),
// 		})
// 	return err
// }

// // Firestoreからトークン取得（復号付き）
// func GetUserToken(userID string) (*oauth2.Token, error) {
// 	ctx := context.Background()

// 	doc, err := services.FirestoreClient.
// 		Collection("users").
// 		Doc(userID).
// 		Collection("tokens").
// 		Doc("google").
// 		Get(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	data := doc.Data()

// 	decAccess, err := decryptWithKMS(ctx, getKmsKeyName(), data["access_token"].(string))
// 	if err != nil {
// 		return nil, err
// 	}

// 	decRefresh, err := decryptWithKMS(ctx, getKmsKeyName(), data["refresh_token"].(string))
// 	if err != nil {
// 		return nil, err
// 	}

// 	token := &oauth2.Token{
// 		AccessToken:  decAccess,
// 		RefreshToken: decRefresh,
// 		TokenType:    data["token_type"].(string),
// 		Expiry:       data["expiry"].(time.Time),
// 	}
// 	return token, nil
// }

// // トークンをJSONにしてSecret Managerに保存
// func SaveUserToken(userID string, token *oauth2.Token) error {
// 	ctx := context.Background()
// 	client, err := secretmanager.NewClient(ctx)
// 	if err != nil {
// 		log.Printf("[ERROR] failed to create SecretManager client: %v", err)
// 		return err
// 	}
// 	defer client.Close()

// 	tokenData, err := json.Marshal(token)
// 	if err != nil {
// 		log.Printf("[ERROR] failed to marshal token: %v", err)
// 		return err
// 	}

// 	secretID := fmt.Sprintf("user_token_%s", userID)
// 	parent := fmt.Sprintf("projects/%s/secrets/%s", "ginapp-88534", secretID)

// 	// シークレット作成
// 	_, err = client.CreateSecret(ctx, &secretmanagerpb.CreateSecretRequest{
// 		Parent:   fmt.Sprintf("projects/%s", "ginapp-88534"),
// 		SecretId: secretID,
// 		Secret: &secretmanagerpb.Secret{
// 			Replication: &secretmanagerpb.Replication{
// 				Replication: &secretmanagerpb.Replication_Automatic_{
// 					Automatic: &secretmanagerpb.Replication_Automatic{},
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		if status.Code(err) != codes.AlreadyExists {
// 			log.Printf("[ERROR] failed to create secret: %v", err)
// 			return err
// 		}
// 		log.Printf("[INFO] secret %s already exists, continuing", secretID)
// 	} else {
// 		log.Printf("[INFO] secret %s created", secretID)
// 	}

// 	// シークレットに新しいバージョンを追加
// 	_, err = client.AddSecretVersion(ctx, &secretmanagerpb.AddSecretVersionRequest{
// 		Parent: parent,
// 		Payload: &secretmanagerpb.SecretPayload{
// 			Data: tokenData,
// 		},
// 	})
// 	if err != nil {
// 		log.Printf("[ERROR] failed to add secret version: %v", err)
// 		return err
// 	}

// 	log.Printf("[INFO] token saved for user %s", userID)
// 	return nil
// }

// func GetUserToken(userID string) (*oauth2.Token, error) {
// 	ctx := context.Background()
// 	client, err := secretmanager.NewClient(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer client.Close()

// 	// ユーザーごとのシークレット名
// 	secretID := fmt.Sprintf("user_token_%s", userID)
// 	secretVersion := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", "ginapp-88534", secretID)

// 	req := &secretmanagerpb.AccessSecretVersionRequest{
// 		Name: secretVersion,
// 	}

// 	result, err := client.AccessSecretVersion(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var token oauth2.Token
// 	if err := json.Unmarshal(result.Payload.Data, &token); err != nil {
// 		return nil, err
// 	}

// 	return &token, nil
// }

// Firestoreからトークン取得
func GetUserToken(userID string) (*oauth2.Token, error) {
	doc, err := services.FirestoreClient.Collection("users").Doc(userID).Collection("tokens").Doc("google").Get(context.Background())
	if err != nil {
		return nil, err
	}

	data := doc.Data()
	token := &oauth2.Token{
		AccessToken:  data["access_token"].(string),
		RefreshToken: data["refresh_token"].(string),
		TokenType:    data["token_type"].(string),
		Expiry:       data["expiry"].(time.Time),
	}
	return token, nil
}

// Firestoreにトークン保存
func SaveUserToken(userID string, token *oauth2.Token) error {
	_, err := services.FirestoreClient.Collection("users").Doc(userID).Collection("tokens").Doc("google").Set(context.Background(), map[string]interface{}{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"expiry":        token.Expiry,
		"token_type":    token.TokenType,
		"updated_at":    time.Now(),
	})
	return err
}

// /start-oauth ハンドラー
func StartOAuth(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing UID"})
		return
	}

	state := GenerateRandomState()

	if os.Getenv("LOCAL_DEV") == "True" {
		// Cookie に保存（1時間有効、httpOnly）
		c.SetCookie("oauth_state", state, 3600, "/", "localhost", false, true)
		c.SetCookie("oauth_uid", uid, 3600, "/", "localhost", false, true)
	} else {
		// 本番環境
		c.SetCookie("oauth_state", state, 3600, "/", "gin-app.vercel.app", true, true)
		c.SetCookie("oauth_uid", uid, 3600, "/", "gin-app.vercel.app", true, true)
	}
	// Google 認可 URL
	url := GetGoogleOauthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	// バックエンドから直接リダイレクト
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// /oauth/callback ハンドラー
func OAuthCallback(c *gin.Context) {
	// クエリから state と code を取得
	stateQuery := c.Query("state")
	code := c.Query("code")

	stateCookie, _ := c.Cookie("oauth_state")
	uidCookie, _ := c.Cookie("oauth_uid")

	if stateQuery != stateCookie || uidCookie == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid state or missing UID"})
		return
	}

	token, err := GetGoogleOauthConfig().Exchange(context.Background(), code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}

	// Firestore 等に保存
	if err := SaveUserToken(uidCookie, token); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	// 成功したらフロントのホームにリダイレクト
	c.Redirect(http.StatusSeeOther, os.Getenv("FRONTEND_URL"))
}

// // OAuth開始
// func StartOAuth(c *gin.Context) {
// 	state := GenerateRandomState()
// 	c.SetCookie("oauth_state", state, 3600, "/", "", false, true) // secureは本番でtrue
// 	url := GetGoogleOauthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
// 	c.Redirect(http.StatusTemporaryRedirect, url)
// }

// // OAuthコールバック
// func OAuthCallback(c *gin.Context) {
// 	stateQuery := c.Query("state")
// 	stateCookie, _ := c.Cookie("oauth_state")
// 	if stateQuery != stateCookie {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
// 		return
// 	}

// 	code := c.Query("code")
// 	token, err := GetGoogleOauthConfig().Exchange(context.Background(), code)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
// 		return
// 	}

// 	// uid はコールバックに来る前にセッション等で保持しておく必要あり
// 	// uid := c.Query("uid") // フロントからクエリで渡す or セッションで保持
// 	uidCookie, err := c.Cookie("oauth_uid")
// 	if err != nil || uidCookie == "" {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing UID"})
// 		return
// 	}
// 	uid := uidCookie

// 	if err := SaveUserToken(uid, token); err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
// 		return
// 	}

// 	c.Redirect(http.StatusSeeOther, "/")
// }
