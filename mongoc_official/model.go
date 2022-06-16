package mongoc_official

import (
	"context"
	"time"

	"github.com/yunxiaoyang01/open_sdk/logger"
	"github.com/yunxiaoyang01/open_sdk/mongoc_official/mongoc_wrapped"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model interface {
	DatabaseName() string
	CollectionName() string
}

func NewBaseModel(clientName, db, collection string) *Base {
	return &Base{
		clientName: clientName,
		db:         db,
		coll:       collection,
	}
}

type Base struct {
	clientName string
	db         string
	coll       string
}

func (b *Base) Client() *mongoc_wrapped.WrappedClient {
	return get(b.clientName)
}

func (b *Base) DatabaseName() string {
	return b.db
}

func (b *Base) CollectionName() string {
	return b.coll
}

func (b *Base) UseSessionWithOptions(
	ctx context.Context,
	opts *options.SessionOptions,
	fn func(mongo.SessionContext) error,
) error {
	return b.Client().UseSessionWithOptions(ctx, opts, fn)
}

func (b *Base) cname(ctx context.Context) string {
	return b.coll
}

func (b *Base) C(ctx context.Context) *mongoc_wrapped.WrappedCollection {
	return b.Client().Database(b.DatabaseName()).Collection(b.cname(ctx))
}

func (b *Base) c(ctx context.Context) *mongo.Collection {

	return b.Client().Database(b.DatabaseName()).Collection(b.CollectionName()).Collection()
}

func (b *Base) Aggregate(
	ctx context.Context,
	pipeline, results interface{},
	opts ...*options.AggregateOptions,
) error {
	cur, err := b.c(ctx).Aggregate(ctx, pipeline, opts...)
	if err == nil {
		err = cur.All(ctx, results)
	}

	return err
}

func (b *Base) BulkWrite(
	ctx context.Context,
	models []mongo.WriteModel,
	opts ...*options.BulkWriteOptions,
) (*mongo.BulkWriteResult, error) {

	bwres, err := b.c(ctx).BulkWrite(ctx, models, opts...)
	return bwres, err
}

func (b *Base) Clone(opts ...*options.CollectionOptions) (*mongo.Collection, error) {
	return b.c(context.Background()).Clone(opts...)
}

func (b *Base) Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	count, err := b.c(ctx).CountDocuments(ctx, filter, opts...)
	return count, err
}

func (b *Base) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {

	count, err := b.c(ctx).CountDocuments(ctx, filter, opts...)
	return count, err
}

func (b *Base) Database() *mongoc_wrapped.WrappedDatabase {
	return get(b.clientName).Database(b.DatabaseName())
}

func (b *Base) DeleteMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.DeleteOptions,
) (*mongo.DeleteResult, error) {
	dmres, err := b.c(ctx).DeleteMany(ctx, filter, opts...)
	return dmres, err
}

func (b *Base) DeleteOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.DeleteOptions,
) (*mongo.DeleteResult, error) {

	dor, err := b.c(ctx).DeleteOne(ctx, filter, opts...)
	return dor, err
}

func (b *Base) Distinct(
	ctx context.Context,
	fieldName string,
	filter interface{},
	opts ...*options.DistinctOptions,
) ([]interface{}, error) {

	distinct, err := b.c(ctx).Distinct(ctx, fieldName, filter, opts...)
	return distinct, err
}

func (b *Base) Drop(ctx context.Context) error {

	err := b.c(ctx).Drop(ctx)
	return err
}

func (b *Base) EstimatedDocumentCount(
	ctx context.Context,
	opts ...*options.EstimatedDocumentCountOptions,
) (int64, error) {

	count, err := b.c(ctx).EstimatedDocumentCount(ctx, opts...)
	return count, err
}

func (b *Base) Find(ctx context.Context, filter interface{}, results interface{}, opts ...*options.FindOptions) error {

	cur, err := b.c(ctx).Find(ctx, filter, opts...)
	if err == nil {
		err = cur.All(ctx, results)
	}

	return err
}

func (b *Base) FindOne(
	ctx context.Context,
	filter interface{}, result interface{},
	opts ...*options.FindOneOptions,
) error {

	return b.c(ctx).FindOne(ctx, filter, opts...).Decode(result)
}

func (b *Base) FindOneAndDelete(
	ctx context.Context,
	filter interface{}, result interface{},
	opts ...*options.FindOneAndDeleteOptions,
) error {

	return b.c(ctx).FindOneAndDelete(ctx, filter, opts...).Decode(result)
}

func (b *Base) FindOneAndReplace(
	ctx context.Context,
	filter, replacement, result interface{},
	opts ...*options.FindOneAndReplaceOptions,
) error {

	return b.c(ctx).FindOneAndReplace(ctx, filter, replacement, opts...).Decode(result)
}

func (b *Base) FindOneAndUpdate(
	ctx context.Context,
	filter, update, result interface{},
	opts ...*options.FindOneAndUpdateOptions,
) error {

	return b.c(ctx).FindOneAndUpdate(ctx, filter, update, opts...).Decode(result)
}

