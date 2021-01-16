package fifo

import (
	"cache"
	"container/list"
)

type fifo struct {
	// 缓存最大存放个数
	maxBytes int

	// 当从缓存中移除一个值时，调用该回调
	onEvicted func(key string,value interface{})

	// 已使用的字节数，不包含key
	usedBytes int

	ll *list.List
	cache map[string]*list.Element
}

// Set 往cache尾部增加一个元素，如果元素已存在则移至尾部
func (f *fifo) Set(key string, value interface{}) {
	if e,ok := f.cache[key];ok {
		f.ll.MoveToBack(e)
		en := e.Value.(*entry)
		f.usedBytes = f.usedBytes - cache.CalcLen(en.value) + cache.CalcLen(value)
		en.value = value
		return
	}

	en := &entry{key,value}
	e := f.ll.PushBack(en)
	f.cache[key] = e
	f.usedBytes += en.Len()
	if f.maxBytes > 0 && f.usedBytes > f.maxBytes {
		f.DelOldest()
	}
}

func (f *fifo) Get(key string) interface{} {
	if e,ok := f.cache[key]; ok {
		return e.Value.(*entry).value
	}

	return nil
}

func (f *fifo) Del(key string) {
	if e,ok := f.cache[key];ok {
		f.removeElement(e)
	}
}

func (f *fifo) DelOldest() {
	f.removeElement(f.ll.Front())
}

func (f *fifo) Len() int {
	return f.ll.Len()
}

func (f *fifo) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	f.ll.Remove(e)
	en := e.Value.(*entry)
	f.usedBytes -= cache.CalcLen(en)
	delete(f.cache,en.key)
	if f.onEvicted != nil {
		f.onEvicted(en.key,en.value)
	}
}

type entry struct {
	key string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}

func New(maxBytes int,onEvicted func(key string,value interface{})) cache.Cache {
	return &fifo{
		maxBytes: maxBytes,
		onEvicted: onEvicted,
		ll: list.New(),
		cache: make(map[string]*list.Element),
	}
}


