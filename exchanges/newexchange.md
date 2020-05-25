# GoCryptoTrader package Exchanges

<img src="https://github.com/thrasher-corp/gocryptotrader/blob/master/web/src/assets/page-logo.png?raw=true" width="350px" height="350px" hspace="70">

[![Build Status](https://travis-ci.org/thrasher-corp/gocryptotrader.svg?branch=master)](https://travis-ci.org/thrasher-corp/gocryptotrader)
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/thrasher-corp/gocryptotrader/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/thrasher-corp/gocryptotrader?status.svg)](https://godoc.org/github.com/thrasher-corp/gocryptotrader/exchanges)
[![Coverage Status](http://codecov.io/github/thrasher-corp/gocryptotrader/coverage.svg?branch=master)](http://codecov.io/github/thrasher-corp/gocryptotrader?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/thrasher-corp/gocryptotrader)](https://goreportcard.com/report/github.com/thrasher-corp/gocryptotrader)

This exchanges package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Current Features for exchanges

+ This package is used to connect and query data from supported exchanges.

+ Please checkout individual exchange README for more information on
implementation

#### How to add a new exchange

+ 1) run exchange_template.go which automatically creates files & inbuilt functions

###### Linux/OSX
GoCryptoTrader is built using [Go Modules](https://github.com/golang/go/wiki/Modules) and requires Go 1.11 or above
Using Go Modules you now clone this repository **outside** your GOPATH

```bash
git clone https://github.com/thrasher-corp/gocryptotrader.git
cd cmd
cd exchange_template
go build
go run exchange_template.go -name Bitmex -ws -rest
```

###### Windows

```bash
git clone https://github.com/thrasher-corp/gocryptotrader.git
cd gocryptotrader\cmd\exchange_template
go build
go run exchange_template.go -name Bitmex -ws -rest
```

+ 2) add exchange struct to config_example.json, configtest.json (in testdata) & to the main config

###### If main config path is unknown the following function can be used:
```go
config.GetDefaultFilePath()
```

```go
  {
   "name": "FTX",
   "enabled": true,
   "verbose": false,
   "httpTimeout": 15000000000,
   "websocketResponseCheckTimeout": 30000000,
   "websocketResponseMaxLimit": 7000000000,
   "websocketTrafficTimeout": 30000000000,
   "websocketOrderbookBufferLimit": 5,
   "baseCurrencies": "USD",
   "currencyPairs": {
    "requestFormat": {
     "uppercase": false,
     "delimiter": "_"
    },
    "configFormat": {
     "uppercase": true,
     "delimiter": "_"
    },
    "useGlobalFormat": true,
    "assetTypes": [
      "spot",
      "futures"
     ],
     "pairs": {
      "futures": {
       "enabled": "BTC-PERP",
       "available": "BTC-PERP",
       "requestFormat": {
        "uppercase": true,
        "delimiter": "-"
       },
       "configFormat": {
        "uppercase": true,
        "delimiter": "-"
       }
      },
      "spot": {
       "enabled": "BTC/USD",
       "available": "BTC/USD",
       "requestFormat": {
        "uppercase": true,
        "delimiter": "/"
       },
       "configFormat": {
        "uppercase": true,
        "delimiter": "/"
       }
      }
     }
    },
   "api": {
    "authenticatedSupport": false,
    "authenticatedWebsocketApiSupport": false,
    "endpoints": {
     "url": "NON_DEFAULT_HTTP_LINK_TO_EXCHANGE_API",
     "urlSecondary": "NON_DEFAULT_HTTP_LINK_TO_EXCHANGE_API",
     "websocketURL": "NON_DEFAULT_HTTP_LINK_TO_WEBSOCKET_EXCHANGE_API"
    },
    "credentials": {
     "key": "Key",
     "secret": "Secret"
    },
    "credentialsValidator": {
     "requiresKey": true,
     "requiresSecret": true
    }
   },
   "features": {
    "supports": {
     "restAPI": true,
     "restCapabilities": {
      "tickerBatching": true,
      "autoPairUpdates": true
     },
     "websocketAPI": true,
     "websocketCapabilities": {}
    },
    "enabled": {
     "autoPairUpdates": true,
     "websocketAPI": false
    }
   },
   "bankAccounts": [
    {
     "enabled": false,
     "bankName": "",
     "bankAddress": "",
     "bankPostalCode": "",
     "bankPostalCity": "",
     "bankCountry": "",
     "accountName": "",
     "accountNumber": "",
     "swiftCode": "",
     "iban": "",
     "supportedCurrencies": ""
    }
   ]
  },
```

###### Available pairs will be automatically filled in the configs when wrapper functions are filled out and gocryptotrader is run with the new exchange enabled

+ 3) Add the currency structs in ftx_wrapper.go:

###### Futures currency support:

Spot pairs' support is inbuilt, for other asset types, struct needs to be manually created

```go
	spot := currency.PairStore{
		RequestFormat: &currency.PairFormat{
			Uppercase: true,
			Delimiter: "/",
		},
		ConfigFormat: &currency.PairFormat{
			Uppercase: true,
			Delimiter: "/",
		},
	}
	futures := currency.PairStore{
		RequestFormat: &currency.PairFormat{
			Uppercase: true,
			Delimiter: "-",
		},
		ConfigFormat: &currency.PairFormat{
			Uppercase: true,
			Delimiter: "-",
		},
	}
```

+ 4) Document the addition of the new exchange (FTX exchange is used as an example below):

