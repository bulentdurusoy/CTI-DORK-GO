package search

import (
	"strings"
)

// MockSearchEngine provides simulated search results for in-app demonstration.
// It generates realistic-looking results based on the dork query pattern.
type MockSearchEngine struct{}

// NewMockSearchEngine creates a new MockSearchEngine instance
func NewMockSearchEngine() *MockSearchEngine {
	return &MockSearchEngine{}
}

// Search generates mock search results based on the query content
func (m *MockSearchEngine) Search(query string) ([]SearchResult, error) {
	domain := extractDomain(query)
	keyword := extractKeyword(query)

	// Determine mock result type from query content
	results := m.generateResults(query, domain, keyword)

	// Validate results strictly
	validResults := ValidateResults(results)

	return validResults, nil
}

// ValidateResults ensures all results meet strict criteria and removes duplicates
func ValidateResults(results []SearchResult) []SearchResult {
	var valid []SearchResult
	seen := make(map[string]bool)

	for _, r := range results {
		urlStr := strings.TrimSpace(r.URL)
		title := strings.TrimSpace(r.Title)

		// Title must not be empty
		if title == "" {
			continue
		}

		// URL must be non-empty and absolute
		if urlStr == "" || (!strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://")) {
			continue
		}

		// Remove duplicate URLs
		if seen[urlStr] {
			continue
		}
		seen[urlStr] = true

		// Assign cleaned values
		r.URL = urlStr
		r.Title = title
		valid = append(valid, r)
	}

	if valid == nil {
		return []SearchResult{}
	}
	return valid
}

// extractKeyword pulls the keyword from a quoted query segment
func extractKeyword(query string) string {
	// Try to extract from quotes like "keyword"
	start := strings.Index(query, `"`)
	if start != -1 {
		end := strings.Index(query[start+1:], `"`)
		if end != -1 {
			return query[start+1 : start+1+end]
		}
	}
	return ""
}

// extractDomain pulls the domain from a "site:X" query segment
func extractDomain(query string) string {
	parts := strings.Fields(query)
	for _, p := range parts {
		if strings.HasPrefix(p, "site:") {
			return strings.TrimPrefix(p, "site:")
		}
	}
	return ""
}

