package network

import (
	"net"
	"bufio"
	"strings"
	"fmt"
)

type Connection struct {
	conn 			net.Conn
	connReader 		*bufio.Reader
	isOpen 			bool
	isSynchronizing bool
	isOutgoing		bool
	onReceive		ReceiveFunc
	onFinish		FinishFunc
}

func newConnection() *Connection {
	return &Connection{
		isOpen:				false,
		isSynchronizing:	false,
		isOutgoing: 		false,
	}
}

func NewOutgoingConnection(remoteAddr string, onRecevie ReceiveFunc, onFinish FinishFunc) (*Connection, error) {
	if !strings.Contains(remoteAddr, ":") {
		return nil, fmt.Errorf("invalid peer address %s", remoteAddr)
	}
	conn, err := net.Dial(TCP, remoteAddr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	c := newConnection()
	c.conn = conn
	c.isOutgoing = true
	c.isOpen = true
	c.onReceive = onRecevie
	c.onFinish = onFinish
	return c, nil
}

func NewIncomingConnection(conn net.Conn, onRecevie ReceiveFunc, onFinish FinishFunc) *Connection {
	c := newConnection()
	c.conn = conn
	c.isOutgoing = false
	c.isOpen = true
	c.onReceive = onRecevie
	c.onFinish = onFinish
	return c
}

func (c *Connection) IsOutgoing() bool {
	return c.isOutgoing
}

func (c *Connection) IsAvailable() bool {
	return c.isOpen && !c.isSynchronizing
}

func (c *Connection) RemoteAddress() string {
	return c.conn.RemoteAddr().String()
}

func (c *Connection) LocalAddress() string {
	return c.conn.LocalAddr().String()
}

func (c *Connection) Send(packet interface{}) error {
	err := error(nil)
	bytes := make([]byte, 0)
	messageType := Handshake
	switch packet.(type) {
	case HandshakePacket:
		bytes, err = MarshalBinary(packet)
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
	data, err := MarshalBinary(message)
	if err != nil {
		return err
	}
	c.conn.Write(data)
	return err
}

func (c *Connection) readLoop() {
	c.connReader = bufio.NewReader(c.conn)
	for {
		message := Message{
			Header:		MessageHeader{},
			Payload: 	make([]byte, 0),
		}
		UnmarshalBinaryMessage(c.connReader, &message)
		receiveMessage := ReceiveMessage{
			Conn: 		c,
			Message: 	message,
		}
		c.onReceive(receiveMessage)
	}
	c.onFinish(c)
}

func (c *Connection) Start() {
	go c.readLoop()
}

func (c *Connection) Close()  {
	c.connReader = nil
	c.isOpen = false
	c.isSynchronizing = false
	c.conn.Close()
}