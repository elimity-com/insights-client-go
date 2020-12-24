package insights

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

// Client represents an authenticated HTTP client for an Elimity Insights server.
type Client struct {
	basePath string
	client   *http.Client
	token    string
}

// NewClient creates a new client that is authenticated with the given token at a server at the given base path.
func NewClient(basePath, token string) (Client, error) {
	client := Client{
		basePath: basePath,
		client:   http.DefaultClient,
		token:    token,
	}
	return client, nil
}

// NewClientDisableTLSCertificateVerification creates a new client that is authenticated with the given token at a
// server at the given base path.
//
// The resulting client does not verify the TLS certificate of the configured server.
func NewClientDisableTLSCertificateVerification(basePath, token string) Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: transport}
	return Client{
		basePath: basePath,
		client:   client,
		token:    token,
	}
}

func (c Client) performRequest(path, requestContentType string, requestBody io.Reader) error {
	url := fmt.Sprintf("%s/%s", c.basePath, path)
	request, err := http.NewRequest(http.MethodPost, url, requestBody)
	if err != nil {
		panic(err)
	}
	header := request.Header
	authorization := fmt.Sprintf("Bearer %s", c.token)
	header.Set("Authorization", authorization)
	header.Set("Content-Type", requestContentType)
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
