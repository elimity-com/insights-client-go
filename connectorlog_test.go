package insights_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/elimity-com/insights-client-go"
)

func TestClientCreateConnectorLogs(t *testing.T) {
	timestamp := time.Now()
	var fun http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/connectorLogs" {
			t.Fatalf(`got path %q, want "/connectorLogs"`, request.URL.Path)
		}

		bs, err := ioutil.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("failed reading request body: %v", err)
		}

		type connectorLog struct {
			Level     string    `json:"level"`
			Message   string    `json:"message"`
			Timestamp time.Time `json:"timestamp"`
		}

		var logs []connectorLog
		if err := json.Unmarshal(bs, &logs); err != nil {
			t.Fatalf("failed unmarshalling body: %v", err)
		}

		if length := len(logs); length != 1 {
			t.Fatalf("got %d logs, want 1", length)
		}

		log := logs[0]
		if log.Level != "info" {
			t.Fatalf(`got level %q, want "info"`, log.Level)
		}

		if log.Message != "foo" {
			t.Fatalf(`got message %q, want "foo"`, log.Message)
		}

		if !log.Timestamp.Equal(timestamp) {
			t.Fatalf("got timestamp %v, want %v", log.Timestamp, timestamp)
		}
	}

	client, server := setup(t, fun)
	defer server.Close()

	logs := []insights.ConnectorLog{
		{
			Level:     insights.Info,
			Message:   "foo",
			Timestamp: timestamp,
		},
	}
	if err := client.CreateConnectorLogs(logs); err != nil {
		t.Fatalf("failed creating connector logs: %v", err)
	}
}
