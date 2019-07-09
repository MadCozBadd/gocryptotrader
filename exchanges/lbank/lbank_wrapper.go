package lbank

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/orderbook"
	"github.com/thrasher-corp/gocryptotrader/exchanges/ticker"
	log "github.com/thrasher-corp/gocryptotrader/logger"
)

// Start starts the Lbank go routine
func (l *Lbank) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		l.Run()
		wg.Done()
	}()
}

// Run implements the Lbank wrapper
func (l *Lbank) Run() {
	if l.Verbose {
		log.Debugf("%s Websocket: %s. (url: %s).\n", l.GetName(), common.IsEnabled(l.Websocket.IsEnabled()), l.Websocket.GetWebsocketURL())
		log.Debugf("%s polling delay: %ds.\n", l.GetName(), l.RESTPollingDelay)
		log.Debugf("%s %d currencies enabled: %s.\n", l.GetName(), len(l.EnabledPairs), l.EnabledPairs)
	}
}

// UpdateTicker updates and returns the ticker for a currency pair
func (l *Lbank) UpdateTicker(p currency.Pair, assetType string) (ticker.Price, error) {
	var tickerPrice ticker.Price
	tickerInfo, err := l.GetTicker(p.String())
	if err != nil {
		return tickerPrice, err
	}
	tickerPrice.Pair = p
	tickerPrice.Last = tickerInfo.Ticker.Latest
	tickerPrice.High = tickerInfo.Ticker.High
	tickerPrice.Volume = tickerInfo.Ticker.Volume
	tickerPrice.Low = tickerInfo.Ticker.Low

	err = ticker.ProcessTicker(l.GetName(), &tickerPrice, assetType)
	if err != nil {
		return tickerPrice, err
	}

	return ticker.GetTicker(l.Name, p, assetType)
}

// GetTickerPrice returns the ticker for a currency pair
func (l *Lbank) GetTickerPrice(p currency.Pair, assetType string) (ticker.Price, error) {
	tickerNew, err := ticker.GetTicker(l.GetName(), p, assetType)
	if err != nil {
		return l.UpdateTicker(p, assetType)
	}
	return tickerNew, nil
}

// GetOrderbookEx returns orderbook base on the currency pair
func (l *Lbank) GetOrderbookEx(currency currency.Pair, assetType string) (orderbook.Base, error) {
	ob, err := orderbook.Get(l.GetName(), currency, assetType)
	if err != nil {
		return l.UpdateOrderbook(currency, assetType)
	}
	return ob, nil
}

// UpdateOrderbook updates and returns the orderbook for a currency pair
func (l *Lbank) UpdateOrderbook(p currency.Pair, assetType string) (orderbook.Base, error) {
	var orderBook orderbook.Base
	a, err := l.GetMarketDepths(p.String(), "60", "1")
	if err != nil {
		return orderBook, err
	}
	for i := range a.Asks {
		orderBook.Asks = append(orderBook.Asks, orderbook.Item{
			Price:  a.Asks[i][0],
			Amount: a.Asks[i][1]})
	}
	for i := range a.Bids {
		orderBook.Bids = append(orderBook.Bids, orderbook.Item{
			Price:  a.Bids[i][0],
			Amount: a.Bids[i][1]})
	}
	orderBook.Pair = p
	orderBook.ExchangeName = l.GetName()
	orderBook.AssetType = assetType
	err = orderBook.Process()
	if err != nil {
		return orderBook, err
	}

	return orderbook.Get(l.Name, p, assetType)
}

