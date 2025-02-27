package groq

import (
	"bytes"
	"encoding/json"
	"net/http"

	"go-api/internal/types"
)

const GroqEndpoint = "https://api.groq.com/openai/v1/chat/completions"

type Client struct {
	apiKey string
	client *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (c *Client) ChatCompletion(req *types.ChatRequest) (*http.Response, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", GroqEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	return c.client.Do(httpReq)
}
