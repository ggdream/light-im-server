package cache

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type ChatConv struct {
	SenderID       string `json:"sender_id" redis:"sender_id"`
	ReceiverID     string `json:"receiver_id" redis:"receiver_id"`
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
	Unread         int64  `json:"unread" redis:"unread"`
}

func (c *ChatConv) Add(senderId, receiverId string) error {
	ctx, cancel := withTimeout()
	defer cancel()

	_, err := client.Pipelined(ctx, func(p redis.Pipeliner) error {
		z1 := redis.Z{
			Score:  float64(c.CreateAt),
			Member: receiverId,
		}
		err := p.ZAdd(ctx, c.joinName(senderId), z1).Err()
		if err != nil {
			return err
		}
		z2 := redis.Z{
			Score:  float64(c.CreateAt),
			Member: senderId,
		}
		err = p.ZAdd(ctx, c.joinName(receiverId), z2).Err()
		if err != nil {
			return err
		}

		err = p.HSet(ctx, c.joinName1(senderId, receiverId), c).Err()
		if err != nil {
			return err
		}

		err = p.HSet(ctx, c.joinName1(receiverId, senderId), c).Err()
		if err != nil {
			return err
		}

		return p.HIncrBy(ctx, c.joinName1(receiverId, senderId), "unread", 1).Err()
	})
	return err
}

func (c *ChatConv) MarkAsRead(userId, conversationId string) error {
	ctx, cancel := withTimeout()
	defer cancel()

	return client.HDel(ctx, c.joinName1(userId, conversationId), "unread").Err()
}

func (c *ChatConv) Del(userId, conversationId string) error {
	ctx, cancel := withTimeout()
	defer cancel()

	_, err := client.Pipelined(ctx, func(p redis.Pipeliner) error {
		err := p.Unlink(ctx, c.joinName1(userId, conversationId)).Err()
		if err != nil {
			return err
		}

		return p.ZRem(ctx, c.joinName(userId), conversationId).Err()
	})
	return err
}

func (c *ChatConv) List(userId string) ([]ChatConv, error) {
	ctx, cancel := withTimeout()
	defer cancel()

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
