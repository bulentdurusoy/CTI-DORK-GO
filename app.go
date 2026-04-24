package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"cti-dork/internal/dorks"
	"cti-dork/internal/search"
	"cti-dork/internal/storage"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ─── Result Types ────────────────────────────────────────────────

// CategoryResult holds all dork results for a single category
type CategoryResult struct {
	Category string       `json:"category"`
	Dorks    []DorkResult `json:"dorks"`
}

// DorkResult holds a single dork's query and its search results
type DorkResult struct {
	Dork    dorks.Dork            `json:"dork"`
	Query   string                `json:"query"`
	Results []search.SearchResult `json:"results"`
}

// SearchResponse wraps the search results with mode metadata
type SearchResponse struct {
	Mode         string           `json:"mode"`         // "WEBSITE_ONLY", "USERNAME_ONLY", "INTERSECTION"
	ModeLabel    string           `json:"modeLabel"`    // "Website Mode", etc.
	IsMockMode   bool             `json:"isMockMode"`   // true only if debug mock engine is active
	IsLiveMode   bool             `json:"isLiveMode"`   // true if API key is active
	FetchedCount int              `json:"fetchedCount"` // total real results fetched
	Categories   []CategoryResult `json:"categories"`
}

// ─── App Struct ──────────────────────────────────────────────────

// App is the main application struct bound to the Wails frontend.
// All public methods are exposed as JavaScript functions.
type App struct {
	ctx         context.Context
	favStore    *storage.FavoriteStore
	engine      search.SearchEngine
	isMockMode  bool
	lastResults []CategoryResult
	lastDomain  string
	lastKeyword string
	lastMode    search.SearchMode
}

// NewApp creates a new App instance with initialized dependencies.
// Uses LiveSearchEngine by default — no fabricated data in live mode.
func NewApp() *App {
	return &App{
		favStore:   storage.NewFavoriteStore(),
		engine:     search.NewLiveSearchEngine(),
		isMockMode: false,
	}
}

// startup is called when the Wails app starts. Receives the application context.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ─── Category & Dork Queries ─────────────────────────────────────

// GetCategories returns all available dork categories
func (a *App) GetCategories() []string {
	return dorks.GetCategories()
}

// GetAllDorks returns all 45 dork definitions
func (a *App) GetAllDorks() []dorks.Dork {
	return dorks.GetAllDorks()
}

// GetDorksByCategory returns dorks filtered by category name
func (a *App) GetDorksByCategory(category string) []dorks.Dork {
	return dorks.GetDorksByCategory(category)
}

// GetCategoryCounts returns a map of category name → number of dorks
func (a *App) GetCategoryCounts() map[string]int {
	counts := make(map[string]int)
	for _, d := range dorks.GetAllDorks() {
		counts[string(d.Category)]++
	}
	return counts
}

// ─── Mode & Status ───────────────────────────────────────────────

// GetSearchMode returns the search mode string for the given inputs.
// Called by frontend before/during search to determine which mode badge to show.
func (a *App) GetSearchMode(websiteOrDomain, username string) string {
	input := search.SearchInput{
		WebsiteOrDomain: strings.TrimSpace(websiteOrDomain),
		Username:        strings.TrimSpace(username),
	}
	return string(search.DetectSearchMode(input))
}

// IsMockMode returns true only if the debug mock engine is active
func (a *App) IsMockMode() bool {
	return a.isMockMode
}

// ─── API Key Management ──────────────────────────────────────────

// SetAPIKey saves the Serper API key to the config
func (a *App) SetAPIKey(key string) error {
	cfg := search.GetCachedConfig()
	cfg.SerperAPIKey = strings.TrimSpace(key)
	return search.SaveConfig(cfg)
}

// GetAPIKey returns a masked version of the API key for UI display
func (a *App) GetAPIKey() string {
	cfg := search.GetCachedConfig()
	return search.MaskAPIKey(cfg.SerperAPIKey)
}

// HasAPIKey checks if an API key is configured
func (a *App) HasAPIKey() bool {
	return search.HasAPIKey()
}

// ─── Search Execution ────────────────────────────────────────────

