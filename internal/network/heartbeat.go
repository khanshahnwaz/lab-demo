package network

import (
	"fmt"
	"time"
)

const (
	HeartbeatInterval = 5 * time.Second
	HeartbeatTimeout  = 15 * time.Second
)

func (n *Node) StartHeartbeat() {
	ticker := time.NewTicker(HeartbeatInterval)

	go func() {
		for range ticker.C {
			n.sendHeartbeats()
			n.checkPeerTimeouts()
		}
	}()
}

func (n *Node) sendHeartbeats() {
	for _, peer := range n.PeerManager.GetPeers() {
		if !peer.Connected {
			continue
		}

		err := SendMessage(peer, Message{Type: "PING", SenderID: n.ID})
		if err != nil {
			fmt.Println("Heartbeat failed:", peer.Address, err)
			n.PeerManager.RemovePeer(peer.Address)
			continue
		}

		fmt.Println("PING ->", peer.Address)
	}
}

func (n *Node) checkPeerTimeouts() {
	now := time.Now()
	for _, peer := range n.PeerManager.GetPeers() {
		if now.Sub(peer.LastSeen) > HeartbeatTimeout {
			fmt.Println("Heartbeat timeout:", peer.Address)
			n.PeerManager.RemovePeer(peer.Address)
		}
	}
}
