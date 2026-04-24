package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type FavoriteStore struct {
	mu        sync.RWMutex
	favorites map[string]bool
	filePath  string
}

func NewFavoriteStore() *FavoriteStore {
	store := &FavoriteStore{
		favorites: make(map[string]bool),
		filePath:  getFavoritesPath(),
	}
	store.load()
	return store
}

func (fs *FavoriteStore) AddFavorite(id string) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.favorites[id] = true
	fs.save()
}

func (fs *FavoriteStore) RemoveFavorite(id string) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	delete(fs.favorites, id)
	fs.save()
}

func (fs *FavoriteStore) IsFavorite(id string) bool {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.favorites[id]
}

func (fs *FavoriteStore) GetFavorites() []string {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	var result []string
	for id := range fs.favorites {
		result = append(result, id)
	}
	return result
}

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

func getFavoritesPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".cti-dork", "favorites.json")
}
