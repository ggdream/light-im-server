package hub

import (
	"sync/atomic"
	"time"
	
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"lim/pkg/packet"
)

var (
	ErrChannelClosed = errors.New("channel closed")
)

type Conn struct {
	userId             string
	conn               *websocket.Conn
	timer              *time.Timer
	wChannel, rChannel chan []byte
	signal             int32
}

func NewConn(userId string, conn *websocket.Conn) *Conn {
	return &Conn{
		userId: userId,
		conn:   conn,
	}
}

func (c *Conn) Dispatch() {
	c.timer = time.NewTimer(time.Second * 15)

	go func() {
		for {
			_, data, err := c.conn.ReadMessage()
			if err != nil {
				return
			}

			if atomic.LoadInt32(&c.signal) == 1 {
				return
			}

			c.rChannel <- data
		}
	}()

	defer c.Close()

	for {
		select {
		case <-c.timer.C:
			return
		case data, ok := <-c.wChannel:
			if !ok {
				return
			}

			err := c.conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				return
			}
		case data, ok := <-c.rChannel:
			if !ok {
				return
			}

			_ = c.timer.Reset(time.Second * 15)
			go c.Read(data)
		}
	}
}

func (c *Conn) Write(data []byte) error {
	if atomic.LoadInt32(&c.signal) == 1 {
		return ErrChannelClosed
	}

	c.wChannel <- data
	return nil
}

func (c *Conn) Read(data []byte) {
	pkt := packet.New()
	if err := pkt.Decode(data); err != nil {
		c.Close()
		return
	}

	var (
		retPkt = &packet.Packet{}
	)

	switch pkt.Type {
	case packet.PingPacketType:
		retPkt.Set(packet.PongPacketType, &packet.PongPktData{})
	default:
		return
	}

	c.write(retPkt)
}

func (c *Conn) write(pkt *packet.Packet) error {
	return c.Write(pkt.Encode())
}

func (c *Conn) Close() {
	if atomic.LoadInt32(&c.signal) == 1 {
		return
	}

	atomic.StoreInt32(&c.signal, 1)
	c.timer.Stop()
	// close(c.wChannel)
	// close(c.rChannel)
	_ = c.conn.Close()
	c.conn = nil
}
