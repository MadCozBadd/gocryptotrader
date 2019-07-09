package lbank

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/config"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
	"github.com/thrasher-corp/gocryptotrader/exchanges/ticker"
	log "github.com/thrasher-corp/gocryptotrader/logger"
)

// Lbank is the overarching type across this package
type Lbank struct {
	exchange.Base
	privateKey    *rsa.PrivateKey
	privKeyLoaded bool
}

const (
	lbankAPIURL     = "https://api.lbkex.com"
	lbankAPIVersion = "1"

	// Public endpoints

	lbankTicker            = "ticker.do"
	lbankCurrencyPairs     = "currencyPairs.do"
	lbankMarketDepths      = "depth.do"
	lbankTrades            = "trades.do"
	lbankKlines            = "kline.do"
	lbankUserInfo          = "user_info.do"
	lbankPlaceOrder        = "create_order.do"
	lbankCancelOrder       = "cancel_order.do"
	lbankQueryOrder        = "orders_info.do"
	lbankQueryHistoryOrder = "orders_info_history.do"
	lbankOpeningOrders     = "orders_info_no_deal.do"
	lbankWithdrawlRecords  = "withdraws.do"
	lbankPairInfo          = "accuracy.do"
	lbankUSD2CNYRate       = "usdToCny.do"
	lbankWithdrawConfig    = "withdrawConfigs.do"
	lbankWithdraw          = "withdraw.do"
	lbankRevokeWithdraw    = "withdrawCancel.do"

	// Authenticated endpoints

)

// SetDefaults sets the basic defaults for Lbank
func (l *Lbank) SetDefaults() {
	l.Name = "Lbank"
	l.Enabled = false
	l.Verbose = false
	l.RESTPollingDelay = 10
	l.RequestCurrencyPairFormat.Delimiter = "_"
	l.RequestCurrencyPairFormat.Uppercase = false
	l.ConfigCurrencyPairFormat.Delimiter = "_"
	l.ConfigCurrencyPairFormat.Uppercase = false
	l.AssetTypes = []string{ticker.Spot}
	l.SupportsAutoPairUpdating = false
	l.SupportsRESTTickerBatching = false
	l.Requester = request.New(l.Name,
		request.NewRateLimit(time.Second, 0),
		request.NewRateLimit(time.Second, 0),
		common.NewHTTPClientWithTimeout(exchange.DefaultHTTPTimeout))
	l.APIUrlDefault = lbankAPIURL
	l.APIUrl = l.APIUrlDefault
	l.WebsocketInit()
}

// Setup takes in the supplied exchange configuration details and sets params
func (l *Lbank) Setup(exch *config.ExchangeConfig) {
	if !exch.Enabled {
		l.SetEnabled(false)
	} else {
		l.Enabled = true
		l.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
		l.SetAPIKeys(exch.APIKey, exch.APISecret, "", false)
		l.SetHTTPClientTimeout(exch.HTTPTimeout)
		l.SetHTTPClientUserAgent(exch.HTTPUserAgent)
		l.RESTPollingDelay = exch.RESTPollingDelay
		l.Verbose = exch.Verbose
		l.Websocket.SetWsStatusAndConnection(exch.Websocket)
		l.BaseCurrencies = exch.BaseCurrencies
		l.AvailablePairs = exch.AvailablePairs
		l.EnabledPairs = exch.EnabledPairs
		err := l.SetCurrencyPairFormat()
		if err != nil {
			log.Fatal(err)
		}
		err = l.SetAssetTypes()
		if err != nil {
			log.Fatal(err)
		}
		err = l.SetAutoPairDefaults()
		if err != nil {
			log.Fatal(err)
		}
		err = l.SetAPIURL(exch)
		if err != nil {
			log.Fatal(err)
		}
		err = l.SetClientProxyAddress(exch.ProxyAddress)
		if err != nil {
			log.Fatal(err)
		}

	}
}

