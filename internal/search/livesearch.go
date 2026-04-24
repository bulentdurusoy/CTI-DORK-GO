package search

// LiveSearchEngine is the production search engine used in live mode.
// It acts as a smart router:
// - If an API key is configured, it uses SerperSearchEngine to fetch real results.
// - If no API key is configured, it falls back to browser mode (returns empty results).
type LiveSearchEngine struct{}

// NewLiveSearchEngine creates a new LiveSearchEngine instance
func NewLiveSearchEngine() *LiveSearchEngine {
	return &LiveSearchEngine{}
}

// Search executes the search query.
// If an API key is present, it calls Serper API.
// Otherwise, it returns empty results (browser mode).
func (l *LiveSearchEngine) Search(query string) ([]SearchResult, error) {
	cfg := GetCachedConfig()
	if cfg.SerperAPIKey != "" {
		// Live API mode
		engine := NewSerperSearchEngine(cfg.SerperAPIKey)
		return engine.Search(query)
	}

	// Fallback browser mode
	return []SearchResult{}, nil
}

// IsLiveMode returns true if the engine will execute real API searches
func (l *LiveSearchEngine) IsLiveMode() bool {
	return HasAPIKey()
}
