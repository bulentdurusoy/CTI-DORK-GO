package search

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)


//  ~/.cti-dork/config.json
type SearchConfig struct {
	SerperAPIKey string `json:"serperApiKey"`
}

var (
	configOnce     sync.Once
	cachedConfig   *SearchConfig
	configFilePath string
)

// returns the path to the config JSON file
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

// search configuration from disk

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

// search configuration to disk
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

//  the last loaded config without rereading disk

func GetCachedConfig() *SearchConfig {
	if cachedConfig != nil {
		return cachedConfig
	}
	return LoadConfig()
}

// true if a Serper API key is configured
func HasAPIKey() bool {
	cfg := GetCachedConfig()
	return cfg.SerperAPIKey != ""
}

// masked version of the API key for display.
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
