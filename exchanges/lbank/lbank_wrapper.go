package lbank

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/thrasher-/gocryptotrader/common"
	"github.com/thrasher-/gocryptotrader/currency"
	exchange "github.com/thrasher-/gocryptotrader/exchanges"
	"github.com/thrasher-/gocryptotrader/exchanges/orderbook"
	"github.com/thrasher-/gocryptotrader/exchanges/ticker"
	log "github.com/thrasher-/gocryptotrader/logger"
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
	exchangeCurrencies, err := l.GetCurrencyPairs()
	if err != nil {
		log.Errorf("%s Failed to get available symbols.\n", l.GetName())
	} else {
		forceUpdate := false
		if common.StringDataCompare(l.AvailablePairs.Strings(), "btc_usdt") {
			log.Warnf("%s contains invalid pair, forcing upgrade of available currencies.\n",
				l.GetName())
			forceUpdate = true
		}

		var newExchangeCurrencies currency.Pairs
		for _, p := range exchangeCurrencies {
			newExchangeCurrencies = append(newExchangeCurrencies,
				currency.NewPairFromString(p))
		}

		err = l.UpdateCurrencies(newExchangeCurrencies, false, forceUpdate)
		if err != nil {
			log.Errorf("%s Failed to update available currencies %s.\n", l.GetName(), err)
		}
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
	return nil, common.ErrFunctionNotSupported
}

// GetExchangeHistory returns historic trade data since exchange opening.
func (l *Lbank) GetExchangeHistory(p currency.Pair, assetType string) ([]exchange.TradeHistory, error) {
	return nil, common.ErrFunctionNotSupported
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
	return err
}

