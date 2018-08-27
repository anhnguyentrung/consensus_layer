package consensus

import "consensus_layer/network"

type ElectionManager struct {
	Role Role
	Term uint64
	NewTerm chan *RequestNewTerm
	VoteResponse chan *RequestVoteResponse
}

func NewElectionManager() *ElectionManager {
	em := &ElectionManager{
		Role: Follower,
		Term: 1,
	}
	return em
}

// election manager inherit base manager interface
var _ network.BaseManager = (*ElectionManager)(nil)

func (em *ElectionManager) Receive(conn *network.Connection, message network.Message) {

}
