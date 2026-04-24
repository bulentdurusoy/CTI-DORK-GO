package search

import (
	"net/url"
	"strings"
)

// ─── Search Mode Types ───────────────────────────────────────────

// SearchMode represents the active search mode based on user input
type SearchMode string

const (
	// ModeWebsiteOnly — only website/domain is provided (A)
	ModeWebsiteOnly SearchMode = "WEBSITE_ONLY"
	// ModeUsernameOnly — only username is provided (B)
	ModeUsernameOnly SearchMode = "USERNAME_ONLY"
	// ModeIntersection — both website/domain and username are provided (A ∩ B)
	ModeIntersection SearchMode = "INTERSECTION"
)

// SearchInput holds the cleaned user inputs for a search
type SearchInput struct {
	WebsiteOrDomain string `json:"websiteOrDomain"`
	Username        string `json:"username"`
}

// ─── Mode Detection ──────────────────────────────────────────────

// DetectSearchMode determines the search mode from the user input.
//   - website only  → ModeWebsiteOnly
//   - username only → ModeUsernameOnly
//   - both          → ModeIntersection
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

// ModeLabel returns a human-readable label for the search mode
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

// ─── Domain Helpers ──────────────────────────────────────────────

// CleanDomain strips protocol prefix and trailing slashes from a domain input.
// "https://example.com/" → "example.com"
func CleanDomain(raw string) string {
	d := strings.TrimSpace(raw)
	d = strings.TrimPrefix(d, "https://")
	d = strings.TrimPrefix(d, "http://")
	d = strings.TrimSuffix(d, "/")
	d = strings.TrimSpace(d)
	return d
}

// NormalizeURL lowercases and trims a URL for deduplication comparisons.
// Removes trailing slashes and fragments.
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

// ─── Result Validation ──────────────────────────────────────────

// ValidateResult checks whether a single search result meets strict quality criteria.
// A result is valid only if:
//   - URL is non-empty and absolute (http:// or https://)
//   - Title is non-empty
//   - Result is relevant to the active search mode
func ValidateResult(result SearchResult, mode SearchMode, input SearchInput) bool {
	urlStr := strings.TrimSpace(result.URL)
	title := strings.TrimSpace(result.Title)

	// Title must not be empty
	if title == "" {
		return false
	}

	// URL must be non-empty and absolute
	if urlStr == "" {
		return false
	}
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return false
	}

	// Mode relevance check
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

// ─── Matching Helpers ────────────────────────────────────────────

// MatchesDomain checks if a result URL belongs to or is related to the given domain.
// Handles subdomains (e.g., "blog.example.com" matches "example.com").
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

	host := parsed.Hostname() // strips port

	// Exact match or subdomain match
	if host == domain {
		return true
	}
	if strings.HasSuffix(host, "."+domain) {
		return true
	}

	return false
}

// MatchesUsername checks if a result (URL, title, or snippet) contains the given username.
// Uses both exact case-insensitive matching AND normalized comparison.
// Normalized comparison strips separators (-, _, ., spaces) so that:
//   - "bulentdurusoy" matches "bulent-durusoy", "bulent_durusoy", "bulent.durusoy"
func MatchesUsername(result SearchResult, username string) bool {
	if username == "" {
		return false
	}

	uLower := strings.ToLower(strings.TrimSpace(username))
	if uLower == "" {
		return false
	}

	// Normalized form for fuzzy matching
	uNormalized := NormalizeForComparison(username)

	// Check URL (exact + normalized)
	urlLower := strings.ToLower(result.URL)
	if strings.Contains(urlLower, uLower) {
		return true
	}
	if uNormalized != "" && strings.Contains(NormalizeForComparison(result.URL), uNormalized) {
		return true
	}

	// Check title (exact + normalized)
	titleLower := strings.ToLower(result.Title)
	if strings.Contains(titleLower, uLower) {
		return true
	}
	if uNormalized != "" && strings.Contains(NormalizeForComparison(result.Title), uNormalized) {
		return true
	}

	// Check snippet (exact + normalized)
	snippetLower := strings.ToLower(result.Snippet)
	if strings.Contains(snippetLower, uLower) {
		return true
	}
	if uNormalized != "" && strings.Contains(NormalizeForComparison(result.Snippet), uNormalized) {
		return true
	}

	return false
}

// ─── Filtering & Deduplication ───────────────────────────────────

// FilterResultsByMode validates and filters results based on the active search mode.
// Only results that pass all validation checks are returned.
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

// DeduplicateResults removes duplicate results by normalized URL.
// The first occurrence of each URL is kept; subsequent duplicates are dropped.
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

// ValidateAndFilterResults applies the full validation pipeline:
// 1. Filter by mode relevance
// 2. Deduplicate by normalized URL
func ValidateAndFilterResults(results []SearchResult, mode SearchMode, input SearchInput) []SearchResult {
	filtered := FilterResultsByMode(results, mode, input)
	return DeduplicateResults(filtered)
}
