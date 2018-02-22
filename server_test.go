package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAdd(t *testing.T) {
	setUpLogger(false)
	cases := []struct {
		name, url     string
		expStatusCode int
		expBody       gin.H
		fCache        *fakeCacheClient
	}{
		{
			name: "case fail valid", url: "/add", expStatusCode: 400,
			expBody: gin.H{"err": errMissX.Error()}, fCache: NewFakeCache(),
		},
		{
			name: "case fail valid", url: "/add?x=1", expStatusCode: 400,
			expBody: gin.H{"err": errMissY.Error()}, fCache: NewFakeCache(),
		},
		{
			name: "case uncached", url: "/add?x=1&y=3", expStatusCode: 200,
			expBody: gin.H{"action": "add", "x": 1, "y": 3, "answer": 4, "cached": false},
			fCache:  NewFakeCache(),
		},
		{
			name: "case cached", url: "/add?x=1&y=3", expStatusCode: 200,
			expBody: gin.H{"action": "add", "x": 1, "y": 3, "answer": 4, "cached": true},
			fCache:  &fakeCacheClient{val: map[string]int{"add:1:3": 4}},
		},
		{
			name: "case cached2", url: "/add?x=3&y=1", expStatusCode: 200,
			expBody: gin.H{"action": "add", "x": 3, "y": 1, "answer": 4, "cached": true},
			fCache:  &fakeCacheClient{val: map[string]int{"add:1:3": 4}},
		},
	}

	// setup router
	router := newRouter()
	// Perform a GET request with that handler.
	for _, c := range cases {
		cache = c.fCache
		w := performRequest(router, "GET", c.url)
		jsonEncoded, _ := json.Marshal(c.expBody)

		if w.Code != c.expStatusCode {
			t.Errorf("error on: %v\ngot code:\n %v \nexp code\n %v \n", c.name, w.Code, c.expStatusCode)
		}
		if w.Body.String() != string(jsonEncoded) {
			t.Errorf("error on: %v\ngot body:\n %v \nexp body\n %v \n", c.name, w.Body.String(), string(jsonEncoded))

		}
	}
}

func TestSubtract(t *testing.T) {
	setUpLogger(false)
	cases := []struct {
		name, url     string
		expStatusCode int
		expBody       gin.H
		fCache        *fakeCacheClient
	}{
		{
			name: "case fail valid", url: "/subtract", expStatusCode: 400,
			expBody: gin.H{"err": errMissX.Error()}, fCache: NewFakeCache(),
		},
		{
			name: "case fail valid", url: "/subtract?x=1", expStatusCode: 400,
			expBody: gin.H{"err": errMissY.Error()}, fCache: NewFakeCache(),
		},
		{
			name: "case uncached", url: "/subtract?x=1&y=3", expStatusCode: 200,
			expBody: gin.H{"action": "subtract", "x": 1, "y": 3, "answer": -2, "cached": false},
			fCache:  NewFakeCache(),
		},
		{
			name: "case cached", url: "/subtract?x=1&y=3", expStatusCode: 200,
			expBody: gin.H{"action": "subtract", "x": 1, "y": 3, "answer": -2, "cached": true},
			fCache:  &fakeCacheClient{val: map[string]int{"sub:1:3": -2}},
		},
	}

	// setup router
	router := newRouter()
	// Perform a GET request with that handler.
	for _, c := range cases {
		cache = c.fCache
		w := performRequest(router, "GET", c.url)
		jsonEncoded, _ := json.Marshal(c.expBody)

		if w.Code != c.expStatusCode {
			t.Errorf("error on: %v\ngot code:\n %v \nexp code\n %v \n", c.name, w.Code, c.expStatusCode)
		}
		if w.Body.String() != string(jsonEncoded) {
			t.Errorf("error on: %v\ngot body:\n %v \nexp body\n %v \n", c.name, w.Body.String(), string(jsonEncoded))

		}
	}
}

func TestMultiply(t *testing.T) {
	setUpLogger(false)
	cases := []struct {
		name, url     string
		expStatusCode int
		expBody       gin.H
		fCache        *fakeCacheClient
	}{
		{
			name: "case fail valid", url: "/multiply", expStatusCode: 400,
			expBody: gin.H{"err": errMissX.Error()}, fCache: NewFakeCache(),
		},
		{
			name: "case fail valid", url: "/multiply?x=1", expStatusCode: 400,
			expBody: gin.H{"err": errMissY.Error()}, fCache: NewFakeCache(),
		},
		{
			name: "case uncached", url: "/multiply?x=1&y=3", expStatusCode: 200,
			expBody: gin.H{"action": "multiply", "x": 1, "y": 3, "answer": 3, "cached": false},
			fCache:  NewFakeCache(),
		},
		{
			name: "case cached", url: "/multiply?x=1&y=3", expStatusCode: 200,
			expBody: gin.H{"action": "multiply", "x": 1, "y": 3, "answer": 3, "cached": true},
			fCache:  &fakeCacheClient{val: map[string]int{"mul:1:3": 3}},
		},
		{
			name: "case cached 2", url: "/multiply?x=3&y=1", expStatusCode: 200,
			expBody: gin.H{"action": "multiply", "x": 3, "y": 1, "answer": 3, "cached": true},
			fCache:  &fakeCacheClient{val: map[string]int{"mul:1:3": 3}},
		},
	}

	// setup router
	router := newRouter()
	// Perform a GET request with that handler.
	for _, c := range cases {
		cache = c.fCache
		w := performRequest(router, "GET", c.url)
		jsonEncoded, _ := json.Marshal(c.expBody)

		if w.Code != c.expStatusCode {
			t.Errorf("error on: %v\ngot code:\n %v \nexp code\n %v \n", c.name, w.Code, c.expStatusCode)
		}
		if w.Body.String() != string(jsonEncoded) {
			t.Errorf("error on: %v\ngot body:\n %v \nexp body\n %v \n", c.name, w.Body.String(), string(jsonEncoded))

		}
	}
}