// ExecuteSearch generates dork queries AND expanded OSINT queries, runs them
// through the search engine concurrently, validates, deduplicates, and returns
// results grouped by category.
//
// When a username is provided, the system automatically:
//  1. Generates platform-specific queries (14+ platforms × 2-3 queries each)
//  2. Generates generic inurl/intitle/quoted queries
//  3. Runs all queries concurrently (8 workers)
//  4. If results < threshold, runs fallback broad queries
//  5. Deduplicates globally across all queries
//
// Category filter only controls UI visibility — all queries always execute.
func (a *App) ExecuteSearch(websiteOrDomain string, username string, categories []string) (*SearchResponse, error) {
	websiteOrDomain = strings.TrimSpace(websiteOrDomain)
	username = strings.TrimSpace(username)

	if websiteOrDomain == "" && username == "" {
		return nil, fmt.Errorf("please enter a website/domain or a username")
	}

	domain := ""
	if websiteOrDomain != "" {
		domain = search.CleanDomain(websiteOrDomain)
	}

	input := search.SearchInput{
		WebsiteOrDomain: domain,
		Username:        username,
	}
	mode := search.DetectSearchMode(input)

	// ── Phase 1: Collect all queries to execute ──────────────────

	type queryJob struct {
		Index    int
		Query    string
		Category string
		Label    string
		Platform string
		Priority int
		IsDork   bool      // true = from hardcoded dork library
		Dork     dorks.Dork // only set if IsDork
	}

	var allJobs []queryJob
	jobIndex := 0

	// 1a. Hardcoded dorks from the dork library (for domain-based searches)
	allDorks := dorks.GetDorksByMode(string(mode), domain, username)
	for _, d := range allDorks {
		query := dorks.BuildQuery(d, domain, username)
		allJobs = append(allJobs, queryJob{
			Index:    jobIndex,
			Query:    query,
			Category: string(d.Category),
			Label:    d.Name,
			Platform: "Dork",
			Priority: 2,
			IsDork:   true,
			Dork:     d,
		})
		jobIndex++
	}

	// 1b. Expanded OSINT queries (only when username is provided)
	if username != "" {
		// Platform-specific queries (Priority 1)
		platformQueries := search.GeneratePlatformQueries(username)
		for _, eq := range platformQueries {
			allJobs = append(allJobs, queryJob{
				Index:    jobIndex,
				Query:    eq.Query,
				Category: string(dorks.CategoryOSINT),
				Label:    eq.Label,
				Platform: eq.Platform,
				Priority: eq.Priority,
			})
			jobIndex++
		}

		// Generic queries (Priority 2-3)
		genericQueries := search.GenerateGenericQueries(username)
		for _, eq := range genericQueries {
			allJobs = append(allJobs, queryJob{
				Index:    jobIndex,
				Query:    eq.Query,
				Category: string(dorks.CategoryOSINT),
				Label:    eq.Label,
				Platform: eq.Platform,
				Priority: eq.Priority,
			})
			jobIndex++
		}

		// Intersection queries (when both domain and username)
		if domain != "" {
			intersectionQueries := search.GenerateIntersectionQueries(domain, username)
			for _, eq := range intersectionQueries {
				allJobs = append(allJobs, queryJob{
					Index:    jobIndex,
					Query:    eq.Query,
					Category: string(dorks.CategoryOSINT),
					Label:    eq.Label,
					Platform: eq.Platform,
					Priority: eq.Priority,
				})
				jobIndex++
			}
		}
	}

	// ── Phase 2: Execute all queries concurrently ────────────────

	type queryResult struct {
		Job     queryJob
		Results []search.SearchResult
	}

	numWorkers := 8
	jobsChan := make(chan queryJob, len(allJobs))
	resultsChan := make(chan queryResult, len(allJobs))
	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobsChan {
				rawResults, err := a.engine.Search(j.Query)
				if err != nil {
					rawResults = []search.SearchResult{}
				}
				// Basic validation only (URL format, title non-empty)
				// Skip mode-specific filtering here — we do it globally after
				var valid []search.SearchResult
				for _, r := range rawResults {
					urlStr := strings.TrimSpace(r.URL)
					title := strings.TrimSpace(r.Title)
					if urlStr == "" || title == "" {
						continue
					}
					if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
						continue
					}
					valid = append(valid, r)
				}
				if valid == nil {
					valid = []search.SearchResult{}
				}
				resultsChan <- queryResult{Job: j, Results: valid}
			}
		}()
	}

	// Send all jobs
	for _, j := range allJobs {
		jobsChan <- j
	}
	close(jobsChan)

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// ── Phase 3: Collect and deduplicate ─────────────────────────

	// Global deduplication across ALL queries
	globalSeen := make(map[string]bool)
	var fetchedCount int32

	// Build category structure
	categoryOrder := []string{}
	resultMap := make(map[string]*CategoryResult)

	// Temporary storage for ordered collection
	allResults := make([]queryResult, 0, len(allJobs))
	for r := range resultsChan {
		allResults = append(allResults, r)
	}

	// Sort by priority (lower = higher priority, gets to "claim" URLs first)
	// Simple insertion sort is fine for <100 items
	for i := 1; i < len(allResults); i++ {
		for j := i; j > 0 && allResults[j].Job.Priority < allResults[j-1].Job.Priority; j-- {
			allResults[j], allResults[j-1] = allResults[j-1], allResults[j]
		}
	}

	// Process results in priority order
	for _, qr := range allResults {
		cat := qr.Job.Category

		// Initialize category
		if _, exists := resultMap[cat]; !exists {
			resultMap[cat] = &CategoryResult{
				Category: cat,
				Dorks:    []DorkResult{},
			}
			categoryOrder = append(categoryOrder, cat)
		}

		// Deduplicate results globally
		var uniqueResults []search.SearchResult
		for _, r := range qr.Results {
			normalized := search.NormalizeURL(r.URL)
			if normalized != "" && !globalSeen[normalized] {
				globalSeen[normalized] = true
				uniqueResults = append(uniqueResults, r)
			}
		}
		if uniqueResults == nil {
			uniqueResults = []search.SearchResult{}
		}

		atomic.AddInt32(&fetchedCount, int32(len(uniqueResults)))

		// Build the DorkResult entry
		if qr.Job.IsDork {
			resultMap[cat].Dorks = append(resultMap[cat].Dorks, DorkResult{
				Dork:    qr.Job.Dork,
				Query:   qr.Job.Query,
				Results: uniqueResults,
			})
		} else {
			// Expanded OSINT query — create a synthetic Dork entry for UI
			syntheticDork := dorks.Dork{
				ID:           fmt.Sprintf("osint-auto-%d", qr.Job.Index),
				Name:         qr.Job.Label,
				Template:     qr.Job.Query,
				Category:     dorks.CategoryOSINT,
				Description:  fmt.Sprintf("Auto-generated %s query", qr.Job.Platform),
				Severity:     "medium",
				NeedsDomain:  false,
				NeedsKeyword: true,
			}
			resultMap[cat].Dorks = append(resultMap[cat].Dorks, DorkResult{
				Dork:    syntheticDork,
				Query:   qr.Job.Query,
				Results: uniqueResults,
			})
		}
	}

	// ── Phase 4: Fallback (if results are low and username provided) ──

	if username != "" && int(fetchedCount) < 5 {
		fallbackQueries := search.GenerateFallbackQueries(username)
		for _, eq := range fallbackQueries {
			rawResults, err := a.engine.Search(eq.Query)
			if err != nil {
				continue
			}
			var valid []search.SearchResult
			for _, r := range rawResults {
				urlStr := strings.TrimSpace(r.URL)
				title := strings.TrimSpace(r.Title)
				if urlStr == "" || title == "" {
					continue
				}
				if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
					continue
				}
				normalized := search.NormalizeURL(r.URL)
				if normalized != "" && !globalSeen[normalized] {
					globalSeen[normalized] = true
					valid = append(valid, r)
				}
			}
			if len(valid) == 0 {
				valid = []search.SearchResult{}
			}
			atomic.AddInt32(&fetchedCount, int32(len(valid)))

			cat := string(dorks.CategoryOSINT)
			if _, exists := resultMap[cat]; !exists {
				resultMap[cat] = &CategoryResult{
					Category: cat,
					Dorks:    []DorkResult{},
				}
				categoryOrder = append(categoryOrder, cat)
			}

			syntheticDork := dorks.Dork{
				ID:           fmt.Sprintf("osint-fallback-%s", eq.Platform),
				Name:         eq.Label,
				Template:     eq.Query,
				Category:     dorks.CategoryOSINT,
				Description:  fmt.Sprintf("Fallback query: %s", eq.Platform),
				Severity:     "low",
				NeedsDomain:  false,
				NeedsKeyword: true,
			}
			resultMap[cat].Dorks = append(resultMap[cat].Dorks, DorkResult{
				Dork:    syntheticDork,
				Query:   eq.Query,
				Results: valid,
			})
		}
	}

	// ── Phase 5: Build final response ────────────────────────────

	var categoryResults []CategoryResult
	for _, cat := range categoryOrder {
		categoryResults = append(categoryResults, *resultMap[cat])
	}

	a.lastResults = categoryResults
	a.lastDomain = domain
	a.lastKeyword = username
	a.lastMode = mode

	liveEngine, isLive := a.engine.(*search.LiveSearchEngine)
	isLiveMode := false
	if isLive {
		isLiveMode = liveEngine.IsLiveMode()
	}

	return &SearchResponse{
		Mode:         string(mode),
		ModeLabel:    search.ModeLabel(mode),
		IsMockMode:   a.isMockMode,
		IsLiveMode:   isLiveMode,
		FetchedCount: int(fetchedCount),
		Categories:   categoryResults,
	}, nil
}

