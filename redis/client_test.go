package redis

import (
	"context"
	"fmt"
	"github.com/yunxiaoyang01/open_sdk/logger"
	"github.com/yunxiaoyang01/open_sdk/public"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	type args struct {
		address string
		db      int
		opts    []Options
	}
	tests := []struct {
		name    string
		args    args
		want    Client
		wantErr bool
	}{
		{
			name: "test normal",
			args: args{
				address: "redis:6379",
				opts: []Options{
					WithConnectTimeout(1),
					WithMaxIdle(1),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.args.address, tt.args.db, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func ExampleNewClient() {
	ctx := context.Background()
	client, err := NewClient(os.Getenv("REDIS_ADDRESS"), 0)
	if err != nil {
		panic(err)
	}
	result, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	if err := client.Set(ctx, "something", "someone", 0).Err(); err != nil {
		panic(err)
	}
	count, err := client.Del(ctx, "something").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(count)
	// Output:
	// PONG
	// 1
}

func TestNewClusterClient(t *testing.T) {
	ctx := context.Background()
	client, err := NewClusterClient(strings.Split(os.Getenv("CLUSTER_NODES"), ","), WithConnectTimeout(1))
	if err != nil {
		t.Error("failed to new client", err)
		return
	}
	result, err := client.Ping(ctx).Result()
	if err != nil {
		t.Error("failed to ping cluster", err)
	}
	t.Logf("ping result is %s", result)
	if err := client.Set(ctx, "something", "someone", 0).Err(); err != nil {
		t.Error("failed to Set", err)
	}
	count, err := client.Del(ctx, "something").Result()
	if err != nil {
		t.Error("failed to Del", err)
	}
	t.Logf("delete count %d", count)
}

func ExamplePubSub() {
	ctx := context.Background()
	rdb, err := NewClient(os.Getenv("REDIS_ADDRESS"), 0)
	if err != nil {
		panic(err)
	}
	pubsub := rdb.Subscribe(ctx, "mychannel1")

	// Wait for confirmation that subscription is created before publishing anything.
	response, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(response)
	// Go channel which receives messages.
	ch := pubsub.Channel() // this will automaticly ping

	// Publish a message.
	err = rdb.Publish(ctx, "mychannel1", "hello").Err()
	if err != nil {
		panic(err)
	}

	time.AfterFunc(time.Second, func() {
		// When pubsub is closed channel is closed too.
		_ = pubsub.Close()
	})

	// Consume messages.
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
	}
	// Output:
	// subscribe: mychannel1
	// mychannel1 hello
}

func TestNewClientSetName(t *testing.T) {
	logger.SetLevel(logger.DebugLevel)
	svcName := "xd_sdk"
	hostName := "localhost"
	os.Setenv(public.EnvSvcName, svcName)
	os.Setenv(public.EnvHostName, hostName)
	cli, err := NewClient(os.Getenv("REDIS_ADDRESS"), 0)
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	cliName := cli.ClientGetName(ctx).Val()
	t.Log("clientName: ", cliName)
	if cliName != svcName+":"+hostName {
		t.Fail()
	}
}
