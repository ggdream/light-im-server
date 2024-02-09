package db

import (
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"lim/config"
)

type UserInfoDoc struct {
	BaseAutoIDModel `bson:",inline"`
	UserID          string `json:"user_id" bson:"user_id"`
	Username        string `json:"username" bson:"username"`
	Avatar          string `json:"avatar" bson:"avatar"`
}

func (m *UserInfoDoc) Create(userId, username, avatar string) error {
	err := client.SearchOne(m.DocName(), bson.M{
		"user_id": userId,
		"status": bson.M{
			"$ne": config.RecordStatusDelete,
		},
	}, nil).Err()
	if err == nil {
		return errors.New("用户已存在")
	}
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
	}

	m.UserID = userId
	m.Username = username
	m.Avatar = avatar
	m.CreateAt = time.Now().UnixMilli()
	m.Status = config.RecordStatusNormal

	return client.Insert(m.DocName(), m)
}

func (m *UserInfoDoc) Search(userId string) error {
	return client.SearchOne(m.DocName(), bson.M{
		"user_id": userId,
		"status": bson.M{
			"$ne": config.RecordStatusDelete,
		},
	}, nil).Decode(m)
}

func (m *UserInfoDoc) Update(userId, username, avatar string) error {
	return client.Update(
		m.DocName(),
		bson.D{
			{
				Key: "user_id", Value: userId,
			},
		},
		bson.M{
			"$set": bson.M{
				"username":  username,
				"avatar":    avatar,
				"update_at": time.Now().UnixMilli(),
			},
		},
	)
}

func (m *UserInfoDoc) Delete(userId string) error {
	return client.Update(
		m.DocName(),
		bson.D{{Key: "user_id", Value: userId}},
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "delete_at", Value: time.Now().UnixMilli()},
		}}},
	)
}

func (m *UserInfoDoc) DocName() string { return "lim_user_info" }
