// Package insights provides functions for connector interactions with an Elimity Insights server.
package insights

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

func request(
	contentType, insightsURL, sourceID, sourceToken, urlFormat string, body io.Reader, skipSSLVerification bool,
) error {
	config := &tls.Config{InsecureSkipVerify: skipSSLVerification}
	transport := &http.Transport{TLSClientConfig: config}
	client := http.Client{Transport: transport}
	url := fmt.Sprintf(urlFormat, insightsURL, sourceID)
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", contentType)
	request.SetBasicAuth(sourceID, sourceToken)
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed sending request: %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return fmt.Errorf("failed closing response body: %v", err)
	}
	return nil
}
