package sonarr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func (s *Sonarr) get(endpoint string, params url.Values) (*http.Response, error) {
	relativeURL, err := url.Parse(endpoint)
	if err != nil {
		return &http.Response{}, err
	}
	endpointURL := s.baseURL.ResolveReference(relativeURL)
	if params == nil {
		params = endpointURL.Query()
	}
	params.Set("apikey", s.apiKey)
	endpointURL.RawQuery = params.Encode()

	fmt.Printf("sonarr GET request: %s\n", endpoint)

	req, err := http.NewRequest("GET", endpointURL.String(), nil)
	if err != nil {
		return &http.Response{}, err
	}

	return s.HTTPClient.Do(req)
}

func (s *Sonarr) put(endpoint string, payload interface{}) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return &http.Response{}, err
	}
	relativeURL, err := url.Parse(endpoint)
	if err != nil {
		return &http.Response{}, err
	}
	endpointURL := s.baseURL.ResolveReference(relativeURL)

	params := endpointURL.Query()
	params.Set("apikey", s.apiKey)
	endpointURL.RawQuery = params.Encode()

	req, err := http.NewRequest("PUT", endpointURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return &http.Response{}, err
	}

	return s.HTTPClient.Do(req)
}

func (s Sonarr) post(query string, body []byte) (*http.Response, error) {
	relativeURL, err := url.Parse(query)

	if err != nil {
		return &http.Response{}, err
	}

	endpointURL := s.baseURL.ResolveReference(relativeURL)

	client := http.Client{
		Timeout: time.Duration(s.Timeout) * time.Second,
	}

	req, err := http.NewRequest("POST", endpointURL.String(), bytes.NewBuffer(body))

	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func (s *Sonarr) del(endpoint string, params url.Values) (*http.Response, error) {
	relativeURL, err := url.Parse(endpoint)
	if err != nil {
		return &http.Response{}, err
	}
	endpointURL := s.baseURL.ResolveReference(relativeURL)
	if params == nil {
		params = endpointURL.Query()
	}
	params.Set("apikey", s.apiKey)
	endpointURL.RawQuery = params.Encode()

	req, err := http.NewRequest("DELETE", endpointURL.String(), nil)
	if err != nil {
		return &http.Response{}, err
	}

	return s.HTTPClient.Do(req)
}