// GetTicker returns a ticker for the specified symbol
// symbol: eth_btc
func (l *Lbank) GetTicker(symbol string) (TickerResponse, error) {
	var t TickerResponse
	params := url.Values{}
	params.Set("symbol", symbol)
	path := fmt.Sprintf("%s/v%s/%s?%s", l.APIUrl, lbankAPIVersion, lbankTicker, params.Encode())
	return t, l.SendHTTPRequest(path, &t)
}

// GetCurrencyPairs returns a list of supported currency pairs by the exchange
func (l *Lbank) GetCurrencyPairs() ([]string, error) {
	path := fmt.Sprintf("%s/v%s/%s", l.APIUrl, lbankAPIVersion,
		lbankCurrencyPairs)

	var result []string
	return result, l.SendHTTPRequest(path, &result)
}

// GetMarketDepths returns arrays of asks, bids and timestamp
func (l *Lbank) GetMarketDepths(symbol string, size string, merge string) (MarketDepthResponse, error) {
	var m MarketDepthResponse
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("size", size)
	params.Set("merge", merge)
	path := fmt.Sprintf("%s/v%s/%s?%s", l.APIUrl, lbankAPIVersion, lbankMarketDepths, params.Encode())
	return m, l.SendHTTPRequest(path, &m)
}

// GetTrades returns an array of structs of available trades regarding a particular exchange
func (l *Lbank) GetTrades(symbol string, size string, time string) ([]TradeResponse, error) {
	var g []TradeResponse
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("size", size)
	params.Set("time", time)
	path := fmt.Sprintf("%s/v%s/%s?%s", l.APIUrl, lbankAPIVersion, lbankTrades, params.Encode())
	log.Println(path)
	return g, l.SendHTTPRequest(path, &g)
}

// GetKlines returns arrays of something
func (l *Lbank) GetKlines(symbol string, size string, klineType string, time string) ([]KlineResponse, error) {
	var klineTemp interface{}
	var k []KlineResponse
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("size", size)
	params.Set("type", klineType)
	params.Set("time", time)
	path := fmt.Sprintf("%s/v%s/%s?%s", l.APIUrl, lbankAPIVersion, lbankKlines, params.Encode())
	err := l.SendHTTPRequest(path, &klineTemp)
	if err != nil {
		return k, err
	}

	resp, ok := klineTemp.([]interface{})
	if !ok {
		return nil, errors.New("response recieved is invalid")
	}

	for i := range resp {
		resp2, ok := resp[i].([]interface{})
		if !ok {
			return nil, errors.New("response recieved is invalid")
		}
		var someResponse KlineResponse
		for x := range resp2 {
			switch x {
			case 0:
				someResponse.TimeStamp = int64(resp2[x].(float64))
			case 1:
				if val, ok := resp2[x].(int64); ok {
					someResponse.OpenPrice = float64(val)
				} else {
					someResponse.OpenPrice = resp2[x].(float64)
				}
			case 2:
				if val, ok := resp2[x].(int64); ok {
					someResponse.HigestPrice = float64(val)
				} else {
					someResponse.HigestPrice = resp2[x].(float64)
				}
			case 3:
				if val, ok := resp2[x].(int64); ok {
					someResponse.ClosePrice = float64(val)
				} else {
					someResponse.ClosePrice = resp2[x].(float64)
				}
			case 4:
				if val, ok := resp2[x].(int64); ok {
					someResponse.TradingVolume = float64(val)
				} else {
					someResponse.TradingVolume = resp2[x].(float64)
				}
			}
		}
		k = append(k, someResponse)
	}
	return k, nil
}

// ErrCapture helps with error info
type ErrCapture struct {
	Error  int64 `json:"error_code"`
	Result bool  `json:"result,string"`
}

