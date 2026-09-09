package network

import (
	"fmt"
	"net"
	"time"
)

type Node struct {
	ID          string
	Address     string
	PeerManager *PeerManager
}

func (n *Node) Start() error {
	n.PeerManager = NewPeerManager()

	listener, err := startServer(n.Address)
	if err != nil {
		return err
	}

	fmt.Println("Node", n.ID, "started")

	go func() {
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Accept error:", err)
				continue
			}

			go n.handleConnection(conn)
		}
	}()

	n.StartHeartbeat()

	return nil
}

func (n *Node) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	peer := &Peer{
		Address:   address,
		Conn:      conn,
		Connected: true,
		LastSeen:  time.Now(),
	}

	n.PeerManager.AddPeer(peer)

	fmt.Println("Node", n.ID, "connected to", address)

	go ReceiveMessages(
		conn,
		func(message Message) {
			n.handleMessage(peer, message)
		},
		func() {
			n.PeerManager.RemovePeer(address)
			fmt.Println("Removed disconnected peer:", address)
		},
	)

	if err := SendMessage(peer, Message{Type: "HELLO", SenderID: n.ID}); err != nil {
		_ = conn.Close()
		return err
	}

	return nil
}

func (n *Node) SendMessage(conn net.Conn, message Message) error {
	for _, peer := range n.PeerManager.GetPeers() {
		if peer.Conn == conn {
			return SendMessage(peer, message)
		}
	}

	return fmt.Errorf("peer connection not found")
}

func (n *Node) handleConnection(conn net.Conn) {
	address := conn.RemoteAddr().String()

	peer := &Peer{
		Address:   address,
		Conn:      conn,
		Connected: true,
		LastSeen:  time.Now(),
	}

	n.PeerManager.AddPeer(peer)

	fmt.Println("Incoming connection from:", address)

	if err := SendMessage(peer, Message{Type: "HELLO", SenderID: n.ID}); err != nil {
		_ = conn.Close()
		return
	}

	go ReceiveMessages(
		conn,
		func(message Message) {
			n.handleMessage(peer, message)
		},
		func() {
			n.PeerManager.RemovePeer(address)
			fmt.Println("Removed disconnected peer:", address)
		},
	)
}

func (n *Node) handleMessage(peer *Peer, message Message) {
	peer.LastSeen = time.Now()

	switch message.Type {
	case "HELLO":
		if message.SenderID == "" {
			fmt.Println("Invalid HELLO from:", peer.Address)
			return
		}

		peer.ID = message.SenderID
		fmt.Println("HELLO received from:", peer.ID)
	case "PING":
		fmt.Println("PING received from:", message.SenderID)
		if err := SendMessage(peer, Message{Type: "PONG", SenderID: n.ID}); err != nil {
			fmt.Println("Failed to send PONG:", err)
		}
	case "PONG":
		fmt.Println("PONG received from:", message.SenderID)
	case "BROADCAST":
		fmt.Println("Message received:")
		fmt.Println("  Type:", message.Type)
		fmt.Println("  Sender:", message.SenderID)
		fmt.Println("  Payload:", message.Payload)
	default:
		fmt.Println("Unknown message type:", message.Type)
	}
}

func (n *Node) Broadcast(message Message) {
	peers := n.PeerManager.GetPeers()
	if len(peers) == 0 {
		fmt.Println("No peers connected")
		return
	}

	fmt.Println("Broadcasting message to", len(peers), "peers")
	for _, peer := range peers {
		if !peer.Connected {
			continue
		}
		if err := SendMessage(peer, message); err != nil {
			fmt.Println("Broadcast failed to", peer.Address, ":", err)
			n.PeerManager.RemovePeer(peer.Address)
			continue
		}
		fmt.Println("Broadcast sent to:", peer.Address)
	}
}
