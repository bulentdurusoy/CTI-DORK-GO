/* ═══════════════════════════════════════════════════════════════
   CTI DORK — Main Application Logic
   Handles state management, Wails bindings, and UI orchestration
   ═══════════════════════════════════════════════════════════════ */

// ── Application State ─────────────────────────────────────────
const AppState = {
    domain: '',
    keyword: '',
    results: null,        // SearchResponse from backend
    categories: null,     // CategoryResult[] extracted from response
    favorites: [],
    selectedCategories: [],
    allCategories: [],
    categoryCounts: {},
    showFavoritesOnly: false,
    showCriticalOnly: false,
    isSearching: false,
    isMockMode: false,
    isLiveMode: false,
    fetchedCount: 0,
    searchMode: '',       // 'WEBSITE_ONLY', 'USERNAME_ONLY', 'INTERSECTION'
    searchModeLabel: '',  // 'Website Mode', etc.
};

// ── Wails Go Bindings (via window.go) ─────────────────────────
const Go = {
    GetCategories:      () => window.go.main.App.GetCategories(),
    GetAllDorks:        () => window.go.main.App.GetAllDorks(),
    GetCategoryCounts:  () => window.go.main.App.GetCategoryCounts(),
    ExecuteSearch:      (domain, keyword, cats) => window.go.main.App.ExecuteSearch(domain, keyword, cats),
    OpenInBrowser:      (query) => window.go.main.App.OpenInBrowser(query),
    ToggleFavorite:     (id) => window.go.main.App.ToggleFavorite(id),
    GetFavorites:       () => window.go.main.App.GetFavorites(),
    IsFavorite:         (id) => window.go.main.App.IsFavorite(id),
    ExportResultsJSON:  () => window.go.main.App.ExportResultsJSON(),
    ExportResultsTXT:   () => window.go.main.App.ExportResultsTXT(),
    GetSearchMode:      (domain, username) => window.go.main.App.GetSearchMode(domain, username),
    IsMockMode:         () => window.go.main.App.IsMockMode(),
    SetAPIKey:          (key) => window.go.main.App.SetAPIKey(key),
    GetAPIKey:          () => window.go.main.App.GetAPIKey(),
    HasAPIKey:          () => window.go.main.App.HasAPIKey(),
};

// ── Initialization ────────────────────────────────────────────
document.addEventListener('DOMContentLoaded', async () => {
    try {
        // Load categories and counts
        AppState.allCategories = await Go.GetCategories();
        AppState.categoryCounts = await Go.GetCategoryCounts();
        AppState.favorites = await Go.GetFavorites();
        AppState.isMockMode = await Go.IsMockMode();

        // All categories selected by default
        AppState.selectedCategories = [...AppState.allCategories];

        // Render sidebar
        renderCategoryFilters();
        updateStats();
        updateAppModeBadge();

        // Setup keyboard shortcut
        document.getElementById('domain-input').addEventListener('keydown', (e) => {
            if (e.key === 'Enter') startSearch();
        });
        document.getElementById('keyword-input').addEventListener('keydown', (e) => {
            if (e.key === 'Enter') startSearch();
        });

        // Initial button state
        updateSearchButton();

        // Focus domain input
        document.getElementById('domain-input').focus();
    } catch (err) {
        console.error('Initialization error:', err);
        showToast('Failed to initialize application', 'error');
    }
});

// ── Category Filters Rendering ────────────────────────────────
function renderCategoryFilters() {
    const container = document.getElementById('category-filters');
    const html = AppState.allCategories.map(cat => {
        const count = AppState.categoryCounts[cat] || 0;
        const isActive = AppState.selectedCategories.includes(cat);
        return createCategoryFilter(cat, count, isActive);
    }).join('');
    container.innerHTML = html;
}

// ── Toggle Category Selection ─────────────────────────────────
function toggleCategory(el) {
    const category = el.dataset.category;
    const idx = AppState.selectedCategories.indexOf(category);

    if (idx > -1) {
        AppState.selectedCategories.splice(idx, 1);
        el.classList.remove('active');
    } else {
        AppState.selectedCategories.push(category);
        el.classList.add('active');
    }

    // If we have results, re-filter displayed results
    if (AppState.categories) {
        renderResults();
    }
}

