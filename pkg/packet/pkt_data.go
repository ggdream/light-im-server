package packet

type PingPktData struct{}

type PongPktData struct{}

type MessagePktData struct {
	SenderID       string `json:"sender_id"`
	ReceiverID     string `json:"receiver_id"`
	ConversationID string `json:"conversation_id"`
	Type           uint8  `json:"type"`
	Text           string `json:"text"`
	Image          string `json:"image"`
	Audio          string `json:"audio"`
	Video          string `json:"video"`
	Custom         string `json:"custom"`
	IsRead         uint8  `json:"is_read"`
	Timestamp      int64  `json:"timestamp"`
	Sequence       int64  `json:"sequence"`
	CreateAt       int64  `json:"create_at"`
}
