package cacheMe

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var redisCounter = "hit"

// RedisClient implemented cacheClient interface
// use redis as cache backend
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient return a new RedisClient
func NewRedisClient(redisURL string) *RedisClient {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}
	c := redis.NewClient(opt)
	return &RedisClient{client: c}
}

// Close will close connection
func (c *RedisClient) Close() {
	c.client.Close()
}

// Get will use pipeline, return value, true if key exist, otherwise 0, false
// get the value if cached and extend TTL for 60 seconds
func (c *RedisClient) Get(key string) (int, bool) {

	pipe := c.client.TxPipeline()

	g := pipe.Get(key)
	pipe.Expire(key, time.Minute)

	_, err := pipe.Exec()
	if err == redis.Nil {
		return 0, false
	} else if err != nil {
		fmt.Printf("get key err: %v\n", err)
	}

	val, err := g.Result()
	if err == redis.Nil {
		fmt.Println("got nil")
		return 0, false
	} else if err != nil {
		fmt.Printf("get key err: %v\n", err)
	}

	v, err := stringToInt(val)
	if err != nil {
		return 0, false
	}
	return v, true
}

// SetWithTTL will set kv in redis with TTL
func (c *RedisClient) SetWithTTL(key string, value int) {
	err := c.client.Set(key, value, time.Minute).Err()
	if err != nil {
		fmt.Printf("set key err: %v\n", err)
	}

}

// Ping will ping redis
func (c *RedisClient) Ping() error {
	err := c.client.Ping().Err()
	return err
}

// IncrCounter will increment hit counter
func (c *RedisClient) IncrCounter() {
	err := c.client.Incr(redisCounter).Err()
	if err != nil {
		fmt.Printf("increment counter err: %v\n", err)
	}
}

// GetCounter will get hit counter
func (c *RedisClient) GetCounter() int {
	val := c.client.Get(redisCounter).Val()
	v, _ := stringToInt(val)
	return v
}

// GetSize will return size of DB, including counter
func (c *RedisClient) GetSize() int {
	val := c.client.DBSize().Val()
	return int(val)
}

// Flush will flush redis db
func (c *RedisClient) Flush() {
	c.client.FlushDB()
}
