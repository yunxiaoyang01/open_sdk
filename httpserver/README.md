# httpserver

`httpserver`封装了 [gin](github.com/gin-gonic/gin)
用户可以方便地启动一个 `http` 服务端。

## 安装

```shell
go get -u github.com/yunxiaoyang01/open_sdk/httpserver
```

## 初始化

```golang

	apiLogger, _ := logger.NewLogger(logger.WithFile("./api_tttt.log"))

	s := httpserver.NewServer(
		httpserver.WithAddress(":8080"),
		httpserver.WithEnableStats(true), // 开启metrics
		httpserver.WithMiddles(
			middles.LoggingRequestWithLogger(apiLogger), // 指定logger
			middles.LoggingResponseWithLogger(apiLogger),
		),
	)

	go metrics.Start(context.Background()) // 测试metrics,默认10001端口，通过 http://127.0.0.1:10001/metrics 查看统计值

...

```

## 添加路由

```golang
v1 := s.GetKernel().Group("/v1")
v1.GET("hello", httpserver.CreateHandlerFunc(handlers, "hello"))
```

## 添加中间层

```golang
// 已经默认集成：flowControlTag、错误恢复，无需再设置。
// 可自定义logger打印位置
apiLogger, _ := logger.NewLogger(logger.WithFile("./api_tttt.log"))
httpserver.WithMiddles(
	middles.LoggingRequestWithLogger(apiLogger), // 指定logger
	middles.LoggingResponseWithLogger(apiLogger),
),
```

## 获取服务名称

```golang
name := s.Name()
```

## 启动服务

```golang
ctx := context.Background()
s.Run(ctx)
```

## 启动Metrics
```golang
	s := httpserver.NewServer(
		httpserver.WithAddress(":8081"),
		httpserver.WithEnableStats(true), // 开启metrics
	)

	go metrics.Start(context.Background()) // 测试metrics,默认10001端口，通过 http://127.0.0.1:10001/metrics 查看统计值

```