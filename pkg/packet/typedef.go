package packet

import "github.com/pkg/errors"

type (
	PacketType = uint16
)

const (
	PingPacketType    PacketType = 1
	PongPacketType    PacketType = 2
	MessagePacketType PacketType = 3
)

var (
	ErrInvalidPacketType = errors.New("invalid: unknown packet type")
)
