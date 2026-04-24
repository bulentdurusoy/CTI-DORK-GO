package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func ExportJSON(data interface{}, filePath string) error {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return os.WriteFile(filePath, jsonBytes, 0644)
}

func ExportTXT(content string, filePath string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}

func FormatResultsAsTXT(domain, keyword string, data interface{}) string {
	var sb strings.Builder

	sb.WriteString("═══════════════════════════════════════════════════════════════\n")
	sb.WriteString("  CTI DORK — Google Dorking Search Report\n")
	sb.WriteString("═══════════════════════════════════════════════════════════════\n\n")
	sb.WriteString(fmt.Sprintf("  Domain:   %s\n", domain))
	if keyword != "" {
		sb.WriteString(fmt.Sprintf("  Keyword:  %s\n", keyword))
	}
	sb.WriteString(fmt.Sprintf("  Date:     %s\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString("\n")

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		sb.WriteString("Error: Could not format results\n")
		return sb.String()
	}

	var categories []struct {
		Category string `json:"category"`
		Dorks    []struct {
			Dork struct {
				Name        string `json:"name"`
				Template    string `json:"template"`
				Severity    string `json:"severity"`
				Description string `json:"description"`
			} `json:"dork"`
			Query   string `json:"query"`
			Results []struct {
				Title   string `json:"title"`
				URL     string `json:"url"`
				Snippet string `json:"snippet"`
			} `json:"results"`
		} `json:"dorks"`
	}

	if err := json.Unmarshal(jsonBytes, &categories); err != nil {
		sb.WriteString("Error: Could not parse results structure\n")
		return sb.String()
	}

	totalDorks := 0
	totalResults := 0

	for _, cat := range categories {
		sb.WriteString("───────────────────────────────────────────────────────────────\n")
		sb.WriteString(fmt.Sprintf("  📁 %s\n", cat.Category))
		sb.WriteString("───────────────────────────────────────────────────────────────\n\n")

		for _, dork := range cat.Dorks {
			totalDorks++
			sb.WriteString(fmt.Sprintf("  ▸ %s [%s]\n", dork.Dork.Name, strings.ToUpper(dork.Dork.Severity)))
			sb.WriteString(fmt.Sprintf("    Query: %s\n", dork.Query))

			if len(dork.Results) == 0 {
				sb.WriteString("    No results found\n")
			} else {
				for i, result := range dork.Results {
					totalResults++
					sb.WriteString(fmt.Sprintf("    [%d] %s\n", i+1, result.Title))
					sb.WriteString(fmt.Sprintf("        %s\n", result.URL))
					if result.Snippet != "" {
						snippet := result.Snippet
						if len(snippet) > 120 {
							snippet = snippet[:117] + "..."
						}
						sb.WriteString(fmt.Sprintf("        %s\n", snippet))
					}
				}
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("═══════════════════════════════════════════════════════════════\n")
	sb.WriteString(fmt.Sprintf("  Summary: %d dorks executed, %d results found\n", totalDorks, totalResults))
	sb.WriteString("═══════════════════════════════════════════════════════════════\n")

	return sb.String()
}
