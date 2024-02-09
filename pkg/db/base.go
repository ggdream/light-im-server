package db

type BaseModel struct {
	ID       string `bson:"_id" json:"id"`
	CreateAt int64  `bson:"create_at" json:"create_at"`
	UpdateAt int64  `bson:"update_at" json:"update_at"`
	DeleteAt int64  `bson:"delete_at" json:"delete_at"`
	Status   int8   `bson:"status" json:"status"`
}

type BaseAutoIDModel struct {
	ID       string `bson:"-" json:"id"`
	CreateAt int64  `bson:"create_at" json:"create_at"`
	UpdateAt int64  `bson:"update_at" json:"update_at"`
	DeleteAt int64  `bson:"delete_at" json:"delete_at"`
	Status   int8   `bson:"status" json:"status"`
}