// GetUserInfo gets users account info
func (l *Lbank) GetUserInfo() (InfoResponse, error) {
	var resp InfoResponse
	path := fmt.Sprintf("%s/v%s/%s?", l.APIUrl, lbankAPIVersion, lbankUserInfo)
	return resp, l.SendAuthHTTPRequest("POST", path, nil, &resp)
}

// CreateOrder creates an order
func (l *Lbank) CreateOrder(pair, side string, amount, price float64) (CreateOrderResponse, error) {
	var resp CreateOrderResponse
	if !strings.EqualFold(side, "buy") && !strings.EqualFold(side, "sell") {
		log.Println(side)
		log.Println(strings.EqualFold(side, "BUY"))
		return resp, errors.New("side type invalid can only be 'buy' or 'sell'")
	}
	if amount <= 0 {
		return resp, errors.New("amount can't be smaller than 0")
	}
	if price <= 0 {
		return resp, errors.New("price can't be smaller than 0")
	}
	params := url.Values{}

	params.Set("symbol", pair)
	params.Set("type", common.StringToLower(side))
	params.Set("price", strconv.FormatFloat(price, 'f', -1, 64))
	params.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	path := fmt.Sprintf("%s/v%s/%s?", l.APIUrl, lbankAPIVersion, lbankPlaceOrder)
	return resp, l.SendAuthHTTPRequest("POST", path, params, &resp)
}

// RemoveOrder cancels a given order
func (l *Lbank) RemoveOrder(pair, orderID string) (RemoveOrderResponse, error) {
	var resp RemoveOrderResponse
	params := url.Values{}
	params.Set("symbol", pair)
	params.Set("order_id", orderID)
	path := fmt.Sprintf("%s/v%s/%s", l.APIUrl, lbankAPIVersion, lbankCancelOrder)
	return resp, l.SendAuthHTTPRequest("POST", path, params, &resp)
}

// QueryOrder finds out information about orders
func (l *Lbank) QueryOrder(pair, orderIDs string) (QueryOrderResponse, error) {
	var resp QueryOrderResponse
	params := url.Values{}
	params.Set("symbol", pair)
	params.Set("order_id", orderIDs)
	path := fmt.Sprintf("%s/v%s/%s?", l.APIUrl, lbankAPIVersion, lbankQueryOrder)
	return resp, l.SendAuthHTTPRequest("POST", path, params, &resp)
}

// QueryOrderHistory finds order info in the past 2 days
func (l *Lbank) QueryOrderHistory(pair, pageNumber, pageLength string) (OrderHistoryResponse, error) {
	var resp OrderHistory
	params := url.Values{}
	params.Set("symbol", pair)
	params.Set("current_page", pageNumber)
	params.Set("page_length", pageLength)
	path := fmt.Sprintf("%s/v%s/%s?", l.APIUrl, lbankAPIVersion, lbankQueryHistoryOrder)
	log.Println(path)
	err := l.SendAuthHTTPRequest("POST", path, params, &resp)
	if err != nil {
		return OrderHistoryResponse{}, err
	}

	var rt OrderHistoryResponse
	rt.CurrentPage = resp.CurrentPage
	rt.ErrorCode = resp.ErrorCode
	rt.PageLength = resp.PageLength
	rt.Result = resp.Result
	rt.Total = resp.Total

	var orders []OrderResponse
	err = json.Unmarshal(resp.Orders, &orders)
	if err == nil {
		rt.Orders = orders
		return rt, nil
	}

	var order OrderResponse
	err = json.Unmarshal(resp.Orders, &order)
	if err == nil {
		rt.Orders = append(rt.Orders, order)
		return rt, nil
	}

	return rt, nil
}

// GetPairInfo finds information about all trading pairs
func (l *Lbank) GetPairInfo() ([]PairInfoResponse, error) {
	var resp []PairInfoResponse
	path := fmt.Sprintf("%s/v%s/%s?", lbankAPIURL, lbankAPIVersion, lbankPairInfo)
	return resp, l.SendHTTPRequest(path, &resp)
}

