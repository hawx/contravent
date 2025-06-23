package consumer

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"hawx.me/code/contravent"
)

func TestHandle(t *testing.T) {
	schema, _ := contravent.LoadJSONSchemaWithExample("ping", "../schema/ping.json", "../schema/ping-example.json")

	contravent.Consumes(t, "my-service", schema, func(event *events.SQSEvent) error {
		result, err := Handle(context.Background(), event)
		if err != nil {
			return err
		}

		if r, ok := result.(map[string]any); ok && r["batchItemFailures"] != nil {
			if l := r["batchItemFailures"].([]map[string]any); len(l) > 0 {
				return errors.New("had batch item failures")
			}
		}

		return nil
	})
}
