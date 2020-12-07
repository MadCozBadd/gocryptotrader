package huobi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/crypto"
	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
)

const (
	huobiAPIURL      = "https://api.huobi.pro"
	huobiURL         = "https://api.hbdm.com/"
	huobiAPIVersion  = "1"
	huobiAPIVersion2 = "2"

	// Futures endpoints
	fContractInfo              = "api/v1/contract_contract_info?"
	fContractIndexPrice        = "api/v1/contract_index?"
	fContractPriceLimitation   = "api/v1/contract_price_limit?"
	fContractOpenInterest      = "api/v1/contract_open_interest?"
	fEstimatedDeliveryPrice    = "api/v1/contract_delivery_price?"
	fContractMarketDepth       = "/market/depth?"
	fContractKline             = "/market/history/kline?"
	fMarketOverview            = "/market/detail/merged?"
	fLastTradeContract         = "/market/trade?"
	fContractBatchTradeRecords = "/market/history/trade?"
	fInsuranceAndClawback      = "api/v1/contract_risk_info?"
	fInsuranceBalanceHistory   = "api/v1/contract_insurance_fund?"
	fTieredAdjustmentFactor    = "api/v1/contract_adjustfactor?"
	fHisContractOpenInterest   = "api/v1/contract_his_open_interest?"
	fSystemStatus              = "api/v1/contract_api_state?"
	fTopAccountsSentiment      = "api/v1/contract_elite_account_ratio?"
	fTopPositionsSentiment     = "api/v1/contract_elite_position_ratio?"
	fLiquidationOrders         = "api/v1/contract_liquidation_orders?"
	fIndexKline                = "/index/market/history/index?"
	fBasisData                 = "/index/market/history/basis?"

	fAccountData               = "api/v1/contract_account_info"
	fPositionInformation       = "api/v1/contract_position_info"
	fAllSubAccountAssets       = "api/v1/contract_sub_account_list"
	fSingleSubAccountAssets    = "api/v1/contract_sub_account_info"
	fSingleSubAccountPositions = "api/v1/contract_sub_position_info"
	fFinancialRecords          = "api/v1/contract_financial_record"
	fSettlementRecords         = "api/v1/contract_user_settlement_records"
	fOrderLimitInfo            = "api/v1/contract_order_limit"
	fContractTradingFee        = "api/v1/contract_fee"
	fTransferLimitInfo         = "api/v1/contract_transfer_limit"
	fPositionLimitInfo         = "api/v1/contract_position_limit"
	fQueryAssetsAndPositions   = "api/v1/contract_account_position_info"
	fTransfer                  = "api/v1/contract_master_sub_transfer"
	fTransferRecords           = "api/v1/contract_master_sub_transfer_record"
	fAvailableLeverage         = "api/v1/contract_available_level_rate"
	fOrder                     = "api/v1/contract_order"
	fBatchOrder                = "api/v1/contract_batchorder"
	fCancelOrder               = "api/v1/contract_cancel"
	fCancelAllOrders           = "api/v1/contract_cancelall"
	fFlashCloseOrder           = "api/v1/lightning_close_position"
	fOrderInfo                 = "api/v1/contract_order_info"
	fOrderDetails              = "api/v1/contract_order_detail"
	fQueryOpenOrders           = "api/v1/contract_openorders"
	fOrderHistory              = "api/v1/contract_hisorders"
	fMatchResult               = "api/v1/contract_matchresults"
	fTriggerOrder              = "api/v1/contract_trigger_order"
	fCancelTriggerOrder        = "api/v1/contract_trigger_cancel"
	fCancelAllTriggerOrders    = "api/v1/contract_trigger_cancelall"
	fTriggerOpenOrders         = "api/v1/contract_trigger_openorders"
	fTriggerOrderHistory       = "api/v1/contract_trigger_hisorders"

	// Coin Margined Swap (perpetual futures) endpoints
	huobiSwapMarkets                     = "swap-api/v1/swap_contract_info?"
	huobiSwapFunding                     = "swap-api/v1/swap_funding_rate?"
	huobiSwapIndexPriceInfo              = "swap-api/v1/swap_index?"
	huobiSwapPriceLimitation             = "swap-api/v1/swap_price_limit?"
	huobiSwapOpenInterestInfo            = "swap-api/v1/swap_open_interest?"
	huobiSwapMarketDepth                 = "swap-ex/market/depth?"
	huobiKLineData                       = "swap-ex/market/history/kline?"
	huobiMarketDataOverview              = "swap-ex/market/detail/merged?"
	huobiLastTradeContract               = "swap-ex/market/trade?"
	huobiRequestBatchOfTradingRecords    = "swap-ex/market/history/trade?"
	huobiInsuranceBalanceAndClawbackRate = "swap-api/v1/swap_risk_info?"
	huobiInsuranceBalanceHistory         = "swap-api/v1/swap_insurance_fund?"
	huobiTieredAdjustmentFactor          = "swap-api/v1/swap_adjustfactor?"
	huobiOpenInterestInfo                = "swap-api/v1/swap_his_open_interest?"
	huobiSwapSystemStatus                = "swap-api/v1/swap_api_state?"
	huobiSwapSentimentAccountData        = "swap-api/v1/swap_elite_account_ratio?"
	huobiSwapSentimentPosition           = "swap-api/v1/swap_elite_position_ratio?"
	huobiSwapLiquidationOrders           = "swap-api/v1/swap_liquidation_orders?"
	huobiSwapHistoricalFundingRate       = "swap-api/v1/swap_historical_funding_rate?"
	huobiPremiumIndexKlineData           = "index/market/history/swap_premium_index_kline?"
	huobiPredictedFundingRateData        = "index/market/history/swap_estimated_rate_kline?"
	huobiBasisData                       = "index/market/history/swap_basis?"
	huobiSwapAccInfo                     = "swap-api/v1/swap_account_info"
	huobiSwapPosInfo                     = "swap-api/v1/swap_position_info"
	huobiSwapAssetsAndPos                = "swap-api/v1/swap_account_position_info" // nolint // false positive gosec
	huobiSwapSubAccList                  = "swap-api/v1/swap_sub_account_list"
	huobiSwapSubAccInfo                  = "swap-api/v1/swap_sub_account_info"
	huobiSwapSubAccPosInfo               = "swap-api/v1/swap_sub_position_info"
	huobiSwapFinancialRecords            = "swap-api/v1/swap_financial_record"
	huobiSwapSettlementRecords           = "swap-api/v1/swap_user_settlement_records"
	huobiSwapAvailableLeverage           = "swap-api/v1/swap_available_level_rate"
	huobiSwapOrderLimitInfo              = "swap-api/v1/swap_order_limit"
	huobiSwapTradingFeeInfo              = "swap-api/v1/swap_fee"
	huobiSwapTransferLimitInfo           = "swap-api/v1/swap_transfer_limit"
	huobiSwapPositionLimitInfo           = "swap-api/v1/swap_position_limit"
	huobiSwapInternalTransferData        = "swap-api/v1/swap_master_sub_transfer"
	huobiSwapInternalTransferRecords     = "swap-api/v1/swap_master_sub_transfer_record"
	huobiSwapPlaceOrder                  = "/swap-api/v1/swap_order"
	huobiSwapPlaceBatchOrder             = "/swap-api/v1/swap_batchorder"
	huobiSwapCancelOrder                 = "/swap-api/v1/swap_cancel"
	huobiSwapCancelAllOrders             = "/swap-api/v1/swap_cancelall"
	huobiSwapLightningCloseOrder         = "/swap-api/v1/swap_lightning_close_position"
	huobiSwapOrderInfo                   = "/swap-api/v1/swap_order_info"
	huobiSwapOrderDetails                = "/swap-api/v1/swap_order_detail"
	huobiSwapOpenOrders                  = "/swap-api/v1/swap_openorders"
	huobiSwapOrderHistory                = "/swap-api/v1/swap_hisorders"
	huobiSwapTradeHistory                = "/swap-api/v1/swap_matchresults"
	huobiSwapTriggerOrder                = "swap-api/v1/swap_trigger_order"
	huobiSwapCancelTriggerOrder          = "/swap-api/v1/swap_trigger_cancel"
	huobiSwapCancelAllTriggerOrders      = "/swap-api/v1/swap_trigger_cancelall"
	huobiSwapTriggerOrderHistory         = "/swap-api/v1/swap_trigger_hisorders"

	// Spot endpoints
	huobiMarketHistoryKline    = "market/history/kline"
	huobiMarketDetail          = "market/detail"
	huobiMarketDetailMerged    = "market/detail/merged"
	huobiMarketDepth           = "market/depth"
	huobiMarketTrade           = "market/trade"
	huobiMarketTickers         = "market/tickers"
	huobiMarketTradeHistory    = "market/history/trade"
	huobiSymbols               = "common/symbols"
	huobiCurrencies            = "common/currencys"
	huobiTimestamp             = "common/timestamp"
	huobiAccounts              = "account/accounts"
	huobiAccountBalance        = "account/accounts/%s/balance"
	huobiAccountDepositAddress = "account/deposit/address"
	huobiAccountWithdrawQuota  = "account/withdraw/quota"
	huobiAggregatedBalance     = "subuser/aggregate-balance"
	huobiOrderPlace            = "order/orders/place"
	huobiOrderCancel           = "order/orders/%s/submitcancel"
	huobiOrderCancelBatch      = "order/orders/batchcancel"
	huobiBatchCancelOpenOrders = "order/orders/batchCancelOpenOrders"
	huobiGetOrder              = "order/orders/getClientOrder"
	huobiGetOrderMatch         = "order/orders/%s/matchresults"
	huobiGetOrders             = "order/orders"
	huobiGetOpenOrders         = "order/openOrders"
	huobiGetOrdersMatch        = "orders/matchresults"
	huobiMarginTransferIn      = "dw/transfer-in/margin"
	huobiMarginTransferOut     = "dw/transfer-out/margin"
	huobiMarginOrders          = "margin/orders"
	huobiMarginRepay           = "margin/orders/%s/repay"
	huobiMarginLoanOrders      = "margin/loan-orders"
	huobiMarginAccountBalance  = "margin/accounts/balance"
	huobiWithdrawCreate        = "dw/withdraw/api/create"
	huobiWithdrawCancel        = "dw/withdraw-virtual/%s/cancel"
	huobiStatusError           = "error"
	huobiMarginRates           = "margin/loan-info"
)

// HUOBI is the overarching type across this package
type HUOBI struct {
	exchange.Base
	AccountID string
}

// Futures Contracts

