package btcmarkets

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/crypto"
	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/orderbook"
	"github.com/thrasher-corp/gocryptotrader/exchanges/websocket/wshandler"
	log "github.com/thrasher-corp/gocryptotrader/logger"
)

const (
	btcMarketsWSURL = "wss://socket.btcmarkets.net/v2"
)

// WsConnect connects to a websocket feed
func (b *BTCMarkets) WsConnect() error {
	if !b.Websocket.IsEnabled() || !b.IsEnabled() {
		return errors.New(wshandler.WebsocketNotEnabled)
	}
	var dialer websocket.Dialer
	err := b.WebsocketConn.Dial(&dialer, http.Header{})
	if err != nil {
		return err
	}
	if b.Verbose {
		log.Debugf(log.ExchangeSys, "%s Connected to Websocket.\n", b.Name)
	}
	go b.WsHandleData()
	if b.GetAuthenticatedAPISupport(exchange.WebsocketAuthentication) {
		err = b.generateAuthSubscriptions()
		if err != nil {
			b.Websocket.DataHandler <- err
			b.Websocket.SetCanUseAuthenticatedEndpoints(false)
		}
	}
	return nil
}

// WsHandleData handles websocket data from WsReadData
func (b *BTCMarkets) WsHandleData() {
	b.Websocket.Wg.Add(1)
	defer func() {
		b.Websocket.Wg.Done()
	}()

	for {
		select {
		case <-b.Websocket.ShutdownC:
			return
		default:
			resp, err := b.WebsocketConn.ReadMessage()
			if err != nil {
				b.Websocket.ReadMessageErrors <- err
				return
			}
			b.Websocket.TrafficAlert <- struct{}{}
			var wsResponse WsMessageType
			err = common.JSONDecode(resp.Raw, &wsResponse)
			if err != nil {
				b.Websocket.DataHandler <- err
				continue
			}
			switch wsResponse.MessageType {
			case "heartbeat":
				if b.Verbose {
					log.Debugf(log.ExchangeSys, "%v - Websocket heartbeat received %s", b.Name, resp.Raw)
				}
			case "orderbook":
				var ob WsOrderbook
				err := common.JSONDecode(resp.Raw, &ob)
				if err != nil {
					b.Websocket.DataHandler <- err
					continue
				}

				p := currency.NewPairFromString(ob.Currency)
				var bids, asks []orderbook.Item
				for x := range ob.Bids {
					var price, amount float64
					price, err = strconv.ParseFloat(ob.Bids[x][0], 64)
					if err != nil {
						b.Websocket.DataHandler <- err
						continue
					}
					amount, err = strconv.ParseFloat(ob.Bids[x][1], 64)
					if err != nil {
						b.Websocket.DataHandler <- err
						continue
					}
					bids = append(bids, orderbook.Item{
						Amount: amount,
						Price:  price,
					})
				}
				for x := range ob.Asks {
					var price, amount float64
					price, err = strconv.ParseFloat(ob.Asks[x][0], 64)
					if err != nil {
						b.Websocket.DataHandler <- err
						continue
					}
					amount, err = strconv.ParseFloat(ob.Asks[x][1], 64)
					if err != nil {
						b.Websocket.DataHandler <- err
						continue
					}
					asks = append(asks, orderbook.Item{
						Amount: amount,
						Price:  price,
					})
				}
				err = b.Websocket.Orderbook.LoadSnapshot(&orderbook.Base{
					Pair:         p,
					Bids:         bids,
					Asks:         asks,
					LastUpdated:  ob.Timestamp,
					AssetType:    asset.Spot,
					ExchangeName: b.Name,
				})
				if err != nil {
					b.Websocket.DataHandler <- err
					continue
				}
				b.Websocket.DataHandler <- wshandler.WebsocketOrderbookUpdate{
					Pair:     p,
					Asset:    asset.Spot,
					Exchange: b.Name,
				}
			case "trade":
				var trade WsTrade
				err := common.JSONDecode(resp.Raw, &trade)
				if err != nil {
					b.Websocket.DataHandler <- err
					continue
				}
				p := currency.NewPairFromString(trade.Currency)
				b.Websocket.DataHandler <- wshandler.TradeData{
					Timestamp:    trade.Timestamp,
					CurrencyPair: p,
					AssetType:    asset.Spot,
					Exchange:     b.Name,
					Price:        trade.Price,
					Amount:       trade.Volume,
				}
			case "tick":
				var tick WsTick
				err := common.JSONDecode(resp.Raw, &tick)
				if err != nil {
					b.Websocket.DataHandler <- err
					continue
				}

				p := currency.NewPairFromString(tick.Currency)
				b.Websocket.DataHandler <- wshandler.TickerData{
					Exchange:  b.Name,
					Volume:    tick.Volume,
					High:      tick.High24,
					Low:       tick.Low24h,
					Bid:       tick.Bid,
					Ask:       tick.Ask,
					Last:      tick.Last,
					Timestamp: tick.Timestamp,
					AssetType: asset.Spot,
					Pair:      p,
				}
			case "fundChange":
				var transferData WsFundTransfer
				err := common.JSONDecode(resp.Raw, &transferData)
				if err != nil {
					b.Websocket.DataHandler <- err
					continue
				}
				b.Websocket.DataHandler <- transferData
			case "orderChange":
				var orderData WsOrderChange
				err := common.JSONDecode(resp.Raw, &orderData)
				if err != nil {
					b.Websocket.DataHandler <- err
					continue
				}
				b.Websocket.DataHandler <- orderData
			case "error":
				var wsErr WsError
				err := common.JSONDecode(resp.Raw, &wsErr)
				if err != nil {
					b.Websocket.DataHandler <- err
					continue
				}
				b.Websocket.DataHandler <- fmt.Errorf("%v websocket error. Code: %v Message: %v", b.Name, wsErr.Code, wsErr.Message)
			default:
				b.Websocket.DataHandler <- fmt.Errorf("%v Unhandled websocket message %s", b.Name, resp.Raw)
			}
		}
	}
}

