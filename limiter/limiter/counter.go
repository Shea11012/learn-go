package limiter

import (
	"sync"
	"time"
)

// Counter 计数器限流会有一个问题，例如某一接口每秒最多允许访问200次，如果有用户在59秒的最后几毫秒瞬间发送200次
// 当59秒结束后 count 清零，下一秒又发送了200次请求。这意味用户在一秒内发送了2倍的请求数，这是计数器限流的缺陷
type Counter struct {
	rate int	// 周期内允许的最大请求数
	begin time.Time	// 开始时间
	cycle time.Duration // 计数周期
	count int	// 周期内累计请求数
	lock sync.Mutex
}

func (c *Counter) Allow() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.count == c.rate-1 {	// 当前累计请求数等于最大频次时，判断当前请求时间是否在设置的时间周期内
		now := time.Now()
		if now.Sub(c.begin) >= c.cycle {	// 当前请求在允许的请求周期内
			// 请求成功，重置请求数和时间
			c.Reset(now)
			return true
		} else {
			return false
		}
	} else {
		c.count++
		return true
	}
}

func (c *Counter) Set(r int,cycle time.Duration) {
	c.rate = r
	c.begin = time.Now()
	c.cycle = cycle
}

func (c *Counter) Reset(t time.Time) {
	c.begin = t
	c.count = 0
}