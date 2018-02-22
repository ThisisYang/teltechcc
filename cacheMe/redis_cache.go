package cacheMe

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var redisCounter = "hit"

type redisClient struct {
	client *redis.Client
}

func NewRedisClient(redisURL string) *redisClient {
	//c, err := redis.DialURL(os.Getenv("REDIS_URL"))
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}
	c := redis.NewClient(opt)
	c.FlushDB()
	return &redisClient{client: c}
}

func (c *redisClient) Close() {
	c.client.Close()
}

// Get will use pipeline, return value, true if key exist, otherwise 0, false
// get the value if cached and extend TTL
func (c *redisClient) Get(key string) (int, bool) {

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

func (c *redisClient) SetWithTTL(key string, value, seconds int) {
	err := c.client.Set(key, value, time.Duration(seconds)*time.Second).Err()
	if err != nil {
		fmt.Printf("set key err: %v\n", err)
	}

}

func (c *redisClient) Ping() error {
	err := c.client.Ping().Err()
	return err
}

func (c *redisClient) IncrCounter() {
	err := c.client.Incr(redisCounter).Err()
	if err != nil {
		fmt.Printf("increment counter err: %v\n", err)
	}
}

func (c *redisClient) GetCounter() int {
	val := c.client.Get(redisCounter).Val()
	v, _ := stringToInt(val)
	return v
}

// GetSize will return size of DB, including counter
func (c *redisClient) GetSize() int {
	val := c.client.DBSize().Val()
	return int(val)
}
