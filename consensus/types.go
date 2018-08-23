package consensus

import (
	"consensus_layer/crypto"
	"consensus_layer/blockchain"
)

type Role uint8

const (
	Follower Role = iota
	Candidate
	Leader
)

type RequestNewTerm struct {
	Term uint64
	Sender string
	Signature crypto.Signature
}

type RequestVote struct {
	Term uint64
	CandidateId blockchain.SHA256Type
}

type RequestVoteResponse struct {
	Term uint64
	Sender string
	GrantVote bool
}