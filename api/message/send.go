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

type sendReqModel struct {
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

func SendController() gin.HandlerFunc {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		var (
			form sendReqModel
			err  error
		)
		if err = c.ShouldBindJSON(&form); err != nil {
			errno.NewFParamInvalid(err.Error()).Reply(c)
			return
		}

		userId := config.CtxKeyManager.GetUserID(c)
		seq := node.Generate().Int64()
		doc := db.ChatMsgDoc{
			SenderID:   userId,
			ReceiverID: *form.ReceiverID,
			Type:       *form.Type,
			Text:       form.Text,
			Image:      form.Image,
			Audio:      form.Audio,
			Video:      form.Video,
			Custom:         form.Custom,
			Timestamp:  *form.Timestamp,
			Sequence:   seq,
		}
		err = doc.Create()
		if err != nil {
			errno.NewF(errno.BaseErrMongo, err.Error(), errno.ErrChatMsgSaveFailed).Reply(c)
			return
		}

		ca := cache.ChatConv{
			SenderID:       userId,
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
		_ = ca.Add(userId, *form.ReceiverID)

		pkt := packet.New()
		data := &packet.MessagePktData{
			SenderID:       userId,
			ReceiverID:     *form.ReceiverID,
			ConversationID: doc.ConversationID,
			Type:           *form.Type,
			Text:           form.Text,
			Image:          form.Image,
			Audio:          form.Audio,
			Video:          form.Video,
			Custom:         form.Custom,
			Timestamp:      *form.Timestamp,
			IsRead:         doc.IsRead,
			Sequence:       seq,
			CreateAt:       doc.CreateTs,
		}
		pkt.Set(packet.MessagePacketType, data)
		_ = hub.Write2Conn(*form.ReceiverID, pkt)

		ret := &sendResModel{
			SenderID:       userId,
			ReceiverID:     *form.ReceiverID,
			ConversationID: doc.ConversationID,
			Type:           *form.Type,
			Text:           form.Text,
			Image:          form.Image,
			Audio:          form.Audio,
			Video:          form.Video,
			Custom:         form.Custom,
			Timestamp:      *form.Timestamp,
			IsRead:         doc.IsRead,
			Sequence:       seq,
			CreateAt:       doc.CreateTs,
		}
		errno.NewS(ret).Reply(c)
	}
}
