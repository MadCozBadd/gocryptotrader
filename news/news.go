package news

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	pathCryptoDaily   = "https://cryptodaily.co.uk/feed/"
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
		pathCointelegraph, pathMicky, pathNulltx, pathCryptoDaily}

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
		words, err := ReadFile("checklist.json")
		if err != nil {
			return err
		}
		for o := range words {
			for y := range newStorage.Items {
				tempStrings := strings.Split(newStorage.Items[y].Title, " ")
				if common.StringDataCompareInsensitive(tempStrings, words[o]) {
					stuff := fmt.Sprintf("%s:\n%s",
						newStorage.Items[y].Title,
						newStorage.Items[y].Link)
					fmt.Println(stuff)
				}
			}
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

// WriteFile write to the json file
func WriteFile(moreWords []string, fileName string) error {
	words, err := ReadFile(fileName)
	if err != nil {
		return err
	}
	for x := range moreWords {
		if !common.StringDataCompareInsensitive(words, moreWords[x]) {
			words = append(words, moreWords[x])
		}
	}
	file, err := json.MarshalIndent(words, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, file, 0770)
}

// ReadFile reads the file
func ReadFile(fileName string) ([]string, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var words []string
	return words, json.Unmarshal(bytes, &words)
}
