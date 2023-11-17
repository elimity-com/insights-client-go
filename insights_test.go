package insights_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elimity-com/insights-client-go/v7"
)

func TestInsights(t *testing.T) {
	handler := handler{t: t}
	server := httptest.NewServer(handler)
	if err := insights.CreateSourceConnectorLogs(server.URL, "foo", "bar", nil, false); err != nil {
		t.Errorf("failed creating logs: %v", err)
	}
	server.Close()
}

type handler struct {
	t *testing.T
}

func (h handler) ServeHTTP(_ http.ResponseWriter, request *http.Request) {
	username, password, _ := request.BasicAuth()
	if username != "foo" || password != "bar" {
		h.t.Error("got invalid credentials")
	}
}
