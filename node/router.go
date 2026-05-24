package node

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Route struct {
	path        string
	destination string
	prefix      Prefix
}

type Router struct {
	routes   map[string]Route
	configDb *clientv3.Client
	mu       sync.RWMutex
}

func (r *Router) AddRoute(ctx context.Context, configDb *clientv3.Client, path string, destination string, metadata interface{}) error {
	_, err := configDb.Put(ctx, path, destination, clientv3.WithPrefix())
	return err
}

func (r *Router) StartWatching(ctx context.Context, endpoints []string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}

	// 1. Initial Load: Fetch all existing routes on startup
	resp, err := cli.Get(ctx, string(RoutePrefix), clientv3.WithPrefix())
	if err != nil {
		slog.Error("failed to load initial routes: ", "error", err)
		return err
	}

	r.mu.Lock()
	for _, kv := range resp.Kvs {
		var route Route
		if err := json.Unmarshal(kv.Value, &route); err == nil {
			r.routes[string(kv.Key)] = route
		}
	}
	r.mu.Unlock()
	slog.Info("Loaded initial routes", "count", len(resp.Kvs))
	// 2. Start the Watcher in a background goroutine
	go func() {
		watchChan := cli.Watch(ctx, string(RoutePrefix), clientv3.WithPrefix())

		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				key := string(event.Kv.Key)

				r.mu.Lock()
				switch event.Type {
				case clientv3.EventTypePut:
					// A route was added or updated
					var route Route
					if err := json.Unmarshal(event.Kv.Value, &route); err == nil {
						r.routes[key] = route
						slog.Info("Route updated: %s -> %s", key, route.destination)
					}
				case clientv3.EventTypeDelete:
					// A route was deleted
					delete(r.routes, key)
					slog.Info("Route deleted: %s", key)
				}
				r.mu.Unlock()
			}
		}
	}()
	return nil
}
func NewRouter(configDb *clientv3.Client) *Router {
	return &Router{
		routes:   map[string]Route{},
		configDb: configDb,
		mu:       sync.RWMutex{},
	}
}
