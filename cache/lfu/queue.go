package lfu

import (
	"cache"
	"container/heap"
)

type entry struct {
	key string
	value interface{}
	weight int // 访问次数越多，权重越高
	index int // 在heap中的索引
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}

type queue []*entry

func (q queue) Len() int {
	return len(q)
}

func (q queue) Less(i, j int) bool {
	return q[i].weight < q[j].weight
}

func (q queue) Swap(i, j int) {
	q[i],q[j] = q[j],q[i]
	q[i].index = i
	q[j].index = j
}

func (q *queue) Push(x interface{})  {
	n := len(*q)
	en := x.(*entry)
	en.index = n
	*q = append(*q,en)
}

func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	en := old[n-1]
	old[n-1] = nil
	en.index = -1
	*q = old[0:n-1]
	return en
}

func (q *queue) update(e *entry, value interface{}, weight int) {
	e.value = value
	e.weight = weight
	heap.Fix(q,e.index)
}
