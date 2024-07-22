package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var CacheChan = make(chan string)
var Broadcast = make(chan map[string]string)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust this for more secure origin checks
	},
}

var clients = make(map[*websocket.Conn]bool)

type CacheData struct {
	Key        string
	Value      string
	ExpiryTime time.Time
}

type LRUCache struct {
	capacity int
	cache    map[string]*CacheData
	lruList  linkedList
	mutex    sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*CacheData),
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
		n.Value = value
		n.ExpiryTime = time.Now().Add(expiration)
		c.removeSpecificNode(key)
		// Move the most recently used key-value pair to the front of the cache list
		newNode := &node{
			key:        key,
			value:      value,
			ExpiryTime: n.ExpiryTime,
		}
		c.lruList.addToFront(newNode)
		return
	}

	// If the cache is full, remove the least recently used entry
	if len(c.cache) >= c.capacity {
		c.deleteLRUEntry()
	}

	// Add the new entry to the cache
	expirytime := time.Now().Add(expiration)
	n := &node{
		key:        key,
		value:      value,
		ExpiryTime: expirytime,
	}
	c.lruList.addToFront(n)
	ch := &CacheData{
		Key:        key,
		Value:      value,
		ExpiryTime: expirytime,
	}
	c.cache[key] = ch
	CacheChan <- "done"
}

func (c *LRUCache) Getcache(key string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if the key exists in the cache
	if n, ok := c.cache[key]; ok {

		n.ExpiryTime = time.Now().Add(5 * time.Second)
		// c.cache[key] = n
		//Once get cache request is made chace epiry limit will be inceased.
		//Example set cache is made a 10.00.00 AM and now expiry time is 10.00.020 AM.
		// Now get cache request is made at 10.00.03 AM and now expiry time will be 10.00.08 AM.

		// Move the entry to the front of the cache list (most recently used)
		c.removeSpecificNode(key)
		newNode := &node{
			key:        key,
			value:      n.Value,
			ExpiryTime: n.ExpiryTime,
		}
		c.lruList.addToFront(newNode)
		CacheChan <- "done"

		// c.lruList.moveToFront(n)
		return n.Value, true
	}
	return "", false
}
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

func HandleSetcache(Key, Value string, Seconds int) string {
	fmt.Println("Time Got: ", Seconds)

	// secondsInt, err := strconv.Atoi(Seconds)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	fmt.Println("SET CACHE TIME", time.Now())
	fmt.Println("Set Cache Value")
	if Seconds < 1 && Seconds > 9 {
		return "Invalid Time Duration"

	}
	duration := time.Duration(Seconds) * time.Second
	cache.Setcache(Key, Value, duration)

	printCacheContents(cache)

	return "Cache Entered"
}

func CleanupExpiredEntries() {
	for {
		// time.Sleep(1 * time.Second) // Wait for 1 second

		cache.mutex.Lock()

		// Get the current time
		if len(cache.cache) == 0 {
			cache.mutex.Unlock()
			continue
		}

		// Iterate over cache entries and remove expired ones
		for key, n := range cache.cache {
			currentTime := time.Now()
			// Check if the entry has expired
			if currentTime.After(n.ExpiryTime) {
				fmt.Printf("Removing expired cache entry with key: %s and time: %s\n", key, currentTime)
				delete(cache.cache, key)
				cache.removeSpecificNode(key)
				CacheChan <- "done"
			}
		}

		cache.mutex.Unlock()
	}
}

func (c *LRUCache) deleteLRUEntry() {
	if c.lruList.tail != nil {
		delete(c.cache, c.lruList.tail.key)
		c.lruList.removeTail()
	}
}

func SendUpdate() {
	for {
		select {
		case <-CacheChan:
			cacheContent := make(map[string]string)
			for key, data := range cache.cache {
				cacheContent[key] = data.Value
			}
			Broadcast <- cacheContent
		}
	}
}
