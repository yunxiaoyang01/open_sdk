package httpclient

import (
	"time"

	"github.com/yunxiaoyang01/open_sdk/httpclient/hooks"
)

type Options struct {
	Address            string
	Timeout            time.Duration
	RetryCount         int           // 重试次数
	RetryWaitTime      time.Duration // 重试间隔等待时间
	RetryMaxWaitTime   time.Duration // 重试间隔最大等待时间
	BeforeRequestHooks []hooks.BeforeRequestHook
	PreRequestHooks    []hooks.PreRequestHook
	AfterResponseHooks []hooks.AfterResponseHook
}

func newOptions(opts ...Option) Options {
	options := Options{
		Address:            "",
		Timeout:            3 * time.Second,
		RetryCount:         0,
		RetryWaitTime:      time.Duration(100) * time.Millisecond,
		RetryMaxWaitTime:   time.Duration(2000) * time.Millisecond,
		BeforeRequestHooks: []hooks.BeforeRequestHook{},
		PreRequestHooks:    []hooks.PreRequestHook{},
		AfterResponseHooks: []hooks.AfterResponseHook{},
	}
	for _, opt := range opts {
		opt(&options)
	}
	options.PreRequestHooks = append(options.PreRequestHooks, hooks.DNSTrace())
	options.AfterResponseHooks = append(options.AfterResponseHooks, hooks.LatencyMetrics())
	return options
}

type Option func(*Options)

func WithAddress(address string) Option {
	return func(options *Options) {
		options.Address = address
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.Timeout = timeout
	}
}

func WithRetryCount(retryCount int) Option {
	return func(options *Options) {
		options.RetryCount = retryCount
	}
}

func WithRetryWaitTime(retryWaitTime time.Duration) Option {
	return func(options *Options) {
		options.RetryWaitTime = retryWaitTime
	}
}

func WithRetryMaxWaitTime(retryMaxWaitTime time.Duration) Option {
	return func(options *Options) {
		options.RetryMaxWaitTime = retryMaxWaitTime
	}
}

func WithBeforeRequestHooks(hs ...hooks.BeforeRequestHook) Option {
	return func(options *Options) {
		options.BeforeRequestHooks = hs
	}
}

func WithPreRequestHooks(hs ...hooks.PreRequestHook) Option {
	return func(options *Options) {
		options.PreRequestHooks = hs
	}
}

func WithAfterResponseHooks(hs ...hooks.AfterResponseHook) Option {
	return func(options *Options) {
		options.AfterResponseHooks = hs
	}
}
