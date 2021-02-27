package main

import (
	"limiter/limiter"
	"log"
	"sync"
	"time"
)

func main() {
	/* c := &limiter.Counter{}
	c.Set(3, time.Second)
	limit(c) */

	/* leakyBucket := &limiter.LeakyBucket{}
	leakyBucket.Set(2, 2)
	limit(leakyBucket) */

	tokenBucket := &limiter.TokenBucket{}
	tokenBucket.Set(1, 3)
	limit(tokenBucket)
}

func limit(limiter limiter.Limiter) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		log.Println("创建请求：",i)
		go func(i int) {
			defer wg.Done()
			if limiter.Allow() {
				log.Println("响应请求：",i)
			}
		}(i)
		time.Sleep(200 * time.Millisecond)
	}

	wg.Wait()
}