package dorks

import "strings"

// dork query category
type Category string

const (
	CategoryFileExposure Category = "File Exposure"
	CategoryLoginPanels  Category = "Login/Admin Panels"
	CategoryDirectories  Category = "Directory Listing"
	CategoryErrors       Category = "Error Messages"
	CategorySensitive    Category = "Sensitive Data"
	CategorySEO          Category = "SEO/Indexing"
	CategoryDatabase     Category = "Database Exposure"
	CategoryConfig       Category = "Configuration Files"
	CategoryOSINT        Category = "OSINT / Username"
)

//  single Google dork 
type Dork struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Template     string   `json:"template"`
	Category     Category `json:"category"`
	Description  string   `json:"description"`
	Severity     string   `json:"severity"`     // "low", "medium", "high", "critical"
	NeedsDomain  bool     `json:"needsDomain"`  // true if this dork requires a domain
	NeedsKeyword bool     `json:"needsKeyword"` // true if this dork requires a keyword
}

//  "domain", "keyword", or "hybrid" 
func (d Dork) DorkType() string {
	if d.NeedsDomain && d.NeedsKeyword {
		return "hybrid"
	}
	if d.NeedsKeyword {
		return "keyword"
	}
	return "domain"
}

func (d Dork) IsCompatible(domain, keyword string) bool {
	if d.NeedsDomain && domain == "" {
		return false
	}
	if d.NeedsKeyword && keyword == "" {
		return false
	}

	if strings.Contains(d.Template, "{domain}") && domain == "" {
		return false
	}
	if strings.Contains(d.Template, "{keyword}") && keyword == "" {
		return false
	}
	return true
}

// GetCategories 
func GetCategories() []string {
	return []string{
		string(CategoryFileExposure),
		string(CategoryLoginPanels),
		string(CategoryDirectories),
		string(CategoryErrors),
		string(CategorySensitive),
		string(CategorySEO),
		string(CategoryDatabase),
		string(CategoryConfig),
		string(CategoryOSINT),
	}
}

