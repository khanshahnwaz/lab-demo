package network

import (
	"net"
	"sync"
	"time"
)

type Peer struct {
	ID        string
	Address   string
	Conn      net.Conn
	Connected bool
	LastSeen  time.Time

	WriteMux sync.Mutex
}
