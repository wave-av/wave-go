// Package waveav is the official WAVE SDK for Go: video infrastructure for
// people and AI agents, one API for live and on-demand video.
//
// The client is a thin, honest wrapper over the WAVE gateway REST API
// (https://api.wave.online): every method maps to one documented route,
// authentication is a bearer API key from https://console.wave.online, and
// errors carry the gateway's typed error codes.
//
// Basic usage:
//
//	client := waveav.NewClient("wav_live_...")
//	models, err := client.ListModels(context.Background())
package waveav

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultBaseURL is the WAVE gateway (the one API fronting every product).
const DefaultBaseURL = "https://api.wave.online"

// Client is a WAVE API client. Safe for concurrent use.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a client with an API key from https://console.wave.online.
func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// APIError is a typed gateway error: the machine-readable code plus message.
type APIError struct {
	Status int
	Code   string
	// Message is the human-readable detail.
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("wave: %s (%d): %s", e.Code, e.Status, e.Message)
}

func (c *Client) do(ctx context.Context, method, path string, body any, out any) error {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("wave: encode body: %w", err)
		}
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, reader)
	if err != nil {
		return fmt.Errorf("wave: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("wave: request: %w", err)
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("wave: read response: %w", err)
	}
	if res.StatusCode >= 400 {
		var e struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		_ = json.Unmarshal(raw, &e)
		return &APIError{Status: res.StatusCode, Code: e.Error.Code, Message: e.Error.Message}
	}
	if out != nil {
		if err := json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("wave: decode response: %w", err)
		}
	}
	return nil
}

// Model is an inference model served by the WAVE funnel (inference.wave.online).
type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by,omitempty"`
}

// ListModels lists the inference models available through the funnel.
func (c *Client) ListModels(ctx context.Context) ([]Model, error) {
	var out struct {
		Data []Model `json:"data"`
	}
	if err := c.do(ctx, http.MethodGet, "/v1/dispatch/models", nil, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// Message is one chat-completions message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionRequest is an OpenAI-compatible chat completion against the funnel.
type CompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Completion is the funnel's chat-completion response (OpenAI-compatible shape).
type Completion struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Choices []struct {
		Index        int     `json:"index"`
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int     `json:"prompt_tokens"`
		CompletionTokens int     `json:"completion_tokens"`
		TotalTokens      int     `json:"total_tokens"`
		Cost             float64 `json:"cost"`
	} `json:"usage"`
}

// Complete runs one chat completion through the measured-routing funnel.
func (c *Client) Complete(ctx context.Context, req CompletionRequest) (*Completion, error) {
	var out Completion
	if err := c.do(ctx, http.MethodPost, "/v1/dispatch/chat/completions", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UsageSummary is the platform usage snapshot (the /usage public surface).
type UsageSummary struct {
	Source        string  `json:"source"`
	TotalCalls    int     `json:"totalCalls"`
	TotalSpentUSD float64 `json:"totalSpentUsd"`
}

// Usage reads the public usage snapshot (the platform's live meter ledger).
func (c *Client) Usage(ctx context.Context) (*UsageSummary, error) {
	var out UsageSummary
	if err := c.do(ctx, http.MethodGet, "/usage", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
