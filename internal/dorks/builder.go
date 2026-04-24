package dorks

import "strings"

func BuildQuery(dork Dork, domain, keyword string) string {
	query := dork.Template

	if domain != "" {
		query = strings.ReplaceAll(query, "{domain}", domain)
	} else {

		query = removeDomainSegments(query)
	}

	if keyword != "" {
		query = strings.ReplaceAll(query, "{keyword}", keyword)
	} else {
		query = removeKeywordSegments(query)
	}

	for strings.Contains(query, "  ") {
		query = strings.ReplaceAll(query, "  ", " ")
	}

	return strings.TrimSpace(query)
}

func removeDomainSegments(query string) string {
	domainPatterns := []string{
		"site:{domain}",
	}

	for _, pattern := range domainPatterns {
		query = strings.ReplaceAll(query, pattern, "")
	}

	return query
}

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

func BuildGoogleURL(query string) string {
	encoded := strings.ReplaceAll(query, " ", "+")
	encoded = strings.ReplaceAll(encoded, `"`, "%22")
	encoded = strings.ReplaceAll(encoded, ":", "%3A")
	return "https://www.google.com/search?q=" + encoded
}
