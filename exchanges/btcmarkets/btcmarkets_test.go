package btcmarkets

import (
	"log"
	"os"
	"testing"

	"github.com/thrasher-corp/gocryptotrader/config"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
)

var b BTCMarkets

// Please supply your own keys here to do better tests
const (
	apiKey                  = ""
	apiSecret               = ""
	canManipulateRealOrders = false
)

func TestMain(m *testing.M) {
	b.SetDefaults()
	cfg := config.GetConfig()
	err := cfg.LoadConfig("../../testdata/configtest.json", true)
	if err != nil {
		log.Fatal("BTC Markets load config error", err)
	}
	bConfig, err := cfg.GetExchangeConfig("BTC Markets")
	if err != nil {
		log.Fatalf("BTC Markets Setup() init error: %v", err)
	}
	bConfig.API.Credentials.Key = apiKey
	bConfig.API.Credentials.Secret = apiSecret
	bConfig.API.AuthenticatedSupport = true

	err = b.Setup(bConfig)
	if err != nil {
		log.Fatal("BTC Markets setup error", err)
	}

	os.Exit(m.Run())
}

func areTestAPIKeysSet() bool {
	if b.API.Credentials.Key != "" && b.API.Credentials.Key != "Key" &&
		b.API.Credentials.Secret != "" && b.API.Credentials.Secret != "Secret" {
		return true
	}
	return false
}

func TestGetMarkets(t *testing.T) {
	t.Parallel()
	_, err := b.GetMarkets()
	if err != nil {
		t.Error("GetTicker() error", err)
	}
}

func TestGetTicker(t *testing.T) {
	t.Parallel()
	_, err := b.GetTicker("BTC-AUD")
	if err != nil {
		t.Error("GetOrderbook() error", err)
	}
}

func TestGetTrades(t *testing.T) {
	t.Parallel()
	_, err := b.GetTrades("BTC-AUD", 0, 0, 5)
	if err != nil {
		t.Error("GetTrades() error", err)
	}
}

func TestGetOrderbook(t *testing.T) {
	t.Parallel()
	_, err := b.GetOrderbook("BTC-AUD", 2)
	if err != nil {
		t.Error("GetTrades() error", err)
	}
}

func TestGetMarketCandles(t *testing.T) {
	t.Parallel()
	_, err := b.GetMarketCandles("BTC-AUD", "", "", "", 0, 0, 5)
	if err != nil {
		t.Error(err)
	}
}

func TestGetTickers(t *testing.T) {
	t.Parallel()
	_, err := b.GetTickers("BTC-AUD,LTC-AUD,ETH-AUD")
	if err != nil {
		t.Error(err)
	}
}

func TestGetMultipleOrderbooks(t *testing.T) {
	t.Parallel()
	_, err := b.GetMultipleOrderbooks("BTC-AUD,LTC-AUD,ETH-AUD")
	if err != nil {
		t.Error(err)
	}
}

func TestGetServerTime(t *testing.T) {
	t.Parallel()
	_, _, err := b.GetServerTime()
	if err != nil {
		t.Error(err)
	}
}

func TestGetAccountBalance(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.GetAccountBalance()
	if err != nil {
		t.Error(err)
	}
}

func TestGetTradingFees(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.GetTradingFees()
	if err != nil {
		t.Error(err)
	}
}

func TestGetTradeHistory(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.GetTradeHistory("BTC-AUD", "", "", "", "")
	if err != nil {
		t.Error(err)
	}
}

func TestNewOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test, either api keys or manipulaterealorders isnt set correctly")
	}
	_, err := b.NewOrder("BTC-AUD", 100, 1, "Limit", "Bid", 0, 0, "", true, "", "")
	if err != nil {
		t.Error(err)
	}
}

func TestGetOrders(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.GetOrders("", "", "", "", "")
	if err != nil {
		t.Error(err)
	}
}

