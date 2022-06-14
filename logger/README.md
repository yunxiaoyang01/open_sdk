# logger

`logger`封装了 [resty](gopkg.in/resty.v1), 为用户提供方便使用的 http 请求发送。

## 安装

```shell
go get -u github.com/yunxiaoyang01/open_sdk/logger
```

## 初始化

```golang
	_ = logger.ResetStandardWithOptions(logger.Options{
		Level:   "debug",
		File:    "/tmp/test.log",
		ErrFile: "/tmp/test.err.log",
	})
```

## 使用

```golang
logger.Infof(ctx,"hello world :%s","open_sdk")
```