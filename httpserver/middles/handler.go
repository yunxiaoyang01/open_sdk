package middles

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	xe "github.com/yunxiaoyang01/open_sdk/error"
	"github.com/yunxiaoyang01/open_sdk/httpserver/io"
	"github.com/yunxiaoyang01/open_sdk/httpserver/status"
	"github.com/yunxiaoyang01/open_sdk/logger"
)

type requestKey struct{}
type responseKey struct{}
type ginContextKey struct{}

type Initer interface {
	Init(ctx context.Context)
}

var (
	EmptyStr = ""

	ErrMustPtr           = errors.New("param must be ptr")
	ErrMustPointToStruct = errors.New("param must point to struct")
	ErrMustHasThreeParam = errors.New("method must has three input")
	ErrMustFunc          = errors.New("method must be func")
	ErrMustValid         = errors.New("method must be valid")
	ErrMustError         = errors.New("method ret must be error or xderror")
	ErrMustOneOut        = errors.New("method must has one out")
	ErrWrongMethodType   = errors.New("method 格式不对")

	initerType       = reflect.TypeOf((*Initer)(nil)).Elem()
	replyErrorType   = reflect.TypeOf((*error)(nil)).Elem()
	replyXDErrorType = reflect.TypeOf((*xe.CommonError)(nil))

	RequestKey    = requestKey{}
	ResponseKey   = responseKey{}
	GinContextKey = ginContextKey{}
)

func NewHandlerFuncCreator(opts ...Option) HandlerFuncCreator {
	return &handlerFuncCreator{
		opts: opts,
		log:  logger.StandardLogger(),
	}
}

// HandlerFuncCreator HTTP处理函数的创建者，使用统一的 logger 和 option
type HandlerFuncCreator interface {
	SetLogger(log logger.Logger)
	WithOption(opt ...Option)
	// CreateHandlerFunc 从给定的 method 函数创建 gin.HandlerFunc
	// method 格式为 Method(ctx context.Context, req *ReqObj, rsp *RspObj) error
	CreateHandlerFunc(method interface{}, opt ...Option) gin.HandlerFunc
	// NewHandlerFuncFrom 从 GRPC handler 中创建 gin handler
	// method 格式为 Method(ctx context.Context, req *ReqObj) (rsp *RspObj, err error)
	NewHandlerFuncFrom(method interface{}, opt ...Option) gin.HandlerFunc
}

type handlerFuncCreator struct {
	log  logger.Logger
	opts []Option
}

func (c *handlerFuncCreator) SetLogger(log logger.Logger) {
	c.log = log
}
func (c *handlerFuncCreator) WithOption(opt ...Option) {
	for _, v := range opt {
		if v != nil {
			c.opts = append(c.opts, v)
		}
	}
}

// CreateHandlerFunc 从给定的 method 函数创建 gin.HandlerFunc
// method 格式为 Method(ctx context.Context, req *ReqObj, rsp *RspObj) error
func (c *handlerFuncCreator) CreateHandlerFunc(method interface{}, opt ...Option) gin.HandlerFunc {
	opts := append(c.opts, opt...)
	return CreateHandlerFuncWithLogger(method, c.log, opts...)
}

// NewHandlerFuncFrom 从 GRPC handler 中创建 gin handler
// method 格式为 Method(ctx context.Context, req *ReqObj) (rsp *RspObj, err error)
func (c *handlerFuncCreator) NewHandlerFuncFrom(method interface{}, opt ...Option) gin.HandlerFunc {
	opts := append(c.opts, opt...)
	return NewHandlerFuncWithLoggerFrom(method, c.log, opts...)
}

// CreateHandlerFunc 从给定的 method 函数创建 gin.HandlerFunc
// method 格式为 Method(ctx context.Context, req *ReqObj, rsp *RspObj) error
func CreateHandlerFunc(method interface{}, opts ...Option) gin.HandlerFunc {
	return CreateHandlerFuncWithLogger(method, logger.StandardLogger(), opts...)
}

// NewHandlerFuncFrom 从 GRPC handler 中创建 gin handler
// method 格式为 Method(ctx context.Context, req *ReqObj) (rsp *RspObj, err error)
func NewHandlerFuncFrom(method interface{}, opts ...Option) gin.HandlerFunc {
	return NewHandlerFuncWithLoggerFrom(method, logger.StandardLogger(), opts...)
}

