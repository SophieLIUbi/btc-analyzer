package patterns

import (
	"btc-analyzer/internal/timeseries"
	"btc-analyzer/internal/types"
	"math"
	"sort"
)

// FindSupportResistanceLevels identifies key support and resistance levels
func FindSupportResistanceLevels(bts *types.BTCTimeSeries, lookbackPeriod int, tolerance float64) types.SupportResistanceData {
	if len(bts.Data) < lookbackPeriod*2 {
		return types.SupportResistanceData{}
	}

	timeseries.Sort(bts)
	
	var supportLevels []float64
	var resistanceLevels []float64
	
	// Find potential support and resistance points
	for i := lookbackPeriod; i < len(bts.Data)-lookbackPeriod; i++ {
		currentPrice := bts.Data[i]
		
		// Check if current point is a local minimum (support)
		isSupport := true
		isResistance := true
		
		for j := i - lookbackPeriod; j <= i+lookbackPeriod; j++ {
			if j != i {
				if bts.Data[j].Low < currentPrice.Low {
					isSupport = false
				}
				if bts.Data[j].High > currentPrice.High {
					isResistance = false
				}
			}
		}
		
		if isSupport {
			supportLevels = append(supportLevels, currentPrice.Low)
		}
		if isResistance {
			resistanceLevels = append(resistanceLevels, currentPrice.High)
		}
	}
	
	// Cluster nearby levels
	supportLevels = clusterLevels(supportLevels, tolerance)
	resistanceLevels = clusterLevels(resistanceLevels, tolerance)
	
	return types.SupportResistanceData{
		SupportLevels:    supportLevels,
		ResistanceLevels: resistanceLevels,
	}
}

// clusterLevels groups nearby price levels together
func clusterLevels(levels []float64, tolerance float64) []float64 {
	if len(levels) == 0 {
		return levels
	}
	
	sort.Float64s(levels)
	clustered := make([]float64, 0)
	
	currentCluster := []float64{levels[0]}
	
	for i := 1; i < len(levels); i++ {
		if math.Abs(levels[i]-levels[i-1])/levels[i-1] <= tolerance {
			currentCluster = append(currentCluster, levels[i])
		} else {
			// Calculate average of current cluster
			sum := 0.0
			for _, level := range currentCluster {
				sum += level
			}
			clustered = append(clustered, sum/float64(len(currentCluster)))
			
			// Start new cluster
			currentCluster = []float64{levels[i]}
		}
	}
	
	// Add last cluster
	if len(currentCluster) > 0 {
		sum := 0.0
		for _, level := range currentCluster {
			sum += level
		}
		clustered = append(clustered, sum/float64(len(currentCluster)))
	}
	
	return clustered
}

// DetectTrend analyzes overall trend direction
func DetectTrend(bts *types.BTCTimeSeries, period int) string {
	if len(bts.Data) < period {
		return "insufficient_data"
	}
	
	prices := timeseries.GetClosePrices(bts)
	startPrice := prices[len(prices)-period]
	endPrice := prices[len(prices)-1]
	
	change := (endPrice - startPrice) / startPrice
	
	if change > 0.05 {
		return "uptrend"
	} else if change < -0.05 {
		return "downtrend"
	}
	return "sideways"
}

// DetectCandlestickPatterns identifies common candlestick patterns
func DetectCandlestickPatterns(bts *types.BTCTimeSeries) map[string][]int {
	patterns := make(map[string][]int)
	
	if len(bts.Data) < 3 {
		return patterns
	}
	
	timeseries.Sort(bts)
	
	for i := 1; i < len(bts.Data)-1; i++ {
		prev := bts.Data[i-1]
		curr := bts.Data[i]
		
		
		// Doji pattern
		if isDoji(curr) {
			patterns["doji"] = append(patterns["doji"], i)
		}
		
		// Hammer pattern
		if isHammer(curr) {
			patterns["hammer"] = append(patterns["hammer"], i)
		}
		
		// Shooting star pattern
		if isShootingStar(curr) {
			patterns["shooting_star"] = append(patterns["shooting_star"], i)
		}
		
		// Engulfing patterns
		if isBullishEngulfing(prev, curr) {
			patterns["bullish_engulfing"] = append(patterns["bullish_engulfing"], i)
		}
		
		if isBearishEngulfing(prev, curr) {
			patterns["bearish_engulfing"] = append(patterns["bearish_engulfing"], i)
		}
		
		// Three-candle patterns
		if i > 1 {
			prevPrev := bts.Data[i-2]
			
			// Morning star
			if isMorningStar(prevPrev, prev, curr) {
				patterns["morning_star"] = append(patterns["morning_star"], i)
			}
			
			// Evening star
			if isEveningStar(prevPrev, prev, curr) {
				patterns["evening_star"] = append(patterns["evening_star"], i)
			}
		}
	}
	
	return patterns
}

// Candlestick pattern helper functions
func isDoji(candle types.BTCPrice) bool {
	body := math.Abs(candle.Close - candle.Open)
	range_ := candle.High - candle.Low
	return range_ > 0 && body/range_ < 0.1
}