// GetOpenOrders gets opening orders
func (l *Lbank) GetOpenOrders(pair string, pageNumber, pageLength int64) (OpenOrderResponse, error) {
	var resp OpenOrderResponse
	params := url.Values{}
	params.Set("symbol", pair)
	params.Set("current_page", strconv.FormatInt(pageNumber, 10))
	params.Set("page_length", strconv.FormatInt(pageLength, 10))
	path := fmt.Sprintf("%s/v%s/%s", l.APIUrl, lbankAPIVersion, lbankOpeningOrders)
	return resp, l.SendAuthHTTPRequest("POST", path, params, &resp)
}

// USD2RMBRate finds USD-CNY Rate
func (l *Lbank) USD2RMBRate() (ExchangeRateResponse, error) {
	var resp ExchangeRateResponse
	path := fmt.Sprintf("%s/v%s/%s", lbankAPIURL, lbankAPIVersion, lbankUSD2CNYRate)
	return resp, l.SendHTTPRequest(path, &resp)
}

// GetWithdrawConfig gets information about withdrawls
func (l *Lbank) GetWithdrawConfig(assetCode string) ([]WithdrawConfigResponse, error) {
	l.Verbose = true
	var resp []WithdrawConfigResponse
	params := url.Values{}
	if assetCode != "" {
		params.Set("assetCode", assetCode)
	}
	path := fmt.Sprintf("%s/v%s/%s?%s", lbankAPIURL, lbankAPIVersion, lbankWithdrawConfig, params.Encode())
	return resp, l.SendHTTPRequest(path, &resp)
}

// Withdraw withdraws
func (l *Lbank) Withdraw(account, assetCode, amount, memo, mark string) (WithdrawResponse, error) {
	var resp WithdrawResponse
	params := url.Values{}
	params.Set("account", account)
	params.Set("assetCode", assetCode)
	params.Set("amount", amount)
	if memo != "" {
		params.Set("memo", memo)
	}
	if mark != "" {
		params.Set("mark", mark)
	}
	path := fmt.Sprintf("%s/v%s/%s", lbankAPIURL, lbankAPIVersion, lbankWithdraw)
	return resp, l.SendAuthHTTPRequest("POST", path, params, &resp)
}

// RevokeWithdraw cancels the withdrawal given the withdrawalID
func (l *Lbank) RevokeWithdraw(withdrawID string) (RevokeWithdrawResponse, error) {
	var resp RevokeWithdrawResponse
	params := url.Values{}
	if withdrawID != "" {
		params.Set("withdrawId", withdrawID)
	}
	path := fmt.Sprintf("%s/v%s/%s?", lbankAPIURL, lbankAPIVersion, lbankRevokeWithdraw)
	return resp, l.SendAuthHTTPRequest("POST", path, params, &resp)
}

// GetWithdrawlRecords gets withdrawl records
func (l *Lbank) GetWithdrawlRecords(assetCode, status, pageNo, pageSize string) (WithdrawlResponse, error) {
	var resp WithdrawlResponse
	params := url.Values{}
	params.Set("assetCode", assetCode)
	params.Set("status", status)
	params.Set("pageNo", pageNo)
	params.Set("pageSize", pageSize)
	path := fmt.Sprintf("%s/v%s/%s", l.APIUrl, lbankAPIVersion, lbankWithdrawlRecords)
	return resp, l.SendAuthHTTPRequest("POST", path, params, &resp)
}

// ErrorCapture captures errors
func ErrorCapture(intermediary json.RawMessage) error {
	var capErr ErrCapture
	err := json.Unmarshal(intermediary, &capErr)
	if err == nil && capErr.Error != 0 {
		msg, ok := errorCodes[capErr.Error]
		if !ok {
			return errors.New("undefined code please check api docs for error code definition")
		}
		return errors.New(msg)
	}
	return nil
}

