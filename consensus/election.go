package consensus

import (
	"consensus_layer/network"
	"fmt"
)

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
	messageType := message.Header.Type
	switch messageType {
	case network.RequestNewTerm:
		fmt.Println("request new term")
	case network.RequestVote:
		fmt.Println("request vote")
	case network.RequestVoteResponse:
		fmt.Println("request vote response")
	default:
		break
	}
}
