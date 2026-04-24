package search

type LiveSearchEngine struct{}

// NewLiveSearchEngine creates a new LiveSearchEngine instance
func NewLiveSearchEngine() *LiveSearchEngine {
	return &LiveSearchEngine{}
}

// Search executes the search query.
func (l *LiveSearchEngine) Search(query string) ([]SearchResult, error) {
	cfg := GetCachedConfig()
	if cfg.SerperAPIKey != "" {
		engine := NewSerperSearchEngine(cfg.SerperAPIKey)
		return engine.Search(query)
	}

	return []SearchResult{}, nil
}

func (l *LiveSearchEngine) IsLiveMode() bool {
	return HasAPIKey()
}
