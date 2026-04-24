package search

// SearchResult represents a single search result entry
type SearchResult struct {
	Title      string `json:"title"`
	URL        string `json:"url"`
	Snippet    string `json:"snippet"`
	ResultType string `json:"resultType,omitempty"` // "website", "username", "intersection"
}


type SearchEngine interface {
	Search(query string) ([]SearchResult, error)
}
