package consistent_hash

import (
	"context"
	"errors"
)

// 用户所注册的闭包函数, 用来执行数据迁移操作
type Migrator func(ctx context.Context, dataKeys map[string]struct{}, from, to string) error

func (c *ConsistentHash) migrateIn(ctx context.Context, virtualScore int32,
	nodeID string) (from, to string, data map[string]struct{}, _err error) {
	if c.migrator == nil {
		return
	}
	nodes, err := c.hashRing.Node(ctx, virtualScore)
	if err != nil {
		_err = err
		return
	}
	if len(nodes) > 1 {
		return
	}

	// 找到下一个节点
	nextScore, err := c.hashRing.Ceiling(ctx, virtualScore)
	if err != nil {
		_err = err
		return
	}
	if nextScore == -1 {
		return
	}

	// 找到上一个节点,确定边界, 注意,其实这个时候如果下一个节点可以找到,
	// 上一个节点是必定可以找到的, 虽然可能和下一个节点重合
	lastScore, err := c.hashRing.Floor(ctx, virtualScore)
	if err != nil {
		_err = err
		return
	}

	nextNodes, err := c.hashRing.Node(ctx, nextScore)
	if err != nil {
		_err = err
		return
	}
	if len(nextNodes) == 0 {
		return
	}

	dataKeys, err := c.hashRing.DataKeys(ctx, c.getNodeID(nextNodes[0]))
	if err != nil {
		_err = err
		return
	}

	datas := make(map[string]struct{})
	for k := range dataKeys {
		virS := c.encryptor.Encrypt(k)
		if check(lastScore, virtualScore, virS) {
			datas[k] = struct{}{}
		}
	}

	if err = c.hashRing.DeleteNodeToDataKeys(ctx, c.getNodeID(nextNodes[0]), datas); err != nil {
		return "", "", nil, err
	}
	if err = c.hashRing.AddNodeToDataKeys(ctx, nodeID, datas); err != nil {
		return "", "", nil, err
	}
	return c.getNodeID(nextNodes[0]), nodeID, datas, nil
}

func check(l, r, x int32) bool {
	if l <= r {
		return l <= x && x <= r
	}
	return x >= l || x <= r
}

func (c *ConsistentHash) migrateOut(ctx context.Context, virtualScore int32,
	nodeID string) (from, to string, datas map[string]struct{}, err error) {
	if c.migrator == nil {
		return
	}
	defer func() {
		if err != nil {
			return
		}
		if to == "" || len(datas) == 0 {
			return
		}
		if err = c.hashRing.DeleteNodeToDataKeys(ctx, nodeID, datas); err != nil {
			return
		}
		err = c.hashRing.AddNodeToDataKeys(ctx, to, datas)
	}()
	from = nodeID

	nodes, _err := c.hashRing.Node(ctx, virtualScore)
	if _err != nil {
		err = _err
		return
	}
	if len(nodes) == 0 {
		return
	}
	if c.getNodeID(nodes[0]) != nodeID {
		return
	}

	var all map[string]struct{}
	if all, err = c.hashRing.DataKeys(ctx, nodeID); err != nil {
		return
	}
	if len(all) == 0 {
		return
	}

	lastScore, _err := c.hashRing.Floor(ctx, virtualScore)
	if _err != nil {
		err = _err
		return
	}
	ok := false
	if lastScore == -1 {
		if len(nodes) == 1 {
			err = errors.New("only one node")
			return
		}
		ok = true
	}
	datas = make(map[string]struct{})
	// 虽然全局只有一个虚拟节点,但是虚拟节点内包含很多节点,可以委托给下一个节点
	if ok {
		for k := range all {
			datas[k] = struct{}{}
		}
		to = c.getNodeID(nodes[1])
		return
	}
	for k := range all {
		virS := c.encryptor.Encrypt(k)
		if check(lastScore, virtualScore, virS) {
			datas[k] = struct{}{}
		}
	}

	// 已经特判掉只有一个节点的情况
	nextScore, _err := c.hashRing.Ceiling(ctx, virtualScore)
	if _err != nil {
		err = _err
		return
	}
	nextnodes, _err := c.hashRing.Node(ctx, nextScore)
	if _err != nil {
		err = _err
		return
	}
	to = c.getNodeID(nextnodes[0])
	return
}
