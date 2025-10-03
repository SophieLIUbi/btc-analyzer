package reporter

import (
	"btc-analyzer/internal/analyzer"
	"btc-analyzer/internal/types"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"time"
)

// GenerateHTMLReport creates an HTML report
func GenerateHTMLReport(bts *types.BTCTimeSeries, analytics types.BTCAnalytics, filename string) error {
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>Bitcoin Analysis Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .header { background-color: #f8f9fa; padding: 20px; border-radius: 5px; }
        .section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }
        .metric { display: inline-block; margin: 10px; padding: 10px; background-color: #e9ecef; border-radius: 3px; }
        .signal-buy { color: #28a745; font-weight: bold; }
        .signal-sell { color: #dc3545; font-weight: bold; }
        .signal-hold { color: #ffc107; font-weight: bold; }
        table { width: 100%; border-collapse: collapse; margin: 10px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Bitcoin Market Analysis Report</h1>
        <p>Symbol: {{.Symbol}} | Generated: {{.GeneratedAt}}</p>
        <p>Data Points: {{.DataPoints}} | Time Range: {{.TimeRange}}</p>
    </div>

    <div class="section">
        <h2>Current Price Information</h2>
        <div class="metric">Latest Price: ${{printf "%.2f" .LatestPrice}}</div>
        <div class="metric">Latest Volume: {{printf "%.0f" .LatestVolume}}</div>
    </div>

    <div class="section">
        <h2>Price Statistics</h2>
        <div class="metric">Mean: ${{printf "%.2f" .PriceStats.Mean}}</div>
        <div class="metric">Median: ${{printf "%.2f" .PriceStats.Median}}</div>
        <div class="metric">Min: ${{printf "%.2f" .PriceStats.Min}}</div>
        <div class="metric">Max: ${{printf "%.2f" .PriceStats.Max}}</div>
        <div class="metric">Std Dev: ${{printf "%.2f" .PriceStats.StdDev}}</div>
    </div>

    <div class="section">
        <h2>Risk Metrics</h2>
        <div class="metric">Volatility: {{printf "%.2f" .Volatility}}%</div>
        <div class="metric">Sharpe Ratio: {{printf "%.3f" .SharpeRatio}}</div>
        <div class="metric">Max Drawdown: {{printf "%.2f" .MaxDrawdown}}%</div>
    </div>

    {{if .Signals}}
    <div class="section">
        <h2>Trading Signals</h2>
        <table>
            <tr><th>Indicator</th><th>Signal</th></tr>
            {{range $indicator, $signal := .Signals}}
            <tr>
                <td>{{$indicator}}</td>
                <td class="{{if contains $signal "BUY"}}signal-buy{{else if contains $signal "SELL"}}signal-sell{{else}}signal-hold{{end}}">{{$signal}}</td>
            </tr>
            {{end}}
        </table>
    </div>
    {{end}}

    <div class="section">
        <h2>Technical Indicators</h2>
        {{if .LatestRSI}}
        <div class="metric">RSI (14): {{printf "%.2f" .LatestRSI}}</div>
        {{end}}
        {{if .LatestMACD}}
        <div class="metric">MACD: {{printf "%.4f" .LatestMACD}}</div>
        {{end}}
    </div>

    <div class="section">
        <h2>Full Text Report</h2>
        <pre>{{.TextReport}}</pre>
    </div>
</body>
</html>`

	// Prepare template data
	data := prepareTemplateData(bts, analytics)
	
	// Create template
	t, err := template.New("report").Funcs(template.FuncMap{
		"contains": func(s, substr string) bool {
			return fmt.Sprintf("%s", s) != fmt.Sprintf("%s", substr) // Simplified for template
		},
	}).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	
	// Create file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create HTML file: %w", err)
	}
	defer file.Close()
	
	// Execute template
	if err := t.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	
	return nil
}

// prepareTemplateData prepares data for HTML template
func prepareTemplateData(bts *types.BTCTimeSeries, analytics types.BTCAnalytics) map[string]interface{} {
	data := make(map[string]interface{})
	
	data["Symbol"] = bts.Symbol
	data["GeneratedAt"] = time.Now().Format("2006-01-02 15:04:05")
	data["DataPoints"] = len(bts.Data)
	
	if len(bts.Data) > 0 {
		latest := bts.Data[len(bts.Data)-1]
		data["LatestPrice"] = latest.Close
		data["LatestVolume"] = latest.Volume
		data["TimeRange"] = fmt.Sprintf("%s to %s", 
			bts.Data[0].Timestamp.Format("2006-01-02"),
			latest.Timestamp.Format("2006-01-02"))
	}
	
	data["PriceStats"] = analytics.PriceStats
	data["Volatility"] = analytics.Volatility * 100
	data["SharpeRatio"] = analytics.SharpeRatio
	data["MaxDrawdown"] = analytics.MaxDrawdown * 100
	
	if len(analytics.RSI) > 0 {
		data["LatestRSI"] = analytics.RSI[len(analytics.RSI)-1]
	}
	
	if len(analytics.MACD.MACD) > 0 {
		data["LatestMACD"] = analytics.MACD.MACD[len(analytics.MACD.MACD)-1]
	}
	
	// Get trading signals
	signals := analyzer.GetTradingSignals(bts, analytics)
	data["Signals"] = signals
	
	// Generate full text report
	data["TextReport"] = analyzer.GenerateReport(bts, analytics)
	
	return data
}

// GenerateJSONReport creates a JSON report
func GenerateJSONReport(bts *types.BTCTimeSeries, analytics types.BTCAnalytics, filename string) error {
	report := map[string]interface{}{
		"metadata": map[string]interface{}{
			"symbol":        bts.Symbol,
			"generated_at":  time.Now().Format(time.RFC3339),
			"data_points":   len(bts.Data),
		},
		"analytics":     analytics,
		"trading_signals": analyzer.GetTradingSignals(bts, analytics),
		"portfolio_metrics": analyzer.CalculatePortfolioMetrics(bts, 10000), // $10k initial
	}
	
	if len(bts.Data) > 0 {
		latest := bts.Data[len(bts.Data)-1]
		report["metadata"].(map[string]interface{})["latest_price"] = latest.Close
		report["metadata"].(map[string]interface{})["latest_volume"] = latest.Volume
		report["metadata"].(map[string]interface{})["time_range"] = map[string]string{
			"start": bts.Data[0].Timestamp.Format("2006-01-02"),
			"end":   latest.Timestamp.Format("2006-01-02"),
		}
	}
	
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create JSON report file: %w", err)
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	if err := encoder.Encode(report); err != nil {
		return fmt.Errorf("failed to encode JSON report: %w", err)
	}
	
	return nil
}

// PrintSummary prints a brief summary to console
func PrintSummary(bts *types.BTCTimeSeries, analytics types.BTCAnalytics) {
	fmt.Println("=== BITCOIN ANALYSIS SUMMARY ===")
	
	if len(bts.Data) > 0 {
		latest := bts.Data[len(bts.Data)-1]
		fmt.Printf("Latest Price: $%.2f\n", latest.Close)
		fmt.Printf("Data Points: %d\n", len(bts.Data))
	}
	
	fmt.Printf("Mean Price: $%.2f\n", analytics.PriceStats.Mean)
	fmt.Printf("Price Range: $%.2f - $%.2f\n", analytics.PriceStats.Min, analytics.PriceStats.Max)
	
	if analytics.Volatility > 0 {
		fmt.Printf("Volatility: %.2f%%\n", analytics.Volatility*100)
		fmt.Printf("Sharpe Ratio: %.3f\n", analytics.SharpeRatio)
	}
	
	if len(analytics.RSI) > 0 {
		fmt.Printf("Latest RSI: %.2f\n", analytics.RSI[len(analytics.RSI)-1])
	}
	
	// Show key signals
	signals := analyzer.GetTradingSignals(bts, analytics)
	fmt.Println("\n=== KEY SIGNALS ===")
	for indicator, signal := range signals {
		fmt.Printf("%s: %s\n", indicator, signal)
	}
	
	fmt.Println("================================")
}