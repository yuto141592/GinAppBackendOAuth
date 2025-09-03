package handlers

import (
	"context"
	"net/http"

	"hello_gin/services"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

// Item構造体
type Item struct {
	Text string `json:"text"`
}

// GET /items
func GetItems(c *gin.Context) {
	ctx := context.Background()
	uid := c.GetString("uid")

	iter := services.FirestoreClient.Collection("users").Doc(uid).Collection("items").Documents(ctx)
	var items []map[string]interface{}

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		data := doc.Data()
		data["id"] = doc.Ref.ID
		items = append(items, data)
	}
	c.JSON(http.StatusOK, items)
}

// POST /items
func CreateItem(c *gin.Context) {
	ctx := context.Background()
	uid := c.GetString("uid")

	var req Item
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, _, err := services.FirestoreClient.Collection("users").Doc(uid).Collection("items").Add(ctx, map[string]interface{}{
		"text": req.Text,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// PUT /items/:id
func UpdateItem(c *gin.Context) {
	ctx := context.Background()
	uid := c.GetString("uid")
	id := c.Param("id")

	var req Item
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := services.FirestoreClient.Collection("users").Doc(uid).Collection("items").Doc(id).Set(ctx, map[string]interface{}{
		"text": req.Text,
	}, firestore.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// DELETE /items/:id
func DeleteItem(c *gin.Context) {
	ctx := context.Background()
	uid := c.GetString("uid")
	id := c.Param("id")

	_, err := services.FirestoreClient.Collection("users").Doc(uid).Collection("items").Doc(id).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