func TestCancelOpenOrders(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test, either api keys or manipulaterealorders isnt set correctly")
	}
	_, err := b.CancelOpenOrders("BTC-AUD,LTC-AUD")
	if err != nil {
		t.Error(err)
	}
}

func TestFetchOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.FetchOrder("4477045999")
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveOrder(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test, either api keys or manipulaterealorders isnt set correctly")
	}
	_, err := b.RemoveOrder("4477045999")
	if err != nil {
		t.Error(err)
	}
}

func TestListWithdrawals(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.ListWithdrawals()
	if err != nil {
		t.Error(err)
	}
}

func TestGetWithdrawal(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.GetWithdrawal("4477381751")
	if err != nil {
		t.Error(err)
	}
}

func TestListDeposits(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.ListDeposits()
	if err != nil {
		t.Error(err)
	}
}

func TestGetDeposit(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.GetDeposit("4476769607")
	if err != nil {
		t.Error(err)
	}
}

func TestListTransfers(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.ListTransfers()
	if err != nil {
		t.Error(err)
	}
}

func TestGetTransfer(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.GetTransfer("4476769607")
	if err != nil {
		t.Error(err)
	}
}

func TestFetchDepositAddress(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.FetchDepositAddress("LTC")
	if err != nil {
		t.Error(err)
	}
}

func TestGetWithdrawalFees(t *testing.T) {
	t.Parallel()
	_, err := b.GetWithdrawalFees()
	if err != nil {
		t.Error(err)
	}
}

func TestListAssets(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.ListAssets()
	if err != nil {
		t.Error(err)
	}
}

func TestGetTransactions(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.GetTransactions("BTC")
	if err != nil {
		t.Error(err)
	}
}

func TestCreateNewReport(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.CreateNewReport("TransactionReport", "json")
	if err != nil {
		t.Error(err)
	}
}

func TestGetReport(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.GetReport("1kv38epne5v7lek9f18m60idg6")
	if err != nil {
		t.Error(err)
	}
}

func TestRequestWithdaw(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.RequestWithdraw("BTC", 1, "sdjflajdslfjld", "", "", "", "")
	if err != nil {
		t.Error(err)
	}
}

func TestBatchPlaceCancelOrders(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test, either api keys or manipulaterealorders isnt set correctly")
	}
	var pOrders []placeBatch
	placeOrders := placeBatch{MarketID: "BTC-AUD",
		Amount:    "11000",
		Price:     "1",
		OrderType: "Market",
		Side:      "Bid"}
	pOrders = append(pOrders, placeOrders)
	_, err := b.BatchPlaceCancelOrders(nil, pOrders)
	if err != nil {
		t.Error(err)
	}
}

func TestGetBatchTrades(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.GetBatchTrades("4477045999,4477381751,4476769607")
	if err != nil {
		t.Error(err)
	}
}

func TestCancelBatchOrders(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() || !canManipulateRealOrders {
		t.Skip("skipping test, either api keys or manipulaterealorders isnt set correctly")
	}
	_, err := b.CancelBatchOrders("4477045999,4477381751,4477381751")
	if err != nil {
		t.Error(err)
	}
}

func TestGetAccountInfo(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	_, err := b.GetAccountInfo()
	if err != nil {
		t.Error(err)
	}
}

func TestGetOrderHistory(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	var input order.GetOrdersRequest
	input.OrderSide = order.Buy
	_, err := b.GetOrderHistory(&input)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateOrderbook(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	cp := currency.NewPairWithDelimiter("BTC", "AUD", "-")
	_, err := b.UpdateOrderbook(cp, "spot")
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateTicker(t *testing.T) {
	t.Parallel()
	if !areTestAPIKeysSet() {
		t.Skip("API keys required but not set, skipping test")
	}
	cp := currency.NewPairWithDelimiter("BTC", "AUD", "-")
	_, err := b.UpdateTicker(cp, "spot")
	if err != nil {
		t.Error(err)
	}
}
