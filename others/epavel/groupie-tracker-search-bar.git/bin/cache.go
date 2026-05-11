package bin

import (
	"fmt"
	"sync"
	"time"
)

type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
	LastUsed  time.Time
}

type Cache struct {
	items    map[string]CacheItem
	mu       sync.RWMutex // Mutex for thread safety
	maxItems int
}

type CachedData struct {
	PaginatedArtists []Artist `json:"paginated_artists"`
	TotalArtists     int      `json:"total_artists"`
	Message          string   `json:"message"`
}

func NewCache(maxItems int) *Cache {
	return &Cache{
		items:    make(map[string]CacheItem),
		maxItems: maxItems,
	}
}

var (
	queryCache      = NewCache(32)
	artistsCache    = NewCache(1)
	artistCache     = NewCache(16)
	locationsCache  = NewCache(1)
	filterDataCache = NewCache(20)
)

// Generate a unique cache key based on query parameters
func generateCacheKey(pagination int, shuffle bool, searchQuery string, filters Filters) string {
	return fmt.Sprintf("%d-%t-%s-%v", pagination, shuffle, searchQuery, filters)
}

// Set adds a cache item to the cache with a specified expiration time
func (c *Cache) Set(key string, data interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the key already exists
	if _, found := c.items[key]; !found && len(c.items) >= c.maxItems {
		// Evict the least recently used item if the cache is full and the key is new
		c.EvictLRU()
	}

	// Add or update the cache item
	c.items[key] = CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(duration),
		LastUsed:  time.Now(),
	}
}

// Get retrieves a cache item from the cache if it exists and is not expired
func (c *Cache) Get(key string, additionalDuration time.Duration) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()
	if !found || time.Now().After(item.ExpiresAt) {
		return nil, false
	}
	item.LastUsed = time.Now()                              // Update the last used time
	item.ExpiresAt = item.ExpiresAt.Add(additionalDuration) // Extend the expiration time
	c.mu.Lock()
	c.items[key] = item
	c.mu.Unlock()
	return item.Data, true
}

// Delete removes a cache item from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// EvictLRU evicts the least recently used cache item
func (c *Cache) EvictLRU() {
	c.mu.Lock()
	defer c.mu.Unlock()

	var oldestKey string
	var oldestTime time.Time

	for key, item := range c.items {
		if oldestKey == "" || item.LastUsed.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.LastUsed
		}
	}

	if oldestKey != "" {
		delete(c.items, oldestKey)
	}
}

// StartLRUEviction starts a goroutine that evicts the least recently used
// cache item every 5 minutes to prevent memory leaks and minimize cache size
func StartLRUEviction() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			queryCache.EvictLRU()
			artistsCache.EvictLRU()
			artistCache.EvictLRU()
			locationsCache.EvictLRU()
			filterDataCache.EvictLRU()
		}
	}()
}

// LRU (Least Recently Used) cache eviction policy
// When the cache is probably densely populated
// the least recently used item is evicted to make space for new items
// This is a simple and efficient way to manage cache size
// and prevent memory leaks
