package checkhost

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	DefaultBaseURL = "https://api.check-host.cc"
)

// CheckHostException mapped as a generic Go error.
var (
	ErrRateLimit   = errors.New("rate limit reached, please provide an API key or slow down your requests")
	ErrBadRequest  = errors.New("problem with your input parameters, please check your payload")
	ErrServerError = errors.New("internal server error, please try again later")
)

type CheckHost struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new CheckHost client. The apiKey is optional.
func NewClient(apiKey string) *CheckHost {
	return &CheckHost{
		APIKey: apiKey,
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs the HTTP request and unmarshals the response.
func (c *CheckHost) doRequest(method, path string, body interface{}, response interface{}) error {
	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error marshaling request body: %w", err)
		}
		buf = bytes.NewBuffer(b)
	}

	url := fmt.Sprintf("%s/%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("network error occurred during request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	// Handle specific HTTP Status Codes based on check-host API spec
	switch resp.StatusCode {
	case 200:
		// Success
	case 400:
		return fmt.Errorf("%w: %s", ErrBadRequest, string(respBody))
	case 429:
		return fmt.Errorf("%w: %s", ErrRateLimit, string(respBody))
	case 500:
		return fmt.Errorf("%w: %s", ErrServerError, string(respBody))
	default:
		if resp.StatusCode >= 400 {
			return fmt.Errorf("check-host api returned an unexpected status code %d: %s", resp.StatusCode, string(respBody))
		}
	}

	// Unmarshal JSON to the provided response pointer, if any
	if response != nil && len(respBody) > 0 {
		err = json.Unmarshal(respBody, response)
		if err != nil {
			// Some endpoints like /myip simply return a string text/plain
			if strResp, ok := response.(*string); ok {
				*strResp = string(respBody)
				return nil
			}
			return fmt.Errorf("error parsing response as json: %w. response text: %s", err, string(respBody))
		}
	}

	return nil
}
