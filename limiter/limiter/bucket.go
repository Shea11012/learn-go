package limiter

import (
	"math"
	"sync"
	"time"
)

// LeakyBucket 漏桶算法限定了流量的流出速率，所以最大的速率就是出水速率，不会出现突发流量。
// 在网络较好的时候，依然没有办法放开量。
type LeakyBucket struct {
	rate float64	// 固定每秒请求速率
	capacity float64 // 桶容量
	water float64 // 桶中当前水量
	lastLeakMs int64  // 桶上次漏水时间 ms

	lock sync.Mutex
}

func (lb *LeakyBucket) Allow() bool {
	lb.lock.Lock()
	defer lb.lock.Unlock()

	now := time.Now().UnixNano() / 1e6	// 获取当前时间毫秒
	eclipse := float64(now - lb.lastLeakMs) * lb.rate / 1000 // 计算当前可出水数量
	lb.water = lb.water - eclipse // 桶内剩余水量
	lb.water = math.Max(float64(0),lb.water)

	lb.lastLeakMs = now

	if (lb.water + 1) < lb.capacity { // 桶未满
		lb.water++
		return true
	} else {
		return false	// 桶满了
	}
}

func (lb *LeakyBucket) Set(rate,capactiy float64) {
	lb.rate = rate
	lb.capacity = capactiy
	lb.water = 0
	lb.lastLeakMs = time.Now().UnixNano() / 1e6
}

// TokenBucket 令牌桶算法限制的是平均流入速率，允许突发请求，会按照固定的速率生成令牌放入桶中
// 最多存放N个令牌，多余的会被丢弃，不足的会补上。允许一定程度上的突发流量
type TokenBucket struct {
	rate int64	// 固定的token放入速率,r/s
	capacity int64 // 桶容量
	tokens int64 // 桶中当前token数量
	lastTokenSec int64 // 桶上次放token的时间戳 s

	lock sync.Mutex
}

func (tb *TokenBucket) Allow() bool {
	tb.lock.Lock()
	defer tb.lock.Unlock()

	now := time.Now().Unix()
	tb.tokens = tb.tokens + (now - tb.lastTokenSec) * tb.rate // 计算出单位时间内的token数量，将token放入桶中

	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	tb.lastTokenSec = now
	if tb.tokens > 0 {	// 当前桶中还有token剩余
		tb.tokens--
		return true
	} else {
		return false
	}
}

func (tb *TokenBucket) Set(rate,capacity int64) {
	tb.rate = rate
	tb.capacity = capacity
	tb.tokens = 0
	tb.lastTokenSec = time.Now().Unix()
}