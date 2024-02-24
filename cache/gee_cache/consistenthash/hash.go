package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash Function
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

func NewMap(replicas int, fn Hash) *Map {
	mp := &Map{
		hash:     fn,
		replicas: replicas,
		keys:     nil,
		hashMap:  make(map[int]string),
	}
	if mp.hash == nil {
		mp.hash = crc32.ChecksumIEEE
	}
	return mp
}

// Add to add peers to the map
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get using binary search to find the peer to get cache
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
