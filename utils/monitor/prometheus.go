package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yiqiang3344/go-lib/utils/config"
	cLog "github.com/yiqiang3344/go-lib/utils/log"
	"net/http"
)

var httpReqsHistory = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace:   "micro",
	Subsystem:   "",
	Name:        "srv_gateway_req_history",
	Help:        "Histogram of response latency (seconds) of http handlers.",
	ConstLabels: nil,
	Buckets:     nil,
}, []string{"method", "code"})

func InitPrometheus() *prometheus.HistogramVec {
	//配置网关请求监控
	http.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(
		httpReqsHistory,
	)
	go func() {
		err := http.ListenAndServe(config.GetCfgString("prometheus.address"), nil)
		if err != nil {
			cLog.FatalLog(err.Error(), "")
		}
	}()
	return httpReqsHistory
}
