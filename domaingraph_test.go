package insights_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/elimity-com/insights-client-go"
	"github.com/google/go-cmp/cmp"
)

func TestClientReloadDomainGraph(t *testing.T) {
	expectedBodyString := `{
		"entities": [
			{
				"active": true,
				"attributeAssignments": [
					{
						"attributeTypeName": "foo",
						"value": {
							"type": "boolean",
							"value": "true"
						}
					},
					{
						"attributeTypeName": "bar",
						"value": {
							"type": "date",
							"value": "2006-01-02"
						}
					}
				],
				"id": "foo",
				"name": "bar",
				"type": "baz"
			},
			{
				"active": false,
				"attributeAssignments": [
					{
						"attributeTypeName": "baz",
						"value": {
							"type": "time",
							"value": "15:04:05Z"
						}
					}
				],
				"id": "bar",
				"name": "baz",
				"type": "foo"
			}
		],
		"relationships": [
			{
				"attributeAssignments": [
					{
						"attributeTypeName": "asd",
						"value": {
							"type": "string",
							"value": "asd"
						}
					}
				],
				"fromId": "foo",
				"fromType": "baz",
				"toId": "bar",
				"toType": "foo"
			}
		]
	}`

	client, server := domainGraphTestClientServer(t, expectedBodyString)
	defer server.Close()

	fooValue := insights.NewBooleanValue(true)
	barTime := time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC)
	barValue := insights.NewDateValue(barTime)
	bazTime := time.Date(0, time.January, 1, 15, 4, 5, 0, time.UTC)
	bazValue := insights.NewTimeValue(bazTime)
	asdValue := insights.NewStringValue("asd")

	domainGraph := insights.DomainGraph{
		Entities: []insights.Entity{
			{
				Active: true,
				AttributeAssignments: []insights.AttributeAssignment{
					{
						AttributeTypeName: "foo",
						Value:             fooValue,
					},
					{
						AttributeTypeName: "bar",
						Value:             barValue,
					},
				},
				ID:   "foo",
				Name: "bar",
				Type: "baz",
			},
			{
				Active: false,
				AttributeAssignments: []insights.AttributeAssignment{
					{
						AttributeTypeName: "baz",
						Value:             bazValue,
					},
				},
				ID:   "bar",
				Name: "baz",
				Type: "foo",
			},
		},
		Relationships: []insights.Relationship{
			{
				AttributeAssignments: []insights.AttributeAssignment{
					{
						AttributeTypeName: "asd",
						Value:             asdValue,
					},
				},
				FromEntityID:   "foo",
				FromEntityType: "baz",
				ToEntityID:     "bar",
				ToEntityType:   "foo",
			},
		},
	}

	if err := client.ReloadDomainGraph(domainGraph); err != nil {
		t.Fatalf("failed reloading domain graph: %v", err)
	}
}

func TestClientReloadDomainGraphTimestamp(t *testing.T) {
	expectedBodyString := `{
		"entities": [],
		"relationships": [],
		"historyTimestamp": "2020-06-15T16:26:10Z"
	}`

	client, server := domainGraphTestClientServer(t, expectedBodyString)
	defer server.Close()

	timestamp := time.Date(2020, time.June, 15, 16, 26, 10, 0, time.UTC)
	domainGraph := insights.DomainGraph{
		Timestamp: &timestamp,
	}

	if err := client.ReloadDomainGraph(domainGraph); err != nil {
		t.Fatalf("failed reloading domain graph: %v", err)
	}
}

func domainGraphTestClientServer(t *testing.T, expectedBodyString string) (insights.Client, *httptest.Server) {
	var handlerFunc http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/custom-connector-domain-graph-reload" {
			t.Fatalf(`got path %q, want "/custom-connector-domain-graph-reload"`, request.URL.Path)
		}

		actualBodyBytes, err := ioutil.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("failed reading request body: %v", err)
		}

		var actualBody interface{}
		if err := json.Unmarshal(actualBodyBytes, &actualBody); err != nil {
			t.Fatalf("failed unmarshalling body: %v", err)
		}

		expectedBodyBytes := []byte(expectedBodyString)

		var expectedBody interface{}
		if err := json.Unmarshal(expectedBodyBytes, &expectedBody); err != nil {
			panic(err)
		}

		if diff := cmp.Diff(expectedBody, actualBody); diff != "" {
			t.Fatalf("body mismatch (-got, +want):\n%s", diff)
		}
	}

	return setup(t, handlerFunc)
}
