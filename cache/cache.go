package cache

import (
	"sync"
)

type Cache interface {
	Set(key string,value interface{})
	Get(key string) interface{}
	Del(key string)
	DelOldest()
	Len() int
}

// DefaultMaxBytes 默认允许占用的最大内存
const DefaultMaxBytes = 1 << 29

type safeCache struct {
	m sync.RWMutex
	cache Cache
	nhit int	// 命中数
	nget int	// 获取数
}

func newSafeCache(cache Cache) *safeCache {
	return &safeCache{
		cache: cache,
	}
}

func (s *safeCache) set(key string, value interface{}) {
	s.m.Lock()
	defer s.m.Unlock()
	s.cache.Set(key,value)
}

func (s *safeCache) get(key string) interface{} {
	s.m.RLock()
	defer s.m.RUnlock()
	s.nget++
	if s.cache == nil {
		return nil
	}

	v := s.cache.Get(key)
	if v != nil {
		s.nhit++
	}
	return v
}

type Stat struct {
	NHit,NGet int
}

func (s *safeCache) stat() *Stat {
	s.m.RLock()
	defer s.m.RUnlock()
	return &Stat{
		NHit: s.nhit,
		NGet: s.nget,
	}
}
