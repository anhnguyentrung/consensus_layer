package consensus

import (
	"consensus_layer/crypto"
	"consensus_layer/blockchain"
)

type RequestNewTerm struct {
	Term uint64
	NodeId blockchain.SHA256Type
	Signature crypto.Signature
}

type RequestVote struct {
	Term uint64
	CandidateId blockchain.SHA256Type
}

type RequestVoteResponse struct {
	Term uint64
	NodeId blockchain.SHA256Type
	GrantVote bool
}