// ── Search ────────────────────────────────────────────────────
async function startSearch() {
    const domainInput = document.getElementById('domain-input');
    const keywordInput = document.getElementById('keyword-input');
    const domain = domainInput.value.trim();
    const keyword = keywordInput.value.trim();

    // Validation: at least one field required
    if (!domain && !keyword) {
        showToast('Please enter a website/domain or a username', 'warning');
        domainInput.focus();
        domainInput.parentElement.style.borderColor = 'var(--severity-medium)';
        keywordInput.parentElement.style.borderColor = 'var(--severity-medium)';
        setTimeout(() => {
            domainInput.parentElement.style.borderColor = '';
            keywordInput.parentElement.style.borderColor = '';
        }, 2000);
        return;
    }

    if (AppState.isSearching) return;
    AppState.isSearching = true;
    AppState.domain = domain;
    AppState.keyword = keyword;

    // Update UI to loading state
    const searchBtn = document.getElementById('search-btn');
    searchBtn.classList.add('searching');
    searchBtn.innerHTML = `
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" class="spinning">
            <circle cx="12" cy="12" r="10" stroke-dasharray="31.4" stroke-dashoffset="10"/></svg>
        <span>Scanning...</span>
    `;

    showState('loading');
    updateProgress(0);

    try {
        // Simulate progressive loading
        let progress = 0;
        const progressInterval = setInterval(() => {
            progress = Math.min(progress + Math.random() * 15, 90);
            updateProgress(progress);
        }, 200);

        // Execute search — returns SearchResponse with mode metadata
        const response = await Go.ExecuteSearch(domain, keyword, []);
        clearInterval(progressInterval);
        updateProgress(100);

        // Store response and extract categories
        AppState.results = response;
        AppState.categories = response.categories || [];
        AppState.searchMode = response.mode || '';
        AppState.searchModeLabel = response.modeLabel || '';
        AppState.isMockMode = response.isMockMode || false;
        AppState.isLiveMode = response.isLiveMode || false;
        AppState.fetchedCount = response.fetchedCount || 0;

        // Brief pause to show 100% progress
        await sleep(300);

        // Calculate total real results (should match fetchedCount from backend)
        const totalResults = AppState.fetchedCount;

        // Update UI
        updateModeBadge(AppState.searchMode);
        updateAppModeBadge();

        // Render results
        renderResults();
        showState('results');

        // Enable export buttons
        document.getElementById('export-json-btn').disabled = false;
        document.getElementById('export-txt-btn').disabled = false;

        // Update stats — show real validated result count only
        document.getElementById('stat-results').textContent = totalResults;

        const totalDorks = AppState.categories.reduce((sum, cat) => sum + cat.dorks.length, 0);
        showToast(`Scan complete — ${totalDorks} dorks across ${AppState.categories.length} categories`, 'success');
    } catch (err) {
        console.error('Search error:', err);
        document.getElementById('error-text').textContent = err.message || 'Search execution failed';
        showState('error');
        showToast('Search failed: ' + (err.message || 'Unknown error'), 'error');
    } finally {
        AppState.isSearching = false;
        searchBtn.classList.remove('searching');
        searchBtn.innerHTML = `
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>
            <span>Search</span>
        `;
    }
}