// CancelAllOrders cancels all orders associated with a currency pair
func (l *Lbank) CancelAllOrders(orders *exchange.OrderCancellation) (exchange.CancelAllOrdersResponse, error) {
	var resp exchange.CancelAllOrdersResponse
	orderIDs, err := l.GetAllOpenOrderID()
	if err != nil {
		return resp, nil
	}
	y := 1
	var tempSlice []string
	for i := range orderIDs {
		if orderIDs[i].CurrencyPair != orders.CurrencyPair.String() {
			continue
		}
		tempSlice = append(tempSlice, orderIDs[y].OrderID)
		if y%3 == 0 {
			input := strings.Join(tempSlice, ",")
			CancelResponse, err2 := l.RemoveOrder(orderIDs[y].CurrencyPair, input)
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
	return resp, nil
}

// GetOrderInfo returns information on a current open order
func (l *Lbank) GetOrderInfo(orderID string) (exchange.OrderDetail, error) {
	var resp exchange.OrderDetail
	orderIDs, err := l.GetAllOpenOrderID()
	if err != nil {
		return resp, err
	}

	for i := range orderIDs {
		if orderIDs[i].OrderID != orderID {
			continue
		}
		tempResp, err := l.QueryOrder(orderIDs[i].CurrencyPair, orderID)
		if err != nil {
			return resp, err
		}
		resp.Exchange = l.GetName()
		resp.CurrencyPair = currency.NewPairFromString(orderIDs[i].CurrencyPair)
		if strings.EqualFold(tempResp.Orders[0].Type, "buy") {
			resp.OrderSide = exchange.BuyOrderSide
		} else {
			resp.OrderSide = exchange.SellOrderSide
		}
		z := tempResp.Orders[0].Status
		switch {
		case z == -1:
			resp.Status = "cancelled"
		case z == 1:
			resp.Status = "on trading"
		case z == 2:
			resp.Status = "filled partially"
		case z == 3:
			resp.Status = "Filled totally"
		case z == 4:
			resp.Status = "Cancelling"
		default:
			return resp, fmt.Errorf("invalid order status: %v", tempResp.Orders[0].Status)
		}
		resp.Price = tempResp.Orders[0].Price
		resp.Amount = tempResp.Orders[0].Amount
		resp.ExecutedAmount = tempResp.Orders[0].DealAmount
		resp.RemainingAmount = tempResp.Orders[0].Price - tempResp.Orders[0].DealAmount
		resp.Fee, err = l.GetFeeByType(&exchange.FeeBuilder{
			FeeType:       exchange.CryptocurrencyTradeFee,
			Amount:        tempResp.Orders[0].Amount,
			PurchasePrice: tempResp.Orders[0].Price})
		if err != nil {
			return resp, err
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
	resp, err := l.Withdraw(withdrawRequest.Address, withdrawRequest.Currency.String(), strconv.FormatFloat(withdrawRequest.Amount, 'f', -1, 64), "", withdrawRequest.Description)
	if err != nil {
		return resp.WithdrawID, err
	}
	return resp.WithdrawID, nil
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
	var finalResp []exchange.OrderDetail
	var resp exchange.OrderDetail
	tempData, err := l.GetAllOpenOrderID()
	if err != nil {
		return finalResp, err
	}

	for _, val := range tempData {
		tempResp, err := l.QueryOrder(val.CurrencyPair, val.OrderID)
		if err != nil {
			return finalResp, err
		}
		resp.Exchange = l.GetName()
		resp.CurrencyPair = currency.NewPairFromString(val.CurrencyPair)
		if strings.EqualFold(tempResp.Orders[0].Type, "buy") {
			resp.OrderSide = exchange.BuyOrderSide
		} else {
			resp.OrderSide = exchange.SellOrderSide
		}
		z := tempResp.Orders[0].Status
		switch {
		case z == -1:
			resp.Status = "cancelled"
		case z == 1:
			resp.Status = "on trading"
		case z == 2:
			resp.Status = "filled partially"
		case z == 3:
			resp.Status = "Filled totally"
		case z == 4:
			resp.Status = "Cancelling"
		default:
			return finalResp, fmt.Errorf("invalid order status: %v", tempResp.Orders[0].Status)
		}
		resp.Price = tempResp.Orders[0].Price
		resp.Amount = tempResp.Orders[0].Amount
		resp.OrderDate = time.Unix(tempResp.Orders[0].CreateTime, 9)
		resp.ExecutedAmount = tempResp.Orders[0].DealAmount
		resp.RemainingAmount = tempResp.Orders[0].Price - tempResp.Orders[0].DealAmount
		resp.Fee, err = l.GetFeeByType(&exchange.FeeBuilder{
			FeeType:       exchange.CryptocurrencyTradeFee,
			Amount:        tempResp.Orders[0].Amount,
			PurchasePrice: tempResp.Orders[0].Price})
		if err != nil {
			return finalResp, err
		}
		for y := int(0); y < len(getOrdersRequest.Currencies); y++ {
			if getOrdersRequest.Currencies[y].String() != val.CurrencyPair {
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
				return resp, err
			}
		}
	}
	return resp, nil
}

// GetFeeByType returns an estimate of fee based on the type of transaction *
func (l *Lbank) GetFeeByType(feeBuilder *exchange.FeeBuilder) (float64, error) {
	var resp float64
	if feeBuilder.FeeType == exchange.CryptocurrencyTradeFee {
		return feeBuilder.Amount * feeBuilder.PurchasePrice * l.Fee, nil
	}
	if feeBuilder.FeeType == exchange.CryptocurrencyWithdrawalFee {
		withdrawalFee, err := l.GetWithdrawConfig(feeBuilder.Pair.Base.Lower().String())
		if err != nil {
			return resp, err
		}

		return withdrawalFee.Fee, nil
	}
	return resp, nil
}

// GetAllOpenOrderID returns map[string][]string -> map[currencypair][]orderIDs
func (l *Lbank) GetAllOpenOrderID() ([]GetAllOpenIDResp, error) {
	allPairs := l.GetEnabledCurrencies()
	var resp []GetAllOpenIDResp

	for i := range allPairs {
		pair := exchange.FormatExchangeCurrency(l.Name, allPairs[i])
		b := int64(1)
		tempResp, err := l.GetOpenOrders(pair.String(), b, 200)
		if err != nil {
			return resp, err
		}
		var x int64
		tempData, err := strconv.ParseInt(tempResp.Total, 10, 64)
		if err != nil {
			return resp, err
		}
		if tempData%200 != 0 {
			tempData = tempData - (tempData % 200)
			x = tempData/200 + 1
		} else {
			x = tempData / 200
		}
		for ; b <= x; b++ {
			tempResp, err = l.GetOpenOrders(pair.String(), b, 200)
			if err != nil {
				return resp, err
			}

			d, err := strconv.ParseInt(tempResp.Total, 10, 64)
			if err != nil {
				return resp, err
			}

			for c := int64(0); c < d; c++ {
				resp = append(resp, GetAllOpenIDResp{
					CurrencyPair: pair.String(),
					OrderID:      tempResp.Orders[c].OrderID})
			}
		}
	}
	return resp, nil
}

// SubscribeToWebsocketChannels appends to ChannelsToSubscribe
// which lets websocket.manageSubscriptions handle subscribing
func (l *Lbank) SubscribeToWebsocketChannels(channels []exchange.WebsocketChannelSubscription) error {
	return common.ErrNotYetImplemented
}

// UnsubscribeToWebsocketChannels removes from ChannelsToSubscribe
// which lets websocket.manageSubscriptions handle unsubscribing
func (l *Lbank) UnsubscribeToWebsocketChannels(channels []exchange.WebsocketChannelSubscription) error {
	return common.ErrNotYetImplemented
}

// AuthenticateWebsocket authenticates it
func (l *Lbank) AuthenticateWebsocket() error {
	return common.ErrNotYetImplemented
}

// GetSubscriptions gets subscriptions
func (l *Lbank) GetSubscriptions() ([]exchange.WebsocketChannelSubscription, error) {
	return nil, common.ErrNotYetImplemented
}
