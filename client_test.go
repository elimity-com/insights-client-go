package insights_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elimity-com/insights-client-go/v3"
)

func TestNewClientDisableTLSCertificateVerification(t *testing.T) {
	fun := http.HandlerFunc(handler)
	server := httptest.NewTLSServer(fun)
	client := insights.NewClientDisableTLSCertificateVerification(server.URL, "foo")
	if err := client.Infof("foo"); err != nil {
		t.Fatalf("failed creating info log: %v", err)
	}
}

func handler(http.ResponseWriter, *http.Request) {}

func setup(t *testing.T, handler http.Handler) (insights.Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	client, err := insights.NewClient(server.URL, "foo")
	if err != nil {
		t.Fatalf("failed creating client: %v", err)
	}
	return client, server
}
