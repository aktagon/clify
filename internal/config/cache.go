package config

import (
	"encoding/json"
	"fmt"
	"clify/internal/models"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultCacheDir  = ".clify"
	DefaultCacheFile = "cache.json"
	CacheExpiry      = 24 * time.Hour
)

type CacheManager struct {
	filePath string
	cache    map[string]models.CacheEntry
}

func NewCacheManager() *CacheManager {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	cacheDir := filepath.Join(home, DefaultCacheDir)
	filePath := filepath.Join(cacheDir, DefaultCacheFile)

	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		fmt.Printf("Warning: failed to create cache directory: %v\n", err)
	}

	cm := &CacheManager{
		filePath: filePath,
		cache:    make(map[string]models.CacheEntry),
	}

	if err := cm.loadCache(); err != nil {
		fmt.Printf("Warning: failed to load cache: %v\n", err)
	}
	return cm
}

func (cm *CacheManager) Get(query string) (string, bool) {
	entry, exists := cm.cache[query]
	if !exists {
		return "", false
	}

	// Check if cache entry is expired
	if time.Since(entry.Timestamp) > CacheExpiry {
		delete(cm.cache, query)
		if err := cm.saveCache(); err != nil {
			fmt.Printf("Warning: failed to save cache after expiry cleanup: %v\n", err)
		}
		return "", false
	}

	return entry.Response, true
}

func (cm *CacheManager) Set(query, response string) error {
	cm.cache[query] = models.CacheEntry{
		Query:     query,
		Response:  response,
		Timestamp: time.Now(),
	}

	return cm.saveCache()
}

func (cm *CacheManager) loadCache() error {
	data, err := os.ReadFile(cm.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Cache file doesn't exist yet, that's okay
		}
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	if err := json.Unmarshal(data, &cm.cache); err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	return nil
}

func (cm *CacheManager) saveCache() error {
	data, err := json.MarshalIndent(cm.cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	if err := os.WriteFile(cm.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

func (cm *CacheManager) Clear() error {
	cm.cache = make(map[string]models.CacheEntry)
	return cm.saveCache()
}

func (cm *CacheManager) Size() int {
	return len(cm.cache)
}

// GetSearchHistory returns a slice of recent search queries
func (cm *CacheManager) GetSearchHistory() []string {
	queries := make([]string, 0, len(cm.cache))
	for query := range cm.cache {
		queries = append(queries, query)
	}
	return queries
}
