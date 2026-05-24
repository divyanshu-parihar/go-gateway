package node

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type NodeMetaData struct {
	Description string
}
type Node struct {
	id       string
	region   string
	strategy TrafficStrategy
	metadata NodeMetaData
	configDb *clientv3.Client
	redisDb  *redis.Client
	Router   *Router
}

func NewNode(context context.Context, id, region string, strategy TrafficStrategy, metadata NodeMetaData, configDb *clientv3.Client, redisDb *redis.Client) *Node {
	Router := NewRouter(configDb)
	return &Node{
		id,
		region,
		strategy,
		metadata,
		configDb,
		redisDb,
		Router,
	}
}
func (node *Node) ChangeStrategy(newStrategy TrafficStrategy) {
	slog.Info("Node", "id", node.id, "status", "Changing strategy", "new_strategy", newStrategy)
	node.strategy = newStrategy
}
func (node *Node) Allow() bool {
	return true
}