// FGetContractInfo gets contract info for futures
func (h *HUOBI) FGetContractInfo(symbol, contractType, code string) (FContractInfoData, error) {
	var resp FContractInfoData
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if common.StringDataCompare(validContractTypes, contractType) {
		params.Set("contract_type", contractType)
	}
	if code != "" {
		params.Set("contract_code", code)
	}
	path := fContractInfo + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FIndexPriceInfo gets index price info for a futures contract
func (h *HUOBI) FIndexPriceInfo(symbol string) (FContractIndexPriceInfo, error) {
	var resp FContractIndexPriceInfo
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	path := fContractIndexPrice + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FContractPriceLimitations gets price limits for a futures contract
func (h *HUOBI) FContractPriceLimitations(symbol, contractType, code string) (FContractIndexPriceInfo, error) {
	var resp FContractIndexPriceInfo
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if common.StringDataCompare(validContractTypes, contractType) {
		params.Set("contract_type", contractType)
	}
	if code != "" {
		params.Set("contract_code", code)
	}
	path := fContractPriceLimitation + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FContractOpenInterest gets open interest data for futures contracts
func (h *HUOBI) FContractOpenInterest(symbol, contractType, code string) (FContractOIData, error) {
	var resp FContractOIData
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if common.StringDataCompare(validContractTypes, contractType) {
		params.Set("contract_type", contractType)
	}
	if code != "" {
		params.Set("contract_code", code)
	}
	path := fContractOpenInterest + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FGetEstimatedDeliveryPrice gets estimated delivery price info for futures
func (h *HUOBI) FGetEstimatedDeliveryPrice(symbol string) (FEstimatedDeliveryPriceInfo, error) {
	var resp FEstimatedDeliveryPriceInfo
	params := url.Values{}
	params.Set("symbol", symbol)
	path := fEstimatedDeliveryPrice + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FGetMarketDepth gets market depth data for futures contracts
func (h *HUOBI) FGetMarketDepth(symbol, dataType string) (OBData, error) {
	var resp OBData
	var tempData FMarketDepth
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("type", dataType)
	path := fContractMarketDepth + params.Encode()
	err := h.SendHTTPRequest(exchange.RestFutures, path, &tempData)
	if err != nil {
		return resp, err
	}
	resp.Symbol = symbol
	for x := range tempData.Tick.Asks {
		resp.Asks = append(resp.Asks, obItem{
			Price:    tempData.Tick.Asks[x][0],
			Quantity: tempData.Tick.Bids[x][1],
		})
	}
	for y := range tempData.Tick.Bids {
		resp.Bids = append(resp.Bids, obItem{
			Price:    tempData.Tick.Bids[y][0],
			Quantity: tempData.Tick.Bids[y][1],
		})
	}
	return resp, nil
}

// FGetKlineData gets kline data for futures
func (h *HUOBI) FGetKlineData(symbol, period string, size int64, startTime, endTime time.Time) (FKlineData, error) {
	var resp FKlineData
	params := url.Values{}
	params.Set("symbol", symbol)
	if !common.StringDataCompare(validFuturesPeriods, period) {
		return resp, fmt.Errorf("invalid period value received")
	}
	params.Set("period", period)
	if !(size > 1) && !(size < 2000) {
		return resp, fmt.Errorf("invalid size")
	}
	params.Set("size", strconv.FormatInt(size, 10))
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return resp, errors.New("startTime cannot be after endTime")
		}
		params.Set("start_time", strconv.FormatInt(startTime.Unix(), 10))
		params.Set("end_time", strconv.FormatInt(endTime.Unix(), 10))
	}
	path := fContractKline + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FGetMarketOverviewData gets market overview data for futures
func (h *HUOBI) FGetMarketOverviewData(symbol string) (FMarketOverviewData, error) {
	var resp FMarketOverviewData
	params := url.Values{}
	params.Set("symbol", symbol)
	path := fMarketOverview + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FLastTradeData gets last trade data for a futures contract
func (h *HUOBI) FLastTradeData(symbol string) (FLastTradeData, error) {
	var resp FLastTradeData
	params := url.Values{}
	params.Set("symbol", symbol)
	path := fLastTradeContract + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FRequestPublicBatchTrades gets public batch trades for a futures contract
func (h *HUOBI) FRequestPublicBatchTrades(symbol string, size int64) (FBatchTradesForContractData, error) {
	var resp FBatchTradesForContractData
	params := url.Values{}
	params.Set("symbol", symbol)
	if size > 1 && size < 2000 {
		params.Set("size", strconv.FormatInt(size, 10))
	}
	path := fContractBatchTradeRecords + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FQueryInsuranceAndClawbackData gets insurance and clawback data for a futures contract
func (h *HUOBI) FQueryInsuranceAndClawbackData(symbol string) (FClawbackRateAndInsuranceData, error) {
	var resp FClawbackRateAndInsuranceData
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	path := fInsuranceAndClawback + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FQueryHistoricalInsuranceData gets insurance data
func (h *HUOBI) FQueryHistoricalInsuranceData(symbol string) (FHistoricalInsuranceRecordsData, error) {
	var resp FHistoricalInsuranceRecordsData
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	path := fInsuranceBalanceHistory + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FQueryTieredAdjustmentFactor gets tiered adjustment factor for futures contracts
func (h *HUOBI) FQueryTieredAdjustmentFactor(symbol string) (FTieredAdjustmentFactorInfo, error) {
	var resp FTieredAdjustmentFactorInfo
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	path := fTieredAdjustmentFactor + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FQueryHisOpenInterest gets open interest for futures contract
func (h *HUOBI) FQueryHisOpenInterest(symbol, contractType, period, amountType string, size int64) (FOIData, error) {
	var resp FOIData
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if !common.StringDataCompare(validContractTypes, contractType) {
		return resp, fmt.Errorf("invalid contract type")
	}
	params.Set("contract_type", contractType)
	if !common.StringDataCompare(validPeriods, period) {
		return resp, fmt.Errorf("invalid period")
	}
	params.Set("period", period)
	if size > 0 || size <= 200 {
		params.Set("size", strconv.FormatInt(size, 10))
	}
	validAmount, ok := validAmountType[amountType]
	if !ok {
		return resp, fmt.Errorf("invalid amountType")
	}
	params.Set("amount_type", strconv.FormatInt(validAmount, 10))
	path := fHisContractOpenInterest + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FQuerySystemStatus gets system status data
func (h *HUOBI) FQuerySystemStatus(symbol string) (FContractOIData, error) {
	var resp FContractOIData
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	path := fSystemStatus + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FQueryTopAccountsRatio gets top accounts' ratio
func (h *HUOBI) FQueryTopAccountsRatio(symbol, period string) (FTopAccountsLongShortRatio, error) {
	var resp FTopAccountsLongShortRatio
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if !common.StringDataCompare(validPeriods, period) {
		return resp, fmt.Errorf("invalid period")
	}
	params.Set("period", period)
	path := fTopAccountsSentiment + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FQueryTopPositionsRatio gets top positions' long/short ratio for futures
func (h *HUOBI) FQueryTopPositionsRatio(symbol, period string) (FTopPositionsLongShortRatio, error) {
	var resp FTopPositionsLongShortRatio
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if !common.StringDataCompare(validPeriods, period) {
		return resp, fmt.Errorf("invalid period")
	}
	params.Set("period", period)
	path := fTopPositionsSentiment + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FLiquidationOrders gets liquidation orders for futures contracts
func (h *HUOBI) FLiquidationOrders(symbol, tradeType string, pageIndex, pageSize, createDate int64) (FLiquidationOrdersInfo, error) {
	var resp FLiquidationOrdersInfo
	params := url.Values{}
	params.Set("symbol", symbol)
	if createDate != 7 && createDate != 90 {
		return resp, fmt.Errorf("invalid createDate. 7 and 90 are the only supported values")
	}
	params.Set("create_date", strconv.FormatInt(createDate, 10))
	tType, ok := validTradeTypes[tradeType]
	if !ok {
		return resp, fmt.Errorf("invalid trade type")
	}
	params.Set("trade_type", strconv.FormatInt(tType, 10))
	if pageIndex != 0 {
		params.Set("page_index", strconv.FormatInt(pageIndex, 10))
	}
	if pageSize != 0 {
		params.Set("page_size", strconv.FormatInt(pageIndex, 10))
	}
	path := fLiquidationOrders + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FIndexKline gets index kline data for futures contracts
func (h *HUOBI) FIndexKline(symbol, period string, size int64) (FIndexKlineData, error) {
	var resp FIndexKlineData
	params := url.Values{}
	params.Set("symbol", symbol)
	if !common.StringDataCompare(validFuturesPeriods, period) {
		return resp, fmt.Errorf("invalid period value received")
	}
	params.Set("period", period)
	if size == 0 && size > 2000 {
		return resp, fmt.Errorf("invalid size")
	}
	params.Set("size", strconv.FormatInt(size, 10))
	path := fIndexKline + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FGetBasisData gets basis data futures contracts
func (h *HUOBI) FGetBasisData(symbol, period, basisPriceType string, size int64) (FBasisData, error) {
	var resp FBasisData
	params := url.Values{}
	params.Set("symbol", symbol)
	if !common.StringDataCompare(validFuturesPeriods, period) {
		return resp, fmt.Errorf("invalid period value received")
	}
	params.Set("period", period)
	if basisPriceType != "" {
		if common.StringDataCompare(validBasisPriceTypes, basisPriceType) {
			params.Set("basis_price_type", basisPriceType)
		}
	}
	if size > 0 && size <= 2000 {
		params.Set("size", strconv.FormatInt(size, 10))
	}
	path := fBasisData + params.Encode()
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// FGetAccountInfo gets user info for futures account
func (h *HUOBI) FGetAccountInfo(symbol string) (FUserAccountData, error) {
	var resp FUserAccountData
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fAccountData, nil, req, &resp)
}

// FGetPositionsInfo gets positions info for futures account
func (h *HUOBI) FGetPositionsInfo(symbol string) (FUserAccountData, error) {
	var resp FUserAccountData
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fPositionInformation, nil, req, &resp)
}

// FGetAllSubAccountAssets gets assets info for all futures subaccounts
func (h *HUOBI) FGetAllSubAccountAssets(symbol string) (FSubAccountAssetsInfo, error) {
	var resp FSubAccountAssetsInfo
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fAllSubAccountAssets, nil, req, &resp)
}

// FGetSingleSubAccountInfo gets assets info for a futures subaccount
func (h *HUOBI) FGetSingleSubAccountInfo(symbol, subUID string) (FSingleSubAccountAssetsInfo, error) {
	var resp FSingleSubAccountAssetsInfo
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	req["sub_uid"] = subUID
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fSingleSubAccountAssets, nil, req, &resp)
}

// FGetSingleSubPositions gets positions info for a single sub account
func (h *HUOBI) FGetSingleSubPositions(symbol, subUID string) (FSingleSubAccountPositionsInfo, error) {
	var resp FSingleSubAccountPositionsInfo
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	req["sub_uid"] = subUID
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fSingleSubAccountPositions, nil, req, &resp)
}

// FGetFinancialRecords gets financial records for futures
func (h *HUOBI) FGetFinancialRecords(symbol, recordType string, createDate, pageIndex, pageSize int64) (FFinancialRecords, error) {
	var resp FFinancialRecords
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	if recordType != "" {
		rType, ok := validFuturesRecordTypes[recordType]
		if !ok {
			return resp, fmt.Errorf("invalid recordType")
		}
		req["type"] = rType
	}
	if createDate > 0 && createDate < 90 {
		req["create_date"] = createDate
	}
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize != 0 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fFinancialRecords, nil, req, &resp)
}

// FGetSettlementRecords gets settlement records for futures
func (h *HUOBI) FGetSettlementRecords(symbol string, pageIndex, pageSize int64, startTime, endTime time.Time) (FSettlementRecords, error) {
	var resp FSettlementRecords
	req := make(map[string]interface{})
	req["symbol"] = symbol
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize != 0 {
		req["page_size"] = pageSize
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return resp, errors.New("startTime cannot be after endTime")
		}
		req["start_time"] = strconv.FormatInt(startTime.Unix(), 10)
		req["end_time"] = strconv.FormatInt(endTime.Unix(), 10)
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fSettlementRecords, nil, req, &resp)
}

// FGetOrderLimits gets order limits for futures contracts
func (h *HUOBI) FGetOrderLimits(symbol, orderPriceType string) (FContractInfoOnOrderLimit, error) {
	var resp FContractInfoOnOrderLimit
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	if orderPriceType != "" {
		if !common.StringDataCompare(validFuturesOrderPriceTypes, orderPriceType) {
			return resp, fmt.Errorf("invalid orderPriceType")
		}
		req["order_price_type"] = orderPriceType
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fOrderLimitInfo, nil, req, &resp)
}

// FContractTradingFee gets futures contract trading fees
func (h *HUOBI) FContractTradingFee(symbol string) (FContractTradingFeeData, error) {
	var resp FContractTradingFeeData
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fContractTradingFee, nil, req, &resp)
}

// FGetTransferLimits gets transfer limits for futures
func (h *HUOBI) FGetTransferLimits(symbol string) (FTransferLimitData, error) {
	var resp FTransferLimitData
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fTransferLimitInfo, nil, req, &resp)
}

// FGetPositionLimits gets position limits for futures
func (h *HUOBI) FGetPositionLimits(symbol string) (FPositionLimitData, error) {
	var resp FPositionLimitData
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fPositionLimitInfo, nil, req, &resp)
}

// FGetAssetsAndPositions gets assets and positions for futures
func (h *HUOBI) FGetAssetsAndPositions(symbol string) (FAssetsAndPositionsData, error) {
	var resp FAssetsAndPositionsData
	req := make(map[string]interface{})
	req["symbol"] = symbol
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fQueryAssetsAndPositions, nil, req, &resp)
}

// FTransfer transfers assets between master and subaccounts
func (h *HUOBI) FTransfer(subUID, symbol, transferType string, amount float64) (FAccountTransferData, error) {
	var resp FAccountTransferData
	req := make(map[string]interface{})
	req["symbol"] = symbol
	req["subUid"] = subUID
	req["amount"] = amount
	if !common.StringDataCompare(validTransferType, transferType) {
		return resp, fmt.Errorf("inavlid transferType received")
	}
	req["type"] = transferType
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fTransfer, nil, req, &resp)
}

// FGetTransferRecords gets transfer records data for futures
func (h *HUOBI) FGetTransferRecords(symbol, transferType string, createDate, pageIndex, pageSize int64) (FTransferRecords, error) {
	var resp FTransferRecords
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	if !common.StringDataCompare(validTransferType, transferType) {
		return resp, fmt.Errorf("inavlid transferType received")
	}
	req["type"] = transferType
	if createDate < 0 || createDate > 90 {
		return resp, fmt.Errorf("invalid create date value: only supports up to 90 days")
	}
	req["create_date"] = strconv.FormatInt(createDate, 10)
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize > 0 && pageSize <= 50 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fTransferRecords, nil, req, &resp)
}

// FGetAvailableLeverage gets available leverage data for futures
func (h *HUOBI) FGetAvailableLeverage(symbol string) (FAvailableLeverageData, error) {
	var resp FAvailableLeverageData
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fAvailableLeverage, nil, req, &resp)
}

// FOrder places an order for futures
func (h *HUOBI) FOrder(symbol, contractType, contractCode, clientOrderID, direction, offset, orderPriceType string, price, volume, leverageRate float64) (FOrderData, error) {
	var resp FOrderData
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	if contractType != "" {
		if !common.StringDataCompare(validContractTypes, contractType) {
			return resp, fmt.Errorf("invalid contractType")
		}
		req["contract_type"] = contractType
	}
	if contractCode != "" {
		req["contract_code"] = contractCode
	}
	if clientOrderID != "" {
		req["client_order_id"] = clientOrderID
	}
	req["direction"] = direction
	if !common.StringDataCompare(validOffsetTypes, offset) {
		return resp, fmt.Errorf("invalid offset amounts")
	}
	if !common.StringDataCompare(validFuturesOrderPriceTypes, orderPriceType) {
		return resp, fmt.Errorf("invalid orderPriceType")
	}
	req["order_price_type"] = orderPriceType
	req["lever_rate"] = leverageRate
	req["volume"] = volume
	req["price"] = price
	req["offset"] = offset
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fOrder, nil, req, &resp)
}

// FPlaceBatchOrder places a batch of orders for futures
func (h *HUOBI) FPlaceBatchOrder(data []fBatchOrderData) (FBatchOrderResponse, error) {
	var resp FBatchOrderResponse
	req := make(map[string]interface{})
	if (len(data) > 10) || (len(data) == 0) {
		return resp, fmt.Errorf("invalid data provided: maximum of 10 batch orders supported")
	}
	for x := range data {
		if data[x].ContractType != "" {
			if !common.StringDataCompare(validContractTypes, data[x].ContractType) {
				return resp, fmt.Errorf("invalid contractType")
			}
		}
		if !common.StringDataCompare(validOffsetTypes, data[x].Offset) {
			return resp, fmt.Errorf("invalid offset amounts")
		}
		if !common.StringDataCompare(validFuturesOrderPriceTypes, data[x].OrderPriceType) {
			return resp, fmt.Errorf("invalid orderPriceType")
		}
	}
	req["orders_data"] = data
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fBatchOrder, nil, req, &resp)
}

// FCancelOrder cancels a futures order
func (h *HUOBI) FCancelOrder(symbol, orderID, clientOrderID string) (FCancelOrderData, error) {
	var resp FCancelOrderData
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	if orderID != "" {
		req["order_id"] = orderID
	}
	if clientOrderID != "" {
		req["client_order_id"] = clientOrderID
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fCancelOrder, nil, req, &resp)
}

// FCancelAllOrders cancels all futures order for a given symbol
func (h *HUOBI) FCancelAllOrders(symbol, contractCode, contractType string) (FCancelOrderData, error) {
	var resp FCancelOrderData
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	if contractType != "" {
		if !common.StringDataCompare(validContractTypes, contractType) {
			return resp, fmt.Errorf("invalid contractType")
		}
		req["contract_type"] = contractType
	}
	if contractCode != "" {
		req["contract_code"] = contractCode
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fCancelAllOrders, nil, req, &resp)
}

// FFlashCloseOrder flash closes a futures order
func (h *HUOBI) FFlashCloseOrder(symbol, contractType, contractCode, direction, orderPriceType, clientOrderID string, volume float64) (FOrderData, error) {
	var resp FOrderData
	req := make(map[string]interface{})
	req["symbol"] = symbol
	if contractType != "" {
		if !common.StringDataCompare(validContractTypes, contractType) {
			return resp, fmt.Errorf("invalid contractType")
		}
		req["contract_type"] = contractType
	}
	if contractCode != "" {
		req["contract_code"] = contractCode
	}
	req["direction"] = direction
	req["volume"] = volume
	if clientOrderID != "" {
		req["client_order_id"] = clientOrderID
	}
	if orderPriceType != "" {
		if !common.StringDataCompare(validOPTypes, orderPriceType) {
			return resp, fmt.Errorf("invalid orderPriceType")
		}
		req["orderPriceType"] = orderPriceType
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fFlashCloseOrder, nil, req, &resp)
}

// FGetOrderInfo gets order info for futures
func (h *HUOBI) FGetOrderInfo(symbol, clientOrderID, orderID string) (FOrderInfo, error) {
	var resp FOrderInfo
	req := make(map[string]interface{})
	req["symbol"] = symbol
	if orderID != "" {
		req["order_id"] = orderID
	}
	if clientOrderID != "" {
		req["client_order_id"] = clientOrderID
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fOrderInfo, nil, req, &resp)
}

// FOrderDetails gets order details for futures orders
func (h *HUOBI) FOrderDetails(symbol, orderID, orderType string, createdAt time.Time, pageIndex, pageSize int64) (FOrderDetailsData, error) {
	var resp FOrderDetailsData
	req := make(map[string]interface{})
	req["symbol"] = symbol
	req["order_id"] = orderID
	req["created_at"] = strconv.FormatInt(createdAt.Unix(), 10)
	oType, ok := validOrderType[orderType]
	if !ok {
		return resp, fmt.Errorf("invalid orderType")
	}
	req["order_type"] = oType
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize != 0 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fOrderDetails, nil, req, &resp)
}

