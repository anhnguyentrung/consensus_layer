package network

import (
	"net"
	"bufio"
)

type Connection struct {
	conn 			net.Conn
	connReader 		*bufio.Reader
	isOpen 			bool
	isSynchronizing bool
	peerAddress		string
}

func NewConnection(peerAddress string) *Connection {
	return &Connection{
		peerAddress:		peerAddress,
		isOpen:				false,
		isSynchronizing:	false,
	}
}

func (c *Connection) sendPacket(packet interface{}) error {
	err := error(nil)
	bytes := make([]byte, 0)
	messageType := Handshake
	switch packet.(type) {
	case HandshakePacket:
		bytes, err = marshalBinary(packet)
		if err != nil {
			return err
		}
	}
	message := Message{
		Header: MessageHeader{},
		Payload: bytes,
	}
	message.Header.Type = messageType
	message.Header.Length = uint32(len(bytes))
	data, err := marshalBinary(message)
	if err != nil {
		return err
	}
	c.conn.Write(data)
	return err
}

func (c *Connection) close()  {
	c.connReader = nil
	c.isOpen = false
	c.isSynchronizing = false
	c.conn.Close()
}