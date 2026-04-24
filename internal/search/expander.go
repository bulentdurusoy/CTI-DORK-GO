package search

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ─── Expanded Query Types ────────────────────────────────────────

// ExpandedQuery represents a single auto-generated OSINT query
type ExpandedQuery struct {
	Query    string `json:"query"`    // the actual Google query string
	Platform string `json:"platform"` // "Instagram", "GitHub", "Generic", "Fallback"
	Priority int    `json:"priority"` // 1=platform, 2=inurl/intitle, 3=generic, 4=fallback
	Label    string `json:"label"`    // human-readable label for UI
}

// ─── Platform Registry ──────────────────────────────────────────

// Platform defines a search target with its query templates
type Platform struct {
	Name      string   // "Instagram", "GitHub", etc.
	Domains   []string // ["instagram.com"]
	Templates []string // query templates with {u} placeholder
}

// GetPlatforms returns the full platform registry for OSINT username searches
func GetPlatforms() []Platform {
	return []Platform{
		{
			Name:    "Facebook",
			Domains: []string{"facebook.com"},
			Templates: []string{
				`site:facebook.com "{u}"`,
				`site:facebook.com inurl:{u}`,
			},
		},
		{
			Name:    "Instagram",
			Domains: []string{"instagram.com"},
			Templates: []string{
				`site:instagram.com "{u}"`,
				`site:instagram.com inurl:{u}`,
			},
		},
		{
			Name:    "Twitter / X",
			Domains: []string{"x.com", "twitter.com"},
			Templates: []string{
				`site:x.com "{u}"`,
				`site:twitter.com "{u}"`,
				`site:x.com inurl:{u}`,
			},
		},
		{
			Name:    "YouTube",
			Domains: []string{"youtube.com"},
			Templates: []string{
				`site:youtube.com "{u}"`,
				`site:youtube.com @{u}`,
				`site:youtube.com inurl:{u}`,
			},
		},
		{
			Name:    "TikTok",
			Domains: []string{"tiktok.com"},
			Templates: []string{
				`site:tiktok.com "{u}"`,
				`site:tiktok.com inurl:{u}`,
			},
		},
		{
			Name:    "Spotify",
			Domains: []string{"open.spotify.com"},
			Templates: []string{
				`site:open.spotify.com "{u}"`,
			},
		},
		{
			Name:    "Apple Music",
			Domains: []string{"music.apple.com"},
			Templates: []string{
				`site:music.apple.com "{u}"`,
			},
		},
		{
			Name:    "SoundCloud",
			Domains: []string{"soundcloud.com"},
			Templates: []string{
				`site:soundcloud.com "{u}"`,
				`site:soundcloud.com inurl:{u}`,
			},
		},
		{
			Name:    "Telegram",
			Domains: []string{"t.me", "telegram.me"},
			Templates: []string{
				`site:t.me "{u}"`,
				`site:telegram.me "{u}"`,
			},
		},
		{
			Name:    "Snapchat",
			Domains: []string{"snapchat.com"},
			Templates: []string{
				`site:snapchat.com "{u}"`,
				`site:snapchat.com inurl:{u}`,
			},
		},
		{
			Name:    "GitHub",
			Domains: []string{"github.com"},
			Templates: []string{
				`site:github.com "{u}"`,
				`site:github.com inurl:{u}`,
				`site:github.com intitle:{u}`,
			},
		},
		{
			Name:    "LinkedIn",
			Domains: []string{"linkedin.com"},
			Templates: []string{
				`site:linkedin.com/in "{u}"`,
			},
		},
		{
			Name:    "Reddit",
			Domains: []string{"reddit.com"},
			Templates: []string{
				`site:reddit.com "{u}"`,
				`site:reddit.com inurl:{u}`,
			},
		},
		{
			Name:    "Pastebin",
			Domains: []string{"pastebin.com"},
			Templates: []string{
				`site:pastebin.com "{u}"`,
			},
		},
		{
			Name:    "GitLab",
			Domains: []string{"gitlab.com"},
			Templates: []string{
				`site:gitlab.com "{u}"`,
				`site:gitlab.com inurl:{u}`,
			},
		},
		// ── Turkey-specific ──
		{
			Name:    "NSosyal",
			Domains: []string{"nsosyal.com"},
			Templates: []string{
				`site:nsosyal.com "{u}"`,
				`site:nsosyal.com inurl:{u}`,
			},
		},
	}
}

// ─── Username Expansion ─────────────────────────────────────────

