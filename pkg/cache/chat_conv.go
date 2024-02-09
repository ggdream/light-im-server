package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type ChatConv struct {
	SenderID       string `json:"sender_id" redis:"sender_id"`
	ReceiverID     string `json:"receiver_id" redis:"receiver_id"`
	GroupID        string `json:"group_id" redis:"group_id"`
	ConversationID string `json:"conversation_id" redis:"conversation_id"`
	Type           uint8  `json:"type" redis:"type"`
	Text           string `json:"text" redis:"text"`
	Image          string `json:"image" redis:"image"`
	Audio          string `json:"audio" redis:"audio"`
	Video          string `json:"video" redis:"video"`
	File           string `json:"file" redis:"file"`
	Custom         string `json:"custom" redis:"custom"`
	Record         string `json:"record" redis:"record"`
	Timestamp      int64  `json:"timestamp" redis:"timestamp"`
	Sequence       int64  `json:"sequence" redis:"sequence"`
	CreateAt       int64  `json:"create_at" redis:"create_at"`
	Unread         int64  `json:"unread" redis:"unread,omitempty"`
}

func (c *ChatConv) Add(ctx context.Context, senderId string, receiverIds []string, conversationId string) error {
	_, err := client.Pipelined(ctx, func(p redis.Pipeliner) error {
		z := redis.Z{
			Score:  float64(c.CreateAt),
			Member: conversationId,
		}
		err := p.ZAdd(ctx, c.joinName(senderId), z).Err()
		if err != nil {
			return err
		}
		if receiverIds != nil {
			for _, receiverId := range receiverIds {
				err = p.ZAdd(ctx, c.joinName(receiverId), z).Err()
				if err != nil {
					return err
				}
			}
		}

		err = p.HSet(ctx, c.joinName1(senderId, conversationId), c).Err()
		if err != nil {
			return err
		}
		if receiverIds != nil {
			for _, receiverId := range receiverIds {
				err = p.HSet(ctx, c.joinName1(receiverId, conversationId), c).Err()
				if err != nil {
					return err
				}
				err = p.HIncrBy(ctx, c.joinName1(receiverId, conversationId), "unread", 1).Err()
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return err
}

func (c *ChatConv) MarkAsRead(ctx context.Context, userId, conversationId string) error {
	return client.HSet(ctx, c.joinName1(userId, conversationId), "unread", 0).Err()
}

func (c *ChatConv) Del(ctx context.Context, userId, conversationId string) error {
	_, err := client.Pipelined(ctx, func(p redis.Pipeliner) error {
		err := p.Unlink(ctx, c.joinName1(userId, conversationId)).Err()
		if err != nil {
			return err
		}

		return p.ZRem(ctx, c.joinName(userId), conversationId).Err()
	})

	return err
}

func (c *ChatConv) List(ctx context.Context, userId string) ([]ChatConv, error) {
	res, err := client.ZRevRange(ctx, c.joinName(userId), 0, -1).Result()
	if err != nil {
		return nil, err
	}

	ret := make([]ChatConv, 0, len(res))
	for _, v := range res {
		var r ChatConv
		err = client.HGetAll(ctx, c.joinName1(userId, v)).Scan(&r)
		if err != nil {
			return nil, err
		}

		ret = append(ret, r)
	}

	return ret, nil
}

func (c *ChatConv) joinName(userId string) string {
	return fmt.Sprintf("lim:chat:conv:list:%s", userId)
}

func (c *ChatConv) joinName1(userId, conversationId string) string {
	return fmt.Sprintf("lim:chat:conv:info:%s:%s", userId, conversationId)
}