// ── Render Results ────────────────────────────────────────────
function renderResults() {
    if (!AppState.categories) return;

    const container = document.getElementById('results-container');
    let filteredResults = [...AppState.categories];

    // Filter by selected categories
    if (AppState.selectedCategories.length < AppState.allCategories.length) {
        filteredResults = filteredResults.filter(
            cat => AppState.selectedCategories.includes(cat.category)
        );
    }

    // Filter favorites only
    if (AppState.showFavoritesOnly) {
        filteredResults = filteredResults.map(cat => ({
            ...cat,
            dorks: cat.dorks.filter(d => AppState.favorites.includes(d.dork.id))
        })).filter(cat => cat.dorks.length > 0);
    }

    // Filter critical only
    if (AppState.showCriticalOnly) {
        filteredResults = filteredResults.map(cat => ({
            ...cat,
            dorks: cat.dorks.filter(d => d.dork.severity === 'critical')
        })).filter(cat => cat.dorks.length > 0);
    }

    // Calculate totals — real validated results only
    const totalDorks = filteredResults.reduce((sum, cat) => sum + cat.dorks.length, 0);
    const totalResults = filteredResults.reduce((sum, cat) =>
        sum + cat.dorks.reduce((s, d) => s + d.results.length, 0), 0);

    // Build HTML
    let html = createResultsSummary(AppState.domain, AppState.keyword, totalDorks, totalResults, AppState.searchModeLabel, AppState.isLiveMode);

    if (filteredResults.length === 0) {
        html += '<div class="no-results" style="padding: 60px 0;">No results match your current filters</div>';
    } else {
        html += filteredResults.map(cat =>
            createCategorySection(cat, AppState.favorites, AppState.isLiveMode)
        ).join('');
    }

    container.innerHTML = html;
}

// ── Open Dork in Browser ──────────────────────────────────────
async function openDorkInBrowser(query) {
    try {
        await Go.OpenInBrowser(query);
        showToast('Opening in browser...', 'info');
    } catch (err) {
        showToast('Failed to open browser', 'error');
    }
}

// ── Open URL directly ─────────────────────────────────────────
function openUrl(url) {
    try {
        window.open(url, '_blank');
    } catch (err) {
        showToast('Failed to open URL', 'error');
    }
}

// ── Toggle Favorite ───────────────────────────────────────────
async function toggleFav(dorkId, btnEl) {
    try {
        const isFav = await Go.ToggleFavorite(dorkId);
        AppState.favorites = await Go.GetFavorites();

        // Update button visual
        if (isFav) {
            btnEl.classList.add('favorited');
            btnEl.querySelector('svg').setAttribute('fill', 'currentColor');
            showToast('Added to favorites', 'success');
        } else {
            btnEl.classList.remove('favorited');
            btnEl.querySelector('svg').setAttribute('fill', 'none');
            showToast('Removed from favorites', 'info');
        }

        // Update favorite count in stats
        document.getElementById('stat-favorites').textContent = AppState.favorites.length;

        // Re-render if favorites filter is active
        if (AppState.showFavoritesOnly) {
            renderResults();
        }
    } catch (err) {
        showToast('Failed to update favorite', 'error');
    }
}

// ── Favorites Filter ──────────────────────────────────────────
function toggleFavoritesFilter() {
    AppState.showFavoritesOnly = !AppState.showFavoritesOnly;
    const btn = document.getElementById('filter-favorites-btn');
    btn.classList.toggle('active', AppState.showFavoritesOnly);

    if (AppState.categories) renderResults();
}

// ── Critical Filter ───────────────────────────────────────────
function toggleCriticalFilter() {
    AppState.showCriticalOnly = !AppState.showCriticalOnly;
    const btn = document.getElementById('filter-critical-btn');
    btn.classList.toggle('active', AppState.showCriticalOnly);

    if (AppState.categories) renderResults();
}

// ── Export Functions ───────────────────────────────────────────
async function exportJSON() {
    try {
        await Go.ExportResultsJSON();
        showToast('Results exported as JSON', 'success');
    } catch (err) {
        if (err.message && err.message.includes('no results')) {
            showToast('No results to export — run a search first', 'warning');
        } else {
            showToast('Export failed: ' + (err.message || 'Unknown error'), 'error');
        }
    }
}

async function exportTXT() {
    try {
        await Go.ExportResultsTXT();
        showToast('Results exported as TXT', 'success');
    } catch (err) {
        if (err.message && err.message.includes('no results')) {
            showToast('No results to export — run a search first', 'warning');
        } else {
            showToast('Export failed: ' + (err.message || 'Unknown error'), 'error');
        }
    }
}

