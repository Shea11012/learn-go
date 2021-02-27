package gee_cache

import (
	"gee_cache/lru"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte,error)
}

type GetterFunc func(key string) ([]byte,error)

func (g GetterFunc) Get(key string) ([]byte,error)  {
	return g(key)
}

type cache struct {
	mu sync.Mutex
	lru *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes,nil)
	}
	c.lru.Add(key,value)
}

func (c *cache) get(key string) (ByteView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes,nil)
	}

	if v,ok := c.lru.Get(key);ok {
		return v.(ByteView),ok
	}

	return ByteView{},false
}
