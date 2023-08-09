package db

import (
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserInfoDoc struct {
	UserID   string `json:"user_id" bson:"user_id"`
	Nickname string `json:"nickname" bson:"nickname"`
	Avatar   string `json:"avatar" bson:"avatar"`
	CreateAt string `json:"create_at" bson:"create_at"`
	DeleteAt string `json:"delete_at" bson:"delete_at"`
	// UpdateAt int64  `json:"update_at" bson:"update_at"`
}

func (m *UserInfoDoc) Create(userId, nickname, avatar string) error {
	err := client.SearchOne(m.DocName(), bson.D{{Key: "user_id", Value: userId}}, nil).Err()
	if err == nil {
		return errors.New("用户已存在")
	}
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
	}

	m.UserID = userId
	m.Nickname = nickname
	m.Avatar = avatar
	m.CreateAt = time.Now().Format(time.RFC3339)

	return client.Insert(m.DocName(), m)
}

func (m *UserInfoDoc) Search(userId string) error {
	return client.SearchOne(m.DocName(), bson.D{{Key: "user_id", Value: userId}, {Key: "delete_at", Value: ""}}, nil).Decode(m)
}

func (m *UserInfoDoc) Update(userId, nickname, avatar string) error {
	return client.Update(
		m.DocName(),
		bson.D{
			{
				Key: "user_id", Value: userId,
			},
		},
		bson.D{
			{
				Key: "nickname", Value: nickname,
			},
			{
				Key: "avatar", Value: avatar,
			},
		},
	)
}

func (m *UserInfoDoc) Delete(userId string) error {
	return client.Update(
		m.DocName(),
		bson.D{{Key: "user_id", Value: userId}},
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "delete_at", Value: time.Now().Format(time.RFC3339)},
		}}},
	)
}

func (m *UserInfoDoc) DocName() string { return "lim_user_info" }