// generateResults creates category-aware mock results
func (m *MockSearchEngine) generateResults(query, domain, keyword string) []SearchResult {
	q := strings.ToLower(query)

	switch {
	case strings.Contains(q, "filetype:pdf"):
		return m.fileResults(domain, "pdf", []mockEntry{
			{"Annual Report 2024", "/documents/annual-report-2024.pdf", "Complete financial statements and business overview for fiscal year 2024..."},
			{"Employee Handbook v3.2", "/hr/employee-handbook.pdf", "Internal employee handbook covering company policies, benefits, and procedures..."},
			{"Product Catalog Q4", "/marketing/product-catalog-q4.pdf", "Full product lineup with specifications, pricing, and availability details..."},
		})

	case strings.Contains(q, "filetype:xls") || strings.Contains(q, "filetype:xlsx"):
		return m.fileResults(domain, "xlsx", []mockEntry{
			{"Budget Planning 2024", "/finance/budget-2024.xlsx", "Departmental budget allocations and expense tracking spreadsheet..."},
			{"Customer List Export", "/data/customers-export.xlsx", "Exported customer database with contact information and account details..."},
			{"Inventory Report", "/ops/inventory-q3.xlsx", "Warehouse inventory levels, SKU data, and reorder thresholds..."},
		})

	case strings.Contains(q, "filetype:sql"):
		return m.fileResults(domain, "sql", []mockEntry{
			{"Database Backup - users_db", "/backups/users_db_backup.sql", "Full database dump containing user tables, roles, and permissions..."},
			{"Schema Migration v45", "/db/migration_045.sql", "ALTER TABLE statements modifying user_accounts and payment_records..."},
		})

	case strings.Contains(q, "filetype:log"):
		return m.fileResults(domain, "log", []mockEntry{
			{"Apache Access Log", "/logs/access.log", "192.168.1.105 - - [15/Mar/2024:10:15:32] \"GET /admin/config HTTP/1.1\" 200..."},
			{"Application Error Log", "/logs/error.log", "[ERROR] 2024-03-15 Failed authentication attempt for user admin@..."},
			{"Debug Trace Log", "/logs/debug.log", "[DEBUG] Database connection string: mysql://root:p4ssw0rd@localhost..."},
		})

	case strings.Contains(q, "filetype:env"):
		return m.fileResults(domain, "env", []mockEntry{
			{"Environment Configuration", "/.env", "DB_HOST=localhost DB_USER=root DB_PASS=*** API_KEY=sk-..."},
			{"Production Environment", "/.env.production", "AWS_ACCESS_KEY_ID=AKIA*** AWS_SECRET_ACCESS_KEY=***..."},
		})

	case strings.Contains(q, "filetype:bak") || strings.Contains(q, "filetype:backup"):
		return m.fileResults(domain, "bak", []mockEntry{
			{"Website Backup Archive", "/backups/site-backup-20240315.bak", "Full website backup including database, uploads, and configuration files..."},
			{"Database Backup", "/backups/db-full-backup.bak", "Complete database backup with all tables, stored procedures, and triggers..."},
		})

	case strings.Contains(q, "inurl:admin"):
		return m.pageResults(domain, []mockEntry{
			{"Administration Panel", "/admin/", "Admin control panel. Please login with your administrator credentials to continue..."},
			{"Admin Dashboard — System Overview", "/admin/dashboard", "System health: OK | Active users: 1,247 | Server load: 23% | Last backup: 2h ago..."},
			{"User Management — Admin", "/admin/users", "Manage user accounts, roles, and permissions. Total registered users: 5,892..."},
		})

	case strings.Contains(q, "inurl:login"):
		return m.pageResults(domain, []mockEntry{
			{"Login — Secure Access", "/login", "Enter your username and password to access your account. Forgot password?..."},
			{"Staff Portal Login", "/staff/login", "Internal staff portal. Use your corporate credentials to sign in..."},
		})

	case strings.Contains(q, "wp-admin") || strings.Contains(q, "wp-login"):
		return m.pageResults(domain, []mockEntry{
			{"WordPress — Log In", "/wp-login.php", "Powered by WordPress. Username or Email Address. Password. Remember Me..."},
			{"Dashboard — WordPress Admin", "/wp-admin/", "WordPress 6.4.2 | At a Glance: 156 Posts, 42 Pages, 12 Comments..."},
		})

	case strings.Contains(q, "cpanel") || strings.Contains(q, "webmail"):
		return m.pageResults(domain, []mockEntry{
			{"cPanel Login", ":2083/login", "cPanel hosting control panel. Enter your credentials to manage your hosting account..."},
			{"Webmail Login", ":2096/", "Roundcube Webmail — access your email from anywhere. Select your email client..."},
		})

	case strings.Contains(q, `"index of /"`):
		return m.pageResults(domain, []mockEntry{
			{"Index of /uploads/", "/uploads/", "Name  Last modified  Size  Description\n../  reports/  images/  documents/  backup-2024.zip 450M..."},
			{"Index of /data/", "/data/", "Name  Last modified  Size\n../  exports/  csv/  config.json 2.3K  users.sql 15M..."},
		})

	case strings.Contains(q, `"parent directory"`):
		return m.pageResults(domain, []mockEntry{
			{"Index of /var/www/html/", "/var/www/html/", "[Parent Directory]  assets/  includes/  wp-content/  .htaccess 1.2K  wp-config.php 3.1K..."},
		})

	case strings.Contains(q, "inurl:ftp"):
		return m.pageResults(domain, []mockEntry{
			{"Index of /pub/", "/ftp/pub/", "FTP Directory listing: software/  docs/  updates/  releases-2024.tar.gz 890M..."},
		})

	case strings.Contains(q, `"backup"`):
		return m.pageResults(domain, []mockEntry{
			{"Index of /backups/", "/backups/", "full-site-backup.tar.gz 2.1G  db-dump-20240315.sql.gz 450M  config-backup.zip 12M..."},
		})

	case strings.Contains(q, `"sql syntax"`):
		return m.pageResults(domain, []mockEntry{
			{"Search Results Error", "/search?q=test%27", "You have an error in your SQL syntax; check the manual... near ''test'' at line 1..."},
			{"Product Details", "/product?id=1337", "ERROR 1064 (42000): You have an error in your SQL syntax... MySQL server version..."},
		})

	case strings.Contains(q, `"warning:"`) || strings.Contains(q, "mysql_"):
		return m.pageResults(domain, []mockEntry{
			{"Products Page", "/products.php", "Warning: mysql_connect(): Access denied for user 'webapp'@'localhost'..."},
		})

	case strings.Contains(q, `"fatal error"`):
		return m.pageResults(domain, []mockEntry{
			{"Internal Error", "/includes/db.php", "Fatal error: Uncaught PDOException in /var/www/html/includes/db.php:42..."},
		})

	case strings.Contains(q, `"stack trace"`):
		return m.pageResults(domain, []mockEntry{
			{"Server Error", "/api/v2/users", "Server Error in '/' Application. Stack Trace: System.NullReferenceException at..."},
		})

	case strings.Contains(q, `"begin rsa private key"`):
		return m.fileResults(domain, "key", []mockEntry{
			{"Server Private Key", "/ssl/server.key", "-----BEGIN RSA PRIVATE KEY----- MIIEpAIBAAKCAQEA0Z3VS5JJcds3xfn..."},
		})

	case strings.Contains(q, "api_key") || strings.Contains(q, "apikey"):
		return m.pageResults(domain, []mockEntry{
			{"Configuration Settings", "/config/settings.js", "const config = { api_key: 'sk-proj-abc123...', endpoint: 'https://api..."},
			{"API Documentation", "/docs/api-reference", "Authentication: Include your API key in the header: X-API-Key: your_key_here..."},
		})

	case strings.Contains(q, "secret_key") || strings.Contains(q, "access_key"):
		return m.pageResults(domain, []mockEntry{
			{"AWS Configuration", "/config/aws.js", "module.exports = { accessKeyId: 'AKIA...', secretAccessKey: 's3cr3t...'..."},
		})

	case strings.Contains(q, "phpmyadmin"):
		return m.pageResults(domain, []mockEntry{
			{"phpMyAdmin 5.2.1", "/phpmyadmin/", "phpMyAdmin — MySQL Database Administration. Server: localhost | Database: webapp_prod..."},
		})

	case strings.Contains(q, "filetype:mdb") || strings.Contains(q, "filetype:accdb"):
		return m.fileResults(domain, "mdb", []mockEntry{
			{"Customer Database", "/data/customers.mdb", "Microsoft Access database containing customer records, orders, and payment history..."},
		})

	case strings.Contains(q, "filetype:yml") || strings.Contains(q, "filetype:yaml"):
		return m.fileResults(domain, "yml", []mockEntry{
			{"Application Config", "/config/app.yml", "database: host: db.internal port: 5432 credentials: user: admin pass: ..."},
			{"CI/CD Pipeline", "/.github/workflows/deploy.yml", "name: Deploy to Production on: push: branches: [main] env: AWS_REGION: ..."},
		})

	case strings.Contains(q, "filetype:json") && strings.Contains(q, "password"):
		return m.fileResults(domain, "json", []mockEntry{
			{"Package Config", "/config/database.json", `{"host":"db.prod","port":5432,"user":"admin","password":"P@ssw0rd!123",...}`},
		})

	case strings.Contains(q, "filetype:ini"):
		return m.fileResults(domain, "ini", []mockEntry{
			{"PHP Configuration", "/config/php.ini", "[database] host=localhost user=root password=toor dbname=production..."},
		})

	case strings.Contains(q, "inurl:.git"):
		return m.pageResults(domain, []mockEntry{
			{"Git Repository Index", "/.git/HEAD", "ref: refs/heads/main — Exposed git repository allowing source code download..."},
			{"Git Config", "/.git/config", "[core] repositoryformatversion=0 [remote \"origin\"] url=https://github.com/..."},
		})

	case strings.Contains(q, "docker-compose"):
		return m.fileResults(domain, "yml", []mockEntry{
			{"Docker Compose Stack", "/docker-compose.yml", "version: '3.8' services: web: image: nginx db: image: postgres environment: POSTGRES_PASSWORD:..."},
		})

	case strings.Contains(q, "filetype:htaccess"):
		return m.fileResults(domain, "htaccess", []mockEntry{
			{"Apache Configuration", "/.htaccess", "RewriteEngine On RewriteRule ^admin$ /secret-admin.php [L] AuthType Basic AuthUserFile /etc/..."},
		})

	case strings.Contains(q, "connectionstring"):
		return m.fileResults(domain, "config", []mockEntry{
			{"Web.config", "/web.config", `<connectionStrings><add name="DefaultConnection" connectionString="Server=db;Database=prod;User=sa;Password=..."/></connectionStrings>`},
		})

	case strings.Contains(q, "filetype:csv"):
		return m.fileResults(domain, "csv", []mockEntry{
			{"User Export Data", "/exports/users.csv", "id,name,email,phone,address\n1,John Doe,john@...,555-0123,123 Main St..."},
			{"Transaction Records", "/data/transactions.csv", "date,amount,account,description\n2024-03-15,4500.00,ACC-001,Wire Transfer..."},
		})

	// ─── OSINT / Keyword-only queries ───────────────────────────
	case strings.Contains(q, "site:twitter.com") || strings.Contains(q, "site:x.com"):
		kw := keyword
		if kw == "" {
			kw = "target_user"
		}
		return []SearchResult{
			{Title: kw + " (@" + kw + ") / X", URL: "https://x.com/" + kw, Snippet: "Posts, replies, and media from @" + kw + " on X (formerly Twitter)..."},
			{Title: kw + " on X: \"Latest post...\"", URL: "https://x.com/" + kw + "/status/1234567890", Snippet: "Tweet by @" + kw + " — discussing recent developments and sharing updates..."},
		}

	case strings.Contains(q, "site:github.com"):
		kw := keyword
		if kw == "" {
			kw = "target_user"
		}
		return []SearchResult{
			{Title: kw + " · GitHub", URL: "https://github.com/" + kw, Snippet: "GitHub profile for " + kw + ". Repositories, contributions, and activity..."},
			{Title: kw + "/dotfiles: Personal config", URL: "https://github.com/" + kw + "/dotfiles", Snippet: "Personal configuration files and shell scripts. Last updated 3 days ago..."},
		}

	case strings.Contains(q, "site:linkedin.com"):
		kw := keyword
		if kw == "" {
			kw = "target_user"
		}
		return []SearchResult{
			{Title: kw + " - LinkedIn", URL: "https://linkedin.com/in/" + kw, Snippet: "View " + kw + "'s profile on LinkedIn. Experience, skills, and network connections..."},
		}

	case strings.Contains(q, "site:pastebin.com"):
		kw := keyword
		if kw == "" {
			kw = "target_user"
		}
		return []SearchResult{
			{Title: kw + " - Pastebin.com", URL: "https://pastebin.com/abc123", Snippet: "Paste containing references to " + kw + ". Posted anonymously 2 weeks ago..."},
			{Title: kw + " credentials dump", URL: "https://pastebin.com/xyz789", Snippet: "Data dump mentioning " + kw + " — email addresses and hashed passwords found..."},
		}

	case strings.Contains(q, "site:reddit.com"):
		kw := keyword
		if kw == "" {
			kw = "target_user"
		}
		return []SearchResult{
			{Title: "u/" + kw + " — Reddit profile", URL: "https://reddit.com/u/" + kw, Snippet: "Reddit user profile for " + kw + ". Karma: 12,345 | Posts: 89 | Comments: 567..."},
			{Title: "Anyone know " + kw + "? : r/OSINT", URL: "https://reddit.com/r/OSINT/comments/abc", Snippet: "Discussion thread about " + kw + " in the OSINT subreddit..."},
		}

	case strings.Contains(q, "site:gist.github.com") || strings.Contains(q, "site:gitlab.com"):
		kw := keyword
		if kw == "" {
			kw = "target_user"
		}
		return []SearchResult{
			{Title: kw + " — Gist", URL: "https://gist.github.com/" + kw + "/abc123", Snippet: "Code snippet posted by " + kw + " containing API integration examples..."},
		}

	case strings.Contains(q, "inurl:forum") || strings.Contains(q, "inurl:thread"):
		kw := keyword
		if kw == "" {
			kw = "target_user"
		}
		return []SearchResult{
			{Title: kw + " — Forum Post", URL: "https://forum.example.com/thread/" + kw, Snippet: "Forum discussion mentioning " + kw + ". Multiple replies and community feedback..."},
		}

	case strings.Contains(q, "intext:password") || strings.Contains(q, "intext:passwd"):
		kw := keyword
		if kw == "" {
			kw = "target_user"
		}
		return []SearchResult{
			{Title: kw + " — Credential Leak", URL: "https://leaked-data.example.com/" + kw, Snippet: "Exposed credentials containing " + kw + ": email and password combinations found in data breach..."},
		}

	// SEO queries
	case strings.Contains(q, "site:") && !strings.Contains(q, "filetype:") && !strings.Contains(q, "inurl:") && !strings.Contains(q, "intitle:") && !strings.Contains(q, "intext:"):
		return m.pageResults(domain, []mockEntry{
			{domain + " — Official Website", "/", "Welcome to " + domain + ". We provide innovative solutions for businesses worldwide..."},
			{"About Us — " + domain, "/about", "Learn about our mission, team, and the values that drive " + domain + " forward..."},
			{"Contact — " + domain, "/contact", "Get in touch with our team. Email: info@" + domain + " | Phone: +1 (555) 000-0000..."},
		})

	default:
		// Do NOT fabricate fake data if no realistic mock matches exist
		return []SearchResult{}
	}
}

// mockEntry holds template data for a mock result
type mockEntry struct {
	title   string
	path    string
	snippet string
}

// fileResults creates results that look like file discoveries
func (m *MockSearchEngine) fileResults(domain, ext string, entries []mockEntry) []SearchResult {
	var results []SearchResult
	for _, e := range entries {
		results = append(results, SearchResult{
			Title:   e.title + " — " + domain,
			URL:     "https://" + domain + e.path,
			Snippet: e.snippet,
		})
	}
	return results
}

// pageResults creates results that look like page discoveries
func (m *MockSearchEngine) pageResults(domain string, entries []mockEntry) []SearchResult {
	var results []SearchResult
	for _, e := range entries {
		results = append(results, SearchResult{
			Title:   e.title,
			URL:     "https://" + domain + e.path,
			Snippet: e.snippet,
		})
	}
	return results
}
