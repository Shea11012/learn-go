package singleflight

import (
	"sync"
)

type call struct {
	wg sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m map[string]*call
}

func (g *Group) Do(key string,fn func() (interface{},error)) (interface{},error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)	// 延迟初始化，提升内存使用效率
	}
	if c,ok := g.m[key];ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val,c.err
	}

	c := new(call)
	c.wg.Add(1)	// 发起请求前加锁
	g.m[key] = c
	g.mu.Unlock()
	c.val,c.err = fn()	// 发起请求
	c.wg.Done()	// 请求结束

	g.mu.Lock()
	delete(g.m,key)
	g.mu.Unlock()

	return c.val,c.err
}
