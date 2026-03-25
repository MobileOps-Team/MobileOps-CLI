package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/MobileOps-Team/mobileops-cli/internal/config"
)

const (
	version       = "0.1.0"
	maxRetries    = 2
	retryInterval = 500 * time.Millisecond
	backoffFactor = 2
)

// Error types matching the Ruby CLI.

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "Authentication failed. Run `mobileops auth login` to configure credentials."
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "Resource not found."
}

type PermissionError struct {
	Message string
}

func (e *PermissionError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "API key does not have permission for this operation."
}

type ApiError struct {
	Message string
	Status  int
}

func (e *ApiError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "API request failed."
}

// Client is the MobileOps API client.
type Client struct {
	host       string
	accessKey  string
	secretKey  string
	httpClient *http.Client
}

// New creates a new Client. It can accept explicit credentials or load from config.
func New(host, accessKey, secretKey, environment string) (*Client, error) {
	env := environment
	if env == "" {
		env = config.CurrentEnvironment()
	}

	h := host
	ak := accessKey
	sk := secretKey

	if ak == "" || sk == "" {
		creds, err := config.Load(env)
		if err == nil && creds != nil {
			if h == "" {
				h = creds.Host
			}
			if ak == "" {
				ak = creds.AccessKey
			}
			if sk == "" {
				sk = creds.SecretKey
			}
		}
	}

	if h == "" {
		h = config.HostFor(env)
	}

	if ak == "" || sk == "" {
		return nil, &AuthError{}
	}

	return &Client{
		host:       h,
		accessKey:  ak,
		secretKey:  sk,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// NewWithKeys creates a client with explicit keys (used during login validation).
func NewWithKeys(host, accessKey, secretKey string) *Client {
	return &Client{
		host:       host,
		accessKey:  accessKey,
		secretKey:  secretKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Get performs a GET request to the API.
func (c *Client) Get(path string, params map[string]string) (map[string]interface{}, error) {
	apiPath := c.apiPath(path)
	u, err := url.Parse(c.host + apiPath)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)

	return c.doWithRetry(req)
}

func (c *Client) apiPath(path string) string {
	path = strings.TrimPrefix(path, "/")
	return "/api/" + path
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("X-Api-Access-Key", c.accessKey)
	req.Header.Set("X-Api-Secret-Key", c.secretKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "mobileops-cli/"+version)
}

func (c *Client) doWithRetry(req *http.Request) (map[string]interface{}, error) {
	var lastErr error
	interval := retryInterval

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(interval)
			interval *= backoffFactor
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		return c.handleResponse(resp)
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

func (c *Client) handleResponse(resp *http.Response) (map[string]interface{}, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &result); err != nil {
			// If body isn't JSON, wrap in a map
			result = map[string]interface{}{"raw": string(body)}
		}
	}

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		// Check for API-level errors returned as 200
		if errField, ok := result["error"]; ok {
			errStr, _ := errField.(string)
			errMsg := extractErrorMessage(result)
			switch errStr {
			case "api_connection_error":
				return nil, &AuthError{Message: errMsg}
			case "api_key_write_error", "api_key_delete_error":
				return nil, &PermissionError{Message: errMsg}
			case "resource_not_found":
				return nil, &NotFoundError{Message: errMsg}
			default:
				return nil, &ApiError{Message: errMsg, Status: resp.StatusCode}
			}
		}
		return result, nil

	case resp.StatusCode == 401:
		return nil, &AuthError{}

	case resp.StatusCode == 403:
		return nil, &PermissionError{}

	case resp.StatusCode == 404:
		msg := extractErrorMessage(result)
		if msg == "" {
			msg = "Resource not found."
		}
		return nil, &NotFoundError{Message: msg}

	default:
		msg := extractErrorMessage(result)
		if msg == "" {
			msg = fmt.Sprintf("API request failed (HTTP %d).", resp.StatusCode)
		}
		return nil, &ApiError{Message: msg, Status: resp.StatusCode}
	}
}

func extractErrorMessage(body map[string]interface{}) string {
	if body == nil {
		return ""
	}
	if msg, ok := body["error_message"].(string); ok {
		return msg
	}
	if msg, ok := body["error"].(string); ok {
		return msg
	}
	if errs, ok := body["errors"]; ok {
		return fmt.Sprintf("%v", errs)
	}
	return ""
}