// method 格式为 Method(ctx context.Context, req *ReqObj, rsp *RspObj) error
func check31Method(method interface{}) (mV reflect.Value, reqT, replyT reflect.Type, err error) {
	mV = reflect.ValueOf(method)
	if !mV.IsValid() {
		err = ErrMustValid
		return
	}

	mT := mV.Type()
	if mT.Kind() != reflect.Func {
		err = ErrMustFunc
		return
	}

	if mT.NumIn() != 3 {
		err = ErrMustHasThreeParam
		return
	}

	reqT = mT.In(1)
	if reqT.Kind() != reflect.Ptr {
		err = ErrMustPtr
		return
	}

	if reqT.Elem().Kind() != reflect.Struct {
		err = ErrMustPointToStruct
		return
	}

	reqT = reqT.Elem()
	replyT = mT.In(2)

	if replyT.Kind() != reflect.Ptr {
		err = ErrMustPtr
		return
	}
	/*if replyT.Elem().Kind() != reflect.Struct {
		err = ErrMustPointToStruct
		return
	}*/
	replyT = replyT.Elem()
	if mT.NumOut() != 1 {
		err = ErrMustOneOut
		return
	}
	retT := mT.Out(0)
	if retT != replyErrorType && retT != replyXDErrorType {
		err = ErrMustError
		return
	}
	return mV, reqT, replyT, err
}

// method 格式为 (h *Handler) Method(ctx context.Context, req *ReqObj) (rsp *RspObj, err error)
func check22Method(method interface{}) (mV reflect.Value, reqT reflect.Type, err error) {
	mV = reflect.ValueOf(method)
	if !mV.IsValid() {
		err = ErrMustValid
		return
	}

	mT := mV.Type()
	if mT.Kind() != reflect.Func {
		err = ErrMustFunc
		return
	}

	if mT.NumIn() != 2 {
		err = ErrWrongMethodType
		return
	}

	reqT = mT.In(1)
	if reqT.Kind() != reflect.Ptr {
		err = ErrMustPtr
		return
	}

	if reqT.Elem().Kind() != reflect.Struct {
		err = ErrMustPointToStruct
		return
	}
	reqT = reqT.Elem()

	// rspT = mT.In(2)
	// if rspT.Kind() != reflect.Ptr {
	// 	err = ErrMustPtr
	// 	return
	// }
	/*if rspT.Elem().Kind() != reflect.Struct {
		err = ErrMustPointToStruct
		return
	}*/
	if mT.NumOut() != 2 {
		err = ErrMustOneOut
		return
	}
	rspT := mT.Out(0)
	if rspT.Kind() != reflect.Ptr || rspT.Elem().Kind() != reflect.Struct {
		err = ErrMustPointToStruct
		return
	}

	errT := mT.Out(1)
	if errT != replyErrorType && errT != replyXDErrorType {
		err = ErrMustError
		return
	}
	return mV, reqT, err
}

func isImplementIniter(v reflect.Value) bool {
	return v.Type().Implements(initerType)
}

func callFieldInit(ctx context.Context, v reflect.Value) {
	elem := v.Elem()
	vT := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		ev := elem.Field(i)
		if isImplementIniter(ev) {
			if ev.CanSet() {
				ev.Set(reflect.New(vT.Field(i).Type.Elem()))
				initer := ev.Interface().(Initer)
				initer.Init(ctx)
			}
		}
	}
}

