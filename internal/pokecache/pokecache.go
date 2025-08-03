package pokecache

import (
	"log"
	"sync"
	"time"
)

type pokeCacheEntry struct {
	createdAt time.Time
	val       []byte
}
type PokeCache struct {
	mu             sync.Mutex
	entries        map[string]pokeCacheEntry
	expiryInterval time.Duration
}

func (pk *PokeCache) reapLoop() {
	ticker := time.NewTicker(pk.expiryInterval)
	defer ticker.Stop()
	for c := range ticker.C {
		pk.mu.Lock()
		for k, v := range pk.entries {
			if c.After(v.createdAt.Add(pk.expiryInterval)) {
				delete(pk.entries, k)
			}
		}
		pk.mu.Unlock()
	}
}

func NewPokeCache(interval time.Duration) *PokeCache {
	emptyEntries := make(map[string]pokeCacheEntry)
	pk := &PokeCache{
		entries:        emptyEntries,
		expiryInterval: interval,
	}
	go pk.reapLoop()
	return pk
}

func DefaultPokeCache() *PokeCache {
	return NewPokeCache(8 * time.Second)
}

// "key" is the previous and next field URL names
func (pk *PokeCache) Add(key string, newData []byte) {
	log.Println("Adding new cache entry...")
	pk.mu.Lock()
	createdAt := time.Now()
	newEntry := pokeCacheEntry{
		createdAt: createdAt,
		val:       newData,
	}
	pk.entries[key] = newEntry
	pk.mu.Unlock()
}

func (pk *PokeCache) Get(key string) ([]byte, bool) {
	log.Println("Getting cache entry....")
	pk.mu.Lock()
	defer pk.mu.Unlock()
	cacheEntry, ok := pk.entries[key]
	if !ok {
		log.Printf("Cache entry is outdated or does not exist. Entry: %s\n", key)
		return nil, false
	}
	return cacheEntry.val, true
}
