package news

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
)

const (
	pathBitcoinist    = "https://bitcoinist.com/feed/"
	pathCoingape      = "https://coingape.com/feed/"
	pathCoindesk      = "https://www.coindesk.com/feed"
	pathCointelegraph = "https://cointelegraph.com/feed"
	pathMicky         = "https://micky.com.au/feed/"
	pathCCN           = "https://www.ccn.com/feed"
	pathNulltx        = "https://nulltx.com/feed/"
	pathSendMessage   = "https://slack.com/api/chat.postMessage"
)

func main() {
	for {
		err := CheckOtherThings()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("HELLLLLLLLOOOOOOOOOOOO\n\n\n\n")
		time.Sleep(time.Minute)
	}
}

// CheckOtherThings checks other things
func CheckOtherThings() error {
	updateTimes := make(map[string]string)
	allPaths := []string{pathBitcoinist, pathCCN, pathCoindesk, pathCoingape,
		pathCointelegraph, pathMicky, pathNulltx}
	for z := range allPaths {
		var newStorage Storage
		newStorage.Path = allPaths[z]
		a, err := common.SendHTTPRequest(http.MethodGet, allPaths[z], nil, nil)
		if err != nil {
			return err
		}
		var q Query
		err = xml.Unmarshal([]byte(a), &q)
		if err != nil {
			log.Fatal(err)
		}
		for y := range q.Channel.Items {
			newStorage.Items = append(newStorage.Items, q.Channel.Items[y])
		}
		_, ok := updateTimes[allPaths[z]]
		if ok {
			for x := range newStorage.Items {
				if newStorage.Items[x].PubTime == updateTimes[allPaths[z]] {
					newStorage.Items = newStorage.Items[:x]
				}
			}
		}
		for y := range newStorage.Items {
			stuff := fmt.Sprintf("%s:\n%s",
				newStorage.Items[y].Title,
				newStorage.Items[y].Link)
			fmt.Println(stuff)
		}
		updateTimes[allPaths[z]] = newStorage.Items[0].PubTime
	}
	return nil
}

// SendMessage sends message to the slack channel
func SendMessage(message string) error {
	headers := make(map[string]string)
	headers["Content-type"] = "application/json"
	headers["Authorization"] = "NO"
	var params AuthParams
	params.Channel = "GTP4246MB"
	params.Text = message
	b, err := json.Marshal(params)
	if err != nil {
		return err
	}
	a, err := common.SendHTTPRequest(http.MethodPost,
		pathSendMessage,
		headers,
		bytes.NewBuffer(b))
	log.Println(a)
	return err
}
