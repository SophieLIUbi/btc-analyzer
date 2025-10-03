package analyzer

import (
	"btc-analyzer/internal/indicators"
	"btc-analyzer/internal/patterns"
	"btc-analyzer/internal/statistics"
	"btc-analyzer/internal/timeseries"
	"btc-analyzer/internal/types"
	"fmt"
	"time"
	"math"
)

// PerformComprehensiveAnalysis runs a full analysis on Bitcoin data
func PerformComprehensiveAnalysis(bts *types.BTCTimeSeries) types.BTCAnalytics {
	analytics := types.BTCAnalytics{}
	
	if len(bts.Data) < 2 {
		return analytics
	}
	
	// Basic price and volume statistics
	prices := timeseries.GetClosePrices(bts)
	volumes := timeseries.GetVolumeData(bts)
	
	analytics.PriceStats = statistics.Calculate(prices)
	analytics.VolumeStats = statistics.Calculate(volumes)
	
	// Calculate returns
	returns, logReturns := statistics.CalculateReturns(bts)
	analytics.Returns = returns
	analytics.LogReturns = logReturns
	
	// Risk metrics
	if len(returns) > 0 {
		analytics.Volatility = statistics.CalculateVolatility(returns, 365)
		analytics.SharpeRatio = statistics.CalculateSharpeRatio(returns, 0.0, 365)
		analytics.MaxDrawdown = statistics.CalculateMaxDrawdown(bts)
	}
	
	// Technical indicators
	if len(bts.Data) >= 14 {
		analytics.RSI = indicators.CalculateRSI(bts, 14)
	}
	
	if len(bts.Data) >= 26 {
		analytics.MACD = indicators.CalculateMACD(bts, 12, 26, 9)
	}
	
	if len(bts.Data) >= 20 {
		analytics.BollingerBands = indicators.CalculateBollingerBands(bts, 20, 2.0)
	}
	
	// Pattern analysis
	if len(bts.Data) >= 10 {
		analytics.SupportResistance = patterns.FindSupportResistanceLevels(bts, 5, 0.02)
	}
	
	return analytics
}

