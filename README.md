# CTI DORK — OSINT Google Dorking Desktop Tool

<p align="center">
  <strong>An intelligent OSINT discovery engine that automates Google dorking for security researchers, penetration testers, and CTI analysts.</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go" alt="Go 1.24">
  <img src="https://img.shields.io/badge/Wails-2.12-purple?style=flat" alt="Wails 2.12">
  <img src="https://img.shields.io/badge/Platform-Windows-0078D6?style=flat&logo=windows" alt="Windows">
  <img src="https://img.shields.io/badge/API-Serper.dev-orange?style=flat" alt="Serper.dev">
  <img src="https://img.shields.io/badge/License-Private-red?style=flat" alt="Private">
</p>

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Search Modes](#search-modes)
- [OSINT Platform Coverage](#osint-platform-coverage)
- [Query Expansion Engine](#query-expansion-engine)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Technical Details](#technical-details)
- [Export Options](#export-options)
- [Security & Data Integrity](#security--data-integrity)
- [Development](#development)
- [FAQ](#faq)

---

## Overview

**CTI DORK** is a Windows desktop application built with [Wails](https://wails.io/) (Go + HTML/CSS/JS) that transforms Google dorking from a manual, error-prone process into an automated, intelligent OSINT workflow.

Instead of manually crafting Google dork queries one by one, CTI DORK:

1. **Generates** 45+ categorized dork queries from a built-in template library
2. **Expands** usernames into 50+ platform-specific and variation queries automatically
3. **Executes** all queries concurrently via the Serper.dev API (real Google results)
4. **Validates** and deduplicates results with a strict data integrity pipeline
5. **Presents** results in a premium dark-themed UI grouped by category

> **No fake data. No mock results. No placeholders.** Every result comes directly from Google via the Serper.dev API.

---

## Features

### 🔍 Intelligent Search Engine
- **3 Search Modes**: Website-only (A), Username-only (B), Intersection (A ∩ B)
- **Automatic query expansion** — generates 50+ queries from a single username
- **16-platform OSINT coverage** — Facebook, Instagram, Twitter/X, GitHub, YouTube, TikTok, Telegram, and more
- **Smart fallback system** — automatically runs broader queries when results are low
- **Normalized fuzzy matching** — `bulentdurusoy` matches `bulent-durusoy`, `bulent_durusoy`

### ⚡ Performance
- **8 concurrent goroutine workers** — searches complete in 5-8 seconds
- **Global deduplication** — no duplicate URLs across queries
- **Priority-based execution** — platform queries run before generic ones

### 🛡️ Data Integrity
- **Zero fabricated data** — every result is from a real Google search
- **Strict validation** — URLs must be absolute (`http(s)://`), titles must be non-empty
- **Browser mode fallback** — if no API key, shows queries with "Open in Browser" buttons

### 🎨 Premium UI
- Dark glassmorphism design with smooth animations
- Category-based filtering with color-coded sections
- Favorites system with persistent storage
- Collapsible result sections
- One-click "Open in Browser" for any query
- Settings modal for API key management

### 📦 Export
- **JSON export** — structured data for further analysis
- **TXT export** — formatted human-readable reports

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    CTI DORK Application                      │
├─────────────┬───────────────────────────────────────────────┤
│   Frontend  │  HTML/CSS/JS (Wails WebView)                  │
│             │  ├── index.html     — Main layout + modals    │
│             │  ├── app.js         — State + orchestration   │
│             │  ├── components.js  — UI component builders   │
│             │  └── styles.css     — Dark theme + animations │
├─────────────┼───────────────────────────────────────────────┤
│   Backend   │  Go (Wails Bindings)                          │
│             │                                               │
│   app.go    │  ┌─ ExecuteSearch() ─────────────────────┐    │
│             │  │  Phase 1: Collect queries              │    │
│             │  │    ├── Hardcoded dorks (45)             │    │
│             │  │    ├── Platform queries (~29)           │    │
│             │  │    ├── Generic queries (~6)             │    │
│             │  │    └── Intersection queries             │    │
│             │  │  Phase 2: Execute (8 workers)           │    │
│             │  │  Phase 3: Deduplicate + sort            │    │
│             │  │  Phase 4: Fallback (if results < 5)     │    │
│             │  │  Phase 5: Build response                │    │
│             │  └────────────────────────────────────────┘    │
├─────────────┼───────────────────────────────────────────────┤
│  Packages   │                                               │
│             │  internal/search/                              │
│             │  ├── engine.go      — SearchEngine interface   │
│             │  ├── serper.go      — Serper.dev API client    │
│             │  ├── livesearch.go  — Smart engine router      │
│             │  ├── expander.go    — Query expansion + OSINT  │
│             │  ├── modes.go       — Modes + validation       │
│             │  ├── config.go      — API key management       │
│             │  └── browser.go     — OS browser launcher      │
│             │                                               │
│             │  internal/dorks/                               │
│             │  ├── dorks.go       — 45 dork templates        │
│             │  └── builder.go     — Query builder            │
│             │                                               │
│             │  internal/storage/                             │
│             │  ├── favorites.go   — Favorites persistence    │
│             │  └── export.go      — JSON/TXT exporters       │
└─────────────┴───────────────────────────────────────────────┘
```

---

## Search Modes

### Mode A — Website Only

**Input:** Only a website/domain is provided (e.g., `example.com`)

Runs domain-focused dorks across 8 categories:

| Category | Examples |
|----------|----------|
| File Exposure | `site:example.com filetype:pdf`, `site:example.com filetype:sql` |
| Login/Admin Panels | `site:example.com inurl:admin`, `site:example.com intitle:"login"` |
| Directory Listing | `site:example.com intitle:"index of"` |
| Error Messages | `site:example.com "SQL syntax"`, `site:example.com "Fatal error"` |
| Sensitive Data | `site:example.com filetype:env`, `site:example.com "password"` |
| SEO/Indexing | `site:example.com -www`, `site:example.com inurl:sitemap.xml` |
| Database Exposure | `site:example.com filetype:sql "INSERT INTO"` |
| Configuration Files | `site:example.com filetype:xml "password"` |

### Mode B — Username Only

**Input:** Only a username is provided (e.g., `bulentdurusoy`)

Runs:
- **45+ hardcoded OSINT dorks** (Twitter, GitHub, LinkedIn, Pastebin, Reddit, etc.)
- **~29 platform-specific queries** across 16 platforms
- **~6 generic queries** (inurl, intitle, quoted)
- **Fallback queries** if results < 5

### Mode A ∩ B — Intersection

**Input:** Both domain and username provided (e.g., `aybu.edu.tr` + `21050111029`)

Runs:
- All domain-based dorks
- All username expansion queries
- **Intersection-specific queries**: `site:aybu.edu.tr "21050111029"`, `site:aybu.edu.tr inurl:21050111029`

---

## OSINT Platform Coverage

When a username is provided, CTI DORK automatically searches across these platforms:

| Platform | Queries | Type |
|----------|---------|------|
| 📘 Facebook | 2 | quoted + inurl |
| 📷 Instagram | 2 | quoted + inurl |
| 🐦 Twitter / X | 3 | x.com + twitter.com + inurl |
| 🎬 YouTube | 3 | quoted + @handle + inurl |
| 🎵 TikTok | 2 | quoted + inurl |
| 🎧 Spotify | 1 | quoted |
| 🍎 Apple Music | 1 | quoted |
| 🔊 SoundCloud | 2 | quoted + inurl |
| ✈️ Telegram | 2 | t.me + telegram.me |
| 👻 Snapchat | 2 | quoted + inurl |
| 🐙 GitHub | 3 | quoted + inurl + intitle |
| 💼 LinkedIn | 1 | quoted |
| 🤖 Reddit | 2 | quoted + inurl |
| 📋 Pastebin | 1 | quoted |
| 🦊 GitLab | 2 | quoted + inurl |
| 🇹🇷 NSosyal | 2 | quoted + inurl |

**Total: ~29 platform-specific queries per username**

---

## Query Expansion Engine

### Username Variations

For input `bulentdurusoy`, the engine automatically generates:

```
bulentdurusoy           → raw username
"bulentdurusoy"         → quoted (via platform templates)
bulent-durusoy          → hyphenated split
bulent_durusoy          → underscore split
bulent durusoy          → space split
"bulent durusoy"        → quoted space split
```

The split heuristic uses:
- Minimum part length (≥3 characters)
- Vowel/consonant transition scoring
- Turkish language support (ö, ü, ı vowels)

### Normalized Matching

Results are validated using normalized comparison:

```
Input:   bulentdurusoy
Matches: bulent-durusoy  ✅ (strips hyphens)
         bulent_durusoy  ✅ (strips underscores)
         bulent.durusoy  ✅ (strips dots)
         BULENTDURUSOY   ✅ (case-insensitive)
```

### Priority System

Queries execute in priority order:

| Priority | Type | Example |
|----------|------|---------|
| 1 | Platform-specific | `site:instagram.com "bulentdurusoy"` |
| 2 | inurl / intitle | `inurl:bulentdurusoy` |
| 3 | Generic variations | `"bulent durusoy"` |
| 4 | Fallback (broad) | `"bulentdurusoy" OR "bulent durusoy"` |

Higher-priority queries claim URLs first during deduplication.

---

## Installation

### Prerequisites

- **Go** 1.24+ — [download](https://go.dev/dl/)
- **Wails CLI** v2.12+ — install via `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **WebView2 Runtime** — pre-installed on Windows 10/11
- **GCC** (for CGo) — install via [MSYS2](https://www.msys2.org/) or [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)

### Build

```bash
# Clone or navigate to the project directory
cd "CTI DORK"

# Install Go dependencies
go mod tidy

# Development mode (hot reload)
wails dev

# Production build
wails build
```

The production binary is output to:
```
build/bin/CTI-DORK.exe
```

---

## Configuration

### API Key Setup

CTI DORK uses [Serper.dev](https://serper.dev) to execute real Google searches.

**Option 1: Via the app UI**
1. Launch CTI DORK
2. Click the **⚙️ Settings** button in the top bar
3. Enter your Serper.dev API key
4. Click **Save Configuration**

**Option 2: Manual config file**

Create or edit `~/.cti-dork/config.json`:

```json
{
  "serperApiKey": "YOUR_SERPER_API_KEY_HERE"
}
```

### Modes

| Mode | Condition | Behavior |
|------|-----------|----------|
| **Live Search Mode** | API key is set | Fetches real results from Google via Serper.dev |
| **Browser Mode** | No API key | Generates queries only; click "Open in Browser" to search manually |

> Serper.dev offers **2,500 free queries** on signup at [serper.dev](https://serper.dev).

---

## Usage

### Basic Search

1. **Website scan**: Enter a domain (e.g., `example.com`) → Click **Search**
   - Runs 35+ domain-focused security dorks
   
2. **Username OSINT**: Enter a username (e.g., `johndoe`) → Click **Search**
   - Runs 50+ queries across 16 platforms automatically

3. **Combined scan**: Enter both domain and username → Click **Search**
   - Runs all domain dorks + all username queries + intersection queries

### Filtering

- **Category filters** (sidebar): Show/hide results by category — does not block query execution
- **Favorites filter**: Show only favorited dorks
- **Critical filter**: Show only critical-severity dorks

### Exporting

- **JSON**: Click the JSON export button → saves structured data
- **TXT**: Click the TXT export button → saves formatted report

### Keyboard Shortcuts

- **Enter**: Execute search (when focused on either input field)

---

## Project Structure

```
CTI DORK/
├── app.go                          # Main app: search orchestration, API bindings
├── main.go                         # Wails app entry point
├── wails.json                      # Wails configuration
├── go.mod                          # Go module definition
│
├── internal/
│   ├── search/
│   │   ├── engine.go               # SearchEngine interface + SearchResult struct
│   │   ├── serper.go               # Serper.dev API client (real Google results)
│   │   ├── livesearch.go           # Smart engine router (API ↔ browser mode)
│   │   ├── expander.go             # OSINT query expansion + 16 platform registry
│   │   ├── modes.go                # Search modes, validation, fuzzy matching
│   │   ├── config.go               # API key config (~/.cti-dork/config.json)
│   │   ├── browser.go              # OS browser launcher
│   │   └── mockapi.go              # Debug-only mock engine (not used in production)
│   │
│   ├── dorks/
│   │   ├── dorks.go                # 45 categorized dork templates
│   │   └── builder.go              # Query builder ({domain}/{keyword} replacement)
│   │
│   └── storage/
│       ├── favorites.go            # Favorites persistence (~/.cti-dork/favorites.json)
│       └── export.go               # JSON/TXT export formatters
│
├── frontend/
│   ├── index.html                  # Main HTML layout + settings modal
│   └── assets/
│       ├── js/
│       │   ├── app.js              # Application state + search orchestration
│       │   └── components.js       # UI component builders (cards, sections, banners)
│       └── css/
│           └── styles.css          # Dark theme + glassmorphism + animations
│
└── build/
    ├── bin/
    │   └── CTI-DORK.exe            # Production binary (after `wails build`)
    └── windows/
        └── wails.exe.manifest      # Windows manifest
```

---

## Technical Details

### Search Pipeline

```
User Input
    │
    ▼
Phase 1: COLLECT QUERIES
    ├── Hardcoded dorks (from dorks.go)           → 10-45 queries
    ├── Platform queries (from expander.go)        → ~29 queries
    ├── Generic queries (inurl/intitle/quoted)      → ~6 queries
    └── Intersection queries (if both inputs)       → ~4 queries
    │
    ▼
Phase 2: EXECUTE CONCURRENTLY
    └── 8 goroutine workers → Serper.dev API
    │
    ▼
Phase 3: VALIDATE + DEDUPLICATE
    ├── URL format validation (must be absolute)
    ├── Title validation (must be non-empty)
    ├── Global URL deduplication (normalized)
    └── Priority-based URL claiming
    │
    ▼
Phase 4: FALLBACK (if results < 5 and username provided)
    ├── "username" profile OR account OR user
    ├── "username" OR "full name guess"
    └── "username" filetype:pdf OR doc OR xlsx
    │
    ▼
Phase 5: BUILD RESPONSE
    └── Group by category → JSON response → Frontend
```

### Concurrency Model

```
                    ┌──── Worker 1 → Serper API → results
                    ├──── Worker 2 → Serper API → results
 Job Channel ───────├──── Worker 3 → Serper API → results
 (~50 queries)      ├──── Worker 4 → Serper API → results
                    ├──── Worker 5 → Serper API → results
                    ├──── Worker 6 → Serper API → results
                    ├──── Worker 7 → Serper API → results
                    └──── Worker 8 → Serper API → results
                                        │
                                 Results Channel
                                        │
                              Main Goroutine (collect)
                                        │
                              Deduplicate + Sort
```

### Dork Categories (45 Templates)

| Category | Count | Severity Range |
|----------|-------|---------------|
| File Exposure | 5 | Medium – Critical |
| Login/Admin Panels | 5 | Medium – High |
| Directory Listing | 5 | Low – High |
| Error Messages | 5 | Low – Medium |
| Sensitive Data | 5 | High – Critical |
| SEO/Indexing | 3 | Low |
| Database Exposure | 5 | High – Critical |
| Configuration Files | 7 | Medium – Critical |
| OSINT / Username | 10 | Low – Critical |

---

## Export Options

### JSON Export

```json
{
  "domain": "example.com",
  "keyword": "username",
  "mode": "INTERSECTION",
  "categories": [
    {
      "category": "File Exposure",
      "dorks": [
        {
          "dork": {
            "id": "file-01",
            "name": "PDF Documents",
            "template": "site:{domain} filetype:pdf",
            "severity": "medium"
          },
          "query": "site:example.com filetype:pdf",
          "results": [
            {
              "title": "Annual Report 2024",
              "url": "https://example.com/docs/report.pdf",
              "snippet": "Company annual report..."
            }
          ]
        }
      ]
    }
  ]
}
```

### TXT Export

```
═══════════════════════════════════════════════════════════════
  CTI DORK — Google Dorking Search Report
═══════════════════════════════════════════════════════════════

  Domain:   example.com
  Keyword:  username
  Date:     2026-04-24 13:54:37

───────────────────────────────────────────────────────────────
  📁 File Exposure
───────────────────────────────────────────────────────────────

  ▸ PDF Documents [MEDIUM]
    Query: site:example.com filetype:pdf
    [1] Annual Report 2024
        https://example.com/docs/report.pdf
        Company annual report...

═══════════════════════════════════════════════════════════════
  Summary: 45 dorks executed, 12 results found
═══════════════════════════════════════════════════════════════
```

---

## Security & Data Integrity

### Strict Rules

| Rule | Implementation |
|------|---------------|
| **No fabricated data** | `LiveSearchEngine` returns only Serper.dev API results or empty `[]` |
| **No mock data in production** | `MockSearchEngine` exists but is never instantiated in `NewApp()` |
| **URL validation** | Every result URL must start with `http://` or `https://` |
| **Title validation** | Every result must have a non-empty title |
| **Global deduplication** | Normalized URLs are tracked across all queries |
| **API key security** | Key is stored locally at `~/.cti-dork/config.json`, masked in UI |

### Data Flow

```
Serper.dev API → Raw Results → URL Validation → Title Validation → Deduplication → Frontend
                                    ↓                                    ↓
                              Rejected results                  Duplicate URLs
                              are discarded                     are discarded
```

---

## Development

### Run in Development Mode

```bash
wails dev
```

This starts the app with hot reload — frontend changes reflect immediately, Go changes trigger a recompile.

### Build for Production

```bash
wails build
```

Outputs `build/bin/CTI-DORK.exe` — a standalone Windows executable.

### Key Files to Know

| File | Purpose |
|------|---------|
| `app.go` | All Wails-bound methods — the "API" between Go and JS |
| `internal/search/expander.go` | OSINT intelligence — platform registry + query expansion |
| `internal/search/serper.go` | Real Google search API client |
| `internal/search/modes.go` | Search modes, validation, fuzzy matching |
| `internal/dorks/dorks.go` | All 45 hardcoded dork templates |
| `frontend/assets/js/app.js` | Frontend state management + UI orchestration |

### Adding a New Platform

Edit `internal/search/expander.go` → `GetPlatforms()`:

```go
{
    Name:    "NewPlatform",
    Domains: []string{"newplatform.com"},
    Templates: []string{
        `site:newplatform.com "{u}"`,
        `site:newplatform.com inurl:{u}`,
    },
},
```

The platform will automatically be included in all username searches.

### Adding a New Dork

Edit `internal/dorks/dorks.go` → `GetAllDorks()`:

```go
{
    ID: "cat-XX", Name: "Dork Name",
    Template: `site:{domain} filetype:xyz "sensitive"`,
    Category: CategorySensitive,
    Description: "Description of what this dork finds",
    Severity: "high",
    NeedsDomain: true, NeedsKeyword: false,
},
```

---

## FAQ

### How many API credits does a single search use?

- **Website-only search**: ~35 API calls (one per dork)
- **Username-only search**: ~50 API calls (dorks + platform + generic)
- **Combined search**: ~55+ API calls

Serper.dev's free tier gives 2,500 queries.

### Why are some results missing even though Google shows them?

Possible reasons:
1. **API limit**: Serper.dev returns up to 10 results per query
2. **Deduplication**: A result found by an earlier query won't appear again
3. **Validation**: Results without valid URLs or titles are discarded

### Can I use a different search API?

Yes. Implement the `SearchEngine` interface in `internal/search/engine.go`:

```go
type SearchEngine interface {
    Search(query string) ([]SearchResult, error)
}
```

Then update `livesearch.go` to use your implementation.

### Is this tool legal?

CTI DORK uses public Google search results via a legitimate API (Serper.dev). However, how you use the results is your responsibility. Always ensure you have proper authorization before conducting security assessments.

---

<p align="center">
  <strong>Built with Go + Wails</strong><br>
  <sub>Powered by Serper.dev Google Search API</sub>
</p>