func TestDivide(t *testing.T) {
	setUpLogger(false)
	cases := []struct {
		name, url     string
		expStatusCode int
		expBody       gin.H
		fCache        *fakeCacheClient
	}{
		{
			name: "case fail valid", url: "/divide", expStatusCode: 400,
			expBody: gin.H{"err": errMissX.Error()}, fCache: NewFakeCache(),
		},
		{
			name: "case fail valid", url: "/divide?x=1", expStatusCode: 400,
			expBody: gin.H{"err": errMissY.Error()}, fCache: NewFakeCache(),
		},
		{
			name: "case divide by zero", url: "/divide?x=4&y=0", expStatusCode: 400,
			expBody: gin.H{"err": errDivideByZero.Error()}, fCache: NewFakeCache(),
		},
		{
			name: "case uncached", url: "/divide?x=1&y=3", expStatusCode: 200,
			expBody: gin.H{"action": "divide", "x": 1, "y": 3, "answer": 0, "cached": false},
			fCache:  NewFakeCache(),
		},
		{
			name: "case cached", url: "/divide?x=1&y=3", expStatusCode: 200,
			expBody: gin.H{"action": "divide", "x": 1, "y": 3, "answer": 0, "cached": true},
			fCache:  &fakeCacheClient{val: map[string]int{"div:1:3": 0}},
		},
		{
			name: "case cached 2", url: "/divide?x=3&y=1", expStatusCode: 200,
			expBody: gin.H{"action": "divide", "x": 3, "y": 1, "answer": 3, "cached": false},
			fCache:  &fakeCacheClient{val: map[string]int{"div:1:3": 3}},
		},
	}

	// setup router
	router := newRouter()
	// Perform a GET request with that handler.
	for _, c := range cases {
		cache = c.fCache
		w := performRequest(router, "GET", c.url)
		jsonEncoded, _ := json.Marshal(c.expBody)

		if w.Code != c.expStatusCode {
			t.Errorf("error on: %v\ngot code:\n %v \nexp code\n %v \n", c.name, w.Code, c.expStatusCode)
		}
		if w.Body.String() != string(jsonEncoded) {
			t.Errorf("error on: %v\ngot body:\n %v \nexp body\n %v \n", c.name, w.Body.String(), string(jsonEncoded))

		}
	}
}

func Test405(t *testing.T) {
	setUpLogger(false)
	cases := []struct {
		name, url     string
		expStatusCode int
		expBody       string
		method        string
	}{
		{
			name: "case post", url: "/divide", expStatusCode: 405,
			expBody: "405 method not allowed", method: "POST",
		},
		{
			name: "case delete", url: "/add", expStatusCode: 405,
			expBody: "405 method not allowed", method: "DELETE",
		},
		{
			name: "case option", url: "/subtract", expStatusCode: 405,
			expBody: "405 method not allowed", method: "OPTION",
		},
	}
	// setup router
	router := newRouter()
	// Perform a GET request with that handler.
	for _, c := range cases {
		w := performRequest(router, c.method, c.url)

		if w.Code != c.expStatusCode {
			t.Errorf("error on: %v\ngot code:\n %v \nexp code\n %v \n", c.name, w.Code, c.expStatusCode)
		}
		if w.Body.String() != c.expBody {
			t.Errorf("error on: %v\ngot body:\n %v \nexp body\n %v \n", c.name, w.Body.String(), c.expBody)

		}
	}
}

func Test404(t *testing.T) {
	setUpLogger(false)
	cases := []struct {
		name, url     string
		expStatusCode int
		expBody       string
		method        string
	}{
		{
			name: "case get", url: "/", expStatusCode: 404,
			expBody: "404 page not found", method: "GET",
		},
		{
			name: "case post", url: "/ok", expStatusCode: 404,
			expBody: "404 page not found", method: "POST",
		},
	}
	// setup router
	router := newRouter()
	// Perform a GET request with that handler.
	for _, c := range cases {
		w := performRequest(router, c.method, c.url)

		if w.Code != c.expStatusCode {
			t.Errorf("error on: %v\ngot code:\n %v \nexp code\n %v \n", c.name, w.Code, c.expStatusCode)
		}
		if w.Body.String() != c.expBody {
			t.Errorf("error on: %v\ngot body:\n %v \nexp body\n %v \n", c.name, w.Body.String(), c.expBody)

		}
	}
}
