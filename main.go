package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	node "github.com/divyanshu-parihar/go-gateway/node"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	slog.Info("SERVER", "status", "Started")

	ctx := context.Background()

	// Starting resouces required by the node
	cli, rdb, err := node.NodeResourceSetup(ctx, []string{"localhost:2379"}, "localhost:6379")
	if err != nil {
		slog.Error("Resource Setup", "status", "Failed to setup resources", "error", err)
		return // Handle error appropriately, e.g., log and exit
	}

	// Creating the node Object handling all the transactions
	node := node.NewNode(ctx, "node1", "us-east-1", node.ROUND_ROBIN, node.NodeMetaData{Description: "Node1"}, cli, rdb)
	node.Router.StartWatching(ctx, []string{"localhost:2379"})

	// routes for the api gateway
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		resp, err := cli.Get(req.Context(), "/gateway/routes", clientv3.WithPrefix())
		if err != nil {
			http.Error(res, "Failed to fetch routes from etcd", http.StatusInternalServerError)
			return
		}

		// 2. Unpack etcd's internal Key-Value pairs into a standard Go map or struct
		routes := make(map[string]json.RawMessage)
		for _, kv := range resp.Kvs {
			// kv.Key is []byte, cast to string. kv.Value is the JSON you stored.
			routes[string(kv.Key)] = kv.Value
		}
		fmt.Printf("Fetched Routes: %v\n", routes)

		// 3. Write the extracted data to the http.ResponseWriter
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)

		// CUSTOM ROUTER
		result := node.Allow()
		slog.Info("Request", "method", req.Method, "url", req.URL.Path, "Verdict", result)
		if err != nil && result == false {
			fmt.Fprintf(res, "Sorry")
			return
		}

		fmt.Fprintf(res, "Welcome to the server")
		return
	})
	//server blocker
	http.ListenAndServe(":8080", nil)
}
