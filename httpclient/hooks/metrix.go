package hooks

import (
	"net/url"

	"github.com/prometheus/client_golang/prometheus"
	resty "gopkg.in/resty.v1"
)

func init() {
	prometheus.MustRegister(latencySummary)
}

var (
	latencySummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "xd_sdk_httpclient_request_latency",
		Help: "Request latency of request that httpclient made",
	}, []string{"url"})
)

func LatencyMetrics() AfterResponseHook {
	return func(c *resty.Client, r *resty.Response) error {
		u, _ := url.Parse(r.Request.URL)
		latencySummary.With(prometheus.Labels{
			"url": u.Host + u.Path, // exclude param from url
		}).Observe(float64(r.ReceivedAt().Sub(r.Request.Time).Milliseconds()))
		return nil
	}
}
