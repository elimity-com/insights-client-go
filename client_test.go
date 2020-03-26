package insights_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elimity-com/insights-client-go"
)

func setup(t *testing.T, handler http.Handler) (insights.Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	client, err := insights.NewClient(server.URL, "foo")
	if err != nil {
		t.Fatalf("failed creating client: %v", err)
	}
	return client, server
}
