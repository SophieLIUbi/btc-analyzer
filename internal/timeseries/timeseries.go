package timeseries

import (
	"btc-analyzer/internal/types"
	"sort"
	"time"
)

// New creates a new Bitcoin time series
func New(symbol string) *types.BTCTimeSeries {
	return &types.BTCTimeSeries{
		Symbol: symbol,
		Data:   make([]types.BTCPrice, 0),
	}
}

// AddPrice adds a price point to the series
func AddPrice(bts *types.BTCTimeSeries, price types.BTCPrice) {
	bts.Data = append(bts.Data, price)
}

// Sort sorts the data by timestamp
func Sort(bts *types.BTCTimeSeries) {
	sort.Slice(bts.Data, func(i, j int) bool {
		return bts.Data[i].Timestamp.Before(bts.Data[j].Timestamp)
	})
}

// GetClosePrices extracts closing prices for analysis
func GetClosePrices(bts *types.BTCTimeSeries) []float64 {
	prices := make([]float64, len(bts.Data))
	for i, data := range bts.Data {
		prices[i] = data.Close
	}
	return prices
}

// GetVolumeData extracts volume data
func GetVolumeData(bts *types.BTCTimeSeries) []float64 {
	volumes := make([]float64, len(bts.Data))
	for i, data := range bts.Data {
		volumes[i] = data.Volume
	}
	return volumes
}

// GetTimeRange returns the time range of the data
func GetTimeRange(bts *types.BTCTimeSeries) (time.Time, time.Time) {
	if len(bts.Data) == 0 {
		return time.Time{}, time.Time{}
	}
	Sort(bts)
	return bts.Data[0].Timestamp, bts.Data[len(bts.Data)-1].Timestamp
}

// GetLatestPrice returns the most recent price data
func GetLatestPrice(bts *types.BTCTimeSeries) types.BTCPrice {
	if len(bts.Data) == 0 {
		return types.BTCPrice{}
	}
	Sort(bts)
	return bts.Data[len(bts.Data)-1]
}

// FilterByDateRange filters data within a specific date range
func FilterByDateRange(bts *types.BTCTimeSeries, start, end time.Time) *types.BTCTimeSeries {
	filtered := New(bts.Symbol + "_filtered")
	
	for _, price := range bts.Data {
		if (price.Timestamp.Equal(start) || price.Timestamp.After(start)) &&
		   (price.Timestamp.Equal(end) || price.Timestamp.Before(end)) {
			AddPrice(filtered, price)
		}
	}
	
	return filtered
}

// ResampleToDaily resamples data to daily intervals
func ResampleToDaily(bts *types.BTCTimeSeries) *types.BTCTimeSeries {
	if len(bts.Data) == 0 {
		return New(bts.Symbol + "_daily")
	}

	Sort(bts)
	resampled := New(bts.Symbol + "_daily")
	
	currentDay := bts.Data[0].Timestamp.Truncate(24 * time.Hour)
	var dayData []types.BTCPrice
	
	for _, price := range bts.Data {
		priceDay := price.Timestamp.Truncate(24 * time.Hour)
		
		if priceDay.Equal(currentDay) {
			dayData = append(dayData, price)
		} else {
			// Process accumulated day data
			if len(dayData) > 0 {
				dailyPrice := aggregateDayData(dayData, currentDay)
				AddPrice(resampled, dailyPrice)
			}
			
			// Start new day
			currentDay = priceDay
			dayData = []types.BTCPrice{price}
		}
	}
	
	// Process last day
	if len(dayData) > 0 {
		dailyPrice := aggregateDayData(dayData, currentDay)
		AddPrice(resampled, dailyPrice)
	}
	
	return resampled
}

// aggregateDayData aggregates multiple price points into a single daily OHLCV
func aggregateDayData(dayData []types.BTCPrice, day time.Time) types.BTCPrice {
	if len(dayData) == 0 {
		return types.BTCPrice{}
	}
	
	open := dayData[0].Open
	close := dayData[len(dayData)-1].Close
	high := dayData[0].High
	low := dayData[0].Low
	volume := 0.0
	
	for _, price := range dayData {
		if price.High > high {
			high = price.High
		}
		if price.Low < low {
			low = price.Low
		}
		volume += price.Volume
	}
	
	return types.BTCPrice{
		Timestamp: day,
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close,
		Volume:    volume,
	}
}