// GetAccountInfo retrieves balances for all enabled currencies for the
// Lbank exchange
func (l *Lbank) GetAccountInfo() (exchange.AccountInfo, error) {
	var info exchange.AccountInfo
	data, err := l.GetUserInfo()
	if err != nil {
		return info, err
	}

	var account exchange.Account
	for key, val := range data.Asset {
		c := currency.NewCode(key)
		hold, ok := data.Freeze[key]
		if !ok {
			return info, fmt.Errorf("hold data not found with %s", key)
		}
		account.Currencies = append(account.Currencies,
			exchange.AccountCurrencyInfo{CurrencyName: c,
				TotalValue: val,
				Hold:       hold})
	}

	info.Accounts = append(info.Accounts, account)
	info.Exchange = l.GetName()
	return info, nil
}

// GetFundingHistory returns funding history, deposits and
// withdrawals
func (l *Lbank) GetFundingHistory() ([]exchange.FundHistory, error) {
	return nil, common.ErrNotYetImplemented
}

// GetExchangeHistory returns historic trade data since exchange opening.
func (l *Lbank) GetExchangeHistory(p currency.Pair, assetType string) ([]exchange.TradeHistory, error) {
	return nil, common.ErrNotYetImplemented
}

// SubmitOrder submits a new order
func (l *Lbank) SubmitOrder(p currency.Pair, side exchange.OrderSide, _ exchange.OrderType, amount, price float64, clientID string) (exchange.SubmitOrderResponse, error) {
	var resp exchange.SubmitOrderResponse
	if side != "BUY" && side != "SELL" {
		return resp, fmt.Errorf("%s orderside is not supported by the exchange", side)
	}
	tempResp, err := l.CreateOrder(p.String(), side.ToString(), amount, price)
	if err != nil {
		return resp, err
	}
	resp.OrderID = tempResp.OrderID
	return resp, nil
}

// ModifyOrder will allow of changing orderbook placement and limit to
// market conversion
func (l *Lbank) ModifyOrder(action *exchange.ModifyOrder) (string, error) {
	return "", common.ErrFunctionNotSupported
}

// CancelOrder cancels an order by its corresponding ID number
func (l *Lbank) CancelOrder(order *exchange.OrderCancellation) error {
	_, err := l.RemoveOrder(order.CurrencyPair.Lower().String(), order.OrderID)
	if err != nil {
		return err
	}
	return nil
}

// CancelAllOrders cancels all orders associated with a currency pair
func (l *Lbank) CancelAllOrders(orders *exchange.OrderCancellation) (exchange.CancelAllOrdersResponse, error) {
	var resp exchange.CancelAllOrdersResponse
	mappymCMapMap, err := l.GetAllOpenOrderID()
	if err != nil {
		return resp, nil
	}
	for key := range mappymCMapMap {
		if key == orders.CurrencyPair.String() {
			var x int64
			x = 0
			for mappymCMapMap[key][x] != "" {
				x++
			}
			var y int64
			y = 0
			for y != x {
				var tempSlice []string
				tempSlice = append(tempSlice, mappymCMapMap[key][y])
				if y%3 == 0 {
					input := strings.Join(tempSlice, ",")
					CancelResponse, err2 := l.RemoveOrder(key, input)
					if err2 != nil {
						return resp, err2
					}
					tempStringSuccess := strings.Split(CancelResponse.Success, ",")
					for k := range tempStringSuccess {
						resp.OrderStatus[tempStringSuccess[k]] = "Cancelled"
					}
					tempStringError := strings.Split(CancelResponse.Error, ",")
					for l := range tempStringError {
						resp.OrderStatus[tempStringError[l]] = "Failed"
					}
					tempSlice = tempSlice[:0]
					y++
				}
				y++
			}
			x++
		}
	}

	// get all exchange trading pairs
	// var allOrders map[string][]string
	return resp, nil
}

