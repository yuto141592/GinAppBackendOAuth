package handlers

import (
	"hello_gin/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListImages(c *gin.Context) {
	uid := c.GetString("uid")

	iter := services.FirestoreClient.Collection("images").Where("uid", "==", uid).Documents(c)
	defer iter.Stop()

	var results []map[string]string // ここを string にする
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		data := doc.Data()
		results = append(results, map[string]string{
			"id":  data["id"].(string),
			"uid": data["uid"].(string),
			"url": data["url"].(string),
		})
	}

	c.JSON(http.StatusOK, results)
}
