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
	"lim/tools/typec"
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
	SenderID   string         `json:"sender_id"`
	ReceiverID string         `json:"receiver_id"`
	GroupID    string         `json:"group_id"`
	Type       *uint8         `json:"type" binding:"required"`
	Text       *db.TextElem   `json:"text"`
	Image      *db.ImageElem  `json:"image"`
	Audio      *db.AudioElem  `json:"audio"`
	Video      *db.VideoElem  `json:"video"`
	File       *db.FileElem   `json:"file"`
	Custom     *db.CustomElem `json:"custom"`
	Record     *db.RecordElem `json:"record"`
	Timestamp  *int64         `json:"timestamp" binding:"required"`
}

type sendResModel struct {
	UserID         string         `json:"user_id"`
	SenderID       string         `json:"sender_id"`
	ReceiverID     string         `json:"receiver_id"`
	GroupID        string         `json:"group_id"`
	ConversationID string         `json:"conversation_id"`
	Type           uint8          `json:"type"`
	Text           *db.TextElem   `json:"text"`
	Image          *db.ImageElem  `json:"image"`
	Audio          *db.AudioElem  `json:"audio"`
	Video          *db.VideoElem  `json:"video"`
	File           *db.FileElem   `json:"file"`
	Custom         *db.CustomElem `json:"custom"`
	Record         *db.RecordElem `json:"record"`
	IsSelf         uint8          `json:"is_self"`
	IsRead         uint8          `json:"is_read"`
	IsPeerRead     uint8          `json:"is_peer_read"`
	Timestamp      int64          `json:"timestamp"`
	Sequence       int64          `json:"sequence"`
	CreateAt       int64          `json:"create_at"`
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

		seq := node.Generate()
		doc := db.ChatMsgDoc{
			UserID:     senderId,
			SenderID:   senderId,
			ReceiverID: form.ReceiverID,
			GroupID:    form.GroupID,
			Type:       *form.Type,
			Text:       form.Text,
			Image:      form.Image,
			Audio:      form.Audio,
			Video:      form.Video,
			File:       form.File,
			Custom:     form.Custom,
			Record:     form.Record,
			Timestamp:  *form.Timestamp,
			Sequence:   seq.Int64(),
		}
		receiverIds, err := doc.Create(c)
		if err != nil {
			errno.NewF(errno.BaseErrMongo, err.Error(), errno.ErrChatMsgSaveFailed).Reply(c)
			return
		}

		ca := cache.ChatConv{
			SenderID:       senderId,
			ReceiverID:     doc.ReceiverID,
			GroupID:        doc.GroupID,
			ConversationID: doc.ConversationID,
			Type:           *form.Type,
			Text:           typec.MapToJson(form.Text),
			Image:          typec.MapToJson(form.Image),
			Audio:          typec.MapToJson(form.Audio),
			Video:          typec.MapToJson(form.Video),
			File:           typec.MapToJson(form.File),
			Custom:         typec.MapToJson(form.Custom),
			Record:         typec.MapToJson(form.Record),
			Timestamp:      *form.Timestamp,
			Sequence:       seq.Int64(),
			CreateAt:       doc.CreateAt,
		}
		_ = ca.Add(c, senderId, receiverIds, doc.ConversationID)

		text, image, audio, video, file, custom, record := (*packet.MessageTextElem)(nil), (*packet.MessageImageElem)(nil), (*packet.MessageAudioElem)(nil), (*packet.MessageVideoElem)(nil), (*packet.MessageFileElem)(nil), (*packet.MessageCustomElem)(nil), (*packet.MessageRecordElem)(nil)
		switch *form.Type {
		case packet.TextMessageType:
			text = &packet.MessageTextElem{
				Text: form.Text.Text,
			}
		case packet.ImageMessageType:
			image = &packet.MessageImageElem{
				Name:         form.Image.Name,
				Size:         form.Image.Size,
				ContentType:  form.Image.ContentType,
				URL:          form.Image.URL,
				ThumbnailURL: form.Image.ThumbnailURL,
			}
		case packet.AudioMessageType:
			audio = &packet.MessageAudioElem{
				Name:        form.Audio.Name,
				Size:        form.Audio.Size,
				ContentType: form.Audio.ContentType,
				Duration:    form.Audio.Duration,
				URL:         form.Audio.URL,
			}
		case packet.VideoMessageType:
			video = &packet.MessageVideoElem{
				Name:         form.Video.Name,
				Size:         form.Video.Size,
				ContentType:  form.Video.ContentType,
				Duration:     form.Video.Duration,
				URL:          form.Video.URL,
				ThumbnailURL: form.Video.ThumbnailURL,
			}
		case packet.FileMessageType:
			file = &packet.MessageFileElem{
				Name:        form.File.Name,
				Size:        form.File.Size,
				ContentType: form.File.ContentType,
				URL:         form.File.URL,
			}
		case packet.CustomMessageType:
			custom = &packet.MessageCustomElem{
				Content: form.Custom.Content,
			}
		case packet.RecordMessageType:
			record = &packet.MessageRecordElem{
				Size:        form.Record.Size,
				ContentType: form.Record.ContentType,
				Duration:    form.Record.Duration,
				URL:         form.Record.URL,
			}
		default:
			return
		}

		pkt1 := packet.New()
		data1 := &packet.MessagePktData{
			UserID:         doc.UserID,
			SenderID:       doc.SenderID,
			ReceiverID:     doc.ReceiverID,
			GroupID:        doc.GroupID,
			ConversationID: doc.ConversationID,
			Type:           *form.Type,
			Text:           text,
			Image:          image,
			Audio:          audio,
			Video:          video,
			File:           file,
			Custom:         custom,
			Record:         record,
			Timestamp:      doc.Timestamp,
			IsSelf:         1,
			IsRead:         1,
			IsPeerRead:     0,
			Sequence:       doc.Sequence,
			CreateAt:       doc.CreateAt,
		}
		pkt1.Set(packet.MessagePacketType, data1)
		err = hub.Write2Conn(senderId, pkt1)
		pkt2 := packet.New()
		data2 := &packet.MessagePktData{
			UserID:         doc.UserID,
			SenderID:       doc.SenderID,
			ReceiverID:     doc.ReceiverID,
			GroupID:        doc.GroupID,
			ConversationID: doc.ConversationID,
			Type:           *form.Type,
			Text:           text,
			Image:          image,
			Audio:          audio,
			Video:          video,
			File:           file,
			Custom:         custom,
			Record:         record,
			Timestamp:      doc.Timestamp,
			IsSelf:         0,
			IsRead:         0,
			IsPeerRead:     1,
			Sequence:       doc.Sequence,
			CreateAt:       doc.CreateAt,
		}
		pkt2.Set(packet.MessagePacketType, data2)
		err = hub.Write2Conn(form.ReceiverID, pkt2)

		ret := &sendResModel{
			UserID:         doc.UserID,
			SenderID:       doc.SenderID,
			ReceiverID:     doc.ReceiverID,
			GroupID:        doc.GroupID,
			ConversationID: doc.ConversationID,
			Type:           *form.Type,
			Text:           form.Text,
			Image:          form.Image,
			Audio:          form.Audio,
			Video:          form.Video,
			File:           form.File,
			Custom:         form.Custom,
			Record:         form.Record,
			Timestamp:      doc.Timestamp,
			IsSelf:         1,
			IsRead:         1,
			IsPeerRead:     0,
			Sequence:       doc.Sequence,
			CreateAt:       doc.CreateAt,
		}
		errno.NewS(ret).Reply(c)
	}
}