// ── UI State Management ───────────────────────────────────────
function showState(state) {
    const states = ['welcome-state', 'loading-state', 'error-state', 'empty-state', 'results-container'];

    // Hide all states first
    states.forEach(id => {
        document.getElementById(id).classList.add('hidden');
    });

    // Show the requested state
    if (state === 'welcome') {
        document.getElementById('welcome-state').classList.remove('hidden');
    } else if (state === 'loading') {
        document.getElementById('loading-state').classList.remove('hidden');
    } else if (state === 'error') {
        document.getElementById('error-state').classList.remove('hidden');
    } else if (state === 'empty') {
        document.getElementById('empty-state').classList.remove('hidden');
    } else if (state === 'results') {
        document.getElementById('results-container').classList.remove('hidden');
    }
}

function resetUI() {
    AppState.results = null;
    AppState.categories = null;
    AppState.isSearching = false;
    AppState.searchMode = '';
    AppState.searchModeLabel = '';
    document.getElementById('export-json-btn').disabled = true;
    document.getElementById('export-txt-btn').disabled = true;
    document.getElementById('stat-results').textContent = '0';
    hideModeBadge();
    showState('welcome');
}

function updateProgress(percent) {
    const fill = document.getElementById('progress-fill');
    if (fill) fill.style.width = percent + '%';
}

function updateStats() {
    const totalDorks = Object.values(AppState.categoryCounts).reduce((a, b) => a + b, 0);
    document.getElementById('stat-total-dorks').textContent = totalDorks;
    document.getElementById('stat-favorites').textContent = AppState.favorites.length;
}

// ── Mode Badge Management ─────────────────────────────────────
function updateModeBadge(mode) {
    const badge = document.getElementById('search-mode-badge');
    if (!badge) return;

    const modeConfig = {
        'WEBSITE_ONLY':  { label: 'Website Mode',      cssClass: 'mode-indicator--website' },
        'USERNAME_ONLY': { label: 'Username Mode',      cssClass: 'mode-indicator--username' },
        'INTERSECTION':  { label: 'Intersection Mode',  cssClass: 'mode-indicator--intersection' },
    };

    const config = modeConfig[mode];
    if (!config) {
        badge.style.display = 'none';
        return;
    }

    badge.textContent = config.label;
    badge.className = 'mode-indicator ' + config.cssClass;
    badge.style.display = 'inline-block';
}

function hideModeBadge() {
    const badge = document.getElementById('search-mode-badge');
    if (badge) badge.style.display = 'none';
}

function updateAppModeBadge() {
    const badge = document.getElementById('app-mode-badge');
    if (!badge) return;

    if (AppState.isMockMode) {
        badge.textContent = 'Mock Data Mode';
        badge.className = 'mode-indicator mode-indicator--mock';
    } else if (AppState.isLiveMode) {
        badge.textContent = 'Live Search Mode';
        badge.className = 'mode-indicator mode-indicator--live';
    } else {
        badge.textContent = 'Browser Mode (No API)';
        badge.className = 'mode-indicator mode-indicator--search'; // using search style for info
    }
}

// ── Settings Modal ────────────────────────────────────────────
async function openSettings() {
    const modal = document.getElementById('settings-modal');
    if (modal) {
        modal.classList.remove('hidden');
        // Fetch current masked API key
        try {
            const currentKey = await Go.GetAPIKey();
            const input = document.getElementById('api-key-input');
            input.value = currentKey;
        } catch (e) {
            console.error("Failed to load API key", e);
        }
    }
}

function closeSettings() {
    const modal = document.getElementById('settings-modal');
    if (modal) modal.classList.add('hidden');
}

async function saveSettings() {
    const input = document.getElementById('api-key-input');
    let newKey = input.value.trim();

    // If it's still masked (contains ***), ignore save
    if (newKey.includes('***')) {
        closeSettings();
        return;
    }

    try {
        await Go.SetAPIKey(newKey);
        AppState.isLiveMode = await Go.HasAPIKey();
        updateAppModeBadge();
        showToast('Settings saved successfully', 'success');
        closeSettings();
    } catch (err) {
        showToast('Failed to save settings', 'error');
    }
}

// ── Utility ───────────────────────────────────────────────────
function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

// ── Search Button State ───────────────────────────────────────
function updateSearchButton() {
    const domain = document.getElementById('domain-input').value.trim();
    const keyword = document.getElementById('keyword-input').value.trim();
    const searchBtn = document.getElementById('search-btn');
    searchBtn.disabled = (!domain && !keyword);
}
