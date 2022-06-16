package mongoc_wrapped

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WrappedSession struct {
	mongo.Session
}

var _ mongo.Session = (*WrappedSession)(nil)

func (ws *WrappedSession) EndSession(ctx context.Context) {

	ws.Session.EndSession(ctx)
}

func (ws *WrappedSession) StartTransaction(topts ...*options.TransactionOptions) error {
	return ws.Session.StartTransaction(topts...)
}

func (ws *WrappedSession) AbortTransaction(ctx context.Context) error {

	err := ws.Session.AbortTransaction(ctx)
	return err
}

func (ws *WrappedSession) CommitTransaction(ctx context.Context) error {

	err := ws.Session.CommitTransaction(ctx)
	return err
}

func (ws *WrappedSession) ClusterTime() bson.Raw {
	return ws.Session.ClusterTime()
}

func (ws *WrappedSession) AdvanceClusterTime(br bson.Raw) error {
	return ws.Session.AdvanceClusterTime(br)
}

func (ws *WrappedSession) OperationTime() *primitive.Timestamp {
	return ws.Session.OperationTime()
}

func (ws *WrappedSession) AdvanceOperationTime(pt *primitive.Timestamp) error {
	return ws.Session.AdvanceOperationTime(pt)
}
