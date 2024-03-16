package consistent_hash

import "context"

type HashRing interface {
	// 保证并发安全的加锁
	Lock(ctx context.Context, expireSeconds int) error
	Unlock(ctx context.Context) error
	// 添加删除节点
	Add(ctx context.Context, virtualScore int32, nodeID string) error
	Rem(ctx context.Context, virtualScore int32, nodeID string) error
	// 找到前驱 or 后继 (< virtualScore) (> virtualScore)
	Ceiling(ctx context.Context, virtualScore int32) (int32, error)
	Floor(ctx context.Context, virtualScore int32) (int32, error)
	// 获取所有节点([id, 虚拟节点个数])
	Nodes(ctx context.Context) (map[string]int, error)
	// 获取某个值上的节点列表, 我们规定只有第一个节点是生效的
	Node(ctx context.Context, virtualScore int32) ([]string, error)
	// 增加和删除节点所对应的虚拟节点
	AddNodeToReplica(ctx context.Context, nodeID string, replicas int) error
	DeleteNodeToReplica(ctx context.Context, nodeID string) error
	// 查看某个节点对应的key的集合
	DataKeys(ctx context.Context, nodeID string) (map[string]struct{}, error)
	// 增加和删除节点所对应的key集合
	AddNodeToDataKeys(ctx context.Context, nodeID string, dataKeys map[string]struct{}) error
	DeleteNodeToDataKeys(ctx context.Context, nodeID string, dataKeys map[string]struct{}) error
}
