package main

import (
	"github.com/bnb-chain/gnfd-qa-test-monitor/checks"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func recordMetrics() {
	go func() {
		for {
			checks.CheckDbShard("testnet", checks.TestNetRpc, checks.TestNetSpHosts)
			checks.CheckDbShard("mainnet", checks.MainNetRpc, checks.MainNetSpHosts)
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