// ExpandUsername generates multiple search variations from a single username.
// For example, "bulentdurusoy" generates:
//   - bulentdurusoy
//   - "bulentdurusoy"
//   - bulent-durusoy (if split point found)
//   - bulent_durusoy (if split point found)
//   - bulent durusoy (if split point found)
//   - "bulent durusoy" (if split point found)
func ExpandUsername(username string) []string {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil
	}

	seen := make(map[string]bool)
	var variations []string

	add := func(v string) {
		v = strings.TrimSpace(v)
		if v == "" {
			return
		}
		lower := strings.ToLower(v)
		if !seen[lower] {
			seen[lower] = true
			variations = append(variations, v)
		}
	}

	// 1. Raw username
	add(username)

	// 2. If the username already contains separators, extract the "clean" form and parts
	parts := splitUsername(username)
	if len(parts) > 1 {
		joined := strings.Join(parts, "")
		add(joined)
		add(strings.Join(parts, "-"))
		add(strings.Join(parts, "_"))
		add(strings.Join(parts, " "))
		add(fmt.Sprintf(`"%s"`, strings.Join(parts, " ")))
	} else {
		// 3. Try to guess split points for concatenated usernames like "bulentdurusoy"
		guessedParts := guessSplitPoints(username)
		if len(guessedParts) > 1 {
			add(strings.Join(guessedParts, "-"))
			add(strings.Join(guessedParts, "_"))
			add(strings.Join(guessedParts, " "))
			add(fmt.Sprintf(`"%s"`, strings.Join(guessedParts, " ")))
		}
	}

	return variations
}

// splitUsername splits a username by common separators (-, _, .)
func splitUsername(username string) []string {
	// Replace common separators with a standard one
	normalized := strings.NewReplacer("-", " ", "_", " ", ".", " ").Replace(username)
	parts := strings.Fields(normalized)
	if len(parts) <= 1 {
		return []string{username}
	}
	return parts
}

// guessSplitPoints tries to find a valid split point in a concatenated username.
// Heuristic: tries splitting at every position and checks if both parts are >= 3 chars.
// Prefers splits where both halves look like name parts (start with consonant+vowel patterns).
func guessSplitPoints(username string) []string {
	lower := strings.ToLower(username)

	// Don't try to split very short usernames or ones with numbers
	if len(lower) < 6 {
		return []string{username}
	}

	// If it has numbers mixed with letters, don't try to split as a name
	hasLetters := false
	hasDigits := false
	for _, r := range lower {
		if unicode.IsLetter(r) {
			hasLetters = true
		}
		if unicode.IsDigit(r) {
			hasDigits = true
		}
	}
	if hasDigits && hasLetters {
		// Could be student ID + name or vice versa; don't guess
		return []string{username}
	}
	if !hasLetters {
		return []string{username}
	}

	// Try splitting at each position
	bestSplit := []string{username}
	bestScore := 0

	for i := 3; i <= len(lower)-3; i++ {
		left := lower[:i]
		right := lower[i:]

		score := scoreSplit(left, right)
		if score > bestScore {
			bestScore = score
			bestSplit = []string{lower[:i], lower[i:]}
		}
	}

	if bestScore >= 2 {
		return bestSplit
	}
	return []string{username}
}

// scoreSplit scores how likely a split is a valid name split.
// Higher = more likely to be a real name boundary.
func scoreSplit(left, right string) int {
	score := 0

	// Both parts should be at least 3 chars
	if len(left) < 3 || len(right) < 3 {
		return 0
	}

	// Bonus for common name lengths (3-10 chars)
	if len(left) >= 3 && len(left) <= 10 {
		score++
	}
	if len(right) >= 3 && len(right) <= 10 {
		score++
	}

	// Bonus if both start with a consonant (common in names)
	vowels := "aeiouyöüı"
	if !strings.ContainsRune(vowels, rune(left[0])) {
		score++
	}
	if !strings.ContainsRune(vowels, rune(right[0])) {
		score++
	}

	// Bonus for vowel-consonant transitions at split point
	lastOfLeft := rune(left[len(left)-1])
	firstOfRight := rune(right[0])
	leftIsVowel := strings.ContainsRune(vowels, lastOfLeft)
	rightIsVowel := strings.ContainsRune(vowels, firstOfRight)
	if leftIsVowel != rightIsVowel {
		score++
	}

	return score
}

// ─── Query Generation ───────────────────────────────────────────

// GeneratePlatformQueries generates platform-specific OSINT queries for a username.
// These are Priority 1 queries — most precise, highest value.
func GeneratePlatformQueries(username string) []ExpandedQuery {
	platforms := GetPlatforms()
	username = strings.TrimSpace(username)
	if username == "" {
		return nil
	}

	var queries []ExpandedQuery

	for _, p := range platforms {
		for _, tmpl := range p.Templates {
			query := strings.ReplaceAll(tmpl, "{u}", username)
			queries = append(queries, ExpandedQuery{
				Query:    query,
				Platform: p.Name,
				Priority: 1,
				Label:    fmt.Sprintf("%s — %s", p.Name, username),
			})
		}
	}

	return queries
}