func (b *BTCMarkets) generateDefaultSubscriptions() {
	var channels = []string{"tick", "trade", "orderbook"}
	enabledCurrencies := b.GetEnabledPairs(asset.Spot)
	var subscriptions []wshandler.WebsocketChannelSubscription
	for i := range channels {
		for j := range enabledCurrencies {
			subscriptions = append(subscriptions, wshandler.WebsocketChannelSubscription{
				Channel:  channels[i],
				Currency: enabledCurrencies[j],
			})
		}
	}
	b.Websocket.SubscribeToChannels(subscriptions)
}

// Subscribe sends a websocket message to receive data from the channel
func (b *BTCMarkets) Subscribe(channelToSubscribe wshandler.WebsocketChannelSubscription) error {
	req := WsSubscribe{
		MarketIDs:   []string{channelToSubscribe.Currency.String()},
		Channels:    []string{channelToSubscribe.Channel},
		MessageType: "subscribe",
	}
	return b.WebsocketConn.SendMessage(req)
}

// Login logs in allowing private ws events
func (b *BTCMarkets) generateAuthSubscriptions() error {
	var authSub WsAuthSubscribe
	signTime := strconv.FormatInt(time.Now().UTC().UnixNano()/1000000, 10)
	strToSign := "/users/self/subscribe" + "\n" + signTime
	tempSign := crypto.GetHMAC(crypto.HashSHA512,
		[]byte(strToSign),
		[]byte(b.API.Credentials.Secret))
	sign := crypto.Base64Encode(tempSign)
	authSub.MarketIDs = []string{"BTC-AUD"}
	authSub.Channels = []string{"fundChange", "heartbeat", "orderChange"}
	authSub.Key = b.API.Credentials.Key
	authSub.Signature = sign
	authSub.Timestamp = signTime
	authSub.MessageType = "subscribe"
	return b.WebsocketConn.SendMessage(authSub)
}