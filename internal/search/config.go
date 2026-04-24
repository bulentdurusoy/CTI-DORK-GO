package search

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// SearchConfig holds application-level search configuration.
// Persisted to ~/.cti-dork/config.json
type SearchConfig struct {
	SerperAPIKey string `json:"serperApiKey"`
}

var (
	configOnce     sync.Once
	cachedConfig   *SearchConfig
	configFilePath string
)

// getConfigPath returns the path to the config JSON file
func getConfigPath() string {
	if configFilePath != "" {
		return configFilePath
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	configFilePath = filepath.Join(home, ".cti-dork", "config.json")
	return configFilePath
}

// LoadConfig reads the search configuration from disk.
// Returns a default (empty) config if the file doesn't exist.
func LoadConfig() *SearchConfig {
	path := getConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		return &SearchConfig{}
	}

	var cfg SearchConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return &SearchConfig{}
	}

	cfg.SerperAPIKey = strings.TrimSpace(cfg.SerperAPIKey)
	cachedConfig = &cfg
	return &cfg
}

// SaveConfig persists the search configuration to disk
func SaveConfig(cfg *SearchConfig) error {
	path := getConfigPath()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	cachedConfig = cfg
	return os.WriteFile(path, data, 0644)
}

// GetCachedConfig returns the last loaded config without re-reading disk.
// Falls back to LoadConfig if no cached version exists.
func GetCachedConfig() *SearchConfig {
	if cachedConfig != nil {
		return cachedConfig
	}
	return LoadConfig()
}

// HasAPIKey returns true if a Serper API key is configured
func HasAPIKey() bool {
	cfg := GetCachedConfig()
	return cfg.SerperAPIKey != ""
}

// MaskAPIKey returns a masked version of the API key for display.
// Shows first 4 and last 4 characters: "691a...bc0c"
func MaskAPIKey(key string) string {
	key = strings.TrimSpace(key)
	if len(key) <= 8 {
		if len(key) == 0 {
			return ""
		}
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
