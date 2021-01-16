package cache

type Getter interface {
	Get(key string) interface{}
}

type GetFunc func(key string) interface{}

func (f GetFunc) Get(key string) interface{} {
	return f(key)
}

type ClientCache struct {
	mainCache *safeCache
	getter Getter
}

func NewClientCache(getter Getter,cache Cache) *ClientCache {
	return &ClientCache{
		mainCache: newSafeCache(cache),
		getter: getter,
	}
}

func (c *ClientCache) Get(key string) interface{} {
	val := c.mainCache.get(key)
	if val != nil {
		return val
	}

	if c.getter != nil {
		val = c.getter.Get(key)
		if val == nil {
			return nil
		}

		c.mainCache.set(key,val)
		return val
	}

	return nil
}

func (c *ClientCache) Set(key string,val interface{})  {
	if val == nil {
		return
	}

	c.mainCache.set(key,val)
}

func (c *ClientCache) Stat() *Stat {
	return c.mainCache.stat()
}
