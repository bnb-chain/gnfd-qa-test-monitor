package main

import (
	"github.com/bnb-chain/gnfd-qa-test-monitor/checks"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func recordMetricsV() {
	go func() {
		for {
			checks.CheckDbShard("testnet", checks.TestNetRpc, checks.TestNetSpHosts)
			time.Sleep(time.Minute * 10)
		}
	}()

	go func() {
		for {
			checks.CheckDbShard("mainnet", checks.MainNetRpc, checks.MainNetSpHosts)
			time.Sleep(time.Minute * 10)
		}
	}()
}

func main() {
	recordMetricsV()
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":24367", nil)
	if err != nil {
		return
	}
}
