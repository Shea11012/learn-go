package cache_test

import (
	"cache"
	"cache/lru"
	"github.com/matryer/is"
	"log"
	"sync"
	"testing"
)

func TestClientCacheGet(t *testing.T) {
	db := map[string]string{
		"key1":"val1",
		"key2":"val2",
		"key3":"val3",
		"key4":"val4",
	}

	getter := cache.GetFunc(func(key string) interface{} {
		log.Println("From db find key",key)
		if val,ok := db[key];ok {
			return val
		}

		return nil
	})

	clientCache := cache.NewClientCache(getter,lru.New(0,nil))
	i := is.New(t)
	var wg sync.WaitGroup
	for k,v := range db {
		wg.Add(1)
		go func(k,v string) {
			defer wg.Done()
			i.Equal(clientCache.Get(k),v)
			i.Equal(clientCache.Get(k),v)
		}(k,v)
	}

	wg.Wait()

	i.Equal(clientCache.Get("unknown"),nil)
	i.Equal(clientCache.Get("unknown"),nil)

	i.Equal(clientCache.Stat().NGet,10)
	i.Equal(clientCache.Stat().NHit,4)
}
