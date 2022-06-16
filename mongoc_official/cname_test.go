package mongoc_official

import (
	"context"
	"testing"
)

func TestCName(t *testing.T) {

	var tests = []struct {
		ctx  context.Context
		name string
		want string
	}{
		{
			context.Background(),
			"test",
			"test",
		},
	}

	for _, test := range tests {
		b := &Base{
			coll: test.name,
		}
		get := b.cname(test.ctx)
		if get != test.want {
			t.Errorf("want: %v  get: %v", test.want, get)
		}
	}
}
