package checks

import (
	"fmt"
	"github.com/bnb-chain/gnfd-qa-test-monitor/abci"
	"github.com/bnb-chain/gnfd-qa-test-monitor/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/tidwall/gjson"
	"strings"
)

type Code uint32

const (
	OK Code = iota
	GetBlockHeightErr
	GetObjectTotalCountErr
	GetObjectSealCountErr
	CheckObjectTotalCountErr
	CheckObjectSealCountErr
)

const (
	TestNetRpc = "https://gnfd-testnet-fullnode-tendermint-us.bnbchain.org:443"
	MainNetRpc = "https://greenfield-chain.bnbchain.org:443"
)

var (
	MainNetSpHosts = []string{
		"greenfield-sp.bnbchain.org",
		"greenfield-sp.defibit.io",
		"greenfield-sp.ninicoin.io",
		"greenfield-sp.nariox.org",
		"greenfield-sp.lumibot.org",
		"greenfield-sp.voltbot.io",
		"greenfield-sp.nodereal.io",
	}

	TestNetSpHosts = []string{
		"gnfd-testnet-sp1.bnbchain.org",
		"gnfd-testnet-sp2.bnbchain.org",
		"gnfd-testnet-sp3.bnbchain.org",
		"gnfd-testnet-sp4.bnbchain.org",
		"gnfd-testnet-sp1.nodereal.io",
		"gnfd-testnet-sp2.nodereal.io",
		"gnfd-testnet-sp3.nodereal.io",
	}
)

type DbShard struct {
	checkEnv         string
	checkRpc         string
	checkSpHosts     []string
	blockMetrics     prometheus.Gauge
	spErrCodeMetrics []prometheus.Gauge
}

func NewCheckDbShard(checkEnv, checkRpc string, checkSpHosts []string) *DbShard {
	checkSpErrCodeMetrics := make([]prometheus.Gauge, len(checkSpHosts))
	for i, spHost := range checkSpHosts {
		metricsSpHost := strings.Replace(spHost, "-", "_", -1)
		metricsSpHost = strings.Replace(metricsSpHost, ".", "_", -1)
		checkSpErrCodeMetrics[i] = promauto.NewGauge(prometheus.GaugeOpts{Name: fmt.Sprintf("%v_sp_db_shard_error_code_%v", checkEnv, metricsSpHost)})
	}

	return &DbShard{
		checkEnv:         checkEnv,
		checkRpc:         checkRpc,
		checkSpHosts:     checkSpHosts,
		blockMetrics:     promauto.NewGauge(prometheus.GaugeOpts{Name: fmt.Sprintf("%v_sp_db_shard_check_block_height", checkEnv)}),
		spErrCodeMetrics: checkSpErrCodeMetrics,
	}
}

func (s *DbShard) CheckDbShard() {

	lastChainHeight, err := abci.LastBlockHeight(s.checkRpc)
	if err != nil {
		s.blockMetrics.Set(float64(GetBlockHeightErr))
		return
	}
	calcHeight := lastChainHeight / 3600 * 3600
	s.blockMetrics.Set(float64(calcHeight))

	objCountArr := make([][]gjson.Result, len(s.checkSpHosts))
	sealObjCountArr := make([][]gjson.Result, len(s.checkSpHosts))
	isErr := false
	for i, spHost := range s.checkSpHosts {
		objCount, sealCount, errCode := getSpDbData(spHost, calcHeight)
		if errCode != OK {
			s.spErrCodeMetrics[i].Set(float64(errCode))
			isErr = true
		}
		objCountArr[i] = objCount
		sealObjCountArr[i] = sealCount
	}

	if isErr {
		return
	}

	spIndex, errCode := checkDbData(objCountArr, sealObjCountArr)
	if errCode != OK {
		s.spErrCodeMetrics[spIndex].Set(float64(errCode))
		return
	}

	for _, metric := range s.spErrCodeMetrics {
		metric.Set(float64(OK))
	}
}

func getSpDbData(spHost string, height int64) (objCount, objSealCount []gjson.Result, errCode Code) {
	xmlResult, err := abci.BsDBInfoBlockHeight(spHost, height)
	if err != nil {
		return nil, nil, GetBlockHeightErr
	}

	objectResString := utils.GetXmlPath(xmlResult, "GfSpGetBsDBInfoResponse/ObjectTotalCount")
	if objectResString == "" {
		fmt.Printf("sp: %v, ObjectTotalCount error\n", spHost)
		return nil, nil, GetObjectTotalCountErr
	} else {
		objectTotalCount := gjson.Parse(objectResString).Array()
		objCount = objectTotalCount
	}

	objectSealResString := utils.GetXmlPath(xmlResult, "GfSpGetBsDBInfoResponse/ObjectSealCount")
	if objectSealResString == "" {
		fmt.Printf("sp: %v, ObjectSealCount error\n", spHost)
		return nil, nil, GetObjectSealCountErr
	} else {
		ObjectSealCount := gjson.Parse(objectSealResString).Array()
		objSealCount = ObjectSealCount
	}

	return objCount, objSealCount, OK
}

func checkDbData(spObjCounts, spObjSealCounts [][]gjson.Result) (spIndex int, errCode Code) {
	for i := 0; i < 64; i++ {
		sumObject := int64(0)
		sumSp1 := int64(0)
		for _, objectCount := range spObjCounts {
			sumObject = sumObject + objectCount[i].Int()
			sumSp1++
		}
		sumSealedObject := int64(0)
		sumSp2 := int64(0)
		for _, sealObjectCount := range spObjSealCounts {
			sumSealedObject = sumSealedObject + sealObjectCount[i].Int()
			sumSp2++
		}

		objectAverage := sumObject / sumSp1
		for spIndex, eachValue := range spObjCounts {
			if objectAverage != eachValue[i].Int() {
				return spIndex, CheckObjectTotalCountErr
			}
		}

		sealObjectAverage := sumSealedObject / sumSp2
		for _, eachValue := range spObjSealCounts {
			if sealObjectAverage != eachValue[i].Int() {
				return spIndex, CheckObjectSealCountErr
			}
		}
	}

	return 0, OK
}
