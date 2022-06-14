package error

import (
	"encoding/json"
	"fmt"
	"runtime"
)

const (
	maxStackDepth = 32
	skipStack     = 2
)

type XDError struct {
	Code     ErrorCodeType `json:"code"`
	Message  string        `json:"message"`
	RTStacks string        `json:"rt_stacks"`
	err      error
}

func (e *XDError) Error() string {
	data, _ := json.Marshal(e)

	return string(data)
}

func (e *XDError) Unwrap() error {
	return e.err
}

func getStacks(skip int) string {
	pc := make([]uintptr, maxStackDepth)
	stacks := ""

	n := runtime.Callers(skip, pc)
	for i := 0; i < n-2; i++ { // skip some basic frames.
		f := runtime.FuncForPC(pc[i])
		file, line := f.FileLine(pc[i])
		stacks += fmt.Sprintf("file: %s:%d func: %s\n", file, line, f.Name())
	}

	return stacks
}

func newError(code ErrorCodeType, err error, msg ...string) error {
	str := ""
	if len(msg) > 0 {
		str = msg[0]
	}

	stacks := getStacks(skipStack)

	return &XDError{
		Code:     code,
		Message:  str,
		RTStacks: stacks,
		err:      err,
	}
}

func New(code ErrorCodeType, msg ...string) error {
	return newError(code, nil, msg...)
}

func ErrToXDError(err error, code ErrorCodeType) error {
	if err == nil {
		return nil
	}

	_, ok := err.(*XDError)
	if ok {
		return err
	}

	return newError(code, err, err.Error())
}

func WrapCode(err error, code ErrorCodeType) error {
	return ErrToXDError(err, code)
}

func Wrap(err error, msg ...string) error {
	if err == nil {
		return nil
	}

	_, ok := err.(*XDError)
	if ok {
		return err
	}

	var info string
	if len(msg) > 0 {
		info = fmt.Sprintf("%s: %s", msg[0], err.Error())
	} else {
		info = err.Error()
	}

	return newError(ErrorCodeUnknown, err, info)
}

func Code(err error) ErrorCodeType {
	xderror, ok := err.(*XDError)
	if !ok {
		return ErrorCodeUnknown
	}

	return xderror.Code
}
