# config

config 通过对 [viper](github.com/spf13/viper) 进行简单封装，解决每一个项目都需要解析配置文件的问题。

## 安装

```shell
go get -u github.com/yunxiaoyang01/open_sdk/config
```
## 加载配置文件

```golang
package main

import (
    "github.com/yunxiaoyang01/open_sdk/config"
)

// 一般在config目录下定义业务的配置结构体。
type Config struct {
    Port int
}

func main() {
    var c Config

    if c, err := config.Load("./demo.config"); err != nil {
        fmt.Println(err)
        return
    }
}

```

## 添加勾子函数

config 加载配置文件时，可以设置勾子函数进行一些初始化操作。

```golang
package main

import (
    "github.com/yunxiaoyang01/open_sdk/config"
)

type Config struct {
    Port int
}

func NewDemoHook() Hook {
	return func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		if from.Kind() != reflect.Map {
			return data, nil
		}
		if to == reflect.TypeOf(logger.Config{}) {
			var config logger.Config
			err := mapstructure.Decode(data, &config.Options)
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode data")
			}
			config.Logger, err = logger.NewLoggerWithOptions(config.Options)
			if err != nil {
				return nil, errors.Wrap(err, "failed to new logger")
			}
			return config, nil
		}
		return data, nil
	}
}

func main() {
    var c Config

    if err := config.LoadWithHooks("./demo.config", &c, NewDemoHook); err != nil {
        fmt.Println(err)
        return
    }
}
```

其中，config 默认添加了对基础库的hook，对以下这些组件进行了初始化，用户可以直接方便地使用。
```golang
var defaultHooks = []Hook{
	NewLoggerHook(),  // logger
	NewMongocHook(),  // mongo
	NewRediscHook(),  // redis
	NewElasticsearchcHook(),  // es
	NewTikvcHook(),  // tikv
}
```

## 配置结构体和配置内容示例
业务开发中，常用的配置结构体和配置内容示例
```golang
// 全局配置结构定义。
type Config struct {
	GrpcListen  string          `mapstructure:"grpc_listen"`
	HTTPListen  string          `mapstructure:"http_listen"`
	ProfPort    int             `mapstructure:"prof_port"`
	StdLog      logger.Standard `mapstructure:"std_log"`
	APILog      logger.Config   `mapstructure:"api_log"`
	OCTrace     config.Config   `mapstructure:"oc_trace"`
	ServiceName string          `mapstructure:"service_name"`
	Account     account.Config  `mapstructure:"account"`
	Ping        ping.Config     `mapstructure:"ping"`
	MongoURI    string          `mapstructure:"mongo_uri"`
}
```

配置文件内容(XX.toml)，与上边结构体对应的toml格式示例配置文件如下：
```
prof_port = 8057

http_listen = "127.0.0.1:8080"

mongo_uri = "mongodb://127.0.0.1:27017/?replicaSet=test"


# 标准日志配置
[std_log]
    level = "debug"
    file = "./ec_template.log"
    err_file = "./ec_template.err.log"

# 接口日志配置
[api_log]
    level = "debug"
    file = "./ec_template.api.log"

[oc_trace]
	sample_type         = 0
	fraction            = 10.0
	agent_end_point     = "localhost:6831"
	service_name        ="ec_template"
    
[account]
    name_prefix = "test"

[ping]
    msg_prefix = "test"

```