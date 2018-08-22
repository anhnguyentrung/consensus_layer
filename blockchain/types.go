package blockchain

import (
	"consensus_layer/crypto"
	"time"
)

type SHA256Type [32]byte

type CommitType byte

const (
	PreCommitment CommitType = iota
	Commitment
)

type BlockHeader struct {
	Id SHA256Type
	Height uint64
	PreviousId SHA256Type
	Producer string
	Timestamp time.Time
}

type SignedHeader struct {
	Header BlockHeader
	Signature crypto.Signature
}

type SignedBlock struct {
	signedHeader SignedHeader
	//transactions []Transaction
}

type Commit struct {
	Type CommitType
	BlockId SHA256Type
	Committer string
	Timestamp time.Time
	Signature crypto.Signature
}
