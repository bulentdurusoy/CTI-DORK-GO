package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ─── Serper.dev API Client ───────────────────────────────────────

const serperEndpoint = "https://google.serper.dev/search"

// SerperSearchEngine executes real Google searches via the Serper.dev API.
// Returns only actual search results — no fabricated or synthetic data.
type SerperSearchEngine struct {
	apiKey     string
	httpClient *http.Client
}

// NewSerperSearchEngine creates a new Serper.dev search engine with the given API key
func NewSerperSearchEngine(apiKey string) *SerperSearchEngine {
	return &SerperSearchEngine{
		apiKey: strings.TrimSpace(apiKey),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// serperRequest is the JSON payload sent to Serper.dev
type serperRequest struct {
	Q   string `json:"q"`
	Num int    `json:"num,omitempty"` // number of results (default 10)
}

// serperResponse is the JSON response from Serper.dev
type serperResponse struct {
	Organic []serperOrganicResult `json:"organic"`
}

// serperOrganicResult represents a single organic search result from Serper.dev
type serperOrganicResult struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}

// Search executes a real Google search via Serper.dev API.
// Returns only actual results from Google. If the API returns no results
// or an error occurs, an empty slice is returned — never fabricated data.
func (s *SerperSearchEngine) Search(query string) ([]SearchResult, error) {
	if s.apiKey == "" {
		return []SearchResult{}, fmt.Errorf("no Serper API key configured")
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return []SearchResult{}, nil
	}

	// Build request payload
	payload := serperRequest{
		Q:   query,
		Num: 10,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return []SearchResult{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", serperEndpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return []SearchResult{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-API-KEY", s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return []SearchResult{}, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []SearchResult{}, fmt.Errorf("failed to read response: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return []SearchResult{}, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var serperResp serperResponse
	if err := json.Unmarshal(body, &serperResp); err != nil {
		return []SearchResult{}, fmt.Errorf("failed to parse API response: %w", err)
	}

	// Map Serper results → SearchResult
	var results []SearchResult
	for _, organic := range serperResp.Organic {
		link := strings.TrimSpace(organic.Link)
		title := strings.TrimSpace(organic.Title)

		// Strict validation: skip invalid results
		if link == "" || title == "" {
			continue
		}
		if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
			continue
		}

		results = append(results, SearchResult{
			Title:   title,
			URL:     link,
			Snippet: strings.TrimSpace(organic.Snippet),
		})
	}

	if results == nil {
		return []SearchResult{}, nil
	}
	return results, nil
}
