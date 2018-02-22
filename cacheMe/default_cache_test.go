package cacheMe

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

// get test DefaultCache without setting map
// map will be set on each test case
func getDC() *DefaultCache {
	return &DefaultCache{
		mutex: &sync.Mutex{},
		done:  make(chan struct{}),
	}
}

func emptyVal() map[string]*valueStruct {
	return make(map[string]*valueStruct)
}

func cachedVal(key string, value int, backOffSecond int64) map[string]*valueStruct {
	return map[string]*valueStruct{
		key: &valueStruct{
			value: value,
			expTS: time.Now().Unix() + backOffSecond,
		},
	}
}

func TestDCGet(t *testing.T) {
	dc := getDC()
	cases := []struct {
		name     string
		existVal map[string]*valueStruct
		getKey   string
		expInt   int
		expBool  bool
		expVal   map[string]*valueStruct
	}{
		{
			name: "not cached", existVal: emptyVal(), getKey: "foo",
			expInt: 0, expBool: false, expVal: emptyVal(),
		},
		{
			name: "cached", existVal: cachedVal("foo", 1, 30), getKey: "foo",
			expInt: 1, expBool: true, expVal: cachedVal("foo", 1, 90),
		},
		{
			name: "expired", existVal: cachedVal("foo", 1, -30), getKey: "foo",
			expInt: 0, expBool: false, expVal: emptyVal(),
		},
	}

	for _, c := range cases {
		dc.val = c.existVal
		gotInt, gotBool := dc.Get(c.getKey)
		if gotInt != c.expInt {
			t.Errorf("error on: %v\ngot int:\n %v \nexp int\n %v \n", c.name, gotInt, c.expInt)
		}
		if gotBool != c.expBool {
			t.Errorf("error on: %v\ngot bool:\n %v \nexp bool\n %v \n", c.name, gotBool, c.expBool)
		}
		if reflect.DeepEqual(dc.val, c.expVal) == false {
			t.Errorf("error on: %v asserting val\ngot val:\n %v \nexp val\n %v \n", c.name, dc.val, c.expVal)
		}
	}
}

func TestDCSetWithTTL(t *testing.T) {
	dc := getDC()
	cases := []struct {
		name     string
		existVal map[string]*valueStruct
		setKey   string
		setVal   int
		expVal   map[string]*valueStruct
	}{
		{
			name: "case 1", existVal: emptyVal(), setKey: "foo",
			setVal: 1, expVal: cachedVal("foo", 1, 60),
		},
	}

	for _, c := range cases {
		dc.val = c.existVal
		dc.SetWithTTL(c.setKey, c.setVal)
		if reflect.DeepEqual(dc.val, c.expVal) == false {
			t.Errorf("error on: %v asserting val\ngot val:\n %v \nexp val\n %v \n", c.name, dc.val, c.expVal)
		}
	}
}

func TestDCIncrCounter(t *testing.T) {
	dc := getDC()
	dc.IncrCounter()
	if dc.hit != 1 {
		t.Errorf("IncrCounter faled, exp: 1, got: %v\n", dc.hit)
	}
}

func TestDCGetCounter(t *testing.T) {
	dc := getDC()
	dc.hit = 10
	gotHit := dc.GetCounter()
	if gotHit != 10 {
		t.Errorf("IncrCounter faled, exp: 10, got: %v\n", gotHit)
	}
}

func TestDCGetSize(t *testing.T) {
	dc := getDC()
	cases := []struct {
		name     string
		existVal map[string]*valueStruct
		expSize  int
	}{
		{name: "case 1", existVal: emptyVal(), expSize: 0},
		{name: "case 2", existVal: cachedVal("foo", 1, 30), expSize: 1},
	}

	for _, c := range cases {
		dc.val = c.existVal
		gotSize := dc.GetSize()
		if gotSize != c.expSize {
			t.Errorf("error on: %v\ngot size:\n %v \nexp size\n %v \n", c.name, gotSize, c.expSize)
		}
	}
}

func TestDCFlush(t *testing.T) {
	dc := getDC()
	cases := []struct {
		name     string
		existVal map[string]*valueStruct
		expVal   map[string]*valueStruct
	}{
		{name: "case 1", existVal: emptyVal(), expVal: emptyVal()},
		{name: "case 2", existVal: cachedVal("foo", 1, 30), expVal: emptyVal()},
	}

	for _, c := range cases {
		dc.val = c.existVal
		dc.Flush()
		if reflect.DeepEqual(dc.val, c.expVal) == false {
			t.Errorf("error on: %v asserting val\ngot val:\n %v \nexp val\n %v \n", c.name, dc.val, c.expVal)
		}
	}
}

func TestIsExpired(t *testing.T) {
	cases := []struct {
		name    string
		expTS   int64
		expBool bool
	}{
		{name: "case 1", expTS: time.Now().Unix(), expBool: false},
		{name: "case 2", expTS: time.Now().Unix() + 10, expBool: false},
		{name: "case 3", expTS: time.Now().Unix() - 10, expBool: true},
	}
	for _, c := range cases {
		gotBool := isExpired(c.expTS)
		if gotBool != c.expBool {
			t.Errorf("error on: %v\ngot bool:\n %v \nexp bool\n %v \n", c.name, gotBool, c.expBool)
		}
	}
}
