package main

import (
	"fmt"
)

func genCacheKey(f string, v1, v2 int) string {
	switch f {
	case "add":
		return genSortedCacheKey(f, v1, v2)
	case "div":
		return genUnSortedCacheKey(f, v1, v2)
	case "mul":
		return genSortedCacheKey(f, v1, v2)
	case "sub":
		return genUnSortedCacheKey(f, v1, v2)
	default:
		panic(fmt.Sprintf("invalid cache key prefix %v", f))
	}
}

// for subtract and divide, order of the query string matters
// x - y != y - x, x / y != y / x
func genSortedCacheKey(f string, v1, v2 int) string {
	if v1 > v2 {
		v1, v2 = v2, v1
	}
	return genUnSortedCacheKey(f, v1, v2)
}

// for add, multiply calculation, order of x, y doesn't matter.
// x + y == y + x , x * y == y * x
func genUnSortedCacheKey(f string, v1, v2 int) string {
	return fmt.Sprintf("%v:%v:%v", f, v1, v2)
}

func calculate(f string, v1, v2 int) int {
	switch f {
	case "add":
		return v1 + v2
	case "div":
		return v1 / v2
	case "mul":
		return v1 * v2
	case "sub":
		return v1 - v2
	default:
		panic(fmt.Sprintf("invalid cache key prefix %v", f))
	}
}

func getResult(f string, x, y int) (int, bool) {
	cacheKey := genCacheKey(f, x, y)
	result, cached := cache.Get(cacheKey)
	if cached {
		cache.IncrCounter()
		return result, cached
	}
	result = calculate(f, x, y)
	cache.SetWithTTL(cacheKey, result, 60)
	return result, cached
}
