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

// Result Types 

type CategoryResult struct {
	Category string       `json:"category"`
	Dorks    []DorkResult `json:"dorks"`
}

type DorkResult struct {
	Dork    dorks.Dork            `json:"dork"`
	Query   string                `json:"query"`
	Results []search.SearchResult `json:"results"`
}

type SearchResponse struct {
	Mode         string           `json:"mode"`         
	ModeLabel    string           `json:"modeLabel"`    
	IsMockMode   bool             `json:"isMockMode"`   
	IsLiveMode   bool             `json:"isLiveMode"`   
	FetchedCount int              `json:"fetchedCount"` 
	Categories   []CategoryResult `json:"categories"`
}

// App Struct 

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

func NewApp() *App {
	return &App{
		favStore:   storage.NewFavoriteStore(),
		engine:     search.NewLiveSearchEngine(),
		isMockMode: false,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Category & Dork Queries 

func (a *App) GetCategories() []string {
	return dorks.GetCategories()
}

func (a *App) GetAllDorks() []dorks.Dork {
	return dorks.GetAllDorks()
}

func (a *App) GetDorksByCategory(category string) []dorks.Dork {
	return dorks.GetDorksByCategory(category)
}

func (a *App) GetCategoryCounts() map[string]int {
	counts := make(map[string]int)
	for _, d := range dorks.GetAllDorks() {
		counts[string(d.Category)]++
	}
	return counts
}

// Mode & Status 

func (a *App) GetSearchMode(websiteOrDomain, username string) string {
	input := search.SearchInput{
		WebsiteOrDomain: strings.TrimSpace(websiteOrDomain),
		Username:        strings.TrimSpace(username),
	}
	return string(search.DetectSearchMode(input))
}

func (a *App) IsMockMode() bool {
	return a.isMockMode
}

//API Key Management 

func (a *App) SetAPIKey(key string) error {
	cfg := search.GetCachedConfig()
	cfg.SerperAPIKey = strings.TrimSpace(key)
	return search.SaveConfig(cfg)
}

func (a *App) GetAPIKey() string {
	cfg := search.GetCachedConfig()
	return search.MaskAPIKey(cfg.SerperAPIKey)
}

func (a *App) HasAPIKey() bool {
	return search.HasAPIKey()
}

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


	type queryJob struct {
		Index    int
		Query    string
		Category string
		Label    string
		Platform string
		Priority int
		IsDork   bool      
		Dork     dorks.Dork 
	}

	var allJobs []queryJob
	jobIndex := 0

	//  Hardcoded dorks from the dork library (for domain-based searches)
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

	if username != "" {
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

	for _, j := range allJobs {
		jobsChan <- j
	}
	close(jobsChan)

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	globalSeen := make(map[string]bool)
	var fetchedCount int32

	categoryOrder := []string{}
	resultMap := make(map[string]*CategoryResult)

	allResults := make([]queryResult, 0, len(allJobs))
	for r := range resultsChan {
		allResults = append(allResults, r)
	}

	for i := 1; i < len(allResults); i++ {
		for j := i; j > 0 && allResults[j].Job.Priority < allResults[j-1].Job.Priority; j-- {
			allResults[j], allResults[j-1] = allResults[j-1], allResults[j]
		}
	}

	for _, qr := range allResults {
		cat := qr.Job.Category

		if _, exists := resultMap[cat]; !exists {
			resultMap[cat] = &CategoryResult{
				Category: cat,
				Dorks:    []DorkResult{},
			}
			categoryOrder = append(categoryOrder, cat)
		}

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

		if qr.Job.IsDork {
			resultMap[cat].Dorks = append(resultMap[cat].Dorks, DorkResult{
				Dork:    qr.Job.Dork,
				Query:   qr.Job.Query,
				Results: uniqueResults,
			})
		} else {
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

func (a *App) OpenInBrowser(query string) error {
	return search.OpenInBrowser(query)
}

func (a *App) ToggleFavorite(dorkID string) bool {
	if a.favStore.IsFavorite(dorkID) {
		a.favStore.RemoveFavorite(dorkID)
		return false
	}
	a.favStore.AddFavorite(dorkID)
	return true
}

func (a *App) GetFavorites() []string {
	favs := a.favStore.GetFavorites()
	if favs == nil {
		return []string{}
	}
	return favs
}

func (a *App) IsFavorite(dorkID string) bool {
	return a.favStore.IsFavorite(dorkID)
}

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
		return nil 
	}

	exportData := map[string]interface{}{
		"domain":     a.lastDomain,
		"keyword":    a.lastKeyword,
		"mode":       string(a.lastMode),
		"categories": a.lastResults,
	}

	return storage.ExportJSON(exportData, filePath)
}

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
		return nil 
	}

	content := storage.FormatResultsAsTXT(a.lastDomain, a.lastKeyword, a.lastResults)
	return storage.ExportTXT(content, filePath)
}
