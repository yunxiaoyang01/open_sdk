package redis

import (
	"context"
	"errors"
	"strings"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/yunxiaoyang01/open_sdk/logger"
	"github.com/yunxiaoyang01/open_sdk/public"
)

// Client client interface
type Client interface {
	goredis.UniversalClient
}

// NewClient creates a redis client
func NewClient(address string, db int, opts ...Options) (Client, error) {
	option := defaultOption()
	for _, opt := range opts {
		opt(option)
	}
	goclient := goredis.NewClient(&goredis.Options{
		Addr:         address,
		DB:           db,
		Username:     option.Username,
		Password:     option.Password,
		DialTimeout:  time.Duration(option.ConnectTimeout) * time.Second,
		ReadTimeout:  time.Duration(option.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(option.WriteTimeout) * time.Second,
		PoolSize:     option.MaxIdle,
		OnConnect: func(ctx context.Context, cn *goredis.Conn) error {
			info := public.GetCallingInfoFromOS()
			if clientName := info.SvcName + ":" + info.HostName; clientName != ":" {
				err := cn.ClientSetName(ctx, info.SvcName+":"+info.HostName).Err()
				if err != nil {
					logger.Warningf(ctx, "redis client set clientName err %v", err)
				}
				return nil
			}
			return nil
		},
	})
	return goclient, goclient.Ping(context.Background()).Err()
}

// NewClusterClient creates a new redis cluster client.
func NewClusterClient(nodes []string, opts ...Options) (Client, error) {
	option := defaultOption()
	for _, opt := range opts {
		opt(option)
	}
	clusterClient := goredis.NewClusterClient(&goredis.ClusterOptions{
		Addrs:        nodes,
		Username:     option.Username,
		Password:     option.Password,
		DialTimeout:  time.Duration(option.ConnectTimeout) * time.Second,
		ReadTimeout:  time.Duration(option.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(option.WriteTimeout) * time.Second,
		PoolSize:     option.MaxIdle,
		OnConnect: func(ctx context.Context, cn *goredis.Conn) error {
			info := public.GetCallingInfoFromOS()
			if clientName := info.SvcName + ":" + info.HostName; clientName != ":" {
				err := cn.ClientSetName(ctx, info.SvcName+":"+info.HostName).Err()
				if err != nil {
					logger.Warningf(ctx, "redis cluster client set clientName err %v", err)
				}
				return nil
			}
			return nil
		},
	})
	return clusterClient, clusterClient.Ping(context.Background()).Err()
}

// NewClusterClientWithIPCheck 在 NewClusterClient 的基础上增加检查所配置的 node 是否都是连接集群中的节点
// 注意：nodes 必须是集群节点的 ip:port，如果是域名或者代理的地址，会报错
func NewClusterClientWithIPCheck(nodes []string, opts ...Options) (Client, error) {
	client, err := NewClusterClient(nodes, opts...)
	if err != nil {
		return nil, err
	}
	ctx := client.Context()
	nodesInfo := client.ClusterNodes(ctx).Val()
	for _, n := range nodes {
		if !strings.Contains(nodesInfo, n) {
			err := errors.New("请求IP不在集群中：" + n)
			logger.Warnf(ctx, "NewClusterClient ClusterNodesInfo:%s fail:%s", nodesInfo, err)
			return nil, err
		}
	}
	return client, nil
}
