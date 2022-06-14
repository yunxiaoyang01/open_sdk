package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yunxiaoyang01/open_sdk/httpclient/hooks"
)

func TestHttpClientRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()
	hc, err := NewClient(
		WithTimeout(3*time.Second),
		WithRetryCount(3),
		WithPreRequestHooks(hooks.LoggingRequest()), // 请求之前打印日志，方便查看header信息
		WithAfterResponseHooks(hooks.LoggingResponse()),
	)
	if err != nil {
		t.Error(err.Error())
	}

	ctx := context.Background()
	req := hc.NewRequest(ctx)
	_, err = req.Get(ts.URL) // 日志里，header信息无压测标记
	if err != nil {
		t.Error(err.Error())
	}
	_, err = hc.NewRequest(ctx).Get(ts.URL)
	if err != nil {
		t.Error(err.Error())
	}
}