// FGetOpenOrders gets order details for futures orders
func (h *HUOBI) FGetOpenOrders(symbol string, pageIndex, pageSize int64) (FOpenOrdersData, error) {
	var resp FOpenOrdersData
	req := make(map[string]interface{})
	req["symbol"] = symbol
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize != 0 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fQueryOpenOrders, nil, req, &resp)
}

// FGetOrderHistory gets order order history for futures
func (h *HUOBI) FGetOrderHistory(symbol, tradeType, reqType, contractCode, orderType string, status []order.Status, createDate, pageIndex, pageSize int64) (FOrderHistoryData, error) {
	var resp FOrderHistoryData
	req := make(map[string]interface{})
	req["symbol"] = symbol
	tType, ok := validFuturesTradeType[tradeType]
	if !ok {
		return resp, fmt.Errorf("invalid tradeType")
	}
	req["trade_type"] = tType
	rType, ok := validFuturesReqType[reqType]
	if !ok {
		return resp, fmt.Errorf("invalid reqType")
	}
	req["type"] = rType
	var reqStatus string = "0"
	if len(status) > 0 {
		var firstTime bool = true
		for x := range status {
			sType, ok := validOrderStatus[status[x]]
			if !ok {
				return resp, fmt.Errorf("invalid status")
			}
			if firstTime {
				firstTime = false
				reqStatus = strconv.FormatInt(sType, 10)
				continue
			}
			reqStatus = reqStatus + "," + strconv.FormatInt(sType, 10)
		}
	}
	req["status"] = reqStatus
	if createDate < 0 || createDate > 90 {
		return resp, fmt.Errorf("invalid createDate")
	}
	req["create_date"] = createDate
	if contractCode != "" {
		req["contract_code"] = contractCode
	}
	if orderType != "" {
		oType, ok := validFuturesOrderTypes[orderType]
		if !ok {
			return resp, fmt.Errorf("invalid orderType")
		}
		req["order_type"] = oType
	}
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize != 0 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fOrderHistory, nil, req, &resp)
}

// FTradeHistory gets trade history data for futures
func (h *HUOBI) FTradeHistory(symbol, tradeType, contractCode string, createDate, pageIndex, pageSize int64) (FOrderHistoryData, error) {
	var resp FOrderHistoryData
	req := make(map[string]interface{})
	req["symbol"] = symbol
	tType, ok := validTradeType[tradeType]
	if !ok {
		return resp, fmt.Errorf("invalid tradeType")
	}
	req["trade_type"] = tType
	if contractCode != "" {
		req["contract_code"] = contractCode
	}
	if createDate <= 0 || createDate > 90 {
		return resp, fmt.Errorf("invalid createDate")
	}
	req["create_date"] = createDate
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize != 0 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fMatchResult, nil, req, &resp)
}

// FPlaceTriggerOrder places a trigger order for futures
func (h *HUOBI) FPlaceTriggerOrder(symbol, contractType, contractCode, triggerType, orderPriceType, direction, offset string, triggerPrice, orderPrice, volume, leverageRate float64) (FTriggerOrderData, error) {
	var resp FTriggerOrderData
	req := make(map[string]interface{})
	if symbol != "" {
		req["symbol"] = symbol
	}
	if contractType != "" {
		if !common.StringDataCompare(validContractTypes, contractType) {
			return resp, fmt.Errorf("invalid contractType")
		}
		req["contract_type"] = contractType
	}
	if contractCode != "" {
		req["contract_code"] = contractCode
	}
	tType, ok := validTriggerType[triggerType]
	if !ok {
		return resp, fmt.Errorf("invalid trigger type")
	}
	req["trigger_type"] = tType
	req["direction"] = direction
	if !common.StringDataCompare(validOffsetTypes, offset) {
		return resp, fmt.Errorf("invalid offset")
	}
	req["offset"] = offset
	req["trigger_price"] = triggerPrice
	req["volume"] = volume
	req["lever_rate"] = leverageRate
	req["order_price"] = orderPrice
	if !common.StringDataCompare(validOrderPriceType, orderPriceType) {
		return resp, fmt.Errorf("invalid order price type")
	}
	req["order_price_type"] = orderPriceType
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fTriggerOrder, nil, req, &resp)
}

