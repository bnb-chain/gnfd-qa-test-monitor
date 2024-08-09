package main

import (
	"fmt"
	"github.com/bnb-chain/gnfd-qa-test-monitor/abci"
	"github.com/bnb-chain/gnfd-qa-test-monitor/checks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"
	"time"
)

func recordMetrics() {
	testNetRpc := "https://gnfd-testnet-fullnode-tendermint-us.bnbchain.org:443"
	mainNetRpc := "https://greenfield-chain.bnbchain.org:443"

	checkTestNetBlock := promauto.NewGauge(prometheus.GaugeOpts{Name: "testnet_sp_db_shard_check_block_height"})
	checkMainNetBlock := promauto.NewGauge(prometheus.GaugeOpts{Name: "mainnet_sp_db_shard_check_block_height"})

	checkTestNetSpErrCodes := make([]prometheus.Gauge, len(checks.TestNetSpHost))
	for i, spHost := range checks.TestNetSpHost {
		spHost = strings.Replace(spHost, "-", "_", -1)
		spHost = strings.Replace(spHost, ".", "_", -1)
		name := fmt.Sprintf("testnet_sp_db_shard_error_code_%v", spHost)
		checkTestNetSpErrCodes[i] = promauto.NewGauge(prometheus.GaugeOpts{Name: name})
	}

	checkMainNetSpErrCodes := make([]prometheus.Gauge, len(checks.TestNetSpHost))
	for i, spHost := range checks.TestNetSpHost {
		spHost = strings.Replace(spHost, "-", "_", -1)
		spHost = strings.Replace(spHost, ".", "_", -1)
		name := fmt.Sprintf("mainnet_sp_db_shard_error_code_%v", spHost)
		checkMainNetSpErrCodes[i] = promauto.NewGauge(prometheus.GaugeOpts{Name: name})
	}

	go func() {
		for {
			// check TestNet
			testNetChainHeight, err := abci.LastBlockHeight(testNetRpc)
			if err != nil {
				fmt.Println(err)
				continue
			}
			testNetCalcHeight := testNetChainHeight / 3600 * 3600
			checkTestNetBlock.Set(float64(testNetCalcHeight))
			for i, spHost := range checks.TestNetSpHost {
				errCode := checks.CheckDbShard(testNetCalcHeight, spHost)
				checkTestNetSpErrCodes[i].Set(float64(errCode))
			}

			// check MainNet
			mainNetChainHeight, err := abci.LastBlockHeight(mainNetRpc)
			if err != nil {
				fmt.Println(err)
				continue
			}
			mainNetCalcHeight := mainNetChainHeight / 3600 * 3600
			checkMainNetBlock.Set(float64(mainNetCalcHeight))
			for i, spHost := range checks.MainNetSpHost {
				errCode := checks.CheckDbShard(mainNetCalcHeight, spHost)
				checkMainNetSpErrCodes[i].Set(float64(errCode))
			}

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
