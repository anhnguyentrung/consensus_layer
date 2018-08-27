package node

import (
	"sync"
	"consensus_layer/crypto"
	"net"
	"fmt"
	"time"
	"crypto/sha256"
	"consensus_layer/blockchain"
	"consensus_layer/consensus"
	"consensus_layer/network"
)

//type receiveBlock struct {
//	conn *network.Connection
//}

type keyPair struct {
	publicKey *crypto.PublicKey
	privateKey *crypto.PrivateKey
}

type Node struct {
	id 					blockchain.SHA256Type
	chainId 			blockchain.SHA256Type
	p2pAddress 			string
	targets				[]string // addresses of specific peers that this node try to connect
	keyPair				keyPair
	conns 				map[string]*network.Connection
	network 			network.NetworkType
	version 			uint16
	newConn 			chan *network.Connection // trigger when a connection is accepted
	doneConn 			chan *network.Connection // trigger when a connection is disconnected
	//receiveBlockQueue 	[]receiveBlock
	managers			map[string]network.BaseManager
	walletAddress		string
	mutex 				sync.Mutex
}

func NewNode(p2pAddress string, outbounds []string) *Node {
	node := &Node {
		p2pAddress: p2pAddress,
		targets: outbounds,
		//keyPairs: make(map[string]*crypto.PrivateKey, 0),
		conns: make(map[string]*network.Connection, 0),
		newConn: make(chan *network.Connection),
		doneConn: make(chan *network.Connection),
		managers: make(map[string]network.BaseManager, 0),
		//newMessage: make(chan *network.ReceiveMessage),
	}
	electionManager := consensus.NewElectionManager(node.Signer, node.walletAddress)
	node.addManager(electionManager, network.ElectionManager)
	return node
}

func (node *Node) Start() {
	go node.connectsToTargets()
	go node.listen()
}

func (node *Node) addManager(manager network.BaseManager, id string) {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	node.managers[id] = manager
}

// connect to specific remote peers
func (node *Node) connectsToTargets() {
	fmt.Println("connecting to peers ...")
	for _, addr := range node.targets {
		c, err := network.NewOutgoingConnection(addr, node.OnReceive, node.OnFinish)
		if err != nil {
			continue
		}
		node.newConn <- c
	}
}

// listen from remote peers
func (node *Node) listen() error {
	listener, err := net.Listen(network.TCP, node.p2pAddress)
	if err != nil {
		return err
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}
			c := network.NewIncomingConnection(conn, node.OnReceive, node.OnFinish)
			node.newConn <- c
		}
	}()
	for {
		select {
		case connection := <-node.newConn:
			fmt.Println("accepted new client from address ", connection.RemoteAddress())
			node.addConnection(connection)
			connection.Start()
		case doneConnection := <-node.doneConn:
			fmt.Println("disconnected client from address ", doneConnection.RemoteAddress())
			node.removeConnection(doneConnection)
		}
	}
	return nil
}

func (node *Node) addConnection(c *network.Connection) {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	node.conns[c.RemoteAddress()] = c
}

func (node *Node) sendHandshake(c *network.Connection) {
	if c.IsOutgoing() {
		handshake := node.newHandshakePacket()
		c.Send(handshake)
	}
}

func (node *Node) removeConnection(c *network.Connection) {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	c.Close()
	delete(node.conns, c.RemoteAddress())
}

func (node *Node) OnReceive(receiveMessage network.ReceiveMessage) {
	messageType := receiveMessage.Message.Header.Type
	c := receiveMessage.Conn
	switch messageType {
	case network.Handshake:
		fmt.Println("handshake")
		handshake := network.HandshakePacket{}
		network.UnmarshalBinary(receiveMessage.Message.Payload, &handshake)
		node.handleHandshake(c, handshake)
	case network.Notice:
	case network.Request:
	case network.Block:
	case network.RequestNewTerm, network.RequestVote, network.RequestVoteResponse:
		node.managers[network.ElectionManager].Receive(c, receiveMessage.Message)
	}
}

func (node *Node) OnFinish(c *network.Connection) {
	node.doneConn <- c
}

func (node *Node) Signer(hash blockchain.SHA256Type) crypto.Signature {
	privateKey := node.keyPair.privateKey
	sign, _ := privateKey.Sign(hash[:])
	return sign
}

func (node *Node) handleHandshake(c *network.Connection, handshake network.HandshakePacket) {
}

func (node *Node) newHandshakePacket() network.HandshakePacket {
	publicKey := node.keyPair.publicKey
	privateKey := node.keyPair.privateKey
	info := network.HandshakeInfo{
		Network:				network.TestNet,
		Version:				1,
		ChainId: 				node.chainId,
		NodeId: 				node.id,
		Key: 					*publicKey,
		OriginAddress: 			node.p2pAddress,
		LastCommitBlockHeight: 	0,
		LastCommitBlockId: 		blockchain.SHA256Type{},
		TopBlockHeight: 		0,
		TopBlockId:				blockchain.SHA256Type{},
		Timestamp:				time.Now(),
	}
	buf, _ := network.MarshalBinary(info)
	hash := sha256.Sum256(buf)
	sign, _ := privateKey.Sign(hash[:])
	return network.HandshakePacket{
		Info: info,
		Sign: sign,
	}
}