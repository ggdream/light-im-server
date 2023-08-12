package packet

import "github.com/pkg/errors"

type (
	PacketType  = uint16
	MessageType = uint8
)

const (
	PingPacketType    PacketType = 1
	PongPacketType    PacketType = 2
	AuthPacketType    PacketType = 3
	PassPacketType    PacketType = 4
	MessagePacketType PacketType = 5
)

const (
	TextMessageType   MessageType = 1
	ImageMessageType  MessageType = 2
	AudioMessageType  MessageType = 3
	VideoMessageType  MessageType = 4
	FileMessageType   MessageType = 5
	CustomMessageType MessageType = 6
	RecordMessageType MessageType = 7
)

var (
	ErrInvalidPacketType = errors.New("invalid: unknown packet type")
)
