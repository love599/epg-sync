package provider

import (
	"context"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/epg-sync/epgsync/pkg/errors"
	"github.com/epg-sync/epgsync/pkg/logger"
)

type HTTPClient struct {
	client  *http.Client
	baseURL string
	headers map[string]string
	timeout time.Duration
	mu      sync.Mutex
}

func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		baseURL: baseURL,
		headers: make(map[string]string),
		timeout: timeout,
	}
}

func (c *HTTPClient) SetHeader(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.headers[key] = value
}

func (c *HTTPClient) SetHeaders(headers map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	maps.Copy(c.headers, headers)
}

func (c *HTTPClient) Get(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	return c.GetWithHeaders(ctx, path, params, nil)
}

func (c *HTTPClient) GetWithHeaders(ctx context.Context, path string, params map[string]string, headers map[string]string) ([]byte, error) {
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Add(k, v)
		}
		path += "?" + values.Encode()
	}
	return c.doRequest(ctx, http.MethodGet, path, nil, headers)
}

func (c *HTTPClient) Post(ctx context.Context, path string, body io.Reader) ([]byte, error) {
	return c.PostWithHeaders(ctx, path, body, nil)
}

func (c *HTTPClient) PostWithHeaders(ctx context.Context, path string, body io.Reader, headers map[string]string) ([]byte, error) {
	return c.doRequest(ctx, http.MethodPost, path, body, headers)
}

func (c *HTTPClient) doRequest(ctx context.Context, method, path string, body io.Reader, headers map[string]string) ([]byte, error) {
	var requestURL string
	if !strings.HasPrefix(path, "http") {
		requestURL = c.baseURL + path
	} else {
		requestURL = path
	}

	logger.Debug("HTTP Request",
		logger.String("method", method),
		logger.String("url", requestURL),
	)
	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, errors.HTTPRequestFailed(err, requestURL, 0, "failed to create request")
	}

	c.mu.Lock()
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	c.mu.Unlock()

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.HTTPRequestFailed(err, requestURL, 0, "failed to do request")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.HTTPRequestFailed(err, requestURL, resp.StatusCode, "failed to read response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.HTTPRequestFailed(nil, requestURL, resp.StatusCode, string(data))
	}

	return data, nil
}

func (c *HTTPClient) Close() {
	c.client.CloseIdleConnections()
}
