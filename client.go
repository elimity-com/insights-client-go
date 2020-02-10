package insights

import (
	"bytes"
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
	} else {
		requestBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			panic(err)
		}
		return bytes.NewReader(requestBodyBytes)
	}
}

// Client represents an authenticated HTTP client for an Elimity Insights server.
type Client struct {
	basePath string
	token    string
}

// NewClient creates a new client that is authenticated with the given credentials at a server at the given base path.
func NewClient(basePath, userID, password string) (Client, error) {
	type authenticateRequestBody struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}
	requestBody := authenticateRequestBody{
		Type:  "password",
		Value: password,
	}

	type authenticateResponseBody struct {
		Token string `json:"token"`
	}
	var responseBody authenticateResponseBody

	params := performRequestParams{
		basePath:       basePath,
		method:         http.MethodPost,
		pathComponents: []string{"authenticate", userID},
		requestBody:    requestBody,
		responseBody:   &responseBody,
	}
	if err := performRequest(params); err != nil {
		return Client{}, err
	}

	client := Client{
		basePath: basePath,
		token:    responseBody.Token,
	}
	return client, nil
}

func (c Client) performRequest(method string, pathComponents []string, requestBody, responseBody interface{}) error {
	authorization := fmt.Sprintf("Bearer %s", c.token)
	params := performRequestParams{
		authorization:  authorization,
		basePath:       c.basePath,
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

	response, err := http.DefaultClient.Do(request)
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
