package packet

import "encoding/json"

type Packet struct {
	Type uint16           `json:"type"`
	Data *json.RawMessage `json:"data"`
}

func New() *Packet {
	return &Packet{}
}

func (p *Packet) Set(packetType PacketType, data any) {
	value, _ := json.Marshal(data)
	tmp := json.RawMessage(value)

	p.Type = packetType
	p.Data = &tmp
}

func (p *Packet) Encode() []byte {
	data, _ := json.Marshal(p)

	return data
}

func (p *Packet) Decode(data []byte) (any, error) {
	err := json.Unmarshal(data, p)
	if err != nil {
		return nil, err
	}

	var body any
	switch p.Type {
	case PingPacketType:
		body = &PingPktData{}
	case AuthPacketType:
		body = &AuthPktData{}
	default:
		return nil, ErrInvalidPacketType
	}

	err = json.Unmarshal(*p.Data, &body)
	return body, err
}
