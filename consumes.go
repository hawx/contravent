package contravent

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func Consumes(t *testing.T, source string, schema SchemaWithExample, fn func(event *events.SQSEvent) error) {
	event, err := schema.sqsEvent(source)
	if err != nil {
		switch v := err.(type) {
		case MatchError:
			msg := "example did not match schema"
			for _, r := range v.Reasons {
				msg += "\n- " + r
			}

			t.Fail()
			t.Log(msg)

		default:
			t.Fail()
			t.Log(err)
		}
		return
	}

	if err := fn(event); err != nil {
		t.Log(err)
		t.Fail()
	}
}
