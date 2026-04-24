# CTI DORK — OSINT Google Dorking Desktop Tool

CTI DORK is a desktop application designed to automate Google dorking workflows. The goal of this project is to reduce manual effort in OSINT processes by generating, executing, and organizing search queries in a structured way.

Instead of manually writing and testing dorks, the application builds queries dynamically based on user input (domain and/or username), executes them concurrently, and aggregates the results.

---

## Features

- Domain-based Google dorking (file exposure, admin panels, indexing issues, etc.)
- Username-based OSINT search across multiple platforms
- Automatic query generation using predefined templates
- Username variation engine (hyphen, underscore, spacing)
- Concurrent query execution using goroutines
- Result validation and filtering
- Global deduplication across all queries
- Categorized output
- Export options (JSON, TXT)
- Desktop UI built with Wails

---

## Architecture Overview

The application is structured into three main layers:

### 1. Frontend
- Built with HTML, CSS, and JavaScript
- Runs inside Wails WebView
- Handles UI rendering and user interaction
- Communicates with backend via Wails bindings

### 2. Backend (Go)
- Core application logic
- Query generation and execution
- Result processing (validation, deduplication)
- API communication (Serper.dev)

### 3. Internal Packages

```
internal/
├── search/     # Query execution, API client, modes
├── dorks/      # Dork templates and query builders
└── storage/    # Export and persistence logic
```

---

## Search Pipeline

The search process is divided into multiple stages:

### 1. Query Collection

Queries are generated from multiple sources:

- Predefined dork templates (domain-based)
- Platform-specific queries (username-based)
- Generic queries (inurl, intitle, quoted)
- Intersection queries (domain + username)

---

### 2. Query Execution

- Queries are executed concurrently using goroutines
- A worker pool handles parallel execution
- Each query is sent to Serper.dev API
- API responses contain real Google search results

---

### 3. Validation

Each result is validated before processing:

- URL must start with `http://` or `https://`
- Title must not be empty
- Invalid or malformed results are discarded

---

### 4. Deduplication

- Results are normalized (case-insensitive, stripped separators)
- Duplicate URLs across different queries are removed
- First occurrence of a result is preserved

---

### 5. Response Construction

- Results are grouped by category
- Structured response is built
- Data is sent to frontend via Wails bindings

---

## Search Modes

### Domain Mode

**Input:**
```
example.com
```

**Purpose:**
Identify potential exposures on a target domain.

**Example queries:**
- `site:example.com filetype:pdf`
- `site:example.com inurl:admin`
- `site:example.com intitle:"index of"`
- `site:example.com "password"`

---

### Username Mode

**Input:**
```
johndoe
```

**Process:**
- Generate variations:
  - `johndoe`
  - `john-doe`
  - `john_doe`
  - `"john doe"`
- Run platform-specific queries:
  - Facebook
  - Instagram
  - Twitter / X
  - YouTube
  - TikTok
  - Spotify
  - Apple Music
  - SoundCloud
  - Telegram
  - Snapchat
  - GitHub
  - LinkedIn
  - Reddit
  - Pastebin
  - GitLab
  - NSosyal (Türkiye's Social Media Platform)
- Run generic queries:
  - `inurl:johndoe`
  - `intitle:johndoe`
  - `"johndoe"`

---

### Combined Mode

**Input:**
```
example.com + johndoe
```

**Behavior:**
- Executes all domain-based queries
- Executes all username-based queries
- Adds intersection queries:
  - `site:example.com "johndoe"`
  - `site:example.com inurl:johndoe`

---

## Concurrency Model

- Fixed number of worker goroutines
- Shared job queue (queries)
- Each worker:
  - Pulls query from queue
  - Sends request to API
  - Returns results to main routine
- Main routine aggregates results

This approach reduces total execution time significantly compared to sequential execution.

---

## Tech Stack

- **Go** — core backend logic
- **Wails** — desktop application framework
- **JavaScript / HTML / CSS** — frontend
- **Serper.dev API** — Google search provider

---

## Project Structure

```
.
├── app.go                # Main application logic and bindings
├── main.go               # Entry point
├── go.mod
├── wails.json
│
├── internal/
│   ├── search/           # Search engine, API client, modes
│   ├── dorks/            # Dork templates and builders
│   └── storage/          # Export logic
│
├── frontend/             # UI files
└── build/                # Compiled output
```

---

## Installation

### Requirements

- Go 1.24+
- Wails CLI

### Run (development)

```bash
go mod tidy
wails dev
```

### Build

```bash
wails build
```

Output binary will be located in:

```
build/bin/
```

---

## Configuration

The application uses Serper.dev API to fetch search results.

### API Key Setup

You can configure the API key in two ways:

1. From the application UI
2. Via a configuration file (e.g. local config path)

---

## Data Handling

- No mock or generated data is used
- All results come from real search queries
- Invalid data is filtered out before display
- Duplicate results are removed globally

---

## Export

### JSON Export
- Structured output
- Suitable for further processing

### TXT Export
- Human-readable format
- Suitable for reporting

---

## Limitations

- Dependent on external API (Serper.dev)
- API result limits apply (e.g. max results per query)
- Some results may be filtered due to validation rules

---

## Purpose

This project was developed to:

- Automate repetitive OSINT tasks
- Improve efficiency in Google dorking workflows
- Practice backend development with Go
- Build a desktop application using Wails

---

## Disclaimer

This tool is intended for educational and research purposes only. Users are responsible for how they use it.