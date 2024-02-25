# SimpleGo
> Learn and easily implement some well-known Go project on github


### 一.GeeCache (FINISHED)
基于GroupCache源码实现的简易版本的分布式缓存，功能如下
+ 单机缓存和基于 HTTP 的分布式缓存
+ 最近最少访问(Least Recently Used, LRU) 缓存策略
+ 使用 Go 锁机制防止缓存击穿
+ 使用一致性哈希选择节点，实现负载均衡
+ 使用 protobuf 优化节点间二进制通信

### 二. GeeWebFrame