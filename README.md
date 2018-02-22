# teltechcc

This repository contains code and notes to get a simple web application in Go using [Gin framework](https://github.com/gin-gonic/gin) which accepts math problems via the URL and returns the response in JSON. 

### Web Server
There are 4 endpoints: `/add`, `/divide`, `/subtract` and `/multiply`.

Currently, only two arguments will ever be accepted and parsed -- x and y. Both variables have to be `Integer`.

Example:
http://localhost/add?x=2&y=5

If success, it will return status code `200` with JSON response:

`{"action": "add", "x": 2, "y": 5, "answer", 7, "cached": false}`

Otherwise, `400` will be returned if data is invalid or miss variable with JSON response that includes error details.
`404` will be returned if route doesn't exist. `405` will be returned if method is not allowed.

### Cache
Cache is used with TTL 60 second. Two type of caches are available: 
1. [redis](https://redis.io/) 
2. local memory (default, use if you don't want to use redis. But not recommended).

This can be configured via `--redis` flag when run the server.

### health check
 `/health` endpoint will return status code `200` with JSON response. 
 
 Each time access `/health` endpoint, server will `Ping` cache to check if cache is connected, also, return cache hit counter and size of cache.
 
example output:
`{cache: OK, hit: 10, size: 20}`

### Flags
4 flags are available:
```go
    var ip = flag.String("ip", "0.0.0.0", "IP server bind to")
    var port = flag.Int("port", 8000, "port server listen on")
    var redisURL = flag.String("redis", "", "redis url. for example: `redis://localhost:6379`. If not set, will use local memory instead of redis as cache")
    var debug = flag.Bool("debug", false, "boolean field, set to enable debug mode")
```
By default, server will bu functional without passing any flag. Local memory will be used as cache. In this way, you don't have to setup redis.

Note: if you are planning to run in docker container, don't set `port` flag, instead, expose docker port via `-p xx:8000`.

### how to run:
1. You can build binary and run the binary directly.
```sh
$ go build -o foo .
$ ./foo
```
2. You can run build docker image and run in container.
```sh
$ docker build -t {tag} .
$ docker run -d -p 80:8000  {image_id} --debug=true --redis redis://{redis_ip}:{redis_port}/{DB}
```

### Others:
If you want to, you can use other cache backend (`memcached`,`redshift` or even persistent database) as long as you implement `cacheClient` interface.

It is not suggested to use default cache (local memory) as there is not limit on the size of internal map. Also, there will be a goroutine running at the background to scan the entire map every second to remove expired keys. This will lock the memory and block other goroutine accessing it.

### Issues:

Currently it only handles `int` operation. Since `int` in golang is 32 bit, it has range -2147483648 through 2147483647. If x or y has value beyond range, server will return `400`. If the result is beyond the range, server will still return `200` with incorrect value.