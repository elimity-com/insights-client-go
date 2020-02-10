package insights_test

import (
	"encoding/json"
	"github.com/elimity-com/insights-client-go"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	var fun http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/authenticate/foo" {
			t.Fatalf(`got path %q, want "/authenticate/foo"`, request.URL.Path)
		}

		bs, err := ioutil.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("failed reading request body: %v", err)
		}

		type testRequestBody struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}
		var body testRequestBody

		if err := json.Unmarshal(bs, &body); err != nil {
			t.Fatalf("failed unmarshalling body: %v", err)
		}

		if body.Type != "password" {
			t.Fatalf(`got type %q, want "password"`, body.Type)
		}

		if body.Value != "bar" {
			t.Fatalf(`got password %q, want "bar"`, body.Type)
		}

		if _, err := io.WriteString(writer, `{ "token": "baz" }`); err != nil {
			t.Fatalf("failed writing response: %v", err)
		}
	}

	server := httptest.NewServer(fun)
	defer server.Close()

	if _, err := insights.NewClient(server.URL, "foo", "bar"); err != nil {
		t.Errorf("failed creating client: %v", err)
	}
}
