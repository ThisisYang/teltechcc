package main

// fakeCacheClient is used for testing purpose only
type fakeCacheClient struct {
	val map[string]int
	err error
}

// NewFakeCache return a new fakeCacheClient
func NewFakeCache() *fakeCacheClient {
	v := make(map[string]int)
	return &fakeCacheClient{
		val: v,
		err: nil,
	}
}
func (f *fakeCacheClient) Get(key string) (int, bool) {
	val, ok := f.val[key]
	return val, ok
}

func (f *fakeCacheClient) SetWithTTL(key string, value, seconds int) {
	f.val[key] = value
}

func (f *fakeCacheClient) Ping() error {
	return f.err
}

func (f *fakeCacheClient) Close() {}

func (f *fakeCacheClient) IncrCounter() {
	f.val["hit"] = f.val["hit"] + 1
}

func (f *fakeCacheClient) GetCounter() int {
	return f.val["hit"]
}

func (f *fakeCacheClient) GetSize() int {
	return len(f.val)
}