###### root Readme.md:
```go
| Exchange | REST API | Streaming API | FIX API |
|----------|------|-----------|-----|
| Alphapoint | Yes  | Yes        | NA  |
| Binance| Yes  | Yes        | NA  |
| Bitfinex | Yes  | Yes        | NA  |
| Bitflyer | Yes  | No      | NA  |
| Bithumb | Yes  | NA       | NA  |
| BitMEX | Yes | Yes | NA |
| Bitstamp | Yes  | Yes       | No  |
| Bittrex | Yes | No | NA |
| BTCMarkets | Yes | No       | NA  |
| BTSE | Yes | Yes | NA |
| COINUT | Yes | Yes | NA |
| Exmo | Yes | NA | NA |
| FTX | Yes | Yes | No |
| CoinbasePro | Yes | Yes | No|
| Coinbene | Yes | No | No |
| GateIO | Yes | Yes | NA |
| Gemini | Yes | Yes | No |
| HitBTC | Yes | Yes | No |
| Huobi.Pro | Yes | Yes | NA |
| ItBit | Yes | NA | No |
| Kraken | Yes | Yes | NA |
| Lbank | Yes | No | NA |
| LakeBTC | Yes | No | NA |
| LocalBitcoins | Yes | NA | NA |
| OKCoin International | Yes | Yes | No |
| OKEX | Yes | Yes | No |
| Poloniex | Yes | Yes | NA |
| Yobit | Yes | NA | NA |
| ZB.COM | Yes | Yes | NA |
```

###### exchanges\support.go:
```go
var Exchanges = []string{
	"binance",
	"bitfinex",
	"bitflyer",
	"bithumb",
	"bitmex",
	"bitstamp",
	"bittrex",
	"btc markets",
	"btse",
	"coinbasepro",
	"coinbene",
	"coinut",
	"exmo",
	"ftx",
	"gateio",
	"gemini",
	"hitbtc",
	"huobi",
	"itbit",
	"kraken",
	"lakebtc",
	"lbank",
	"localbitcoins",
	"okcoin international",
	"okex",
	"poloniex",
	"yobit",
    "zb",
```

###### exchanges\exchange_test.go:
```go
func TestExchange_Exchanges(t *testing.T) {
	t.Parallel()
	x := exchangeTest.Exchanges(false)
	y := len(x)
	if y != 28 { // add 1 here (before FTX was added it was 27, so 28 now)
		t.Fatalf("expected 28 received %v", y) // add 1 here
	}
}
```

