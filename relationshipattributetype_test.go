package insights_test

import (
	"encoding/json"
	"github.com/elimity-com/insights-client-go"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestClientCreateRelationshipAttributeType(t *testing.T) {
	var fun http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/relationshipAttributeTypes" {
			t.Fatalf(`got path %q, want "/relationshipAttributeTypes"`, request.URL.Path)
		}

		bs, err := ioutil.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("failed reading request body: %v", err)
		}

		type testRequestBody struct {
			ChildType   string `json:"childType"`
			Description string `json:"description"`
			Name        string `json:"name"`
			ParentType  string `json:"parentType"`
			Type        string `json:"type"`
		}

		var body testRequestBody
		if err := json.Unmarshal(bs, &body); err != nil {
			t.Fatalf("failed unmarshalling body: %v", err)
		}

		if body.ChildType != "foo" {
			t.Fatalf(`got category %q, want "foo"`, body.ChildType)
		}

		if body.Description != "bar" {
			t.Fatalf(`got description %q, want "bar"`, body.Description)
		}

		if body.Name != "baz" {
			t.Fatalf(`got name %q, want "baz"`, body.Name)
		}

		if body.ParentType != "asd" {
			t.Fatalf(`got category %q, want "asd"`, body.ParentType)
		}

		if body.Type != "boolean" {
			t.Fatalf(`got type %q, want "boolean"`, body.Type)
		}
	}

	client, server := setup(t, fun)
	defer server.Close()

	typ := insights.NewBooleanType()
	relationshipAttributeType := insights.RelationshipAttributeType{
		FromCategory: "asd",
		Description:  "bar",
		Name:         "baz",
		ToCategory:   "foo",
		Type:         typ,
	}

	if err := client.CreateRelationshipAttributeType(relationshipAttributeType); err != nil {
		t.Fatalf("failed creating attribute type: %v", err)
	}
}
