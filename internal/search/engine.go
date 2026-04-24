package search

// SearchResult represents a single search result entry
type SearchResult struct {
	Title      string `json:"title"`
	URL        string `json:"url"`
	Snippet    string `json:"snippet"`
	ResultType string `json:"resultType,omitempty"` // "website", "username", "intersection"
}

// SearchEngine defines the interface for executing search queries.
// Implementations:
//   - LiveSearchEngine  — production engine (browser mode, no fabrication)
//   - MockSearchEngine  — debug/demo only (must be explicitly enabled)
type SearchEngine interface {
	Search(query string) ([]SearchResult, error)
}
