package packet

type PingPktData struct{}

type PongPktData struct{}

type AuthPktData struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

type PassPktData struct {
	Code uint16 `json:"code"`
	Desc string `json:"desc"`
}

type MessagePktData struct {
	SenderID       string             `json:"sender_id"`
	ReceiverID     string             `json:"receiver_id"`
	UserID         string             `json:"user_id"`
	ConversationID string             `json:"conversation_id"`
	Type           uint8              `json:"type"`
	Text           *MessageTextElem   `json:"text"`
	Image          *MessageImageElem  `json:"image"`
	Audio          *MessageAudioElem  `json:"audio"`
	Video          *MessageVideoElem  `json:"video"`
	File           *MessageFileElem   `json:"file"`
	Custom         *MessageCustomElem `json:"custom"`
	Record         *MessageRecordElem `json:"record"`
	IsSelf         uint8              `json:"is_self"`
	IsRead         uint8              `json:"is_read"`
	IsPeerRead     uint8              `json:"is_peer_read"`
	Timestamp      int64              `json:"timestamp"`
	Sequence       int64              `json:"sequence"`
	CreateAt       int64              `json:"create_at"`
}

type MessageTextElem struct {
	Text string `json:"text"`
}

type MessageImageElem struct {
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	ContentType  string `json:"content_type"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type MessageAudioElem struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	Duration    int64  `json:"duration"`
	URL         string `json:"url"`
}

type MessageVideoElem struct {
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	ContentType  string `json:"content_type"`
	Duration     int64  `json:"duration"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type MessageFileElem struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

type MessageCustomElem struct {
	Content string `json:"content"`
}

type MessageRecordElem struct {
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	Duration    int64  `json:"duration"`
	URL         string `json:"url"`
}
