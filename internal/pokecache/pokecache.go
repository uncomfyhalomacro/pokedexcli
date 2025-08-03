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
	pk.mu.Lock()
	defer pk.mu.Unlock()
	tick := time.Tick(pk.expiryInterval)
	for range tick {
		currentTime := time.Now()
		var toBeDeleted []string
		for k, v := range pk.entries {
			addedTime := v.createdAt.Add(pk.expiryInterval)
			if currentTime.After(addedTime) {
				log.Printf("Cache entry has expired. Entry: %s\n", k)
				toBeDeleted = append(toBeDeleted, k)
			}
		}
		for _, k := range toBeDeleted {
			log.Printf("Deleting cache entry: %s\n", k)
			delete(pk.entries, k)
			log.Printf("Deleted cache entry: %s\n", k)
		}
	}
}

func NewPokeCache(interval time.Duration) PokeCache {
    	emptyEntries := make(map[string]pokeCacheEntry)
	pk := PokeCache{
		entries:        emptyEntries,
		expiryInterval: interval,
	}
	go pk.reapLoop()
	return pk
}

func DefaultPokeCache() PokeCache {
	return NewPokeCache(8)
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
