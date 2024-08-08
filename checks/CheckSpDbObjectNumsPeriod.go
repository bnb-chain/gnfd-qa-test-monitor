package checks

import (
	"fmt"
	"github.com/bnb-chain/gnfd-qa-test-monitor/abci"
	"github.com/bnb-chain/gnfd-qa-test-monitor/utils"
	"github.com/tidwall/gjson"
)

var SelfMainNetSpHost = []string{
	"greenfield-sp.bnbchain.org",
	"greenfield-sp.defibit.io",
	"greenfield-sp.ninicoin.io",
	"greenfield-sp.nariox.org",
	"greenfield-sp.lumibot.org",
	"greenfield-sp.voltbot.io",
	"greenfield-sp.nodereal.io",
}

var SelfTestNetSpHost = []string{
	"gnfd-testnet-sp1.bnbchain.org",
	"gnfd-testnet-sp2.bnbchain.org",
	"gnfd-testnet-sp3.bnbchain.org",
	"gnfd-testnet-sp4.bnbchain.org",
	"gnfd-testnet-sp1.nodereal.io",
	"gnfd-testnet-sp2.nodereal.io",
	"gnfd-testnet-sp3.nodereal.io",
}

type Code float64

const (
	OK Code = iota
	GetObjectTotalCountErr
	GetObjectSealCountErr
	CheckObjectTotalCountErr
	CheckObjectSealCountErr
)

func CheckSpDbObjectNumsPeriod(spHostArray []string) (block int64, errCode Code) {
	chainHeight, err := abci.LastBlockHeight()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("height: %d\n", chainHeight)
	calcHeight := chainHeight / 3600 * 3600
	fmt.Printf("calcHeight: %d\n", calcHeight)

	resObjectCount := make(map[string][]gjson.Result)
	resObjectSealCount := make(map[string][]gjson.Result)

	// get everyone sp object count
	for _, spHost := range spHostArray {
		xmlResult, err := abci.BsDBInfoBlockHeight(spHost, calcHeight)
		if err != nil {
			fmt.Println(err)
		}

		objectResString := utils.GetXmlPath(xmlResult, "GfSpGetBsDBInfoResponse/ObjectTotalCount")
		if objectResString == "" {
			fmt.Printf("sp: %v, ObjectTotalCount error\n", spHost)
			return calcHeight, GetObjectTotalCountErr
		} else {
			ObjectTotalCount := gjson.Parse(objectResString).Array()
			resObjectCount[spHost] = ObjectTotalCount
		}

		objectSealResString := utils.GetXmlPath(xmlResult, "GfSpGetBsDBInfoResponse/ObjectSealCount")
		if objectSealResString == "" {
			fmt.Printf("sp: %v, ObjectSealCount error\n", spHost)
			return calcHeight, GetObjectSealCountErr
		} else {
			ObjectSealCount := gjson.Parse(objectSealResString).Array()
			resObjectSealCount[spHost] = ObjectSealCount
		}
	}

	// check sp object count
	for i := 0; i < 64; i++ {
		sumObject := int64(0)
		sumSp1 := int64(0)
		for _, objectCount := range resObjectCount {
			sumObject = sumObject + objectCount[i].Int()
			sumSp1++
		}
		sumSealedObject := int64(0)
		sumSp2 := int64(0)
		for _, sealObjectCount := range resObjectSealCount {
			sumSealedObject = sumSealedObject + sealObjectCount[i].Int()
			sumSp2++
		}

		objectAverage := sumObject / sumSp1
		sealObjectAverage := sumSealedObject / sumSp2
		for _, eachValue := range resObjectCount {
			if objectAverage != eachValue[i].Int() {
				return calcHeight, CheckObjectTotalCountErr
			}
		}
		for _, eachValue := range resObjectSealCount {
			if sealObjectAverage != eachValue[i].Int() {
				return calcHeight, CheckObjectSealCountErr
			}
		}
	}

	return calcHeight, OK
}