// GetAllDorks returns 
func GetAllDorks() []Dork {
	return []Dork{
		
		// FILE EXPOSURE 
		
		{
			ID: "file-01", Name: "PDF Documents",
			Template: "site:{domain} filetype:pdf", Category: CategoryFileExposure,
			Description: "Find publicly accessible PDF files on the target domain",
			Severity: "low", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "file-02", Name: "Excel Spreadsheets",
			Template: "site:{domain} filetype:xls OR filetype:xlsx", Category: CategoryFileExposure,
			Description: "Discover Excel files that may contain sensitive data",
			Severity: "medium", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "file-03", Name: "Word Documents",
			Template: "site:{domain} filetype:doc OR filetype:docx", Category: CategoryFileExposure,
			Description: "Find Word documents with potential internal information",
			Severity: "low", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "file-04", Name: "SQL Dump Files",
			Template: "site:{domain} filetype:sql", Category: CategoryFileExposure,
			Description: "Locate exposed SQL database dump files",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "file-05", Name: "Log Files",
			Template: "site:{domain} filetype:log", Category: CategoryFileExposure,
			Description: "Find application and server log files",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "file-06", Name: "Backup Files",
			Template: "site:{domain} filetype:bak OR filetype:backup OR filetype:old", Category: CategoryFileExposure,
			Description: "Discover backup files that may contain sensitive data",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "file-07", Name: "Environment Files",
			Template: "site:{domain} filetype:env", Category: CategoryFileExposure,
			Description: "Find .env files with credentials and secrets",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "file-08", Name: "CSV Data Files",
			Template: "site:{domain} filetype:csv", Category: CategoryFileExposure,
			Description: "Locate CSV files that may contain exported data",
			Severity: "medium", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "file-09", Name: "XML Data Files",
			Template: "site:{domain} filetype:xml", Category: CategoryFileExposure,
			Description: "Find XML configuration or data files",
			Severity: "medium", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "file-10", Name: "Configuration Files",
			Template: "site:{domain} filetype:conf OR filetype:cfg", Category: CategoryFileExposure,
			Description: "Discover server and application configuration files",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},

		// LOGIN / ADMIN PANELS 
		{
			ID: "login-01", Name: "Admin URL Paths",
			Template: "site:{domain} inurl:admin", Category: CategoryLoginPanels,
			Description: "Find pages with 'admin' in the URL path",
			Severity: "medium", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "login-02", Name: "Login Pages",
			Template: "site:{domain} inurl:login", Category: CategoryLoginPanels,
			Description: "Discover login entry points on the domain",
			Severity: "medium", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "login-03", Name: "Admin Panel Titles",
			Template: `site:{domain} intitle:"admin panel"`, Category: CategoryLoginPanels,
			Description: "Find pages titled 'admin panel'",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "login-04", Name: "WordPress Admin",
			Template: "site:{domain} inurl:wp-admin OR inurl:wp-login", Category: CategoryLoginPanels,
			Description: "Locate WordPress administration panels",
			Severity: "medium", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "login-05", Name: "cPanel Access",
			Template: "site:{domain} inurl:cpanel OR inurl:webmail", Category: CategoryLoginPanels,
			Description: "Find cPanel or webmail login interfaces",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "login-06", Name: "Dashboard Pages",
			Template: `site:{domain} intitle:"dashboard" inurl:admin`, Category: CategoryLoginPanels,
			Description: "Discover admin dashboard pages",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},

		// DIRECTORY LISTING 

		{
			ID: "dir-01", Name: "Open Directory Index",
			Template: `site:{domain} intitle:"index of /"`, Category: CategoryDirectories,
			Description: "Find open directory listings on the server",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "dir-02", Name: "Parent Directory Listing",
			Template: `site:{domain} intitle:"index of" "parent directory"`, Category: CategoryDirectories,
			Description: "Discover directory listings with parent navigation",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "dir-03", Name: "FTP Directory Listing",
			Template: `site:{domain} intitle:"index of" inurl:ftp`, Category: CategoryDirectories,
			Description: "Find exposed FTP directory indices",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "dir-04", Name: "Backup Directories",
			Template: `site:{domain} intitle:"index of" "backup"`, Category: CategoryDirectories,
			Description: "Locate directory listings containing backup files",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "dir-05", Name: "Upload Directories",
			Template: `site:{domain} intitle:"index of" "upload"`, Category: CategoryDirectories,
			Description: "Find exposed upload directories",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},

		// ERROR MESSAGES

		{
			ID: "err-01", Name: "MySQL Syntax Errors",
			Template: `site:{domain} "sql syntax" "mysql"`, Category: CategoryErrors,
			Description: "Find pages exposing MySQL syntax error messages",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "err-02", Name: "PHP MySQL Warnings",
			Template: `site:{domain} "Warning:" "mysql_" filetype:php`, Category: CategoryErrors,
			Description: "Discover PHP pages with MySQL warning messages",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "err-03", Name: "PHP Fatal Errors",
			Template: `site:{domain} "Fatal error" filetype:php`, Category: CategoryErrors,
			Description: "Find PHP fatal error pages revealing server paths",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "err-04", Name: "Server Stack Traces",
			Template: `site:{domain} "Stack Trace" "Server Error"`, Category: CategoryErrors,
			Description: "Locate pages with exposed server stack traces",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "err-05", Name: "Debug Error Pages",
			Template: `site:{domain} intext:"error occurred" "debug"`, Category: CategoryErrors,
			Description: "Find debug-mode error pages with detailed info",
			Severity: "medium", NeedsDomain: true, NeedsKeyword: false,
		},

	
		// SENSITIVE DATA

		{
			ID: "sens-01", Name: "Credentials in Logs",
			Template: `site:{domain} intext:"username" "password" filetype:log`, Category: CategorySensitive,
			Description: "Find log files containing username/password combos",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "sens-02", Name: "Private RSA Keys",
			Template: `site:{domain} "BEGIN RSA PRIVATE KEY"`, Category: CategorySensitive,
			Description: "Discover exposed RSA private key files",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "sens-03", Name: "Credential URLs",
			Template: "site:{domain} inurl:credentials", Category: CategorySensitive,
			Description: "Find URLs containing credential-related paths",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "sens-04", Name: "SSL/TLS Key Files",
			Template: "site:{domain} filetype:pem OR filetype:key", Category: CategorySensitive,
			Description: "Locate exposed SSL certificates and key files",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "sens-05", Name: "API Keys in Source",
			Template: `site:{domain} intext:"api_key" OR intext:"apikey"`, Category: CategorySensitive,
			Description: "Find pages exposing API keys in source code",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "sens-06", Name: "Secret/Access Keys",
			Template: `site:{domain} intext:"secret_key" OR intext:"access_key"`, Category: CategorySensitive,
			Description: "Discover exposed secret and access key values",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "sens-07", Name: "Password Files",
			Template: `site:{domain} inurl:password.txt OR inurl:passwords.txt`, Category: CategorySensitive,
			Description: "Find exposed password text files",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},


		// SEO / INDEXING

		{
			ID: "seo-01", Name: "Full Site Index",
			Template: "site:{domain}", Category: CategorySEO,
			Description: "View all indexed pages for the target domain",
			Severity: "low", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "seo-02", Name: "Non-WWW Pages",
			Template: "site:{domain} -www", Category: CategorySEO,
			Description: "Find indexed pages outside the www subdomain",
			Severity: "low", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "seo-03", Name: "Sitemap Files",
			Template: "site:{domain} inurl:sitemap.xml", Category: CategorySEO,
			Description: "Locate XML sitemap files for the domain",
			Severity: "low", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "seo-04", Name: "XML Sitemaps",
			Template: "site:{domain} filetype:xml inurl:sitemap", Category: CategorySEO,
			Description: "Find all XML sitemap variations",
			Severity: "low", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "seo-05", Name: "Keyword on Domain",
			Template: "site:{domain} intext:{keyword}", Category: CategorySEO,
			Description: "Search for keyword mentions across the domain",
			Severity: "low", NeedsDomain: true, NeedsKeyword: true,
		},


		// DATABASE EXPOSURE 

		{
			ID: "db-01", Name: "SQL INSERT Dumps",
			Template: `site:{domain} filetype:sql "INSERT INTO"`, Category: CategoryDatabase,
			Description: "Find SQL dump files containing INSERT statements",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "db-02", Name: "SQL CREATE Dumps",
			Template: `site:{domain} filetype:sql "CREATE TABLE"`, Category: CategoryDatabase,
			Description: "Discover SQL schema dump files",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "db-03", Name: "phpMyAdmin Panels",
			Template: "site:{domain} inurl:phpmyadmin", Category: CategoryDatabase,
			Description: "Locate exposed phpMyAdmin interfaces",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "db-04", Name: "phpMyAdmin Welcome",
			Template: `site:{domain} intitle:"phpMyAdmin" "Welcome to"`, Category: CategoryDatabase,
			Description: "Find phpMyAdmin welcome/login pages",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "db-05", Name: "Access DB Files",
			Template: "site:{domain} filetype:mdb OR filetype:accdb", Category: CategoryDatabase,
			Description: "Discover exposed Microsoft Access database files",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},

	
		// CONFIGURATION FILES 
	
		{
			ID: "cfg-01", Name: "YAML Config Files",
			Template: "site:{domain} filetype:yml OR filetype:yaml", Category: CategoryConfig,
			Description: "Find YAML configuration files",
			Severity: "medium", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "cfg-02", Name: "JSON with Passwords",
			Template: `site:{domain} filetype:json "password"`, Category: CategoryConfig,
			Description: "Locate JSON files containing password fields",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "cfg-03", Name: "INI Database Config",
			Template: `site:{domain} filetype:ini "[database]"`, Category: CategoryConfig,
			Description: "Find INI files with database configuration sections",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "cfg-04", Name: "Exposed Git Repos",
			Template: "site:{domain} inurl:.git", Category: CategoryConfig,
			Description: "Discover exposed .git directories with source code",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "cfg-05", Name: "Docker Compose Files",
			Template: `site:{domain} filetype:yml "docker-compose"`, Category: CategoryConfig,
			Description: "Find Docker Compose configuration files",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "cfg-06", Name: "htaccess Files",
			Template: `site:{domain} filetype:htaccess`, Category: CategoryConfig,
			Description: "Locate exposed Apache .htaccess configuration files",
			Severity: "high", NeedsDomain: true, NeedsKeyword: false,
		},
		{
			ID: "cfg-07", Name: "Web.config Files",
			Template: `site:{domain} filetype:config "connectionString"`, Category: CategoryConfig,
			Description: "Find ASP.NET web.config files with connection strings",
			Severity: "critical", NeedsDomain: true, NeedsKeyword: false,
		},


		// OSINT / USERNAME 

		{
			ID: "osint-01", Name: "Username Search",
			Template: `"{keyword}" site:twitter.com OR site:x.com`, Category: CategoryOSINT,
			Description: "Search for a username or keyword on Twitter/X",
			Severity: "low", NeedsDomain: false, NeedsKeyword: true,
		},
		{
			ID: "osint-02", Name: "GitHub Profiles",
			Template: `"{keyword}" site:github.com`, Category: CategoryOSINT,
			Description: "Find GitHub accounts or repos matching the keyword",
			Severity: "low", NeedsDomain: false, NeedsKeyword: true,
		},
		{
			ID: "osint-03", Name: "LinkedIn Profiles",
			Template: `"{keyword}" site:linkedin.com/in`, Category: CategoryOSINT,
			Description: "Locate LinkedIn profiles matching the keyword",
			Severity: "low", NeedsDomain: false, NeedsKeyword: true,
		},
		{
			ID: "osint-04", Name: "Pastebin Leaks",
			Template: `"{keyword}" site:pastebin.com`, Category: CategoryOSINT,
			Description: "Search for keyword mentions on Pastebin paste dumps",
			Severity: "high", NeedsDomain: false, NeedsKeyword: true,
		},
		{
			ID: "osint-05", Name: "Reddit Mentions",
			Template: `"{keyword}" site:reddit.com`, Category: CategoryOSINT,
			Description: "Find Reddit posts and comments mentioning the keyword",
			Severity: "low", NeedsDomain: false, NeedsKeyword: true,
		},
		{
			ID: "osint-06", Name: "Exposed Credentials",
			Template: `"{keyword}" intext:password OR intext:passwd`, Category: CategoryOSINT,
			Description: "Search for password leaks containing the keyword",
			Severity: "critical", NeedsDomain: false, NeedsKeyword: true,
		},
		{
			ID: "osint-07", Name: "Email Address Leaks",
			Template: `"{keyword}" filetype:txt OR filetype:log "email"`, Category: CategoryOSINT,
			Description: "Find text or log files leaking email addresses with keyword",
			Severity: "high", NeedsDomain: false, NeedsKeyword: true,
		},
		{
			ID: "osint-08", Name: "Public Documents",
			Template: `"{keyword}" filetype:pdf OR filetype:docx`, Category: CategoryOSINT,
			Description: "Discover PDF/Word documents mentioning the keyword",
			Severity: "medium", NeedsDomain: false, NeedsKeyword: true,
		},
		{
			ID: "osint-09", Name: "Forum Posts",
			Template: `"{keyword}" inurl:forum OR inurl:thread OR inurl:topic`, Category: CategoryOSINT,
			Description: "Find forum discussions mentioning the keyword",
			Severity: "low", NeedsDomain: false, NeedsKeyword: true,
		},
		{
			ID: "osint-10", Name: "Code Leaks",
			Template: `"{keyword}" site:gist.github.com OR site:gitlab.com`, Category: CategoryOSINT,
			Description: "Search for code snippets and gists matching the keyword",
			Severity: "medium", NeedsDomain: false, NeedsKeyword: true,
		},
	}
}


func GetDorksByCategory(category string) []Dork {
	var result []Dork
	for _, d := range GetAllDorks() {
		if string(d.Category) == category {
			result = append(result, d)
		}
	}
	return result
}

func GetCompatibleDorks(domain, keyword string) []Dork {
	var result []Dork
	for _, d := range GetAllDorks() {
		if d.IsCompatible(domain, keyword) {
			result = append(result, d)
		}
	}
	return result
}

func GetDorksByMode(mode string, domain, keyword string) []Dork {
	all := GetAllDorks()
	var result []Dork

	for _, d := range all {
		switch mode {
		case "WEBSITE_ONLY":
			if d.NeedsKeyword {
				continue
			}
			if !d.IsCompatible(domain, "") {
				continue
			}
		case "USERNAME_ONLY":
			if d.NeedsDomain {
				continue
			}
			if !d.IsCompatible("", keyword) {
				continue
			}
		case "INTERSECTION":
			
			if !d.IsCompatible(domain, keyword) {
				continue
			}
		default:
			if !d.IsCompatible(domain, keyword) {
				continue
			}
		}
		result = append(result, d)
	}

	return result
}
