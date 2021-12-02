package main

import (
	"net/http"
	"os"

	"frank.com/gcp/publisher"
	"github.com/gin-gonic/gin"
)

type Message struct {
	Id       string `json:"id" binding:"required"`
	Key      string `json:"key" binding:"required"`
	User     string `json:"user" binding:"required"`
	Abstract string `json:"abstract" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Filename string `json:"filename" binding:"required"`
}

/*
  export GOOGLE_APPLICATION_CREDENTIALS=/Users/ETORRIFUT/work/SA.json

  Da eseguire con ./web training-gcp-309207 topic-frank

  Per pubblicare

  curl -X POST  http://localhost:8080/content \
    -d '{"id":"23846582137", "key":"key1", "user":"frank", "abstract":"an abstract", "content":"some interesting content", "filename":"content-x.json"}'

  curl -X POST  http://localhost:8080/content/bulk \
    -d '[
		{"id":"0000001", "key":"key1", "user":"frank", "abstract":"an abstract 1", "content":"some interesting content part 1", "filename":"content-frank-1.json"},
		{"id":"0000002", "key":"key1", "user":"frank", "abstract":"an abstract 2", "content":"some interesting content part 2", "filename":"content-frank-2.json"},
		{"id":"0000003", "key":"key2", "user":"mike", "abstract":"shiny launder", "content":"foo bar boo", "filename":"content-mike-2.json"},
		{"id":"0000004", "key":"key3", "user":"john", "abstract":"mistic pizza", "content":"foo bar boo", "filename":"content-john-2.json"},
		{"id":"0000005", "key":"key1", "user":"frank", "abstract":"an abstract 3", "content":"some interesting content part 3", "filename":"content-frank-3.json"}
		]'

*/
func main() {

	router := gin.Default()
	projectId := os.Args[1] // training-gcp-309207
	topic := os.Args[2]     // topic-frank

	// Just a sample
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/content", func(c *gin.Context) {
		var json Message
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// fmt.Println("Received message with ID " + json.Id)

		publishedMessage := publisher.PublishedMessage{Id: json.Id, Key: json.Key, Content: json.Content, Filename: json.Filename}
		publisher.PublishSimple(os.Stdout, projectId, topic, publishedMessage)

		c.JSON(http.StatusOK, gin.H{"status": "message published"})
	})

	router.POST("/content/bulk", func(c *gin.Context) {
		var messages []Message
		if err := c.ShouldBindJSON(&messages); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// fmt.Println("Received message with ID " + json.Id)

		// Declare a slice
		var publishedMessages []publisher.PublishedMessage
		for _, m := range messages {
			publishedMessage := publisher.PublishedMessage{Id: m.Id, Key: m.Key, Content: m.Content, Filename: m.Filename}
			publishedMessages = append(publishedMessages, publishedMessage)
		}

		publisher.PublishOrdered(os.Stdout, projectId, topic, publishedMessages)

		c.JSON(http.StatusOK, gin.H{"status": "published bulk"})
	})

	router.Run(":8080")
}