// GenerateReport creates a comprehensive text report
func GenerateReport(bts *types.BTCTimeSeries, analytics types.BTCAnalytics) string {
	var report string
	
	report += "=== BITCOIN MARKET ANALYSIS REPORT ===\n\n"
	
	// Basic information
	report += fmt.Sprintf("Symbol: %s\n", bts.Symbol)
	report += fmt.Sprintf("Data Points: %d\n", len(bts.Data))
	
	if len(bts.Data) > 0 {
		start, end := timeseries.GetTimeRange(bts)
		report += fmt.Sprintf("Time Range: %s to %s\n", 
			start.Format("2006-01-02"), 
			end.Format("2006-01-02"))
		
		latest := timeseries.GetLatestPrice(bts)
		report += fmt.Sprintf("Latest Price: $%.2f\n", latest.Close)
		report += fmt.Sprintf("Latest Volume: %.0f\n\n", latest.Volume)
	}
	
	// Price statistics
	report += "=== PRICE STATISTICS ===\n"
	report += fmt.Sprintf("Mean Price: $%.2f\n", analytics.PriceStats.Mean)
	report += fmt.Sprintf("Median Price: $%.2f\n", analytics.PriceStats.Median)
	report += fmt.Sprintf("Price Range: $%.2f - $%.2f\n", analytics.PriceStats.Min, analytics.PriceStats.Max)
	report += fmt.Sprintf("Standard Deviation: $%.2f\n", analytics.PriceStats.StdDev)
	report += fmt.Sprintf("Price Variance: %.2f\n", analytics.PriceStats.Variance)
	
	if analytics.PriceStats.Skewness != 0 {
		report += fmt.Sprintf("Skewness: %.3f\n", analytics.PriceStats.Skewness)
		report += fmt.Sprintf("Kurtosis: %.3f\n", analytics.PriceStats.Kurtosis)
	}
	report += "\n"
	
	// Risk metrics
	if analytics.Volatility > 0 {
		report += "=== RISK METRICS ===\n"
		report += fmt.Sprintf("Annualized Volatility: %.2f%%\n", analytics.Volatility*100)
		report += fmt.Sprintf("Sharpe Ratio: %.3f\n", analytics.SharpeRatio)
		report += fmt.Sprintf("Maximum Drawdown: %.2f%%\n", analytics.MaxDrawdown*100)
		report += "\n"
	}
	
	// Volume statistics
	report += "=== VOLUME STATISTICS ===\n"
	report += fmt.Sprintf("Mean Volume: %.0f\n", analytics.VolumeStats.Mean)
	report += fmt.Sprintf("Median Volume: %.0f\n", analytics.VolumeStats.Median)
	report += fmt.Sprintf("Volume Range: %.0f - %.0f\n", analytics.VolumeStats.Min, analytics.VolumeStats.Max)
	report += fmt.Sprintf("Volume Std Dev: %.0f\n", analytics.VolumeStats.StdDev)
	report += "\n"
	
	// Technical indicators
	if len(analytics.RSI) > 0 {
		report += "=== TECHNICAL INDICATORS ===\n"
		latestRSI := analytics.RSI[len(analytics.RSI)-1]
		report += fmt.Sprintf("Latest RSI (14): %.2f", latestRSI)
		
		if latestRSI > 70 {
			report += " (Overbought)\n"
		} else if latestRSI < 30 {
			report += " (Oversold)\n"
		} else {
			report += " (Neutral)\n"
		}
	}
	
	if len(analytics.MACD.MACD) > 0 {
		latestMACD := analytics.MACD.MACD[len(analytics.MACD.MACD)-1]
		latestSignal := analytics.MACD.Signal[len(analytics.MACD.Signal)-1]
		report += fmt.Sprintf("Latest MACD: %.4f\n", latestMACD)
		report += fmt.Sprintf("MACD Signal: %.4f", latestSignal)
		
		if latestMACD > latestSignal {
			report += " (Bullish)\n"
		} else {
			report += " (Bearish)\n"
		}
	}
	
	if len(analytics.BollingerBands.Middle) > 0 {
		latest := len(analytics.BollingerBands.Middle) - 1
		latestPrice := timeseries.GetLatestPrice(bts).Close
		upper := analytics.BollingerBands.Upper[latest]
		middle := analytics.BollingerBands.Middle[latest]
		lower := analytics.BollingerBands.Lower[latest]
		
		report += fmt.Sprintf("Bollinger Bands - Upper: %.2f, Middle: %.2f, Lower: %.2f\n", upper, middle, lower)
		
		if latestPrice > upper {
			report += "Price is above upper band (potentially overbought)\n"
		} else if latestPrice < lower {
			report += "Price is below lower band (potentially oversold)\n"
		} else {
			report += "Price is within normal range\n"
		}
	}
	report += "\n"
	
	// Support and resistance
	if len(analytics.SupportResistance.SupportLevels) > 0 || len(analytics.SupportResistance.ResistanceLevels) > 0 {
		report += "=== SUPPORT & RESISTANCE LEVELS ===\n"
		
		if len(analytics.SupportResistance.SupportLevels) > 0 {
			report += "Support Levels: "
			for i, level := range analytics.SupportResistance.SupportLevels {
				if i > 0 {
					report += ", "
				}
				report += fmt.Sprintf("$%.2f", level)
			}
			report += "\n"
		}
		
		if len(analytics.SupportResistance.ResistanceLevels) > 0 {
			report += "Resistance Levels: "
			for i, level := range analytics.SupportResistance.ResistanceLevels {
				if i > 0 {
					report += ", "
				}
				report += fmt.Sprintf("$%.2f", level)
			}
			report += "\n"
		}
		report += "\n"
	}
	
	// Trend analysis
	trend := patterns.DetectTrend(bts, 30)
	report += "=== TREND ANALYSIS ===\n"
	report += fmt.Sprintf("30-Day Trend: %s\n", trend)
	
	// Pattern detection
	candlestickPatterns := patterns.DetectCandlestickPatterns(bts)
	volumePatterns := patterns.DetectVolumePatterns(bts)
	
	if len(candlestickPatterns) > 0 {
		report += "\n=== RECENT CANDLESTICK PATTERNS ===\n"
		for pattern, indices := range candlestickPatterns {
			if len(indices) > 0 {
				// Show only recent patterns (last 10 occurrences)
				recent := indices
				if len(indices) > 10 {
					recent = indices[len(indices)-10:]
				}
				report += fmt.Sprintf("%s: %d recent occurrences\n", pattern, len(recent))
			}
		}
	}
	
	if len(volumePatterns) > 0 {
		report += "\n=== RECENT VOLUME PATTERNS ===\n"
		for pattern, indices := range volumePatterns {
			if len(indices) > 0 {
				recent := indices
				if len(indices) > 5 {
					recent = indices[len(indices)-5:]
				}
				report += fmt.Sprintf("%s: %d recent occurrences\n", pattern, len(recent))
			}
		}
	}
	
	// Pivot points
	pivots := patterns.FindPivotPoints(bts)
	if len(pivots) > 0 {
		report += "\n=== PIVOT POINTS ===\n"
		if pivot, exists := pivots["pivot"]; exists {
			report += fmt.Sprintf("Pivot Point: $%.2f\n", pivot)
		}
		if r1, exists := pivots["r1"]; exists {
			report += fmt.Sprintf("Resistance 1: $%.2f\n", r1)
		}
		if s1, exists := pivots["s1"]; exists {
			report += fmt.Sprintf("Support 1: $%.2f\n", s1)
		}
	}
	
	// Fibonacci retracements
	fibs := patterns.CalculateFibonacciRetracements(bts, 30)
	if len(fibs) > 0 {
		report += "\n=== FIBONACCI RETRACEMENTS (30-day) ===\n"
		fibLevels := []string{"high", "fib_23_6", "fib_38_2", "fib_50", "fib_61_8", "fib_76_4", "low"}
		for _, level := range fibLevels {
			if price, exists := fibs[level]; exists {
				report += fmt.Sprintf("%s: $%.2f\n", level, price)
			}
		}
	}
	
	report += "\n=== END OF REPORT ===\n"
	report += fmt.Sprintf("Generated at: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	
	return report
}

// GetTradingSignals analyzes data and provides trading signals
func GetTradingSignals(bts *types.BTCTimeSeries, analytics types.BTCAnalytics) map[string]string {
	signals := make(map[string]string)
	
	// RSI signals
	if len(analytics.RSI) > 0 {
		latestRSI := analytics.RSI[len(analytics.RSI)-1]
		if latestRSI > 70 {
			signals["RSI"] = "SELL - Overbought"
		} else if latestRSI < 30 {
			signals["RSI"] = "BUY - Oversold"
		} else {
			signals["RSI"] = "HOLD - Neutral"
		}
	}
	
	// MACD signals
	if len(analytics.MACD.MACD) > 1 && len(analytics.MACD.Signal) > 1 {
		latestMACD := analytics.MACD.MACD[len(analytics.MACD.MACD)-1]
		prevMACD := analytics.MACD.MACD[len(analytics.MACD.MACD)-2]
		latestSignal := analytics.MACD.Signal[len(analytics.MACD.Signal)-1]
		prevSignal := analytics.MACD.Signal[len(analytics.MACD.Signal)-2]
		
		// Check for crossovers
		if prevMACD <= prevSignal && latestMACD > latestSignal {
			signals["MACD"] = "BUY - Bullish crossover"
		} else if prevMACD >= prevSignal && latestMACD < latestSignal {
			signals["MACD"] = "SELL - Bearish crossover"
		} else if latestMACD > latestSignal {
			signals["MACD"] = "HOLD - Bullish"
		} else {
			signals["MACD"] = "HOLD - Bearish"
		}
	}
	
	// Bollinger Bands signals
	if len(analytics.BollingerBands.Upper) > 0 {
		latestPrice := timeseries.GetLatestPrice(bts).Close
		latest := len(analytics.BollingerBands.Upper) - 1
		upper := analytics.BollingerBands.Upper[latest]
		lower := analytics.BollingerBands.Lower[latest]
		
		if latestPrice > upper {
			signals["Bollinger"] = "SELL - Price above upper band"
		} else if latestPrice < lower {
			signals["Bollinger"] = "BUY - Price below lower band"
		} else {
			signals["Bollinger"] = "HOLD - Price in normal range"
		}
	}
	
	// Trend signals
	trend := patterns.DetectTrend(bts, 30)
	switch trend {
	case "uptrend":
		signals["Trend"] = "BUY - Uptrend detected"
	case "downtrend":
		signals["Trend"] = "SELL - Downtrend detected"
	default:
		signals["Trend"] = "HOLD - Sideways movement"
	}
	
	// Support/Resistance signals
	if len(analytics.SupportResistance.SupportLevels) > 0 || len(analytics.SupportResistance.ResistanceLevels) > 0 {
		latestPrice := timeseries.GetLatestPrice(bts).Close
		
		// Check if price is near support (buy signal)
		for _, support := range analytics.SupportResistance.SupportLevels {
			if math.Abs(latestPrice-support)/support < 0.02 { // Within 2%
				signals["Support"] = "BUY - Near support level"
				break
			}
		}
		
		// Check if price is near resistance (sell signal)
		for _, resistance := range analytics.SupportResistance.ResistanceLevels {
			if math.Abs(latestPrice-resistance)/resistance < 0.02 { // Within 2%
				signals["Resistance"] = "SELL - Near resistance level"
				break
			}
		}
	}
	
	return signals
}

// CalculatePortfolioMetrics calculates portfolio-level metrics
func CalculatePortfolioMetrics(bts *types.BTCTimeSeries, initialInvestment float64) map[string]interface{} {
	metrics := make(map[string]interface{})
	
	if len(bts.Data) < 2 {
		return metrics
	}
	
	// Basic portfolio metrics
	backtest := statistics.PerformBacktest(bts, initialInvestment)
	for key, value := range backtest {
		metrics[key] = value
	}
	
	// Risk metrics
	riskMetrics := statistics.GetRiskMetrics(bts)
	for key, value := range riskMetrics {
		metrics[key] = value
	}
	
	// Performance ratios
	if volatility, exists := riskMetrics["volatility_annual"]; exists && volatility > 0 {
		if totalReturn, exists := backtest["annualized_return"]; exists {
			metrics["information_ratio"] = totalReturn / volatility
		}
	}
	
	return metrics
}