package types

import "time"

// BTCPrice represents Bitcoin price data with OHLCV format
type BTCPrice struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

// BTCTimeSeries represents Bitcoin time series data
type BTCTimeSeries struct {
	Symbol string
	Data   []BTCPrice
}

// Statistics represents basic statistical measures
type Statistics struct {
	Count    int
	Mean     float64
	Median   float64
	StdDev   float64
	Min      float64
	Max      float64
	Variance float64
	Skewness float64
	Kurtosis float64
}

// MACDData holds MACD indicator values
type MACDData struct {
	MACD      []float64
	Signal    []float64
	Histogram []float64
}

// BollingerBandsData holds Bollinger Bands values
type BollingerBandsData struct {
	Upper  []float64
	Middle []float64
	Lower  []float64
}

// SupportResistanceData holds support and resistance levels
type SupportResistanceData struct {
	SupportLevels    []float64
	ResistanceLevels []float64
}

// BTCAnalytics holds comprehensive Bitcoin market analytics
type BTCAnalytics struct {
	PriceStats        Statistics
	VolumeStats       Statistics
	Volatility        float64
	SharpeRatio       float64
	MaxDrawdown       float64
	Returns           []float64
	LogReturns        []float64
	RSI               []float64
	MACD              MACDData
	BollingerBands    BollingerBandsData
	SupportResistance SupportResistanceData
}

// PriceAlert represents a price alert condition
type PriceAlert struct {
	Type      string // "above", "below", "change"
	Threshold float64
	Triggered bool
	Timestamp time.Time
}

// CoinGeckoResponse represents API response from CoinGecko
type CoinGeckoResponse struct {
	Prices       [][]float64 `json:"prices"`
	MarketCaps   [][]float64 `json:"market_caps"`
	TotalVolumes [][]float64 `json:"total_volumes"`
}