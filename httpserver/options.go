package httpserver

import (
	"github.com/yunxiaoyang01/open_sdk/httpserver/middles"
)

type Options struct {
	Name    string
	Address string
	Middles []middles.Middle
}

func newOptions(opts ...Option) Options {
	options := Options{
		Name:    "httpserver",
		Address: ":8080",
		Middles: nil,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

type Option func(*Options)

func WithName(name string) Option {
	return func(options *Options) {
		options.Name = name
	}
}

func WithAddress(address string) Option {
	return func(options *Options) {
		options.Address = address
	}
}

func WithMiddles(ms ...middles.Middle) Option {
	return func(options *Options) {
		options.Middles = ms
	}
}
