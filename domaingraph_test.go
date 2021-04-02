package insights_test

import (
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/elimity-com/insights-client-go/v5"
	"github.com/google/go-cmp/cmp"
)

func TestClientReloadDomainGraph(t *testing.T) {
	expectedBodyString := `{
		"entities": [
			{
				"attributeAssignments": [
					{
						"attributeTypeID": "foo",
						"value": {
							"type": "boolean",
							"value": true
						}
					},
					{
						"attributeTypeID": "bar",
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
						"attributeTypeID": "baz",
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
						"attributeTypeID": "asd",
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
						AttributeTypeID: "foo",
						Value:           fooValue,
					},
					{
						AttributeTypeID: "bar",
						Value:           barValue,
					},
				},
				ID:   "foo",
				Name: "bar",
				Type: "baz",
			},
			{
				AttributeAssignments: []insights.AttributeAssignment{
					{
						AttributeTypeID: "baz",
						Value:           bazValue,
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
						AttributeTypeID: "asd",
						Value:           asdValue,
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
		"historyTimestamp": {
			"day": 1,
			"hour": 2,
			"minute": 3,
			"month": 4,
			"second": 5,
			"year": 6
		},
		"relationships": []
	}`

	client, server := domainGraphTestClientServer(t, expectedBodyString)
	defer server.Close()

	timestamp := time.Date(6, time.April, 1, 2, 3, 5, 0, time.UTC)
	domainGraph := insights.DomainGraph{
		Timestamp: &timestamp,
	}

	if err := client.ReloadDomainGraph(domainGraph); err != nil {
		t.Fatalf("failed reloading domain graph: %v", err)
	}
}

func domainGraphTestClientServer(t *testing.T, expectedBodyString string) (insights.Client, *httptest.Server) {
	sourceID := 5
	var handlerFunc http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		path := fmt.Sprintf("/custom-sources/%d/snapshots", sourceID)
		if request.URL.Path != path {
			t.Fatalf(`got path %q, want %s`, request.URL.Path, path)
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
			t.Fatalf("failed closing zlib reader: %v", err)
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
			t.Fatalf("body mismatch (-want, +got):\n%s", diff)
		}
	}

	return setup(t, handlerFunc, sourceID)
}
