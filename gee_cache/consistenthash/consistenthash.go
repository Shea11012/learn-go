package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func([]byte) uint32

type Map struct {
	hash     Hash	// 哈希函数
	replicas int	// 虚拟节点倍数
	keys     []int	// 哈希环
	hashMap  map[int]string	// 虚拟节点与真实节点映射，键是虚拟节点，值是真实节点
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// Add 添加真实节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 每个真实节点创建 replicas 个虚拟节点
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}

	// 对环上的哈希值进行排序
	sort.Ints(m.keys)
}

// Get 获取对应key的真实节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// 根据key 的hash值，在环中查询索引
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 因为是环所以可以取余，找到真实节点
	return m.hashMap[m.keys[idx % len(m.keys)]]
}
