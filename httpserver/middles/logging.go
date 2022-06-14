package middles

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yunxiaoyang01/open_sdk/logger"
)

const maxBodyLen = 2 << 15

func LoggingRequest() Middle {
	return LoggingRequestWithLogger(logger.StandardLogger())
}

func LoggingRequestWithLogger(l logger.Logger) Middle {
	return func(c *gin.Context) {
		l.WithFields(logger.Fields{
			"method":    c.Request.Method,
			"uri":       c.Request.URL.RequestURI(),
			"client_ip": c.Request.RemoteAddr,
			//"headers": composeHeaders(c.Request.Header),
			"body": requestBody(c),
		}).Infof(c.Request.Context(), "incoming http request")
	}
}

func LoggingResponse() Middle {
	return LoggingResponseWithLogger(logger.StandardLogger())
}

func LoggingResponseWithLogger(l logger.Logger) Middle {
	return func(c *gin.Context) {
		rw := &responseWriter{Body: new(bytes.Buffer), ResponseWriter: c.Writer}
		c.Writer = rw
		now := time.Now()

		c.Next()

		usedTime := time.Since(now)
		statusCode := c.Writer.Status()
		statusText := http.StatusText(statusCode)
		body := rw.Body.Bytes()
		cnt := len(body)
		if cnt > maxBodyLen {
			cnt = maxBodyLen
		}
		l.WithFields(logger.Fields{
			"status": fmt.Sprintf("%v %s", statusCode, statusText),
			"body":   string(body[:cnt]),
			//"headers": composeHeaders(c.Writer.Header()),
			"costms": usedTime.Milliseconds(),
		}).Infof(c.Request.Context(), "outgoing http response")
	}
}

func composeHeaders(headers http.Header) string {
	pairs := make([]string, 0, len(headers))
	for key, values := range headers {
		pairs = append(pairs, fmt.Sprintf("%s: %s", key, strings.Join(values, ", ")))
	}
	sort.Strings(pairs)
	return strings.Join(pairs, "; ")
}

func requestBody(c *gin.Context) string {
	if c.Request.Body == nil || c.Request.Body == http.NoBody {
		return ""
	}
	body, _ := ioutil.ReadAll(c.Request.Body)
	_ = c.Request.Body.Close()
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
	cnt := len(body)
	if cnt > maxBodyLen {
		cnt = maxBodyLen
	}
	return string(body[:cnt])
}

type responseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w responseWriter) Write(body []byte) (int, error) {
	w.Body.Write(body)
	return w.ResponseWriter.Write(body)
}
