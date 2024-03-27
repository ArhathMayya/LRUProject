package handlers

import (
	"fmt"
	"sync"
	"time"
)

type CacheData struct {
	Key        string
	Value      string
	ExpiryTime time.Time
}

type LRUCache struct {
	capacity int
	cache    map[string]*node
	lruList  linkedList
	mutex    sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*node),
		lruList:  linkedList{},
	}
}

var cache *LRUCache

func (c *LRUCache) Setcache(key string, value string, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if the key is present in the cache
	if n, ok := c.cache[key]; ok {
		// Update the value and expiration time
		n.value = value
		n.ExpiryTime = time.Now().Add(expiration)

		// Move the most recently used key-value pair to the front of the cache list
		c.lruList.moveToFront(n)
		return
	}

	// If the cache is full, remove the least recently used entry
	if len(c.cache) >= c.capacity {
		c.deleteLRUEntry()
	}

	// Add the new entry to the cache
	n := &node{
		key:        key,
		value:      value,
		ExpiryTime: time.Now().Add(expiration),
	}
	c.lruList.addToFront(n)
	c.cache[key] = n
}

func (c *LRUCache) Getcache(key string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if the key exists in the cache
	if n, ok := c.cache[key]; ok {

		n.ExpiryTime = time.Now().Add(5 * time.Second)
		//Once get cache request is made chace epiry limit will be inceased.
		//Example set cache is made a 10.00.00 AM and now expiry time is 10.00.05 AM.
		// Now get cache request is made at 10.00.03 AM and now expiry time will be 10.00.08 AM.

		// Move the entry to the front of the cache list (most recently used)
		c.lruList.moveToFront(n)
		return n.value, true
	}
	return "", false
}

func (c *LRUCache) deleteLRUEntry() {
	if c.lruList.tail != nil {
		delete(c.cache, c.lruList.tail.key)
		c.lruList.removeTail()
	}
}

// Here we will print the cache list for every get and set request
func printCacheContents(cache *LRUCache) {

	fmt.Println("Current Cache Contents:")
	current := cache.lruList.head
	for current != nil {
		fmt.Printf("Key: %s, Value: %v, Expiration: %s\n", current.key, current.value, current.ExpiryTime)
		current = current.next
	}

	fmt.Println("--------------")
}

// I'm using init() just to be cautious about the initialization of cache variable.

func init() {
	cache = NewLRUCache(1024)
}

func HandleGetcache(Key string) (string, bool) {

	fmt.Println("Get Cached Value")
	data, flag := cache.Getcache(Key)
	printCacheContents(cache)

	return data, flag
}

func HandleSetcache(Key, Value string) string {
	fmt.Println("SET CACHE TIME", time.Now())
	fmt.Println("Set Cache Value")
	cache.Setcache(Key, Value, 5*time.Second)
	printCacheContents(cache)

	return "Cache Entered"
}

func CleanupExpiredEntries() {
	for {
		time.Sleep(1 * time.Second) // Wait for 1 second

		cache.mutex.Lock()

		if len(cache.cache) == 0 {
			cache.mutex.Unlock()
			continue
		}

		// Iterate over cache entries and remove expired ones
		for key, n := range cache.cache {

			// check time for every iteration
			currentTime := time.Now()
			// Check if the entry has expired
			if currentTime.After(n.ExpiryTime) {
				fmt.Printf("Removing expired cache entry with key: %s and time: %s\n", key, currentTime)
				delete(cache.cache, key)
				cache.removeNode(n)
			}
		}

		cache.mutex.Unlock()
	}
}