func isHammer(candle types.BTCPrice) bool {
	body := math.Abs(candle.Close - candle.Open)
	lowerShadow := math.Min(candle.Open, candle.Close) - candle.Low
	upperShadow := candle.High - math.Max(candle.Open, candle.Close)
	range_ := candle.High - candle.Low
	
	return range_ > 0 && lowerShadow > 2*body && upperShadow < body*0.5
}

func isShootingStar(candle types.BTCPrice) bool {
	body := math.Abs(candle.Close - candle.Open)
	lowerShadow := math.Min(candle.Open, candle.Close) - candle.Low
	upperShadow := candle.High - math.Max(candle.Open, candle.Close)
	range_ := candle.High - candle.Low
	
	return range_ > 0 && upperShadow > 2*body && lowerShadow < body*0.5
}

func isBullishEngulfing(prev, curr types.BTCPrice) bool {
	prevBearish := prev.Close < prev.Open
	currBullish := curr.Close > curr.Open
	
	return prevBearish && currBullish && 
		   curr.Open < prev.Close && 
		   curr.Close > prev.Open
}

func isBearishEngulfing(prev, curr types.BTCPrice) bool {
	prevBullish := prev.Close > prev.Open
	currBearish := curr.Close < curr.Open
	
	return prevBullish && currBearish && 
		   curr.Open > prev.Close && 
		   curr.Close < prev.Open
}

func isMorningStar(first, second, third types.BTCPrice) bool {
	firstBearish := first.Close < first.Open
	secondSmall := math.Abs(second.Close-second.Open) < math.Abs(first.Close-first.Open)*0.3
	thirdBullish := third.Close > third.Open
	
	return firstBearish && secondSmall && thirdBullish &&
		   second.High < first.Low &&
		   third.Close > (first.Open+first.Close)/2
}

func isEveningStar(first, second, third types.BTCPrice) bool {
	firstBullish := first.Close > first.Open
	secondSmall := math.Abs(second.Close-second.Open) < math.Abs(first.Close-first.Open)*0.3
	thirdBearish := third.Close < third.Open
	
	return firstBullish && secondSmall && thirdBearish &&
		   second.Low > first.High &&
		   third.Close < (first.Open+first.Close)/2
}

// DetectVolumePatterns analyzes volume patterns
func DetectVolumePatterns(bts *types.BTCTimeSeries) map[string][]int {
	patterns := make(map[string][]int)
	
	if len(bts.Data) < 20 {
		return patterns
	}
	
	volumes := timeseries.GetVolumeData(bts)
	
	// Calculate average volume for comparison
	sum := 0.0
	for _, vol := range volumes {
		sum += vol
	}
	avgVolume := sum / float64(len(volumes))
	
	for i := 1; i < len(bts.Data); i++ {
		curr := bts.Data[i]
		prev := bts.Data[i-1]
		
		// Volume spike with price increase
		if curr.Volume > avgVolume*2 && curr.Close > prev.Close*1.02 {
			patterns["volume_breakout"] = append(patterns["volume_breakout"], i)
		}
		
		// Volume spike with price decrease
		if curr.Volume > avgVolume*2 && curr.Close < prev.Close*0.98 {
			patterns["volume_selloff"] = append(patterns["volume_selloff"], i)
		}
		
		// Low volume drift
		if curr.Volume < avgVolume*0.5 {
			patterns["low_volume"] = append(patterns["low_volume"], i)
		}
	}
	
	return patterns
}

// FindPivotPoints calculates pivot points for the day
func FindPivotPoints(bts *types.BTCTimeSeries) map[string]float64 {
	pivots := make(map[string]float64)
	
	if len(bts.Data) == 0 {
		return pivots
	}
	
	// Use the latest complete day's data
	latest := bts.Data[len(bts.Data)-1]
	high := latest.High
	low := latest.Low
	close := latest.Close
	
	// Standard pivot point calculation
	pivot := (high + low + close) / 3
	
	pivots["pivot"] = pivot
	pivots["r1"] = 2*pivot - low      // Resistance 1
	pivots["s1"] = 2*pivot - high     // Support 1
	pivots["r2"] = pivot + (high - low) // Resistance 2
	pivots["s2"] = pivot - (high - low) // Support 2
	pivots["r3"] = high + 2*(pivot - low) // Resistance 3
	pivots["s3"] = low - 2*(high - pivot) // Support 3
	
	return pivots
}

// CalculateFibonacciRetracements calculates Fibonacci retracement levels
func CalculateFibonacciRetracements(bts *types.BTCTimeSeries, period int) map[string]float64 {
	fibs := make(map[string]float64)
	
	if len(bts.Data) < period {
		return fibs
	}
	
	// Find high and low in the period
	recentData := bts.Data[len(bts.Data)-period:]
	high := recentData[0].High
	low := recentData[0].Low
	
	for _, data := range recentData {
		if data.High > high {
			high = data.High
		}
		if data.Low < low {
			low = data.Low
		}
	}
	
	range_ := high - low
	
	// Standard Fibonacci retracement levels
	fibLevels := []float64{0.0, 0.236, 0.382, 0.5, 0.618, 0.786, 1.0}
	fibNames := []string{"high", "fib_76_4", "fib_61_8", "fib_50", "fib_38_2", "fib_23_6", "low"}
	
	for i, level := range fibLevels {
		fibs[fibNames[i]] = high - (range_ * level)
	}
	
	return fibs
}