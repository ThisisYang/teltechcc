package main

import (
	"context"
	"flag"
	"github.com/ThisisYang/teltechcc/cacheMe"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cache cacheClient

func main() {
	var (
		ip       = flag.String("ip", "0.0.0.0", "IP server bind to")
		port     = flag.Int("port", 8000, "port server listen on")
		redisURL = flag.String("redis", "", "redis url. for example: `redis://localhost:6379`. If not set, will use local memory instead of redis as cache")
		debug    = flag.Bool("debug", false, "boolean field, set to enable debug mode for both gin server and app")
		flush    = flag.Bool("flush", false, "boolean, set true if to flush db on boot")
	)
	flag.Parse()

	setUpLogger(*debug)

	// if debug is false, set gin server to release mode as well
	if *debug == false {
		gin.SetMode(gin.ReleaseMode)
	}

	if *redisURL != "" {
		cache = cacheMe.NewRedisClient(*redisURL)
	} else {
		cache = cacheMe.NewDefaultClient()
	}

	if *flush {
		cache.Flush()
	}

	defer cache.Close()

	s := newServer(*ip, *port)

	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			Warning.Println("http server closed with err: ", err)
		}
	}()

	waitSingal()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		Warning.Println("Server shutdown failure: ", err)
	}

}

func waitSingal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs
	Info.Println("received signal")
}
