package db

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"lim/config"
)

type GroupInfoDoc struct {
	BaseModel   `bson:",inline"`
	UserID      string   // 由谁创建
	Name        string   // 名称
	Desc        string   // 介绍
	Avatar      string   // 头像
	Tags        []string // 标签
	Capacity    uint64   // 容量
	MemberCount int64    // 成员数
}

func (m *GroupInfoDoc) First(ctx context.Context) error {
	one := client.Database().Collection(m.DocName()).FindOne(ctx, bson.M{"_id": m.ID})
	if err := one.Err(); err != nil {
		return err
	}

	return one.Decode(m)
}

func (m *GroupInfoDoc) Create(ctx context.Context) error {
	m.ID = ksuid.New().String()
	m.CreateAt = time.Now().UnixMilli()
	m.Status = config.RecordStatusNormal

	err := client.Client().UseSession(ctx, func(sCtx mongo.SessionContext) (err error) {
		err = sCtx.StartTransaction()
		if err != nil {
			return
		}

		defer func() {
			if err != nil {
				_ = sCtx.AbortTransaction(ctx)
			}
		}()

		_, err = client.Database().Collection(m.DocName()).InsertOne(sCtx, m)
		if err != nil {
			return err
		}

		m1 := &GroupMemberDoc{
			UserID:  m.UserID,
			GroupID: m.ID,
			Role:    config.GroupMemberRoleOwner,
		}
		err = m1.Create(sCtx)
		if err != nil {
			return err
		}

		return sCtx.CommitTransaction(ctx)
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *GroupInfoDoc) Delete(ctx context.Context, isAdmin bool) error {
	filter := bson.M{
		"_id": m.ID,
	}
	if !isAdmin {
		filter["user_id"] = m.UserID
		filter["status"] = config.RecordStatusNormal
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

func (m *GroupInfoDoc) Update(ctx context.Context, isAdmin bool) error {
	filter := bson.M{
		"_id": m.ID,
	}
	if !isAdmin {
		filter["user_id"] = m.UserID
		filter["status"] = config.RecordStatusNormal
	}

	up := bson.M{}
	inc := bson.M{}
	if m.Name != "" {
		up["name"] = m.Name
	}
	if m.Desc != "" {
		up["desc"] = m.Desc
	}
	if m.Avatar != "" {
		up["avatar"] = m.Avatar
	}
	if m.Tags != nil {
		up["tags"] = m.Tags
	}
	if m.Status != 0 {
		up["status"] = m.Status
	}
	if m.Capacity != 0 {
		up["capacity"] = m.Capacity
	}
	if m.MemberCount != 0 {
		inc["member_count"] = m.MemberCount
	}

	if len(up) == 0 && len(inc) == 0 {
		return nil
	}

	up["update_at"] = time.Now().UnixMilli()
	_, err := client.Database().Collection(m.DocName()).UpdateOne(ctx, filter, bson.M{
		"$set": up,
		"$inc": inc,
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *GroupInfoDoc) DocName() string {
	return GetGroupInfoDocName()
}

func GetGroupInfoDocName() string {
	return "group_info"
}
