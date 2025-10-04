package main

import (
	"btc-analyzer/internal/analyzer"
	"btc-analyzer/internal/types"
	"btc-analyzer/internal/dataloader"
	"btc-analyzer/internal/reporter"
	"btc-analyzer/internal/visualizer"
	"encoding/base64"  // Move this to the top with other imports
	"flag"
	"fmt"
	"log"
	"os"
)

// generateSingleChart creates just the technical indicators chart
func generateSingleChart(bts *types.BTCTimeSeries, analytics types.BTCAnalytics, outputDir string) {
	fmt.Println("\n📊 Generating Technical Indicators Chart...")
	
	// Create charts directory
	chartsDir := fmt.Sprintf("%s/charts", outputDir)
	if err := os.MkdirAll(chartsDir, 0755); err != nil {
		fmt.Printf("Error creating charts directory: %v\n", err)
		return
	}
	
	// Generate just the technical indicators chart
	chartData, err := visualizer.GenerateIndicatorChart(bts, analytics)
	if err != nil {
		fmt.Printf("Error generating technical indicators chart: %v\n", err)
		return
	}
	
	// Save chart as PNG file
	chartPath := fmt.Sprintf("%s/technical_indicators.png", chartsDir)
	if err := os.WriteFile(chartPath, chartData, 0644); err != nil {
		fmt.Printf("Error saving chart: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Technical indicators chart saved: %s\n", chartPath)
	
	// Generate simple HTML report with just this chart
	htmlReport := generateSimpleHTMLReport(bts, analytics, chartData)
	htmlPath := fmt.Sprintf("%s/technical_analysis.html", outputDir)
	if err := os.WriteFile(htmlPath, []byte(htmlReport), 0644); err != nil {
		fmt.Printf("Error saving HTML report: %v\n", err)
	} else {
		fmt.Printf("✅ HTML report with chart: %s\n", htmlPath)
	}
	
	fmt.Println("📈 Technical indicators visualization complete!")
	fmt.Println("🌐 Open the HTML file in your browser to view the chart")
}

// generateSimpleHTMLReport creates a basic HTML report with the single chart
// generateSimpleHTMLReport creates a basic HTML report with the single chart and data tables
func generateSimpleHTMLReport(bts *types.BTCTimeSeries, analytics types.BTCAnalytics, chartData []byte) string {
	// Convert chart to base64
	base64Chart := ""
	if len(chartData) > 0 {
		base64Chart = base64.StdEncoding.EncodeToString(chartData)
	}
	
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Bitcoin Technical Indicators Analysis</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { 
            font-family: 'Segoe UI', Arial, sans-serif; 
            margin: 0; 
            padding: 20px; 
            background: #f5f5f5;
        }
        .container { 
            max-width: 1400px; 
            margin: 0 auto; 
            background: white; 
            padding: 30px; 
            border-radius: 10px; 
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .header { 
            text-align: center; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); 
            color: white; 
            padding: 30px; 
            border-radius: 10px; 
            margin-bottom: 30px;
        }
        .header h1 { margin: 0; font-size: 2.2em; }
        .stats-grid { 
            display: grid; 
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); 
            gap: 20px; 
            margin: 30px 0; 
        }
        .stat-card { 
            background: #f8f9fa; 
            padding: 20px; 
            border-radius: 8px; 
            text-align: center;
            border-left: 4px solid #667eea;
        }
        .stat-value { font-size: 1.8em; font-weight: bold; color: #333; }
        .stat-label { color: #666; margin-top: 5px; }
        .chart-container { 
            text-align: center; 
            margin: 30px 0; 
            padding: 20px; 
            background: #f8f9fa; 
            border-radius: 10px;
        }
        .chart-title { 
            font-size: 1.5em; 
            color: #333; 
            margin-bottom: 20px; 
        }
        img { 
            max-width: 100%; 
            height: auto; 
            border: 1px solid #ddd; 
            border-radius: 8px;
        }
        .data-section {
            margin: 30px 0;
            background: #f8f9fa;
            padding: 20px;
            border-radius: 10px;
        }
        .data-section h3 {
            color: #333;
            margin-top: 0;
        }
        .data-table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
            background: white;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        .data-table th,
        .data-table td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        .data-table th {
            background: #667eea;
            color: white;
            font-weight: 600;
        }
        .data-table tr:hover {
            background: #f5f5f5;
        }
        .data-table td.number {
            text-align: right;
            font-family: 'Courier New', monospace;
        }
        .data-table td.date {
            font-weight: 500;
        }
        .indicators { 
            background: #e3f2fd; 
            padding: 20px; 
            border-radius: 10px; 
            margin: 20px 0;
        }
        .indicators h3 { margin-top: 0; color: #1976d2; }
        .indicator-item { 
            display: inline-block; 
            margin: 10px 15px; 
            padding: 10px; 
            background: white; 
            border-radius: 5px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .summary-stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 15px;
            margin: 20px 0;
        }
        .summary-item {
            background: white;
            padding: 15px;
            border-radius: 8px;
            text-align: center;
            border-left: 3px solid #667eea;
        }
        .scrollable {
            max-height: 400px;
            overflow-y: auto;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📊 Bitcoin Technical Analysis</h1>
            <p>RSI & MACD Indicators with Raw Data</p>
        </div>

        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-value">` + fmt.Sprintf("%d", len(bts.Data)) + `</div>
                <div class="stat-label">Data Points</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">$` + fmt.Sprintf("%.2f", analytics.PriceStats.Mean) + `</div>
                <div class="stat-label">Average Price</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">` + fmt.Sprintf("%.2f%%", analytics.Volatility*100) + `</div>
                <div class="stat-label">Volatility</div>
            </div>`

	// Add current RSI if available
	if len(analytics.RSI) > 0 {
		currentRSI := analytics.RSI[len(analytics.RSI)-2]
		html += `
            <div class="stat-card">
                <div class="stat-value">` + fmt.Sprintf("%.1f", currentRSI) + `</div>
                <div class="stat-label">Current RSI</div>
            </div>`
	}

	html += `
        </div>`

	// Add chart if available
	if base64Chart != "" {
		html += `
        <div class="chart-container">
            <div class="chart-title">📈 Technical Indicators Chart</div>
            <img src="data:image/png;base64,` + base64Chart + `" alt="Technical Indicators Chart">
        </div>`
	}

	// Add Price Data Table
	html += `
        <div class="data-section">
            <h3>💰 Price Data (Last 20 Records)</h3>
            <div class="scrollable">
                <table class="data-table">
                    <thead>
                        <tr>
                            <th>Date</th>
                            <th>Open</th>
                            <th>High</th>
                            <th>Low</th>
                            <th>Close</th>
                            <th>Volume</th>
                        </tr>
                    </thead>
                    <tbody>`

	// Show last 20 price records
	start := len(bts.Data) - 20
	if start < 0 {
		start = 0
	}
	
	for i := start; i < len(bts.Data); i++ {
		data := bts.Data[i]
		html += `
                        <tr>
                            <td class="date">` + data.Timestamp.Format("Jan 02, 2006") + `</td>
                            <td class="number">$` + fmt.Sprintf("%.2f", data.Open) + `</td>
                            <td class="number">$` + fmt.Sprintf("%.2f", data.High) + `</td>
                            <td class="number">$` + fmt.Sprintf("%.2f", data.Low) + `</td>
                            <td class="number">$` + fmt.Sprintf("%.2f", data.Close) + `</td>
                            <td class="number">` + fmt.Sprintf("%.0f", data.Volume) + `</td>
                        </tr>`
	}

	html += `
                    </tbody>
                </table>
            </div>
        </div>`

	// Add RSI Data Table if available
	if len(analytics.RSI) > 0 {
		html += `
        <div class="data-section">
            <h3>📊 RSI Values (Last 20 Records)</h3>
            <div class="summary-stats">
                <div class="summary-item">
                    <strong>` + fmt.Sprintf("%.1f", analytics.RSI[len(analytics.RSI)-2]) + `</strong><br>
                    <small>Current RSI</small>
                </div>
                <div class="summary-item">
                    <strong>` + fmt.Sprintf("%d", len(analytics.RSI)) + `</strong><br>
                    <small>Total RSI Points</small>
                </div>`
		
		// Calculate RSI average
		rsiSum := 0.0
		for _, rsi := range analytics.RSI {
			rsiSum += rsi
		}
		rsiAvg := rsiSum / float64(len(analytics.RSI))
		
		html += `
                <div class="summary-item">
                    <strong>` + fmt.Sprintf("%.1f", rsiAvg) + `</strong><br>
                    <small>Average RSI</small>
                </div>
            </div>
            <div class="scrollable">
                <table class="data-table">
                    <thead>
                        <tr>
                            <th>Index</th>
                            <th>RSI Value</th>
                            <th>Status</th>
                        </tr>
                    </thead>
                    <tbody>`

		// Show last 20 RSI values
		rsiStart := len(analytics.RSI) - 20
		if rsiStart < 0 {
			rsiStart = 0
		}
		
		for i := rsiStart; i < len(analytics.RSI); i++ {
			rsi := analytics.RSI[i]
			status := "Neutral"
			if rsi < 30 {
				status = "Oversold"
			} else if rsi > 70 {
				status = "Overbought"
			}
			
			html += `
                        <tr>
                            <td class="number">` + fmt.Sprintf("%d", i+1) + `</td>
                            <td class="number">` + fmt.Sprintf("%.2f", rsi) + `</td>
                            <td>` + status + `</td>
                        </tr>`
		}

		html += `
                    </tbody>
                </table>
            </div>
        </div>`
	}

	// Add MACD Data Table if available
	if len(analytics.MACD.MACD) > 0 {
		html += `
        <div class="data-section">
            <h3>📈 MACD Values (Last 20 Records)</h3>
            <div class="summary-stats">
                <div class="summary-item">
                    <strong>` + fmt.Sprintf("%.3f", analytics.MACD.MACD[len(analytics.MACD.MACD)-1]) + `</strong><br>
                    <small>Current MACD</small>
                </div>`
		
		if len(analytics.MACD.Signal) > 0 {
			html += `
                <div class="summary-item">
                    <strong>` + fmt.Sprintf("%.3f", analytics.MACD.Signal[len(analytics.MACD.Signal)-1]) + `</strong><br>
                    <small>Current Signal</small>
                </div>`
		}
		
		html += `
                <div class="summary-item">
                    <strong>` + fmt.Sprintf("%d", len(analytics.MACD.MACD)) + `</strong><br>
                    <small>Total MACD Points</small>
                </div>
            </div>
            <div class="scrollable">
                <table class="data-table">
                    <thead>
                        <tr>
                            <th>Index</th>
                            <th>MACD</th>
                            <th>Signal</th>
                            <th>Histogram</th>
                            <th>Trend</th>
                        </tr>
                    </thead>
                    <tbody>`

		// Show last 20 MACD values
		macdStart := len(analytics.MACD.MACD) - 20
		if macdStart < 0 {
			macdStart = 0
		}
		
		for i := macdStart; i < len(analytics.MACD.MACD); i++ {
			macd := analytics.MACD.MACD[i]
			signal := ""
			histogram := ""
			trend := "Neutral"
			
			if i < len(analytics.MACD.Signal) {
				signalVal := analytics.MACD.Signal[i]
				signal = fmt.Sprintf("%.3f", signalVal)
				
				if macd > signalVal {
					trend = "Bullish"
				} else if macd < signalVal {
					trend = "Bearish"
				}
			}
			
			if i < len(analytics.MACD.Histogram) {
				histogram = fmt.Sprintf("%.3f", analytics.MACD.Histogram[i])
			}
			
			html += `
                        <tr>
                            <td class="number">` + fmt.Sprintf("%d", i+1) + `</td>
                            <td class="number">` + fmt.Sprintf("%.3f", macd) + `</td>
                            <td class="number">` + signal + `</td>
                            <td class="number">` + histogram + `</td>
                            <td>` + trend + `</td>
                        </tr>`
		}

		html += `
                    </tbody>
                </table>
            </div>
        </div>`
	}

	// Add indicator explanations
	html += `
        <div class="indicators">
            <h3>📋 Current Indicator Status</h3>`

	if len(analytics.RSI) > 0 {
		currentRSI := analytics.RSI[len(analytics.RSI)-1]
		rsiStatus := "Neutral"
		if currentRSI < 30 {
			rsiStatus = "Oversold (Buy Signal)"
		} else if currentRSI > 70 {
			rsiStatus = "Overbought (Sell Signal)"
		}
		html += `
            <div class="indicator-item">
                <strong>RSI (` + fmt.Sprintf("%.1f", currentRSI) + `):</strong> ` + rsiStatus + `
            </div>`
	}

	if len(analytics.MACD.MACD) > 0 && len(analytics.MACD.Signal) > 0 {
		currentMACD := analytics.MACD.MACD[len(analytics.MACD.MACD)-1]
		currentSignal := analytics.MACD.Signal[len(analytics.MACD.Signal)-1]
		macdStatus := "Neutral"
		if currentMACD > currentSignal {
			macdStatus = "Bullish Trend"
		} else if currentMACD < currentSignal {
			macdStatus = "Bearish Trend"
		}
		html += `
            <div class="indicator-item">
                <strong>MACD:</strong> ` + macdStatus + ` (` + fmt.Sprintf("%.3f", currentMACD) + `)
            </div>`
	}

	html += `
        </div>
    </div>
</body>
</html>`

	return html
}

func main() {
	// Command line flags
	var (
		source         = flag.String("source", "api", "Data source: 'api', 'csv', 'json', or 'sample'")
		days           = flag.Int("days", 30, "Number of days for API data")
		csvFile        = flag.String("csv", "", "CSV file path")
		jsonFile       = flag.String("json", "", "JSON file path")
		outputDir      = flag.String("output", ".", "Output directory for reports")
		htmlReport     = flag.Bool("html", true, "Generate HTML report")
		jsonReport     = flag.Bool("json-report", true, "Generate JSON report")
		generateChart  = flag.Bool("chart", true, "Generate technical indicators chart")
		verbose        = flag.Bool("verbose", false, "Verbose output")
	)
	flag.Parse()

	fmt.Println("🚀 Bitcoin Market Analyzer Starting...")

	// Load data based on source
	var bts *types.BTCTimeSeries
	var err error

	switch *source {
	case "api":
		fmt.Printf("📡 Fetching %d days of data from CoinGecko API...\n", *days)
		bts, err = dataloader.LoadFromCoinGecko(*days)
		if err != nil {
			log.Fatalf("Failed to load data from API: %v", err)
		}

	case "csv":
		if *csvFile == "" {
			log.Fatal("CSV file path required when using -source=csv")
		}
		fmt.Printf("📄 Loading data from CSV file: %s\n", *csvFile)
		bts, err = dataloader.LoadFromCSV(*csvFile)
		if err != nil {
			log.Fatalf("Failed to load CSV data: %v", err)
		}

	case "json":
		if *jsonFile == "" {
			log.Fatal("JSON file path required when using -source=json")
		}
		fmt.Printf("📄 Loading data from JSON file: %s\n", *jsonFile)
		bts, err = dataloader.LoadFromJSON(*jsonFile)
		if err != nil {
			log.Fatalf("Failed to load JSON data: %v", err)
		}

	case "sample":
		fmt.Println("🎲 Generating sample data for demonstration...")
		bts = dataloader.GenerateSampleData(*days, 50000.0)

	default:
		log.Fatalf("Invalid source: %s. Use 'api', 'csv', 'json', or 'sample'", *source)
	}

	if bts == nil {
		log.Fatal("Failed to load data")
	}

	// Validate data
	fmt.Println("🔍 Validating data...")
	issues := dataloader.ValidateData(bts)
	if len(issues) > 0 {
		fmt.Printf("⚠️  Data validation warnings:\n")
		for _, issue := range issues {
			fmt.Printf("  - %s\n", issue)
		}
	} else {
		fmt.Println("✅ Data validation passed")
	}

	// Perform analysis
	fmt.Println("📊 Performing comprehensive analysis...")
	analytics := analyzer.PerformComprehensiveAnalysis(bts)

	// Print summary to console
	reporter.PrintSummary(bts, analytics)

	// Generate technical indicators chart
	if *generateChart {
		generateSingleChart(bts, analytics, *outputDir)
	}

	// Generate reports
	if *htmlReport {
		htmlPath := fmt.Sprintf("%s/btc_analysis_report.html", *outputDir)
		fmt.Printf("📝 Generating HTML report: %s\n", htmlPath)
		if err := reporter.GenerateHTMLReport(bts, analytics, htmlPath); err != nil {
			log.Printf("Failed to generate HTML report: %v", err)
		} else {
			fmt.Printf("✅ HTML report generated successfully\n")
		}
	}

	if *jsonReport {
		jsonPath := fmt.Sprintf("%s/btc_analysis_report.json", *outputDir)
		fmt.Printf("📝 Generating JSON report: %s\n", jsonPath)
		if err := reporter.GenerateJSONReport(bts, analytics, jsonPath); err != nil {
			log.Printf("Failed to generate JSON report: %v", err)
		} else {
			fmt.Printf("✅ JSON report generated successfully\n")
		}
	}

	// Save processed data
	csvPath := fmt.Sprintf("%s/btc_data.csv", *outputDir)
	fmt.Printf("💾 Saving data to CSV: %s\n", csvPath)
	if err := dataloader.SaveToCSV(bts, csvPath); err != nil {
		log.Printf("Failed to save CSV: %v", err)
	}

	if *verbose {
		fmt.Println("\n" + analyzer.GenerateReport(bts, analytics))
	}

	fmt.Println("🎉 Analysis complete! Check the output directory for reports and charts.")
}