func (b *Base) FindOneAndUpsert(
	ctx context.Context,
	filter, update, result interface{},
	opts ...*options.FindOneAndUpdateOptions,
) error {

	rd := options.After
	optUpsert := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(rd)
	opts = append(opts, optUpsert)

	return b.c(ctx).FindOneAndUpdate(ctx, filter, update, opts...).Decode(result)
}

func (b *Base) Indexes() mongo.IndexView { return b.c(context.Background()).Indexes() }

func (b *Base) InsertMany(
	ctx context.Context,
	documents []interface{},
	opts ...*options.InsertManyOptions,
) (*mongo.InsertManyResult, error) {

	insmres, err := b.c(ctx).InsertMany(ctx, documents, opts...)
	return insmres, err
}

func (b *Base) InsertOne(
	ctx context.Context,
	document interface{},
	opts ...*options.InsertOneOptions,
) (*mongo.InsertOneResult, error) {

	insores, err := b.c(ctx).InsertOne(ctx, document, opts...)
	return insores, err
}

func (b *Base) Name() string { return b.CollectionName() }

func (b *Base) ReplaceOne(
	ctx context.Context,
	filter, replacement interface{},
	opts ...*options.ReplaceOptions,
) (*mongo.UpdateResult, error) {

	repres, err := b.c(ctx).ReplaceOne(ctx, filter, replacement, opts...)
	return repres, err
}

func (b *Base) UpdateMany(ctx context.Context, filter, replacement interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {

	umres, err := b.c(ctx).UpdateMany(ctx, filter, replacement, opts...)
	return umres, err
}

func (b *Base) UpdateOne(
	ctx context.Context,
	filter, replacement interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {

	uores, err := b.c(ctx).UpdateOne(ctx, filter, replacement, opts...)
	return uores, err
}

func (b *Base) Upsert(
	ctx context.Context,
	filter, replacement interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {

	optUpsert := options.Update().SetUpsert(true)
	opts = append(opts, optUpsert)
	uores, err := b.c(ctx).UpdateOne(ctx, filter, replacement, opts...)
	return uores, err
}

func (b *Base) Watch(
	ctx context.Context,
	pipeline interface{},
	opts ...*options.ChangeStreamOptions,
) (*mongo.ChangeStream, error) {

	cs, err := b.c(ctx).Watch(ctx, pipeline, opts...)
	return cs, err
}

type AckType int

const (
	ResumeAfter AckType = 1
	Timestamp   AckType = 2
	StartAfter  AckType = 3
)

type StreamEvent interface {
	AckInfo() (typ AckType, val interface{})
	EventOpType() string
	Handle(ctx context.Context) error
}

type Decoder interface {
	Decode(stream *mongo.ChangeStream) (StreamEvent, error)
}

func (b *Base) reWatch(ctx context.Context, ackType AckType, point interface{}, opt *options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	opt.SetStartAfter(nil)
	opt.SetResumeAfter(nil)
	opt.SetStartAtOperationTime(nil)

	switch ackType {
	case StartAfter:
		opt.SetStartAfter(point)
	case ResumeAfter:
		opt.SetResumeAfter(point)
	case Timestamp:
		t, ok := point.(*primitive.Timestamp)
		if ok {
			opt.SetStartAtOperationTime(t)
		}
	}
	return b.Watch(ctx, mongo.Pipeline{}, opt)
}

const (
	watchErrSleep = time.Second * 3
)

func (b *Base) Sub(ctx context.Context, pipeline interface{}, opt *options.ChangeStreamOptions, dec Decoder) error {
	stream, err := b.Watch(ctx, mongo.Pipeline{}, opt)
	if err != nil {
		logger.Fatal(ctx, err)
	}

	defer func() {
		if stream != nil {
			stream.Close(context.Background())
		}
	}()

	var ackType AckType

	var point interface{}

	logger.Info(ctx, "start watching mongo changeStream")

	reconnect := func() {
		stream.Close(ctx)
		stream, err = b.reWatch(ctx, ackType, point, opt)
		if err != nil {
			logger.Errorf(ctx, "try mongo stream watch failed, err : %v", err)
			time.Sleep(watchErrSleep)
		} else {
			logger.Info(ctx, "try watching mongo changeStream ok !")
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			ok := stream.Next(ctx)
			if ok {
				event, err := dec.Decode(stream)
				if err != nil {
					continue
				}

				err = b.dispatchEvent(event)
				if err != nil {
					continue
				}

				if event.EventOpType() == "invalidate" {
					reconnect()
					continue
				}

				ackType, point = event.AckInfo()

				logger.Debugf(ctx, "get event: %+v", event)
			} else { // 异常处理
				logger.Errorf(ctx, "mongo stream.Next failed, err : %v", stream.Err())
				reconnect()
			}
		}
	}
}

func (b *Base) dispatchEvent(event StreamEvent) (err error) {
	ctx := context.Background()
	return event.Handle(ctx)
}
