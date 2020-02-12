package insights_test

import (
	"encoding/json"
	"github.com/elimity-com/insights-client-go"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestClientCreateAttributeType(t *testing.T) {
	var fun http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/attributeTypes" {
			t.Fatalf(`got path %q, want "/attributeTypes"`, request.URL.Path)
		}

		bs, err := ioutil.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("failed reading request body: %v", err)
		}

		type testRequestBody struct {
			Category    string `json:"category"`
			Description string `json:"description"`
			Name        string `json:"name"`
			Type        string `json:"type"`
		}

		var body testRequestBody
		if err := json.Unmarshal(bs, &body); err != nil {
			t.Fatalf("failed unmarshalling body: %v", err)
		}

		if body.Category != "foo" {
			t.Fatalf(`got category %q, want "foo"`, body.Category)
		}

		if body.Description != "bar" {
			t.Fatalf(`got description %q, want "bar"`, body.Category)
		}

		if body.Name != "baz" {
			t.Fatalf(`got name %q, want "baz"`, body.Category)
		}

		if body.Type != "boolean" {
			t.Fatalf(`got type %q, want "boolean"`, body.Category)
		}
	}

	client, server := setup(t, fun)
	defer server.Close()

	typ := insights.NewBooleanType()
	attributeType := insights.AttributeType{
		Category:    "foo",
		Description: "bar",
		Name:        "baz",
		Type:        typ,
	}

	if err := client.CreateAttributeType(attributeType); err != nil {
		t.Fatalf("failed creating attribute type: %v", err)
	}
}