type RespStruct struct {
	Code    int64       `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type createHandlerOptions struct {
	hasDataKey   bool
	notWriteResp bool

	withoutCodeMetric bool                   // 关闭 metric
	parseCodeFromErr  ParseMetricCodeFromErr // 从返回err中解析业务code
	parseCodeFromRsp  ParseMetricCodeFromRsp // 从返回rsp中解析业务code，不推荐
	replaceMetricPath string                 // 替换监控 path，用于 path 中带变量参数时
}

type Option func(opt *createHandlerOptions)

func getOption(opts ...Option) createHandlerOptions {
	option := createHandlerOptions{
		hasDataKey: true,
	}

	for _, opt := range opts {
		opt(&option)
	}

	return option
}
func WithDataKey(has bool) Option {
	return func(opt *createHandlerOptions) {
		opt.hasDataKey = has
	}
}

func WithNotWriteResp(notWriteResp bool) Option {
	return func(opt *createHandlerOptions) {
		opt.notWriteResp = notWriteResp
	}
}

// WithoutCodeMetric 设置为：关闭metric
func WithoutCodeMetric() Option {
	return func(opt *createHandlerOptions) {
		opt.withoutCodeMetric = true
	}
}

// WithReplaceMetricPath 用于替换监控的 path 项
func WithReplaceMetricPath(path string) Option {
	return func(opt *createHandlerOptions) {
		opt.replaceMetricPath = path
	}
}

// WithParseMetricCodeFromErr 设置自定义解析 http 接口返回 code 的函数
func WithParseMetricCodeFromErr(f ParseMetricCodeFromErr) Option {
	return func(opt *createHandlerOptions) {
		opt.parseCodeFromErr = f
	}
}

// WithParseMetricCodeFromRsp 设置自定义解析 http 接口返回 code 的函数
func WithParseMetricCodeFromRsp(f ParseMetricCodeFromRsp) Option {
	return func(opt *createHandlerOptions) {
		opt.parseCodeFromRsp = f
	}
}

// MutateRequest mutate request
var MutateRequest func(r *http.Request) = func(r *http.Request) {}

// CreateHandlerFuncWithLogger 中 method 格式为 Method(ctx context.Context, req *ReqObj, rsp *RspObj) error
func CreateHandlerFuncWithLogger(method interface{}, l logger.Logger, opts ...Option) gin.HandlerFunc {
	option := getOption(opts...)

	mV, reqT, replyT, err := check31Method(method)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		MutateRequest(c.Request)
		ctx := c.Request.Context()
		req := reflect.New(reqT)

		if err := c.ShouldBind(req.Interface()); err != nil {
			l.WithFields(logger.Fields{"req": c.Request.URL.Path, "err": err}).Warn(ctx, "bind param failed")
			c.JSON(http.StatusBadRequest, RespStruct{Code: int64(xe.ErrorCodeWrongParam), Message: err.Error()})
			return
		}

		ctx = context.WithValue(ctx, RequestKey, c.Request)
		ctx = context.WithValue(ctx, ResponseKey, c.Writer)
		ctx = context.WithValue(ctx, GinContextKey, c)

		callFieldInit(ctx, req)

		reply := reflect.New(replyT)
		l.WithFields(logger.Fields{"req": req, "func": mV.Type().String()}).Debug(ctx, "invoke handler")

		results := mV.Call([]reflect.Value{reflect.ValueOf(ctx), req, reply})
		errValue := results[0]
		if errValue.IsValid() && errValue.CanInterface() && !errValue.IsZero() && errValue.Elem().IsValid() && !errValue.Elem().IsZero() {
			l.WithFields(logger.Fields{"url": c.Request.URL.Path}).Warn(ctx, "handler err: ", errValue)
			err := errValue.Interface().(error) // checkMethod 中已经保证是 error 类型
			if code, msg, ok := parseDefaultServiceCode(ctx, err); ok {
				status := http.StatusOK
				c.JSON(status, RespStruct{Code: code, Message: msg})
				return
			}
		}
		ret := reply.Interface()
		if option.notWriteResp {
			return
		}
		if option.hasDataKey {
			ret = RespStruct{Code: 0, Data: ret}
		}
		if statusCode := c.Writer.Status(); statusCode != 0 {
			c.PureJSON(statusCode, ret)
			return
		}
		c.PureJSON(http.StatusOK, ret)
	}
}

/*
func getRetData(value reflect.Value, option createHandlerOptions) interface{} {
	if option.hasDataKey {
		ret := &RespStruct{Code: 0, Data: value.Interface()}
		return ret
	}
	return value.Interface()
}
*/

// CreateObjectHandlerFuncWithLogger method 格式为 (h *Handler) Method(ctx context.Context, req *ReqObj) (rsp *RspObj, err error)
func CreateObjectHandlerFuncWithLogger(handler interface{}, method string, l logger.Logger, opts ...Option) gin.HandlerFunc {
	//option := getOption(opts...)

	hV := reflect.ValueOf(handler)
	mV := hV.MethodByName(method)
	if !mV.IsValid() {
		panic(fmt.Errorf("method(%s) not found", method))
	}
	mT := mV.Type()
	if mT.NumIn() != 2 {
		panic(fmt.Errorf("method(%s) must has 2 ins", method))
	}
	if mT.NumOut() != 2 {
		panic(fmt.Errorf("method(%s) must has 2 out", method))
	}
	if errT := mT.Out(1); errT != replyErrorType && errT != replyXDErrorType {
		panic(fmt.Errorf("method(%s) second ret must be err", method))
	}
	reqT := mT.In(1).Elem()

	return func(c *gin.Context) {
		MutateRequest(c.Request)
		var (
			ctx = c.Request.Context()
			req = reflect.New(reqT)
			err error
		)

		// bind request data
		if binder, ok := handler.(io.Binder); ok {
			err = binder.Bind(c, method, req.Interface())
		} else {
			err = c.ShouldBind(req.Interface())
		}
		if err != nil {
			l.Errorf(ctx, "method(%s) failed to bind: %v", method, err)
			c.String(http.StatusBadRequest, err.Error()) // todo: 此处 c.String 是个历史遗留的坑，应该替换为 JSON
			// c.JSON(http.StatusBadRequest, RespStruct{Code: 1, Message: err.Error()})
			return
		}

		results := mV.Call([]reflect.Value{reflect.ValueOf(ctx), req})
		errValue := results[1]
		if errValue.IsValid() && errValue.CanInterface() && !errValue.IsZero() && errValue.Elem().IsValid() && !errValue.Elem().IsZero() {
			l.WithFields(logger.Fields{"url": c.Request.URL.Path}).Warn(ctx, "handler err: ", errValue)
			err := errValue.Interface().(error)
			if _, msg, ok := parseDefaultServiceCode(ctx, err); ok {
				status := status.GetCode(err) // 兼容历史版本
				// c.JSON(status, RespStruct{Code: code, Message: msg})
				c.String(status, msg) // todo: 此处 c.String 是个历史遗留的坑，应该替换为 JSON
				return
			}
		}
		response := results[0].Interface()

		// rend result
		if render, ok := handler.(io.Render); ok {
			render.Rend(c, method, response)
		} else {
			c.PureJSON(http.StatusOK, response)
		}
	}
}

// NewHandlerFuncWithLoggerFrom 从 method 创建 gin.HandlerFunc
// method 格式为 (h *Handler) Method(ctx context.Context, req *ReqObj) (rsp *RspObj, err error)
func NewHandlerFuncWithLoggerFrom(method interface{}, l logger.Logger, opts ...Option) gin.HandlerFunc {
	option := getOption(opts...)

	mV, reqT, err := check22Method(method)
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		MutateRequest(c.Request)
		ctx := c.Request.Context()
		req := reflect.New(reqT)

		if err := c.ShouldBind(req.Interface()); err != nil {
			l.WithFields(logger.Fields{"req": c.Request.URL.Path, "err": err}).Warn(ctx, "bind param failed")
			c.JSON(http.StatusBadRequest, RespStruct{Code: 1, Message: err.Error()})
			return
		}
		ctx = context.WithValue(ctx, RequestKey, c.Request)
		ctx = context.WithValue(ctx, ResponseKey, c.Writer)
		ctx = context.WithValue(ctx, GinContextKey, c)
		callFieldInit(ctx, req)

		// reply := reflect.New(replyT)
		l.WithFields(logger.Fields{"req": req, "func": mV.Type().String()}).Debug(ctx, "invoke handler")

		results := mV.Call([]reflect.Value{reflect.ValueOf(ctx), req})
		errValue := results[1]
		if errValue.IsValid() && !errValue.IsZero() && errValue.CanInterface() && errValue.Elem().IsValid() && !errValue.Elem().IsZero() {
			l.WithFields(logger.Fields{"url": c.Request.URL.Path}).Warn(ctx, "handler err: ", errValue)
			err := errValue.Interface().(error)
			if code, msg, ok := parseDefaultServiceCode(ctx, err); ok {
				status := status.GetCode(err) // 兼容历史版本
				c.JSON(status, RespStruct{Code: code, Message: msg})
				return
			}
		}
		ret := results[0].Interface()
		if option.notWriteResp {
			return
		}
		if option.hasDataKey {
			ret = RespStruct{Code: 0, Data: ret}
		}
		if statusCode := c.Writer.Status(); statusCode != 0 {
			c.PureJSON(statusCode, ret)
			return
		}
		c.PureJSON(http.StatusOK, ret)
	}
}
