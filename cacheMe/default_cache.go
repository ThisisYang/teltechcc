package cacheMe

import (
	"strconv"
	"sync"
	"time"
)

type valueStruct struct {
	value int
	expTS int64
}

// DefaultCache will use local memory.
// All kv will be stored in a map
// key of the map is the key value
// value is pointer to struct valueStruct which store the value and expiration info
// expiration ts will be the number of seconds elapsed since January 1, 1970 UTC
// kv can expired (deleted) in 2 ways
// 1. when accessing the cache via Get method, delete the kv if expired
// 2. there will be a goroutine running in background and scan the map in every 5 seconds
type DefaultCache struct {
	mutex *sync.Mutex
	val   map[string]*valueStruct
	done  chan struct{}
	hit   int
}

// NewDefaultClient return a new defaultCache
// Also create a goroutine that periodically expire keys
func NewDefaultClient() *DefaultCache {
	v := make(map[string]*valueStruct)
	done := make(chan struct{})
	mutex := &sync.Mutex{}
	// init cahce, also, start a goroutine, periodically check and expire cache data
	c := &DefaultCache{
		mutex: mutex,
		val:   v,
		done:  done,
	}
	go c.cronJob()
	return c
}

// Get will get value and extend TTL if exist.
// If not, return 0 and false
func (c *DefaultCache) Get(key string) (int, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	val, ok := c.val[key]
	if ok == false {
		return 0, false
	}
	if isExpired(val.expTS) {
		delete(c.val, key)
		return 0, false
	}

	val.expTS += 60
	return val.value, ok
}

// SetWithTTL will set the key value, and set expiration to 60 second
func (c *DefaultCache) SetWithTTL(key string, value int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	exp := time.Now().Unix() + 60
	c.val[key] = &valueStruct{value: value, expTS: exp}
}

// Ping return nil
func (c *DefaultCache) Ping() error { return nil }

// Close will simply close done channel.
// all goroutines should monitor done channel and exit when done channel closed
func (c *DefaultCache) Close() {
	close(c.done)
}

// IncrCounter will increment hit counter
func (c *DefaultCache) IncrCounter() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.hit++
}

// GetCounter will return hit counter
func (c *DefaultCache) GetCounter() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	count := c.hit
	return count
}

// GetSize return number of keys
// This might not be accurate since the cronJob run every minute
func (c *DefaultCache) GetSize() int {
	return len(c.val)
}

// Flush assign new map to val
func (c *DefaultCache) Flush() {
	c.val = make(map[string]*valueStruct)
}

// cronJob will run periodically in background
// scan all values in map and check if they are expired
// if so, delete them
func (c *DefaultCache) cronJob() {
	// scan the map every 5 second
	tickCh := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-c.done:
			tickCh.Stop()
			return
		case <-tickCh.C:
			c.mutex.Lock()
			for k, v := range c.val {
				if isExpired(v.expTS) {
					delete(c.val, k)
				}
			}
			c.mutex.Unlock()

		}
	}
}

func isExpired(expTS int64) bool {
	if expTS < time.Now().Unix() {
		return true
	}
	return false
}

func stringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}
