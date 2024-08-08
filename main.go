package main

import (
	"fmt"
	"github.com/bnb-chain/gnfd-qa-test-monitor/checks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func recordMetrics() {
	checkBlock := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sp_db_object_nums_period_check_block_height",
	})
	checkErrCode := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sp_db_object_nums_period_check_error_code",
	})

	go func() {
		for {
			block, result := checks.CheckSpDbObjectNumsPeriod(checks.SelfTestNetSpHost)
			fmt.Printf("block: %v, result: %v \n", block, result)
			checkBlock.Set(float64(block))
			checkErrCode.Set(float64(result))
			time.Sleep(time.Minute * 10)
		}
	}()
}

func main() {
	recordMetrics()
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":24367", nil)
	if err != nil {
		return
	}
}
