package handlers

import (
	"context"
	"fmt"
	"hello_gin/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func UploadImage(c *gin.Context) {
	uid := c.GetString("uid")

	// Firestoreからそのユーザーのrefresh_tokenを取得
	doc, err := services.FirestoreClient.
		Collection("users").
		Doc(uid).
		Collection("tokens").
		Doc("google").
		Get(c)
	if err != nil || !doc.Exists() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no oauth token"})
		return
	}

	refreshToken, _ := doc.Data()["refresh_token"].(string)
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}

	// oauth2.Config（Google OAuth設定）を取得
	conf := GetGoogleOauthConfig()

	// refresh_tokenからトークンを生成
	token := &oauth2.Token{RefreshToken: refreshToken}
	ts := conf.TokenSource(context.Background(), token)

	// ユーザーのDriveクライアントを作成
	srv, err := drive.NewService(context.Background(), option.WithTokenSource(ts))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "drive init failed"})
		return
	}

	// ファイル取得
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	f, _ := file.Open()
	defer f.Close()

	// ユーザーのDriveにアップロード（ルート直下）
	driveFile, err := srv.Files.Create(&drive.File{
		Name: file.Filename,
	}).Media(f).Do()
	if err != nil {
		fmt.Println("Drive upload error:", err) // まずログ出力
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 公開リンク作成
	_, _ = srv.Permissions.Create(driveFile.Id, &drive.Permission{
		Role: "reader", Type: "anyone",
	}).Do()

	fileInfo := map[string]interface{}{
		"id":  driveFile.Id,
		"url": "https://drive.google.com/uc?id=" + driveFile.Id,
		"uid": uid,
	}

	// Firestoreに保存（imagesコレクション）
	_, err = services.FirestoreClient.Collection("images").Doc(driveFile.Id).Set(c, fileInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save to Firestore"})
		return
	}

	c.JSON(http.StatusOK, fileInfo)
}
