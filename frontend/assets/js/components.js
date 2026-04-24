/* ═══════════════════════════════════════════════════════════════
   CTI DORK — UI Component Builders
   Pure functions that return DOM element strings
   ═══════════════════════════════════════════════════════════════ */

// ── Category Color Map ────────────────────────────────────────
const CATEGORY_COLORS = {
    'File Exposure':        'var(--cat-file)',
    'Login/Admin Panels':   'var(--cat-login)',
    'Directory Listing':    'var(--cat-dir)',
    'Error Messages':       'var(--cat-error)',
    'Sensitive Data':       'var(--cat-sensitive)',
    'SEO/Indexing':         'var(--cat-seo)',
    'Database Exposure':    'var(--cat-db)',
    'Configuration Files':  'var(--cat-config)',
    'OSINT / Username':     'var(--cat-osint)',
};

// Category Icons (SVG paths)
const CATEGORY_ICONS = {
    'File Exposure':        '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>',
    'Login/Admin Panels':   '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>',
    'Directory Listing':    '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>',
    'Error Messages':       '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>',
    'Sensitive Data':       '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>',
    'SEO/Indexing':         '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>',
    'Database Exposure':    '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/></svg>',
    'Configuration Files':  '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>',
    'OSINT / Username':     '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>',
};

// ── Create Category Filter Item ───────────────────────────────
function createCategoryFilter(category, count, isActive) {
    const color = CATEGORY_COLORS[category] || 'var(--text-muted)';
    const activeClass = isActive ? ' active' : '';

    return `
        <div class="category-filter${activeClass}" data-category="${escapeHtml(category)}" onclick="toggleCategory(this)">
            <div class="category-filter__checkbox"></div>
            <div class="category-filter__dot" style="background: ${color}; box-shadow: 0 0 6px ${color}"></div>
            <span class="category-filter__name">${escapeHtml(category)}</span>
            <span class="category-filter__count">${count}</span>
        </div>
    `;
}

// ── Create Category Section (collapsible result group) ────────
function createCategorySection(categoryResult, favorites, isLiveMode) {
    const { category, dorks } = categoryResult;
    const color = CATEGORY_COLORS[category] || 'var(--text-muted)';
    const totalResults = dorks.reduce((sum, d) => sum + d.results.length, 0);

    const dorksHtml = dorks.map(d => createDorkCard(d, favorites, isLiveMode)).join('');

    return `
        <div class="category-section" data-category="${escapeHtml(category)}">
            <div class="category-section__header" onclick="toggleSection(this.parentElement)">
                <div class="category-section__dot" style="color: ${color}; background: ${color}"></div>
                <span class="category-section__title">${escapeHtml(category)}</span>
                <span class="category-section__count">${dorks.length} dorks · ${totalResults} results</span>
                <span class="category-section__toggle">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="6 9 12 15 18 9"/></svg>
                </span>
            </div>
            <div class="category-section__body">
                ${dorksHtml}
            </div>
        </div>
    `;
}

// ── Create Dork Result Card ───────────────────────────────────
function createDorkCard(dorkResult, favorites, isLiveMode) {
    const { dork, query, results } = dorkResult;
    const severityClass = `badge--severity-${dork.severity}`;
    const isFav = favorites && favorites.includes(dork.id);
    const favClass = isFav ? 'favorited' : '';
    const favFill = isFav ? 'currentColor' : 'none';

    let emptyMessage = '<div class="no-results">No results found for this query</div>';
    if (!isLiveMode) {
        emptyMessage = '<div class="no-results">No results fetched (Browser Mode). Click "Open in Browser" to view real results.</div>';
    }

    const resultsHtml = results.length > 0
        ? results.map(r => createResultItem(r)).join('')
        : emptyMessage;

    return `
        <div class="dork-card" data-dork-id="${escapeHtml(dork.id)}">
            <div class="dork-card__header">
                <span class="dork-card__name">${escapeHtml(dork.name)}</span>
                <span class="dork-card__type badge--type-${getDorkType(dork)}">${getDorkType(dork).toUpperCase()}</span>
                <span class="dork-card__severity ${severityClass}">${escapeHtml(dork.severity)}</span>
            </div>
            <div class="dork-card__description">${escapeHtml(dork.description)}</div>
            <div class="dork-card__query">
                <code>${escapeHtml(query)}</code>
            </div>
            <div class="dork-card__results">
                ${resultsHtml}
            </div>
            <div class="dork-card__actions">
                <button class="btn btn--small btn--ghost" onclick="openDorkInBrowser('${escapeJs(query)}')" title="Open in Google">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/><polyline points="15 3 21 3 21 9"/><line x1="10" y1="14" x2="21" y2="3"/></svg>
                    <span>Open in Browser</span>
                </button>
                <button class="btn btn--icon btn--ghost ${favClass}" onclick="toggleFav('${escapeJs(dork.id)}', this)" title="Toggle Favorite">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="${favFill}" stroke="currentColor" stroke-width="2"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>
                </button>
            </div>
        </div>
    `;
}

