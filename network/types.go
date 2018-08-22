package network

type MessageType byte
const (
	Handshake MessageType = iota
	Notice
	Request
	Block
)

type MessageHeader struct {
	Type 	MessageType // 1 byte
	Length  uint32  	// 4 bytes
}
type Message struct {
	Header MessageHeader
	Payload []byte
}
