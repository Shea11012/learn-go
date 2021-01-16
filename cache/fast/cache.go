package fast

type fastCache struct {
	shards []*cacheShard
	shardMask uint64
	hash fnv64a
}

func NewFastCache(maxEntries int,shardsNum int,onEvicted func(key string,value interface{})) *fastCache {
	fastCache := &fastCache{
		hash: newDefaultHasher(),
		shards: make([]*cacheShard,shardsNum),
		shardMask: uint64(shardsNum-1),
	}

	for i:=0;i<shardsNum;i++ {
		fastCache.shards[i] = newCacheShard(maxEntries,onEvicted)
	}

	return fastCache
}

func (f *fastCache) getShard(key string) *cacheShard {
	hashKey := f.hash.Sum64(key)
	return f.shards[hashKey&f.shardMask]	// key % n == key & (n-1),n为2的幂时，使用位运算开销更小
}

func (f *fastCache) Set(key string, value interface{}) {
	f.getShard(key).set(key,value)
}

func (f *fastCache) Get(key string) interface{} {
	return f.getShard(key).get(key)
}

func (f *fastCache) Del(key string) {
	f.getShard(key).del(key)
}

func (f *fastCache) Len() int {
	length := 0
	for _,shard := range f.shards {
		length += shard.Len()
	}

	return length
}



