package middles

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/yunxiaoyang01/open_sdk/logger"

	"github.com/gin-gonic/gin"
)

const (
	HEALTHAPI       = "/health"
	HELATHSWITCHAPI = "/health_switch"
)

var DefaultHealthCheck = &HealthSwitch{}

func WrapHealthCheck(engine *gin.Engine) {
	engine.GET(HEALTHAPI, DefaultHealthCheck.Health)
	engine.GET(HELATHSWITCHAPI, DefaultHealthCheck.Switch)
}

func HealthCheck() Middle {
	return func(c *gin.Context) {
		// add health check, to avoid 502 bad gateway error
		switch c.Request.URL.Path {
		case HEALTHAPI:
			DefaultHealthCheck.Health(c)
			c.Abort()
			return
		case HELATHSWITCHAPI:
			DefaultHealthCheck.Switch(c)
			c.Abort()
			return
		}
		logger.Debugf(c.Request.Context(), "get path:%v", c.Request.URL.Path)
	}
}

type HealthSwitch struct {
	disable bool
	sync.RWMutex
}

func (h *HealthSwitch) isLocal(c *gin.Context) bool {
	const localIp = "127.0.0.1"
	ip, _, _ := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	return ip == localIp
}

func (h *HealthSwitch) Health(c *gin.Context) {
	h.RLock()
	defer h.RUnlock()

	if !h.disable {
		c.Status(http.StatusOK)
		return
	}

	c.Status(http.StatusForbidden)
	return
}

func (h *HealthSwitch) Switch(c *gin.Context) {
	if !h.isLocal(c) {
		c.Status(http.StatusForbidden)
		return
	}

	enable := c.DefaultQuery("switch", "")
	if enable == "false" {
		h.Lock()
		h.disable = true
		h.Unlock()
	}

	sleepString := c.DefaultQuery("sleep", "")
	sleep, _ := strconv.Atoi(sleepString)

	time.Sleep(time.Duration(sleep) * time.Second)

	c.Status(http.StatusOK)
}
