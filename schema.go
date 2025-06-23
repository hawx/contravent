package contravent

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/xeipuuv/gojsonschema"
)

type Schema interface {
	CanMatch(string) bool
	Matches(string) error
}

type SchemaWithExample interface {
	Schema
	sqsEvent(string) (*events.SQSEvent, error)
}

type JSONSchema struct {
	detail       string
	schemaLoader gojsonschema.JSONLoader
	example      string
}

func LoadJSONSchema(detail, schemaPath string) (Schema, error) {
	schemaBody, err := os.ReadFile(schemaPath)
	if err != nil {
		return JSONSchema{}, fmt.Errorf("reading file (%s): %w", schemaPath, err)
	}

	return JSONSchema{
		detail:       detail,
		schemaLoader: gojsonschema.NewBytesLoader(schemaBody),
	}, nil
}

func LoadJSONSchemaWithExample(detail, schemaPath, examplePath string) (JSONSchema, error) {
	schemaBody, err := os.ReadFile(schemaPath)
	if err != nil {
		return JSONSchema{}, fmt.Errorf("reading file (%s): %w", schemaPath, err)
	}

	exampleBody, err := os.ReadFile(examplePath)
	if err != nil {
		return JSONSchema{}, fmt.Errorf("reading file (%s): %w", examplePath, err)
	}

	return JSONSchema{
		detail:       detail,
		schemaLoader: gojsonschema.NewBytesLoader(schemaBody),
		example:      string(exampleBody),
	}, nil
}

func (s JSONSchema) CanMatch(detailType string) bool {
	return s.detail == detailType
}

func (s JSONSchema) Matches(event string) error {
	eventLoader := gojsonschema.NewStringLoader(event)

	result, err := gojsonschema.Validate(s.schemaLoader, eventLoader)
	if err != nil {
		return fmt.Errorf("schema validation: %w", err)
	}

	if !result.Valid() {
		err := MatchError{}

		for _, desc := range result.Errors() {
			err.Reasons = append(err.Reasons, desc.String())
		}

		return err
	}

	return nil
}

func (s JSONSchema) sqsEvent(source string) (*events.SQSEvent, error) {
	if err := s.Matches(s.example); err != nil {
		return nil, err
	}

	body, _ := json.Marshal(events.CloudWatchEvent{
		Source:     source,
		DetailType: s.detail,
		Detail:     json.RawMessage(s.example),
	})

	return &events.SQSEvent{
		Records: []events.SQSMessage{{
			MessageId: "test-msg-id",
			Body:      string(body),
		}},
	}, nil
}
