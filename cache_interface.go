package main

type cacheClient interface {
	Getter
	Setter
	Counter
	Flusher
	Ping() error
	Close()
}

// Getter interface implement method of Get
// Get will get the value and renew TTL if key exist
type Getter interface {
	Get(key string) (int, bool)
}

// Setter interface implement method of Set
type Setter interface {
	SetWithTTL(key string, value int)
}

// Flusher implement Flush method
type Flusher interface {
	Flush()
}

// Counter implement IncrCounter, GetCounter and GetSize
type Counter interface {
	IncrCounter()
	GetCounter() int
	GetSize() int
}
