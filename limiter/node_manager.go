package limiter

import (
	"sync"
)

type Node struct {
	ID     string  `json:"id"`
	Weight float64 `json:"weight"`
}

type NodeManager struct {
	mu    sync.Mutex
	nodes map[string]Node
	self  string
}

func NewNodeManager(selfID string) *NodeManager {
	return &NodeManager{
		nodes: make(map[string]Node),
		self:  selfID,
	}
}

func (nm *NodeManager) UpdateNodes(newNodes []Node) {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	nm.nodes = make(map[string]Node)
	for _, n := range newNodes {
		nm.nodes[n.ID] = n
	}
}

func (nm *NodeManager) CalcSelfRate(globalQPS float64) float64 {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	var totalWeight float64
	for _, n := range nm.nodes {
		totalWeight += n.Weight
	}
	selfNode, ok := nm.nodes[nm.self]
	if !ok || totalWeight == 0 {
		return 0
	}
	return globalQPS * (selfNode.Weight / totalWeight)
}