// FCancelTriggerOrder cancels trigger order for futures
func (h *HUOBI) FCancelTriggerOrder(symbol, orderID string) (FCancelOrderData, error) {
	var resp FCancelOrderData
	req := make(map[string]interface{})
	req["symbol"] = symbol
	req["order_id"] = orderID
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fCancelTriggerOrder, nil, req, &resp)
}

// FCancelAllTriggerOrders cancels all trigger order for futures
func (h *HUOBI) FCancelAllTriggerOrders(symbol, contractCode, contractType string) (FCancelOrderData, error) {
	var resp FCancelOrderData
	req := make(map[string]interface{})
	req["symbol"] = symbol
	if contractCode != "" {
		req["contract_code"] = contractCode
	}
	if contractType != "" {
		if !common.StringDataCompare(validContractTypes, contractType) {
			return resp, nil
		}
		req["contract_type"] = contractType
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fCancelAllTriggerOrders, nil, req, &resp)
}

// FQueryTriggerOpenOrders queries open trigger orders for futures
func (h *HUOBI) FQueryTriggerOpenOrders(symbol, contractCode string, pageIndex, pageSize int64) (FTriggerOpenOrders, error) {
	var resp FTriggerOpenOrders
	req := make(map[string]interface{})
	req["symbol"] = symbol
	if contractCode != "" {
		req["contract_code"] = contractCode
	}
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize != 0 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fTriggerOpenOrders, nil, req, &resp)
}

// FQueryTriggerOrderHistory queries trigger order history for futures
func (h *HUOBI) FQueryTriggerOrderHistory(symbol, contractCode, tradeType, status string, createDate, pageIndex, pageSize int64) (FTriggerOrderHistoryData, error) {
	var resp FTriggerOrderHistoryData
	req := make(map[string]interface{})
	req["symbol"] = symbol
	if contractCode != "" {
		req["contract_code"] = contractCode
	}
	if tradeType != "" {
		tType, ok := validTradeType[tradeType]
		if !ok {
			return resp, fmt.Errorf("invalid tradeType")
		}
		req["trade_type"] = tType
	}
	validStatus, ok := validStatusTypes[status]
	if !ok {
		return resp, fmt.Errorf("invalid status")
	}
	req["status"] = validStatus
	if createDate <= 0 || createDate > 90 {
		return resp, fmt.Errorf("invalid createDate")
	}
	req["create_date"] = createDate
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize != 0 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, fTriggerOrderHistory, nil, req, &resp)
}

// Coin Margined Swaps

// QuerySwapIndexPriceInfo gets perpetual swap index's price info
func (h *HUOBI) QuerySwapIndexPriceInfo(code string) (SwapIndexPriceData, error) {
	var resp SwapIndexPriceData
	path := huobiSwapIndexPriceInfo
	if code != "" {
		params := url.Values{}
		params.Set("contract_code", code)
		path = huobiSwapIndexPriceInfo + params.Encode()
	}
	return resp, h.SendHTTPRequest(exchange.RestFutures, path, &resp)
}

// GetSwapPriceLimits gets price caps for perpetual futures
func (h *HUOBI) GetSwapPriceLimits(code string) (SwapPriceLimitsData, error) {
	var resp SwapPriceLimitsData
	params := url.Values{}
	params.Set("contract_code", code)
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiSwapPriceLimitation+params.Encode(),
		&resp)
}

// SwapOpenInterestInformation gets open interest data for perpetual futures
func (h *HUOBI) SwapOpenInterestInformation(code string) (SwapOpenInterestData, error) {
	var resp SwapOpenInterestData
	params := url.Values{}
	params.Set("contract_code", code)
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiSwapOpenInterestInfo+params.Encode(), &resp)
}

// GetSwapMarketDepth gets market depth for perpetual futures
func (h *HUOBI) GetSwapMarketDepth(code, dataType string) (SwapMarketDepthData, error) {
	var resp SwapMarketDepthData
	params := url.Values{}
	params.Set("contract_code", code)
	params.Set("type", dataType)
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiSwapMarketDepth+params.Encode(), &resp)
}

// GetSwapKlineData gets kline data for perpetual futures
func (h *HUOBI) GetSwapKlineData(code, period string, size int64, startTime, endTime time.Time) (SwapKlineData, error) {
	var resp SwapKlineData
	params := url.Values{}
	params.Set("contract_code", code)
	if !common.StringDataCompare(validPeriods, period) {
		return resp, fmt.Errorf("invalid period value received")
	}
	params.Set("period", period)
	if (size == 1) || (size > 2000) {
		return resp, fmt.Errorf("invalid size")
	}
	params.Set("size", strconv.FormatInt(size, 10))
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return resp, errors.New("startTime cannot be after endTime")
		}
		params.Set("start_time", strconv.FormatInt(startTime.Unix(), 10))
		params.Set("end_time", strconv.FormatInt(endTime.Unix(), 10))
	}
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiKLineData+params.Encode(), &resp)
}

// GetSwapMarketOverview gets market data overview for perpetual futures
func (h *HUOBI) GetSwapMarketOverview(code string) (MarketOverviewData, error) {
	var resp MarketOverviewData
	params := url.Values{}
	params.Set("contract_code", code)
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiMarketDataOverview+params.Encode(), &resp)
}

// GetLastTrade gets the last trade for a given perpetual contract
func (h *HUOBI) GetLastTrade(code string) (LastTradeData, error) {
	var resp LastTradeData
	params := url.Values{}
	params.Set("contract_code", code)
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiLastTradeContract+params.Encode(), &resp)
}

// GetBatchTrades gets batch trades for a specified contract (fetching size cannot be bigger than 2000)
func (h *HUOBI) GetBatchTrades(code string, size int64) (BatchTradesData, error) {
	var resp BatchTradesData
	params := url.Values{}
	params.Set("contract_code", code)
	if !(size > 1) && !(size < 2000) {
		return resp, fmt.Errorf("invalid size")
	}
	params.Set("size", strconv.FormatInt(size, 10))
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiRequestBatchOfTradingRecords+params.Encode(), &resp)
}

// GetInsuranceData gets insurance fund data and clawback rates
func (h *HUOBI) GetInsuranceData(code string) (InsuranceAndClawbackData, error) {
	var resp InsuranceAndClawbackData
	params := url.Values{}
	params.Set("contract_code", code)
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiInsuranceBalanceAndClawbackRate+params.Encode(), &resp)
}

// GetHistoricalInsuranceData gets historical insurance fund data and clawback rates
func (h *HUOBI) GetHistoricalInsuranceData(code string, pageIndex, pageSize int64) (HistoricalInsuranceFundBalance, error) {
	var resp HistoricalInsuranceFundBalance
	params := url.Values{}
	params.Set("contract_code", code)
	if pageIndex != 0 {
		params.Set("page_index", strconv.FormatInt(pageIndex, 10))
	}
	if pageSize != 0 {
		params.Set("page_size", strconv.FormatInt(pageIndex, 10))
	}
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiInsuranceBalanceHistory+params.Encode(), &resp)
}

// GetTieredAjustmentFactorInfo gets tiered adjustment factor data
func (h *HUOBI) GetTieredAjustmentFactorInfo(code string) (TieredAdjustmentFactorData, error) {
	var resp TieredAdjustmentFactorData
	params := url.Values{}
	params.Set("contract_code", code)
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiTieredAdjustmentFactor+params.Encode(), &resp)
}

// GetOpenInterestInfo gets open interest data
func (h *HUOBI) GetOpenInterestInfo(code, period, amountType string, size int64) (OpenInterestData, error) {
	var resp OpenInterestData
	params := url.Values{}
	params.Set("contract_code", code)
	if !common.StringDataCompare(validPeriods, period) {
		return resp, fmt.Errorf("invalid period value received")
	}
	params.Set("period", period)
	if !(size > 0 && size <= 1200) {
		return resp, fmt.Errorf("invalid size provided values from 1-1200 supported")
	}
	params.Set("size", strconv.FormatInt(size, 10))
	aType, ok := validAmountType[amountType]
	if !ok {
		return resp, fmt.Errorf("invalid trade type")
	}
	params.Set("amount_type", strconv.FormatInt(aType, 10))
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiOpenInterestInfo+params.Encode(), &resp)
}

// GetSystemStatusInfo gets system status data
func (h *HUOBI) GetSystemStatusInfo(code, period, amountType string, size int64) (SystemStatusData, error) {
	var resp SystemStatusData
	params := url.Values{}
	params.Set("contract_code", code)
	if !common.StringDataCompare(validPeriods, period) {
		return resp, fmt.Errorf("invalid period value received")
	}
	params.Set("period", period)
	if size > 0 && size <= 1200 {
		params.Set("size", strconv.FormatInt(size, 10))
	}
	aType, ok := validAmountType[amountType]
	if !ok {
		return resp, fmt.Errorf("invalid trade type")
	}
	params.Set("amount_type", strconv.FormatInt(aType, 10))
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiSwapSystemStatus+params.Encode(), &resp)
}

// GetTraderSentimentIndexAccount gets top trader sentiment function-account
func (h *HUOBI) GetTraderSentimentIndexAccount(code, period string) (TraderSentimentIndexAccountData, error) {
	var resp TraderSentimentIndexAccountData
	params := url.Values{}
	params.Set("contract_code", code)
	if !common.StringDataCompare(validPeriods, period) {
		return resp, fmt.Errorf("invalid period value received")
	}
	params.Set("period", period)
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiSwapSentimentAccountData+params.Encode(), &resp)
}

// GetTraderSentimentIndexPosition gets top trader sentiment function-position
func (h *HUOBI) GetTraderSentimentIndexPosition(code, period string) (TraderSentimentIndexPositionData, error) {
	var resp TraderSentimentIndexPositionData
	params := url.Values{}
	params.Set("contract_code", code)
	if !common.StringDataCompare(validPeriods, period) {
		return resp, fmt.Errorf("invalid period value received")
	}
	params.Set("period", period)
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiSwapSentimentPosition+params.Encode(), &resp)
}

// GetLiquidationOrders gets liquidation orders for a given perp
func (h *HUOBI) GetLiquidationOrders(code, tradeType string, pageIndex, pageSize, createDate int64) (LiquidationOrdersData, error) {
	var resp LiquidationOrdersData
	params := url.Values{}
	params.Set("contract_code", code)
	if createDate != 7 && createDate != 90 {
		return resp, fmt.Errorf("invalid createDate. 7 and 90 are the only supported values")
	}
	params.Set("create_date", strconv.FormatInt(createDate, 10))
	tType, ok := validTradeTypes[tradeType]
	if !ok {
		return resp, fmt.Errorf("invalid trade type")
	}
	params.Set("trade_type", strconv.FormatInt(tType, 10))
	if pageIndex != 0 {
		params.Set("page_index", strconv.FormatInt(pageIndex, 10))
	}
	if pageSize != 0 {
		params.Set("page_size", strconv.FormatInt(pageIndex, 10))
	}
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiSwapLiquidationOrders+params.Encode(), &resp)
}