// ── Create Result Item ────────────────────────────────────────
function createResultItem(result) {
    return `
        <div class="result-item" onclick="openUrl('${escapeJs(result.url)}')">
            <div class="result-item__title">${escapeHtml(result.title)}</div>
            <div class="result-item__url">${escapeHtml(result.url)}</div>
            <div class="result-item__snippet">${escapeHtml(result.snippet)}</div>
        </div>
    `;
}

// ── Create Results Summary Banner ─────────────────────────────
function createResultsSummary(domain, keyword, totalDorks, totalResults, modeLabel, isLiveMode) {
    let targetText = '';
    if (domain && keyword) {
        targetText = `for <strong>${escapeHtml(domain)}</strong> with username "<strong>${escapeHtml(keyword)}</strong>"`;
    } else if (domain) {
        targetText = `for <strong>${escapeHtml(domain)}</strong>`;
    } else if (keyword) {
        targetText = `for username "<strong>${escapeHtml(keyword)}</strong>"`;
    }

    const modeBadge = modeLabel
        ? `<span class="results-summary__mode">${escapeHtml(modeLabel)}</span>`
        : '';

    let resultsText = `Found <strong>${totalResults}</strong> results fetched across <strong>${totalDorks}</strong> dorks ${targetText}`;
    if (!isLiveMode) {
        resultsText = `Generated <strong>${totalDorks}</strong> dorks ${targetText} — 0 results fetched (Browser Mode)`;
    }

    return `
        <div class="results-summary">
            <div class="results-summary__icon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
            </div>
            <div class="results-summary__text">
                <div class="results-summary__title">Scan Complete ${modeBadge}</div>
                <div class="results-summary__detail">
                    ${resultsText}
                </div>
            </div>
        </div>
    `;
}

// ── Create Toast Notification ─────────────────────────────────
function showToast(message, type = 'info') {
    const iconMap = {
        success: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>',
        error:   '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>',
        info:    '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>',
        warning: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>',
    };

    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `toast toast--${type}`;
    toast.innerHTML = `
        <span class="toast__icon">${iconMap[type] || iconMap.info}</span>
        <span>${escapeHtml(message)}</span>
    `;
    container.appendChild(toast);

    // Remove after animation
    setTimeout(() => {
        if (toast.parentNode) toast.parentNode.removeChild(toast);
    }, 3200);
}

// ── Utility: HTML Escape ──────────────────────────────────────
function escapeHtml(text) {
    if (text === null || text === undefined) return '';
    const div = document.createElement('div');
    div.textContent = String(text);
    return div.innerHTML;
}

// ── Utility: JS String Escape ─────────────────────────────────
function escapeJs(str) {
    if (!str) return '';
    return String(str)
        .replace(/\\/g, '\\\\')
        .replace(/'/g, "\\'")
        .replace(/"/g, '\\"')
        .replace(/\n/g, '\\n');
}

// ── Toggle Section Collapse ───────────────────────────────────
function toggleSection(sectionEl) {
    sectionEl.classList.toggle('collapsed');
}

// ── Dork Type Helper ──────────────────────────────────────────
function getDorkType(dork) {
    if (dork.needsDomain && dork.needsKeyword) return 'hybrid';
    if (dork.needsKeyword) return 'keyword';
    return 'domain';
}