// SendHTTPRequest sends an unauthenticated HTTP request
func (l *Lbank) SendHTTPRequest(path string, result interface{}) error {
	var intermediary json.RawMessage
	err := l.SendPayload(http.MethodGet, path, nil, nil, &intermediary, false, false, l.Verbose, l.HTTPDebugging)
	if err != nil {
		return err
	}

	err = ErrorCapture(intermediary)
	if err != nil {
		return err
	}
	return json.Unmarshal(intermediary, result)
}

func (l *Lbank) loadPrivKey() error {
	if l.privKeyLoaded {
		return nil
	}

	key := strings.Join([]string{
		"-----BEGIN RSA PRIVATE KEY-----",
		l.APISecret,
		"-----END RSA PRIVATE KEY-----",
	}, "\n")

	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return fmt.Errorf("pem block is nil")
	}

	p, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("unable to decode priv key: %s", err)
	}

	var ok bool
	l.privateKey, ok = p.(*rsa.PrivateKey)
	if !ok {
		return fmt.Errorf("unable to parse RSA private key")
	}
	l.privKeyLoaded = true
	return nil
}

func (l *Lbank) sign(data string, p *rsa.PrivateKey) (string, error) {
	if p == nil {
		return "", errors.New("p cannot be nil")
	}
	md5hash := common.GetMD5([]byte(data))
	m := common.StringToUpper(common.HexEncodeToString(md5hash))
	s := common.GetSHA256([]byte(m))
	r, err := rsa.SignPKCS1v15(rand.Reader, p, crypto.SHA256, s)
	return common.Base64Encode(r), err
}

// SendAuthHTTPRequest sends an authenticated request
func (l *Lbank) SendAuthHTTPRequest(method, endpoint string, vals url.Values, result interface{}) error {
	headers := make(map[string]string)

	if vals == nil {
		vals = url.Values{}
	}

	if !l.privKeyLoaded {
		err := l.loadPrivKey()
		if err != nil {
			return err
		}
	}

	vals.Set("api_key", l.APIKey)
	sig, err := l.sign(vals.Encode(), l.privateKey)
	if err != nil {
		return err
	}

	vals.Set("sign", sig)
	payload := vals.Encode()
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	var intermediary json.RawMessage
	err = l.SendPayload(method, endpoint, headers, bytes.NewBufferString(payload), &intermediary, false, false, l.Verbose, l.HTTPDebugging)
	if err != nil {
		return err
	}

	err = ErrorCapture(intermediary)
	if err != nil {
		return err
	}
	return json.Unmarshal(intermediary, result)
}

var errorCodes = map[int64]string{
	10000: "Internal error",
	10001: "The required parameters can not be empty",
	10002: "Validation Failed",
	10003: "Invalid parameter",
	10004: "Request too frequent",
	10005: "Secret key does not exist",
	10006: "User does not exist",
	10007: "Invalid signature",
	10008: "Invalid Trading Pair",
	10009: "Price and/or Amount are required for limit order",
	10010: "Price and/or Amount must be more than 0",
	10013: "The amount is too small",
	10014: "Insufficient amount of money in account",
	10015: "Invalid order type",
	10016: "Insufficient account balance",
	10017: "Server Error",
	10018: "Page size should be between 1 and 50",
	10019: "Cancel NO more than 3 orders in one request",
	10020: "Volume < 0.001",
	10021: "Price < 0.01",
	10022: "Access denied",
	10023: "Market Order is not supported yet.",
	10024: "User cannot trade on this pair",
	10025: "Order has been filled",
	10026: "Order has been cancelld",
	10027: "Order is cancelling",
	10028: "Trading is paused",
	10100: "Has no privilege to withdraw",
	10101: "Invalid fee rate to withdraw",
	10102: "Too little to withdraw",
	10103: "Exceed daily limitation of withdraw",
	10104: "Cancel was rejected",
	10105: "Request has been cancelled",
}
