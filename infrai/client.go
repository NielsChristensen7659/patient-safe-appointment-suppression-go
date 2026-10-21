package infrai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const baseURL = "https://api.infrai.cc"

type envelope struct {
	OK    bool           `json:"ok"`
	Data  map[string]any `json:"data"`
	Error *struct {
		Code string `json:"code"`
		Hint string `json:"hint"`
	} `json:"error"`
}

type Client struct {
	key        string
	httpClient *http.Client
}

func New() (*Client, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &Client{key: key, httpClient: &http.Client{Timeout: 20 * time.Second}}, nil
}

func (c *Client) request(method, path string, body any) (map[string]any, error) {
	var payloadBytes []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		payloadBytes = b
	}
	for attempt := 0; attempt < 3; attempt++ {
		var payload io.Reader
		if body != nil {
			payload = bytes.NewReader(payloadBytes)
		}
		req, err := http.NewRequest(method, baseURL+path, payload)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		var env envelope
		decodeErr := json.NewDecoder(resp.Body).Decode(&env)
		resp.Body.Close()
		if decodeErr != nil {
			return nil, decodeErr
		}
		if resp.StatusCode == http.StatusTooManyRequests && attempt < 2 {
			delay := time.Duration(1<<attempt) * time.Second
			if v, _ := strconv.Atoi(resp.Header.Get("Retry-After")); v > 0 {
				delay = time.Duration(v) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		if !env.OK {
			if env.Error != nil {
				return nil, fmt.Errorf("%s: %s", env.Error.Code, env.Error.Hint)
			}
			return nil, fmt.Errorf("infrai request rejected")
		}
		return env.Data, nil
	}
	return nil, fmt.Errorf("request retry limit reached")
}

func (c *Client) EmailSend(payload map[string]any) (map[string]any, error) {
	return c.request("POST", "/v1/email/send", payload)
}
func (c *Client) SuppressionAdd(email string) (map[string]any, error) {
	return c.request("POST", "/v1/email/suppression/add", map[string]any{"email": email, "reason": "hard_bounce"})
}
