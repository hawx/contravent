package producer

import (
	"testing"

	"hawx.me/code/contravent"
)

func TestRun(t *testing.T) {
	schema, _ := contravent.LoadJSONSchema("ping", "../schema/ping.json")

	contravent.Produces(t, schema, func(url string) error {
		Run(url)
		return nil
	})
}
