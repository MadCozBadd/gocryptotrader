package coindesk

import (
	"log"
	"net/http"

	"github.com/thrasher-corp/gocryptotrader/common"
)

const (
	path = "http://feeds.feedburner.com/CoinDesk"
)

// GetData gets data from a given path
func GetData() error {
	a, err := common.SendHTTPRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		return err
	}
	log.Println(a)
	return err
}
