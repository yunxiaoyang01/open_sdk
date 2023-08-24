package public

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	callerInfo = GetCallingInfoFromOS()
)

const (
	HeaderSvcName  string = "Header-SvcName"  // 头信息中传递 svcName 的字段
	HeaderHostName string = "Header-HostName" // 头信息中传递 hostName 的字段
)

// CallerInfo 调用本服务的服务名、主机名
type CallerInfo struct {
	SvcName  string
	HostName string
}

type callerInfoKey struct{}

var ctxKeyCallerInfo = callerInfoKey{}

// GinParseCallInfo 是处理掉用关系的 gin 中间件
// 会从 HTTP Header 中读取掉用本服务的服务名和主机名，写入 ctx
func GinParseCallInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		info := CallerInfo{
			HostName: c.Request.Header.Get(HeaderHostName),
			SvcName:  c.Request.Header.Get(HeaderSvcName),
		}
		if info.HostName != "" && info.SvcName == "" {
			vv := strings.Split(info.HostName, "-deploy")
			if len(vv) > 0 {
				info.SvcName = vv[0]
			}
		}
		ctx = context.WithValue(ctx, ctxKeyCallerInfo, info)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetCallInfoFromCtx 从 ctx 中获取 CallerInfo
func GetCallInfoFromCtx(ctx context.Context) (info CallerInfo) {
	v := ctx.Value(ctxKeyCallerInfo)
	if vv, ok := v.(CallerInfo); ok {
		info = vv
	}
	return
}

func GetCallingInfo() (info *CallerInfo) {
	return &callerInfo
}

// GetCallingInfoFromOS 获本该服务的基本信息
func GetCallingInfoFromOS() (info CallerInfo) {
	info = CallerInfo{
		HostName: os.Getenv(EnvHostName),
		SvcName:  os.Getenv(EnvSvcName),
	}
	if info.HostName != "" && info.SvcName == "" {
		vv := strings.Split(info.HostName, "-deploy-")
		if len(vv) > 1 {
			info.SvcName = vv[0]
		}
	}
	return info
}

func GetCallingInfoFromOSStr() string {
	info := GetCallingInfoFromOS()
	bs, err := json.Marshal(info)
	if err != nil {
		return ""
	}
	return string(bs)
}
func UnmarshalCallingInfo(bs string) (info CallerInfo) {
	json.Unmarshal([]byte(bs), &info)
	return
}

// SetCallingInfoToCtx 将本服务的 CallerInfo 设置到 ctx，传递到被本服务调用的其他服务
func SetCallingInfoToCtx(ctx context.Context, info CallerInfo) (ctxOut context.Context) {
	return context.WithValue(ctx, ctxKeyCallerInfo, info)
}
