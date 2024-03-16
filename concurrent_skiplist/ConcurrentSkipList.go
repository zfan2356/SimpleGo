package concurrent_skiplist

import "sync/atomic"

type ConcurrentSkipList struct {
	cap atomic.Int32
}
