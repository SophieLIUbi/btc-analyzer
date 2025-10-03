package indicators

import (
	"btc-analyzer/internal/timeseries"
	"btc-analyzer/internal/types"
	"math"
)

// CalculateRSI calculates Relative Strength Index
func CalculateRSI(bts *types.BTCTimeSeries, period int) []float64 {
	if len(bts.Data) < period+1 {
		return nil
	}

	prices := timeseries.GetClosePrices(bts)
	rsi := make([]float64, len(prices)-period)

	// Calculate price changes
	changes := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		changes[i-1] = prices[i] - prices[i-1]
	}

	// Initial RS calculation
	avgGain := 0.0
	avgLoss := 0.0
	
	for i := 0; i < period; i++ {
		if changes[i] > 0 {
			avgGain += changes[i]
		} else {
			avgLoss += math.Abs(changes[i])
		}
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// Calculate RSI
	for i := period; i < len(changes); i++ {
		change := changes[i]
		
		gain := 0.0
		loss := 0.0
		if change > 0 {
			gain = change
		} else {
			loss = math.Abs(change)
		}

		avgGain = (avgGain*float64(period-1) + gain) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + loss) / float64(period)

		var rsiValue float64
		if avgLoss == 0 {
			rsiValue = 100
		} else {
			rs := avgGain / avgLoss
			rsiValue = 100 - (100 / (1 + rs))
		}
		rsi[i-period] = rsiValue
	}

	return rsi
}

// CalculateMACD calculates MACD indicator
func CalculateMACD(bts *types.BTCTimeSeries, fastPeriod, slowPeriod, signalPeriod int) types.MACDData {
	prices := timeseries.GetClosePrices(bts)
	if len(prices) < slowPeriod {
		return types.MACDData{}
	}

	// Calculate EMAs
	fastEMA := calculateEMA(prices, fastPeriod)
	slowEMA := calculateEMA(prices, slowPeriod)

	// Align arrays (slow EMA starts later)
	startIdx := slowPeriod - fastPeriod
	alignedFastEMA := fastEMA[startIdx:]

	// Calculate MACD line
	macdLine := make([]float64, len(slowEMA))
	for i := range slowEMA {
		if i < len(alignedFastEMA) {
			macdLine[i] = alignedFastEMA[i] - slowEMA[i]
		}
	}

	// Calculate signal line (EMA of MACD)
	signalLine := calculateEMA(macdLine, signalPeriod)

	// Calculate histogram
	histogram := make([]float64, len(signalLine))
	startIdx2 := len(macdLine) - len(signalLine)
	for i := range signalLine {
		if startIdx2+i < len(macdLine) {
			histogram[i] = macdLine[startIdx2+i] - signalLine[i]
		}
	}

	return types.MACDData{
		MACD:      macdLine,
		Signal:    signalLine,
		Histogram: histogram,
	}
}

// calculateEMA calculates Exponential Moving Average
func calculateEMA(prices []float64, period int) []float64 {
	if len(prices) < period {
		return nil
	}

	ema := make([]float64, len(prices)-period+1)
	multiplier := 2.0 / (float64(period) + 1.0)

	// Start with SMA for first value
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema[0] = sum / float64(period)

	// Calculate EMA for remaining values
	for i := 1; i < len(ema); i++ {
		ema[i] = (prices[period-1+i] * multiplier) + (ema[i-1] * (1 - multiplier))
	}

	return ema
}

// CalculateBollingerBands calculates Bollinger Bands
func CalculateBollingerBands(bts *types.BTCTimeSeries, period int, stdDevFactor float64) types.BollingerBandsData {
	prices := timeseries.GetClosePrices(bts)
	if len(prices) < period {
		return types.BollingerBandsData{}
	}

	middle := make([]float64, len(prices)-period+1)
	upper := make([]float64, len(prices)-period+1)
	lower := make([]float64, len(prices)-period+1)

	for i := period - 1; i < len(prices); i++ {
		// Calculate SMA
		sum := 0.0
		for j := i - period + 1; j <= i; j++ {
			sum += prices[j]
		}
		sma := sum / float64(period)
		middle[i-period+1] = sma

		// Calculate standard deviation
		sumSquaredDiff := 0.0
		for j := i - period + 1; j <= i; j++ {
			diff := prices[j] - sma
			sumSquaredDiff += diff * diff
		}
		stdDev := math.Sqrt(sumSquaredDiff / float64(period))

		upper[i-period+1] = sma + (stdDevFactor * stdDev)
		lower[i-period+1] = sma - (stdDevFactor * stdDev)
	}

	return types.BollingerBandsData{
		Upper:  upper,
		Middle: middle,
		Lower:  lower,
	}
}

// CalculateMovingAverage calculates simple moving average
func CalculateMovingAverage(bts *types.BTCTimeSeries, period int) []float64 {
	if len(bts.Data) < period {
		return nil
	}

	prices := timeseries.GetClosePrices(bts)
	ma := make([]float64, len(prices)-period+1)
	
	for i := period - 1; i < len(prices); i++ {
		sum := 0.0
		for j := i - period + 1; j <= i; j++ {
			sum += prices[j]
		}
		ma[i-period+1] = sum / float64(period)
	}
	
	return ma
}

// CalculateStochasticOscillator calculates Stochastic Oscillator
func CalculateStochasticOscillator(bts *types.BTCTimeSeries, kPeriod int) []float64 {
	if len(bts.Data) < kPeriod {
		return nil
	}

	stochastic := make([]float64, len(bts.Data)-kPeriod+1)

	for i := kPeriod - 1; i < len(bts.Data); i++ {
		// Find highest high and lowest low in the period
		highestHigh := bts.Data[i-kPeriod+1].High
		lowestLow := bts.Data[i-kPeriod+1].Low

		for j := i - kPeriod + 1; j <= i; j++ {
			if bts.Data[j].High > highestHigh {
				highestHigh = bts.Data[j].High
			}
			if bts.Data[j].Low < lowestLow {
				lowestLow = bts.Data[j].Low
			}
		}

		// Calculate %K
		currentClose := bts.Data[i].Close
		if highestHigh-lowestLow != 0 {
			stochastic[i-kPeriod+1] = ((currentClose - lowestLow) / (highestHigh - lowestLow)) * 100
		} else {
			stochastic[i-kPeriod+1] = 50 // Default to midpoint if no range
		}
	}

	return stochastic
}