// GetOrderInfo returns information on a current open order
func (l *Lbank) GetOrderInfo(orderID string) (exchange.OrderDetail, error) {
	var resp exchange.OrderDetail
	mappymCMapMap, err := l.GetAllOpenOrderID()
	if err != nil {
		return resp, err
	}

	for key, val := range mappymCMapMap {
		for i := range val {
			if val[i] == orderID {
				tempResp, err := l.QueryOrder(key, orderID)
				if err != nil {
					return resp, err
				}
				resp.Exchange = l.GetName()
				resp.CurrencyPair = currency.NewPairFromString(key)
				if strings.EqualFold(tempResp.Orders[0].Type, "buy") {
					resp.OrderSide = exchange.BuyOrderSide
				} else {
					resp.OrderSide = exchange.SellOrderSide
				}
				if tempResp.Orders[0].Status == -1 {
					resp.Status = "cancelled"
				}
				if tempResp.Orders[0].Status == 1 {
					resp.Status = "on trading"
				}
				if tempResp.Orders[0].Status == 2 {
					resp.Status = "filled partially"
				}
				if tempResp.Orders[0].Status == 3 {
					resp.Status = "Filled totally"
				}
				if tempResp.Orders[0].Status == 4 {
					resp.Status = "Cancelling"
				}
				resp.Price = tempResp.Orders[0].Price
				resp.Amount = tempResp.Orders[0].Amount
				resp.ExecutedAmount = tempResp.Orders[0].DealAmount
				resp.RemainingAmount = tempResp.Orders[0].Price - tempResp.Orders[0].DealAmount
				resp.Fee = 0.001
			}
		}
	}

	return resp, nil
}

// GetDepositAddress returns a deposit address for a specified currency
func (l *Lbank) GetDepositAddress(cryptocurrency currency.Code, accountID string) (string, error) {
	return "", common.ErrFunctionNotSupported
}

// WithdrawCryptocurrencyFunds returns a withdrawal ID when a withdrawal is
// submitted
func (l *Lbank) WithdrawCryptocurrencyFunds(withdrawRequest *exchange.WithdrawRequest) (string, error) {
	var resp string
	tempResp, err := l.Withdraw(withdrawRequest.Address, withdrawRequest.Currency.String(), strconv.FormatFloat(withdrawRequest.Amount, 'f', -1, 64), "", withdrawRequest.Description)
	if err != nil {
		return resp, err
	}
	resp = tempResp.WithdrawID
	return resp, nil
}

// WithdrawFiatFunds returns a withdrawal ID when a withdrawal is
// submitted
func (l *Lbank) WithdrawFiatFunds(withdrawRequest *exchange.WithdrawRequest) (string, error) {
	return "", common.ErrFunctionNotSupported
}

// WithdrawFiatFundsToInternationalBank returns a withdrawal ID when a withdrawal is
// submitted
func (l *Lbank) WithdrawFiatFundsToInternationalBank(withdrawRequest *exchange.WithdrawRequest) (string, error) {
	return "", common.ErrFunctionNotSupported
}

// GetWebsocket returns a pointer to the exchange websocket
func (l *Lbank) GetWebsocket() (*exchange.Websocket, error) {
	return nil, common.ErrNotYetImplemented
}

