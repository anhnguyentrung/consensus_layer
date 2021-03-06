package consensus

import (
	"consensus_layer/crypto"
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
	Candidate string // address
	SignedNewTerms []RequestNewTerm
	Signature crypto.Signature
}

type GrantVote struct {
	Term uint64
	Sender string
	Signature crypto.Signature
}

type Producer struct {
	Address string
	PublicKey *crypto.PublicKey
}

type TermVote map[uint64]uint32 // [term]vote