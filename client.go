package insights

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
)

func makeRequestBodyReader(requestBody interface{}) io.Reader {
	if requestBody == nil {
		return nil
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(requestBodyBytes)
}

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

func (c Client) performRequest(method string, pathComponents []string, requestBody, responseBody interface{}) error {
	authorization := fmt.Sprintf("Bearer %s", c.token)
	params := performRequestParams{
		authorization:  authorization,
		basePath:       c.basePath,
		client:         c.client,
		method:         method,
		pathComponents: pathComponents,
		requestBody:    requestBody,
		responseBody:   responseBody,
	}
	return performRequest(params)
}

type performRequestParams struct {
	authorization  string
	basePath       string
	client         *http.Client
	method         string
	pathComponents []string
	requestBody    interface{}
	responseBody   interface{}
}

func performRequest(params performRequestParams) error {
	p := path.Join(params.pathComponents...)
	url := fmt.Sprintf("%s/%s", params.basePath, p)
	requestBodyReader := makeRequestBodyReader(params.requestBody)
	request, err := http.NewRequest(params.method, url, requestBodyReader)
	if err != nil {
		panic(err)
	}

	if params.authorization != "" {
		request.Header.Set("Authorization", params.authorization)
	}
	if params.requestBody != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := params.client.Do(request)
	if err != nil {
		return fmt.Errorf("failed sending request: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode > 300 {
		return fmt.Errorf("response has non-success status code %d", response.StatusCode)
	}

	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		if err := response.Body.Close(); err != nil {
			return fmt.Errorf("failed closing response body: %w", err)
		}
		return fmt.Errorf("failed reading from response body: %w", err)
	}
	if err := response.Body.Close(); err != nil {
		return fmt.Errorf("failed closing response body: %w", err)
	}

	if params.responseBody != nil {
		if err := json.Unmarshal(bs, params.responseBody); err != nil {
			panic(err)
		}
	}

	return nil
}