// GetHistoricalFundingRates gets historical funding rates for perpetual futures
func (h *HUOBI) GetHistoricalFundingRates(code string, pageSize, pageIndex int64) (HistoricalFundingRateData, error) {
	var resp HistoricalFundingRateData
	params := url.Values{}
	params.Set("contract_code", code)
	if pageIndex != 0 {
		params.Set("page_index", strconv.FormatInt(pageIndex, 10))
	}
	if pageSize != 0 {
		params.Set("page_size", strconv.FormatInt(pageIndex, 10))
	}
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiSwapHistoricalFundingRate+params.Encode(), &resp)
}

// GetPremiumIndexKlineData gets kline data for premium index
func (h *HUOBI) GetPremiumIndexKlineData(code, period string, size int64) (PremiumIndexKlineData, error) {
	var resp PremiumIndexKlineData
	params := url.Values{}
	params.Set("contract_code", code)
	if !common.StringDataCompare(validPeriods, period) {
		return resp, fmt.Errorf("invalid period value received")
	}
	params.Set("period", period)
	if !(size > 1) && !(size < 2000) {
		return resp, fmt.Errorf("invalid size")
	}
	params.Set("size", strconv.FormatInt(size, 10))
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiPremiumIndexKlineData+params.Encode(), &resp)
}

// GetEstimatedFundingRates gets estimated funding rates for perpetual futures
func (h *HUOBI) GetEstimatedFundingRates(code, period string, size int64) (EstimatedFundingRateData, error) {
	var resp EstimatedFundingRateData
	params := url.Values{}
	params.Set("contract_code", code)
	if !common.StringDataCompare(validPeriods, period) {
		return resp, fmt.Errorf("invalid period value received")
	}
	params.Set("period", period)
	if !(size > 0 && size <= 1200) {
		return resp, fmt.Errorf("invalid size provided values from 1-1200 supported")
	}
	params.Set("size", strconv.FormatInt(size, 10))
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiPredictedFundingRateData+params.Encode(), &resp)
}

// GetBasisData gets basis data for perpetual futures
func (h *HUOBI) GetBasisData(code, period, basisPriceType string, size int64) (BasisData, error) {
	var resp BasisData
	params := url.Values{}
	params.Set("contract_code", code)
	if !common.StringDataCompare(validPeriods, period) {
		return resp, fmt.Errorf("invalid period value received")
	}
	params.Set("period", period)
	if !(size > 0 && size <= 1200) {
		return resp, fmt.Errorf("invalid size provided values from 1-1200 supported")
	}
	params.Set("size", strconv.FormatInt(size, 10))
	if !common.StringDataCompare(validBasisPriceTypes, basisPriceType) {
		return resp, fmt.Errorf("invalid period value received")
	}
	return resp, h.SendHTTPRequest(exchange.RestFutures, huobiBasisData+params.Encode(), &resp)
}

// GetSwapAccountInfo gets swap account info
func (h *HUOBI) GetSwapAccountInfo(code string) (SwapAccountInformation, error) {
	var resp SwapAccountInformation
	req := make(map[string]interface{})
	req["contract_code"] = code
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapAccInfo, nil, req, &resp)
}

// GetSwapPositionsInfo gets swap positions' info
func (h *HUOBI) GetSwapPositionsInfo(code string) (SwapPositionInfo, error) {
	var resp SwapPositionInfo
	req := make(map[string]interface{})
	req["contract_code"] = code
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapPosInfo, nil, req, &resp)
}

// GetSwapAssetsAndPositions gets swap positions and asset info
func (h *HUOBI) GetSwapAssetsAndPositions(code string) (SwapAssetsAndPositionsData, error) {
	var resp SwapAssetsAndPositionsData
	req := make(map[string]interface{})
	req["contract_code"] = code
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapAssetsAndPos, nil, req, &resp)
}

// GetSwapAllSubAccAssets gets asset info for all subaccounts
func (h *HUOBI) GetSwapAllSubAccAssets(code string) (SubAccountsAssetData, error) {
	var resp SubAccountsAssetData
	req := make(map[string]interface{})
	if code != "" {
		req["contract_code"] = code
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapSubAccList, nil, req, &resp)
}

// SwapSingleSubAccAssets gets a subaccount's assets info
func (h *HUOBI) SwapSingleSubAccAssets(code string, subUID int64) (SingleSubAccountAssetsInfo, error) {
	var resp SingleSubAccountAssetsInfo
	req := make(map[string]interface{})
	req["contract_code"] = code
	req["sub_uid"] = subUID
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapSubAccInfo, nil, req, &resp)
}

// GetSubAccPositionInfo gets a subaccount's positions info
func (h *HUOBI) GetSubAccPositionInfo(code string, subUID int64) (SingleSubAccountPositionsInfo, error) {
	var resp SingleSubAccountPositionsInfo
	req := make(map[string]interface{})
	req["contract_code"] = code
	req["sub_uid"] = subUID
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapSubAccPosInfo, nil, req, &resp)
}

// GetAccountFinancialRecords gets the account's financial records
func (h *HUOBI) GetAccountFinancialRecords(code, orderType string, createDate, pageIndex, pageSize int64) (FinancialRecordData, error) {
	var resp FinancialRecordData
	req := make(map[string]interface{})
	req["contract_code"] = code
	if orderType != "" {
		req["type"] = orderType
	}
	if createDate != 0 {
		req["create_date"] = createDate
	}
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize != 0 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapFinancialRecords, nil, req, &resp)
}

// GetSwapSettlementRecords gets the swap account's settlement records
func (h *HUOBI) GetSwapSettlementRecords(code string, startTime, endTime time.Time, pageIndex, pageSize int64) (FinancialRecordData, error) {
	var resp FinancialRecordData
	req := make(map[string]interface{})
	req["contract_code"] = code
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return resp, errors.New("startTime cannot be after endTime")
		}
		req["start_time"] = strconv.FormatInt(startTime.Unix(), 10)
		req["end_time"] = strconv.FormatInt(endTime.Unix(), 10)
	}
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize != 0 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapSettlementRecords, nil, req, &resp)
}

// GetAvailableLeverage gets user's available leverage data
func (h *HUOBI) GetAvailableLeverage(code string) (AvailableLeverageData, error) {
	var resp AvailableLeverageData
	req := make(map[string]interface{})
	if code != "" {
		req["contract_code"] = code
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapAvailableLeverage, nil, req, &resp)
}

// GetSwapOrderLimitInfo gets order limit info for swaps
func (h *HUOBI) GetSwapOrderLimitInfo(code, orderType string) (SwapOrderLimitInfo, error) {
	var resp SwapOrderLimitInfo
	req := make(map[string]interface{})
	req["contract_code"] = code
	if !common.StringDataCompare(validOrderTypes, orderType) {
		return resp, fmt.Errorf("inavlid ordertype provided")
	}
	req["order_price_type"] = orderType
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapOrderLimitInfo, nil, req, &resp)
}

// GetSwapTradingFeeInfo gets trading fee info for swaps
func (h *HUOBI) GetSwapTradingFeeInfo(code string) (SwapTradingFeeData, error) {
	var resp SwapTradingFeeData
	req := make(map[string]interface{})
	req["contract_code"] = code
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapTradingFeeInfo, nil, req, &resp)
}

// GetSwapTransferLimitInfo gets transfer limit info for swaps
func (h *HUOBI) GetSwapTransferLimitInfo(code string) (TransferLimitData, error) {
	var resp TransferLimitData
	req := make(map[string]interface{})
	req["contract_code"] = code
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapTransferLimitInfo, nil, req, &resp)
}

// GetSwapPositionLimitInfo gets transfer limit info for swaps
func (h *HUOBI) GetSwapPositionLimitInfo(code string) (PositionLimitData, error) {
	var resp PositionLimitData
	req := make(map[string]interface{})
	req["contract_code"] = code
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapPositionLimitInfo, nil, req, &resp)
}

// AccountTransferData gets asset transfer data between master and subaccounts
func (h *HUOBI) AccountTransferData(code, subUID, transferType string, amount float64) (InternalAccountTransferData, error) {
	var resp InternalAccountTransferData
	req := make(map[string]interface{})
	req["contract_code"] = code
	req["subUid"] = subUID
	req["amount"] = amount
	if !common.StringDataCompare(validTransferType, transferType) {
		return resp, fmt.Errorf("inavlid transferType received")
	}
	req["type"] = transferType
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapInternalTransferData, nil, req, &resp)
}

// AccountTransferRecords gets asset transfer records between master and subaccounts
func (h *HUOBI) AccountTransferRecords(code, transferType string, createDate, pageIndex, pageSize int64) (InternalAccountTransferData, error) {
	var resp InternalAccountTransferData
	req := make(map[string]interface{})
	req["contract_code"] = code
	if !common.StringDataCompare(validTransferType, transferType) {
		return resp, fmt.Errorf("inavlid transferType received")
	}
	req["type"] = transferType
	if createDate > 90 {
		return resp, fmt.Errorf("invalid create date value: only supports up to 90 days")
	}
	req["create_date"] = strconv.FormatInt(createDate, 10)
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize > 0 && pageSize <= 50 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapInternalTransferRecords, nil, req, &resp)
}

// PlaceSwapOrders places orders for swaps
func (h *HUOBI) PlaceSwapOrders(code, clientOrderID, direction, offset, orderPriceType string, price, volume, leverage float64) (SwapOrderData, error) {
	var resp SwapOrderData
	req := make(map[string]interface{})
	req["contract_code"] = code
	if clientOrderID != "" {
		req["client_order_id"] = clientOrderID
	}
	req["direction"] = direction
	req["offset"] = offset
	if !common.StringDataCompare(validOrderTypes, orderPriceType) {
		return resp, fmt.Errorf("inavlid ordertype provided")
	}
	req["order_price_type"] = orderPriceType
	req["price"] = price
	req["volume"] = volume
	req["lever_rate"] = leverage
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapPlaceOrder, nil, req, &resp)
}

// PlaceSwapBatchOrders places a batch of orders for swaps
func (h *HUOBI) PlaceSwapBatchOrders(data BatchOrderRequestType) (BatchOrderData, error) {
	var resp BatchOrderData
	req := make(map[string]interface{})
	if (len(data.Data) > 10) || len(data.Data) == 0 {
		return resp, fmt.Errorf("invalid data provided: maximum of 10 batch orders supported")
	}
	req["orders_data"] = data.Data
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapPlaceBatchOrder, nil, req, &resp)
}

// CancelSwapOrder sends a request to cancel an order
func (h *HUOBI) CancelSwapOrder(orderID, clientOrderID, contractCode string) (CancelOrdersData, error) {
	var resp CancelOrdersData
	req := make(map[string]interface{})
	if orderID != "" {
		req["order_id"] = orderID
	}
	if clientOrderID != "" {
		req["client_order_id"] = clientOrderID
	}
	req["contract_code"] = contractCode
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapCancelOrder, nil, req, &resp)
}

// CancelAllSwapOrders sends a request to cancel an order
func (h *HUOBI) CancelAllSwapOrders(contractCode string) (CancelOrdersData, error) {
	var resp CancelOrdersData
	req := make(map[string]interface{})
	req["contract_code"] = contractCode
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapCancelAllOrders, nil, req, &resp)
}

