# CTI DORK — OSINT Google Dorking Desktop Tool

CTI DORK is a desktop application I built to automate Google dorking. The main goal of the project is to make manual OSINT searches faster and more structured.

Instead of writing dorks one by one, the app generates queries automatically, runs them, and collects the results in a single interface.

---

## Features

- Domain-based Google dorking
- Username-based OSINT search
- Automatic query (dork) generation
- Concurrent search execution (goroutines)
- Categorized results
- JSON and TXT export
- Simple desktop UI (Wails)

---

## How it works

The application works in 5 main steps:

1. **Query Generation**
   - Uses predefined dork templates
   - Generates username variations (underscore, hyphen, etc.)

2. **Query Execution**
   - Runs all queries concurrently
   - Uses Serper.dev API to fetch real Google results

3. **Validation**
   - Filters invalid URLs
   - Removes empty or broken results

4. **Deduplication**
   - Prevents duplicate links across queries

5. **Response Build**
   - Groups results by category and sends them to the frontend

---

## Search Modes

### 1. Domain Mode

**Input:**
```
example.com
```

**Example queries:**
- `site:example.com filetype:pdf`
- `site:example.com inurl:admin`
- `site:example.com "password"`

---

### 2. Username Mode

**Input:**
```
johndoe
```

**Process:**
- Generates username variations
- Searches across multiple platforms (GitHub, Instagram, Reddit, etc.)
- Runs generic dorks

---

### 3. Combined Mode

**Input:**
```
example.com + johndoe
```

- Runs both domain and username queries
- Adds intersection-based queries

---

## Tech Stack

- **Go** → backend
- **Wails** → desktop app framework
- **JavaScript / HTML / CSS** → frontend
- **Serper.dev API** → Google search

---

## Project Structure

```
.
├── app.go
├── main.go
├── internal/
│   ├── search/
│   ├── dorks/
│   └── storage/
├── frontend/
└── build/
```

---

## Installation

### Requirements

- Go 1.24+
- Wails

### Run

```bash
go mod tidy
wails dev
```

### Build

```bash
wails build
```

---

## Configuration

The app uses Serper.dev API.

You can set your API key:
- from the UI
- or via a config file

---

## Notes

- Without an API key → only generates queries
- With an API key → fetches real results
- All data comes from real Google searches

---

## Purpose

I built this project to:

- Improve OSINT workflows
- Practice Go and Wails

---

## Disclaimer

This tool is intended for educational and research purposes only. The user is responsible for how it is used.