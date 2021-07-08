package deribit

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/thrasher-corp/gocryptotrader/config"
)

// Please supply your own keys here to do authenticated endpoint testing
const (
	apiKey                  = ""
	apiSecret               = ""
	canManipulateRealOrders = false
)

var d Deribit

func TestMain(m *testing.M) {
	d.SetDefaults()
	cfg := config.GetConfig()
	err := cfg.LoadConfig("../../testdata/configtest.json", true)
	if err != nil {
		log.Fatal(err)
	}

	exchCfg, err := cfg.GetExchangeConfig("Deribit")
	if err != nil {
		log.Fatal(err)
	}

	exchCfg.API.AuthenticatedSupport = true
	exchCfg.API.AuthenticatedWebsocketSupport = true
	exchCfg.API.Credentials.Key = apiKey
	exchCfg.API.Credentials.Secret = apiSecret
	d.API.Credentials.Key = apiKey
	d.API.Credentials.Secret = apiSecret

	err = d.Setup(exchCfg)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func areTestAPIKeysSet() bool {
	return d.ValidateAPICredentials()
}

// Implement tests for API endpoints below

func TestGetBookSummaryByCurrency(t *testing.T) {
	t.Parallel()
	_, err := d.GetBookSummaryByCurrency("BTC", "")
	if err != nil {
		t.Error(err)
	}
}

func TestGetBookSummaryByInstrument(t *testing.T) {
	t.Parallel()
	data, err := d.GetBookSummaryByInstrument("BTC-25MAR22")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)
}

func TestGetContractSize(t *testing.T) {
	t.Parallel()
	data, err := d.GetContractSize("BTC-25MAR22")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)
}

func TestGetCurrencies(t *testing.T) {
	t.Parallel()
	_, err := d.GetCurrencies()
	if err != nil {
		t.Error(err)
	}
}

func TestGetFundingChartData(t *testing.T) {
	t.Parallel()
	_, err := d.GetFundingChartData("BTC-PERPETUAL", "8h")
	if err != nil {
		t.Error(err)
	}
}

func TestGetFundingRateValue(t *testing.T) {
	t.Parallel()
	_, err := d.GetFundingRateValue("BTC-PERPETUAL", time.Now().Add(-time.Hour*8), time.Now())
	if err != nil {
		t.Error(err)
	}
	_, err = d.GetFundingRateValue("BTC-PERPETUAL", time.Now(), time.Now().Add(-time.Hour*8))
	if !errors.Is(err, errStartTimeCannotBeAfterEndTime) {
		t.Errorf("expected: %v, received %v", errStartTimeCannotBeAfterEndTime, err)
	}
}

func TestGetHistoricalVolatility(t *testing.T) {
	t.Parallel()
	_, err := d.GetHistoricalVolatility("BTC")
	if err != nil {
		t.Error(err)
	}
}

func TestGetIndexPrice(t *testing.T) {
	t.Parallel()
	_, err := d.GetIndexPrice("btc_usd")
	if err != nil {
		t.Error(err)
	}
}

func TestGetIndexPriceNames(t *testing.T) {
	t.Parallel()
	_, err := d.GetIndexPriceNames()
	if err != nil {
		t.Error(err)
	}
}

func TestGetInstrumentData(t *testing.T) {
	t.Parallel()
	_, err := d.GetInstrumentData("BTC-25MAR22")
	if err != nil {
		t.Error(err)
	}
}

func TestGetInstrumentsData(t *testing.T) {
	t.Parallel()
	_, err := d.GetInstrumentsData("BTC", "", false)
	if err != nil {
		t.Error(err)
	}
	_, err = d.GetInstrumentsData("BTC", "option", true)
	if err != nil {
		t.Error(err)
	}
}

func TestGetLastSettlementsByCurrency(t *testing.T) {
	t.Parallel()
	_, err := d.GetLastSettlementsByCurrency("BTC", "", "", 0, time.Time{})
	if err != nil {
		t.Error(err)
	}
	_, err = d.GetLastSettlementsByCurrency("BTC", "delivery", "5", 0, time.Now().Add(-time.Hour))
	if err != nil {
		t.Error(err)
	}
}

func TestGetLastSettlementsByInstrument(t *testing.T) {
	t.Parallel()
	_, err := d.GetLastSettlementsByInstrument("BTC-25MAR22", "", "", 0, time.Time{})
	if err != nil {
		t.Error(err)
	}
	_, err = d.GetLastSettlementsByInstrument("BTC-25MAR22", "settlement", "5", 0, time.Now().Add(-2*time.Hour))
	if err != nil {
		t.Error(err)
	}
}

func TestGetLastTradesByCurrency(t *testing.T) {
	t.Parallel()
	_, err := d.GetLastTradesByCurrency("BTC", "", "", "", "", 0, false)
	if err != nil {
		t.Error(err)
	}
	_, err = d.GetLastTradesByCurrency("BTC", "option", "36798", "36799", "asc", 0, true)
	if err != nil {
		t.Error(err)
	}
}