// GenerateGenericQueries generates broad inurl/intitle/quoted queries.
// These are Priority 2 queries — wider net, still targeted.
func GenerateGenericQueries(username string) []ExpandedQuery {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil
	}

	variations := ExpandUsername(username)
	var queries []ExpandedQuery

	// For the raw username: inurl, intitle, and quoted
	queries = append(queries, ExpandedQuery{
		Query:    fmt.Sprintf(`inurl:%s`, username),
		Platform: "Generic",
		Priority: 2,
		Label:    fmt.Sprintf("inurl:%s", username),
	})
	queries = append(queries, ExpandedQuery{
		Query:    fmt.Sprintf(`intitle:%s`, username),
		Platform: "Generic",
		Priority: 2,
		Label:    fmt.Sprintf("intitle:%s", username),
	})
	queries = append(queries, ExpandedQuery{
		Query:    fmt.Sprintf(`"%s"`, username),
		Platform: "Generic",
		Priority: 2,
		Label:    fmt.Sprintf(`Quoted: "%s"`, username),
	})

	// For expanded variations (space-separated name etc.)
	for _, v := range variations {
		if v == username {
			continue
		}
		// Only add quoted version for space-containing variations
		if strings.Contains(v, " ") && !strings.HasPrefix(v, `"`) {
			queries = append(queries, ExpandedQuery{
				Query:    fmt.Sprintf(`"%s"`, v),
				Platform: "Generic",
				Priority: 3,
				Label:    fmt.Sprintf(`Variation: "%s"`, v),
			})
		}
	}

	return queries
}

// GenerateFallbackQueries generates broad OR-based queries for when results are low.
// These are Priority 4 queries — widest net.
func GenerateFallbackQueries(username string) []ExpandedQuery {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil
	}

	var queries []ExpandedQuery

	// Try username with common profile-related terms
	queries = append(queries, ExpandedQuery{
		Query:    fmt.Sprintf(`"%s" profile OR account OR user`, username),
		Platform: "Fallback",
		Priority: 4,
		Label:    fmt.Sprintf("Broad: %s + profile keywords", username),
	})

	// If we can guess a name split, add an OR query
	parts := splitUsername(username)
	if len(parts) <= 1 {
		parts = guessSplitPoints(username)
	}
	if len(parts) > 1 {
		spacedName := strings.Join(parts, " ")
		queries = append(queries, ExpandedQuery{
			Query:    fmt.Sprintf(`"%s" OR "%s"`, username, spacedName),
			Platform: "Fallback",
			Priority: 4,
			Label:    fmt.Sprintf(`Fallback: "%s" OR "%s"`, username, spacedName),
		})
	}

	// Document search
	queries = append(queries, ExpandedQuery{
		Query:    fmt.Sprintf(`"%s" filetype:pdf OR filetype:doc OR filetype:xlsx`, username),
		Platform: "Fallback",
		Priority: 4,
		Label:    fmt.Sprintf("Documents: %s", username),
	})

	return queries
}

// ─── Normalized Comparison ──────────────────────────────────────

// NormalizeForComparison strips separators (-, _, ., spaces) and lowercases
// for fuzzy username matching.
// "bulent-durusoy" → "bulentdurusoy"
// "Bulent_Durusoy" → "bulentdurusoy"
var separatorRegex = regexp.MustCompile(`[-_.\s]+`)

func NormalizeForComparison(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return separatorRegex.ReplaceAllString(s, "")
}

// ─── Intersection Queries ───────────────────────────────────────

// GenerateIntersectionQueries generates queries combining domain + username.
// Used when both inputs are provided (A ∩ B mode).
func GenerateIntersectionQueries(domain, username string) []ExpandedQuery {
	domain = strings.TrimSpace(domain)
	username = strings.TrimSpace(username)
	if domain == "" || username == "" {
		return nil
	}

	var queries []ExpandedQuery

	// Direct site + username
	queries = append(queries, ExpandedQuery{
		Query:    fmt.Sprintf(`site:%s "%s"`, domain, username),
		Platform: "Intersection",
		Priority: 1,
		Label:    fmt.Sprintf(`%s + "%s"`, domain, username),
	})
	queries = append(queries, ExpandedQuery{
		Query:    fmt.Sprintf(`site:%s inurl:%s`, domain, username),
		Platform: "Intersection",
		Priority: 1,
		Label:    fmt.Sprintf(`%s inurl:%s`, domain, username),
	})
	queries = append(queries, ExpandedQuery{
		Query:    fmt.Sprintf(`site:%s intitle:%s`, domain, username),
		Platform: "Intersection",
		Priority: 2,
		Label:    fmt.Sprintf(`%s intitle:%s`, domain, username),
	})

	// Variations
	variations := ExpandUsername(username)
	for _, v := range variations {
		if v == username {
			continue
		}
		if strings.Contains(v, " ") && !strings.HasPrefix(v, `"`) {
			queries = append(queries, ExpandedQuery{
				Query:    fmt.Sprintf(`site:%s "%s"`, domain, v),
				Platform: "Intersection",
				Priority: 3,
				Label:    fmt.Sprintf(`%s + "%s"`, domain, v),
			})
		}
	}

	return queries
}
