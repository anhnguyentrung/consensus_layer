package network

import (
	"sync"
	"consensus_layer/crypto"
	"net"
	"fmt"
	"bufio"
	"time"
	"crypto/sha256"
	"consensus_layer/blockchain"
	"strings"
	"consensus_layer/consensus"
)

type receiveMessage struct {
	conn 	*Connection
	message Message
}

type receiveBlock struct {
	conn *Connection
}

type Node struct {
	id 					blockchain.SHA256Type
	chainId 			blockchain.SHA256Type
	address 			string
	targets				[]string // addresses of specific peers that this node try to connect
	keyPairs 			map[string]*crypto.PrivateKey
	conns 				map[string]*Connection
	network 			NetworkType
	version 			uint16
	newConn 			chan *Connection // trigger when a connection is accepted
	doneConn 			chan *Connection // trigger when a connection is disconnected
	newMessage 			chan *receiveMessage // trigger when received a message
	receiveBlockQueue 	[]receiveBlock
	managers			map[string]BaseManager
	mutex 				sync.Mutex
}

func NewNode(address string, outbounds []string) *Node {
	node := &Node {
		address: address,
		targets: outbounds,
		keyPairs: make(map[string]*crypto.PrivateKey,0),
		conns: make(map[string]*Connection,0),
		newConn: make(chan *Connection),
		doneConn: make(chan *Connection),
		newMessage: make(chan *receiveMessage),
	}
	electionManager := consensus.NewElectionManager()
	node.addManager(electionManager, ElectionManager)
	return node
}

func (node *Node) addManager(manager BaseManager, id string) {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	node.managers[id] = manager
}

// connect to specific remote peers
func (node *Node) connectsToTargets() {
	fmt.Println("connecting to peer ...")
	for _, addr := range node.targets {
		c := NewConnection(addr)
		err := node.dial(c)
		if err != nil {
			fmt.Println("connecting to peer: ", err)
		}
	}
}

func (node *Node) dial(c *Connection) error {
	if !strings.Contains(c.peerAddress, ":") {
		return fmt.Errorf("invalid peer address %s", c.peerAddress)
	}
	conn, err := net.Dial(TCP, c.peerAddress)
	if err != nil {
		fmt.Println(err)
		return err
	}
	c.conn = conn
	node.newConn <- c
	return nil
}

// listen from remote peers
func (node *Node) listen() error {
	listener, err := net.Listen(TCP, node.address)
	if err != nil {
		return err
	}
	go func() {
		for {
			inboundConn, err := listener.Accept()
			if err != nil {
				panic(err)
			}
			c := NewConnection("")
			c.conn = inboundConn
			node.newConn <- c
		}
	}()
	for {
		select {
		case connection := <-node.newConn:
			fmt.Println("accepted new client from address ", connection.conn.RemoteAddr().String())
			node.addConnection(connection)
			go node.handleConnection(connection)
		case receiveMessage := <-node.newMessage:
			go node.handleMessage(receiveMessage)
		case doneConnection := <-node.doneConn:
			fmt.Println("disconnected client from address ", doneConnection.conn.RemoteAddr().String())
			node.removeConnection(doneConnection)
		}
	}
	return nil
}

func (node *Node) addConnection(c *Connection) {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	if c.peerAddress == "" {
		node.conns[c.conn.RemoteAddr().String()] = c
	} else {
		node.conns[c.peerAddress] = c
	}
	c.isOpen = true
}

func (node *Node) sendHandshake(c *Connection) {
	if c.peerAddress != "" {
		handshake := node.newHandshakePacket()
		c.sendPacket(handshake)
	}
}

func (node *Node) removeConnection(c *Connection) {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	c.isOpen = false
	if c.peerAddress == "" {
		delete(node.conns, c.conn.RemoteAddr().String())
	} else {
		delete(node.conns, c.peerAddress)
	}
}

func (node *Node) handleConnection(c *Connection) {
	c.connReader = bufio.NewReader(c.conn)
	for {
		message := Message{
			Header:		MessageHeader{},
			Payload: 	make([]byte, 0),
		}
		unmarshalBinaryMessage(c.connReader, &message)
		receiveMessage := &receiveMessage{
			conn: 		c,
			message: 	message,
		}
		node.newMessage <- receiveMessage
	}
	node.doneConn <- c
}

func (node *Node) handleMessage(receiveMessage *receiveMessage) {
	messageType := receiveMessage.message.Header.Type
	c := receiveMessage.conn
	switch messageType {
	case Handshake:
		fmt.Println("handshake")
		handshake := HandshakePacket{}
		unmarshalBinary(receiveMessage.message.Payload, &handshake)
		node.handleHandshake(c, handshake)
	case Notice:
	case Request:
	case Block:
	case RequestNewTerm, RequestVote, RequestVoteResponse:
		node.managers[ElectionManager].Receive(c, receiveMessage.message)
	}
}

func (node *Node) handleHandshake(c *Connection, handshake HandshakePacket) {
}

func (node *Node) newHandshakePacket() HandshakePacket {
	publicKey := &crypto.PublicKey{}
	if len(node.keyPairs) > 0 {
		for k := range node.keyPairs {
			publicKey, _ = crypto.NewPublicKey(k)
			break
		}
	}
	info := HandshakeInfo{
		Network:				TestNet,
		Version:				1,
		ChainId: 				node.chainId,
		NodeId: 				node.id,
		Key: 					*publicKey,
		originAddress: 			node.address,
		LastCommitBlockHeight: 	0,
		LastCommitBlockId: 		blockchain.SHA256Type{},
		TopBlockHeight: 		0,
		TopBlockId:				blockchain.SHA256Type{},
		Timestamp:				time.Now(),
	}
	buf, _ := marshalBinary(info)
	hash := sha256.Sum256(buf)
	sign := crypto.Signature{}
	if privateKey, ok := node.keyPairs[publicKey.String()]; ok {
		sign, _ = privateKey.Sign(hash[:])
	}
	return HandshakePacket{
		Info: info,
		Sign: sign,
	}
}