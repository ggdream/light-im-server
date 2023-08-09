package message

import (
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"

	"lim/config"
	"lim/hub"
	"lim/pkg/cache"
	"lim/pkg/db"
	"lim/pkg/errno"
	"lim/pkg/packet"
)

var (
	node *snowflake.Node
)

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

type sendReqModel struct {
	SenderID   string  `json:"sender_id"`
	ReceiverID *string `json:"receiver_id" binding:"required"`
	Type       *uint8  `json:"type" binding:"required"`
	Text       string  `json:"text"`
	Image      string  `json:"image"`
	Audio      string  `json:"audio"`
	Video      string  `json:"video"`
	Custom     string  `json:"custom"`
	Timestamp  *int64  `json:"timestamp" binding:"required"`
}

type sendResModel struct {
	SenderID       string `json:"sender_id"`
	ReceiverID     string `json:"receiver_id"`
	UserID         string `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	Type           uint8  `json:"type"`
	Text           string `json:"text"`
	Image          string `json:"image"`
	Audio          string `json:"audio"`
	Video          string `json:"video"`
	Custom         string `json:"custom"`
	IsSelf         uint8  `json:"is_self"`
	IsRead         uint8  `json:"is_read"`
	IsPeerRead     uint8  `json:"is_peer_read"`
	Timestamp      int64  `json:"timestamp"`
	Sequence       int64  `json:"sequence"`
	CreateAt       int64  `json:"create_at"`
}

func SendController(isAdmin bool) gin.HandlerFunc {

	return func(c *gin.Context) {
		var (
			form sendReqModel
			err  error
		)
		if err = c.ShouldBindJSON(&form); err != nil {
			errno.NewFParamInvalid(err.Error()).Reply(c)
			return
		}

		var senderId string
		if isAdmin {
			senderId = form.SenderID
		} else {
			senderId = config.CtxKeyManager.GetUserID(c)
		}

		seq := node.Generate().Int64()
		doc := db.ChatMsgDoc{
			SenderID:   senderId,
			ReceiverID: *form.ReceiverID,
			Type:       *form.Type,
			Text:       form.Text,
			Image:      form.Image,
			Audio:      form.Audio,
			Video:      form.Video,
			Custom:     form.Custom,
			Timestamp:  *form.Timestamp,
			Sequence:   seq,
		}
		err = doc.Create()
		if err != nil {
			errno.NewF(errno.BaseErrMongo, err.Error(), errno.ErrChatMsgSaveFailed).Reply(c)
			return
		}

		ca := cache.ChatConv{
			SenderID:       senderId,
			ReceiverID:     *form.ReceiverID,
			ConversationID: doc.ConversationID,
			Type:           *form.Type,
			Text:           form.Text,
			Image:          form.Image,
			Audio:          form.Audio,
			Video:          form.Video,
			Custom:         form.Custom,
			Timestamp:      *form.Timestamp,
			Sequence:       seq,
			CreateAt:       doc.CreateTs,
		}
		_ = ca.Add(senderId, *form.ReceiverID)

		pkt1 := packet.New()
		data1 := &packet.MessagePktData{
			SenderID:       senderId,
			ReceiverID:     *form.ReceiverID,
			UserID:         *form.ReceiverID,
			ConversationID: doc.ConversationID,
			Type:           *form.Type,
			Text:           form.Text,
			Image:          form.Image,
			Audio:          form.Audio,
			Video:          form.Video,
			Custom:         form.Custom,
			Timestamp:      *form.Timestamp,
			IsSelf:         1,
			IsRead:         1,
			IsPeerRead:     0,
			Sequence:       seq,
			CreateAt:       doc.CreateTs,
		}
		pkt1.Set(packet.MessagePacketType, data1)
		_ = hub.Write2Conn(senderId, pkt1)
		pkt2 := packet.New()
		data2 := &packet.MessagePktData{
			SenderID:       senderId,
			ReceiverID:     *form.ReceiverID,
			UserID:         senderId,
			ConversationID: doc.ConversationID,
			Type:           *form.Type,
			Text:           form.Text,
			Image:          form.Image,
			Audio:          form.Audio,
			Video:          form.Video,
			Custom:         form.Custom,
			Timestamp:      *form.Timestamp,
			IsSelf:         0,
			IsRead:         0,
			IsPeerRead:     1,
			Sequence:       seq,
			CreateAt:       doc.CreateTs,
		}
		pkt2.Set(packet.MessagePacketType, data2)
		_ = hub.Write2Conn(*form.ReceiverID, pkt2)

		ret := &sendResModel{
			SenderID:       senderId,
			ReceiverID:     *form.ReceiverID,
			ConversationID: doc.ConversationID,
			Type:           *form.Type,
			Text:           form.Text,
			Image:          form.Image,
			Audio:          form.Audio,
			Video:          form.Video,
			Custom:         form.Custom,
			Timestamp:      *form.Timestamp,
			IsSelf:         1,
			IsRead:         1,
			IsPeerRead:     0,
			Sequence:       seq,
			CreateAt:       doc.CreateTs,
		}
		errno.NewS(ret).Reply(c)
	}
}
