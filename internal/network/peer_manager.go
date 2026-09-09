package network

import "sync"

type PeerManager struct {
	Peers map[string]*Peer
	Mux   sync.RWMutex
}

func NewPeerManager() *PeerManager {
	return &PeerManager{
		Peers: make(map[string]*Peer),
	}
}

func (pm *PeerManager) AddPeer(peer *Peer) {
	pm.Mux.Lock()
	defer pm.Mux.Unlock()

	if existing, exists := pm.Peers[peer.Address]; exists && existing.Conn != nil {
		_ = existing.Conn.Close()
	}

	pm.Peers[peer.Address] = peer
}

func (pm *PeerManager) RemovePeer(address string) {
	pm.Mux.Lock()
	defer pm.Mux.Unlock()

	peer, exists := pm.Peers[address]
	if !exists {
		return
	}

	peer.Connected = false
	if peer.Conn != nil {
		_ = peer.Conn.Close()
	}

	delete(pm.Peers, address)
}

func (pm *PeerManager) GetPeers() []*Peer {
	pm.Mux.RLock()
	defer pm.Mux.RUnlock()

	peers := make([]*Peer, 0, len(pm.Peers))

	for _, peer := range pm.Peers {
		peers = append(peers, peer)
	}

	return peers
}
