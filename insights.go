// Package insights provides functions for connector interactions with an Elimity Insights server.
package insights

import (
	"fmt"
	"io"
	"net/http"
)

func request(contentType, insightsURL, sourceID, sourceToken, urlFormat string, body io.Reader) error {
	url := fmt.Sprintf(urlFormat, insightsURL, sourceID)
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", contentType)
	request.SetBasicAuth(sourceID, sourceToken)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed sending request: %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return fmt.Errorf("failed closing response body: %v", err)
	}
	return nil
}
