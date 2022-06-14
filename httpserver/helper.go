package httpserver

import (
	"context"
	"net/http"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yunxiaoyang01/open_sdk/httpserver/io"
	"github.com/yunxiaoyang01/open_sdk/httpserver/middles"
	"github.com/yunxiaoyang01/open_sdk/httpserver/status"
	"github.com/yunxiaoyang01/open_sdk/logger"
)

// see comment on http.Request.Context()
func changeRequest(r *http.Request) {
	*r = *r.WithContext(context.Background())
}

// there will be a `/.dockerenv` file if running inside docker.
func outsideContainer() bool {
	_, err := os.Stat("/.dockerenv")
	if err == nil {
		return false
	}
	return os.Getenv("CommonRUNTIME") == "" // default is true, set CommonRUNTIME to disable this.
}

func CreateHandlerFunc(handler interface{}, method string) gin.HandlerFunc {
	return CreateHandlerFuncWithLogger(handler, method, logger.StandardLogger())
}

func CreateHandlerFuncWithLogger(handler interface{}, method string, l logger.Logger) gin.HandlerFunc {
	hV := reflect.ValueOf(handler)
	mV := hV.MethodByName(method)
	if !mV.IsValid() {
		panic(errors.Errorf("method(%s) not found", method))
	}
	mT := mV.Type()
	if mT.NumIn() != 2 {
		panic(errors.Errorf("method(%s) must has 2 ins", method))
	}
	reqT := mT.In(1).Elem()
	if mT.NumOut() != 2 {
		panic(errors.Errorf("method(%s) must has 2 out", method))
	}
	return func(c *gin.Context) {
		middles.MutateRequest(c.Request)
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
			c.String(http.StatusBadRequest, err.Error())
			l.Errorf(ctx, "method(%s) failed to bind: %v", method, err)
			return
		}
		results := mV.Call([]reflect.Value{reflect.ValueOf(ctx), req})
		response, e := results[0].Interface(), results[1].Interface()
		if e != nil {
			err := e.(error)
			c.String(status.GetCode(err), err.Error())
			l.Errorf(ctx, "method(%s) failed to call: %v", method, err)
			return
		}

		// rend result
		if render, ok := handler.(io.Render); ok {
			render.Rend(c, method, response)
		} else {
			c.PureJSON(http.StatusOK, response)
		}
	}
}

func Route(routes gin.IRoutes, method string, path string, function interface{}, opts ...middles.Option) {
	routes.Handle(method, path, middles.CreateHandlerFunc(function, opts...))
}
