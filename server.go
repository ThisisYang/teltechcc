package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func newServer(ip string, port int) *http.Server {
	addr := fmt.Sprintf("%v:%v", ip, port)
	r := newRouter()
	return &http.Server{Addr: addr, Handler: r}
}

func newRouter() *gin.Engine {
	r := gin.Default()
	// if route match but not get method. return 405
	r.HandleMethodNotAllowed = true
	r.GET("/add", add)
	r.GET("/subtract", subtract)
	r.GET("/multiply", multiply)
	r.GET("/divide", divide)
	r.GET("/health", health)
	return r
}

//add, subtract, multiply, and divide
func add(ctx *gin.Context) {
	x, y := ctx.Query("x"), ctx.Query("y")
	intX, intY, err := addValidation(x, y)
	if err != nil {
		ctx.JSON(400, gin.H{"err": err.Error()})
		return
	}
	Debug.Println("recieved:", intX, intY)
	result, cached := getResult("add", intX, intY)
	ctx.JSON(200, gin.H{"action": "add", "x": intX, "y": intY, "answer": result, "cached": cached})
}

func subtract(ctx *gin.Context) {
	x, y := ctx.Query("x"), ctx.Query("y")
	intX, intY, err := subValidation(x, y)
	if err != nil {
		ctx.JSON(400, gin.H{"err": err.Error()})
		return
	}
	Debug.Println("recieved:", intX, intY)
	result, cached := getResult("sub", intX, intY)

	ctx.JSON(200, gin.H{"action": "subtract", "x": intX, "y": intY, "answer": result, "cached": cached})
}

func multiply(ctx *gin.Context) {
	x, y := ctx.Query("x"), ctx.Query("y")
	intX, intY, err := mulValidation(x, y)
	if err != nil {
		ctx.JSON(400, gin.H{"err": err.Error()})
		return
	}
	Debug.Println("recieved:", intX, intY)
	result, cached := getResult("mul", intX, intY)

	ctx.JSON(200, gin.H{"action": "multiply", "x": intX, "y": intY, "answer": result, "cached": cached})
}

// divide is floor function.
// example: 1 / 3 =0, 4 / 3 = 1
func divide(ctx *gin.Context) {
	x, y := ctx.Query("x"), ctx.Query("y")
	intX, intY, err := divValidation(x, y)
	if err != nil {
		ctx.JSON(400, gin.H{"err": err.Error()})
		return
	}
	Debug.Println("recieved:", intX, intY)
	result, cached := getResult("div", intX, intY)

	ctx.JSON(200, gin.H{"action": "divide", "x": intX, "y": intY, "answer": result, "cached": cached})
}

// health endpoint. return 200 and cache status
func health(ctx *gin.Context) {
	err := cache.Ping()
	if err != nil {
		ctx.JSON(200, gin.H{"cache": err.Error()})
		return
	}
	hit := cache.GetCounter()
	size := cache.GetSize()
	ctx.JSON(200, gin.H{"cache": "OK", "hit": hit, "size": size})
}
