package consensus

import (
	"consensus_layer/network"
	"fmt"
	"consensus_layer/crypto"
	"crypto/sha256"
)

type ElectionManager struct {
	role Role
	term uint64
	signer network.SignFunc
	address string
	producers []Producer
}

func NewElectionManager(signer network.SignFunc, address string) *ElectionManager {
	em := &ElectionManager{
		role: Follower,
		term: 0,
		signer: signer,
		address: address,
	}
	return em
}

// election manager inherit base manager interface
var _ network.BaseManager = (*ElectionManager)(nil)

func (em *ElectionManager) Receive(conn *network.Connection, message network.Message) {
	messageType := message.Header.Type
	switch messageType {
	case network.RequestNewTerm:
		fmt.Println("receive new term request")
	case network.RequestVote:
		fmt.Println("receive vote request")
	case network.RequestVoteResponse:
		fmt.Println("receive vote response")
	default:
		break
	}
}

func (em *ElectionManager) Send(conn *network.Connection, messageType network.MessageType) {
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

func (em *ElectionManager) sendNewTermRequest(conn *network.Connection) {
	newTerm := RequestNewTerm{
		em.term + 1,
		em.address,
		crypto.Signature{},
	}
	buf, _ := network.MarshalBinary(newTerm)
	hash := sha256.Sum256(buf)
	sig := em.signer(hash)
	newTerm.Signature = sig
	conn.Send(newTerm)
}
