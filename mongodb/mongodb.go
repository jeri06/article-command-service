package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	Connect(ctx context.Context) (err error)
	Database(name string, opts ...*options.DatabaseOptions) (db Database)
	Disconnect(ctx context.Context) (err error)
}

type Database interface {
	Collection(name string, opts ...*options.CollectionOptions) (col Collection)
}

type Collection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (result *mongo.InsertOneResult, err error)
}

type SingleResult interface {
	Decode(v interface{}) error
	Err() error
}

// Cursor is a collection of function of mongodb cursor.
type Cursor interface {
	ID() int64
	Next(ctx context.Context) bool
	Decode(val interface{}) error
	Close(ctx context.Context) error
}
