package config

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/yunxiaoyang01/open_sdk/logger"
)

func NewLoggerHook() Hook {
	return func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		if from.Kind() != reflect.Map {
			return data, nil
		}
		if to == reflect.TypeOf(logger.Config{}) {
			var config logger.Config
			err := mapstructure.Decode(data, &config.Options)
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode data to logger.Config")
			}
			config.Logger, err = logger.NewLoggerWithOptions(config.Options)
			if err != nil {
				return nil, errors.Wrap(err, "failed to new logger")
			}
			return config, nil
		}
		if to == reflect.TypeOf(logger.Standard{}) {
			var standard logger.Standard
			err := mapstructure.Decode(data, &standard)
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode data to logger.Standard")
			}
			err = logger.ResetStandardWithOptions(logger.Options(standard))
			if err != nil {
				return nil, errors.Wrap(err, "failed to reset standard")
			}
			return standard, nil
		}

		return data, nil
	}
}
