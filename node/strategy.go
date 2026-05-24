package node

type TrafficStrategy string

// strategy enum
const (
	ROUND_ROBIN TrafficStrategy = "ROUND_ROBIN"
)

type DistributionStrategy interface {
	allow()
}

