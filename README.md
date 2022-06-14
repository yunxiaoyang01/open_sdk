# open_sdk

```
open_sdk 对通用库做统一的封装和管理，集成通用中间件等。
具体库使用示例见各自README、example目录以及对应的test测试用例文件。
```

## 目录结构
每个库单独目录，见下方目录树：

```Shell
tree -d -L 1 ./
./
├── config  # 配置解析
├── es_official  # ES ElasticSearch，对“olivere”版本库的封装
├── example  # 一些库使用示例
├── grpcclient
├── grpcserver
├── httpclient
├── httpserver
├── logger  # Logger 统一日志库
├── metadata
├── metrics
├── mongoc  # 对mgo的封装，由于mgo已暂停维护，【不建议使用此库】
├── mongoc_official # 对mongo官方的golang库的封装。【推荐使用此库】
├── octrace
├── pubsub
├── redisc # Redis
└── xd_error
```

各层级目录信息如下：

```Shell
tree -d ./
./
├── config
├── es_official
│   └── es_wrapped
├── example
│   ├── grpc
│   │   ├── client
│   │   └── server
│   ├── mongo_official
│   ├── pubsub
│   │   ├── pulsar
│   │   └── redispubsub
│   └── redis
├── grpcclient
│   └── hooks
├── grpcserver
│   └── middles
├── httpclient
│   └── hooks
├── httpserver
│   ├── middles
│   └── status
├── logger
├── metadata
├── metrics
│   └── internal
│       └── lv
├── mongoc
├── mongoc_official
│   └── mongoc_wrapped
├── octrace
│   ├── config
│   ├── examples
│   │   ├── grpc
│   │   │   └── proto
│   │   ├── http
│   │   └── stdhttp
│   ├── filter
│   ├── integration
│   │   ├── grpc
│   │   └── http
│   └── propagation
│       └── textheader
├── pubsub
│   ├── pulsar
│   └── redisbroadcast
├── redisc
└── xd_error

```

