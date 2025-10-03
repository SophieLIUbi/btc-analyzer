package statistics

import (
	"btc-analyzer/internal/timeseries"
	"btc-analyzer/internal/types"
	"math"
	"sort"
)

// Calculate calculates comprehensive statistics
func Calculate(values []float64) types.Statistics {
	if len(values) == 0 {
		return types.Statistics{}
	}

	// Create a copy for sorting
	sortedValues := make([]float64, len(values))
	copy(sortedValues, values)
	sort.Float64s(sortedValues)
	
	n := len(values)
	
	// Basic stats
	sum := 0.0
	min := sortedValues[0]
	max := sortedValues[n-1]
	
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(n)

	// Median
	var median float64
	if n%2 == 0 {
		median = (sortedValues[n/2-1] + sortedValues[n/2]) / 2
	} else {
		median = sortedValues[n/2]
	}

	// Variance and standard deviation
	sumSquaredDiff := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}
	variance := sumSquaredDiff / float64(n)
	stdDev := math.Sqrt(variance)

	// Skewness and kurtosis
	sumCubedDiff := 0.0
	sumQuartedDiff := 0.0
	for _, v := range values {
		diff := v - mean
		cubedDiff := diff * diff * diff
		quartedDiff := cubedDiff * diff
		sumCubedDiff += cubedDiff
		sumQuartedDiff += quartedDiff
	}
	
	skewness := 0.0
	kurtosis := 0.0
	if stdDev > 0 {
		skewness = (sumCubedDiff / float64(n)) / math.Pow(stdDev, 3)
		kurtosis = (sumQuartedDiff / float64(n)) / math.Pow(stdDev, 4) - 3
	}

	return types.Statistics{
		Count:    n,
		Mean:     mean,
		Median:   median,
		StdDev:   stdDev,
		Min:      min,
		Max:      max,
		Variance: variance,
		Skewness: skewness,
		Kurtosis: kurtosis,
	}
}

// CalculateReturns calculates simple and log returns
func CalculateReturns(bts *types.BTCTimeSeries) ([]float64, []float64) {
	if len(bts.Data) < 2 {
		return nil, nil
	}

	returns := make([]float64, len(bts.Data)-1)
	logReturns := make([]float64, len(bts.Data)-1)

	for i := 1; i < len(bts.Data); i++ {
		prevPrice := bts.Data[i-1].Close
		currPrice := bts.Data[i].Close
		
		if prevPrice > 0 {
			returns[i-1] = (currPrice - prevPrice) / prevPrice
			logReturns[i-1] = math.Log(currPrice / prevPrice)
		}
	}

	return returns, logReturns
}

// CalculateVolatility calculates annualized volatility
func CalculateVolatility(returns []float64, periodsPerYear int) float64 {
	if len(returns) == 0 {
		return 0
	}

	stats := Calculate(returns)
	volatility := stats.StdDev * math.Sqrt(float64(periodsPerYear))
	
	return volatility
}

// CalculateMaxDrawdown calculates maximum drawdown
func CalculateMaxDrawdown(bts *types.BTCTimeSeries) float64 {
	prices := timeseries.GetClosePrices(bts)
	if len(prices) == 0 {
		return 0
	}

	maxDrawdown := 0.0
	peak := prices[0]

	for _, price := range prices {
		if price > peak {
			peak = price
		}
		drawdown := (peak - price) / peak
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	return maxDrawdown
}

// CalculateSharpeRatio calculates Sharpe ratio
func CalculateSharpeRatio(returns []float64, riskFreeRate float64, periodsPerYear int) float64 {
	if len(returns) == 0 {
		return 0
	}

	stats := Calculate(returns)
	if stats.StdDev == 0 {
		return 0
	}

	annualizedReturn := stats.Mean * float64(periodsPerYear)
	annualizedVolatility := stats.StdDev * math.Sqrt(float64(periodsPerYear))
	
	return (annualizedReturn - riskFreeRate) / annualizedVolatility
}

// CalculateCorrelation calculates correlation between two series
func CalculateCorrelation(x, y []float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0
	}

	n := len(x)
	sumX, sumY, sumXY, sumX2, sumY2 := 0.0, 0.0, 0.0, 0.0, 0.0

	for i := 0; i < n; i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
		sumY2 += y[i] * y[i]
	}

	numerator := float64(n)*sumXY - sumX*sumY
	denominator := math.Sqrt((float64(n)*sumX2 - sumX*sumX) * (float64(n)*sumY2 - sumY*sumY))

	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}

