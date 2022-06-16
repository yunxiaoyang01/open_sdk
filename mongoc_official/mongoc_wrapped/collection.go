package mongoc_wrapped

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WrappedCollection struct {
	coll *mongo.Collection
}

func (wc *WrappedCollection) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {

	cur, err := wc.coll.Aggregate(ctx, pipeline, opts...)
	return cur, err
}

func (wc *WrappedCollection) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {

	bwres, err := wc.coll.BulkWrite(ctx, models, opts...)
	return bwres, err
}

func (wc *WrappedCollection) Clone(opts ...*options.CollectionOptions) (*mongo.Collection, error) {
	return wc.coll.Clone(opts...)
}

func (wc *WrappedCollection) Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	count, err := wc.coll.CountDocuments(ctx, filter, opts...)
	return count, err
}

func (wc *WrappedCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {

	count, err := wc.coll.CountDocuments(ctx, filter, opts...)
	return count, err
}

func (wc *WrappedCollection) Database() *mongo.Database { return wc.coll.Database() }

func (wc *WrappedCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	dmres, err := wc.coll.DeleteMany(ctx, filter, opts...)
	return dmres, err
}

func (wc *WrappedCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {

	dor, err := wc.coll.DeleteOne(ctx, filter, opts...)

	return dor, err
}

func (wc *WrappedCollection) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {

	distinct, err := wc.coll.Distinct(ctx, fieldName, filter, opts...)
	return distinct, err
}

func (wc *WrappedCollection) Drop(ctx context.Context) error {

	err := wc.coll.Drop(ctx)
	return err
}

func (wc *WrappedCollection) EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {

	count, err := wc.coll.EstimatedDocumentCount(ctx, opts...)
	return count, err
}

func (wc *WrappedCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {

	cur, err := wc.coll.Find(ctx, filter, opts...)
	return cur, err
}

func (wc *WrappedCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {

	return wc.coll.FindOne(ctx, filter, opts...)
}

func (wc *WrappedCollection) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	return wc.coll.FindOneAndDelete(ctx, filter, opts...)
}

func (wc *WrappedCollection) FindOneAndReplace(ctx context.Context, filter, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	return wc.coll.FindOneAndReplace(ctx, filter, replacement, opts...)
}

func (wc *WrappedCollection) FindOneAndUpdate(ctx context.Context, filter, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {

	return wc.coll.FindOneAndUpdate(ctx, filter, update, opts...)
}

func (wc *WrappedCollection) Indexes() mongo.IndexView { return wc.coll.Indexes() }

func (wc *WrappedCollection) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {

	insmres, err := wc.coll.InsertMany(ctx, documents, opts...)
	return insmres, err
}

func (wc *WrappedCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {

	insores, err := wc.coll.InsertOne(ctx, document, opts...)
	return insores, err
}

func (wc *WrappedCollection) Name() string { return wc.coll.Name() }

func (wc *WrappedCollection) ReplaceOne(ctx context.Context, filter, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {

	repres, err := wc.coll.ReplaceOne(ctx, filter, replacement, opts...)
	return repres, err
}

func (wc *WrappedCollection) UpdateMany(ctx context.Context, filter, replacement interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {

	umres, err := wc.coll.UpdateMany(ctx, filter, replacement, opts...)
	return umres, err
}

func (wc *WrappedCollection) UpdateOne(ctx context.Context, filter, replacement interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	uores, err := wc.coll.UpdateOne(ctx, filter, replacement, opts...)
	return uores, err
}

func (wc *WrappedCollection) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {

	cs, err := wc.coll.Watch(ctx, pipeline, opts...)
	return cs, err
}

func (wc *WrappedCollection) Collection() *mongo.Collection {
	return wc.coll
}
