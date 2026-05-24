# Distributed API Gateway

## Features (Functional Requirements)

[] Dynamic Route Resolution: The gateway must parse incoming HTTP/REST requests and map them to downstream backend services based on URL paths
[] Load Balancing Strategies: Implement multiple algorithms to distribute traffic across backend instances. You'll need at least Round Robin (for uniform workloads)
[] Hot-Reloadable Configuration: The routing table must be updatable dynamically (e.g., reading from a YAML file or a key-value store like etcd/Consul) without dropping active connections or requiring a server restart.
