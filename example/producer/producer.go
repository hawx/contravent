package producer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

type Ping struct {
	ID string `json:"id"`
}

func Run(baseURL string) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Print("failed to load aws config:", err)
		return
	}
	cfg.BaseEndpoint = aws.String(baseURL)

	svc := eventbridge.NewFromConfig(cfg)

	v, _ := json.Marshal(Ping{ID: "ping-01"})

	if _, err := svc.PutEvents(ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{{
			EventBusName: aws.String("my-bus"),
			Source:       aws.String("my-service"),
			DetailType:   aws.String("ping"),
			Detail:       aws.String(string(v)),
		}},
	}); err != nil {
		log.Print("failed to send event:", err)
		return
	}
}