// GetRiskMetrics calculates comprehensive risk metrics
func GetRiskMetrics(bts *types.BTCTimeSeries) map[string]float64 {
	metrics := make(map[string]float64)
	
	if len(bts.Data) < 30 {
		return metrics
	}

	returns, _ := CalculateReturns(bts)
	if len(returns) == 0 {
		return metrics
	}

	volatility := CalculateVolatility(returns, 365)
	maxDrawdown := CalculateMaxDrawdown(bts)
	sharpeRatio := CalculateSharpeRatio(returns, 0.0, 365)
	
	// Basic risk metrics
	metrics["volatility_annual"] = volatility
	metrics["max_drawdown"] = maxDrawdown
	metrics["sharpe_ratio"] = sharpeRatio
	
	// Value at Risk (VaR) - 95% confidence level
	returnStats := Calculate(returns)
	metrics["var_95"] = returnStats.Mean - 1.645*returnStats.StdDev // Daily VaR
	metrics["var_95_annual"] = metrics["var_95"] * math.Sqrt(365)
	
	// Conditional Value at Risk (CVaR)
	sortedReturns := make([]float64, len(returns))
	copy(sortedReturns, returns)
	sort.Float64s(sortedReturns)
	
	var5Index := int(0.05 * float64(len(sortedReturns)))
	if var5Index < len(sortedReturns) {
		cvarSum := 0.0
		for i := 0; i <= var5Index; i++ {
			cvarSum += sortedReturns[i]
		}
		metrics["cvar_95"] = cvarSum / float64(var5Index+1)
	}
	
	// Sortino ratio (downside deviation)
	downsideReturns := make([]float64, 0)
	for _, ret := range returns {
		if ret < 0 {
			downsideReturns = append(downsideReturns, ret)
		}
	}
	
	if len(downsideReturns) > 0 {
		downsideStats := Calculate(downsideReturns)
		downsideDeviation := downsideStats.StdDev * math.Sqrt(365)
		if downsideDeviation > 0 {
			metrics["sortino_ratio"] = (returnStats.Mean * 365) / downsideDeviation
		}
	}
	
	// Beta (if we had market data, for now use volatility ratio)
	marketVolatility := 0.16 // Assume 16% market volatility
	metrics["beta_estimate"] = volatility / marketVolatility
	
	return metrics
}

// PerformBacktest performs simple buy-and-hold backtest
func PerformBacktest(bts *types.BTCTimeSeries, startAmount float64) map[string]float64 {
	results := make(map[string]float64)
	
	if len(bts.Data) < 2 {
		return results
	}

	timeseries.Sort(bts)
	startPrice := bts.Data[0].Close
	endPrice := bts.Data[len(bts.Data)-1].Close
	
	btcAmount := startAmount / startPrice
	endValue := btcAmount * endPrice
	
	totalReturn := (endValue - startAmount) / startAmount
	
	days := float64(bts.Data[len(bts.Data)-1].Timestamp.Sub(bts.Data[0].Timestamp).Hours() / 24)
	annualizedReturn := math.Pow(1+totalReturn, 365/days) - 1
	
	results["start_amount"] = startAmount
	results["end_value"] = endValue
	results["total_return"] = totalReturn
	results["annualized_return"] = annualizedReturn
	results["btc_purchased"] = btcAmount
	results["days_held"] = days
	results["start_price"] = startPrice
	results["end_price"] = endPrice
	
	return results
}