// PlaceLightningCloseOrder places a lightning close order
func (h *HUOBI) PlaceLightningCloseOrder(contractCode, direction, orderPriceType string, volume float64, clientOrderID int64) (LightningCloseOrderData, error) {
	var resp LightningCloseOrderData
	req := make(map[string]interface{})
	req["contract_code"] = contractCode
	req["volume"] = volume
	req["direction"] = direction
	if clientOrderID != 0 {
		req["client_order_id"] = clientOrderID
	}
	if orderPriceType != "" {
		if !common.StringDataCompare(validLightningOrderPriceType, orderPriceType) {
			return resp, fmt.Errorf("invalid orderPriceType")
		}
		req["order_price_type"] = orderPriceType
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapLightningCloseOrder, nil, req, &resp)
}

// GetSwapOrderDetails gets order info
func (h *HUOBI) GetSwapOrderDetails(contractCode, orderID, createdAt, orderType string, pageIndex, pageSize int64) (SwapOrderData, error) {
	var resp SwapOrderData
	req := make(map[string]interface{})
	req["contract_code"] = contractCode
	req["order_id"] = orderID
	req["created_at"] = createdAt
	oType, ok := validOrderType[orderType]
	if !ok {
		return resp, fmt.Errorf("invalid ordertype")
	}
	req["order_type"] = oType
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize > 0 && pageSize <= 50 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapOrderDetails, nil, req, &resp)
}

// GetSwapOrderInfo gets info on a swap order
func (h *HUOBI) GetSwapOrderInfo(contractCode, orderID, clientOrderID string) (SwapOrderInfo, error) {
	var resp SwapOrderInfo
	req := make(map[string]interface{})
	if contractCode != "" {
		req["contract_code"] = contractCode
	}
	if orderID != "" {
		req["order_id"] = orderID
	}
	if clientOrderID != "" {
		req["client_order_id"] = clientOrderID
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapOrderInfo, nil, req, &resp)
}

// GetSwapOpenOrders gets open orders for swap
func (h *HUOBI) GetSwapOpenOrders(contractCode string, pageIndex, pageSize int64) (SwapOpenOrdersData, error) {
	var resp SwapOpenOrdersData
	req := make(map[string]interface{})
	req["contract_code"] = contractCode
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize > 0 && pageSize <= 50 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapOpenOrders, nil, req, &resp)
}

// GetSwapOrderHistory gets swap order history
func (h *HUOBI) GetSwapOrderHistory(contractCode, tradeType, reqType string, status []order.Status, createDate, pageIndex, pageSize int64) (SwapOrderHistory, error) {
	var resp SwapOrderHistory
	req := make(map[string]interface{})
	req["contract_code"] = contractCode
	tType, ok := validFuturesTradeType[tradeType]
	if !ok {
		return resp, fmt.Errorf("invalid tradeType")
	}
	req["trade_type"] = tType
	rType, ok := validFuturesReqType[reqType]
	if !ok {
		return resp, fmt.Errorf("invalid reqType")
	}
	req["type"] = rType
	reqStatus := "0"
	if len(status) > 0 {
		firstTime := true
		for x := range status {
			sType, ok := validOrderStatus[status[x]]
			if !ok {
				return resp, fmt.Errorf("invalid status")
			}
			if firstTime {
				firstTime = false
				reqStatus = strconv.FormatInt(sType, 10)
				continue
			}
			reqStatus = reqStatus + "," + strconv.FormatInt(sType, 10)
		}
	}
	req["status"] = reqStatus
	if createDate < 0 || createDate > 90 {
		return resp, fmt.Errorf("invalid createDate")
	}
	req["create_date"] = createDate
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize != 0 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapOrderHistory, nil, req, &resp)
}

// GetSwapTradeHistory gets swap trade history
func (h *HUOBI) GetSwapTradeHistory(contractCode, tradeType string, createDate, pageIndex, pageSize int64) (AccountTradeHistoryData, error) {
	var resp AccountTradeHistoryData
	req := make(map[string]interface{})
	req["contract_code"] = contractCode
	if createDate > 90 {
		return resp, fmt.Errorf("invalid create date value: only supports up to 90 days")
	}
	tType, ok := validTradeType[tradeType]
	if !ok {
		return resp, fmt.Errorf("invalid trade type")
	}
	req["trade_type"] = tType
	req["create_date"] = strconv.FormatInt(createDate, 10)
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize > 0 && pageSize <= 50 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapTradeHistory, nil, req, &resp)
}

// PlaceSwapTriggerOrder places a trigger order for a swap
func (h *HUOBI) PlaceSwapTriggerOrder(contractCode, triggerType, direction, offset, orderPriceType string, triggerPrice, orderPrice, volume, leverageRate float64) (AccountTradeHistoryData, error) {
	var resp AccountTradeHistoryData
	req := make(map[string]interface{})
	req["contract_code"] = contractCode
	tType, ok := validTriggerType[triggerType]
	if !ok {
		return resp, fmt.Errorf("invalid trigger type")
	}
	req["trigger_type"] = tType
	req["direction"] = direction
	req["offset"] = offset
	req["trigger_price"] = triggerPrice
	req["volume"] = volume
	req["lever_rate"] = leverageRate
	req["order_price"] = orderPrice
	if !common.StringDataCompare(validOrderPriceType, orderPriceType) {
		return resp, fmt.Errorf("invalid order price type")
	}
	req["order_price_type"] = orderPriceType
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapTriggerOrder, nil, req, &resp)
}

// CancelSwapTriggerOrder cancels swap trigger order
func (h *HUOBI) CancelSwapTriggerOrder(contractCode, orderID string) (CancelTriggerOrdersData, error) {
	var resp CancelTriggerOrdersData
	req := make(map[string]interface{})
	req["contract_code"] = contractCode
	req["order_id"] = orderID
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapCancelTriggerOrder, nil, req, &resp)
}

// CancelAllSwapTriggerOrders cancels all swap trigger orders
func (h *HUOBI) CancelAllSwapTriggerOrders(contractCode string) (CancelTriggerOrdersData, error) {
	var resp CancelTriggerOrdersData
	req := make(map[string]interface{})
	req["contract_code"] = contractCode
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapCancelAllTriggerOrders, nil, req, &resp)
}

// GetSwapTriggerOrderHistory gets history for swap trigger orders
func (h *HUOBI) GetSwapTriggerOrderHistory(contractCode, status, tradeType string, createDate, pageIndex, pageSize int64) (TriggerOrderHistory, error) {
	var resp TriggerOrderHistory
	req := make(map[string]interface{})
	req["contract_code"] = contractCode
	req["status"] = status
	tType, ok := validTradeType[tradeType]
	if !ok {
		return resp, fmt.Errorf("invalid trade type")
	}
	req["trade_type"] = tType
	if createDate > 90 {
		return resp, fmt.Errorf("invalid create date value: only supports up to 90 days")
	}
	req["create_date"] = strconv.FormatInt(createDate, 10)
	if pageIndex != 0 {
		req["page_index"] = pageIndex
	}
	if pageSize > 0 && pageSize <= 50 {
		req["page_size"] = pageSize
	}
	return resp, h.FuturesAuthenticatedHTTPRequest(exchange.RestFutures, http.MethodPost, huobiSwapTriggerOrderHistory, nil, req, &resp)
}

// GetSwapMarkets gets data of swap markets
func (h *HUOBI) GetSwapMarkets(contract string) ([]SwapMarketsData, error) {
	vals := url.Values{}
	vals.Set("contract_code", contract)
	type response struct {
		Response
		Data []SwapMarketsData `json:"data"`
	}
	var result response
	err := h.SendHTTPRequest(exchange.RestFutures, huobiSwapMarkets+vals.Encode(), &result)
	if result.ErrorMessage != "" {
		return nil, errors.New(result.ErrorMessage)
	}
	return result.Data, err
}

// GetSwapFundingRates gets funding rates data
func (h *HUOBI) GetSwapFundingRates(contract string) (FundingRatesData, error) {
	vals := url.Values{}
	vals.Set("contract_code", contract)
	type response struct {
		Response
		Data FundingRatesData `json:"data"`
	}
	var result response
	err := h.SendHTTPRequest(exchange.RestFutures, huobiSwapFunding+vals.Encode(), &result)
	if result.ErrorMessage != "" {
		return FundingRatesData{}, errors.New(result.ErrorMessage)
	}
	return result.Data, err
}

// SPOT section below

// GetMarginRates gets margin rates
func (h *HUOBI) GetMarginRates(symbol string) (MarginRatesData, error) {
	vals := url.Values{}
	if symbol != "" {
		vals.Set("symbols", symbol)
	}
	var resp MarginRatesData
	return resp, h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodGet, huobiMarginRates, vals, nil, &resp, false)
}

// GetSpotKline returns kline data
// KlinesRequestParams contains symbol, period and size
func (h *HUOBI) GetSpotKline(arg KlinesRequestParams) ([]KlineItem, error) {
	vals := url.Values{}
	vals.Set("symbol", arg.Symbol)
	vals.Set("period", arg.Period)

	if arg.Size != 0 {
		vals.Set("size", strconv.Itoa(arg.Size))
	}

	type response struct {
		Response
		Data []KlineItem `json:"data"`
	}

	var result response
	urlPath := fmt.Sprintf("/%s", huobiMarketHistoryKline)

	err := h.SendHTTPRequest(exchange.RestSpot, common.EncodeURLValues(urlPath, vals), &result)
	if result.ErrorMessage != "" {
		return nil, errors.New(result.ErrorMessage)
	}
	return result.Data, err
}

// GetTickers returns the ticker for the specified symbol
func (h *HUOBI) GetTickers() (Tickers, error) {
	var result Tickers
	urlPath := fmt.Sprintf("/%s", huobiMarketTickers)
	return result, h.SendHTTPRequest(exchange.RestSpot, urlPath, &result)
}

// GetMarketDetailMerged returns the ticker for the specified symbol
func (h *HUOBI) GetMarketDetailMerged(symbol string) (DetailMerged, error) {
	vals := url.Values{}
	vals.Set("symbol", symbol)

	type response struct {
		Response
		Tick DetailMerged `json:"tick"`
	}

	var result response
	urlPath := fmt.Sprintf("/%s", huobiMarketDetailMerged)

	err := h.SendHTTPRequest(exchange.RestSpot, common.EncodeURLValues(urlPath, vals), &result)
	if result.ErrorMessage != "" {
		return result.Tick, errors.New(result.ErrorMessage)
	}
	return result.Tick, err
}

// GetDepth returns the depth for the specified symbol
func (h *HUOBI) GetDepth(obd OrderBookDataRequestParams) (Orderbook, error) {
	vals := url.Values{}
	vals.Set("symbol", obd.Symbol)

	if obd.Type != OrderBookDataRequestParamsTypeNone {
		vals.Set("type", string(obd.Type))
	}

	type response struct {
		Response
		Depth Orderbook `json:"tick"`
	}

	var result response
	urlPath := fmt.Sprintf("/%s", huobiMarketDepth)

	err := h.SendHTTPRequest(exchange.RestSpot, common.EncodeURLValues(urlPath, vals), &result)
	if result.ErrorMessage != "" {
		return result.Depth, errors.New(result.ErrorMessage)
	}
	return result.Depth, err
}

// GetTrades returns the trades for the specified symbol
func (h *HUOBI) GetTrades(symbol string) ([]Trade, error) {
	vals := url.Values{}
	vals.Set("symbol", symbol)

	type response struct {
		Response
		Tick struct {
			Data []Trade `json:"data"`
		} `json:"tick"`
	}

	var result response
	urlPath := fmt.Sprintf("/%s", huobiMarketTrade)

	err := h.SendHTTPRequest(exchange.RestSpot, common.EncodeURLValues(urlPath, vals), &result)
	if result.ErrorMessage != "" {
		return nil, errors.New(result.ErrorMessage)
	}
	return result.Tick.Data, err
}

