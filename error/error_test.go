package error

import (
	"errors"
	"regexp"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		code ErrorCodeType
		msg  []string
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name: "normal test",
			args: args{
				code: ErrorCodeOK,
				msg:  []string{"something"},
			},
			expected: `{"code":0,"message":"something","rt_stacks":".*:30.*"}`, // 30 is the line of where New been called.
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.args.code, tt.args.msg...)
			if ok, _ := regexp.MatchString(tt.expected, err.Error()); !ok {
				t.Errorf("New() error = %v, expected %v", err, tt.expected)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	ori := errors.New("")
	tests := []struct {
		ori  error
		wrap error
		want bool
	}{
		{
			ori:  ori,
			wrap: Wrap(ori),
			want: true,
		},
		{
			ori:  ori,
			wrap: WrapCode(ori, ErrorCodeUnknown),
			want: true,
		},
		{
			ori:  ori,
			wrap: ErrToXDError(ori, ErrorCodeUnknown),
			want: true,
		},
		{
			ori:  errors.New(""),
			wrap: Wrap(ori),
			want: false,
		},
	}

	for k, test := range tests {
		if is := errors.Is(test.wrap, test.ori); is != test.want {
			t.Errorf("test err: idx: %v", k)
		}
	}
}
