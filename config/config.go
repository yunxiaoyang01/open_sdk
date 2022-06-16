package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Hook = mapstructure.DecodeHookFunc

var defaultHooks = []Hook{
	NewLoggerHook(),
	// NewMongocHook(),
	// NewRediscHook(),
	// NewElasticsearchcHook(),
	// NewTikvcHook(),
	// NewProfHook(),
}

func RegisterDefaultHook(hook Hook) {
	defaultHooks = append(defaultHooks, hook)
}

func Load(path string, conf interface{}) error {
	return LoadWithHooks(path, conf, defaultHooks...)
}

func LoadWithHooks(path string, conf interface{}, hooks ...Hook) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "failed to read config")
	}
	err := viper.Unmarshal(conf, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(hooks...)))
	if err != nil {
		return errors.Wrap(err, "failed to scan config")
	}
	return nil
}
