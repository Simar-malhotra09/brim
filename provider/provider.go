package provider

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"resty.dev/v3"
)

const (
	defaultBaseURL = "http://localhost:8080/search"
	httpTimeout    = 10 * time.Second
	retryCount     = 3
	retryWaitTime  = 200 * time.Millisecond
	retryMaxWait   = 2 * time.Second
)

type SearxResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Content string `json:"content"`
	Engine  string `json:"engine"`
}

type SearxResponse struct {
	Query   string        `json:"query"`
	Results []SearxResult `json:"results"`
}

type Service struct {
	client  *resty.Client
	baseURL string
}

func NewService() *Service {
	client := resty.New()
	client.SetTimeout(httpTimeout)
	client.SetRetryCount(retryCount)
	client.SetRetryWaitTime(retryWaitTime)
	client.SetRetryMaxWaitTime(retryMaxWait)
	client.AddRetryConditions(func(resp *resty.Response, _ error) bool {
		return resp != nil && resp.StatusCode() >= http.StatusInternalServerError
	})
	return &Service{client: client, baseURL: defaultBaseURL}
}

func (s *Service) FetchWebResults(query string, itemsToFetch int) (*SearxResponse, error) {
	fullURL := fmt.Sprintf("%s?q=%s&format=json", s.baseURL, url.QueryEscape(query))

	var result SearxResponse
	resp, err := s.client.R().
		SetResult(&result).
		Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode())
	}

	if itemsToFetch > 0 && len(result.Results) > itemsToFetch {
		result.Results = result.Results[:itemsToFetch]
	}
	return &result, nil
}