// ─── Browser ─────────────────────────────────────────────────────

// OpenInBrowser opens a Google search query in the default web browser
func (a *App) OpenInBrowser(query string) error {
	return search.OpenInBrowser(query)
}

// ─── Favorites ───────────────────────────────────────────────────

// ToggleFavorite adds or removes a dork from favorites. Returns true if now a favorite.
func (a *App) ToggleFavorite(dorkID string) bool {
	if a.favStore.IsFavorite(dorkID) {
		a.favStore.RemoveFavorite(dorkID)
		return false
	}
	a.favStore.AddFavorite(dorkID)
	return true
}

// GetFavorites returns all favorite dork IDs
func (a *App) GetFavorites() []string {
	favs := a.favStore.GetFavorites()
	if favs == nil {
		return []string{}
	}
	return favs
}

// IsFavorite checks if a specific dork is a favorite
func (a *App) IsFavorite(dorkID string) bool {
	return a.favStore.IsFavorite(dorkID)
}

// ─── Export ──────────────────────────────────────────────────────

// ExportResultsJSON exports the last search results to a JSON file via save dialog
func (a *App) ExportResultsJSON() error {
	if a.lastResults == nil {
		return fmt.Errorf("no results to export — run a search first")
	}

	filePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Export Results as JSON",
		DefaultFilename: "dork-results.json",
		Filters: []runtime.FileFilter{
			{DisplayName: "JSON Files (*.json)", Pattern: "*.json"},
		},
	})
	if err != nil {
		return err
	}
	if filePath == "" {
		return nil // User cancelled
	}

	exportData := map[string]interface{}{
		"domain":     a.lastDomain,
		"keyword":    a.lastKeyword,
		"mode":       string(a.lastMode),
		"categories": a.lastResults,
	}

	return storage.ExportJSON(exportData, filePath)
}

// ExportResultsTXT exports the last search results to a formatted text file
func (a *App) ExportResultsTXT() error {
	if a.lastResults == nil {
		return fmt.Errorf("no results to export — run a search first")
	}

	filePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Export Results as TXT",
		DefaultFilename: "dork-results.txt",
		Filters: []runtime.FileFilter{
			{DisplayName: "Text Files (*.txt)", Pattern: "*.txt"},
		},
	})
	if err != nil {
		return err
	}
	if filePath == "" {
		return nil // User cancelled
	}

	content := storage.FormatResultsAsTXT(a.lastDomain, a.lastKeyword, a.lastResults)
	return storage.ExportTXT(content, filePath)
}
