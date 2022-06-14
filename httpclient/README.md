# httpclient

`httpclient`封装了 [resty](gopkg.in/resty.v1), 为用户提供方便使用的 http 请求发送。

## 安装

```shell
go get -u github.com/yunxiaoyang01/open_sdk/httpclient
```

## 初始化

```golang
c := httpclient.NewClient(
        httpclient.WithAddress(s.config.Address),
        httpclient.WithTimeout(3 * time.Second),
        httpclient.WithPreRequestHooks(httpclient.hooks.LoggingRequest()),
        httpclient.WithAfterResponseHooks(httpclient.hooks.LoggingResponse()),
    ),
)
```

## 获取客户端

```golang
client := c.GetKernel()
```

## 发送请求

```golang
// get
response, err := c.NewRequest().Get("http://127.0.0.1:/v1/hello")

// post
response, err := c.NewRequest().Post("http://127.0.0.1:/v1/hello")
```
更多接口见 [Request](https://godoc.org/gopkg.in/resty.v1#Request)