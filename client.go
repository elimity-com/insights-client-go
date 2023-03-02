package insights

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Client represents an authenticated HTTP client for an Elimity Insights server.
type Client struct {
	basePath string
	client   *http.Client
	token    string
	sourceID int
}

// NewClient creates a new client that is authenticated with the given token at a server at the given base path.
func NewClient(basePath, token string, sourceID int) (Client, error) {
	client := Client{
		basePath: basePath,
		client:   http.DefaultClient,
		token:    token,
		sourceID: sourceID,
	}
	return client, nil
}

// NewClientDisableTLSCertificateVerification creates a new client that is authenticated with the given token at a
// server at the given base path.
//
// The resulting client does not verify the TLS certificate of the configured server.
func NewClientDisableTLSCertificateVerification(basePath, token string, sourceID int) Client {
	config := &tls.Config{InsecureSkipVerify: true}
	return NewClientWithTLSConfig(basePath, token, config, sourceID)
}

// NewClientWithTLSConfig creates a new client that uses the given TLS configuration.
func NewClientWithTLSConfig(basePath, token string, config *tls.Config, sourceID int) Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = config
	client := &http.Client{Transport: transport}
	return Client{
		basePath: basePath,
		client:   client,
		token:    token,
		sourceID: sourceID,
	}
}

func (c Client) performRequest(path, requestContentType string, requestBody io.Reader) error {
	url := fmt.Sprintf("%s/%s", c.basePath, path)
	request, err := http.NewRequest(http.MethodPost, url, requestBody)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", requestContentType)
	request.SetBasicAuth(strconv.Itoa(c.sourceID), c.token)

	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("failed performing request: %w", err)
	}
	if err := response.Body.Close(); err != nil {
		return fmt.Errorf("failed closing response body: %w", err)
	}
	if statusCode := response.StatusCode; statusCode < 200 || statusCode > 299 {
		return fmt.Errorf("got non-success status code %d", statusCode)
	}
	return nil
}
