package publisher

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"
	"sync/atomic"

	"encoding/json"

	"cloud.google.com/go/pubsub"
)

type PublishedMessage struct {
	Id       string `json:"id" binding:"required"`
	Key      string `json:"key" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Filename string `json:"filename" binding:"required"`
}

/*
  Pubblicazione semplice, step 1
*/
func PublishSimple(w io.Writer, projectID, topicID string, message PublishedMessage) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close() // come finally

	t := client.Topic(topicID)

	stringJson, err := json.Marshal(message)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(stringJson),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("Get: %v", err)
	}
	fmt.Fprintf(w, "Published a message; msg ID: %v\n", id)
	return nil
}

/*
  Pubblicazione con ordering key, step 2
*/
func PublishOrdered(w io.Writer, projectID, topicID string, messages []PublishedMessage) {
	ctx := context.Background()

	// Sending messages to the same region ensures they are received in order
	// even when multiple publishers are used.
	// client, err := pubsub.NewClient(ctx, projectID,
	//	 option.WithEndpoint("us-east1-pubsub.googleapis.com:443"))

	client, err := pubsub.NewClient(ctx, projectID) // TODO Il client va creato tutte le volte?
	if err != nil {
		fmt.Fprintf(w, "pubsub.NewClient: %v", err)
		return
	}
	defer client.Close()

	var wg sync.WaitGroup
	var totalErrors uint64
	topic := client.Topic(topicID)
	topic.EnableMessageOrdering = true

	for _, m := range messages {
		jsonMessage, parseErr := json.Marshal(m)
		if parseErr != nil {
			fmt.Printf("Failed to parse message: %s\n", err)
			continue
		}

		res := topic.Publish(ctx, &pubsub.Message{
			Data:        []byte(jsonMessage),
			OrderingKey: m.Key,
		})

		wg.Add(1)
		go func(res *pubsub.PublishResult) {
			defer wg.Done()
			// The Get method blocks until a server-generated ID or
			// an error is returned for the published message.
			_, err := res.Get(ctx)
			if err != nil {
				// Error handling code can be added here.
				fmt.Printf("Failed to publish: %s\n", err)
				atomic.AddUint64(&totalErrors, 1)
				return
			}
		}(res)
	}

	wg.Wait()

	if totalErrors > 0 {
		fmt.Fprintf(w, "%d messages did not publish successfully", totalErrors)
		return
	}

	fmt.Fprint(w, "Published messages with ordering keys successfully\n")
}

func PublishThatScales(w io.Writer, projectID, topicID string, n int) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	var wg sync.WaitGroup
	var totalErrors uint64
	t := client.Topic(topicID)

	for i := 0; i < n; i++ {
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte("Message " + strconv.Itoa(i)),
		})

		wg.Add(1)
		go func(i int, res *pubsub.PublishResult) {
			defer wg.Done()
			// The Get method blocks until a server-generated ID or
			// an error is returned for the published message.
			id, err := res.Get(ctx)
			if err != nil {
				// Error handling code can be added here.
				fmt.Fprintf(w, "Failed to publish: %v", err)
				atomic.AddUint64(&totalErrors, 1)
				return
			}
			fmt.Fprintf(w, "Published message %d; msg ID: %v\n", i, id)
		}(i, result)
	}

	wg.Wait()

	if totalErrors > 0 {
		return fmt.Errorf("%d of %d messages did not publish successfully", totalErrors, n)
	}
	return nil
}
