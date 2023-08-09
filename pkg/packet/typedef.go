package packet

import "github.com/pkg/errors"

type (
	PacketType = uint16
)

const (
	PingPacketType    PacketType = 1
	PongPacketType    PacketType = 2
	AuthPacketType    PacketType = 3
	PassPacketType    PacketType = 4
	MessagePacketType PacketType = 5
)

var (
	ErrInvalidPacketType = errors.New("invalid: unknown packet type")
)
