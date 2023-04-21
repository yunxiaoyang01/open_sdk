package middles

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	oError "github.com/yunxiaoyang01/open_sdk/error"
	"github.com/yunxiaoyang01/open_sdk/logger"
	gcodes "google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

// ParseMetricCodeFromErr 自定义解析业务 code 方法，ok为true时使用解析到的 code
// 本函数支持从 err 解析业务 code
// 运行时调用时，调用方会保证 err != nil
// 注意：如果 err 对象不为空（有类型）但是底层值为空时，调用 err.Error() 可能会 panic。
// 自定义 error 需要保证 err 为 nil 时： err.Error() 不会 panic
type ParseMetricCodeFromErr = func(err error) (code int, ok bool)

// ParseMetricCodeFromRsp 自定义解析业务 code 方法，ok为true时使用解析到的 code
// 本函数支持从返回值 rsp 中解析业务 code 的情况（不推荐）
// 运行时调用时，调用方会保证 rsp != nil
type ParseMetricCodeFromRsp = func(rsp interface{}) (code int, ok bool)

// ParseMetricCodeFromBody 是默认从 http body 解析业务 code 的方法，ok为true时使用解析到的 code
type ParseMetricCodeFromBody = func(in []byte) (code int, ok bool)

// ParseMetricCodeFromBodyFunc 是默认从 http body 解析业务 code 的方法，有性能消耗，慎用
var ParseMetricCodeFromBodyFunc ParseMetricCodeFromBody = func(in []byte) (code int, ok bool) {
	if len(in) == 0 {
		return
	}
	result := &RespStruct{}
	if err := json.Unmarshal(in, &result); err != nil {
		return
	}
	return int(result.Code), true
}

// ParseMetricCodeFromBodyFunc1 是默认从 http body 解析业务 code 的方法(兜底方式，body太大时消耗大)
var ParseMetricCodeFromBodyFunc1 ParseMetricCodeFromBody = func(in []byte) (code int, ok bool) {
	if in == nil {
		return
	}
	result := make(map[string]int)
	if err := json.Unmarshal(in, &result); err != nil {
		return
	}
	if v, ok := result["code"]; ok {
		return v, true
	}
	if v, ok := result["Code"]; ok {
		return v, true
	}
	return code, false
}

// DemoParseMetricCodeFromRsp 是从 err 解析 metric code 的示例函数
var DemoParseMetricCodeFromErr = func(err error) (code int, ok bool) {
	if e, ok := err.(CodeError); ok && e != nil {
		return e.GetCode(), true
	}
	if e, ok := err.(*oError.CommonError); ok && e != nil {
		return int(e.Code), true
	}
	return
}

// DemoParseMetricCodeFromRsp 是从 rsp 解析 metric code 的示例函数
var DemoParseMetricCodeFromRsp = func(rsp interface{}) (code int, ok bool) {
	if rsp == nil {
		return
	}
	if vv, ok := rsp.(*RespStruct); ok && vv != nil {
		return int(vv.Code), true
	}
	return
}

type MetricEchoReq struct {
	// kind 是类型, 无值：参数解析错误，1:正常，2:rsp中返回自定义错误，3:返回oError错误，4:grpc错误，5:自定义CodeErr错误，6:底层值为空的err
	Kind  int    `json:"kind" form:"kind" binding:"required"`
	Msg   string `json:"msg" form:"msg"`
	Sleep int    `json:"sleep" form:"sleep"` // 延迟时间，单位 ms(-1:不延迟，0:随机延迟)
}

// MetricEcho 用于监控调试
func MetricEcho(ctx context.Context, req *MetricEchoReq, rsp *RespStruct) error {
	var sleep int
	if req.Sleep > 0 {
		sleep = req.Sleep
	} else if req.Sleep == 0 {
		sleep = rand.Intn(1000)
	}
	if sleep > 0 {
		time.Sleep(time.Millisecond * time.Duration(sleep))
	}
	rsp.Data = fmt.Sprintf("Sleep %vms", sleep)
	switch req.Kind {
	case 0:
		// 参数为空，由上游直接返回
		return nil
	case 1:
		// 正常
		rsp.Message = "OK:" + req.Msg
		return nil
	case 2:
		// rsp中返回自定义错误
		rsp.Code = 2
		rsp.Message = "定制错误" + req.Msg
		return nil
	case 3:
		// oError错误
		return oError.New(oError.ErrorCodeUnknown, "oError错误:"+req.Msg)
	case 4:
		// grpc错误
		err := gstatus.Errorf(gcodes.Internal, "grpc err")
		return err

	case 5:
		// 自定义CodeErr错误
		return &selfErr{ErrCode: -1, ErrMsg: "self err"}
	case 6:
		// 返回有类型，但是底层为空的 error
		rsp.Code = 6
		rsp.Message = "底层值为空的错误"
		var err *oError.CommonError
		return err
	case 7:
		rsp.Code = 6
		rsp.Message = "底层值为空的错误"
		var err *selfErr
		return err
	default:
		return errors.New("未定义错误:" + req.Msg)
	}
}

type selfErr struct {
	ErrCode int64
	ErrMsg  string
}

func (e *selfErr) Error() string {
	return e.ErrMsg
}

func (e *selfErr) GetCode() int64 {
	return e.ErrCode
}

// parseDefaultServiceCode 从返回 error 中解析服务 code（默认解析规则）
// 需要调用放保证 err 底层值非空
func parseDefaultServiceCode(ctx context.Context, err error) (code int64, msg string, ok bool) {
	if err == nil {
		return
	}
	defer func() {
		// 以防万一：处理接口底层值为空时 panic 的坑
		if e := recover(); e != nil {
			code = int64(oError.ErrorCodeSystem)
			ok = false
			logger.Errorf(ctx, "parseDefaultServiceCode recover:%v", e)
		}
	}()
	switch e := err.(type) {
	case *oError.CommonError:
		if e != nil {
			return int64(e.Code), e.Message, true
		}
	case CodeError:
		if e != nil {
			return int64(e.GetCode()), e.Error(), true
		}
	case CodeError64:
		if e != nil {
			return e.GetCode(), e.Error(), true
		}
	default:
		if e != nil {
			// grpc err
			if s, ok := gstatus.FromError(err); ok && s != nil && s.Code() != gcodes.Unknown {
				return int64(s.Code()), s.Message(), true
			}
			return int64(oError.ErrorCodeSystem), e.Error(), true
		}
	}
	return
}

// 支持业务 code 的 error 类型
type CodeError interface {
	GetCode() int
	error
}
type CodeError64 interface {
	GetCode() int64
	error
}
