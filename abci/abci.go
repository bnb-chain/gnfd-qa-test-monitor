package abci

import (
	"fmt"
	"github.com/bnb-chain/gnfd-qa-test-monitor/utils"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

func LastBlockHeight() (int64, error) {
	url := `https://gnfd-testnet-fullnode-tendermint-us.bnbchain.org/abci_info?last_block_height`
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer utils.CloseBody(resp.Body)
	body, err := io.ReadAll(resp.Body)
	result := gjson.GetBytes(body, "result.response.last_block_height")
	return result.Int(), nil
}

func BsDBInfoBlockHeight(spHost string, height int64) (string, error) {
	// https://gnfd-testnet-sp1.bnbchain.org/?bsdb-info&block_height=11037600
	url := fmt.Sprintf("https://%v/?bsdb-info&block_height=%v", spHost, height)
	//fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer utils.CloseBody(resp.Body)
	body, err := io.ReadAll(resp.Body)
	return string(body), nil
}
