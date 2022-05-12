package broker

import (
	"net"
	"sync"
	"time"
)

type Peer struct {
	Key       string    `json:"key"`
	UpdatedAt time.Time `json:"updated_at"`
	Expire    time.Time `json:"expire"`
	PeerIP    string    `json:"peer_ip"`
	PeerPort  int       `json:"peer_port"`
}

func NewPeerManager() *PeerManager {
	return &PeerManager{
		peers: map[string]peer{},
		mtx:   sync.RWMutex{},
	}
}

type PeerManager struct {
	peers map[string]peer
	mtx   sync.RWMutex
}

type peer struct {
	UpdatedAt time.Time
	Expire    time.Time
	Addr      *net.UDPAddr
}

func (u *PeerManager) Get(key string) Peer {
	u.mtx.RLock()
	defer u.mtx.RUnlock()

	d := u.peers[key]

	return Peer{
		Key:       key,
		UpdatedAt: d.UpdatedAt,
		Expire:    d.Expire,
		PeerIP:    d.Addr.IP.String(),
		PeerPort:  d.Addr.Port,
	}
}

func (u *PeerManager) Set(key string, addr *net.UDPAddr) {
	u.mtx.Lock()
	defer u.mtx.Unlock()

	now := time.Now()

	u.peers[key] = peer{
		UpdatedAt: now,
		Expire:    now.Add(60 * time.Second),
		Addr:      addr,
	}
}

func (u *PeerManager) List() []Peer {
	u.mtx.RLock()
	defer u.mtx.RUnlock()

	peers := make([]Peer, 0, len(u.peers))
	for k, v := range u.peers {
		peers = append(peers, Peer{
			Key:       k,
			UpdatedAt: v.UpdatedAt,
			Expire:    v.Expire,
			PeerIP:    v.Addr.IP.String(),
			PeerPort:  v.Addr.Port,
		})
	}

	return peers
}

func (u *PeerManager) Clean() {
	u.mtx.Lock()
	defer u.mtx.Unlock()

	now := time.Now()

	for k, v := range u.peers {
		if now.After(v.Expire) {
			delete(u.peers, k)
		}
	}
}
