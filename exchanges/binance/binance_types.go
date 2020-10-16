package binance

import (
	"time"

	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
)

// USDT Margined Futures

// UOBData stores ob data for umargined futures
type UOBData struct {
	LastUpdateID int64      `json:"lastUpdateID"`
	Timestamp    int64      `json:"T"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// UPublicTradesData stores trade data
type UPublicTradesData struct {
	ID           int64   `json:"id"`
	Price        float64 `json:"price,string"`
	Qty          float64 `json:"qty,string"`
	QuoteQty     float64 `json:"quoteQty,string"`
	Time         int64   `json:"time"`
	IsBuyerMaker bool    `json:"isBuyerMaker"`
}

// UCompressedTradeData stores compressed trade data
type UCompressedTradeData struct {
	AggregateTradeID int64   `json:"a"`
	Price            float64 `json:"p,string"`
	Quantity         float64 `json:"q,string"`
	FirstTradeID     int64   `json:"f"`
	LastTradeID      int64   `json:"l"`
	Timestamp        int64   `json:"t"`
	IsBuyerMaker     bool    `json:"m"`
}

// UMarkPrice stores mark price data
type UMarkPrice struct {
	Symbol          string  `json:"symbol"`
	MarkPrice       float64 `json:"markPrice,string"`
	IndexPrice      float64 `json:"indexPrice,string"`
	LastFundingRate float64 `json:"lastFundingRate,string"`
	NextFundingTime int64   `json:"nextFundingTime"`
	Time            int64   `json:"time"`
}

// UFundingRateHistory stores funding rate history
type UFundingRateHistory struct {
	Symbol      string  `json:"symbol"`
	FundingRate float64 `json:"fundingRate,string"`
	FundingTime int64   `json:"fundingTime"`
}

// U24HrPriceChangeStats stores price change stats data
type U24HrPriceChangeStats struct {
	Symbol             string  `json:"symbol"`
	PriceChange        float64 `json:"priceChange,string"`
	PriceChangePercent float64 `json:"priceChangePercent,string"`
	WeightedAvgPrice   float64 `json:"weightedAvgPrice,string"`
	PrevClosePrice     float64 `json:"prevClosePrice,string"`
	LastPrice          float64 `json:"lastPrice,string"`
	LastQty            float64 `json:"lastQty,string"`
	OpenPrice          float64 `json:"openPrice,string"`
	HighPrice          float64 `json:"highPrice,string"`
	LowPrice           float64 `json:"lowPrice,string"`
	Volume             float64 `json:"volume,string"`
	QuoteVolume        float64 `json:"quoteVolume,string"`
	OpenTime           int64   `json:"openTime"`
	CloseTime          int64   `json:"closeTime"`
	FirstID            int64   `json:"firstId"`
	LastID             int64   `json:"lastId"`
	Count              int64   `json:"count"`
}

// USymbolPriceTicker stores symbol price ticker data
type USymbolPriceTicker struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
	Time   int64   `json:"time"`
}

// USymbolOrderbookTicker stores symbol orderbook ticker data
type USymbolOrderbookTicker struct {
	Symbol   string  `json:"symbol"`
	BidPrice float64 `json:"bidPrice,string"`
	BidQty   float64 `json:"bidQty,string"`
	AskPrice float64 `json:"askPrice,string"`
	AskQty   float64 `json:"askQty,string"`
	Time     int64   `json:"time"`
}

// ULiquidationOrdersData stores liquidation orders data
type ULiquidationOrdersData struct {
	Symbol       string  `json:"symbol"`
	Price        float64 `json:"price,string"`
	OrigQty      float64 `json:"origQty,string"`
	ExecutedQty  float64 `json:"executedQty,string"`
	AveragePrice float64 `json:"averagePrice,string"`
	Status       string  `json:"status"`
	TimeInForce  string  `json:"timeInForce"`
	OrderType    string  `json:"type"`
	Side         string  `json:"side"`
	Time         int64   `json:"time"`
}

// UOpenInterestData stores open interest data
type UOpenInterestData struct {
	OpenInterest float64 `json:"openInterest,string"`
	Symbol       string  `json:"symbol"`
	Time         int64   `json:"time"`
}

// UOpenInterestStats stores open interest stats data
type UOpenInterestStats struct {
	Symbol               string  `json:"symbol"`
	SumOpenInterest      float64 `json:"sumOpenInterest,string"`
	SumOpenInterestValue float64 `json:"sumOpenInterestValue,string"`
	Timestamp            int64   `json:"timestamp"`
}

// ULongShortRatio stores top trader accounts' or positions' or global long/short ratio data
type ULongShortRatio struct {
	Symbol         string  `json:"symbol"`
	LongShortRatio float64 `json:"longShortRatio,string"`
	LongAccount    float64 `json:"longAccount,string"`
	ShortAccount   float64 `json:"shortAccount,string"`
	Timestamp      int64   `json:"timestamp"`
}

// UTakerVolumeData stores volume data on buy/sell side from takers
type UTakerVolumeData struct {
	BuySellRatio float64 `json:"buySellRatio,string"`
	BuyVol       float64 `json:"buyVol,string"`
	SellVol      float64 `json:"sellVol,string"`
	Timestamp    int64   `json:"timestamp"`
}

// UOrderData stores order data
type UOrderData struct {
	ClientOrderID string  `json:"clientOrderId"`
	CumQty        float64 `json:"cumQty,string"`
	CumQuote      float64 `json:"cumQuote,string"`
	ExecutedQty   float64 `json:"executedQty,string"`
	OrderID       int64   `json:"orderId"`
	AvgPrice      float64 `json:"avgPrice,string"`
	OrigQty       float64 `json:"origQty,string"`
	Price         float64 `json:"price,string"`
	ReduceOnly    bool    `json:"reduceOnly"`
	Side          string  `json:"side"`
	PositionSide  string  `json:"positionSide"`
	Status        string  `json:"status"`
	StopPrice     float64 `json:"stopPrice,string"`
	ClosePosition bool    `json:"closePosition"`
	Symbol        string  `json:"symbol"`
	TimeInForce   string  `json:"timeInForce"`
	OrderType     string  `json:"type"`
	OrigType      string  `json:"origType"`
	ActivatePrice float64 `json:"activatePrice,string"`
	PriceRate     float64 `json:"priceRate,string"`
	UpdateTime    int64   `json:"updateTime"`
	WorkingType   string  `json:"workingType"`
	Code          int64   `json:"code"`
	Msg           string  `json:"msg"`
}

// UFuturesOrderData stores order data for ufutures
type UFuturesOrderData struct {
	AvgPrice      float64 `json:"avgPrice,string"`
	ClientOrderID string  `json:"clientOrderId"`
	CumQuote      string  `json:"cumQuote"`
	ExecutedQty   float64 `json:"executedQty,string"`
	OrderID       int64   `json:"orderId"`
	OrigQty       float64 `json:"origQty,string"`
	OrigType      string  `json:"origType"`
	Price         float64 `json:"price,string"`
	ReduceOnly    bool    `json:"reduceOnly"`
	Side          string  `json:"side"`
	PositionSide  string  `json:"positionSide"`
	Status        string  `json:"status"`
	StopPrice     float64 `json:"stopPrice,string"`
	ClosePosition bool    `json:"closePosition"`
	Symbol        string  `json:"symbol"`
	Time          int64   `json:"time"`
	TimeInForce   string  `json:"timeInForce"`
	OrderType     string  `json:"type"`
	ActivatePrice float64 `json:"activatePrice,string"`
	PriceRate     float64 `json:"priceRate,string"`
	UpdateTime    int64   `json:"updateTime"`
	WorkingType   string  `json:"workingType"`
}

// UAccountBalanceV2Data stores account balance data for ufutures
type UAccountBalanceV2Data struct {
	AccountAlias       string  `json:"accountAlias"`
	Asset              string  `json:"asset"`
	Balance            float64 `json:"balance,string"`
	CrossWalletBalance float64 `json:"crossWalletBalance,string"`
	CrossUnrealizedPNL float64 `json:"crossUnPnl,string"`
	AvailableBalance   float64 `json:"availableBalance,string"`
	MaxWithdrawAmount  float64 `json:"maxWithdrawAmount,string"`
}

// UAccountInformationV2Data stores account info for ufutures
type UAccountInformationV2Data struct {
	FeeTier                     int64   `json:"feeTier"`
	CanTrade                    bool    `json:"canTrade"`
	CanDeposit                  bool    `json:"canDeposit"`
	CanWithdraw                 bool    `json:"canWithdraw"`
	UpdateTime                  int64   `json:"updateTime"`
	TotalInitialMargin          float64 `json:"totalInitialMargin,string"`
	TotalMaintenance            float64 `json:"totalMaintMargin,string"`
	TotalWalletBalance          float64 `json:"totalWalletBalance,string"`
	TotalUnrealizedProfit       float64 `json:"totalUnrealizedProfit,string"`
	TotalMarginBalance          float64 `json:"totalMarginBalance,string"`
	TotalPositionInitialMargin  float64 `json:"totalPositionInitialMargin,string"`
	TotalOpenOrderInitialMargin float64 `json:"totalOpenOrderInitialMargin,string"`
	TotalCrossWalletBalance     float64 `json:"totalCrossWalletBalance,string"`
	TotalCrossUnrealizedPNL     float64 `json:"totalCrossUnPnl,string"`
	AvailableBalance            float64 `json:"availableBalance,string"`
	MaxWithdrawAmount           float64 `json:"maxWithdrawAmount,string"`
	Assets                      []struct {
		Asset                  string  `json:"asset"`
		WalletBalance          float64 `json:"walletBalance,string"`
		UnrealizedProfit       float64 `json:"unrealizedProfit,string"`
		MarginBalance          float64 `json:"marginBalance,string"`
		MaintMargin            float64 `json:"maintMargin,string"`
		InitialMargin          float64 `json:"initialMargin,string"`
		PositionInitialMargin  float64 `json:"positionInitialMargin,string"`
		OpenOrderInitialMargin float64 `json:"openOrderInitialMargin,string"`
		CrossWalletBalance     float64 `json:"crossWalletBalance,string"`
		CrossUnPnl             float64 `json:"crossUnPnl,string"`
		AvailableBalance       float64 `json:"availableBalance,string"`
		MaxWithdrawAmount      float64 `json:"maxWithdrawAmount,string"`
	} `json:"assets"`
	Positions []struct {
		Symbol                 string  `json:"symbol"`
		InitialMargin          float64 `json:"initialMargin,string"`
		MaintenanceMargin      float64 `json:"maintMargin,string"`
		UnrealizedProfit       float64 `json:"unrealizedProfit,string"`
		PositionInitialMargin  float64 `json:"positionInitialMargin,string"`
		OpenOrderInitialMargin float64 `json:"openOrderInitialMargin,string"`
		Leverage               float64 `json:"leverage,string"`
		Isolated               bool    `json:"isolated"`
		EntryPrice             float64 `json:"entryPrice,string"`
		MaxNotional            float64 `json:"maxNotional,string"`
		PositionSide           string  `json:"positionSide"`
	} `json:"positions"`
}

// UChangeInitialLeverage stores leverage change data
type UChangeInitialLeverage struct {
	Leverage         int64   `json:"leverage"`
	MaxNotionalValue float64 `json:"maxNotionalValue,string"`
	Symbol           string  `json:"symbol"`
}

// UModifyIsolatedPosMargin stores modified isolated margin positions' data
type UModifyIsolatedPosMargin struct {
	Amount     float64 `json:"amount,string"`
	MarginType int64   `json:"type"`
}

// UPositionMarginChangeHistoryData gets position margin change history data
type UPositionMarginChangeHistoryData struct {
	Amount       float64 `json:"amount,string"`
	Asset        string  `json:"asset"`
	Symbol       string  `json:"symbol"`
	Time         int64   `json:"time"`
	MarginType   int64   `json:"type"`
	PositionSide string  `json:"positionSide"`
}

// UPositionInformationV2 stores positions' data
type UPositionInformationV2 struct {
	EntryPrice           float64 `json:"entryPrice,string"`
	MarginType           string  `json:"marginType"`
	AutoAddMarginEnabled bool    `json:"isAutoAddMargin"`
	IsolatedMargin       float64 `json:"isolatedMargin,string"`
	Leverage             float64 `json:"leverage,string"`
	LiquidationPrice     float64 `json:"liquidationPrice,string"`
	MarkPrice            float64 `json:"markPrice,string"`
	MaxNotionalValue     float64 `json:"maxNotionalValue,string"`
	PositionAmount       float64 `json:"positionAmt,string"`
	Symbol               string  `json:"symbol"`
	UnrealizedProfit     float64 `json:"unrealizedProfit,string"`
	PositionSide         string  `json:"positionSide"`
}

// UAccountTradeHistory stores trade data for the users account
type UAccountTradeHistory struct {
	Buyer           bool    `json:"buyer"`
	Commission      float64 `json:"commission,string"`
	CommissionAsset string  `json:"commissionAsset"`
	ID              int64   `json:"id"`
	Maker           bool    `json:"maker"`
	OrderID         int64   `json:"orderId"`
	Price           float64 `json:"price,string"`
	Qty             float64 `json:"qty,string"`
	QuoteQty        float64 `json:"quoteQty"`
	RealizedPNL     float64 `json:"realizedPnl,string"`
	Side            string  `json:"side"`
	PositionSide    string  `json:"positionSide"`
	Symbol          string  `json:"symbol"`
	Time            int64   `json:"time"`
}

// UAccountIncomeHistory stores income history data
type UAccountIncomeHistory struct {
	Symbol     string  `json:"symbol"`
	IncomeType string  `json:"incomeType"`
	Income     float64 `json:"income,string"`
	Asset      string  `json:"asset"`
	Info       string  `json:"info"`
	Time       int64   `json:"time"`
	TranID     int64   `json:"tranId"`
	TradeID    string  `json:"tradeId"`
}

// UNotionalLeverageAndBrakcetsData stores notional and leverage brackets data for the account
type UNotionalLeverageAndBrakcetsData struct {
	Symbol   string `json:"symbol"`
	Brackets []struct {
		Bracket                int64   `json:"bracket"`
		InitialLeverage        float64 `json:"initialLeverage,string"`
		NotionalCap            float64 `json:"notionalCap,string"`
		NotionalFloor          float64 `json:"notionalFloor,string"`
		MaintenanceMarginRatio float64 `json:"maintMarginRatio,string"`
	} `json:"brackets"`
}

// UPositionADLEstimationData stores ADL estimation data for a position
type UPositionADLEstimationData struct {
	Symbol      string `json:"symbol"`
	ADLQuantile struct {
		Long  int64 `json:"LONG"`
		Short int64 `json:"SHORT"`
		Hedge int64 `json:"HEDGE"`
	} `json:"adlQuantile"`
}

// UForceOrdersData stores liquidation orders data for the account
type UForceOrdersData struct {
	OrderID       int64   `json:"orderId"`
	Symbol        string  `json:"symbol"`
	Status        string  `json:"status"`
	ClientOrderID string  `json:"clientOrderId"`
	Price         float64 `json:"price,string"`
	AvgPrice      float64 `json:"avgPrice,string"`
	OrigQty       float64 `json:"origQty,string"`
	ExecutedQty   float64 `json:"executedQty,string"`
	CumQuote      float64 `json:"cumQuote,string"`
	TimeInForce   string  `json:"timeInForce"`
	OrderType     string  `json:"type"`
	ReduceOnly    bool    `json:"reduceOnly"`
	ClosePosition bool    `json:"closePosition"`
	Side          string  `json:"side"`
	PositionSide  string  `json:"positionSide"`
	StopPrice     float64 `json:"stopPrice,string"`
	WorkingType   string  `json:"workingType"`
	PriceProtect  bool    `json:"priceProtect,string"`
	OrigType      string  `json:"origType"`
	Time          int64   `json:"time"`
	UpdateTime    int64   `json:"updateTime"`
}

// Coin Margined Futures

// Response holds basic binance api response data
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// FuturesOBData stores orderbook data for futures
type FuturesOBData struct {
	Time int64      `json:"T"`
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

// FuturesPublicTradesData stores recent public trades for futures
type FuturesPublicTradesData struct {
	ID           int64   `json:"id"`
	Price        float64 `json:"price,string"`
	Qty          float64 `json:"qty,string"`
	QuoteQty     float64 `json:"quoteQty,string"`
	Time         int64   `json:"time"`
	IsBuyerMaker bool    `json:"isBuyerMaker"`
}

// CompressedTradesData stores futures trades data in a compressed format
type CompressedTradesData struct {
	TradeID      int64   `json:"a"`
	Price        float64 `json:"p"`
	Quantity     float64 `json:"q"`
	FirstTradeID int64   `json:"f"`
	LastTradeID  int64   `json:"l"`
	Timestamp    int64   `json:"t"`
	BuyerMaker   bool    `json:"b"`
}

// MarkPriceData stores mark price data for futures
type MarkPriceData struct {
	Symbol          string  `json:"symbol"`
	MarkPrice       float64 `json:"markPrice"`
	LastFundingRate float64 `json:"lastFundingRate"`
	NextFundingTime int64   `json:"nextFundingTime"`
	Time            int64   `json:"time"`
}

// SymbolPriceTicker stores ticker price stats
type SymbolPriceTicker struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
	Time   int64   `json:"time"`
}

// SymbolOrderBookTicker stores orderbook ticker data
type SymbolOrderBookTicker struct {
	Symbol   string  `json:"symbol"`
	BidPrice float64 `json:"bidPrice,string"`
	AskPrice float64 `json:"askPrice,string"`
	BidQty   float64 `json:"bidQty,string"`
	AskQty   float64 `json:"askQty,string"`
	Time     int64   `json:"time"`
}

// FuturesCandleStick holds kline data
type FuturesCandleStick struct {
	OpenTime                time.Time
	Open                    float64
	High                    float64
	Low                     float64
	Close                   float64
	Volume                  float64
	CloseTime               time.Time
	BaseAssetVolume         float64
	NumberOfTrades          int64
	TakerBuyVolume          float64
	TakerBuyBaseAssetVolume float64
}

// AllLiquidationOrders gets all liquidation orders
type AllLiquidationOrders struct {
	Symbol       string  `json:"symbol"`
	Price        float64 `json:"price,string"`
	OrigQty      float64 `json:"origQty,string"`
	ExecutedQty  float64 `json:"executedQty,string"`
	AveragePrice float64 `json:"averagePrice,string"`
	Status       string  `json:"status"`
	TimeInForce  string  `json:"timeInForce"`
	OrderType    string  `json:"type"`
	Side         string  `json:"side"`
	Time         int64   `json:"time"`
}

// OpenInterestData stores open interest data
type OpenInterestData struct {
	Symbol       string  `json:"symbol"`
	Pair         string  `json:"pair"`
	OpenInterest float64 `json:"openInterest,string"`
	ContractType string  `json:"contractType"`
	Time         int64   `json:"time"`
}

// OpenInterestStats stores stats for open interest data
type OpenInterestStats struct {
	Pair                 string  `json:"pair"`
	ContractType         string  `json:"contractType"`
	SumOpenInterest      float64 `json:"sumOpenInterest,string"`
	SumOpenInterestValue float64 `json:"sumOpenInterestValue,string"`
	Timestamp            int64   `json:"timestamp"`
}

// TopTraderAccountRatio stores account ratio data for top traders
type TopTraderAccountRatio struct {
	Pair           string  `json:"pair"`
	LongShortRatio float64 `json:"longShortRatio,string"`
	LongAccount    float64 `json:"longAccount,string"`
	ShortAccount   float64 `json:"shortAccount,string"`
	Timestamp      int64   `json:"timestamp"`
}

// TopTraderPositionRatio stores positons' ratio for top trader accounts
type TopTraderPositionRatio struct {
	Pair           string  `json:"pair"`
	LongShortRatio float64 `json:"longShortRatio,string"`
	LongPosition   float64 `json:"longPosition,string"`
	ShortPosition  float64 `json:"shortPosition,string"`
	Timestamp      int64   `json:"timestamp"`
}

// GlobalLongShortRatio stores ratio data of all longs vs shorts
type GlobalLongShortRatio struct {
	Symbol         string  `json:"symbol"`
	LongShortRatio float64 `json:"longShortRatio"`
	LongAccount    float64 `json:"longAccount"`
	ShortAccount   float64 `json:"shortAccount"`
	Timestamp      string  `json:"timestamp"`
}

// TakerBuySellVolume stores taker buy sell volume
type TakerBuySellVolume struct {
	Pair           string  `json:"pair"`
	ContractType   string  `json:"contractType"`
	TakerBuyVolume float64 `json:"takerBuyVol,string"`
	BuySellRatio   float64 `json:"takerSellVol,string"`
	BuyVol         float64 `json:"takerBuyVolValue,string"`
	SellVol        float64 `json:"takerSellVolValue,string"`
	Timestamp      int64   `json:"timestamp"`
}

// FuturesBasisData gets futures basis data
type FuturesBasisData struct {
	Pair         string  `json:"pair"`
	ContractType string  `json:"contractType"`
	FuturesPrice float64 `json:"futuresPrice,string"`
	IndexPrice   float64 `json:"indexPrice,string"`
	Basis        float64 `json:"basis,string,string"`
	BasisRate    float64 `json:"basisRate,string,string"`
	Timestamp    int64   `json:"timestamp"`
}

// PlaceBatchOrderData stores batch order data for placing
type PlaceBatchOrderData struct {
	Symbol           string  `json:"symbol"`
	Side             string  `json:"side"`
	PositionSide     string  `json:"positionSide,omitempty"`
	OrderType        string  `json:"type"`
	TimeInForce      string  `json:"timeInForce,omitempty"`
	Quantity         float64 `json:"quantity"`
	ReduceOnly       string  `json:"reduceOnly,omitempty"`
	Price            float64 `json:"price"`
	NewClientOrderID string  `json:"newClientOrderId,omitempty"`
	StopPrice        float64 `json:"stopPrice,omitempty"`
	ActivationPrice  float64 `json:"activationPrice,omitempty"`
	CallbackRate     float64 `json:"callbackRate,omitempty"`
	WorkingType      string  `json:"workingType,omitempty"`
	PriceProtect     string  `json:"priceProtect,omitempty"`
	NewOrderRespType string  `json:"newOrderRespType,omitempty"`
}

// BatchCancelOrderData stores batch cancel order data
type BatchCancelOrderData struct {
	ClientOrderID string  `json:"clientOrderID"`
	CumQty        float64 `json:"cumQty,string"`
	CumBase       float64 `json:"cumBase,string"`
	ExecuteQty    float64 `json:"executeQty,string"`
	OrderID       int64   `json:"orderID,string"`
	AvgPrice      float64 `json:"avgPrice,string"`
	OrigQty       float64 `json:"origQty,string"`
	Price         float64 `json:"price,string"`
	ReduceOnly    bool    `json:"reduceOnly"`
	Side          string  `json:"side"`
	PositionSide  string  `json:"positionSide"`
	Status        string  `json:"status"`
	StopPrice     int64   `json:"stopPrice"`
	ClosePosition bool    `json:"closePosition"`
	Symbol        string  `json:"symbol"`
	Pair          string  `json:"pair"`
	TimeInForce   string  `json:"TimeInForce"`
	OrderType     string  `json:"type"`
	OrigType      string  `json:"origType"`
	ActivatePrice float64 `json:"activatePrice,string"`
	PriceRate     float64 `json:"priceRate,string"`
	UpdateTime    int64   `json:"updateTime"`
	WorkingType   string  `json:"workingType"`
	PriceProtect  bool    `json:"priceProtect"`
	Code          int64   `json:"code"`
	Msg           string  `json:"msg"`
}

// FuturesOrderPlaceData stores futures order data
type FuturesOrderPlaceData struct {
	ClientOrderID string  `json:"clientOrderID"`
	CumQty        float64 `json:"cumQty,string"`
	CumBase       float64 `json:"cumBase,string"`
	ExecuteQty    float64 `json:"executeQty,string"`
	OrderID       int64   `json:"orderID,string"`
	AvgPrice      float64 `json:"avgPrice,string"`
	OrigQty       float64 `json:"origQty,string"`
	Price         float64 `json:"price,string"`
	ReduceOnly    bool    `json:"reduceOnly"`
	Side          string  `json:"side"`
	PositionSide  string  `json:"positionSide"`
	Status        string  `json:"status"`
	StopPrice     int64   `json:"stopPrice"`
	ClosePosition bool    `json:"closePosition"`
	Symbol        string  `json:"symbol"`
	Pair          string  `json:"pair"`
	TimeInForce   string  `json:"TimeInForce"`
	OrderType     string  `json:"type"`
	OrigType      string  `json:"origType"`
	ActivatePrice float64 `json:"activatePrice,string"`
	PriceRate     float64 `json:"priceRate,string"`
	UpdateTime    int64   `json:"updateTime"`
	WorkingType   string  `json:"workingType"`
	PriceProtect  bool    `json:"priceProtect"`
}

// FuturesOrderGetData stores futures order data for get requests
type FuturesOrderGetData struct {
	AvgPrice      float64 `json:"avgPrice,string"`
	ClientOrderID string  `json:"clientOrderID"`
	CumQty        float64 `json:"cumQty,string"`
	CumBase       float64 `json:"cumBase,string"`
	ExecutedQty   float64 `json:"executedQty,string"`
	OrderID       int64   `json:"orderId"`
	OrigQty       float64 `json:"origQty,string"`
	OrigType      string  `json:"origType"`
	Price         float64 `json:"price,string"`
	ReduceOnly    bool    `json:"reduceOnly"`
	Side          string  `json:"buy"`
	PositionSide  string  `json:"positionSide"`
	Status        string  `json:"status"`
	StopPrice     float64 `json:"stopPrice,string"`
	ClosePosition bool    `json:"closePosition"`
	Symbol        string  `json:"symbol"`
	Pair          string  `json:"pair"`
	TimeInForce   string  `json:"timeInForce"`
	OrderType     string  `json:"type"`
	ActivatePrice float64 `json:"activatePrice,string"`
	PriceRate     float64 `json:"priceRate,string"`
	UpdateTime    int64   `json:"updateTime"`
	WorkingType   string  `json:"workingType"`
	PriceProtect  bool    `json:"priceProtect"`
}

// FuturesOrderData stores order data for futures
type FuturesOrderData struct {
	AvgPrice      float64 `json:"avgPrice,string"`
	ClientOrderID string  `json:"clientOrderId"`
	CumBase       string  `json:"cumBase"`
	ExecutedQty   float64 `json:"executedQty,string"`
	OrderID       int64   `json:"orderId"`
	OrigQty       float64 `json:"origQty,string"`
	OrigType      string  `json:"origType"`
	Price         float64 `json:"price,string"`
	ReduceOnly    bool    `json:"reduceOnly"`
	Side          string  `json:"side"`
	PositionSide  string  `json:"positionSide"`
	Status        string  `json:"status"`
	StopPrice     float64 `json:"stopPrice,string"`
	ClosePosition bool    `json:"closePosition"`
	Symbol        string  `json:"symbol"`
	Pair          string  `json:"pair"`
	Time          int64   `json:"time"`
	TimeInForce   string  `json:"timeInForce"`
	OrderType     string  `json:"type"`
	ActivatePrice float64 `json:"activatePrice,string"`
	PriceRate     float64 `json:"priceRate,string"`
	UpdateTime    int64   `json:"updateTime"`
	WorkingType   string  `json:"workingType"`
	PriceProtect  bool    `json:"priceProtect"`
}

// OrderVars stores side, status and type for any order/trade
type OrderVars struct {
	Side      order.Side
	Status    order.Status
	OrderType order.Type
	Fee       float64
}

// AutoCancelAllOrdersData gives data of auto cancelling all open orders
type AutoCancelAllOrdersData struct {
	Symbol        string `json:"symbol"`
	CountdownTime int64  `json:"countdownTime,string"`
}

// LevelDetail stores level detail data
type LevelDetail struct {
	Level         string  `json:"level"`
	MaxBorrowable float64 `json:"maxBorrowable,string"`
	InterestRate  float64 `json:"interestRate,string"`
}

// MarginInfoData stores margin info data
type MarginInfoData struct {
	Data []struct {
		MarginRatio string `json:"marginRatio"`
		Base        struct {
			AssetName    string        `json:"assetName"`
			LevelDetails []LevelDetail `json:"levelDetails"`
		} `json:"base"`
		Quote struct {
			AssetName    string        `json:"assetName"`
			LevelDetails []LevelDetail `json:"levelDetails"`
		} `json:"quote"`
	} `json:"data"`
}

// FuturesAccountBalanceData stores account balance data for futures
type FuturesAccountBalanceData struct {
	AccountAlias       string  `json:"accountAlias"`
	Asset              string  `json:"asset"`
	Balance            float64 `json:"balance,string"`
	WithdrawAvailable  float64 `json:"withdrawAvailable,string"`
	CrossWalletBalance float64 `json:"crossWalletBalance,string"`
	CrossUnPNL         float64 `json:"crossUnPNL,string"`
	AvailableBalance   float64 `json:"availableBalance,string"`
	UpdateTime         int64   `json:"updateTime"`
}

// FuturesAccountInformation stores account information for futures account
type FuturesAccountInformation struct {
	Assets []struct {
		Asset                  string  `json:"asset"`
		WalletBalance          float64 `json:"walletBalance,string"`
		UnrealizedProfit       float64 `json:"unrealizedProfit,string"`
		MarginBalance          float64 `json:"marginBalance,string"`
		MaintMargin            float64 `json:"maintMargin,string"`
		InitialMargin          float64 `json:"initialMargin,string"`
		PositionInitialMargin  float64 `json:"positionInitialMargin,string"`
		OpenOrderInitialMargin float64 `json:"openOrderInitialMargin,string"`
		Leverage               float64 `json:"leverage,string"`
		Isolated               bool    `json:"isolated"`
		PositionSide           string  `json:"positionSide"`
		EntryPrice             float64 `json:"entryPrice,string"`
		MaxQty                 float64 `json:"maxQty,string"`
	} `json:"assets"`
	Positions []struct {
		Symbol                 string  `json:"symbol"`
		InitialMargin          float64 `json:"initialMargin,string"`
		MaintMargin            float64 `json:"maintMargin,string"`
		UnrealizedProfit       float64 `json:"unrealizedProfit,string"`
		PositionInitialMargin  float64 `json:"positionInitialMargin,string"`
		OpenOrderInitialMargin float64 `json:"openOrderInitialMargin,string"`
		Leverage               float64 `json:"leverage,string"`
		Isolated               bool    `json:"isolated,false"`
		PositionSide           string  `json:"positionSide"`
		EntryPrice             float64 `json:"entryPrice,string"`
		MaxQty                 float64 `json:"maxQty,string"`
	} `json:"positions"`
	CanDeposit  bool  `json:"canDeposit"`
	CanTrade    bool  `json:"canTrade"`
	CanWithdraw bool  `json:"canWithdraw"`
	FeeTier     int64 `json:"feeTier"`
	UpdateTime  int64 `json:"updateTime"`
}

// GenericAuthResponse is a general data response for a post auth request
type GenericAuthResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

// FuturesLeverageData stores leverage data for futures
type FuturesLeverageData struct {
	Leverage int64   `json:"leverage"`
	MaxQty   float64 `json:"maxQty,string"`
	Symbol   string  `json:"symbol"`
}

// ModifyIsolatedMarginData stores margin modification data
type ModifyIsolatedMarginData struct {
	Amount  float64 `json:"amount"`
	Code    int64   `json:"code"`
	Msg     string  `json:"msg"`
	ModType string  `json:"modType"`
}

// GetPositionMarginChangeHistoryData gets margin change history for positions
type GetPositionMarginChangeHistoryData struct {
	Amount           float64 `json:"amount"`
	Asset            string  `json:"asset"`
	Symbol           string  `json:"symbol"`
	Timestamp        int64   `json:"time"`
	MarginChangeType int64   `json:"type"`
	PositionSide     string  `json:"positionSide"`
}

// FuturesPositionInformation stores futures positon info
type FuturesPositionInformation struct {
	Symbol           string  `json:"symbol"`
	PositionAmount   float64 `json:"positionAmt,string"`
	EntryPrice       float64 `json:"entryPrice,string"`
	MarkPrice        float64 `json:"markPrice,string"`
	UnrealizedProfit float64 `json:"unRealizedProfit,string"`
	LiquidationPrice float64 `json:"liquidation,string"`
	Leverage         int64   `json:"leverage"`
	MaxQty           float64 `json:"maxQty"`
	MarginType       string  `json:"marginType"`
	IsolatedMargin   float64 `json:"isolatedMargin,string"`
	IsAutoAddMargin  bool    `json:"isAutoAddMargin"`
	PositionSide     string  `json:"positionSide"`
}

// FuturesAccountTradeList stores account trade list data
type FuturesAccountTradeList struct {
	Symbol          string  `json:"symbol"`
	ID              int64   `json:"id"`
	OrderID         int64   `json:"orderID"`
	Pair            string  `json:"pair"`
	Side            string  `json:"side"`
	Price           string  `json:"price"`
	Qty             float64 `json:"qty"`
	RealizedPNL     float64 `json:"realizedPNL"`
	MarginAsset     string  `json:"marginAsset"`
	BaseQty         float64 `json:"baseQty"`
	Commission      float64 `json:"commission"`
	CommissionAsset string  `json:"commissionAsset"`
	Timestamp       int64   `json:"timestamp"`
	PositionSide    string  `json:"positionSide"`
	Buyer           bool    `json:"buyer"`
	Maker           bool    `json:"maker"`
}

// FuturesIncomeHistoryData stores futures income history data
type FuturesIncomeHistoryData struct {
	Symbol     string  `json:"symbol"`
	IncomeType string  `json:"incomeType,string"`
	Income     float64 `json:"income,string"`
	Asset      string  `json:"asset"`
	Info       string  `json:"info"`
	Timestamp  int64   `json:"time"`
}

// NotionalBracketData stores notional bracket data
type NotionalBracketData struct {
	Pair     string `json:"pair"`
	Brackets []struct {
		Bracket          int64   `json:"bracket"`
		InitialLeverage  float64 `json:"initialLeverage"`
		QtyCap           float64 `json:"qtyCap"`
		QtylFloor        float64 `json:"qtyFloor"`
		MaintMarginRatio float64 `json:"maintMarginRatio"`
	}
}

// ForcedOrdersData stores forced orders data
type ForcedOrdersData struct {
	OrderID       int64   `json:"orderId"`
	Symbol        string  `json:"symbol"`
	Status        string  `json:"status"`
	ClientOrderID string  `json:"clientOrderId"`
	Price         float64 `json:"price,string"`
	AvgPrice      float64 `json:"avgPrice,string"`
	OrigQty       float64 `json:"origQty,string"`
	ExecutedQty   float64 `json:"executedQty,string"`
	CumQuote      float64 `json:"cumQuote"`
	TimeInForce   string  `json:"timeInForce,string"`
	OrderType     string  `json:"orderType"`
	ReduceOnly    bool    `json:"reduceOnly"`
	ClosePosition bool    `json:"closePosition"`
	Side          string  `json:"side"`
	PositionSide  string  `json:"positionSide"`
	StopPrice     float64 `json:"stopPrice,string"`
	WorkingType   string  `json:"workingType,string"`
	PriceProtect  float64 `json:"priceProtect"`
	OrigType      string  `json:"origType"`
	Time          int64   `json:"time"`
	UpdateTime    int64   `json:"updateTime"`
}

// ADLEstimateData stores data for ADL estimates
type ADLEstimateData struct {
	Symbol      string `json:"symbol"`
	ADLQuantile struct {
		Long  float64 `json:"LONG"`
		Short float64 `json:"SHORT"`
		Hedge float64 `json:"HEDGE"`
	} `json:"adlQuantile"`
}

// InterestHistoryData gets interest history data
type InterestHistoryData struct {
	Asset       string  `json:"asset"`
	Interest    float64 `json:"interest"`
	LendingType string  `json:"lendingType"`
	ProductName string  `json:"productName"`
	Time        string  `json:"time"`
}

// FundingRateData stores funding rates data
type FundingRateData struct {
	Symbol      string  `json:"symbol"`
	FundingRate float64 `json:"fundingRate,string"`
	FundingTime int64   `json:"fundingTime"`
}

// SymbolsData stores perp futures' symbols
type SymbolsData struct {
	Symbol string `json:"symbol"`
}

// PerpsExchangeInfo stores data for perps
type PerpsExchangeInfo struct {
	Symbols []SymbolsData `json:"symbols"`
}

// UFuturesExchangeInfo stores exchange info for ufutures
type UFuturesExchangeInfo struct {
	RateLimits []struct {
		Interval      string `json:"interval"`
		IntervalNum   int64  `json:"intervalNum"`
		Limit         int64  `json:"limit"`
		RateLimitType string `json:"rateLimitType"`
	} `json:"rateLimits"`
	ServerTime int64 `json:"serverTime"`
	Symbols    []struct {
		Symbol                   string  `json:"symbol"`
		Status                   string  `json:"status"`
		MaintenanceMarginPercent float64 `json:"maintMarginPercent,string"`
		RequiredMarginPercent    float64 `json:"requiredMarginPercent,string"`
		BaseAsset                string  `json:"baseAsset"`
		QuoteAsset               string  `json:"quoteAsset"`
		PricePrecision           int64   `json:"pricePrecision"`
		QuantityPrecision        int64   `json:"quantityPrecision"`
		BaseAssetPrecision       int64   `json:"baseAssetPrecision"`
		QuotePrecision           int64   `json:"quotePrecision"`
		Filters                  []struct {
			MinPrice          float64 `json:"minPrice,string,omitempty"`
			MaxPrice          float64 `json:"maxPrice,string,omitempty"`
			FilterType        string  `json:"filterType,omitempty"`
			TickSize          float64 `json:"tickSize,string,omitempty"`
			StepSize          float64 `json:"stepSize,string,omitempty"`
			MaxQty            float64 `json:"maxQty,string,omitempty"`
			MinQty            float64 `json:"minQty,string,omitempty"`
			Limit             int64   `json:"limit,omitempty"`
			MultiplierDown    float64 `json:"multiplierDown,string,omitempty"`
			MultiplierUp      float64 `json:"multiplierUp,string,omitempty"`
			MultiplierDecimal float64 `json:"multiplierDecimal,string,omitempty"`
		} `json:"filters"`
		OrderTypes  []string `json:"orderTypes"`
		TimeInForce []string `json:"timeInForce"`
	} `json:"symbols"`
	Timezone string `json:"timezone"`
}

// CExchangeInfo stores exchange info for cfutures
type CExchangeInfo struct {
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	RateLimits      []struct {
		Interval      string `json:"interval"`
		IntervalNum   int64  `json:"intervalNul"`
		Limit         int64  `json:"limit"`
		RateLimitType string `json:"rateLimitType"`
	} `json:"rateLimits"`
	ServerTime int64 `json:"serverTime"`
	Symbols    []struct {
		Filters []struct {
			FilterType        string  `json:"filterType,omitempty"`
			MinPrice          float64 `json:"minPrice,string,omitempty"`
			MaxPrice          float64 `json:"maxPrice,string,omitempty"`
			StepSize          float64 `json:"stepSize,string,omitempty"`
			MaxQty            float64 `json:"maxQty,string,omitempty"`
			MinQty            float64 `json:"minQty,string,omitempty"`
			Limit             int64   `json:"limit,omitempty"`
			MultiplierDown    float64 `json:"multiplierDown,string,omitempty"`
			MultiplierUp      float64 `json:"multiplierUp,string,omitempty"`
			MultiplierDecimal float64 `json:"multiplierDecimal,string,omitempty"`
		} `json:"filters"`
		OrderTypes            []string `json:"orderType"`
		TimeInForce           []string `json:"timeInForce"`
		Symbol                string   `json:"symbol"`
		Pair                  string   `json:"pair"`
		ContractType          string   `json:"contractType"`
		DeliveryDate          int64    `json:"deliveryDate"`
		OnboardDate           int64    `json:"onboardDate"`
		ContractStatus        string   `json:"contractStatus"`
		ContractSize          int64    `json:"contractSize"`
		QuoteAsset            string   `json:"quoteAsset"`
		BaseAsset             string   `json:"baseAsset"`
		MarginAsset           string   `json:"marginAsset"`
		PricePrecision        int64    `json:"pricePrecision"`
		QuantityPrecision     int64    `json:"quantityPrecision"`
		BaseAssetPrecision    int64    `json:"baseAssetPrecision"`
		QuotePrecision        int64    `json:"quotePrecision"`
		MaintMarginPercent    float64  `json:"maintMarginPercent,string"`
		RequiredMarginPercent float64  `json:"requiredMarginPercent,string"`
	} `json:"symbols"`
	Timezone string `json:"timezone"`
}

// ExchangeInfo holds the full exchange information type
type ExchangeInfo struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	Timezone   string `json:"timezone"`
	Servertime int64  `json:"serverTime"`
	RateLimits []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		Limit         int    `json:"limit"`
	} `json:"rateLimits"`
	ExchangeFilters interface{} `json:"exchangeFilters"`
	Symbols         []struct {
		Symbol                     string   `json:"symbol"`
		Status                     string   `json:"status"`
		BaseAsset                  string   `json:"baseAsset"`
		BaseAssetPrecision         int      `json:"baseAssetPrecision"`
		QuoteAsset                 string   `json:"quoteAsset"`
		QuotePrecision             int      `json:"quotePrecision"`
		OrderTypes                 []string `json:"orderTypes"`
		IcebergAllowed             bool     `json:"icebergAllowed"`
		OCOAllowed                 bool     `json:"ocoAllowed"`
		QuoteOrderQtyMarketAllowed bool     `json:"quoteOrderQtyMarketAllowed"`
		IsSpotTradingAllowed       bool     `json:"isSpotTradingAllowed"`
		IsMarginTradingAllowed     bool     `json:"isMarginTradingAllowed"`
		Filters                    []struct {
			FilterType          string  `json:"filterType"`
			MinPrice            float64 `json:"minPrice,string"`
			MaxPrice            float64 `json:"maxPrice,string"`
			TickSize            float64 `json:"tickSize,string"`
			MultiplierUp        float64 `json:"multiplierUp,string"`
			MultiplierDown      float64 `json:"multiplierDown,string"`
			AvgPriceMins        int64   `json:"avgPriceMins"`
			MinQty              float64 `json:"minQty,string"`
			MaxQty              float64 `json:"maxQty,string"`
			StepSize            float64 `json:"stepSize,string"`
			MinNotional         float64 `json:"minNotional,string"`
			ApplyToMarket       bool    `json:"applyToMarket"`
			Limit               int64   `json:"limit"`
			MaxNumAlgoOrders    int64   `json:"maxNumAlgoOrders"`
			MaxNumIcebergOrders int64   `json:"maxNumIcebergOrders"`
		} `json:"filters"`
	} `json:"symbols"`
}

// OrderBookDataRequestParams represents Klines request data.
type OrderBookDataRequestParams struct {
	Symbol string `json:"symbol"` // Required field; example LTCBTC,BTCUSDT
	Limit  int    `json:"limit"`  // Default 100; max 1000. Valid limits:[5, 10, 20, 50, 100, 500, 1000]
}

// OrderbookItem stores an individual orderbook item
type OrderbookItem struct {
	Price    float64
	Quantity float64
}

// OrderBookData is resp data from orderbook endpoint
type OrderBookData struct {
	Code         int         `json:"code"`
	Msg          string      `json:"msg"`
	LastUpdateID int64       `json:"lastUpdateId"`
	Bids         [][2]string `json:"bids"`
	Asks         [][2]string `json:"asks"`
}

// OrderBook actual structured data that can be used for orderbook
type OrderBook struct {
	Symbol       string
	LastUpdateID int64
	Code         int
	Msg          string
	Bids         []OrderbookItem
	Asks         []OrderbookItem
}

// DepthUpdateParams is used as an embedded type for WebsocketDepthStream
type DepthUpdateParams []struct {
	PriceLevel float64
	Quantity   float64
	ingnore    []interface{}
}

// WebsocketDepthStream is the difference for the update depth stream
type WebsocketDepthStream struct {
	Event         string          `json:"e"`
	Timestamp     int64           `json:"E"`
	Pair          string          `json:"s"`
	FirstUpdateID int64           `json:"U"`
	LastUpdateID  int64           `json:"u"`
	UpdateBids    [][]interface{} `json:"b"`
	UpdateAsks    [][]interface{} `json:"a"`
}

// RecentTradeRequestParams represents Klines request data.
type RecentTradeRequestParams struct {
	Symbol string `json:"symbol"` // Required field. example LTCBTC, BTCUSDT
	Limit  int    `json:"limit"`  // Default 500; max 500.
}

// RecentTrade holds recent trade data
type RecentTrade struct {
	ID           int64   `json:"id"`
	Price        float64 `json:"price,string"`
	Quantity     float64 `json:"qty,string"`
	Time         float64 `json:"time"`
	IsBuyerMaker bool    `json:"isBuyerMaker"`
	IsBestMatch  bool    `json:"isBestMatch"`
}

// TradeStream holds the trade stream data
type TradeStream struct {
	EventType      string `json:"e"`
	EventTime      int64  `json:"E"`
	Symbol         string `json:"s"`
	TradeID        int64  `json:"t"`
	Price          string `json:"p"`
	Quantity       string `json:"q"`
	BuyerOrderID   int64  `json:"b"`
	SellerOrderID  int64  `json:"a"`
	TimeStamp      int64  `json:"T"`
	Maker          bool   `json:"m"`
	BestMatchPrice bool   `json:"M"`
}

// KlineStream holds the kline stream data
type KlineStream struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	Kline     struct {
		StartTime                int64   `json:"t"`
		CloseTime                int64   `json:"T"`
		Symbol                   string  `json:"s"`
		Interval                 string  `json:"i"`
		FirstTradeID             int64   `json:"f"`
		LastTradeID              int64   `json:"L"`
		OpenPrice                float64 `json:"o,string"`
		ClosePrice               float64 `json:"c,string"`
		HighPrice                float64 `json:"h,string"`
		LowPrice                 float64 `json:"l,string"`
		Volume                   float64 `json:"v,string"`
		NumberOfTrades           int64   `json:"n"`
		KlineClosed              bool    `json:"x"`
		Quote                    float64 `json:"q,string"`
		TakerBuyBaseAssetVolume  float64 `json:"V,string"`
		TakerBuyQuoteAssetVolume float64 `json:"Q,string"`
	} `json:"k"`
}

// TickerStream holds the ticker stream data
type TickerStream struct {
	EventType              string  `json:"e"`
	EventTime              int64   `json:"E"`
	Symbol                 string  `json:"s"`
	PriceChange            float64 `json:"p,string"`
	PriceChangePercent     float64 `json:"P,string"`
	WeightedAvgPrice       float64 `json:"w,string"`
	ClosePrice             float64 `json:"x,string"`
	LastPrice              float64 `json:"c,string"`
	LastPriceQuantity      float64 `json:"Q,string"`
	BestBidPrice           float64 `json:"b,string"`
	BestBidQuantity        float64 `json:"B,string"`
	BestAskPrice           float64 `json:"a,string"`
	BestAskQuantity        float64 `json:"A,string"`
	OpenPrice              float64 `json:"o,string"`
	HighPrice              float64 `json:"h,string"`
	LowPrice               float64 `json:"l,string"`
	TotalTradedVolume      float64 `json:"v,string"`
	TotalTradedQuoteVolume float64 `json:"q,string"`
	OpenTime               int64   `json:"O"`
	CloseTime              int64   `json:"C"`
	FirstTradeID           int64   `json:"F"`
	LastTradeID            int64   `json:"L"`
	NumberOfTrades         int64   `json:"n"`
}

// HistoricalTrade holds recent trade data
type HistoricalTrade struct {
	Code         int     `json:"code"`
	Msg          string  `json:"msg"`
	ID           int64   `json:"id"`
	Price        float64 `json:"price,string"`
	Quantity     float64 `json:"qty,string"`
	Time         int64   `json:"time"`
	IsBuyerMaker bool    `json:"isBuyerMaker"`
	IsBestMatch  bool    `json:"isBestMatch"`
}

// AggregatedTrade holds aggregated trade information
type AggregatedTrade struct {
	ATradeID       int64   `json:"a"`
	Price          float64 `json:"p,string"`
	Quantity       float64 `json:"q,string"`
	FirstTradeID   int64   `json:"f"`
	LastTradeID    int64   `json:"l"`
	TimeStamp      int64   `json:"T"`
	BuyerMaker     bool    `json:"m"`
	BestMatchPrice bool    `json:"M"`
}

// IndexMarkPrice stores data for index and mark prices
type IndexMarkPrice struct {
	Symbol               string  `json:"symbol"`
	Pair                 string  `json:"pair"`
	MarkPrice            float64 `json:"markPrice,string"`
	IndexPrice           float64 `json:"indexPrice,string"`
	EstimatedSettlePrice float64 `json:"estimatedSettlePrice,string"`
	LastFundingRate      string  `json:"lastFundingRate"`
	NextFundingTime      int64   `json:"nextFundingTime"`
	Time                 int64   `json:"time"`
}

// CandleStick holds kline data
type CandleStick struct {
	OpenTime                 time.Time
	Open                     float64
	High                     float64
	Low                      float64
	Close                    float64
	Volume                   float64
	CloseTime                time.Time
	QuoteAssetVolume         float64
	TradeCount               float64
	TakerBuyAssetVolume      float64
	TakerBuyQuoteAssetVolume float64
}

// AveragePrice holds current average symbol price
type AveragePrice struct {
	Mins  int64   `json:"mins"`
	Price float64 `json:"price,string"`
}

// PriceChangeStats contains statistics for the last 24 hours trade
type PriceChangeStats struct {
	Symbol             string  `json:"symbol"`
	PriceChange        float64 `json:"priceChange,string"`
	PriceChangePercent float64 `json:"priceChangePercent,string"`
	WeightedAvgPrice   float64 `json:"weightedAvgPrice,string"`
	PrevClosePrice     float64 `json:"prevClosePrice,string"`
	LastPrice          float64 `json:"lastPrice,string"`
	LastQty            float64 `json:"lastQty,string"`
	BidPrice           float64 `json:"bidPrice,string"`
	AskPrice           float64 `json:"askPrice,string"`
	OpenPrice          float64 `json:"openPrice,string"`
	HighPrice          float64 `json:"highPrice,string"`
	LowPrice           float64 `json:"lowPrice,string"`
	Volume             float64 `json:"volume,string"`
	QuoteVolume        float64 `json:"quoteVolume,string"`
	OpenTime           int64   `json:"openTime"`
	CloseTime          int64   `json:"closeTime"`
	FirstID            int64   `json:"firstId"`
	LastID             int64   `json:"lastId"`
	Count              int64   `json:"count"`
}

// SymbolPrice holds basic symbol price
type SymbolPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

// BestPrice holds best price data
type BestPrice struct {
	Symbol   string  `json:"symbol"`
	BidPrice float64 `json:"bidPrice,string"`
	BidQty   float64 `json:"bidQty,string"`
	AskPrice float64 `json:"askPrice,string"`
	AskQty   float64 `json:"askQty,string"`
}

// NewOrderRequest request type
type NewOrderRequest struct {
	// Symbol (currency pair to trade)
	Symbol string
	// Side Buy or Sell
	Side string
	// TradeType (market or limit order)
	TradeType RequestParamsOrderType
	// TimeInForce specifies how long the order remains in effect.
	// Examples are (Good Till Cancel (GTC), Immediate or Cancel (IOC) and Fill Or Kill (FOK))
	TimeInForce RequestParamsTimeForceType
	// Quantity is the total base qty spent or received in an order.
	Quantity float64
	// QuoteOrderQty is the total quote qty spent or received in a MARKET order.
	QuoteOrderQty    float64
	Price            float64
	NewClientOrderID string
	StopPrice        float64 // Used with STOP_LOSS, STOP_LOSS_LIMIT, TAKE_PROFIT, and TAKE_PROFIT_LIMIT orders.
	IcebergQty       float64 // Used with LIMIT, STOP_LOSS_LIMIT, and TAKE_PROFIT_LIMIT to create an iceberg order.
	NewOrderRespType string
}

// NewOrderResponse is the return structured response from the exchange
type NewOrderResponse struct {
	Code            int     `json:"code"`
	Msg             string  `json:"msg"`
	Symbol          string  `json:"symbol"`
	OrderID         int64   `json:"orderId"`
	ClientOrderID   string  `json:"clientOrderId"`
	TransactionTime int64   `json:"transactTime"`
	Price           float64 `json:"price,string"`
	OrigQty         float64 `json:"origQty,string"`
	ExecutedQty     float64 `json:"executedQty,string"`
	// The cumulative amount of the quote that has been spent (with a BUY order) or received (with a SELL order).
	CumulativeQuoteQty float64 `json:"cummulativeQuoteQty,string"`
	Status             string  `json:"status"`
	TimeInForce        string  `json:"timeInForce"`
	Type               string  `json:"type"`
	Side               string  `json:"side"`
	Fills              []struct {
		Price           float64 `json:"price,string"`
		Qty             float64 `json:"qty,string"`
		Commission      float64 `json:"commission,string"`
		CommissionAsset string  `json:"commissionAsset"`
	} `json:"fills"`
}

// CancelOrderResponse is the return structured response from the exchange
type CancelOrderResponse struct {
	Symbol            string `json:"symbol"`
	OrigClientOrderID string `json:"origClientOrderId"`
	OrderID           int64  `json:"orderId"`
	ClientOrderID     string `json:"clientOrderId"`
}

// QueryOrderData holds query order data
type QueryOrderData struct {
	Code          int     `json:"code"`
	Msg           string  `json:"msg"`
	Symbol        string  `json:"symbol"`
	OrderID       int64   `json:"orderId"`
	ClientOrderID string  `json:"clientOrderId"`
	Price         float64 `json:"price,string"`
	OrigQty       float64 `json:"origQty,string"`
	ExecutedQty   float64 `json:"executedQty,string"`
	Status        string  `json:"status"`
	TimeInForce   string  `json:"timeInForce"`
	Type          string  `json:"type"`
	Side          string  `json:"side"`
	StopPrice     float64 `json:"stopPrice,string"`
	IcebergQty    float64 `json:"icebergQty,string"`
	Time          float64 `json:"time"`
	IsWorking     bool    `json:"isWorking"`
}

// Balance holds query order data
type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// Account holds the account data
type Account struct {
	MakerCommission  int       `json:"makerCommission"`
	TakerCommission  int       `json:"takerCommission"`
	BuyerCommission  int       `json:"buyerCommission"`
	SellerCommission int       `json:"sellerCommission"`
	CanTrade         bool      `json:"canTrade"`
	CanWithdraw      bool      `json:"canWithdraw"`
	CanDeposit       bool      `json:"canDeposit"`
	UpdateTime       int64     `json:"updateTime"`
	Balances         []Balance `json:"balances"`
}

// RequestParamsTimeForceType Time in force
type RequestParamsTimeForceType string

var (
	// BinanceRequestParamsTimeGTC GTC
	BinanceRequestParamsTimeGTC = RequestParamsTimeForceType("GTC")

	// BinanceRequestParamsTimeIOC IOC
	BinanceRequestParamsTimeIOC = RequestParamsTimeForceType("IOC")

	// BinanceRequestParamsTimeFOK FOK
	BinanceRequestParamsTimeFOK = RequestParamsTimeForceType("FOK")
)

// RequestParamsOrderType trade order type
type RequestParamsOrderType string

var (
	// BinanceRequestParamsOrderLimit Limit order
	BinanceRequestParamsOrderLimit = RequestParamsOrderType("LIMIT")

	// BinanceRequestParamsOrderMarket Market order
	BinanceRequestParamsOrderMarket = RequestParamsOrderType("MARKET")

	// BinanceRequestParamsOrderStopLoss STOP_LOSS
	BinanceRequestParamsOrderStopLoss = RequestParamsOrderType("STOP_LOSS")

	// BinanceRequestParamsOrderStopLossLimit STOP_LOSS_LIMIT
	BinanceRequestParamsOrderStopLossLimit = RequestParamsOrderType("STOP_LOSS_LIMIT")

	// BinanceRequestParamsOrderTakeProfit TAKE_PROFIT
	BinanceRequestParamsOrderTakeProfit = RequestParamsOrderType("TAKE_PROFIT")

	// BinanceRequestParamsOrderTakeProfitLimit TAKE_PROFIT_LIMIT
	BinanceRequestParamsOrderTakeProfitLimit = RequestParamsOrderType("TAKE_PROFIT_LIMIT")

	// BinanceRequestParamsOrderLimitMarker LIMIT_MAKER
	BinanceRequestParamsOrderLimitMarker = RequestParamsOrderType("LIMIT_MAKER")
)

// KlinesRequestParams represents Klines request data.
type KlinesRequestParams struct {
	Symbol    string // Required field; example LTCBTC, BTCUSDT
	Interval  string // Time interval period
	Limit     int    // Default 500; max 500.
	StartTime int64
	EndTime   int64
}

// WithdrawalFees the large list of predefined withdrawal fees
// Prone to change
var WithdrawalFees = map[currency.Code]float64{
	currency.BNB:     0.13,
	currency.BTC:     0.0005,
	currency.NEO:     0,
	currency.ETH:     0.01,
	currency.LTC:     0.001,
	currency.QTUM:    0.01,
	currency.EOS:     0.1,
	currency.SNT:     35,
	currency.BNT:     1,
	currency.GAS:     0,
	currency.BCC:     0.001,
	currency.BTM:     5,
	currency.USDT:    3.4,
	currency.HCC:     0.0005,
	currency.OAX:     6.5,
	currency.DNT:     54,
	currency.MCO:     0.31,
	currency.ICN:     3.5,
	currency.ZRX:     1.9,
	currency.OMG:     0.4,
	currency.WTC:     0.5,
	currency.LRC:     12.3,
	currency.LLT:     67.8,
	currency.YOYO:    1,
	currency.TRX:     1,
	currency.STRAT:   0.1,
	currency.SNGLS:   54,
	currency.BQX:     3.9,
	currency.KNC:     3.5,
	currency.SNM:     25,
	currency.FUN:     86,
	currency.LINK:    4,
	currency.XVG:     0.1,
	currency.CTR:     35,
	currency.SALT:    2.3,
	currency.MDA:     2.3,
	currency.IOTA:    0.5,
	currency.SUB:     11.4,
	currency.ETC:     0.01,
	currency.MTL:     2,
	currency.MTH:     45,
	currency.ENG:     2.2,
	currency.AST:     14.4,
	currency.DASH:    0.002,
	currency.BTG:     0.001,
	currency.EVX:     2.8,
	currency.REQ:     29.9,
	currency.VIB:     30,
	currency.POWR:    8.2,
	currency.ARK:     0.2,
	currency.XRP:     0.25,
	currency.MOD:     2,
	currency.ENJ:     26,
	currency.STORJ:   5.1,
	currency.KMD:     0.002,
	currency.RCN:     47,
	currency.NULS:    0.01,
	currency.RDN:     2.5,
	currency.XMR:     0.04,
	currency.DLT:     19.8,
	currency.AMB:     8.9,
	currency.BAT:     8,
	currency.ZEC:     0.005,
	currency.BCPT:    14.5,
	currency.ARN:     3,
	currency.GVT:     0.13,
	currency.CDT:     81,
	currency.GXS:     0.3,
	currency.POE:     134,
	currency.QSP:     36,
	currency.BTS:     1,
	currency.XZC:     0.02,
	currency.LSK:     0.1,
	currency.TNT:     47,
	currency.FUEL:    79,
	currency.MANA:    18,
	currency.BCD:     0.01,
	currency.DGD:     0.04,
	currency.ADX:     6.3,
	currency.ADA:     1,
	currency.PPT:     0.41,
	currency.CMT:     12,
	currency.XLM:     0.01,
	currency.CND:     58,
	currency.LEND:    84,
	currency.WABI:    6.6,
	currency.SBTC:    0.0005,
	currency.BCX:     0.5,
	currency.WAVES:   0.002,
	currency.TNB:     139,
	currency.GTO:     20,
	currency.ICX:     0.02,
	currency.OST:     32,
	currency.ELF:     3.9,
	currency.AION:    3.2,
	currency.CVC:     10.9,
	currency.REP:     0.2,
	currency.GNT:     8.9,
	currency.DATA:    37,
	currency.ETF:     1,
	currency.BRD:     3.8,
	currency.NEBL:    0.01,
	currency.VIBE:    17.3,
	currency.LUN:     0.36,
	currency.CHAT:    60.7,
	currency.RLC:     3.4,
	currency.INS:     3.5,
	currency.IOST:    105.6,
	currency.STEEM:   0.01,
	currency.NANO:    0.01,
	currency.AE:      1.3,
	currency.VIA:     0.01,
	currency.BLZ:     10.3,
	currency.SYS:     1,
	currency.NCASH:   247.6,
	currency.POA:     0.01,
	currency.ONT:     1,
	currency.ZIL:     37.2,
	currency.STORM:   152,
	currency.XEM:     4,
	currency.WAN:     0.1,
	currency.WPR:     43.4,
	currency.QLC:     1,
	currency.GRS:     0.2,
	currency.CLOAK:   0.02,
	currency.LOOM:    11.9,
	currency.BCN:     1,
	currency.TUSD:    1.35,
	currency.ZEN:     0.002,
	currency.SKY:     0.01,
	currency.THETA:   24,
	currency.IOTX:    90.5,
	currency.QKC:     24.6,
	currency.AGI:     29.81,
	currency.NXS:     0.02,
	currency.SC:      0.1,
	currency.EON:     10,
	currency.NPXS:    897,
	currency.KEY:     223,
	currency.NAS:     0.1,
	currency.ADD:     100,
	currency.MEETONE: 300,
	currency.ATD:     100,
	currency.MFT:     175,
	currency.EOP:     5,
	currency.DENT:    596,
	currency.IQ:      50,
	currency.ARDR:    2,
	currency.HOT:     1210,
	currency.VET:     100,
	currency.DOCK:    68,
	currency.POLY:    7,
	currency.VTHO:    21,
	currency.ONG:     0.1,
	currency.PHX:     1,
	currency.HC:      0.005,
	currency.GO:      0.01,
	currency.PAX:     1.4,
	currency.EDO:     1.3,
	currency.WINGS:   8.9,
	currency.NAV:     0.2,
	currency.TRIG:    49.1,
	currency.APPC:    12.4,
	currency.PIVX:    0.02,
}

// WithdrawResponse contains status of withdrawal request
type WithdrawResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	ID      string `json:"id"`
}

// UserAccountStream contains a key to maintain an authorised
// websocket connection
type UserAccountStream struct {
	ListenKey string `json:"listenKey"`
}

type wsAccountInfo struct {
	Stream string `json:"stream"`
	Data   struct {
		CanDeposit       bool    `json:"D"`
		CanTrade         bool    `json:"T"`
		CanWithdraw      bool    `json:"W"`
		EventTime        int64   `json:"E"`
		LastUpdated      int64   `json:"u"`
		BuyerCommission  float64 `json:"b"`
		MakerCommission  float64 `json:"m"`
		SellerCommission float64 `json:"s"`
		TakerCommission  float64 `json:"t"`
		EventType        string  `json:"e"`
		Currencies       []struct {
			Asset     string  `json:"a"`
			Available float64 `json:"f,string"`
			Locked    float64 `json:"l,string"`
		} `json:"B"`
	} `json:"data"`
}

type wsAccountPosition struct {
	Stream string `json:"stream"`
	Data   struct {
		Currencies []struct {
			Asset     string  `json:"a"`
			Available float64 `json:"f,string"`
			Locked    float64 `json:"l,string"`
		} `json:"B"`
		EventTime   int64  `json:"E"`
		LastUpdated int64  `json:"u"`
		EventType   string `json:"e"`
	} `json:"data"`
}

type wsBalanceUpdate struct {
	Stream string `json:"stream"`
	Data   struct {
		EventTime    int64   `json:"E"`
		ClearTime    int64   `json:"T"`
		BalanceDelta float64 `json:"d,string"`
		Asset        string  `json:"a"`
		EventType    string  `json:"e"`
	} `json:"data"`
}

type wsOrderUpdate struct {
	Stream string `json:"stream"`
	Data   struct {
		ClientOrderID                     string  `json:"C"`
		EventTime                         int64   `json:"E"`
		IcebergQuantity                   float64 `json:"F,string"`
		LastExecutedPrice                 float64 `json:"L,string"`
		CommissionAsset                   float64 `json:"N"`
		OrderCreationTime                 int64   `json:"O"`
		StopPrice                         float64 `json:"P,string"`
		QuoteOrderQuantity                float64 `json:"Q,string"`
		Side                              string  `json:"S"`
		TransactionTime                   int64   `json:"T"`
		OrderStatus                       string  `json:"X"`
		LastQuoteAssetTransactedQuantity  float64 `json:"Y,string"`
		CumulativeQuoteTransactedQuantity float64 `json:"Z,string"`
		CancelledClientOrderID            string  `json:"c"`
		EventType                         string  `json:"e"`
		TimeInForce                       string  `json:"f"`
		OrderListID                       int64   `json:"g"`
		OrderID                           int64   `json:"i"`
		LastExecutedQuantity              float64 `json:"l,string"`
		IsMaker                           bool    `json:"m"`
		Commission                        float64 `json:"n,string"`
		OrderType                         string  `json:"o"`
		Price                             float64 `json:"p,string"`
		Quantity                          float64 `json:"q,string"`
		RejectionReason                   string  `json:"r"`
		Symbol                            string  `json:"s"`
		TradeID                           int64   `json:"t"`
		IsOnOrderBook                     bool    `json:"w"`
		CurrentExecutionType              string  `json:"x"`
		CumulativeFilledQuantity          float64 `json:"z,string"`
	} `json:"data"`
}

type wsListStatus struct {
	Stream string `json:"stream"`
	Data   struct {
		ListClientOrderID string `json:"C"`
		EventTime         int64  `json:"E"`
		ListOrderStatus   string `json:"L"`
		Orders            []struct {
			ClientOrderID string `json:"c"`
			OrderID       int64  `json:"i"`
			Symbol        string `json:"s"`
		} `json:"O"`
		TransactionTime int64  `json:"T"`
		ContingencyType string `json:"c"`
		EventType       string `json:"e"`
		OrderListID     int64  `json:"g"`
		ListStatusType  string `json:"l"`
		RejectionReason string `json:"r"`
		Symbol          string `json:"s"`
	} `json:"data"`
}

// WsPayload defines the payload through the websocket connection
type WsPayload struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int64    `json:"id"`
}
