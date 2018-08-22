package network

import (
	"net"
	"bufio"
	"fmt"
)

type Connection struct {
	conn net.Conn
	connReader *bufio.Reader
	isOpen bool
	isSynchronized bool
}

func (c *Connection) sendMessage(message Message) {
	buf, err := marshalBinary(message.Content)
	if err != nil {
		fmt.Println("Encode content: ", err)
		return err
	}
	message.Header.Length = uint32(len(buf))
	data, err := marshalBinary(message)
	c.Conn.Write(data)
	return err
}