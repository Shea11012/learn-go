package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes int64
	nbytes int64
	ll *list.List
	cache map[string]*list.Element
	onEvicted func(key string,value Value)
}

type entry struct {
	key string
	value Value
}

func (e *entry) len() int64 {
	return int64(len(e.key)) + int64(e.value.Len())
}

type Value interface {
	Len() int
}

func New(maxBytes int64,onEvicted func(string,Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (Value,bool) {
	if ele,ok := c.cache[key];ok {
		c.ll.MoveToFront(ele)
		en := ele.Value.(*entry)
		return en.value,true
	}

	return nil, false
}

func (c *Cache) Add(key string, value Value) {
	if ele,ok := c.cache[key];ok {
		en := ele.Value.(*entry)
		c.nbytes -= en.len()
		c.ll.MoveToFront(ele)
		en.value = value
		c.nbytes += en.len()
		return
	}

	en := &entry{
		key: key,
		value: value,
	}
	ele := c.ll.PushFront(en)
	c.cache[key] = ele
	c.nbytes += en.len()

	for c.maxBytes != 0 && c.nbytes > c.maxBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	c.removeElement(ele)
}

func (c *Cache) removeElement(ele *list.Element)  {
	if ele == nil {
		return
	}

	c.ll.Remove(ele)
	en := ele.Value.(*entry)
	delete(c.cache,en.key)
	c.nbytes -= en.len()
	if c.onEvicted != nil {
		c.onEvicted(en.key,en.value)
	}
}
