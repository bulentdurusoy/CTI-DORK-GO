package dorks

import "strings"

// BuildQuery replaces {domain} and {keyword} placeholders in a dork template
// with the actual user-provided values and returns the final search query.
// Returns empty string if a required placeholder cannot be filled.
func BuildQuery(dork Dork, domain, keyword string) string {
	query := dork.Template

	// Replace {domain} placeholder — only if we have a value
	if domain != "" {
		query = strings.ReplaceAll(query, "{domain}", domain)
	} else {
		// Remove site:{domain} and other domain segments if no domain provided
		query = removeDomainSegments(query)
	}

	// Replace {keyword} placeholder — only if we have a value
	if keyword != "" {
		query = strings.ReplaceAll(query, "{keyword}", keyword)
	} else {
		// Remove keyword-dependent segments if no keyword provided
		query = removeKeywordSegments(query)
	}

	// Clean up any double spaces and trim
	for strings.Contains(query, "  ") {
		query = strings.ReplaceAll(query, "  ", " ")
	}

	return strings.TrimSpace(query)
}

// removeDomainSegments strips dork operators that reference {domain}
// when no domain is provided by the user.
func removeDomainSegments(query string) string {
	domainPatterns := []string{
		"site:{domain}",
	}

	for _, pattern := range domainPatterns {
		query = strings.ReplaceAll(query, pattern, "")
	}

	return query
}

// removeKeywordSegments strips dork operators that reference {keyword}
// when no keyword is provided by the user.
func removeKeywordSegments(query string) string {
	keywordPatterns := []string{
		"intext:{keyword}",
		"intitle:{keyword}",
		"inurl:{keyword}",
		`"{keyword}"`,
		"{keyword}",
	}

	for _, pattern := range keywordPatterns {
		query = strings.ReplaceAll(query, pattern, "")
	}

	return query
}

// BuildGoogleURL creates a full Google search URL for a given query
func BuildGoogleURL(query string) string {
	encoded := strings.ReplaceAll(query, " ", "+")
	encoded = strings.ReplaceAll(encoded, `"`, "%22")
	encoded = strings.ReplaceAll(encoded, ":", "%3A")
	return "https://www.google.com/search?q=" + encoded
}
