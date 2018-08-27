package network

type BaseManager interface {
	Receive(conn *Connection, message Message)
}
