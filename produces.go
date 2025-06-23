package contravent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
)

func Produces(t *testing.T, schema Schema, fn func(string) error) {
	var rerr error
	var hadMatch bool

	closer := make(chan struct{})

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Amz-Target") == "AWSEvents.PutEvents" {
			var v eventbridge.PutEventsInput
			if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
				rerr = err
				closer <- struct{}{}
				return
			}

			for _, entry := range v.Entries {
				if schema.CanMatch(*entry.DetailType) {
					hadMatch = true
					rerr = schema.Matches(*entry.Detail)
					closer <- struct{}{}
					return
				}
			}
		}
	}))
	defer s.Close()

	go func() {
		select {
		case <-closer:
			s.Close()
		case <-time.After(time.Second):
			rerr = ErrTimeout
			s.Close()
		}
	}()

	verr := fn(s.URL)

	if rerr != nil {
		switch v := rerr.(type) {
		case MatchError:
			msg := v.Error()
			for _, r := range v.Reasons {
				msg += "\n- " + r
			}

			t.Log(msg)
		default:
			t.Log(rerr)
		}

		t.Fail()
		return
	}

	if verr != nil {
		t.Log(verr)
		t.Fail()
		return
	}

	if !hadMatch {
		t.Log("no matching event")
		t.Fail()
		return
	}
}
