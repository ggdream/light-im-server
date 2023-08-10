package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"lim/config"
)

var client *Mongo

type Mongo struct {
	client   *mongo.Client
	database string
	context  context.Context
	timeout  time.Duration
}

func New(cl *mongo.Client) *Mongo {
	return &Mongo{
		client:   cl,
		database: "lim",
		context:  context.TODO(),
		timeout:  time.Second * 1,
	}
}

func Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	cfg := config.GetMongo()
	opt := options.Client().ApplyURI(cfg.URI)
	cli, err := mongo.Connect(ctx, opt)
	if err != nil {
		return err
	}

	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel1()
	err = cli.Ping(ctx1, readpref.Primary())
	if err != nil {
		return err
	}

	client = &Mongo{
		client:   cli,
		database: cfg.Database,
		context:  context.TODO(),
		timeout:  time.Second * 1,
	}

	return nil
}

func Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return client.client.Disconnect(ctx)
}

func (m *Mongo) SearchOne(collection string, filter interface{}, opt *options.FindOneOptions) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(m.context, m.timeout*2)
	defer cancel()

	return m.client.Database(m.database).Collection(collection).FindOne(ctx, filter, opt)
}

func (m *Mongo) SearchMany(collection string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(m.context, m.timeout*2)
	defer cancel()

	return m.client.Database(m.database).Collection(collection).Find(ctx, filter, opts...)
}

func (m *Mongo) MatchSearch(collection string, filter interface{}, opt *options.FindOptions) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(m.context, m.timeout*2)
	defer cancel()

	return m.client.Database(m.database).Collection(collection).Find(ctx, filter, opt)
}

func (m *Mongo) Insert(collection string, document interface{}) error {
	ctx, cancel := context.WithTimeout(m.context, m.timeout)
	defer cancel()

	_, err := m.client.Database(m.database).Collection(collection).InsertOne(ctx, document)
	return err
}

func (m *Mongo) Update(collection string, f, u bson.D) error {
	ctx, cancel := context.WithTimeout(m.context, m.timeout)
	defer cancel()

	_, err := m.client.Database(m.database).Collection(collection).UpdateOne(ctx, f, u)
	return err
}

func (m *Mongo) UpdateMany(collection string, f, u bson.D) error {
	ctx, cancel := context.WithTimeout(m.context, m.timeout)
	defer cancel()

	_, err := m.client.Database(m.database).Collection(collection).UpdateMany(ctx, f, u)
	return err
}

func (m *Mongo) InsertAll(collection string, documents []interface{}) error {
	ctx, cancel := context.WithTimeout(m.context, m.timeout)
	defer cancel()

	_, err := m.client.Database(m.database).Collection(collection).InsertMany(ctx, documents)
	return err
}

func (m *Mongo) DeleteAll(collection string) error {
	ctx, cancel := context.WithTimeout(m.context, m.timeout)
	defer cancel()

	_, err := m.client.Database(m.database).Collection(collection).DeleteMany(ctx, bson.D{})
	return err
}

func (m *Mongo) Delete(collection string, d bson.D) error {
	ctx, cancel := context.WithTimeout(m.context, m.timeout)
	defer cancel()

	_, err := m.client.Database(m.database).Collection(collection).DeleteOne(ctx, d)
	return err
}

func (m *Mongo) Count(collection string, f bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(m.context, m.timeout)
	defer cancel()

	return m.client.Database(m.database).Collection(collection).CountDocuments(ctx, f)
}
