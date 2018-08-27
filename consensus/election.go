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
	voteCounter map[network.MessageType]TermVote
}

func NewElectionManager(signer network.SignFunc, address string) *ElectionManager {
	em := &ElectionManager{
		role: Follower,
		term: 0,
		signer: signer,
		address: address,
		voteCounter: make(map[network.MessageType]TermVote, 0),
	}
	em.voteCounter[network.RequestNewTerm] = make(TermVote, 0)
	return em
}

// election manager inherit base manager interface
var _ network.BaseManager = (*ElectionManager)(nil)

func (em *ElectionManager) Receive(conn *network.Connection, message network.Message) {
	messageType := message.Header.Type
	switch messageType {
	case network.RequestNewTerm:
		fmt.Println("receive new term request")
		newTerm := RequestNewTerm{}
		network.UnmarshalBinary(message.Payload, &newTerm)
		em.receivedNewTerm(newTerm)
	case network.RequestVote:
		fmt.Println("receive vote request")
	case network.RequestVoteResponse:
		fmt.Println("receive vote response")
	default:
		break
	}
}

func (em *ElectionManager) receivedNewTerm(conn *network.Connection, newTerm RequestNewTerm) {
	if em.role == Follower {
		producerIndex := int(newTerm.Term) % len(em.producers)
		if em.producers[producerIndex].Address == em.address {
			em.voteCounter[network.RequestNewTerm][newTerm.Term] += 1
		}
		if em.voteCounter[network.RequestNewTerm][newTerm.Term] > uint32(len(em.producers)) * 2/3 {
			em.becomeCandidate(newTerm.Term)
			//em.sendVoteRequest(conn)
		}
	}
}

func (em *ElectionManager) becomeCandidate(term uint64) {
	em.role = Candidate
	em.term = term
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
