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
	newTerms []RequestNewTerm
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
		em.receivedNewTerm(conn, newTerm)
	case network.RequestVote:
		fmt.Println("receive vote request")
		voteRequest := RequestVote{}
		network.UnmarshalBinary(message.Payload, &voteRequest)
		em.receivedVoteRequest(conn, voteRequest)
	case network.GrantVote:
		fmt.Println("receive vote response")
	default:
		break
	}
}

func (em *ElectionManager) receivedNewTerm(conn *network.Connection, newTerm RequestNewTerm) {
	if em.role == Follower {
		// signature is invalid
		if !em.verifyNewTerm(newTerm) {
			return
		}
		em.newTerms = append(em.newTerms, newTerm)
		producerIndex := int(newTerm.Term) % len(em.producers)
		if em.producers[producerIndex].Address == em.address {
			em.voteCounter[network.RequestNewTerm][newTerm.Term] += 1
		}
		if em.voteCounter[network.RequestNewTerm][newTerm.Term] > uint32(len(em.producers)) * 2/3 {
			em.becomeCandidate(newTerm.Term)
			em.sendVoteRequest(conn)
		}
	}
}

func (em *ElectionManager) receivedVoteRequest(conn *network.Connection, voteRequest RequestVote) {
	if em.term >= voteRequest.Term {
		fmt.Println("the term of vote request should be higher than the local term")
		return
	} else {
		if !em.validateVoteRequest(voteRequest) {
			fmt.Println("vote request is invalid")
			return
		}
		em.role = Follower
		em.sendGrantVote(conn)
	}
}

func (em *ElectionManager) validateVoteRequest(voteRequest RequestVote) bool {
	var candidatePub *crypto.PublicKey = nil
	for _, p := range em.producers {
		if p.Address == voteRequest.Candidate {
			candidatePub = p.PublicKey
			break
		}
	}
	if candidatePub == nil {
		return false
	}
	// verify signature of candidate
	if !em.verifySignatureOfVoteRequest(voteRequest, candidatePub) {
		return false
	}
	// verify signed new term requests that the candidate received
	for _, newTerm := range voteRequest.SignedNewTerms {
		if !em.verifyNewTerm(newTerm) {
			return false
		}
	}
	return true
}

func (em *ElectionManager) verifySignatureOfVoteRequest(voteRequest RequestVote, candidatePub *crypto.PublicKey) bool {
	voteRequestWithoutSignature := RequestVote{
		voteRequest.Term,
		voteRequest.Candidate,
		voteRequest.SignedNewTerms,
		crypto.Signature{},
	}
	buf, _ := network.MarshalBinary(voteRequestWithoutSignature)
	hash := sha256.Sum256(buf)
	return voteRequest.Signature.Verify(*candidatePub, hash[:])
}

func (em *ElectionManager) verifyNewTerm(newTerm RequestNewTerm) bool {
	var senderPub *crypto.PublicKey = nil
	for _, p := range em.producers {
		if p.Address == newTerm.Sender {
			senderPub = p.PublicKey
			break
		}
	}
	if senderPub == nil {
		return false
	}
	newTermWithoutSignature := RequestNewTerm{
		newTerm.Term,
		newTerm.Sender,
		crypto.Signature{},
	}
	buf, _ := network.MarshalBinary(newTermWithoutSignature)
	hash := sha256.Sum256(buf)
	return newTerm.Signature.Verify(*senderPub, hash[:])
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
	case network.GrantVote:
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

func (em *ElectionManager) sendVoteRequest(conn *network.Connection) {
	signedNewTerms := make([]RequestNewTerm, 0)
	for _, newTerm := range em.newTerms {
		if newTerm.Term == em.term {
			signedNewTerms = append(signedNewTerms, newTerm)
		}
	}
	requestVote := RequestVote{
		Term: em.term,
		Candidate: em.address,
		SignedNewTerms:	signedNewTerms,
		Signature: crypto.Signature{},
	}
	buf, _ := network.MarshalBinary(requestVote)
	hash := sha256.Sum256(buf)
	sig := em.signer(hash)
	requestVote.Signature = sig
	conn.Send(requestVote)
}

func (em *ElectionManager) sendGrantVote(conn *network.Connection) {
	grantVote := GrantVote{
		em.term,
		em.address,
		crypto.Signature{},
	}
	buf, _ := network.MarshalBinary(grantVote)
	hash := sha256.Sum256(buf)
	sig := em.signer(hash)
	grantVote.Signature = sig
	conn.Send(grantVote)
}
