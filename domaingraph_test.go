package insights_test

import (
	"compress/zlib"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/elimity-com/insights-client-go/v3"
	"github.com/google/go-cmp/cmp"
)

func TestClientReloadDomainGraph(t *testing.T) {
	expectedBodyString := `{
		"entities": [
			{
				"attributeAssignments": [
					{
						"attributeTypeName": "foo",
						"value": {
							"type": "boolean",
							"value": true
						}
					},
					{
						"attributeTypeName": "bar",
						"value": {
							"type": "date",
							"value": {
								"day": 2,
								"month": 1,
								"year": 2006
							}
						}
					}
				],
				"id": "foo",
				"name": "bar",
				"type": "baz"
			},
			{
				"attributeAssignments": [
					{
						"attributeTypeName": "baz",
						"value": {
							"type": "time",
							"value": {
								"hour": 15,
								"minute": 4,
								"second": 5
							}
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
				"fromEntityId": "foo",
				"fromEntityType": "baz",
				"toEntityId": "bar",
				"toEntityType": "foo"
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
		if request.URL.Path != "/custom-connector-domain-graphs" {
			t.Fatalf(`got path %q, want "/custom-connector-domain-graphs"`, request.URL.Path)
		}

		reader, err := zlib.NewReader(request.Body)
		if err != nil {
			t.Fatalf("failed creating zlib reader: %v", err)
		}

		actualBodyBytes, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Fatalf("failed reading request body: %v", err)
		}

		if err := reader.Close(); err != nil {
			t.Fatalf("failed closing zlib reader: %w", err)
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