###### cmd\documentation\exchange_templates:

- Create a new file named <exchangename>.tmpl
- Copy contents of template from another exchange example here being Exmo
- Replace names and variables as shown:

```go
{{define "exchanges exmo" -}} // exmo -> ftx
{{template "header" .}}
## Exmo Exchange

### Current Features

+ REST Support // if websocket or fix are supported, add that in too
```

```go
main.go
var e exchange.IBotExchange // e -> f

for i := range bot.Exchanges {
  if bot.Exchanges[i].GetName() == "Exmo" { // Exmo -> FTX
    e = bot.Exchanges[i] // e -> f
  }
}

// Public calls - wrapper functions

// Fetches current ticker information
tick, err := e.FetchTicker() // e -> f 
if err != nil {
  // Handle error
}

// Fetches current orderbook information
ob, err := e.FetchOrderbook() // e -> f (do so for the rest of the functions too)
if err != nil {
  // Handle error
}
```

- Run documentation.go to generate readme file for the exchange:
```bash
cd gocryptotrader\cmd\documentation
go build
go run documentation.exe
```

This will generate a readme file for the exchange which can be found in the new exchange's folder

+ 5) Create functions supported by the exchange:

###### Requester functions:

```go
// SendHTTPRequest sends an unauthenticated HTTP request
func (f *FTX) SendHTTPRequest(path string, result interface{}) error {
	return f.SendPayload(context.Background(), &request.Item{
		Method:        http.MethodGet,
		Path:          path,
		Result:        result,
		Verbose:       f.Verbose,
		HTTPDebugging: f.HTTPDebugging,
		HTTPRecording: f.HTTPRecording,
	})
}

Authenticated request function is created based on the way the exchange documentation specifies: https://docs.ftx.com/#authentication
// SendAuthHTTPRequest sends an authenticated request
func (f *FTX) SendAuthHTTPRequest(method, path string, data, result interface{}) error {
	ts := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	var body io.Reader
	var hmac, payload []byte
	var err error
	if data != nil {
		payload, err = json.Marshal(data)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(payload)
		sigPayload := ts + method + "/api" + path + string(payload)
		hmac = crypto.GetHMAC(crypto.HashSHA256, []byte(sigPayload), []byte(f.API.Credentials.Secret))
	} else {
		sigPayload := ts + method + "/api" + path
		hmac = crypto.GetHMAC(crypto.HashSHA256, []byte(sigPayload), []byte(f.API.Credentials.Secret))
	}
	headers := make(map[string]string)
	headers["FTX-KEY"] = f.API.Credentials.Key
	headers["FTX-SIGN"] = crypto.HexEncodeToString(hmac)
	headers["FTX-TS"] = ts
	headers["Content-Type"] = "application/json"
	return f.SendPayload(context.Background(), &request.Item{
		Method:        method,
		Path:          ftxAPIURL + path,
		Headers:       headers,
		Body:          body,
		Result:        result,
		AuthRequest:   true,
		Verbose:       f.Verbose,
		HTTPDebugging: f.HTTPDebugging,
		HTTPRecording: f.HTTPRecording,
	})
}
```

###### Unauthenticated Functions:

https://docs.ftx.com/#get-markets

Create a type struct in types.go for the response type shown on the documentation website:

```go
// MarketData stores market data
type MarketData struct {
	Name           string  `json:"name"`
	BaseCurrency   string  `json:"baseCurrency"`
	QuoteCurrency  string  `json:"quoteCurrency"`
	MarketType     string  `json:"type"`
	Underlying     string  `json:"underlying"`
	Enabled        bool    `json:"enabled"`
	Ask            float64 `json:"ask"`
	Bid            float64 `json:"bid"`
	Last           float64 `json:"last"`
	PriceIncrement float64 `json:"priceIncrement"`
	SizeIncrement  float64 `json:"sizeIncrement"`
}
```

