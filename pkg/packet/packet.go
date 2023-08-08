package packet

import "encoding/json"

type Packet struct {
	Type uint16           `json:"type"`
	Data *json.RawMessage `json:"data"`
}

func New() *Packet {
	return &Packet{}
}

func (p *Packet) Set(packetType PacketType, data any) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	tmp := json.RawMessage(value)

	p.Type = packetType
	p.Data = &tmp

	return nil
}

func (p *Packet) Encode() []byte {
	data, _ := json.Marshal(p)

	return data
}

func (p *Packet) Decode(data []byte) error {
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	switch p.Type {
	case PingPacketType:
		return nil
	default:
		return ErrInvalidPacketType
	}
}
