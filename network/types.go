package network

import (
	"consensus_layer/crypto"
	"consensus_layer/blockchain"
	"time"
)

type MessageType byte
const (
	Handshake MessageType = iota
	Notice
	Request
	Block
)

type NetworkType byte
const (
	TestNet NetworkType = iota
	MainNet
)

type MessageHeader struct {
	Type 	MessageType // 1 byte
	Length  uint32  	// 4 bytes
}
type Message struct {
	Header MessageHeader
	Payload []byte
}

type HandshakeInfo struct {
	Network					NetworkType
	Version					uint16
	ChainId                 blockchain.SHA256Type
	NodeId                  blockchain.SHA256Type
	Key                     crypto.PublicKey
	originAddress           string
	LastCommitBlockHeight   uint32
	LastCommitBlockId  		blockchain.SHA256Type
	TopBlockHeight          uint32
	TopBlockId              blockchain.SHA256Type
	Timestamp 				time.Time
}

type HandshakePacket struct {
	Info HandshakeInfo
	Sign crypto.Signature
}