Create a new variable, they are created at the top of ftx.go file:
```go
const (
	ftxAPIURL = "https://ftx.com/api"

	// Public endpoints
	getMarkets           = "/markets"
	getMarket            = "/markets/"
	getOrderbook         = "/markets/%s/orderbook?depth=%s"
	getTrades            = "/markets/%s/trades?"
	getHistoricalData    = "/markets/%s/candles?"
	getFutures           = "/futures"
	getFuture            = "/futures/"
	getFutureStats       = "/futures/%s/stats"
	getFundingRates      = "/funding_rates"
  getAllWalletBalances = "/wallet/all_balances"
  ```

Create a get function in ftx.go file and unmarshall the data in the created type
```go
// GetMarkets gets market data
func (f *FTX) GetMarkets() (Markets, error) {
	var resp Markets
	return resp, f.SendHTTPRequest(ftxAPIURL+getMarkets, &resp)
}
```

Create a test function in ftx_test.go to see if the data is received and unmarshalled correctly
```go
func TestGetMarket(t *testing.T) {
  t.Parallel() // Tests have a 30s timeout set in GCT so please add t.Parallel() so when the test package is run, testing is done much faster
	a, err := f.GetMarket(spotPair) // spotPair is just a const set for ease of use
  t.Log(a)
	if err != nil {
		t.Error(err)
	}
}
```
Verbose can be set to true to see the data received if there are errors unmarshalling
Once testing is done remove verbose, variable a and t.Log(a) since they produce unnecessary output when GCT is run
```go
_, err := f.GetMarket(spotPair)
```

Create the rest of the unauthenticated functions and their tests similarly

###### Authenticated functions:

For authenticated functions to work, authenticated request function should be configured correctly
Create authenticated functions and test along the way similar to the functions above:

https://docs.ftx.com/#get-account-information:

```go
// GetAccountInfo gets account info
func (f *FTX) GetAccountInfo() (AccountData, error) {
	var resp AccountData
	return resp, f.SendAuthHTTPRequest(http.MethodGet, getAccountInfo, nil, &resp)
}
```

Get Request params for authenticated requests are sent through url.Values{}:

https://docs.ftx.com/#get-withdrawal-history:

```go
// GetTriggerOrderHistory gets trigger orders that are currently open
func (f *FTX) GetTriggerOrderHistory(marketName string, startTime, endTime time.Time, side, orderType, limit string) (TriggerOrderHistory, error) {
	var resp TriggerOrderHistory
	params := url.Values{}
	if marketName != "" {
		params.Set("market", marketName)
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		params.Set("start_time", strconv.FormatInt(startTime.Unix(), 10))
		params.Set("end_time", strconv.FormatInt(endTime.Unix(), 10))
		if !startTime.Before(endTime) {
			return resp, errors.New("startTime cannot be bigger than endTime")
		}
	}
	if side != "" {
		params.Set("side", side)
	}
	if orderType != "" {
		params.Set("type", orderType)
	}
	if limit != "" {
		params.Set("limit", limit)
	}
	return resp, f.SendAuthHTTPRequest(http.MethodGet, getTriggerOrderHistory+params.Encode(), nil, &resp)
}
```

For post or delete requests params are sent through map[string]interface{}:

https://docs.ftx.com/#place-order


Structs for unmarshalling the data are made exactly the same way as the previous functions
```go
type OrderData struct {
	CreatedAt     time.Time `json:"createdAt"`
	FilledSize    float64   `json:"filledSize"`
	Future        string    `json:"future"`
	ID            int64     `json:"id"`
	Market        string    `json:"market"`
	Price         float64   `json:"price"`
	AvgFillPrice  float64   `json:"avgFillPrice"`
	RemainingSize float64   `json:"remainingSize"`
	Side          string    `json:"side"`
	Size          float64   `json:"size"`
	Status        string    `json:"status"`
	OrderType     string    `json:"type"`
	ReduceOnly    bool      `json:"reduceOnly"`
	IOC           bool      `json:"ioc"`
	PostOnly      bool      `json:"postOnly"`
	ClientID      string    `json:"clientId"`
}

// PlaceOrder stores data of placed orders
type PlaceOrder struct {
	Success bool      `json:"success"`
	Result  OrderData `json:"result"`
}
```

