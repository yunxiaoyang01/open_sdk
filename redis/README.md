# redis

本库为对[go-redis](https://github.com/go-redis/redis)的封装，提供了一些列易用的redis方法

## Get started

### 初始化

#### 1. 定义配置文件初始化

定义以下配置
```json
{
    "redis": {
        "address": "localhost:6379"
    },
    "redis_cluster": {
        "startup_nodes": ["redis-cluster:6379"]
    },
}
```

使用以下代码
```go
type Config struct {
    Redis redis.Config `mapstructure:"redis"`
    RedisCluster redis.ClusterConfig `mapstructure:"redis_cluster"`
}
---
var c Config
config.Load("PATH",&c)

c.Ping(ctx)
```

#### 2. 直接初始化

```go
client, err := redis.NewClient("localhost:6379",0,WithPassword("password"))
if err !=nil {
    panic(err)
}
client.Ping(ctx)

clusterClient, err := redis.NewClusterClient([]string{"cluster:6379"},WithPassword("password"))
if err !=nil {
    panic(err)
}
clusterClient.Ping(ctx)
```

### 使用

```go
func ExampleNewClient() {
	ctx := context.Background()
	client, err := NewClient("localhost:6379", 0)
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
	shadowCtx := metadata.WithMetadata(ctx, metadata.Metadata{
		"flow_color": "pt",
	})
	if err := client.Set(shadowCtx, "something", "someone", 0).Err(); err != nil {
		panic(err)
	}
	keys, err := client.Keys(shadowCtx, "shadow*").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(keys)
	count, err = client.Del(shadowCtx, "something").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(count)
	// Output:
	// PONG
	// 1
	// [shadow_something]
	// 1
}
```