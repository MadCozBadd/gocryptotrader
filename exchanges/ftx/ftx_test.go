package ftx

import (
	"log"
	"os"
	"testing"

	"github.com/thrasher-corp/gocryptotrader/config"
)

// Please supply your own keys here to do authenticated endpoint testing
const (
	apiKey                  = ""
	apiSecret               = ""
	canManipulateRealOrders = false
)

var f Ftx

func TestMain(m *testing.M) {
	f.SetDefaults()
	cfg := config.GetConfig()
	err := cfg.LoadConfig("../../testdata/configtest.json", true)
	if err != nil {
		log.Fatal(err)
	}

	exchCfg, err := cfg.GetExchangeConfig("Ftx")
	if err != nil {
		log.Fatal(err)
	}

	exchCfg.API.AuthenticatedSupport = true
	exchCfg.API.AuthenticatedWebsocketSupport = true
	exchCfg.API.Credentials.Key = apiKey
	exchCfg.API.Credentials.Secret = apiSecret

	err = f.Setup(exchCfg)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func areTestAPIKeysSet() bool {
	return f.ValidateAPICredentials()
}

// Implement tests for API endpoints below

func TestGetMarkets(t *testing.T) {
	a, err := f.GetMarkets()
	t.Log(a)
	if err != nil {
		t.Error(err)
	}
}

func TestGetMarket(t *testing.T) {
	a, err := f.GetMarket("FTT/BTC")
	t.Log(a)
	if err != nil {
		t.Error(err)
	}
}

func TestGetOrderbook(t *testing.T) {
	a, err := f.GetOrderbook("FTT/BTC", 5)
	t.Log(a)
	if err != nil {
		t.Error(err)
	}
}

func TestGetTrades(t *testing.T) {
	a, err := f.GetTrades("FTT/BTC", "10234032", "5234343433", 5)
	t.Log(a)
	if err != nil {
		t.Error(err)
	}
}

func TestGetHistoricalData(t *testing.T) {
	a, err := f.GetHistoricalData("FTT/BTC", "86400", "5", "", "")
	t.Log(a)
	if err != nil {
		t.Error(err)
	}
}

func TestGetFutures(t *testing.T) {
	a, err := f.GetFutures()
	t.Log(a)
	if err != nil {
		t.Error(err)
	}
}

func TestGetFuture(t *testing.T) {
	a, err := f.GetFuture("LEO-0327")
	t.Log(a)
	if err != nil {
		t.Error(err)
	}
}

func TestGetFutureStats(t *testing.T) {
	a, err := f.GetFutureStats("LEO-0327")
	t.Log(a)
	if err != nil {
		t.Error(err)
	}
}

func TestGetFundingRates(t *testing.T) {
	a, err := f.GetFundingRates()
	t.Log(a)
	if err != nil {
		t.Error(err)
	}
}

func TestSendAuthHTTPRequest(t *testing.T) {
	err := f.SendAuthHTTPRequest()
	if err != nil {
		t.Error(err)
	}
}
