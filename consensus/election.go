package consensus

import "consensus_layer/network"

type ElectionManager struct {
	node *network.Node
	Role Role
	Term uint64
}

func NewElectionManager(node *network.Node) *ElectionManager {
	em := &ElectionManager{
		node: node,
		Role: Follower,
		Term: 1,
	}
	return em
}
