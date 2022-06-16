package mongoc_official

import (
	"context"
	"sync"
	"time"

	"github.com/yunxiaoyang01/open_sdk/logger"
	"github.com/yunxiaoyang01/open_sdk/mongoc_official/mongoc_wrapped"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient uri format: mongodb://username:password@example1.com,example2.com,example3.com/?replicaSet=test&w=majority&wtimeoutMS=5000
func NewClient(ctx context.Context, name, uri string) (*mongoc_wrapped.WrappedClient, error) {
	c, err := mongoc_wrapped.NewClient(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err = c.Connect(ctx); err != nil {
		return nil, err
	}

	if err = c.Ping(ctx, nil); err != nil { // ping before use
		return nil, err
	}

	return c, nil
}

type clientMgr struct {
	lock    sync.RWMutex
	clients map[string]*mongoc_wrapped.WrappedClient
}

func (mgr *clientMgr) NewClient(ctx context.Context, name, uri string) error {
	c, e := NewClient(ctx, name, uri)
	if e != nil {
		return e
	}

	mgr.Add(name, c)

	return e
}

func (mgr *clientMgr) Add(name string, client *mongoc_wrapped.WrappedClient) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	mgr.clients[name] = client
}

func (mgr *clientMgr) Delete(name string) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	delete(mgr.clients, name)
}

func (mgr *clientMgr) Get(name string) *mongoc_wrapped.WrappedClient {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()

	client, ok := mgr.clients[name]
	if !ok {
		return nil
	}

	return client
}

const (
	serverCheckInterval = time.Second * 10
)

// testing
func (mgr *clientMgr) ping() {
	errCount := map[string]int{}

	for {
		ctx := context.Background()
		time.Sleep(serverCheckInterval)

		mgr.lock.RLock()

		clients := map[string]*mongoc_wrapped.WrappedClient{}
		for name, v := range mgr.clients {
			clients[name] = v
		}

		mgr.lock.RUnlock()

		for name, v := range clients {
			if err := v.Ping(ctx, nil); err != nil {
				logger.WithField("name", name).
					Errorf(ctx, "mongo client ping failed, err : %v errCount: %v ", err, errCount[name])
				errCount[name]++

				continue
			}

			// succ ping, then reset errcount to zero
			errCount[name] = 0
		}
	}
}

var _clientsMgr *clientMgr

func mgr() *clientMgr {
	if _clientsMgr == nil {
		_clientsMgr = &clientMgr{}
		_clientsMgr.clients = map[string]*mongoc_wrapped.WrappedClient{}

		go _clientsMgr.ping()
	}
	return _clientsMgr
}

// AddClientAddress 添加一个Mongo链接，指定name和uri. 使用时需要通过name从clientMgr获取
func AddClientAddress(ctx context.Context, name, uri string) error {
	c, e := NewClient(ctx, name, uri)
	if e != nil {
		return e
	}

	mgr().Add(name, c)

	return e
}

func AddClient(name string, client *mongoc_wrapped.WrappedClient) {
	mgr().Add(name, client)
}

func Delete(name string) {
	mgr().Delete(name)
}

func get(name string) *mongoc_wrapped.WrappedClient {
	return mgr().Get(name)
}
