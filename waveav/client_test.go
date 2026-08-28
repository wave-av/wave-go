package waveav

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompleteParsesTypedError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{"code": "AUTH_MISSING", "message": "no Authorization header"},
		})
	}))
	defer srv.Close()
	c := NewClient("k")
	c.BaseURL = srv.URL
	_, err := c.Complete(context.Background(), CompletionRequest{Model: "m", Messages: []Message{{Role: "user", Content: "hi"}}})
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("want *APIError, got %T: %v", err, err)
	}
	if apiErr.Code != "AUTH_MISSING" {
		t.Fatalf("code = %q", apiErr.Code)
	}
}

func TestUsageDecodes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer k" {
			t.Errorf("missing bearer")
		}
		_, _ = w.Write([]byte(`{"source":"live","totalCalls":1000,"totalSpentUsd":0.022}`))
	}))
	defer srv.Close()
	c := NewClient("k")
	c.BaseURL = srv.URL
	u, err := c.Usage(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if u.TotalCalls != 1000 {
		t.Fatalf("totalCalls = %d", u.TotalCalls)
	}
}
