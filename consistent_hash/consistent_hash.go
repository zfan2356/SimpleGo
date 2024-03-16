package consistent_hash

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type ConsistentHash struct {
	// hash环，真正存储数据的地方
	hashRing HashRing
	// 数据迁移器
	migrator Migrator
	// 哈希散列函数
	encryptor Encryptor
	// 用户自定义配置
	opts ConsistentHashOptions
}

func NewConsistentHash(hashRing HashRing, encryptor Encryptor,
	migrator Migrator, opts ...ConsistentHashOption) *ConsistentHash {
	ch := ConsistentHash{
		hashRing:  hashRing,
		migrator:  migrator,
		encryptor: encryptor,
	}
	for _, opt := range opts {
		opt(&ch.opts)
	}
	repair(&ch.opts)
	return &ch
}
func (c *ConsistentHash) getVaildWeight(w int) int {
	if w <= 0 {
		return 1
	}
	if w >= 10 {
		return 10
	}
	return w
}
func (c *ConsistentHash) getRawNodeKey(nodeID string, index int) string {
	return fmt.Sprintf("%s_%d", nodeID, index)
}
func (c *ConsistentHash) getNodeID(rawNodeKey string) string {
	index := strings.LastIndex(rawNodeKey, "_")
	return rawNodeKey[:index]
}
func (c *ConsistentHash) batchExecuteMigrator(migrateTasks []func()) {
	// 执行所有的数据迁移任务
	var wg sync.WaitGroup
	for _, migrateTask := range migrateTasks {
		//migrateTask := migrateTask
		wg.Add(1)
		migrateTask := migrateTask
		go func() {
			defer func() {
				if err := recover(); err != nil {
				}
				wg.Done()
			}()
			migrateTask()
		}()
	}
	wg.Wait()
}

// 添加节点需要触发数据迁移
func (c *ConsistentHash) AddNode(ctx context.Context, nodeID string, weight int) error {
	if err := c.hashRing.Lock(ctx, c.opts.lockExpireSeconds); err != nil {
		return err
	}
	defer func() {
		_ = c.hashRing.Unlock(ctx)
	}()

	// 判断节点是否重复
	nodes, err := c.hashRing.Nodes(ctx)
	if err != nil {
		return err
	}
	for node := range nodes {
		if node == nodeID {
			return errors.New("repeat node")
		}
	}

	replicas := c.getVaildWeight(weight) * c.opts.replicas
	if err := c.hashRing.AddNodeToReplica(ctx, nodeID, replicas); err != nil {
		return err
	}

	var migrateTasks []func()
	for i := 0; i < replicas; i++ {
		nodeKey := c.getRawNodeKey(nodeID, i)
		virtualScore := c.encryptor.Encrypt(nodeKey)

		if err := c.hashRing.Add(ctx, virtualScore, nodeKey); err != nil {
			return err
		}

		from, to, datas, err := c.migrateIn(ctx, virtualScore, nodeID)
		if err != nil {
			return err
		}
		if len(datas) == 0 {
			continue
		}
		migrateTasks = append(migrateTasks, func() {
			_ = c.migrator(ctx, datas, from, to)
		})
	}

	c.batchExecuteMigrator(migrateTasks)
	return nil
}

func (c *ConsistentHash) RemoveNode(ctx context.Context, nodeID string) error {
	if err := c.hashRing.Lock(ctx, c.opts.lockExpireSeconds); err != nil {
		return err
	}
	defer func() {
		_ = c.hashRing.Unlock(ctx)
	}()

	nodes, err := c.hashRing.Nodes(ctx)
	if err != nil {
		return err
	}

	var (
		nodeExist bool
		replicas  int
	)
	for node, _replicas := range nodes {
		if node == nodeID {
			nodeExist = true
			replicas = _replicas
			break
		}
	}

	if !nodeExist {
		return errors.New("invalid node id")
	}
	if err = c.hashRing.DeleteNodeToReplica(ctx, nodeID); err != nil {
		return err
	}

	var migrateTasks []func()
	for i := 0; i < replicas; i++ {
		virtualScore := c.encryptor.Encrypt(c.getRawNodeKey(nodeID, i))
		from, to, datas, err := c.migrateOut(ctx, virtualScore, nodeID)
		if err != nil {
			return err
		}
		nodeKey := c.getRawNodeKey(nodeID, i)
		if err = c.hashRing.Rem(ctx, virtualScore, nodeKey); err != nil {
			return err
		}
		if len(datas) == 0 {
			continue
		}

		migrateTasks = append(migrateTasks, func() {
			_ = c.migrator(ctx, datas, from, to)
		})
	}
	c.batchExecuteMigrator(migrateTasks)
	return nil
}

func (c *ConsistentHash) GetNode(ctx context.Context, dataKey string) (string, error) {
	if err := c.hashRing.Lock(ctx, c.opts.lockExpireSeconds); err != nil {
		return "", err
	}
	defer func() {
		_ = c.hashRing.Unlock(ctx)
	}()

	dataScore := c.encryptor.Encrypt(dataKey)
	ceilingScore, err := c.hashRing.Ceiling(ctx, dataScore)
	if err != nil {
		return "", err
	}
	if ceilingScore == -1 {
		return "", errors.New("no node available")
	}
	nodes, err := c.hashRing.Node(ctx, ceilingScore)
	if err != nil {
		return "", err
	}
	if len(nodes) == 0 {
		return "", errors.New("no node available with empty score")
	}
	if err = c.hashRing.AddNodeToDataKeys(ctx, c.getNodeID(nodes[0]), map[string]struct{}{
		dataKey: {},
	}); err != nil {
		return "", err
	}
	return nodes[0], nil
}
