package main

import (
	"testing"
)

func TestGenCacheKey(t *testing.T) {
	cases := []struct {
		name, f, expected string
		v1, v2            int
	}{
		{name: "case 1", f: "add", v1: 4, v2: 0, expected: "add:0:4"},
		{name: "case 2", f: "add", v1: 0, v2: 4, expected: "add:0:4"},
		{name: "case 3", f: "div", v1: 1, v2: 4, expected: "div:1:4"},
		{name: "case 4", f: "div", v1: 4, v2: 1, expected: "div:4:1"},
		{name: "case 5", f: "mul", v1: 4, v2: 0, expected: "mul:0:4"},
		{name: "case 6", f: "mul", v1: 0, v2: 4, expected: "mul:0:4"},
		{name: "case 7", f: "sub", v1: 4, v2: 0, expected: "sub:4:0"},
		{name: "case 8", f: "sub", v1: 0, v2: 4, expected: "sub:0:4"},
	}
	for _, c := range cases {
		got := genCacheKey(c.f, c.v1, c.v2)
		if got != c.expected {
			t.Errorf("error on: %v\ngot:\n %v \nexpected\n %v \n", c.name, got, c.expected)
		}
	}
}

func TestCalculate(t *testing.T) {
	cases := []struct {
		name, f          string
		v1, v2, expected int
	}{
		{name: "case 1", f: "add", v1: 4, v2: 0, expected: 4},
		{name: "case 2", f: "div", v1: 8, v2: 3, expected: 2},
		{name: "case 3", f: "mul", v1: 1, v2: 4, expected: 4},
		{name: "case 4", f: "sub", v1: 4, v2: 1, expected: 3},
	}
	for _, c := range cases {
		got := calculate(c.f, c.v1, c.v2)
		if got != c.expected {
			t.Errorf("error on: %v\ngot:\n %v \nexpected\n %v \n", c.name, got, c.expected)
		}
	}
}

func TestGetResult(t *testing.T) {
	cases := []struct {
		f, name      string
		x, y, expInt int
		expBool      bool
		fCache       *fakeCacheClient
	}{
		{
			name: "case 1", f: "add", x: 4, y: 0, expInt: 4,
			fCache: &fakeCacheClient{val: make(map[string]int)}, expBool: false,
		},
		{
			name: "case 2", f: "add", x: 4, y: 0, expInt: 4,
			fCache: &fakeCacheClient{val: map[string]int{"add:0:4": 4}}, expBool: true,
		},
		{
			name: "case 3", f: "div", x: 3, y: 3, expInt: 1,
			fCache: &fakeCacheClient{val: map[string]int{"add:3:3": 6}}, expBool: false,
		},
	}
	for _, c := range cases {
		// overwrite cache global variable
		cache = c.fCache

		gotInt, gotBool := getResult(c.f, c.x, c.y)
		if gotInt != c.expInt {
			t.Errorf("error on: %v\ngot int:\n %v \nexp int\n %v \n", c.name, gotInt, c.expInt)
		}
		if gotBool != c.expBool {
			t.Errorf("error on: %v\ngot bool:\n %v \nexp bool\n %v \n", c.name, gotBool, c.expBool)
		}
	}
}
