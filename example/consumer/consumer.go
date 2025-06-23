package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type Ping struct {
	ID string `json:"id"`
}

func Handle(ctx context.Context, event *events.SQSEvent) (any, error) {
	batchItemFailures := []map[string]any{}
	for _, record := range event.Records {
		if errmsg, err := handleEvent(record); err != nil {
			log.Printf("%s (msgid=%s): %v", errmsg, record.MessageId, err)
			batchItemFailures = append(batchItemFailures, map[string]any{"itemIdentifier": record.MessageId})
		}
	}

	return map[string]any{
		"batchItemFailures": batchItemFailures,
	}, nil
}

func handleEvent(msg events.SQSMessage) (string, error) {
	var event *events.CloudWatchEvent
	if err := json.Unmarshal([]byte(msg.Body), &event); err != nil {
		return "unmarshal event", err
	}

	if event.Source == "my-service" && event.DetailType == "ping" {
		var v Ping
		if err := json.Unmarshal(event.Detail, &v); err != nil {
			return "unmarshal event detail", err
		}

		if v.ID == "" {
			return "handle event", fmt.Errorf("invalid id")
		}
	}

	return "", nil
}