// GetLatestSpotPrice returns latest spot price of symbol
//
// symbol: string of currency pair
func (h *HUOBI) GetLatestSpotPrice(symbol string) (float64, error) {
	list, err := h.GetTradeHistory(symbol, 1)

	if err != nil {
		return 0, err
	}
	if len(list) == 0 {
		return 0, errors.New("the length of the list is 0")
	}

	return list[0].Trades[0].Price, nil
}

// GetTradeHistory returns the trades for the specified symbol
func (h *HUOBI) GetTradeHistory(symbol string, size int64) ([]TradeHistory, error) {
	vals := url.Values{}
	vals.Set("symbol", symbol)

	if size > 0 {
		vals.Set("size", strconv.FormatInt(size, 10))
	}

	type response struct {
		Response
		TradeHistory []TradeHistory `json:"data"`
	}

	var result response
	urlPath := fmt.Sprintf("/%s", huobiMarketTradeHistory)

	err := h.SendHTTPRequest(exchange.RestSpot, common.EncodeURLValues(urlPath, vals), &result)
	if result.ErrorMessage != "" {
		return nil, errors.New(result.ErrorMessage)
	}
	return result.TradeHistory, err
}

// GetMarketDetail returns the ticker for the specified symbol
func (h *HUOBI) GetMarketDetail(symbol string) (Detail, error) {
	vals := url.Values{}
	vals.Set("symbol", symbol)

	type response struct {
		Response
		Tick Detail `json:"tick"`
	}

	var result response
	urlPath := fmt.Sprintf("/%s", huobiMarketDetail)

	err := h.SendHTTPRequest(exchange.RestSpot, common.EncodeURLValues(urlPath, vals), &result)
	if result.ErrorMessage != "" {
		return result.Tick, errors.New(result.ErrorMessage)
	}
	return result.Tick, err
}

// GetSymbols returns an array of symbols supported by Huobi
func (h *HUOBI) GetSymbols() ([]Symbol, error) {
	type response struct {
		Response
		Symbols []Symbol `json:"data"`
	}

	var result response
	urlPath := fmt.Sprintf("/v%s/%s", huobiAPIVersion, huobiSymbols)

	err := h.SendHTTPRequest(exchange.RestSpot, urlPath, &result)
	if result.ErrorMessage != "" {
		return nil, errors.New(result.ErrorMessage)
	}
	return result.Symbols, err
}

// GetCurrencies returns a list of currencies supported by Huobi
func (h *HUOBI) GetCurrencies() ([]string, error) {
	type response struct {
		Response
		Currencies []string `json:"data"`
	}

	var result response
	urlPath := fmt.Sprintf("/v%s/%s", huobiAPIVersion, huobiCurrencies)

	err := h.SendHTTPRequest(exchange.RestSpot, urlPath, &result)
	if result.ErrorMessage != "" {
		return nil, errors.New(result.ErrorMessage)
	}
	return result.Currencies, err
}

// GetTimestamp returns the Huobi server time
func (h *HUOBI) GetTimestamp() (int64, error) {
	type response struct {
		Response
		Timestamp int64 `json:"data"`
	}

	var result response
	urlPath := fmt.Sprintf("/v%s/%s", huobiAPIVersion, huobiTimestamp)

	err := h.SendHTTPRequest(exchange.RestSpot, urlPath, &result)
	if result.ErrorMessage != "" {
		return 0, errors.New(result.ErrorMessage)
	}
	return result.Timestamp, err
}

// GetAccounts returns the Huobi user accounts
func (h *HUOBI) GetAccounts() ([]Account, error) {
	result := struct {
		Accounts []Account `json:"data"`
	}{}
	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodGet, huobiAccounts, url.Values{}, nil, &result, false)
	return result.Accounts, err
}

// GetAccountBalance returns the users Huobi account balance
func (h *HUOBI) GetAccountBalance(accountID string) ([]AccountBalanceDetail, error) {
	result := struct {
		AccountBalanceData AccountBalance `json:"data"`
	}{}
	endpoint := fmt.Sprintf(huobiAccountBalance, accountID)
	v := url.Values{}
	v.Set("account-id", accountID)
	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodGet, endpoint, v, nil, &result, false)
	return result.AccountBalanceData.AccountBalanceDetails, err
}

// GetAggregatedBalance returns the balances of all the sub-account aggregated.
func (h *HUOBI) GetAggregatedBalance() ([]AggregatedBalance, error) {
	result := struct {
		AggregatedBalances []AggregatedBalance `json:"data"`
	}{}
	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot,
		http.MethodGet,
		huobiAggregatedBalance,
		nil,
		nil,
		&result,
		false,
	)
	return result.AggregatedBalances, err
}

// SpotNewOrder submits an order to Huobi
func (h *HUOBI) SpotNewOrder(arg SpotNewOrderRequestParams) (int64, error) {
	data := struct {
		AccountID int    `json:"account-id,string"`
		Amount    string `json:"amount"`
		Price     string `json:"price"`
		Source    string `json:"source"`
		Symbol    string `json:"symbol"`
		Type      string `json:"type"`
	}{
		AccountID: arg.AccountID,
		Amount:    strconv.FormatFloat(arg.Amount, 'f', -1, 64),
		Symbol:    arg.Symbol,
		Type:      string(arg.Type),
	}

	// Only set price if order type is not equal to buy-market or sell-market
	if arg.Type != SpotNewOrderRequestTypeBuyMarket && arg.Type != SpotNewOrderRequestTypeSellMarket {
		data.Price = strconv.FormatFloat(arg.Price, 'f', -1, 64)
	}

	if arg.Source != "" {
		data.Source = arg.Source
	}

	result := struct {
		OrderID int64 `json:"data,string"`
	}{}
	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot,
		http.MethodPost,
		huobiOrderPlace,
		nil,
		data,
		&result,
		false,
	)
	return result.OrderID, err
}

// CancelExistingOrder cancels an order on Huobi
func (h *HUOBI) CancelExistingOrder(orderID int64) (int64, error) {
	resp := struct {
		OrderID int64 `json:"data,string"`
	}{}
	endpoint := fmt.Sprintf(huobiOrderCancel, strconv.FormatInt(orderID, 10))
	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodPost, endpoint, url.Values{}, nil, &resp, false)
	return resp.OrderID, err
}

// CancelOrderBatch cancels a batch of orders -- to-do
func (h *HUOBI) CancelOrderBatch(_ []int64) ([]CancelOrderBatch, error) {
	type response struct {
		Response
		Data []CancelOrderBatch `json:"data"`
	}

	var result response
	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodPost, huobiOrderCancelBatch, url.Values{}, nil, &result, false)

	if result.ErrorMessage != "" {
		return nil, errors.New(result.ErrorMessage)
	}
	return result.Data, err
}

// CancelOpenOrdersBatch cancels a batch of orders -- to-do
func (h *HUOBI) CancelOpenOrdersBatch(accountID, symbol string) (CancelOpenOrdersBatch, error) {
	params := url.Values{}

	params.Set("account-id", accountID)
	var result CancelOpenOrdersBatch

	data := struct {
		AccountID string `json:"account-id"`
		Symbol    string `json:"symbol"`
	}{
		AccountID: accountID,
		Symbol:    symbol,
	}

	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodPost, huobiBatchCancelOpenOrders, url.Values{}, data, &result, false)
	if result.Data.FailedCount > 0 {
		return result, fmt.Errorf("there were %v failed order cancellations", result.Data.FailedCount)
	}

	return result, err
}

// GetOrder returns order information for the specified order
func (h *HUOBI) GetOrder(orderID int64) (OrderInfo, error) {
	resp := struct {
		Order OrderInfo `json:"data"`
	}{}
	urlVal := url.Values{}
	urlVal.Set("clientOrderId", strconv.FormatInt(orderID, 10))
	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodGet,
		huobiGetOrder,
		urlVal,
		nil,
		&resp,
		false)
	return resp.Order, err
}

// GetOrderMatchResults returns matched order info for the specified order
func (h *HUOBI) GetOrderMatchResults(orderID int64) ([]OrderMatchInfo, error) {
	resp := struct {
		Orders []OrderMatchInfo `json:"data"`
	}{}
	endpoint := fmt.Sprintf(huobiGetOrderMatch, strconv.FormatInt(orderID, 10))
	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodGet, endpoint, url.Values{}, nil, &resp, false)
	return resp.Orders, err
}

// GetOrders returns a list of orders
func (h *HUOBI) GetOrders(symbol, types, start, end, states, from, direct, size string) ([]OrderInfo, error) {
	resp := struct {
		Orders []OrderInfo `json:"data"`
	}{}

	vals := url.Values{}
	vals.Set("symbol", symbol)
	vals.Set("states", states)

	if types != "" {
		vals.Set("types", types)
	}

	if start != "" {
		vals.Set("start-date", start)
	}

	if end != "" {
		vals.Set("end-date", end)
	}

	if from != "" {
		vals.Set("from", from)
	}

	if direct != "" {
		vals.Set("direct", direct)
	}

	if size != "" {
		vals.Set("size", size)
	}

	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodGet, huobiGetOrders, vals, nil, &resp, false)
	return resp.Orders, err
}

// GetOpenOrders returns a list of orders
func (h *HUOBI) GetOpenOrders(accountID, symbol, side string, size int64) ([]OrderInfo, error) {
	resp := struct {
		Orders []OrderInfo `json:"data"`
	}{}

	vals := url.Values{}
	vals.Set("symbol", symbol)
	vals.Set("accountID", accountID)
	if len(side) > 0 {
		vals.Set("side", side)
	}
	vals.Set("size", strconv.FormatInt(size, 10))

	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodGet, huobiGetOpenOrders, vals, nil, &resp, false)
	return resp.Orders, err
}

// GetOrdersMatch returns a list of matched orders
func (h *HUOBI) GetOrdersMatch(symbol, types, start, end, from, direct, size string) ([]OrderMatchInfo, error) {
	resp := struct {
		Orders []OrderMatchInfo `json:"data"`
	}{}

	vals := url.Values{}
	vals.Set("symbol", symbol)

	if types != "" {
		vals.Set("types", types)
	}

	if start != "" {
		vals.Set("start-date", start)
	}

	if end != "" {
		vals.Set("end-date", end)
	}

	if from != "" {
		vals.Set("from", from)
	}

	if direct != "" {
		vals.Set("direct", direct)
	}

	if size != "" {
		vals.Set("size", size)
	}

	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodGet, huobiGetOrdersMatch, vals, nil, &resp, false)
	return resp.Orders, err
}

// MarginTransfer transfers assets into or out of the margin account
func (h *HUOBI) MarginTransfer(symbol, currency string, amount float64, in bool) (int64, error) {
	data := struct {
		Symbol   string `json:"symbol"`
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	}{
		Symbol:   symbol,
		Currency: currency,
		Amount:   strconv.FormatFloat(amount, 'f', -1, 64),
	}

	path := huobiMarginTransferIn
	if !in {
		path = huobiMarginTransferOut
	}

	resp := struct {
		TransferID int64 `json:"data"`
	}{}
	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodPost, path, nil, data, &resp, false)
	return resp.TransferID, err
}

// MarginOrder submits a margin order application
func (h *HUOBI) MarginOrder(symbol, currency string, amount float64) (int64, error) {
	data := struct {
		Symbol   string `json:"symbol"`
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	}{
		Symbol:   symbol,
		Currency: currency,
		Amount:   strconv.FormatFloat(amount, 'f', -1, 64),
	}

	resp := struct {
		MarginOrderID int64 `json:"data"`
	}{}
	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodPost, huobiMarginOrders, nil, data, &resp, false)
	return resp.MarginOrderID, err
}