// GetActiveOrders retrieves any orders that are active/open
func (l *Lbank) GetActiveOrders(getOrdersRequest *exchange.GetOrdersRequest) ([]exchange.OrderDetail, error) {
	finalResp := make([]exchange.OrderDetail, 0)
	var resp exchange.OrderDetail
	tempData, err := l.GetAllOpenOrderID()
	if err != nil {
		return finalResp, err
	}

	for key, val := range tempData {
		for x := range val {
			tempResp, err := l.QueryOrder(key, val[x])
			if err != nil {
				return finalResp, err
			}
			resp.Exchange = l.GetName()
			resp.CurrencyPair = currency.NewPairFromString(key)
			if strings.EqualFold(tempResp.Orders[0].Type, "buy") {
				resp.OrderSide = exchange.BuyOrderSide
			} else {
				resp.OrderSide = exchange.SellOrderSide
			}
			if tempResp.Orders[0].Status == -1 {
				resp.Status = "cancelled"
			}
			if tempResp.Orders[0].Status == 1 {
				resp.Status = "on trading"
			}
			if tempResp.Orders[0].Status == 2 {
				resp.Status = "filled partially"
			}
			if tempResp.Orders[0].Status == 3 {
				resp.Status = "Filled totally"
			}
			if tempResp.Orders[0].Status == 4 {
				resp.Status = "Cancelling"
			}
			resp.Price = tempResp.Orders[0].Price
			resp.Amount = tempResp.Orders[0].Amount
			resp.OrderDate = time.Unix(tempResp.Orders[0].CreateTime, 9)
			resp.ExecutedAmount = tempResp.Orders[0].DealAmount
			resp.RemainingAmount = tempResp.Orders[0].Price - tempResp.Orders[0].DealAmount
			resp.Fee = 0.001
			for y := int(0); y < len(getOrdersRequest.Currencies); y++ {
				if getOrdersRequest.Currencies[y].String() != key {
					continue
				}
				if getOrdersRequest.OrderSide == "ANY" {
					finalResp = append(finalResp, resp)
					continue
				}
				if strings.EqualFold(getOrdersRequest.OrderSide.ToString(), tempResp.Orders[0].Type) {
					finalResp = append(finalResp, resp)
				}

			}

		}
	}
	return finalResp, nil
}

// GetOrderHistory retrieves account order information *
// Can Limit response to specific order status
func (l *Lbank) GetOrderHistory(getOrdersRequest *exchange.GetOrdersRequest) ([]exchange.OrderDetail, error) {
	var resp []exchange.OrderDetail
	for a := range getOrdersRequest.Currencies {
		p := exchange.FormatExchangeCurrency(l.Name, getOrdersRequest.Currencies[a])
		b := int64(1)
		tempResp, err := l.QueryOrderHistory(p.String(), strconv.FormatInt(b, 10), "200")
		if err != nil {
			return resp, err
		}
		tempData := tempResp.PageLength
		for tempData == 200 {
			tempResp, err = l.QueryOrderHistory(p.String(), strconv.FormatInt(b, 10), "200")
			if err != nil {

			}
		}
	}
	return resp, nil
}

// GetFeeByType returns an estimate of fee based on the type of transaction *
func (l *Lbank) GetFeeByType(feeBuilder *exchange.FeeBuilder) (float64, error) {
	resp := float64(0.001)
	return resp, nil
}

// GetAllOpenOrderID returns map[string][]string -> map[currencypair][]orderIDs
func (l *Lbank) GetAllOpenOrderID() (map[string][]string, error) {
	allPairs := l.GetEnabledCurrencies()
	resp := make(map[string][]string)

	for a := range allPairs {
		p := exchange.FormatExchangeCurrency(l.Name, allPairs[a])
		b := int64(1)
		tempResp, err := l.GetOpenOrders(p.String(), b, 200)
		if err != nil {
			return resp, err
		}
		tempData := tempResp.PageLength
		for tempData == 200 {
			tempResp, err = l.GetOpenOrders(p.String(), b, 200)
			if err != nil {
				return resp, err
			}

			for c := int64(0); c < tempData; c++ {
				resp[allPairs[a].String()] = append(resp[p.String()], tempResp.Orders[c].OrderID)
			}

			b++
		}
	}
	return resp, nil
}

// SubscribeToWebsocketChannels appends to ChannelsToSubscribe
// which lets websocket.manageSubscriptions handle subscribing
func (l *Lbank) SubscribeToWebsocketChannels(channels []exchange.WebsocketChannelSubscription) error {
	return common.ErrFunctionNotSupported
}

// UnsubscribeToWebsocketChannels removes from ChannelsToSubscribe
// which lets websocket.manageSubscriptions handle unsubscribing
func (l *Lbank) UnsubscribeToWebsocketChannels(channels []exchange.WebsocketChannelSubscription) error {
	return common.ErrFunctionNotSupported
}
