package main

import (
	"github.com/bnb-chain/gnfd-qa-test-monitor/checks"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func recordMetrics() {
	testNet := checks.NewCheckDbShard("testnet", checks.TestNetRpc, checks.TestNetSpHosts)
	mainNet := checks.NewCheckDbShard("mainnet", checks.MainNetRpc, checks.MainNetSpHosts)
	go func() {
		for {
			testNet.CheckDbShard()
			mainNet.CheckDbShard()
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
