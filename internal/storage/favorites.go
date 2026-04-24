package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// FavoriteStore manages user's favorite dork IDs with JSON file persistence
type FavoriteStore struct {
	mu        sync.RWMutex
	favorites map[string]bool
	filePath  string
}

// NewFavoriteStore creates a new FavoriteStore and loads existing favorites from disk
func NewFavoriteStore() *FavoriteStore {
	store := &FavoriteStore{
		favorites: make(map[string]bool),
		filePath:  getFavoritesPath(),
	}
	store.load()
	return store
}

// AddFavorite marks a dork ID as a favorite
func (fs *FavoriteStore) AddFavorite(id string) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.favorites[id] = true
	fs.save()
}

// RemoveFavorite removes a dork ID from favorites
func (fs *FavoriteStore) RemoveFavorite(id string) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	delete(fs.favorites, id)
	fs.save()
}

// IsFavorite checks if a dork ID is marked as a favorite
func (fs *FavoriteStore) IsFavorite(id string) bool {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.favorites[id]
}

// GetFavorites returns all favorite dork IDs
func (fs *FavoriteStore) GetFavorites() []string {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	var result []string
	for id := range fs.favorites {
		result = append(result, id)
	}
	return result
}

// save persists favorites to the JSON file
func (fs *FavoriteStore) save() {
	dir := filepath.Dir(fs.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}

	ids := make([]string, 0, len(fs.favorites))
	for id := range fs.favorites {
		ids = append(ids, id)
	}

	data, err := json.MarshalIndent(ids, "", "  ")
	if err != nil {
		return
	}

	_ = os.WriteFile(fs.filePath, data, 0644)
}

// load reads favorites from the JSON file on disk
func (fs *FavoriteStore) load() {
	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		return
	}

	var ids []string
	if err := json.Unmarshal(data, &ids); err != nil {
		return
	}

	for _, id := range ids {
		fs.favorites[id] = true
	}
}

// getFavoritesPath returns the path to the favorites JSON file in user's home directory
func getFavoritesPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".cti-dork", "favorites.json")
}
