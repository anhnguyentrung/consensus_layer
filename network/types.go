package network

import "consensus_layer/crypto"

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

type SHA256Type [32]byte

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
	ChainId                 SHA256Type
	NodeId                  SHA256Type
	Key                     crypto.PublicKey
	P2PAddress              string
	LastCommitBlockHeight   uint32
	LastCommitBlockId  		SHA256Type
	TopBlockHeight          uint32
	TopBlockId              SHA256Type
}

type HandshakePacket struct {
	Info HandshakeInfo
	Sign crypto.Signature
}
