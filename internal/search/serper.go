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

//Serper.dev API Client 

const serperEndpoint = "https://google.serper.dev/search"

// SerperSearchEngine executes real Google searches via the Serper.dev API.
type SerperSearchEngine struct {
	apiKey     string
	httpClient *http.Client
}

func NewSerperSearchEngine(apiKey string) *SerperSearchEngine {
	return &SerperSearchEngine{
		apiKey: strings.TrimSpace(apiKey),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type serperRequest struct {
	Q   string `json:"q"`
	Num int    `json:"num,omitempty"` 
}

type serperResponse struct {
	Organic []serperOrganicResult `json:"organic"`
}

type serperOrganicResult struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}

// Search executes a real Google search via Serper.dev API
func (s *SerperSearchEngine) Search(query string) ([]SearchResult, error) {
	if s.apiKey == "" {
		return []SearchResult{}, fmt.Errorf("no Serper API key configured")
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return []SearchResult{}, nil
	}

	payload := serperRequest{
		Q:   query,
		Num: 10,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return []SearchResult{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", serperEndpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return []SearchResult{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-API-KEY", s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return []SearchResult{}, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []SearchResult{}, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return []SearchResult{}, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var serperResp serperResponse
	if err := json.Unmarshal(body, &serperResp); err != nil {
		return []SearchResult{}, fmt.Errorf("failed to parse API response: %w", err)
	}

	var results []SearchResult
	for _, organic := range serperResp.Organic {
		link := strings.TrimSpace(organic.Link)
		title := strings.TrimSpace(organic.Title)

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
