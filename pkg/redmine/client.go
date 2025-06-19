package redmine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	Debug      bool
}

func NewClient(baseURL, apiKey string) *Client {
	// URLの末尾のスラッシュを除去
	baseURL = strings.TrimRight(baseURL, "/")
	
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doRequest(method, path string, params url.Values, body interface{}) ([]byte, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	// パスを結合
	u.Path = strings.TrimRight(u.Path, "/") + path
	if params != nil {
		u.RawQuery = params.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Redmine-API-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if c.Debug {
		fmt.Printf("[DEBUG] %s %s\n", method, u.String())
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// エラーハンドリング
	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("authentication failed: invalid API key or unauthorized access\nURL: %s", u.String())
	}
	
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("not found: the requested resource does not exist\nURL: %s", u.String())
	}
	
	// HTMLが返ってきた場合（JSONではない）
	if strings.HasPrefix(strings.TrimSpace(string(respBody)), "<") {
		return nil, fmt.Errorf("invalid response: expected JSON but got HTML. Please check your REDMINE_URL is correct and includes the protocol (http:// or https://)\nURL: %s", u.String())
	}
	
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s\nURL: %s", resp.StatusCode, string(respBody), u.String())
	}

	return respBody, nil
}

func (c *Client) Get(path string, params url.Values, result interface{}) error {
	body, err := c.doRequest("GET", path, params, nil)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func (c *Client) Post(path string, params url.Values, reqBody interface{}, result interface{}) error {
	body, err := c.doRequest("POST", path, params, reqBody)
	if err != nil {
		return err
	}

	if result != nil && len(body) > 0 {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) Put(path string, params url.Values, reqBody interface{}) error {
	_, err := c.doRequest("PUT", path, params, reqBody)
	return err
}