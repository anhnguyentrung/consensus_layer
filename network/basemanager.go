package network

type BaseManager interface {
	Send(conn *Connection, messageType MessageType)
	Receive(conn *Connection, message Message)
}
