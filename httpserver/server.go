package httpserver

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yunxiaoyang01/open_sdk/httpserver/middles"
	"github.com/yunxiaoyang01/open_sdk/logger"

	"github.com/gin-gonic/gin"
)

type Server interface {
	Name() (name string)
	Run(ctx context.Context) (err error)
	AddMiddles(ms ...middles.Middle)
	GetKernel() (kernel *gin.Engine)
	RegisterOnShutdown(f func())
}

type server struct {
	kernel     *gin.Engine
	options    Options
	onShutdown []func()
}

func NewServer(opts ...Option) Server {
	return NewServerWithOptions(newOptions(opts...))
}

func NewServerWithOptions(options Options) Server {
	kernel := gin.New()
	// default Use middles
	// 默认中间件， FlowControlTag，错误恢复
	kernel.Use(
		middles.Recovery(),
	)

	kernel.Use(options.Middles...) // user set middles
	s := &server{kernel: kernel, options: options}
	return s
}

func (s *server) Name() string {
	return s.options.Name
}

func (s *server) Run(ctx context.Context) error {
	srv := &http.Server{Handler: s.kernel, Addr: s.options.Address}

	for _, f := range s.onShutdown { // 注册用户自定义的关闭函数
		srv.RegisterOnShutdown(f)
	}

	// check if inside docker, for local debug.
	if outsideContainer() { // 如果检测出运行环境不是在容器内部，则替换request.ctx为context.Background()，方便在本地调试
		middles.MutateRequest = changeRequest
	}

	// server listenAndServe
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Errorf(ctx, "httpServer ListenAndServe err:%v", err)
			}
			logger.Infof(ctx, "http server is closed")
		}
	}()

	// handle signal
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	logger.Infof(ctx, "got signal %v, exit", <-ch)
	srv.Shutdown(context.Background())
	return nil
}

func (s *server) RegisterOnShutdown(f func()) {
	s.onShutdown = append(s.onShutdown, f)
}

func (s *server) AddMiddles(ms ...middles.Middle) {
	s.kernel.Use(ms...)
}

func (s *server) GetKernel() *gin.Engine {
	return s.kernel
}
