package news

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/thrasher-corp/gocryptotrader/common"
)

const (
	pathBitcoinist    = "https://bitcoinist.com/feed/"
	pathCoingape      = "https://coingape.com/feed/"
	pathCoindesk      = "https://www.coindesk.com/feed"
	pathCointelegraph = "https://cointelegraph.com/feed"
	pathMicky         = "https://micky.com.au/feed/"
	pathCNN           = "https://www.ccn.com/feed"
	pathNulltx        = "https://nulltx.com/feed/"
)

// GetData gets data from a given path
func GetData() error {
	regexpStr := "<title>[^\n]+</title>\n	<link>"
	r, err := regexp.Compile(regexpStr)
	if err != nil {
		return err
	}
	allPaths := []string{pathBitcoinist}
	for x := range allPaths {
		a, err := common.SendHTTPRequest(http.MethodGet, allPaths[x], nil, nil)
		if err != nil {
			return err
		}
		arr := r.FindAllString(a, -1)
		log.Println(arr)
		log.Println(a)
	}
	return nil
}

// Check checks things
func Check() error {
	var XML = []byte(`
	<Channel>
	<title>test</title>
			<hello>test123</hello>
	<link>why</link>
	<link>thisisatest</link>
	</Channel>
	`)
	var c Channel
	err := xml.Unmarshal(XML, &c)
	if err != nil {
		return err
	}
	log.Println(XML)
	fmt.Println(c)
	return nil
}

// CheckOtherThings checks other things
func CheckOtherThings() error {
	a, err := common.SendHTTPRequest(http.MethodGet, "https://coingape.com/feed/", nil, nil)
	if err != nil {
		return err
	}
	var q Query
	err = xml.Unmarshal([]byte(a), &q)
	if err != nil {
		log.Println(err)
	}
	log.Println(q)
	return nil
}
