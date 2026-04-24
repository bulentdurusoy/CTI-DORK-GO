package search

import (
	"net/url"
	"strings"
)

// Search Mode Types 
type SearchMode string

const (
	ModeWebsiteOnly SearchMode = "WEBSITE_ONLY"
	ModeUsernameOnly SearchMode = "USERNAME_ONLY"
	ModeIntersection SearchMode = "INTERSECTION"
)

type SearchInput struct {
	WebsiteOrDomain string `json:"websiteOrDomain"`
	Username        string `json:"username"`
}

// Mode Detection 
func DetectSearchMode(input SearchInput) SearchMode {
	hasDomain := strings.TrimSpace(input.WebsiteOrDomain) != ""
	hasUsername := strings.TrimSpace(input.Username) != ""

	switch {
	case hasDomain && hasUsername:
		return ModeIntersection
	case hasUsername:
		return ModeUsernameOnly
	default:
		return ModeWebsiteOnly
	}
}

func ModeLabel(mode SearchMode) string {
	switch mode {
	case ModeWebsiteOnly:
		return "Website Mode"
	case ModeUsernameOnly:
		return "Username Mode"
	case ModeIntersection:
		return "Intersection Mode"
	default:
		return "Unknown Mode"
	}
}

//Domain Helpers

// CleanDomain strips 
func CleanDomain(raw string) string {
	d := strings.TrimSpace(raw)
	d = strings.TrimPrefix(d, "https://")
	d = strings.TrimPrefix(d, "http://")
	d = strings.TrimSuffix(d, "/")
	d = strings.TrimSpace(d)
	return d
}

func NormalizeURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	raw = strings.ToLower(raw)
	raw = strings.TrimSuffix(raw, "/")
	// Remove fragment
	if idx := strings.Index(raw, "#"); idx != -1 {
		raw = raw[:idx]
	}
	return raw
}

// Result Validation 

func ValidateResult(result SearchResult, mode SearchMode, input SearchInput) bool {
	urlStr := strings.TrimSpace(result.URL)
	title := strings.TrimSpace(result.Title)

	if title == "" {
		return false
	}

	if urlStr == "" {
		return false
	}
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return false
	}

	switch mode {
	case ModeWebsiteOnly:
		if !MatchesDomain(urlStr, input.WebsiteOrDomain) {
			return false
		}
	case ModeUsernameOnly:
		if !MatchesUsername(result, input.Username) {
			return false
		}
	case ModeIntersection:
		if !MatchesDomain(urlStr, input.WebsiteOrDomain) {
			return false
		}
		if !MatchesUsername(result, input.Username) {
			return false
		}
	}

	return true
}

// Matching Helpers 

func MatchesDomain(resultURL, domain string) bool {
	if domain == "" {
		return false
	}

	domain = strings.ToLower(CleanDomain(domain))
	if domain == "" {
		return false
	}

	parsed, err := url.Parse(strings.ToLower(resultURL))
	if err != nil || parsed.Host == "" {
		return false
	}

	host := parsed.Hostname() 

	if host == domain {
		return true
	}
	if strings.HasSuffix(host, "."+domain) {
		return true
	}

	return false
}

// MatchesUsername checks

func MatchesUsername(result SearchResult, username string) bool {
	if username == "" {
		return false
	}

	uLower := strings.ToLower(strings.TrimSpace(username))
	if uLower == "" {
		return false
	}

	uNormalized := NormalizeForComparison(username)

	urlLower := strings.ToLower(result.URL)
	if strings.Contains(urlLower, uLower) {
		return true
	}
	if uNormalized != "" && strings.Contains(NormalizeForComparison(result.URL), uNormalized) {
		return true
	}

	titleLower := strings.ToLower(result.Title)
	if strings.Contains(titleLower, uLower) {
		return true
	}
	if uNormalized != "" && strings.Contains(NormalizeForComparison(result.Title), uNormalized) {
		return true
	}

	snippetLower := strings.ToLower(result.Snippet)
	if strings.Contains(snippetLower, uLower) {
		return true
	}
	if uNormalized != "" && strings.Contains(NormalizeForComparison(result.Snippet), uNormalized) {
		return true
	}

	return false
}

//  Filtering & Deduplication 

func FilterResultsByMode(results []SearchResult, mode SearchMode, input SearchInput) []SearchResult {
	var valid []SearchResult
	for _, r := range results {
		if ValidateResult(r, mode, input) {
			valid = append(valid, r)
		}
	}
	if valid == nil {
		return []SearchResult{}
	}
	return valid
}

func DeduplicateResults(results []SearchResult) []SearchResult {
	seen := make(map[string]bool)
	var unique []SearchResult

	for _, r := range results {
		normalized := NormalizeURL(r.URL)
		if normalized == "" || seen[normalized] {
			continue
		}
		seen[normalized] = true
		unique = append(unique, r)
	}

	if unique == nil {
		return []SearchResult{}
	}
	return unique
}

// full validation pipeline:

func ValidateAndFilterResults(results []SearchResult, mode SearchMode, input SearchInput) []SearchResult {
	filtered := FilterResultsByMode(results, mode, input)
	return DeduplicateResults(filtered)
}
