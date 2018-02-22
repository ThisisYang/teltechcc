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

// DefaultCache will use local memory
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

// Get will get value and extend TTL
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

// SetWithTTL will set the key value, and set expiration ts
func (c *DefaultCache) SetWithTTL(key string, value, seconds int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	exp := time.Now().Unix() + 60
	c.val[key] = &valueStruct{value: value, expTS: exp}
}

// Ping return nil
func (c *DefaultCache) Ping() error { return nil }

// Close will simply close done channel. all goroutines should exit when done channel closed
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

// cronJob will run periodically in background
// check if scan all values in map check if any is expired then delete them
func (c *DefaultCache) cronJob() {
	tickCh := time.NewTicker(1 * time.Second)

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