```go
// Order places an order
func (f *FTX) Order(marketName, side, orderType, reduceOnly, ioc, postOnly, clientID string, price, size float64) (PlaceOrder, error) {
	req := make(map[string]interface{})
	req["market"] = marketName
	req["side"] = side
	req["price"] = price
	req["type"] = orderType
	req["size"] = size
	if reduceOnly != "" {
		req["reduceOnly"] = reduceOnly
	}
	if ioc != "" {
		req["ioc"] = ioc
	}
	if postOnly != "" {
		req["postOnly"] = postOnly
	}
	if clientID != "" {
		req["clientID"] = clientID
	}
	var resp PlaceOrder
	return resp, f.SendAuthHTTPRequest(http.MethodPost, placeOrder, req, &resp)
}
```

+ 6) Implementing wrapper functions:

Wrapper functions are the interface through which GCT bot communicates with exchange for gathering and sending data
The exchanges may not support all the functionality in the wrapper, so fill out the ones that are supported as shown in the examples below

```go
// FetchTradablePairs returns a list of the exchanges tradable pairs
func (f *FTX) FetchTradablePairs(a asset.Item) ([]string, error) {
	if !f.SupportsAsset(a) {
		return nil, fmt.Errorf("asset type of %s is not supported by %s", a, f.Name)
	}
	markets, err := f.GetMarkets()
	if err != nil {
		return nil, err
	}
	var pairs []string
	switch a {
	case asset.Spot:
		for x := range markets.Result {
			if markets.Result[x].MarketType == spotString {
				pairs = append(pairs, markets.Result[x].Name)
			}
		}
	case asset.Futures:
		for x := range markets.Result {
			if markets.Result[x].MarketType == futuresString {
				pairs = append(pairs, markets.Result[x].Name)
			}
		}
	}
	return pairs, nil
}
```

Wrapper functions on most exchanges are written in similar ways so other exchanges can be used as a referrence

Alot of useful helper methods can be found in exchange.go, some examples are given below:

```go
f.FormatExchangeCurrency(p, a) // Formats the currency pair to the style accepted by the exchange. p is the currency pair & a is the asset type

f.SupportsAsset(a) // Checks is asset type is supported by the bot

f.GetPairAssetType(p) // Returns the asset type of currency pair p
```

Currency package also has alot of helpful methods

+ 6) Websocket addition if exchange supports it:

###### Add websocket to exchange struct in ftx.go

```go
// FTX is the overarching type across this package
type FTX struct {
	exchange.Base
	WebsocketConn *wshandler.WebsocketConnection // Add this line
}
```

###### Create functions as explained in the documentation:

- Set the websocket url in ftx_websocket.go that is provided in the documentation:

	ftxWSURL          = "wss://ftx.com/ws/"

```go
	ftxWSURL          = "wss://ftx.com/ws/"
```

- Set channel names as variables for ease of use:

```go
	wsTicker          = "ticker"
	wsTrades          = "trades"
	wsOrderbook       = "orderbook"
	wsMarkets         = "markets"
	wsFills           = "fills"
	wsOrders          = "orders"
	wsUpdate          = "update"
	wsPartial         = "partial"
```

- Create types given in the documentation to unmarshall the streamed data:

https://docs.ftx.com/#fills-2

