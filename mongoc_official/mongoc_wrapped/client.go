package mongoc_wrapped

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type WrappedClient struct {
	cc *mongo.Client
}

func NewClient(ctx context.Context, opts ...*options.ClientOptions) (*WrappedClient, error) {

	client, err := mongo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return &WrappedClient{cc: client}, nil
}

func (wc *WrappedClient) Connect(ctx context.Context) error {

	err := wc.cc.Connect(ctx)
	return err
}

func (wc *WrappedClient) Database(name string, opts ...*options.DatabaseOptions) *WrappedDatabase {
	db := wc.cc.Database(name, opts...)
	if db == nil {
		return nil
	}

	return &WrappedDatabase{db: db}
}

func (wc *WrappedClient) Disconnect(ctx context.Context) error {

	err := wc.cc.Disconnect(ctx)

	return err
}

func (wc *WrappedClient) ListDatabaseNames(
	ctx context.Context,
	filter interface{},
	opts ...*options.ListDatabasesOptions,
) ([]string, error) {

	dbs, err := wc.cc.ListDatabaseNames(ctx, filter, opts...)

	return dbs, err
}

func (wc *WrappedClient) ListDatabases(
	ctx context.Context,
	filter interface{},
	opts ...*options.ListDatabasesOptions,
) (mongo.ListDatabasesResult, error) {
	dbr, err := wc.cc.ListDatabases(ctx, filter, opts...)

	return dbr, err
}

func (wc *WrappedClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {

	err := wc.cc.Ping(ctx, rp)

	return err
}

func (wc *WrappedClient) StartSession(ctx context.Context, opts ...*options.SessionOptions) (mongo.Session, error) {

	ss, err := wc.cc.StartSession(opts...)
	if err != nil {
		return nil, err
	}

	return &WrappedSession{Session: ss}, nil
}

func (wc *WrappedClient) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	return wc.cc.UseSession(ctx, fn)
}

func (wc *WrappedClient) UseSessionWithOptions(
	ctx context.Context,
	opts *options.SessionOptions,
	fn func(mongo.SessionContext) error,
) error {
	return wc.cc.UseSessionWithOptions(ctx, opts, fn)
}

func (wc *WrappedClient) Client() *mongo.Client { return wc.cc }
