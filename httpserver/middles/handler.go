package middles

import (
	"context"
	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	xe "github.com/yunxiaoyang01/open_sdk/error"
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
	ErrMustError         = errors.New("method ret must be error or commonerror")
	ErrMustOneOut        = errors.New("method must has one out")

	initerType           = reflect.TypeOf((*Initer)(nil)).Elem()
	replyErrorType       = reflect.TypeOf((*error)(nil)).Elem()
	replyCommonErrorType = reflect.TypeOf((*xe.CommonError)(nil))

	RequestKey    = requestKey{}
	ResponseKey   = responseKey{}
	GinContextKey = ginContextKey{}
)

func CreateHandlerFunc(method interface{}, opts ...Option) gin.HandlerFunc {
	return CreateHandlerFuncWithLogger(method, logger.StandardLogger(), opts...)
}

func checkMethod(method interface{}) (mV reflect.Value, reqT, replyT reflect.Type, err error) {
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
	if retT != replyErrorType && retT != replyCommonErrorType {
		err = ErrMustError
		return
	}
	return mV, reqT, replyT, err
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
	Code int64       `json:"code"`
	Data interface{} `json:"data"`
}

type createHandlerOptions struct {
	hasDataKey   bool
	notWriteResp bool
}

type Option func(opt *createHandlerOptions)

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

func getOption(opts ...Option) createHandlerOptions {
	option := createHandlerOptions{
		hasDataKey: true,
	}

	for _, opt := range opts {
		opt(&option)
	}

	return option
}

// MutateRequest mutate request
var MutateRequest func(r *http.Request) = func(r *http.Request) {}

func CreateHandlerFuncWithLogger(method interface{}, l logger.Logger, opts ...Option) gin.HandlerFunc {
	option := getOption(opts...)

	mV, reqT, replyT, err := checkMethod(method)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		MutateRequest(c.Request)
		ctx := c.Request.Context()
		req := reflect.New(reqT)

		if err := c.ShouldBind(req.Interface()); err != nil {
			l.WithFields(logger.Fields{
				"req": c.Request.URL,
				"err": err,
			}).Warn(ctx, "bind param failed")

			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    xe.ErrorCodeWrongParam,
				"message": err.Error(),
			})
			return
		}

		ctx = context.WithValue(ctx, RequestKey, c.Request)
		ctx = context.WithValue(ctx, ResponseKey, c.Writer)
		ctx = context.WithValue(ctx, GinContextKey, c)

		callFieldInit(ctx, req)

		reply := reflect.New(replyT)
		l.WithFields(logger.Fields{
			"req":  req,
			"func": mV.Type().String(),
		}).Debug(ctx, "invoke handler")

		results := mV.Call([]reflect.Value{reflect.ValueOf(ctx), req, reply})
		errValue := results[0]
		if errValue.Interface() != nil {
			switch v := errValue.Interface().(type) {
			case *xe.CommonError:
				if v != nil {
					l.WithFields(logger.Fields{
						"req": c.Request.URL,
						"err": v,
					}).Warn(ctx, "handler err")
					c.JSON(http.StatusOK, map[string]interface{}{
						"code":    v.Code,
						"message": v.Message,
					})
					return
				}
			case error:
				if v != nil {
					l.WithFields(logger.Fields{
						"req": c.Request.URL,
						"err": v,
					}).Warn(ctx, "handler err")
					c.JSON(http.StatusOK, map[string]interface{}{
						"code":    xe.ErrorCodeSystem,
						"message": v.Error(),
					})
					return
				}
			}
		}

		ret := getRetData(reply, option)

		if option.notWriteResp {
			return
		}

		statusCode := c.Writer.Status()
		if statusCode != 0 {
			c.PureJSON(statusCode, ret)
			return
		}

		c.PureJSON(http.StatusOK, ret)
	}
}

func getRetData(value reflect.Value, option createHandlerOptions) interface{} {
	if option.hasDataKey {
		ret := &RespStruct{Code: 0, Data: value.Interface()}
		return ret
	}
	return value.Interface()
}
