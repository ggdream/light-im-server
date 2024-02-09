package db

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"lim/config"
)

type GroupMemberDoc struct {
	BaseModel `bson:",inline"`
	UserID    string `bson:"user_id" json:"user_id"`
	GroupID   string `bson:"group_id" json:"group_id"`
	Role      int8   `bson:"role" json:"role"`
}

func (m *GroupMemberDoc) Create(ctx context.Context) error {
	m.ID = ksuid.New().String()
	m.CreateAt = time.Now().UnixMilli()
	m.Status = config.RecordStatusNormal

	return client.Client().UseSession(ctx, func(sCtx mongo.SessionContext) (err error) {
		err = sCtx.StartTransaction()
		if err != nil {
			return
		}

		defer func() {
			if err != nil {
				_ = sCtx.AbortTransaction(ctx)
			}
		}()

		one := client.Database().Collection(m.DocName()).FindOne(sCtx, bson.M{
			"group_id": m.GroupID,
			"user_id":  m.UserID,
		})
		err = one.Err()
		if err == nil {
			return nil
		}
		if err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				return err
			}
		}

		_, err = client.Database().Collection(m.DocName()).InsertOne(sCtx, m)
		if err != nil {
			return err
		}

		return sCtx.CommitTransaction(ctx)
	})
}

func (m *GroupMemberDoc) Delete(ctx context.Context) error {
	filter := bson.M{
		"_id": m.ID,
	}

	_, err := client.Database().Collection(m.DocName()).UpdateOne(ctx, filter, bson.M{
		"$set": bson.M{
			"status":    config.RecordStatusDelete,
			"delete_at": time.Now().UnixMilli(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *GroupMemberDoc) Update(ctx context.Context) error {
	filter := bson.M{
		"group_id": m.GroupID,
		"user_id":  m.UserID,
	}

	up := bson.M{}
	if m.Role != 0 {
		up["role"] = m.Role
	}

	up["update_at"] = time.Now().UnixMilli()
	_, err := client.Database().Collection(m.DocName()).UpdateOne(ctx, filter, bson.M{
		"$set": up,
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *GroupMemberDoc) DocName() string {
	return GetGroupMemberDocName()
}

func GetGroupMemberDocName() string {
	return "group_member"
}
