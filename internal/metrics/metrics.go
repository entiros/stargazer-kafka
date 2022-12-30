package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var SyncCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "stargazer_sync_count",
	Help: "Number of sync operations performed",
})

var ErrCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "stargazer_error_count",
	Help: "Number of errors while performing sync",
})

func init() {
	prometheus.MustRegister(SyncCount)
	prometheus.MustRegister(ErrCount)

}

func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
