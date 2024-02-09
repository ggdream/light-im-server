package hub

import (
	"context"
	"sync"
	"unsafe"

	"github.com/gorilla/websocket"

	"lim/config"
	"lim/pkg/cache"
	"lim/pkg/packet"
	"lim/tools/token"
)

var (
	h = &Hub{
		pool:    make(map[string]*Conn),
		locker:  &sync.Mutex{},
		rootCtx: context.TODO(),
	}

	authTokenCa cache.AuthToken
)

type Hub struct {
	pool    map[string]*Conn
	locker  *sync.Mutex
	rootCtx context.Context
}

func (h *Hub) SetConn(conn *websocket.Conn) {
	c := NewConn(conn)
	_, data, err := c.conn.ReadMessage() // Fist Pkt: auth packet
	if err != nil {
		return
	}

	pkt := packet.New()
	body, err := pkt.Decode(data)
	if err != nil || pkt.Type != packet.AuthPacketType {
		c.Close()
		return
	}
	bodyData := body.(*packet.AuthPktData)
	tokenKey := unsafe.Slice(unsafe.StringData(config.GetApp().TokenKey), len(config.GetApp().TokenKey))
	userId, err := token.Parse(tokenKey, bodyData.Token)
	if err != nil {
		retPkt := packet.New()
		retPkt.Set(packet.PassPacketType, &packet.PassPktData{
			Code: 1,
			Desc: "凭证认证失败",
		})
		c.conn.WriteMessage(websocket.TextMessage, retPkt.Encode())
		c.Close()
		return
	}

	isPass, err := authTokenCa.Verify(bodyData.UserID, bodyData.Token)

	if err != nil {
		// TODO: log redis error
		retPkt := packet.New()
		retPkt.Set(packet.PassPacketType, &packet.PassPktData{
			Code: 2,
			Desc: "缓存读取失败",
		})
		c.conn.WriteMessage(websocket.TextMessage, retPkt.Encode())
		c.Close()
		return
	}
	if !isPass {
		retPkt := packet.New()
		retPkt.Set(packet.PassPacketType, &packet.PassPktData{
			Code: 1,
			Desc: "凭证认证失败",
		})
		c.conn.WriteMessage(websocket.TextMessage, retPkt.Encode())
		c.Close()
		return
	}

	retPkt := packet.New()
	retPkt.Set(packet.PassPacketType, &packet.PassPktData{
		Code: 0,
		Desc: "登录成功",
	})
	c.conn.WriteMessage(websocket.TextMessage, retPkt.Encode())

	c.userId = userId
	ctx, cancel := context.WithCancel(h.rootCtx)
	c.cancelFunc = cancel

	go c.Dispatch(ctx)

	h.locker.Lock()
	defer h.locker.Unlock()
	h.pool[userId] = c
}

func (h *Hub) GetConn(userId string) *Conn {
	h.locker.Lock()
	defer h.locker.Unlock()

	return h.pool[userId]
}

func (h *Hub) DelConn(userId string) {
	h.locker.Lock()
	defer h.locker.Unlock()

	c, ok := h.pool[userId]
	if !ok {
		return
	}

	c.cancelFunc()
	delete(h.pool, userId)
}

func SetConn2Hub(conn *websocket.Conn) {
	h.SetConn(conn)
}

func DelConn4Hub(userId string) {
	h.DelConn(userId)
}

func Write2Conn(userId string, pkt *packet.Packet) error {
	conn, ok := h.pool[userId]
	if !ok {
		return nil
	}

	return conn.Write(pkt.Encode())
}