```go
// WsFills stores websocket fills' data
type WsFills struct {
	Fee       float64   `json:"fee"`
	FeeRate   float64   `json:"feeRate"`
	Future    string    `json:"future"`
	ID        int64     `json:"id"`
	Liquidity string    `json:"liquidity"`
	Market    string    `json:"market"`
	OrderID   int64     `json:"int64"`
	TradeID   int64     `json:"tradeID"`
	Price     float64   `json:"price"`
	Side      string    `json:"side"`
	Size      float64   `json:"size"`
	Time      time.Time `json:"time"`
	OrderType string    `json:"orderType"`
}

// WsFillsDataStore stores ws fills' data
type WsFillsDataStore struct {
	Channel     string  `json:"channel"`
	MessageType string  `json:"type"`
	FillsData   WsFills `json:"fills"`
}
```

- Create the authentication function based on specifications provided in the documentation:

https://docs.ftx.com/#private-channels

```go
// WsAuth sends an authentication message to receive auth data
func (f *FTX) WsAuth() error {
	intNonce := time.Now().UnixNano() / 1000000
	strNonce := strconv.FormatInt(intNonce, 10)
	hmac := crypto.GetHMAC(
		crypto.HashSHA256,
		[]byte(strNonce+"websocket_login"),
		[]byte(f.API.Credentials.Secret),
	)
	sign := crypto.HexEncodeToString(hmac)
	req := Authenticate{Operation: "login",
		Args: AuthenticationData{
			Key:  f.API.Credentials.Key,
			Sign: sign,
			Time: intNonce,
		},
	}
	return f.WebsocketConn.SendJSONMessage(req)
}
```

- Create function to generate default subscriptions

```go
// GenerateDefaultSubscriptions generates default subscription
func (f *FTX) GenerateDefaultSubscriptions() {
	var channels = []string{wsTicker, wsTrades, wsOrderbook, wsMarkets, wsFills, wsOrders} // All the channels that can be subscribed to
	var subscriptions []wshandler.WebsocketChannelSubscription
	for a := range f.CurrencyPairs.AssetTypes {
		pairs := f.GetEnabledPairs(f.CurrencyPairs.AssetTypes[a])
		newPair := currency.NewPairWithDelimiter(pairs[0].Base.String(), pairs[0].Quote.String(), "-")
		for x := range channels {
			subscriptions = append(subscriptions, wshandler.WebsocketChannelSubscription{
				Channel:  channels[x],
				Currency: newPair,
			})
		}
	}
	f.Websocket.SubscribeToChannels(subscriptions)
}
```

- Create subscribe function with the data provided by the exchange documentation:

https://docs.ftx.com/#request-process

Create a struct required to subscribe to channels:

```go
// WsSub has the data used to subscribe to a channel
type WsSub struct {
	Channel   string `json:"channel,omitempty"`
	Market    string `json:"market,omitempty"`
	Operation string `json:"op,omitempty"`
}
```

Create the subscription function:

```go
// Subscribe sends a websocket message to receive data from the channel
func (f *FTX) Subscribe(channelToSubscribe wshandler.WebsocketChannelSubscription) error {
	var sub WsSub
	a, err := f.GetPairAssetType(channelToSubscribe.Currency)
	if err != nil {
		return err
	}
	switch channelToSubscribe.Channel {
	case wsFills, wsOrders:
		sub.Operation = "subscribe"
		sub.Channel = channelToSubscribe.Channel
	default:
		sub.Operation = "subscribe"
		sub.Channel = channelToSubscribe.Channel
		sub.Market = f.FormatExchangeCurrency(channelToSubscribe.Currency, a).String()
	}
	return f.WebsocketConn.SendJSONMessage(sub)
}
```

- Create an unsubscribe function if the exchange has the functionality:

```go
// Unsubscribe sends a websocket message to stop receiving data from the channel
func (f *FTX) Unsubscribe(channelToSubscribe wshandler.WebsocketChannelSubscription) error {
	var unSub WsSub
	a, err := f.GetPairAssetType(channelToSubscribe.Currency)
	if err != nil {
		return err
	}
	unSub.Operation = "unsubscribe"
	unSub.Channel = channelToSubscribe.Channel
	unSub.Market = f.FormatExchangeCurrency(channelToSubscribe.Currency, a).String()
	return f.WebsocketConn.SendJSONMessage(unSub)
}
```