package lru

import (
	"github.com/matryer/is"
	"testing"
)

func TestOnEvicted(t *testing.T) {
	i := is.New(t)
	keys := make([]string,0,8)
	onEvicted := func(key string,value interface{}) {
		keys = append(keys,key)
	}

	cache := New(16,onEvicted)
	cache.Set("k1",1)
	cache.Set("k2",2)
	cache.Get("k1")
	cache.Set("k3",3)
	cache.Get("k1")
	cache.Set("k4",4)

	expected := []string{"k2","k3"}
	i.Equal(expected,keys)
	i.Equal(2,cache.Len())
}