// MarginRepayment repays a margin amount for a margin ID
func (h *HUOBI) MarginRepayment(orderID int64, amount float64) (int64, error) {
	data := struct {
		Amount string `json:"amount"`
	}{
		Amount: strconv.FormatFloat(amount, 'f', -1, 64),
	}

	resp := struct {
		MarginOrderID int64 `json:"data"`
	}{}

	endpoint := fmt.Sprintf(huobiMarginRepay, strconv.FormatInt(orderID, 10))
	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodPost, endpoint, nil, data, &resp, false)
	return resp.MarginOrderID, err
}

// GetMarginLoanOrders returns the margin loan orders
func (h *HUOBI) GetMarginLoanOrders(symbol, currency, start, end, states, from, direct, size string) ([]MarginOrder, error) {
	vals := url.Values{}
	vals.Set("symbol", symbol)
	vals.Set("currency", currency)

	if start != "" {
		vals.Set("start-date", start)
	}

	if end != "" {
		vals.Set("end-date", end)
	}

	if states != "" {
		vals.Set("states", states)
	}

	if from != "" {
		vals.Set("from", from)
	}

	if direct != "" {
		vals.Set("direct", direct)
	}

	if size != "" {
		vals.Set("size", size)
	}

	resp := struct {
		MarginLoanOrders []MarginOrder `json:"data"`
	}{}
	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodGet, huobiMarginLoanOrders, vals, nil, &resp, false)
	return resp.MarginLoanOrders, err
}

// GetMarginAccountBalance returns the margin account balances
func (h *HUOBI) GetMarginAccountBalance(symbol string) ([]MarginAccountBalance, error) {
	resp := struct {
		Balances []MarginAccountBalance `json:"data"`
	}{}
	vals := url.Values{}
	if symbol != "" {
		vals.Set("symbol", symbol)
	}
	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodGet, huobiMarginAccountBalance, vals, nil, &resp, false)
	return resp.Balances, err
}

// Withdraw withdraws the desired amount and currency
func (h *HUOBI) Withdraw(c currency.Code, address, addrTag string, amount, fee float64) (int64, error) {
	resp := struct {
		WithdrawID int64 `json:"data"`
	}{}

	data := struct {
		Address  string `json:"address"`
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
		Fee      string `json:"fee,omitempty"`
		AddrTag  string `json:"addr-tag,omitempty"`
	}{
		Address:  address,
		Currency: c.Lower().String(),
		Amount:   strconv.FormatFloat(amount, 'f', -1, 64),
	}

	if fee > 0 {
		data.Fee = strconv.FormatFloat(fee, 'f', -1, 64)
	}

	if c == currency.XRP {
		data.AddrTag = addrTag
	}

	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodPost, huobiWithdrawCreate, nil, data, &resp.WithdrawID, false)
	return resp.WithdrawID, err
}

// CancelWithdraw cancels a withdraw request
func (h *HUOBI) CancelWithdraw(withdrawID int64) (int64, error) {
	resp := struct {
		WithdrawID int64 `json:"data"`
	}{}
	vals := url.Values{}
	vals.Set("withdraw-id", strconv.FormatInt(withdrawID, 10))

	endpoint := fmt.Sprintf(huobiWithdrawCancel, strconv.FormatInt(withdrawID, 10))
	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodPost, endpoint, vals, nil, &resp, false)
	return resp.WithdrawID, err
}

// QueryDepositAddress returns the deposit address for a specified currency
func (h *HUOBI) QueryDepositAddress(cryptocurrency string) (DepositAddress, error) {
	resp := struct {
		DepositAddress []DepositAddress `json:"data"`
	}{}

	vals := url.Values{}
	vals.Set("currency", cryptocurrency)

	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodGet, huobiAccountDepositAddress, vals, nil, &resp, true)
	if err != nil {
		return DepositAddress{}, err
	}
	if len(resp.DepositAddress) == 0 {
		return DepositAddress{}, errors.New("deposit address data isn't populated")
	}
	return resp.DepositAddress[0], nil
}

// QueryWithdrawQuotas returns the users cryptocurrency withdraw quotas
func (h *HUOBI) QueryWithdrawQuotas(cryptocurrency string) (WithdrawQuota, error) {
	resp := struct {
		WithdrawQuota WithdrawQuota `json:"data"`
	}{}

	vals := url.Values{}
	vals.Set("currency", cryptocurrency)

	err := h.SendAuthenticatedHTTPRequest(exchange.RestSpot, http.MethodGet, huobiAccountWithdrawQuota, vals, nil, &resp, true)
	if err != nil {
		return WithdrawQuota{}, err
	}
	return resp.WithdrawQuota, nil
}

// SendHTTPRequest sends an unauthenticated HTTP request
func (h *HUOBI) SendHTTPRequest(ep exchange.URL, path string, result interface{}) error {
	endpoint, err := h.API.Endpoints.GetURL(ep)
	if err != nil {
		return err
	}
	var tempResp json.RawMessage
	var errCap errorCapture
	err = h.SendPayload(context.Background(), &request.Item{
		Method:        http.MethodGet,
		Path:          endpoint + path,
		Result:        &tempResp,
		Verbose:       h.Verbose,
		HTTPDebugging: h.HTTPDebugging,
		HTTPRecording: h.HTTPRecording,
	})
	if err != nil {
		return err
	}
	if err := json.Unmarshal(tempResp, &errCap); err == nil {
		if errCap.Code != 200 && errCap.ErrMsg != "" {
			return errors.New(errCap.ErrMsg)
		}
	}
	return json.Unmarshal(tempResp, result)
}

// FuturesAuthenticatedHTTPRequest sends authenticated requests to the HUOBI API
func (h *HUOBI) FuturesAuthenticatedHTTPRequest(ep exchange.URL, method, endpoint string, values url.Values, data, result interface{}) error {
	if !h.AllowAuthenticatedRequest() {
		return fmt.Errorf(exchange.WarningAuthenticatedRequestWithoutCredentialsSet, h.Name)
	}
	epoint, err := h.API.Endpoints.GetURL(ep)
	if err != nil {
		return err
	}
	if values == nil {
		values = url.Values{}
	}
	now := time.Now()
	values.Set("AccessKeyId", h.API.Credentials.Key)
	values.Set("SignatureMethod", "HmacSHA256")
	values.Set("SignatureVersion", "2")
	values.Set("Timestamp", now.UTC().Format("2006-01-02T15:04:05"))
	sigPath := fmt.Sprintf("%s\napi.hbdm.com\n/%s\n%s",
		method, endpoint, values.Encode())
	headers := make(map[string]string)
	if method == http.MethodGet {
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	} else {
		headers["Content-Type"] = "application/json"
	}
	hmac := crypto.GetHMAC(crypto.HashSHA256, []byte(sigPath), []byte(h.API.Credentials.Secret))
	sigValues := url.Values{}
	sigValues.Add("Signature", crypto.Base64Encode(hmac))
	urlPath :=
		common.EncodeURLValues(epoint, values) + "&" + sigValues.Encode()

	var body io.Reader
	var payload []byte
	if data != nil {
		payload, err = json.Marshal(data)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(payload)
	}

	var tempResp json.RawMessage
	var errCap errorCapture

	ctx, cancel := context.WithDeadline(context.Background(), now.Add(15*time.Second))
	defer cancel()
	if err := h.SendPayload(ctx, &request.Item{
		Method:        method,
		Path:          urlPath,
		Headers:       headers,
		Body:          body,
		Result:        &tempResp,
		AuthRequest:   true,
		Verbose:       h.Verbose,
		HTTPDebugging: h.HTTPDebugging,
		HTTPRecording: h.HTTPRecording,
	}); err != nil {
		return err
	}

	if err := json.Unmarshal(tempResp, &errCap); err == nil {
		if errCap.Code != 200 && errCap.ErrMsg != "" {
			return errors.New(errCap.ErrMsg)
		}
	}
	return json.Unmarshal(tempResp, result)
}

// SendAuthenticatedHTTPRequest sends authenticated requests to the HUOBI API
func (h *HUOBI) SendAuthenticatedHTTPRequest(ep exchange.URL, method, endpoint string, values url.Values, data, result interface{}, isVersion2API bool) error {
	if !h.AllowAuthenticatedRequest() {
		return fmt.Errorf(exchange.WarningAuthenticatedRequestWithoutCredentialsSet, h.Name)
	}
	epoint, err := h.API.Endpoints.GetURL(ep)
	if err != nil {
		return err
	}
	if values == nil {
		values = url.Values{}
	}

	now := time.Now()
	values.Set("AccessKeyId", h.API.Credentials.Key)
	values.Set("SignatureMethod", "HmacSHA256")
	values.Set("SignatureVersion", "2")
	values.Set("Timestamp", now.UTC().Format("2006-01-02T15:04:05"))

	if isVersion2API {
		endpoint = fmt.Sprintf("/v%s/%s", huobiAPIVersion2, endpoint)
	} else {
		endpoint = fmt.Sprintf("/v%s/%s", huobiAPIVersion, endpoint)
	}

	payload := fmt.Sprintf("%s\napi.huobi.pro\n%s\n%s",
		method, endpoint, values.Encode())

	headers := make(map[string]string)

	if method == http.MethodGet {
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	} else {
		headers["Content-Type"] = "application/json"
	}

	hmac := crypto.GetHMAC(crypto.HashSHA256, []byte(payload), []byte(h.API.Credentials.Secret))
	values.Set("Signature", crypto.Base64Encode(hmac))
	urlPath := epoint + common.EncodeURLValues(endpoint, values)

	var body []byte
	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}

	// Time difference between your timestamp and standard should be less than 1 minute.
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(time.Minute))
	defer cancel()
	interim := json.RawMessage{}
	err = h.SendPayload(ctx, &request.Item{
		Method:        method,
		Path:          urlPath,
		Headers:       headers,
		Body:          bytes.NewReader(body),
		Result:        &interim,
		AuthRequest:   true,
		Verbose:       h.Verbose,
		HTTPDebugging: h.HTTPDebugging,
		HTTPRecording: h.HTTPRecording,
	})
	if err != nil {
		return err
	}

	if isVersion2API {
		var errCap ResponseV2
		if err = json.Unmarshal(interim, &errCap); err == nil {
			if errCap.Code != 200 && errCap.Message != "" {
				return errors.New(errCap.Message)
			}
		}
	} else {
		var errCap Response
		if err = json.Unmarshal(interim, &errCap); err == nil {
			if errCap.Status == huobiStatusError && errCap.ErrorMessage != "" {
				return errors.New(errCap.ErrorMessage)
			}
		}
	}
	return json.Unmarshal(interim, result)
}

// GetFee returns an estimate of fee based on type of transaction
func (h *HUOBI) GetFee(feeBuilder *exchange.FeeBuilder) (float64, error) {
	var fee float64
	if feeBuilder.FeeType == exchange.OfflineTradeFee || feeBuilder.FeeType == exchange.CryptocurrencyTradeFee {
		fee = calculateTradingFee(feeBuilder.Pair, feeBuilder.PurchasePrice, feeBuilder.Amount)
	}
	if fee < 0 {
		fee = 0
	}

	return fee, nil
}

func calculateTradingFee(c currency.Pair, price, amount float64) float64 {
	if c.IsCryptoFiatPair() {
		return 0.001 * price * amount
	}
	return 0.002 * price * amount
}
