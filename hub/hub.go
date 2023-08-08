package hub

import (
	"sync"

	"github.com/gorilla/websocket"

	"lim/pkg/packet"
)

var (
	h = &Hub{
		pool:   make(map[string]*Conn),
		locker: &sync.Mutex{},
	}
)

type Hub struct {
	pool   map[string]*Conn
	locker *sync.Mutex
}

func (h *Hub) SetConn(userId string, conn *websocket.Conn) {
	h.locker.Lock()
	defer h.locker.Unlock()

	c := NewConn(userId, conn)
	go c.Dispatch()

	h.pool[userId] = c
}

func (h *Hub) GetConn(userId string) *Conn {
	h.locker.Lock()
	defer h.locker.Unlock()

	return h.pool[userId]
}

func SetConn2Hub(userId string, conn *websocket.Conn) {
	h.SetConn(userId, conn)
}

func Write2Conn(userId string, pkt *packet.Packet) error {
	return h.pool[userId].Write(pkt.Encode())
}
