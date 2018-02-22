package cacheMe

import (
	"github.com/alicebob/miniredis"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	add := "redis://" + s.Addr()
	redisC := NewRedisClient(add)

	s.Set("foo", "5")
	s.SetTTL("foo", 60*time.Second)
	s.Set("bar", "a")

	cases := []struct {
		name    string
		key     string
		expVal  int
		expBool bool
	}{
		{name: "case 1", key: "foo", expVal: 5, expBool: true},
		{name: "case 2", key: "bar", expVal: 0, expBool: false},
		{name: "case 3", key: "foobar", expVal: 0, expBool: false},
	}

	for _, c := range cases {
		gotVal, gotBool := redisC.Get(c.key)
		if gotVal != c.expVal {
			t.Errorf("error on: %v\ngot val:\n %v \nexp val\n %v \n", c.name, gotVal, c.expVal)
		}
		if gotBool != c.expBool {
			t.Errorf("error on: %v\ngot bool:\n %v \nexp bool\n %v \n", c.name, gotBool, c.expBool)
		}
	}
}

func TestSetWithTTL(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	add := "redis://" + s.Addr()
	redisC := NewRedisClient(add)

	s.Set("foo", "5")
	s.SetTTL("foo", 30*time.Second)

	cases := []struct {
		name string
		key  string
		val  int
		sec  int
	}{
		{name: "case 1", key: "foo", val: 5},
		{name: "case 2", key: "bar", val: 0},
		{name: "case 3", key: "foo", val: 2},
	}

	for _, c := range cases {
		redisC.SetWithTTL(c.key, c.val)
		ttl := s.TTL(c.key)
		if s.Exists(c.key) == false {
			t.Errorf("error on: %v\nkey: %v not set", c.name, c.key)
		}
		if ttl != time.Minute {
			t.Errorf("error on: %v\ngot ttl:\n %v \nexp ttl\n %v \n", c.name, ttl, time.Minute)
		}
	}
}

func TestIncrCounter(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()
	add := "redis://" + s.Addr()
	redisC := NewRedisClient(add)
	s.Set(redisCounter, "5")

	redisC.IncrCounter()
	got, err := s.Get(redisCounter)
	if err != nil {
		t.Errorf("error get counter %v", err)
	}
	if got != "6" {
		t.Errorf("increment counter err, exp: 6, got: %v\n", got)
	}
}

func TestGetCounter(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()
	add := "redis://" + s.Addr()
	redisC := NewRedisClient(add)
	s.Set(redisCounter, "5")
	got := redisC.GetCounter()
	if got != 5 {
		t.Errorf("redis get counter err, exp: 5, got: %v\n", got)
	}
}

func TestGetSize(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	add := "redis://" + s.Addr()
	redisC := NewRedisClient(add)

	s.Set("foo", "5")
	s.Set("bar", "6")
	got := redisC.GetSize()
	if got != 2 {
		t.Errorf("redis get size err, exp: 2, got: %v\n", got)
	}
}

func TestFlush(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	add := "redis://" + s.Addr()
	redisC := NewRedisClient(add)
	s.Set("foo", "5")
	s.Set("bar", "6")
	redisC.Flush()
	if redisC.GetSize() != 0 {
		t.Errorf("redis flush err, exp size: 0, got: %v\n", redisC.GetSize())
	}
}