func TestGetLastTradesByCurrencyAndTime(t *testing.T) {
	t.Parallel()
	_, err := d.GetLastTradesByCurrencyAndTime("BTC", "", "", 0, false,
		time.Now().Add(-8*time.Hour), time.Now())
	if err != nil {
		t.Error(err)
	}
	_, err = d.GetLastTradesByCurrencyAndTime("BTC", "option", "asc", 25, false,
		time.Now().Add(-8*time.Hour), time.Now())
	if err != nil {
		t.Error(err)
	}
}

func TestGetLastTradesByInstrument(t *testing.T) {
	t.Parallel()
	_, err := d.GetLastTradesByInstrument("BTC-25MAR22", "", "", "", 0, false)
	if err != nil {
		t.Error(err)
	}
	_, err = d.GetLastTradesByInstrument("ETH-25MAR22", "30500", "31500", "desc", 0, true)
	if err != nil {
		t.Error(err)
	}
}

func TestGetLastTradesByInstrumentAndTime(t *testing.T) {
	t.Parallel()
	_, err := d.GetLastTradesByInstrumentAndTime("BTC-25MAR22", "", 0, false,
		time.Now().Add(-8*time.Hour), time.Now())
	if err != nil {
		t.Error(err)
	}
	_, err = d.GetLastTradesByInstrumentAndTime("BTC-25MAR22", "asc", 0, false,
		time.Now().Add(-8*time.Hour), time.Now())
	if err != nil {
		t.Error(err)
	}
}

func TestGetOrderbookData(t *testing.T) {
	t.Parallel()
	_, err := d.GetOrderbookData("BTC-25MAR22", 0)
	if err != nil {
		t.Error(err)
	}
}

func TestGetTradeVolumes(t *testing.T) {
	t.Parallel()
	_, err := d.GetTradeVolumes(false)
	if err != nil {
		t.Error(err)
	}
}

func TestGetTradingViewChartData(t *testing.T) {
	t.Parallel()
	_, err := d.GetTradingViewChartData("BTC-25MAR22", "60", time.Now().Add(-time.Hour), time.Now())
	if err != nil {
		t.Error(err)
	}
}

func TestGetVolatilityIndexData(t *testing.T) {
	t.Parallel()
	_, err := d.GetVolatilityIndexData("BTC", "60", time.Now().Add(-time.Hour), time.Now())
	if err != nil {
		t.Error(err)
	}
}

func TestGetPublicTicker(t *testing.T) {
	t.Parallel()
	_, err := d.GetPublicTicker("BTC-25MAR22")
	if err != nil {
		t.Error(err)
	}
}

func TestGetAccountSummary(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.GetAccountSummary("BTC", false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestCancelTransferByID(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.CancelTransferByID("BTC", "", 23487)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestGetTransfers(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.GetTransfers("BTC", 0, 0)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestCancelWithdrawal(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.CancelWithdrawal("BTC", 123844)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestCreateDepositAddress(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.CreateDepositAddress("BTC")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestGetCurrentDepositAddress(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.GetCurrentDepositAddress("BTC")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestGetDeposits(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.GetDeposits("BTC", 25, 0)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestGetWithdrawals(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.GetWithdrawals("BTC", 25, 0)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestSubmitTransferToSubAccount(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.SubmitTransferToSubAccount("BTC", 0.01, 13434)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestSubmitTransferToUser(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.SubmitTransferToUser("BTC", "", 0.001, 13434)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestSubmitWithdraw(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.SubmitWithdraw("BTC", "incorrectAddress", "", "", 0.001)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestChangeAPIKeyName(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.ChangeAPIKeyName(1, "TestKey123")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestChangeScopeInAPIKey(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.ChangeScopeInAPIKey(1, "account:read_write")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestChangeSubAccountName(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.ChangeSubAccountName(1, "TestingSubAccount")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestCreateAPIKey(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.CreateAPIKey("account:read_write", "TestingSubAccount", false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestCreateSubAccount(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.CreateSubAccount()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestDisableAPIKey(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.DisableAPIKey(1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestDisableTFAForSubAccount(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	// Use with caution will reduce the security of the account
	a, err := d.DisableTFAForSubAccount(1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestEnableAffiliateProgram(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.EnableAffiliateProgram()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestEnableAPIKey(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.EnableAPIKey(1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestGetAffiliateProgramInfo(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.GetAffiliateProgramInfo(1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestGetEmailLanguage(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.GetEmailLanguage()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestGetNewAnnouncements(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.GetNewAnnouncements()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestGetPosition(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.GetPosition("BTC-25MAR22")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestGetSubAccounts(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.GetSubAccounts(false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestGetPositions(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.GetPositions("BTC", "option")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
	_, err = d.GetPositions("ETH", "")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestGetTransactionLog(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.GetTransactionLog("BTC", "trade", time.Now().Add(-24*time.Hour), time.Now(), 5, 0)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
	a, err = d.GetTransactionLog("BTC", "trade", time.Now().Add(-24*time.Hour), time.Now(), 0, 0)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestListAPIKeys(t *testing.T) {
	d.Verbose = true
	t.Parallel()
	a, err := d.ListAPIKeys("")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}